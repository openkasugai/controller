package combination_filters

import (
	"context"
	"encoding/json"
	"strings"

	"golang.org/x/exp/slices"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/combination_filters/topo"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/scheduler_common"
)

const k = 5

const (
	inputInterfaceTypeKey  = "inputInterfaceType"
	outputInterfaceTypeKey = "outputInterfaceType"
)

const (
	innerConnKeyWord = "mem"
	outerConnKeyWord = "ether"
)

type deployableItem struct {
	Name                string `json:"name"`
	RegionType          string `json:"regionType"`
	InputInterfaceType  string `json:"inputInterfaceType"`
	OutputInterfaceType string `json:"outputInterfaceType"`
}

type connectionDetailInfo struct {
	fromFuncKey        string
	toFuncKey          string
	connType           string
	fromPort           int32
	toPort             int32
	fromInterfaceTypes map[string]string
	toInterfaceTypes   map[string]string
	fromRegionType     string
	toRegionType       string
}

type routeRelatedInfo struct {
	fromFuncKey    string
	toFuncKey      string
	fromRegionType string
	toRegionType   string
}

func (r *CombinationFilters) generateRouteFilter(
	ctx context.Context,
	req ctrl.Request,
	sd *v1.SchedulingData,
	df *v1.DataFlow,
) (ctrl.Result, error) {

	// Get TargetCombinations from SchedulingData
	combs := sd.Status.TargetCombinations

	// Create a Graph from TopologyInfo to search for the shortest path
	// Get TopologyInfo
	ti, err := r.getTopologyInfo(ctx, sd)
	if err != nil {
		return r.Abort(ctx, sd, err)
	}
	// I don't use FJ's WbTopology, but I think the CreateGraph part can be used as is.
	// (There are two types of Graphs: one with FPGA NIC and one without)
	graph := NewGraphManager()
	graph.CreateGraphs(ctx, ti)

	// Create a map[<Function name>][<RegionType name>][]<Corresponding Interface name> to be used to verify connection availability.
	// Get FunctionInfo(ConfigMap). name : fk.Spec.FunctionName, namespace : fk.Spec.FunctionInfoNameSpace
	// (Isn't it no longer a ConfigMap? For now, I'll follow the FJ one.)
	// FunctionInfo is included in map [<RegionType name>][]<interface name>
	interfacesPerRegionTypesPerFunctions, _ := r.createUsableInterfaceTypeList(ctx, df.Status.FunctionType) // FIXME: need to handle 2nd return value (error)

	// Create a map[ funcInstFrom ] ConnectionType
	// 		ConnectionType ⊂ {inner, outer, both}
	// conTypes := setConTypes(df.Status.FunctionChain, interfacesPerRegionTypesPerFunctions)
	connDtlInfoMap := createConnectionDetailInfoMap(df.Status.FunctionChain, interfacesPerRegionTypesPerFunctions)

	// Get sorted FunctionInstance information
	sortedFuncInstances := GetSortedFunctionInstanceOrder(df.Status.FunctionChain)

	// Get the correspondence between functionInstName and functionName
	funcInstNameToFuncNameMap := make(map[string]string, len(sortedFuncInstances))
	for _, funcInstance := range sortedFuncInstances {
		for key, function := range df.Status.FunctionChain.Spec.Functions {
			if funcInstance.FuncKey == key {
				funcInstNameToFuncNameMap[funcInstance.FuncKey] = function.FunctionName
				break
			}
		}
	}

	// Create map[FunctionTarget name]*FunctionTarget
	// Get FunctionTarget(FT)
	fts, err := r.FetchFunctionTargets(ctx, req, sd, df, nil)
	if err != nil {
		return r.Abort(ctx, sd, err)
	}
	// create
	ftMap := make(map[string]*v1.FunctionTarget)
	for _, ft := range fts {
		// Save a copy because data inconsistency may occur
		ftCopy := *ft
		ftMap[ft.Name] = &ftCopy
	}

	resCombs := make([]v1.TargetCombinationStruct, 0, len(sd.Status.TargetCombinations)*k)
	for _, comb := range combs {
		// Create a route for each combination and return it with ConnectionScheduleInfo attached.
		// Since multiple paths are generated (aligned with FJ), multiple TargetCombinations are returned.
		res := r.addScheduledConnectionsInfo(
			ctx, comb, ftMap, connDtlInfoMap,
			sortedFuncInstances, funcInstNameToFuncNameMap, interfacesPerRegionTypesPerFunctions, graph)
		resCombs = append(resCombs, res...)
	}

	// Update the SchedulingData combination
	sd.Status.TargetCombinations = resCombs

	return r.Finalize(ctx, sd)
}

