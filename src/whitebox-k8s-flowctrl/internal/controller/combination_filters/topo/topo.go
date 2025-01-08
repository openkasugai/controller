/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package top

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	"github.com/go-logr/logr"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type GraphManager struct {
	// Graph of all nodes and links *Not used in the current scheduling
	g *simple.UndirectedGraph

	// For internal connection route search
	// A graph that contains only the links within each node (i.e. does not include links used for external connections)
	nodeNameToGraph map[string]*simple.UndirectedGraph

	// EntityId of memory in each node for each node name. Not supported when there are multiple memory entities in one node.
	nodeNameToEntityIdOfMemory map[string]string

	// EntityId of the host NIC in that node for each node name. Not supported when one node has multiple host NIC entities.
	nodeNameToEntityIdOfHostNic map[string]string

	// For searching external routes
	// For cases where the link's From/To is FPGA. A weighted graph that increases the weight (cost) of internal links so that external paths are searched with priority.
	wg *simple.WeightedUndirectedGraph

	// For searching external routes
	// For when the link's From/To is GPU/CPU. A weighted graph in which the weight (cost) of internal links is increased so that external routes are searched with priority. The FPGA NIC link is deleted.
	wgExcludeFpgaNic *simple.WeightedUndirectedGraph

	// A map for converting topology information EntityIds to Graph NodeIds (int64)
	entityIdToNodeId map[string]int64

	// A map for converting Graph NodeIds (int64) to topology information EntityIds
	nodeIdToEntityId map[int64]string

	// The key is the Graph's NodeId (int64) and the value is the topology information's EntityInfo map.
	// Temporarily saved information on device I/F and network usage status inside and outside the node that is updated during scheduling
	nodeIdToEntityInfo map[int64]v1.EntityInfo

	// EntityId of the external server where the FC's wb-strt-of-chain and wb-end-of-chain are located
	entityIdOfExternalNode string
}

var l logr.Logger

// Constructor function for the structure
func NewGraphManager() *GraphManager {
	return &GraphManager{}
}

