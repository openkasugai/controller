package combination_filters

import (
	"context"
	"fmt"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/filter_template"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/scheduler_common"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *CombinationFilters) targetResourceFitFilter(
	ctx context.Context,
	req ctrl.Request,
	sd *v1.SchedulingData,
	df *v1.DataFlow,
) (ctrl.Result, error) {

	l := log.FromContext(ctx)

	// Get TargetCombinations from SchedulingData
	combs := sd.Status.TargetCombinations

	// Get FunctionTarget(FT)
	fts, err := r.FetchFunctionTargets(ctx, req, sd, df, nil)
	if err != nil {
		return r.Abort(ctx, sd, err)
	}

	// Create map[FunctionTarget name]*FunctionTarget from []*FunctionTarget
	ftMap := make(map[string]*v1.FunctionTarget)
	for _, ft := range fts {
		ftMap[ft.Name] = ft
	}

	// Get the requested capacity
	cap := r.getRequireCapacity(df)

	excludeIndexes := make([]int32, 0, len(combs))
	for comb_i, comb := range combs {

		// Calculate the amount of resources used by the combination of deployment candidates
		// Get the FunctionTarget name and FunctionIndex of the deployment candidate from comb, and the FunctionInstance that uses it.
		// Calculate how many resources are required for the {FunctionTarget, FunctionIndex} over the entire combination
		ftCurrfuncs := make(map[*v1.FunctionTarget]int32)
		ftCurrCap := make(map[*v1.FunctionTarget]int32)
		fiCurrDfs := make(map[FunctionIndexStruct]int32)
		fiCurrCap := make(map[FunctionIndexStruct]int32)
		for _, v := range comb.ScheduledFunctions {
			ft := ftMap[getFunctionTargetName(v)]
			fi := v.FunctionIndex
			ftCurrCap[ft] += cap
			if fi != nil {
				key := FunctionIndexStruct{FunctionTarget: ft, FunctionIndex: fi}
				fiCurrDfs[key]++
				fiCurrCap[key] += cap
			} else {
				ftCurrfuncs[ft]++
			}
		}

		// Check if the FunctionTarget has enough resources
		isValid := true
		for _, v := range comb.ScheduledFunctions {
			ft := ftMap[getFunctionTargetName(v)]
			fi := v.FunctionIndex

			// If FunctionTarget.Status.Status is "NotReady", that FT is target of child bit automatic writing
			// Skip capacity check because each Capacity value is nil
			if ft.Status.Status == v1.WBRegionStatusNotReady {
				continue
				// Normaly this shuld not happen, but just in case, handle it as a fallback
			} else if ft.Status.Status != v1.WBRegionStatusReady {
				isValid = false
				break
			}

			if *ft.Status.CurrentFunctions+ftCurrfuncs[ft] > *ft.Status.MaxFunctions {
				l.Info("function target MaxFunctions capacity over. FunctionTarget=" + ft.Name +
					" MaxFunctions=" + fmt.Sprintf("%d", *ft.Status.MaxFunctions) +
					" postDeployFunctions=" + fmt.Sprintf("%d", *ft.Status.CurrentFunctions+ftCurrfuncs[ft]))
				isValid = false
				break
			} else if *ft.Status.CurrentCapacity+ftCurrCap[ft] > *ft.Status.MaxCapacity {
				l.Info("function target MaxCapacity capacity over. FunctionTarget=" + ft.Name +
					" MaxCapacity=" + fmt.Sprintf("%d", *ft.Status.MaxCapacity) +
					" postDeployCapacity=" + fmt.Sprintf("%d", *ft.Status.CurrentCapacity+ftCurrCap[ft]))
				isValid = false
				break
			}

			if fi != nil {
				key := FunctionIndexStruct{FunctionTarget: ft, FunctionIndex: fi}
				f := SearchFunction(*fi, &ft.Status.Functions)
				if *f.CurrentDataFlows+fiCurrDfs[key] > *f.MaxDataFlows {
					l.Info("function MaxDataflows capacity over. FunctionTarget=" + ft.Name +
						" FunctionIndex=" + fmt.Sprintf("%d", *fi) +
						" MaxDataflows=" + fmt.Sprintf("%d", *f.MaxDataFlows) +
						" postDeployDataflows=" + fmt.Sprintf("%d", *f.CurrentDataFlows+fiCurrDfs[key]))
					isValid = false
					break
				} else if *f.CurrentCapacity+fiCurrCap[key] > *f.MaxCapacity {
					l.Info("function MaxCapacity capacity over. FunctionTarget=" + ft.Name +
						" FunctionIndex=" + fmt.Sprintf("%d", *fi) +
						" MaxCapacity=" + fmt.Sprintf("%d", *f.MaxCapacity) +
						" postDeployCapacity=" + fmt.Sprintf("%d", *f.CurrentCapacity+fiCurrCap[key]))
					isValid = false
					break
				}
			}

		}

		// If not enough, register it as a combination to play.
		if !isValid {
			excludeIndexes = append(excludeIndexes, int32(comb_i))
		}

	}

	// Eliminate combinations that do not satisfy the conditions
	exclude(&sd.Status.TargetCombinations, excludeIndexes)

	return r.Finalize(ctx, sd)
}