func (r *CombinationFilters) createUsableInterfaceTypeList(ctx context.Context, functionTypes []*v1.FunctionType) (map[string]map[string]map[string][]string, error) {

	l := log.FromContext(ctx)

	funcNameToInterfaceTypeMap := make(map[string]map[string]map[string][]string)

	for _, ft := range functionTypes {
		var functionInfo corev1.ConfigMap
		if err := r.Get(ctx, types.NamespacedName{Namespace: ft.Spec.FunctionInfoCMRef.Namespace, Name: ft.Spec.FunctionInfoCMRef.Name}, &functionInfo); err != nil {
			l.Error(err, "unable to fetch functioninfo (configmap). funcName="+ft.Spec.FunctionName)
			return nil, err
		}

		var deployableItems []deployableItem
		if err := json.Unmarshal([]byte(functionInfo.Data["deployableItems"]), &deployableItems); err != nil {
			l.Error(err, "failed to unmarshal deployableItems")
			return nil, err
		}

		for _, item := range deployableItems {
			interfaceTypeMap := funcNameToInterfaceTypeMap[ft.Spec.FunctionName]
			if interfaceTypeMap == nil {
				interfaceTypeMap = make(map[string]map[string][]string)
				funcNameToInterfaceTypeMap[ft.Spec.FunctionName] = interfaceTypeMap
			}
			regionTypeMap := interfaceTypeMap[item.RegionType]
			if regionTypeMap == nil {
				regionTypeMap = make(map[string][]string)
				regionTypeMap[inputInterfaceTypeKey] = []string{}
				regionTypeMap[outputInterfaceTypeKey] = []string{}
				interfaceTypeMap[item.RegionType] = regionTypeMap
			}

			if !slices.Contains(regionTypeMap[inputInterfaceTypeKey], item.InputInterfaceType) {
				regionTypeMap[inputInterfaceTypeKey] = append(regionTypeMap[inputInterfaceTypeKey], item.InputInterfaceType)
			}
			if !slices.Contains(regionTypeMap[outputInterfaceTypeKey], item.OutputInterfaceType) {
				regionTypeMap[outputInterfaceTypeKey] = append(regionTypeMap[outputInterfaceTypeKey], item.OutputInterfaceType)
			}
		}
	}

	return funcNameToInterfaceTypeMap, nil
}

// create possible connection type and related information based on the available InterfaceType list
func createConnectionDetailInfoMap(fc *v1.FunctionChain,

	funcNameToInterfaceTypeMap map[string]map[string]map[string][]string) map[string]connectionDetailInfo {

	funcInstNametoInterfaceTypeMap := map[string]map[string]map[string][]string{}
	connectionDetailInfoMap := make(map[string]connectionDetailInfo)

	for funcKey, funcInfo := range fc.Spec.Functions {
		interfaceTypeInfo := funcNameToInterfaceTypeMap[funcInfo.FunctionName]
		funcInstNametoInterfaceTypeMap[funcKey] = interfaceTypeInfo
	}

	for _, connection := range fc.Spec.Connections {

		// add startPoint and endPoint available InterfaceType to funcInstNametoInterfaceTypeMap
		addStartEndInterfaceTypeInfo(&connection, funcInstNametoInterfaceTypeMap)

		// create possible connection type and other information
		for fromRegionType, fromRegionData := range funcInstNametoInterfaceTypeMap[connection.From.FunctionKey] {
			for toRegionType, toRegionData := range funcInstNametoInterfaceTypeMap[connection.To.FunctionKey] {

				fromInterfaceTypes, toInterfaceTypes := fromRegionData[outputInterfaceTypeKey], toRegionData[inputInterfaceTypeKey]

				fromInterfaceTypeMap, toInterfaceTypeMap := createInterfaceTypeMap(fromInterfaceTypes), createInterfaceTypeMap(toInterfaceTypes)

				availableConnType := getAvailableConnType(fromInterfaceTypeMap, toInterfaceTypeMap)
				if availableConnType == "" {
					continue
				}

				connectionDetailInfo := connectionDetailInfo{
					fromFuncKey: connection.From.FunctionKey, toFuncKey: connection.To.FunctionKey,
					connType: availableConnType,
					fromPort: connection.From.Port, toPort: connection.To.Port,
					fromInterfaceTypes: fromInterfaceTypeMap, toInterfaceTypes: toInterfaceTypeMap,
					fromRegionType: fromRegionType, toRegionType: toRegionType,
				}

				key := connectionDetailInfo.fromFuncKey + "-" +
					connectionDetailInfo.fromRegionType + "-" +
					connectionDetailInfo.toFuncKey + "-" +
					connectionDetailInfo.toRegionType
				connectionDetailInfoMap[key] = connectionDetailInfo
			}
		}
	}

	return connectionDetailInfoMap
}