func (t *GraphManager) CreateGraphs(ctx context.Context, ti *v1.TopologyInfo) {
	l = log.FromContext(ctx)
	l.Info("start CreateGraphs", "TopologyinfoNamespace", ti.Namespace, "TopologyinfoName", ti.Name)

	t.entityIdToNodeId = map[string]int64{}
	t.nodeIdToEntityInfo = map[int64]v1.EntityInfo{}
	t.nodeNameToEntityIdOfMemory = map[string]string{}
	t.nodeNameToEntityIdOfHostNic = map[string]string{}
	fpgaNicList := []string{}
	for _, e := range ti.Status.Entities {
		nId := int64(len(t.entityIdToNodeId))
		t.entityIdToNodeId[e.ID] = nId
		t.nodeIdToEntityInfo[nId] = e
		if e.Type == "interface" && e.InterfaceInfo.DeviceType != nil && strings.Contains(*e.InterfaceInfo.DeviceType, "alveo") &&
			strings.Contains(e.InterfaceInfo.InterfaceType, "ether") {
			fpgaNicList = append(fpgaNicList, e.ID)
		} else if e.Type == "node" {
			t.entityIdOfExternalNode = e.ID
		} else if e.Type == "device" && e.DeviceInfo.DeviceType == "memory" {
			t.nodeNameToEntityIdOfMemory[e.DeviceInfo.NodeName] = e.ID
		} else if e.Type == "device" && e.DeviceInfo.DeviceType == "nic" {
			t.nodeNameToEntityIdOfHostNic[e.DeviceInfo.NodeName] = e.ID
		}
	}

	t.nodeIdToEntityId = map[int64]string{}
	for eId, nId := range t.entityIdToNodeId {
		t.nodeIdToEntityId[nId] = eId
	}

	// Graph of all nodes and all links
	t.g = simple.NewUndirectedGraph()

	// Graph of intranode links for each node (for searching for internal connections only)
	t.nodeNameToGraph = map[string]*simple.UndirectedGraph{}

	for _, r := range ti.Status.Relations {
		if !r.Available {
			continue
		}
		t.g.SetEdge(simple.Edge{F: simple.Node(t.entityIdToNodeId[r.From]), T: simple.Node(t.entityIdToNodeId[r.To])})

		if r.Type == "direct" || r.Type == "pcie" || r.Type == "qpi" {
			fromEntity := t.nodeIdToEntityInfo[t.entityIdToNodeId[r.From]]
			toEntity := t.nodeIdToEntityInfo[t.entityIdToNodeId[r.To]]
			fromNodeName := t.getNodeNameInEntity(fromEntity)
			toNodeName := t.getNodeNameInEntity(toEntity)
			if fromNodeName != toNodeName {
				l.Error(errors.New("graph error"), "nodeName does not match in from and to of inner relation")
				continue
			}
			g, ok := t.nodeNameToGraph[fromNodeName]
			if !ok {
				g = simple.NewUndirectedGraph()
				t.nodeNameToGraph[fromNodeName] = g
			}
			g.SetEdge(simple.Edge{F: simple.Node(t.entityIdToNodeId[r.From]), T: simple.Node(t.entityIdToNodeId[r.To])})
		}

	}

	t.wg = simple.NewWeightedUndirectedGraph(0, math.Inf(1))
	for _, r := range ti.Status.Relations {
		if !r.Available {
			continue
		}
		var weight float64
		if r.Type == "pcie" || r.Type == "qpi" {
			weight = 1000
		} else if r.Type == "direct" {
			fromEntity := t.nodeIdToEntityInfo[t.entityIdToNodeId[r.From]]
			toEntity := t.nodeIdToEntityInfo[t.entityIdToNodeId[r.To]]
			if (fromEntity.Type == "interface" && fromEntity.InterfaceInfo.InterfaceSideType != "outside") ||
				(toEntity.Type == "interface" && toEntity.InterfaceInfo.InterfaceSideType != "outside") {
				// Links where the I/F passing through the internal connection is from or to
				weight = 1000
			} else if fromEntity.Type == "interface" && fromEntity.InterfaceInfo.InterfaceSideType == "outside" &&
				fromEntity.InterfaceInfo.DeviceType != nil && strings.Contains(*fromEntity.InterfaceInfo.DeviceType, "alveo") ||
				toEntity.Type == "interface" && toEntity.InterfaceInfo.InterfaceSideType == "outside" &&
					toEntity.InterfaceInfo.DeviceType != nil && strings.Contains(*toEntity.InterfaceInfo.DeviceType, "alveo") {
				// If the link is an external I/F and the I/F is an FPGA NIC, give it a higher weight (cost) than the link that goes through the host NIC.
				weight = 10
			}
		} else {
			// For external links, the weight (cost) is small.
			weight = 1
		}
		t.wg.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(t.entityIdToNodeId[r.From]), T: simple.Node(t.entityIdToNodeId[r.To]), W: weight})
	}

	// When searching for external connections only and the From is GPU or CPU (paths via FPGA NIC are deleted)
	t.wgExcludeFpgaNic = simple.NewWeightedUndirectedGraph(0, math.Inf(1))
	graph.CopyWeighted(t.wgExcludeFpgaNic, t.wg)
	for _, eId := range fpgaNicList {
		l.V(2).Info("remove from wgCopy: " + eId)
		t.wgExcludeFpgaNic.RemoveNode(t.entityIdToNodeId[eId])
	}

	l.Info("end CreateGraphs", "TopologyinfoNamespace", ti.Namespace, "TopologyinfoName", ti.Name)
}

func (t *GraphManager) GetEntityId(ft *v1.FunctionTarget) string {
	eId := ft.Status.NodeName + "." + ft.Status.DeviceType + "-" + strconv.Itoa(int(ft.Status.DeviceIndex))
	_, ok := t.entityIdToNodeId[eId]
	if ok {
		return eId
	}
	eId2 := eId + "." + ft.Status.RegionName
	_, ok = t.entityIdToNodeId[eId2]
	if ok {
		return eId2
	} else {
		return ""
	}
}

func (t *GraphManager) GetEntityIdOfExternalNode() string {
	return t.entityIdOfExternalNode
}

