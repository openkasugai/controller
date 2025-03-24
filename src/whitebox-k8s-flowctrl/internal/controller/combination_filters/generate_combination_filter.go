package combination_filters

import (
	"context"
	"fmt"
	"sort"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/filter_template"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/scheduler_common"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	RequestDeviceIndexKey     = "requestDeviceIndex"
	RequestRegionNamesKey     = "requestRegionNames"
	RequestFunctionIndexesKey = "requestFunctionIndexes"
)

const funcInstFieldName = "FuncKey"

func (r *CombinationFilters) generateCombinationsFilter(
	ctx context.Context,
	req ctrl.Request,
	sd *v1.SchedulingData,
	df *v1.DataFlow,
) (ctrl.Result, error) {

	l := log.FromContext(ctx)

	// Get the sorted functionInstance names
	funcInstNames := GetFieldValuesFromStructs[string](
		GetSortedFunctionInstanceOrder(df.Status.FunctionChain),
		funcInstFieldName,
	)

	// Get the requested capacity
	cap := r.getRequireCapacity(df)

	// Get the FunctionType.
	fTyps := df.Status.FunctionType

	// Create a map[FunctionType name]FunctionType
	fTypMap := make(map[string]*v1.FunctionType)
	for _, fTyp := range fTyps {
		fTypMap[fTyp.Spec.FunctionName] = fTyp
	}

	// Get userInput from userRequirement
	fTgtFilters := r.getUserInput(sd, funcInstNames)

	combPerFuncInst := make([][]FunctionIndexStruct, len(funcInstNames))
	for funcInstIdx, funcInstName := range funcInstNames {

		// Get the FunctionType corresponding to the FunctionInstance for TargetKindFit, AddReusableFunction
		var fTyp *v1.FunctionType
		for key, function := range df.Status.FunctionChain.Spec.Functions {

			if key == funcInstName {
				fTyp = fTypMap[function.FunctionName]
				break
			}
		}

		// The equivalent of TargetKindFit is in FilterTemplate
		// The following process is abolished because it overwrites the information set from UserInput with the information from FunctionType.
		// fTgtFilters[funcInstIdx].RegionTypes = &fTyp.Status.RegionTypeCandidates

		// Processing equivalent to FTUserInput
		// Get the FunctionTarget for each FunctionInstance (only those that are available and follow the userInput value)
		fTgts, err := r.FetchFunctionTargets(ctx, req, sd, df, &fTgtFilters[funcInstIdx])
		if err != nil {
			l.Error(err, "Failed to fetch FunctionTargets")
			return r.Abort(ctx, sd, err)
		}

		// Processing equivalent to TargetKindFit
		// If FunctionTarget.Status.RegionType matches the target regionTypeCandidates,
		// Extract that FunctionTarget as a candidate
		r.getRegionTypeFitFunctionTargets(ctx, &fTgts, fTyp.Status.RegionTypeCandidates)

		// Processing equivalent to FIUserInput
		// Here, the conversion from []FunctionTarget to []{FunctionTarget, FunctionIndex} takes place.
		functionIndexStructs := r.FetchFunctionIndexStructs(ctx, &fTgts, &fTgtFilters[funcInstIdx])

		// Processing corresponding to FunctionName
		// If functionIndexStructs[i].FunctionTarget.Status.Functions has the same name as the target FunctionName,
		// Extract the Function of that FunctionIndex as a candidate
		r.getFunctionNameFitFunctionIndexStructs(&functionIndexStructs, fTyp.Spec.FunctionName)

		// Processing equivalent to TargetResourceFit
		// When the target function is deployed, exclude combinations that exceed resources in {functionTarget, functionIndex}.
		r.filterTargetResourceFit(ctx, &functionIndexStructs, &cap)

		combPerFuncInst[funcInstIdx] = functionIndexStructs
	}

	// Sort by name (match output to FJ's)
	for comb_i := range combPerFuncInst {
		sort.SliceStable(combPerFuncInst[comb_i],
			func(i, j int) bool {
				return combPerFuncInst[comb_i][i].FunctionTarget.Name < combPerFuncInst[comb_i][j].FunctionTarget.Name
			})
	}

	// Create combinations by using the Cartesian product of candidates for each function
	functionIndexStructCombs := generateCombinations(combPerFuncInst)

	// Create a FunctionScheduleInfo from a combination of FunctionIndexStructs,
	// Create a TargetCombination
	combs := make([]v1.TargetCombinationStruct, 0)
	for _, fisComb := range functionIndexStructCombs {

		fsiMap := make(map[string]v1.FunctionScheduleInfo)
		for fInst_i, funcInstName := range funcInstNames {
			fsiMap[funcInstName] = createScheduledFunctionsInfo(fisComb[fInst_i].FunctionTarget, fisComb[fInst_i].FunctionIndex)
		}

		// In the case of FunctionTarget.Status.Status is v1.WBRegionStatusNotReady
		// Candidates with multiple functions are deployed on the same device are excluded
		if hasDuplicateFpgaDevice(fisComb, fsiMap, funcInstNames) {
			continue
		}

		combs = append(combs, v1.TargetCombinationStruct{
			ScheduledFunctions: fsiMap,
		})
	}

	// Write to SchedulingData
	sd.Status.TargetCombinations = combs

	return r.Finalize(ctx, sd)
}