func addStartEndInterfaceTypeInfo(connection *v1.ConnectionStruct, funcInstNametoInterfaceTypeMap map[string]map[string]map[string][]string) {
	if strings.HasPrefix(connection.From.FunctionKey, StartPointKeyword) {
		if _, exists := funcInstNametoInterfaceTypeMap[connection.From.FunctionKey]; !exists {
			funcInstNametoInterfaceTypeMap[connection.From.FunctionKey] = map[string]map[string][]string{
				StartPointKeyword: {
					outputInterfaceTypeKey: GetStartEndPointAvailableInterfaceType(StartPointKeyword),
				},
			}
		}
	} else if strings.HasPrefix(connection.To.FunctionKey, EndPointKeyword) {
		if _, exists := funcInstNametoInterfaceTypeMap[connection.To.FunctionKey]; !exists {
			funcInstNametoInterfaceTypeMap[connection.To.FunctionKey] = map[string]map[string][]string{
				EndPointKeyword: {
					inputInterfaceTypeKey: GetStartEndPointAvailableInterfaceType(EndPointKeyword),
				},
			}
		}
	}
}

func createInterfaceTypeMap(interfaceTypes []string) map[string]string {
	interfaceTypeMap := make(map[string]string)
	for _, interfaceType := range interfaceTypes {
		switch {
		case strings.Contains(interfaceType, innerConnKeyWord):
			interfaceTypeMap["inner"] = interfaceType
		case strings.Contains(interfaceType, outerConnKeyWord):
			interfaceTypeMap["outer"] = interfaceType
		}
	}
	return interfaceTypeMap
}

func getAvailableConnType(fromInterfaceTypeMap, toInterfaceTypeMap map[string]string) string {

	_, fromIsInner := fromInterfaceTypeMap["inner"]
	_, fromIsOuter := fromInterfaceTypeMap["outer"]
	_, toIsInner := toInterfaceTypeMap["inner"]
	_, toIsOuter := toInterfaceTypeMap["outer"]

	switch {
	case fromIsInner && fromIsOuter && toIsInner && toIsOuter:
		return "both"
	case fromIsInner && toIsInner:
		return "inner"
	case fromIsOuter && toIsOuter:
		return "outer"
	default:
		return ""
	}
}