func (t *GraphManager) GetInterfaceTypeAndNodeName(entityId string) (string, string) {
	eInfo, ok := t.nodeIdToEntityInfo[t.entityIdToNodeId[entityId]]
	if ok && eInfo.Type == "interface" {
		return eInfo.InterfaceInfo.InterfaceType, eInfo.InterfaceInfo.NodeName
	} else {
		return "", ""
	}
}

// Processing to determine whether the incoming, outgoing, both (IncomingAndOutgoing), or none
func (t *GraphManager) GetIOUsedType(routeStr []string) []v1.WBConnectionPath {
	wbcRoute := []v1.WBConnectionPath{}
	var usedType v1.WBIOUsedType
	change := false
	for i, eId := range routeStr {
		eInfo := t.nodeIdToEntityInfo[t.entityIdToNodeId[eId]]
		if eInfo.Type == "interface" {
			wbcRoute = append(wbcRoute, v1.WBConnectionPath{EntityID: eId, UsedType: usedType})
		} else if eInfo.Type == "device" {
			wbcRoute = append(wbcRoute, v1.WBConnectionPath{EntityID: eId, UsedType: v1.WBIOUsedTypeNone})
			change = true
		} else if eInfo.Type == "network" {
			wbcRoute = append(wbcRoute, v1.WBConnectionPath{EntityID: eId, UsedType: v1.WBIOUsedTypeIncomingAndOutgoing})
			change = true
		} else {
			wbcRoute = append(wbcRoute, v1.WBConnectionPath{EntityID: eId, UsedType: v1.WBIOUsedTypeNone})
		}
		if i == 0 {
			usedType = v1.WBIOUsedTypeOutgoing
			change = false
			continue
		}
		if change {
			if usedType == v1.WBIOUsedTypeIncoming {
				usedType = v1.WBIOUsedTypeOutgoing
			} else if usedType == v1.WBIOUsedTypeOutgoing {
				usedType = v1.WBIOUsedTypeIncoming
			}
			change = false
		}
	}
	return wbcRoute
}

func (t *GraphManager) getNodeNameInEntity(entity v1.EntityInfo) string {
	entityKind := entity.Type
	switch entityKind {
	case "device":
		return entity.DeviceInfo.NodeName
	case "interface":
		return entity.InterfaceInfo.NodeName
	case "network":
		return *entity.NetworkInfo.NodeName
	case "node":
		return entity.NodeInfo.NodeName
	default:
		return ""
	}
}

func (t *GraphManager) FindOuterShortestPaths(src string, dst string, k int) [][]string {
	pths := path.YenKShortestPaths(t.wg, k, simple.Node(t.entityIdToNodeId[src]), simple.Node(t.entityIdToNodeId[dst]))
	routeArry := make([][]string, k)
	for i, nodes := range pths {
		for _, node := range nodes {
			routeArry[i] = append(routeArry[i], t.nodeIdToEntityId[node.ID()])
		}
	}
	return routeArry
}

func (t *GraphManager) FindOuterExcludeFpgaNicShortestPaths(src string, dst string, k int) [][]string {
	pths := path.YenKShortestPaths(t.wgExcludeFpgaNic, k, simple.Node(t.entityIdToNodeId[src]), simple.Node(t.entityIdToNodeId[dst]))
	routeList := make([][]string, k)
	for i, nodes := range pths {
		for _, node := range nodes {
			routeList[i] = append(routeList[i], t.nodeIdToEntityId[node.ID()])
		}
	}
	return routeList
}

func (t *GraphManager) FindInnerShortestPaths(nodeName string, src string, dst string, k int) [][]string {
	pths := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], k, simple.Node(t.entityIdToNodeId[src]), simple.Node(t.entityIdToNodeId[dst]))
	routeList := make([][]string, k)
	for i, nodes := range pths {
		for _, node := range nodes {
			routeList[i] = append(routeList[i], t.nodeIdToEntityId[node.ID()])
		}
	}
	return routeList
}