// Create the combination
// Here, we will use the Filter in the FJ version, which can be narrowed down by FunctionInstance:FunctionTarget, and does not require a combination.
func (r *CombinationFilters) getUserInput(sd *v1.SchedulingData, funcInstNames []string) []FunctionTargetFilter {

	ret := make([]FunctionTargetFilter, len(funcInstNames))

	var requestNodeNames map[string][]string
	if r.DoesUserRequirementHaveKey(sd, RequestNodeNamesKey) {
		r.LoadUserRequirementParameter(sd, RequestNodeNamesKey, &requestNodeNames)
	}

	var requestDeviceTypes map[string][]string
	if r.DoesUserRequirementHaveKey(sd, RequestDeviceTypesKey) {
		r.LoadUserRequirementParameter(sd, RequestDeviceTypesKey, &requestDeviceTypes)
	}

	// The FJ version has a DeviceIndex specification,
	// DeviceIndex is not specified alone, but is set with DeviceType.
	// If both are specified, a unique FunctionTarget can be selected.
	// Abolish DeviceIndex and integrate it into FunctionTarget specification.
	var requestFunctionTargets map[string][]string
	if r.DoesUserRequirementHaveKey(sd, RequestFunctionTargetKey) {
		r.LoadUserRequirementParameter(sd, RequestFunctionTargetKey, &requestFunctionTargets)
	}

	var requestRegionNames map[string][]string
	if r.DoesUserRequirementHaveKey(sd, RequestRegionNamesKey) {
		r.LoadUserRequirementParameter(sd, RequestRegionNamesKey, &requestRegionNames)
	}

	// Added when FJ integrated with NTT Scheduler
	var requestFunctionIndexes map[string][]string
	if r.DoesUserRequirementHaveKey(sd, RequestFunctionIndexesKey) {
		r.LoadUserRequirementParameter(sd, RequestFunctionIndexesKey, &requestFunctionIndexes)
	}

	for i, funcInstName := range funcInstNames {

		fil := FunctionTargetFilter{IncludesNotAvailable: valToAddr[bool](false)}

		if v, ok := requestDeviceTypes[funcInstName]; ok {
			fil.DeviceTypes = &v
		}

		if v, ok := requestNodeNames[funcInstName]; ok {
			fil.NodeNames = &v
		}

		if v, ok := requestFunctionTargets[funcInstName]; ok {
			fil.FunctionTargets = &v
		}

		if v, ok := requestRegionNames[funcInstName]; ok {
			fil.RegionNames = &v
		}

		// Added when FJ integrated with NTT Scheduler
		if v, ok := requestFunctionIndexes[funcInstName]; ok {
			fil.FunctionIndexes = &v
		}

		ret[i] = fil
	}

	return ret

}