func (r *CombinationFilters) addScheduledConnectionsInfo(
	ctx context.Context,
	comb v1.TargetCombinationStruct,
	ftMap map[string]*v1.FunctionTarget,
	connDtlInfoMap map[string]connectionDetailInfo,
	sortedFuncInstances []FuncOrderInfoStruct,
	funcInstNameToFuncNameMap map[string]string,
	interfacesPerRegionTypesPerFunction map[string]map[string]map[string][]string,
	graph *GraphManager,
) []v1.TargetCombinationStruct {

	// Array of EntityIds
	// m routes (Array of Entities) are stored for each nth connection.
	var routeIndex int
	routesPerConnection := make([][][]string, 0)
	routeRelatedInfoMap := make(map[int]routeRelatedInfo)

	// Create route information and store it in routesPerConnection
	r.createRoutesPerConnection(ctx, comb, ftMap, connDtlInfoMap, sortedFuncInstances,
		funcInstNameToFuncNameMap, interfacesPerRegionTypesPerFunction, graph,
		&routesPerConnection, routeRelatedInfoMap, &routeIndex)

	// If there is a connection with only an empty route, this combination is invalid.
	for ri := range routesPerConnection {
		if len(routesPerConnection[ri]) == 0 {
			return nil
		}
	}

	// Create a combination from a [][][]string
	routeCombs := generateCombinations[[]string](routesPerConnection)

	// Add UsedType to Route and make it a WBConnectionRoute
	routeWithTypesPerComb := make([][][]v1.WBConnectionPath, len(routeCombs))
	for ri, routeComb := range routeCombs {
		routeWithTypesPerComb[ri] = make([][]v1.WBConnectionPath, len(routeComb))
		for rj, route := range routeComb {
			routeWithTypesPerComb[ri][rj] = graph.GetIOUsedType(route)
		}
	}

	// Create ScheduledConnections based on the combination of Routes.
	// Add ScheduledConnections to comb
	combs := make([]v1.TargetCombinationStruct, len(routeWithTypesPerComb))
	for ri, routeWithTypeComb := range routeWithTypesPerComb {
		combs[ri] = comb
		combs[ri].ScheduledConnections =
			make([]v1.ConnectionScheduleInfo, len(routeWithTypeComb))

		for rj, routeWithUsedType := range routeWithTypeComb {

			combs[ri].ScheduledConnections[rj] =
				createScheduledConnectionsInfo(
					routeWithUsedType,
					connDtlInfoMap[routeRelatedInfoMap[rj].fromFuncKey+"-"+
						routeRelatedInfoMap[rj].fromRegionType+"-"+
						routeRelatedInfoMap[rj].toFuncKey+"-"+
						routeRelatedInfoMap[rj].toRegionType],
					*graph,
				)
		}
	}

	return combs
}

func (r *CombinationFilters) createRoutesPerConnection(
	ctx context.Context,
	comb v1.TargetCombinationStruct,
	ftMap map[string]*v1.FunctionTarget,
	connDtlInfoMap map[string]connectionDetailInfo,
	sortedFuncInstances []FuncOrderInfoStruct,
	funcInstNameToFuncNameMap map[string]string,
	interfacesPerRegionTypesPerFunction map[string]map[string]map[string][]string,
	graph *GraphManager,
	routesPerConnection *[][][]string,
	routeRelatedInfoMap map[int]routeRelatedInfo,
	routeIndex *int,
) {

	for _, sortedFuncInstance := range sortedFuncInstances {
		fromFuncInstName := sortedFuncInstance.FuncKey
		fromFuncName := funcInstNameToFuncNameMap[fromFuncInstName]
		fromFt := ftMap[getFunctionTargetName(comb.ScheduledFunctions[fromFuncInstName])]
		interfacePerRegionTypesForFrom := interfacesPerRegionTypesPerFunction[fromFuncName]

		// Create a route from the StartPoint to the first Function
		for _, startPointCandidate := range sortedFuncInstance.FromFunctions {
			if strings.HasPrefix(startPointCandidate, StartPointKeyword) {
				connType := connDtlInfoMap[startPointCandidate+"-"+StartPointKeyword+"-"+
					sortedFuncInstance.FuncKey+"-"+fromFt.Status.RegionType].connType
				addRoute(ctx, nil, fromFt, nil, interfacePerRegionTypesForFrom,
					connType, graph, routesPerConnection)
				addRouteRelatedInfo(startPointCandidate, fromFuncInstName,
					StartPointKeyword, fromFt.Status.RegionType, routeRelatedInfoMap, routeIndex)
			}
		}

		// Create routes between Functions and from the last Function to the EndPoint
		for _, toFuncInstName := range sortedFuncInstance.ToFunctions {
			toFuncName := funcInstNameToFuncNameMap[toFuncInstName]
			toFt := ftMap[getFunctionTargetName(comb.ScheduledFunctions[toFuncInstName])]
			interfacePerRegionTypesForTo := interfacesPerRegionTypesPerFunction[toFuncName]

			var toRegionType string
			if strings.HasPrefix(toFuncInstName, EndPointKeyword) {
				// Create a route from the last Function to the EndPoint.
				toRegionType = EndPointKeyword
				connType := connDtlInfoMap[fromFuncInstName+"-"+fromFt.Status.RegionType+"-"+
					toFuncInstName+"-"+toRegionType].connType
				addRoute(ctx, fromFt, nil, interfacePerRegionTypesForFrom, nil,
					connType, graph, routesPerConnection)
				addRouteRelatedInfo(fromFuncInstName, toFuncInstName,
					fromFt.Status.RegionType, toRegionType, routeRelatedInfoMap, routeIndex)
			} else {
				// Create a route between functions
				toRegionType = toFt.Status.RegionType
				connType := connDtlInfoMap[fromFuncInstName+"-"+fromFt.Status.RegionType+"-"+
					toFuncInstName+"-"+toRegionType].connType
				addRoute(ctx, fromFt, toFt, interfacePerRegionTypesForFrom, interfacePerRegionTypesForTo,
					connType, graph, routesPerConnection)
				addRouteRelatedInfo(fromFuncInstName, toFuncInstName,
					fromFt.Status.RegionType, toRegionType, routeRelatedInfoMap, routeIndex)
			}
		}
	}
}

