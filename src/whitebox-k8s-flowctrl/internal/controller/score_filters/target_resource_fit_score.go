package score_filters

import (
	"context"
	"fmt"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/filter_template"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/scheduler_common"
	ctrl "sigs.k8s.io/controller-runtime"
)

const targetResourceFitWeightsKey = "Weights"

const funcInstFieldName = "FuncKey"

type targetResourceFitWeight struct {
	FunctionTargetScore    int64
	FunctionIndexScore     int64
	FunctionTargetCapacity int64
	FunctionIndexCapacity  int64
}

func (r *ScoreFilters) targetResourceFitScore(
	ctx context.Context,
	req ctrl.Request,
	sd *v1.SchedulingData,
	df *v1.DataFlow,
) error {

	// Get TargetCombinations from SchedulingData
	combs := *&sd.Status.TargetCombinations //nolint:staticcheck // SA4001: *&x will be simplified to x. FIXME: fix if no side effecnts exist.

	// Get FunctionTarget(FT)
	fts, err := r.FetchFunctionTargets(ctx, req, sd, df, nil)
	if err != nil {
		return err
	}

	// Create map[FunctionTarget name]*FunctionTarget from []*FunctionTarget
	ftMap := make(map[string]*v1.FunctionTarget)
	for _, ft := range fts {
		ftMap[ft.Name] = ft
	}

	// Get the weight for each score calculation from the strategy
	// The items are Function, FunctionIndex, FunctionTargetCapacity, FunctionIndexCapacity under Filter.Weight
	// If not present, use default value
	weights := r.getTargetResourceFitWeights(sd)

	// Get the requested capacity
	cap := int64(r.getRequireCapacity(df))

	// Get the sorted FunctionInstance names (because the Score of FJ changes depending on the order)
	sortedFuncInstNames := GetFieldValuesFromStructs[string](
		GetSortedFunctionInstanceOrder(df.Status.FunctionChain),
		funcInstFieldName,
	)

	calcFtScore := func(ft *v1.FunctionTarget, weight int64) int64 {
		// Returns 0 if Status.Status of FunctionTarget is "NotReady"
		if ft.Status.Status == v1.WBRegionStatusNotReady {
			return int64(0)
			// Normaly this shuld not happen, but just in case, handle it as a fallback
		} else if ft.Status.Status != v1.WBRegionStatusReady {
			return int64(0)
		} else {
			// In cases other than the above, calculate FTScore
			return int64(*ft.Status.CurrentCapacity)*weight + int64(*ft.Status.MaxCapacity)
		}
	}

	// Calculate the points at ft, {ft, fi} alone
	ftScores := make(map[*v1.FunctionTarget]int64)
	fiScores := make(map[FunctionIndexStruct]int64)
	for _, comb := range combs {
		for _, fName := range sortedFuncInstNames {

			fsi := comb.ScheduledFunctions[fName]
			ft := ftMap[getFunctionTargetName(fsi)]
			fi := fsi.FunctionIndex

			// Normaly this shuld not happen, but just in case, handle it as a fallback
			if !(ft.Status.Status == v1.WBRegionStatusReady || ft.Status.Status == v1.WBRegionStatusNotReady) {
				return fmt.Errorf("functionTarget %s is not eligible for scoring", ft.Name)
			}

			if _, ok := ftScores[ft]; !ok {
				ftScores[ft] = calcFtScore(ft, weights.FunctionTargetCapacity)
			}

			if fi != nil {
				key := FunctionIndexStruct{FunctionTarget: ft, FunctionIndex: fi}
				if _, ok := fiScores[key]; !ok {
					function := SearchFunction(*fi, &ft.Status.Functions)
					fiScores[key] =
						int64(*function.CurrentCapacity)*weights.FunctionIndexCapacity + int64(*function.MaxCapacity)
				}
			}
		}
	}

	// Calculate nodeScore
	nodeScores := make(map[string]int64)
	for _, ft := range fts {
		nodeName := ft.Status.NodeName
		if v, ok := ftScores[ft]; ok {
			nodeScores[nodeName] += v
		} else {
			nodeScores[nodeName] += calcFtScore(ft, weights.FunctionTargetCapacity)
		}
	}

	// Calculate the score for the combination
	for comb_i, comb := range combs {
		nodeDiffs := make(map[string]int64)
		ftDiffs := make(map[*v1.FunctionTarget]int64)
		fiDiffs := make(map[FunctionIndexStruct]int64)
		for _, fName := range sortedFuncInstNames {
			fsi := comb.ScheduledFunctions[fName]
			ft := ftMap[getFunctionTargetName(fsi)]
			fi := fsi.FunctionIndex
			fiKey := FunctionIndexStruct{FunctionTarget: ft, FunctionIndex: fi}

			// Calculate each score
			ftScore := ftScores[ft] + ftDiffs[ft]

			var fiScore int64 = 0
			if fi != nil {
				fiScore = fiScores[fiKey] + fiDiffs[fiKey]
			}

			nodeScore := nodeScores[ft.Status.NodeName] + nodeDiffs[ft.Status.NodeName]

			// Add score
			var initialScore int64 = 0
			if combs[comb_i].Score == nil {
				combs[comb_i].Score = &initialScore
			}

			*combs[comb_i].Score +=
				(fiScore*weights.FunctionIndexScore+ftScore)*weights.FunctionTargetScore + nodeScore

			// Update the score based on usage
			ftDiff := cap * weights.FunctionTargetCapacity
			ftDiffs[ft] += ftDiff
			fiDiffs[fiKey] += cap * weights.FunctionIndexCapacity
			nodeDiffs[ft.Status.NodeName] += ftDiff

		}
	}

	return nil
}

func (r *ScoreFilters) getTargetResourceFitWeights(sd *v1.SchedulingData) targetResourceFitWeight {

	ret := targetResourceFitWeight{
		FunctionTargetScore:    1000 * 1000,
		FunctionIndexScore:     1000 * 1000,
		FunctionTargetCapacity: 1000,
		FunctionIndexCapacity:  1000,
	}

	if r.DoesStrategyHaveKey(sd, targetResourceFitWeightsKey) {
		var rcv map[string]int64
		r.LoadStrategyParameter(sd, targetResourceFitWeightsKey, &rcv)
		if v, ok := rcv["FunctionTarget"]; ok {
			ret.FunctionTargetScore = v
		}
		if v, ok := rcv["FunctionIndex"]; ok {
			ret.FunctionIndexScore = v
		}
		if v, ok := rcv["FunctionTargetCapacity"]; ok {
			ret.FunctionTargetCapacity = v
		}
		if v, ok := rcv["FunctionIndexCapacity"]; ok {
			ret.FunctionTargetCapacity = v
		}
	}

	return ret
}