func (r *CombinationFilters) getRegionTypeFitFunctionTargets(
	ctx context.Context,
	fts *[]*v1.FunctionTarget,
	targetRegionTypeCandidates []string,
) {

	l := log.FromContext(ctx)

	retFTs := make([]*v1.FunctionTarget, 0)

	// If ft.Status.RegionType has the same name as the RegionType in the target regionTypeCandidates,
	// Extract that FunctionTarget as a candidate.
	found := false
	for _, ft := range *fts {
		for _, targetRegionType := range targetRegionTypeCandidates {
			if ft.Status.RegionType == targetRegionType {
				retFTs = append(retFTs, ft)
				found = true
			}
		}
	}

	if !found {
		l.Info("no FunctionTarget found for RegionTypeCandidates:" +
			fmt.Sprintf("%v", targetRegionTypeCandidates))
	}

	*fts = retFTs
}

func (r *CombinationFilters) getFunctionNameFitFunctionIndexStructs(
	cands *[]FunctionIndexStruct,
	targetFuncName string,
) {

	retCands := make([]FunctionIndexStruct, 0)

	// If fis.ft.Status.Functions.FunctionName has the same name as the target FunctionName,
	// Extract the FunctionIndex of that Function as a candidate
	for _, fis := range *cands {
		if fis.FunctionIndex == nil {
			retCands = append(retCands, fis)
		} else {
			for _, function := range fis.FunctionTarget.Status.Functions {
				if *fis.FunctionIndex == function.FunctionIndex && function.FunctionName == targetFuncName {
					retCands = append(retCands, fis)
				}
			}
		}
	}

	*cands = retCands
}

// func (r *CombinationFilters) addReusableFunction(fts []*v1.FunctionTarget, funcName string) []functionIndexStruct {

// 	ret := make([]functionIndexStruct, 0)
// 	for _, ft := range fts {
// // Add {ft, nil} as a candidate
// 		ret = append(ret, functionIndexStruct{ft, nil})

// // If there is a function with the same name as the target Function in ft.Status.Functions and it is Available, the one with that FunctionIndex will also be added as one of the candidates.
// 		for i := range ft.Status.Functions {
// 			function := ft.Status.Functions[i]
// 			if function.FunctionName == funcName {
// 				ret = append(ret, functionIndexStruct{ft, &function.FunctionIndex})
// 			}
// 		}
// 	}

// 	return ret
// }

// Check the resource amount of functionTarget
func (r *CombinationFilters) filterTargetResourceFit(ctx context.Context, candidatesAdr *[]FunctionIndexStruct, requireCapacity *int32) {

	l := log.FromContext(ctx)

	candidates := *candidatesAdr

	excludeIndexes := make([]int32, int32(0), len(candidates))
	for i, candidate := range candidates {

		ft := candidate.FunctionTarget
		fi := candidate.FunctionIndex

		// If FunctionTarget.Status.Status is "NotReady", that FT is target of child bit automatic writing
		// Skip capacity check because each Capacity value is nil
		if ft.Status.Status == v1.WBRegionStatusNotReady {
			continue
			// Normaly this shuld not happen, but just in case, handle it as a fallback
		} else if ft.Status.Status != v1.WBRegionStatusReady {
			excludeIndexes = append(excludeIndexes, int32(i))
			continue
		}

		if fi == nil {
			if *ft.Status.MaxFunctions == *ft.Status.CurrentFunctions {
				excludeIndexes = append(excludeIndexes, int32(i))
				l.Info("function target MaxFunctions is full. FunctionTarget=" + ft.Name)
			} else if *ft.Status.MaxCapacity < *ft.Status.CurrentCapacity+*requireCapacity {
				excludeIndexes = append(excludeIndexes, int32(i))
				l.Info("function target MaxCapacity will inevitably result in capacity over. FunctionTarget=" + ft.Name)
			}
		} else {

			// Get the Function corresponding to FI
			var function v1.FunctionCapStruct
			for _, f := range ft.Status.Functions {
				if f.FunctionIndex == *fi {
					function = f
					break
				}
			}

			// - FT.Status.MaxCapacity < FT.Status.CurrentCapacity + requireCapacity
			// - FT.Status.Function[FunctionIndex].MaxDataFlows == FT.Status.Function[FunctionIndex].CurrentDataFlows
			// - FT.Status.Function[FunctionIndex].MaxCapacity < FT.Status.Function[FunctionIndex].CurrentCapacity + requireCapacity
			// Register Index as a candidate for exclusion
			if *ft.Status.MaxCapacity < *ft.Status.CurrentCapacity+*requireCapacity {
				excludeIndexes = append(excludeIndexes, int32(i))
				l.Info("function target MaxCapacity will inevitably result in capacity over. FunctionTarget=" + ft.Name + " FunctionIndex=" + fmt.Sprintf("%d", *fi))
			} else if *function.MaxDataFlows == *function.CurrentDataFlows {
				excludeIndexes = append(excludeIndexes, int32(i))
				l.Info("function MaxDataFlows is full. FunctionTarget=" + ft.Name + " FunctionIndex=" + fmt.Sprintf("%d", *fi))
			} else if *function.MaxCapacity < *function.CurrentCapacity+*requireCapacity {
				excludeIndexes = append(excludeIndexes, int32(i))
				l.Info("function MaxCapacity will inevitably result in capacity over. FunctionTarget=" + ft.Name + " FunctionIndex=" + fmt.Sprintf("%d", *fi))
			}
		}
	}

	// Exclude
	exclude(candidatesAdr, excludeIndexes)
}