func addRoute(
	ctx context.Context,
	fromFt, toFt *v1.FunctionTarget,
	interfaceserRegionTypesForFrom, interfaceserRegionTypesForTo map[string]map[string][]string,
	connType string,
	graph *GraphManager,
	routesPerConnection *[][][]string,
) {
	*routesPerConnection = append(*routesPerConnection, generateRoutes(
		ctx,
		fromFt,
		toFt,
		interfaceserRegionTypesForFrom,
		interfaceserRegionTypesForTo,
		connType,
		graph,
	))
}

func addRouteRelatedInfo(
	fromFuncKey,
	toFuncKey,
	fromRegionType,
	toRegionType string,
	routeRelatedInfoMap map[int]routeRelatedInfo,
	routeIndex *int,
) {
	routeRelatedInfoMap[*routeIndex] = routeRelatedInfo{
		fromFuncKey:    fromFuncKey,
		toFuncKey:      toFuncKey,
		fromRegionType: fromRegionType,
		toRegionType:   toRegionType,
	}
	*routeIndex++
}

func generateRoutes(ctx context.Context,
	fromFt *v1.FunctionTarget,
	toFt *v1.FunctionTarget,
	interfacePerRegionTypesForFrom map[string]map[string][]string,
	interfacePerRegionTypesForTo map[string]map[string][]string,
	connType string,
	graph *GraphManager,
) [][]string {

	l := log.FromContext(ctx)

	// Get the From and To EntityIds for route search
	var fromEntityId string
	if fromFt == nil {
		fromEntityId = graph.GetEntityIdOfExternalNode()
	} else {
		fromEntityId = graph.GetEntityId(fromFt)
	}
	if fromEntityId == "" {
		return make([][]string, 0)
	}

	var toEntityId string
	if toFt == nil {
		toEntityId = graph.GetEntityIdOfExternalNode()
	} else {
		toEntityId = graph.GetEntityId(toFt)
	}
	if toEntityId == "" {
		return make([][]string, 0)
	}

	// Route calculation method differs depending on the connection method
	switch connType {
	case "inner":
		return getInnerRoute(ctx, fromFt, toFt, fromEntityId, toEntityId, graph)
	case "outer":
		return getOuterRoute(
			ctx, fromFt, toFt, fromEntityId, toEntityId,
			interfacePerRegionTypesForFrom, interfacePerRegionTypesForTo, graph)
	case "both":
		mergeRouteList := getInnerRoute(ctx, fromFt, toFt, fromEntityId, toEntityId, graph)
		outerRouteList := getOuterRoute(
			ctx, fromFt, toFt, fromEntityId, toEntityId,
			interfacePerRegionTypesForFrom, interfacePerRegionTypesForTo, graph)
		for _, outerRoute := range outerRouteList {
			if len(outerRoute) != 0 {
				mergeRouteList = append(mergeRouteList, outerRoute)
			}
		}
		if len(mergeRouteList) == 0 {
			return make([][]string, 0)
		}
		return mergeRouteList

	default:
		l.Info("unsupported value for connType=" + connType)
		return make([][]string, 0)
	}

}

