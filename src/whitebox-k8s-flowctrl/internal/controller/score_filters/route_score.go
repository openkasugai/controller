package score_filters

import (
	"context"
	"strings"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ScoreFilters) routeScore(
	ctx context.Context,
	req ctrl.Request,
	sd *v1.SchedulingData,
	df *v1.DataFlow,
) error {

	// Get TopologyInfo
	ti, err := r.getTopologyInfo(ctx, sd)
	if err != nil {
		return err
	}

	// Create EntityInfo map[EntityId]
	entityMap := make(map[string]v1.EntityInfo)
	for _, e := range ti.Status.Entities {
		entityMap[e.ID] = e
	}

	// Calculate and add score
	for i, comb := range sd.Status.TargetCombinations {
		*sd.Status.TargetCombinations[i].Score += r.calcRouteScore(comb.ScheduledConnections, entityMap)
	}

	return nil
}

// Calculate the Score based on TargetCombinations.ConnectionScheduleInfo.Route
func (r *ScoreFilters) calcRouteScore(connectionScheduleInfos []v1.ConnectionScheduleInfo, entityMap map[string]v1.EntityInfo) int64 {

	var score int64 = 0
	for _, csi := range connectionScheduleInfos {

		route := csi.ConnectionPath

		// Score for internally connected routes
		routeTypeScore := 1000
		innerRoute := true

		for _, wbcRoute := range route {
			if strings.HasPrefix(wbcRoute.EntityID, "global.ether-network") {
				// If the route is an external route, the score is 10.
				routeTypeScore = 10
				innerRoute = false
				break
			} else if strings.Contains(wbcRoute.EntityID, "nic") {
				// If the route is an external connection within a node, the score is 100
				routeTypeScore = 100
				innerRoute = false
			}
		}

		numaNodesScore := 0
		if innerRoute {
			tmpPcieNwEntityId := ""

			for _, wbcRoute := range route {

				// Get entityInfo from entityId
				eInfo, ok := entityMap[wbcRoute.EntityID]
				isPcieNetwork := ok && eInfo.Type == "network" && eInfo.NetworkInfo != nil && eInfo.NetworkInfo.NetworkType == "pcie"

				if isPcieNetwork {
					if tmpPcieNwEntityId == "" {
						tmpPcieNwEntityId = wbcRoute.EntityID
					} else if tmpPcieNwEntityId != "" && tmpPcieNwEntityId != wbcRoute.EntityID {
						numaNodesScore = -100
						break
					}
				}

			}

		}

		score += int64(routeTypeScore + numaNodesScore)

	}
	return score
}