func (t *GraphManager) FindReturnPathAtHostNic(nodeName string, src string, dst string) [][]string {
	routeList := make([][]string, 1)
	eIdOfHostNic, ok := t.nodeNameToEntityIdOfHostNic[nodeName]
	if !ok {
		return routeList
	}

	pth1 := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], 1,
		simple.Node(t.entityIdToNodeId[src]), simple.Node(t.entityIdToNodeId[eIdOfHostNic]))
	pth2 := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], 1,
		simple.Node(t.entityIdToNodeId[dst]), simple.Node(t.entityIdToNodeId[eIdOfHostNic]))
	if len(pth1) != 0 && len(pth2) != 0 {
		for _, n := range pth1[0] {
			routeList[0] = append(routeList[0], t.nodeIdToEntityId[n.ID()])
		}
		// Add from the end of pth2. Do not add the last element (host NIC) as it is duplicated.
		for i := len(pth2[0]) - 2; i >= 0; i-- {
			n := pth2[0][i]
			routeList[0] = append(routeList[0], t.nodeIdToEntityId[n.ID()])
		}
		return routeList
	} else {
		return routeList
	}
}

func (t *GraphManager) FindReturnPathAtMemory(nodeName string, src string, dst string) [][]string {
	routeList := make([][]string, 1)
	eIdOfMemory, ok := t.nodeNameToEntityIdOfMemory[nodeName]
	if !ok {
		return routeList
	}

	pth1 := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], 1,
		simple.Node(t.entityIdToNodeId[src]), simple.Node(t.entityIdToNodeId[eIdOfMemory]))
	pth2 := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], 1,
		simple.Node(t.entityIdToNodeId[dst]), simple.Node(t.entityIdToNodeId[eIdOfMemory]))
	if len(pth1) != 0 && len(pth2) != 0 {
		for _, n := range pth1[0] {
			routeList[0] = append(routeList[0], t.nodeIdToEntityId[n.ID()])
		}
		// Add from the end of pth2. The last element (Memory) is duplicated so it is not added.
		for i := len(pth2[0]) - 2; i >= 0; i-- {
			n := pth2[0][i]
			routeList[0] = append(routeList[0], t.nodeIdToEntityId[n.ID()])
		}
		return routeList
	} else {
		return routeList
	}
}

func (t *GraphManager) HasRouteToMemory(nodeName string, src string, dst string) bool {
	eIdOfMemory, ok := t.nodeNameToEntityIdOfMemory[nodeName]
	if !ok {
		return false
	}

	pth1 := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], 1, simple.Node(t.entityIdToNodeId[src]), simple.Node(t.entityIdToNodeId[eIdOfMemory]))
	pth2 := path.YenKShortestPaths(t.nodeNameToGraph[nodeName], 1, simple.Node(t.entityIdToNodeId[dst]), simple.Node(t.entityIdToNodeId[eIdOfMemory]))
	if len(pth1) != 0 && len(pth2) != 0 {
		return true
	} else {
		return false
	}
}

// Check if it is an external route
func (t *GraphManager) IsOuterRoute(route []string, isFromFpga, isToFpga, isRoundTrip, isSameNode bool) bool {
	for _, eId := range route {
		if isFromFpga && isToFpga {
			// Check if it is passing through the Ether switch ("global.ether-network-N")
			if strings.HasPrefix(eId, "global.ether-network") {
				return true
			}
		} else {
			if isRoundTrip || isSameNode {
				// Check if it is passing through the host's NIC
				if strings.Contains(eId, "nic") {
					return true
				}
			} else {
				// Check if it is passing through the Ether switch ("global.ether-network-N")
				if strings.HasPrefix(eId, "global.ether-network") {
					return true
				}
			}
		}
	}
	return false
}

// Check if it is an external route
func (t *GraphManager) IsOuterWBConnRoute(wbcRoute []v1.WBConnectionPath) bool {
	for _, wbcr := range wbcRoute {
		// Check if it is passing through the Ether switch ("global.ether-network-N")
		if strings.HasPrefix(wbcr.EntityID, "global.ether-network") {
			return true
			// Check if it is passing through the host's NIC
		} else if strings.Contains(wbcr.EntityID, "nic") {
			return true
		}
	}
	return false
}

func (t *GraphManager) GetEntityInfo(entityId string) v1.EntityInfo {
	return t.nodeIdToEntityInfo[t.entityIdToNodeId[entityId]]
}