func getInnerRoute(
	ctx context.Context,
	fromFt *v1.FunctionTarget,
	toFt *v1.FunctionTarget,
	fromEntityId string,
	toEntityId string,
	graph *GraphManager,
) [][]string {

	l := log.FromContext(ctx)

	var routeList [][]string
	if fromEntityId == toEntityId {
		routeList = graph.FindReturnPathAtMemory(toFt.Status.NodeName, fromEntityId, toEntityId)
	} else {
		routeList = graph.FindInnerShortestPaths(toFt.Status.NodeName, fromEntityId, toEntityId, k)
	}

	// Remove empty routes. If none remain, exit.
	excludeIndexes := make([]int32, 0)
	for ri, route := range routeList {
		if len(route) == 0 {
			excludeIndexes = append(excludeIndexes, int32(ri))
		}
	}

	exclude(&routeList, excludeIndexes)
	if len(routeList) == 0 {
		l.Info("route not found", "FT", toFt.Name, "fromEntityId", fromEntityId, "toEntityId", toEntityId)
		return make([][]string, 0)
	}
	// Check whether there is a path from the previous deployment destination to the memory of the deployment candidate node, and whether there is a path from the deployment candidate to the memory of the current node
	// If from and to are the same, the route search has confirmed that the route is via memory, so the available I/F check is skipped.
	if !(fromEntityId == toEntityId) && !graph.HasRouteToMemory(toFt.Status.NodeName, fromEntityId, toEntityId) {
		l.Info("Usable I/F check is NG", "fromEntityId", fromEntityId, "toEntityId", toEntityId,
			"fromFt", fromFt.Name, "toFT", toFt.Name,
			"fromRegionType", fromFt.Status.RegionType, "toRgionType", toFt.Status.RegionType,
			"connType", "inner")
		return make([][]string, 0)
	}
	return routeList
}

func getOuterRoute(
	ctx context.Context,
	fromFt, toFt *v1.FunctionTarget,
	fromEntityId, toEntityId string,
	interfacePerRegionTypesForFrom,
	interfacePerRegionTypesForTo map[string]map[string][]string,
	graph *GraphManager,
) [][]string {

	l := log.FromContext(ctx)

	// Check "Back stage is FPGA", "Front stage is FPGA", "Route is round trip route", and "Deployment to the same node"
	isFromFpga := fromFt != nil && strings.Contains(fromFt.Status.RegionType, "alveo")
	isToFpga := toFt != nil && strings.Contains(toFt.Status.RegionType, "alveo")
	isRoundTrip := fromEntityId == toEntityId
	isSameNode := fromFt != nil && toFt != nil && fromFt.Status.NodeName == toFt.Status.NodeName

	// route retrieval
	var routeList [][]string
	if isFromFpga || isToFpga {
		// If the previous or next stage is an FPGA, get the route including the FPGA NIC.
		routeList = graph.FindOuterShortestPaths(fromEntityId, toEntityId, k)
	} else {
		if isRoundTrip || isSameNode {
			routeList = graph.FindReturnPathAtHostNic(toFt.Status.NodeName, fromEntityId, toEntityId)
		} else {
			routeList = graph.FindOuterExcludeFpgaNicShortestPaths(fromEntityId, toEntityId, k)
		}
	}

	// Remove empty routes. If none remain, exit.
	excludeIndexes := make([]int32, 0)
	for ri, route := range routeList {
		if len(route) == 0 {
			excludeIndexes = append(excludeIndexes, int32(ri))
		}
	}
	exclude(&routeList, excludeIndexes)
	if len(routeList) == 0 {
		l.Info("route not found", "fromEntityId", fromEntityId, "toEntityId", toEntityId)
		return make([][]string, 0)
	}

	// If the destination is nil, it is the last stage so exit
	if toFt == nil {
		return routeList
	}

	// Validate each route
	excludeIndexes = make([]int32, 0)
	for ri, route := range routeList {

		// Check if it is an external connection (TCP connections within a node are treated as external connections)
		if !graph.IsOuterRoute(route, isFromFpga, isToFpga, isRoundTrip, isSameNode) {
			l.Info("Not an outer connection route", "FT", toFt.Name, "fromEntityId", fromEntityId, "toEntityId", toEntityId)
			excludeIndexes = append(excludeIndexes, int32(ri))
			continue
		}

		// If the previous stage is an FPGA, check if the InterfaceType (e.g. dev25gether) is included in the available I/F for the output of the previous function.
		if isFromFpga {
			if !checkInterfaceType(route, route[1], fromFt, interfacePerRegionTypesForFrom, outputInterfaceTypeKey, graph) {
				l.Info("Usable I/F check is NG", "fromEntityId", fromEntityId, "toEntityId", toEntityId,
					"fromFt", fromFt.Name, "toFT", toFt.Name,
					"fromRegionType", fromFt.Status.RegionType, "toRgionType", toFt.Status.RegionType,
					"connType", "outer")
				excludeIndexes = append(excludeIndexes, int32(ri))
				continue
			}
		}

		if isToFpga {
			// The one before the end point (to) of the route (the I/F (adjacent entity) on the incoming side of the deployment candidate)
			// Check if the InterfaceType (e.g. dev25gether) is included in the available I/F for the input of the target function.
			if !checkInterfaceType(route, route[len(route)-2], toFt, interfacePerRegionTypesForTo, inputInterfaceTypeKey, graph) {
				l.Info("Usable I/F check is NG", "fromEntityId", fromEntityId, "toEntityId", toEntityId,
					"connType", "outer")
				excludeIndexes = append(excludeIndexes, int32(ri))
			}
		} else if !(isRoundTrip || isSameNode) {
			// Check if the route passes through the same node NIC (same InterfaceType (ex. host100gether)) as the deployment candidate.
			isValid := false
			for _, eId := range route {
				if checkInterfaceType(route, eId, toFt, interfacePerRegionTypesForTo, inputInterfaceTypeKey, graph) {
					isValid = true
					break
				}
			}
			if !isValid {
				l.Info("Usable I/F check is NG", "fromEntityId", fromEntityId, "toEntityId", toEntityId,
					"fromFt", fromFt.Name, "toFT", toFt.Name,
					"fromRegionType", fromFt.Status.RegionType, "toRgionType", toFt.Status.RegionType,
					"connType", "outer")
				excludeIndexes = append(excludeIndexes, int32(ri))
			}
		}
	}

	// Remove invalid routes
	exclude(&routeList, excludeIndexes)

	return routeList
}

