package combination_filters

import (
	"context"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *CombinationFilters) connectionResourceFitFilter(
	ctx context.Context,
	req ctrl.Request,
	sd *v1.SchedulingData,
	df *v1.DataFlow,
) (ctrl.Result, error) {

	l := log.FromContext(ctx)

	l.Info("ABC")

	// Get TargetCombinations from SchedulingData
	combs := sd.Status.TargetCombinations

	// Get TopologyInfo
	ti, err := r.getTopologyInfo(ctx, sd)
	if err != nil {
		return r.Abort(ctx, sd, err)
	}

	// Create a map for checking usage
	entityMap := make(map[string]*v1.EntityInfo)
	for _, e := range ti.Status.Entities {
		// Save a copy because data inconsistency may occur
		entityCopy := e
		entityMap[e.ID] = &entityCopy
	}

	// Get the requested capacity
	cap := r.getRequireCapacity(df)

	excludeIndexes := make([]int32, 0, len(combs))
	for comb_i, comb := range combs {

		incomings := make(map[string]int32)
		outgoings := make(map[string]int32)

		for _, csi := range comb.ScheduledConnections {
			// Calculate the amount of resources for the routes used by the combination of deployment candidates
			for _, p := range csi.ConnectionPath {
				switch p.UsedType {
				case v1.WBIOUsedTypeIncomingAndOutgoing:
					incomings[p.EntityID] += cap
					outgoings[p.EntityID] += cap
				case v1.WBIOUsedTypeIncoming:
					incomings[p.EntityID] += cap
				case v1.WBIOUsedTypeOutgoing:
					outgoings[p.EntityID] += cap
				}
			}
		}

		for _, csi := range comb.ScheduledConnections {
			// Verify that the resources of the entities that make up the route are sufficient.
			isValid := true
			for _, p := range csi.ConnectionPath {
				entity := entityMap[p.EntityID]
				capInfo := entity.CapacityInfo
				switch entity.Type {
				case "network":
					if capInfo.CurrentIncomingCapacity+incomings[p.EntityID] > capInfo.MaxIncomingCapacity {
						isValid = false
						l.Info("network incoming capacity over. EntityId=" + p.EntityID + " Route=" + r.getRouteStr(csi.ConnectionPath))
						break
					} else if capInfo.CurrentOutgoingCapacity+outgoings[p.EntityID] > capInfo.MaxOutgoingCapacity {
						l.Info("network outgoing capacity over. EntityId=" + p.EntityID + " Route=" + r.getRouteStr(csi.ConnectionPath))
						break
					}
				case "interface":
					switch p.UsedType {
					case v1.WBIOUsedTypeIncoming:
						if capInfo.CurrentIncomingCapacity+incomings[p.EntityID] > capInfo.MaxIncomingCapacity {
							isValid = false
							l.Info("device interface incoming capacity over. EntityId=" + p.EntityID + " Route=" + r.getRouteStr(csi.ConnectionPath))
							break
						}
					case v1.WBIOUsedTypeOutgoing:
						if capInfo.CurrentOutgoingCapacity+outgoings[p.EntityID] > capInfo.MaxOutgoingCapacity {
							isValid = false
							l.Info("device interface outgoing capacity over. EntityId=" + p.EntityID + " Route=" + r.getRouteStr(csi.ConnectionPath))
							break
						}
					}
				}
			}

			if !isValid {
				excludeIndexes = append(excludeIndexes, int32(comb_i))
				break
			}

		}

	}

	// Eliminate combinations that do not satisfy the conditions
	exclude(&sd.Status.TargetCombinations, excludeIndexes)

	return r.Finalize(ctx, sd)
}