func generateCombinations[T any](in [][]T) [][]T {
	indexLengths := make([]int, len(in))
	for i, elms := range in {
		indexLengths[i] = len(elms)
	}
	indexesComb := generateIndexComb(0, indexLengths, []int{})
	ret := make([][]T, len(indexesComb))
	for i, indexes := range indexesComb {
		ret[i] = make([]T, len(indexes))
		for j, index := range indexes {
			ret[i][j] = in[j][index]
		}
	}
	return ret
}

// Create an index combination
func generateIndexComb(curIdx int, indexLengths []int, comb []int) [][]int {
	if curIdx == len(indexLengths) {
		newComb := make([]int, len(comb))
		copy(newComb, comb)
		return [][]int{newComb}
	}
	ret := make([][]int, 0)
	for i := 0; i < indexLengths[curIdx]; i++ {
		// Create and pass a new slice for each recursive call to prevent the slice contents from changing unexpectedly.
		newComb := append(append([]int{}, comb...), i)
		ret = append(ret, generateIndexComb(curIdx+1, indexLengths, newComb)...)
	}
	return ret
}

func createScheduledFunctionsInfo(functionTarget *v1.FunctionTarget, functionIndex *int32) v1.FunctionScheduleInfo {
	scheduledFunctions := v1.FunctionScheduleInfo{
		NodeName:      functionTarget.Status.NodeName,
		DeviceType:    functionTarget.Status.DeviceType,
		DeviceIndex:   functionTarget.Status.DeviceIndex,
		RegionName:    functionTarget.Status.RegionName,
		FunctionIndex: functionIndex,
	}
	return scheduledFunctions
}

func valToAddr[T any](in T) *T {
	return &in
}

func hasDuplicateFpgaDevice(fisComb []FunctionIndexStruct, fsiMap map[string]v1.FunctionScheduleInfo, funcInstNames []string) bool {

	var previousFtStatus, currentFtStatus v1.WBRegionStatus
	var previousFtKey, currentFtKey string

	for fInst_i, functfuncInstName := range funcInstNames {
		scheduledFuncInfo := fsiMap[functfuncInstName]

		currentFtKey = fmt.Sprintf("%s-%s-%d",
			scheduledFuncInfo.NodeName,
			scheduledFuncInfo.DeviceType,
			scheduledFuncInfo.DeviceIndex)

		currentFtStatus = fisComb[fInst_i].FunctionTarget.Status.Status

		// Exclude if the previous or current FunctionTarget.Status.Status is v1.WBRegionStatusNotReady and the device is duplicated
		if (previousFtStatus == v1.WBRegionStatusNotReady || currentFtStatus == v1.WBRegionStatusNotReady) && previousFtKey == currentFtKey {
			return true
		}

		previousFtStatus = currentFtStatus
		previousFtKey = currentFtKey
	}
	return false
}