func checkInterfaceType(
	route []string,
	entityId string,
	functionTarget *v1.FunctionTarget,
	interfacePerRegionTypes map[string]map[string][]string,
	interfaceTypeKey string,
	graph *GraphManager,
) bool {

	interfaceType, nodeName := graph.GetInterfaceTypeAndNodeName(entityId)

	if nodeName != functionTarget.Status.NodeName {
		return false
	}
	if interfacePerRegionTypes == nil {
		return false
	}
	for _, availableInterfaceType := range interfacePerRegionTypes[functionTarget.Status.RegionType][interfaceTypeKey] {
		if interfaceType == availableInterfaceType {
			return true
		}
	}
	return false
}

func createScheduledConnectionsInfo(
	route []v1.WBConnectionPath,
	cdi connectionDetailInfo,
	graph GraphManager,
) v1.ConnectionScheduleInfo {
	var connMethod, fromInterfaceType, toInterfaceType string
	switch cdi.connType {
	case "inner":
		connMethod = connMethodHostMem
		fromInterfaceType = cdi.fromInterfaceTypes["inner"]
		toInterfaceType = cdi.toInterfaceTypes["inner"]
	case "outer":
		connMethod = connMethodHost100g
		fromInterfaceType = cdi.fromInterfaceTypes["outer"]
		toInterfaceType = cdi.toInterfaceTypes["outer"]
	case "both":
		// Check whether it is an inside connection or an outside connection based on the determined route
		if graph.IsOuterWBConnRoute(route) {
			connMethod = connMethodHost100g
			fromInterfaceType = cdi.fromInterfaceTypes["outer"]
			toInterfaceType = cdi.toInterfaceTypes["outer"]
		} else {
			connMethod = connMethodHostMem
			fromInterfaceType = cdi.fromInterfaceTypes["inner"]
			toInterfaceType = cdi.toInterfaceTypes["inner"]
		}
	}
	csi := v1.ConnectionScheduleInfo{
		From: v1.FromToFunctionScheduleInfo{
			FunctionKey:   cdi.fromFuncKey,
			Port:          &cdi.fromPort,
			InterfaceType: &fromInterfaceType,
		},
		To: v1.FromToFunctionScheduleInfo{
			FunctionKey:   cdi.toFuncKey,
			Port:          &cdi.toPort,
			InterfaceType: &toInterfaceType,
		},
		ConnectionMethod: connMethod,
		ConnectionPath:   route,
	}
	return csi
}
