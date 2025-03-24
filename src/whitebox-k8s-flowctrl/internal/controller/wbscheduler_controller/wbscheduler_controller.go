/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wbschedulercontroller

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	dataflowstatus "github.com/compsysg/whitebox-k8s-flowctrl/lib/dataflowStatus"
	"github.com/compsysg/whitebox-k8s-flowctrl/lib/filter_template"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/scheduler_common"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	errStatusIsNotScheduling        = errors.New("Status is not Scheduling")
	errValidCombinationDoesNotExist = errors.New("Valid combination does not exist")
	errFailedToFetchConfigMap       = errors.New("CantFindConfigMap")
	errInvalidConfigMapParamter     = errors.New("Invalid ConfigMap Parameter")
	errDataFlowIsDeleted            = errors.New("DataFlow is deleted")
	errSchedulingDataIsProcessing   = errors.New("SchedulingData for another dataflow is still processing")
)

// Controller Name
const controllerName = "wbscheduler"

// Logger Setting
const (
	loggerKeyController      = "wbscheduler"
	loggerKeyControllerGroup = "example.com"
	loggerKeyControllerKind  = "WBScheduler"
)

const (
	// The connection method set in Type of ConnectionScheduleInfo is host-mem
	connMethodHostMem string = "host-mem"
	// The connection method set in Type of ConnectionScheduleInfo is host-100gether
	connMethodHost100g string = "host-100gether"
	// Prefix to identify FPGA from region type (RegionType)
	fpgaPrefix string = "alveo"
)

const initialFilterIndex int32 = 0

// WBschedulerReconciler reconciles a WBscheduler object
type WBschedulerReconciler struct {
	client.Client
	Scheme                *runtime.Scheme
	Recorder              record.EventRecorder
	DefaultFilterPipeline string
	WaitTimeSec           int
	RequeueTimeSec        int
}

//+kubebuilder:rbac:groups=ntt-hpc.example.com,resources=dataflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ntt-hpc.example.com,resources=dataflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ntt-hpc.example.com,resources=dataflows/finalizers,verbs=update
//+kubebuilder:rbac:groups=ntt-hpc.example.com,resources=schedulingdata,verbs=get;list;create;update;delete
//+kubebuilder:rbac:groups=ntt-hpc.example.com,resources=schedulingdata/status,verbs=watch;get;update;
//+kubebuilder:rbac:groups=ntt-hpc.example.com,resources=schedulingdata/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WBscheduler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *WBschedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	l := log.FromContext(ctx)
	l.Info("start Scheduler " + req.NamespacedName.Namespace)

	df, sd, err := r.initialize(ctx, req)
	if err == errSchedulingDataIsProcessing {
		return ctrl.Result{RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, nil
	}

	if !r.isSchedulingFinish(sd) {
		return ctrl.Result{}, err
	}

	resComb, err := r.checkAndGetResult(ctx, sd, df)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * time.Duration(r.RequeueTimeSec)}, nil
	}

	r.allocate(df, resComb)

	if err := r.finalize(ctx, df, sd); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WBschedulerReconciler) initialize(ctx context.Context, req ctrl.Request) (*ntthpcv1.DataFlow, *ntthpcv1.SchedulingData, error) {

	df, err := r.getDataFlow(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	sd, err := r.getOrCreateSchedulingData(ctx, df)
	if err != nil {
		return nil, nil, err
	}

	return df, sd, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WBschedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	p := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(controllerName).
		For(&ntthpcv1.DataFlow{}).
		Owns(&ntthpcv1.SchedulingData{}).
		WithLogConstructor(func(req *reconcile.Request) logr.Logger {
			return mgr.GetLogger().WithValues("controller", loggerKeyController, "controllerGroup",
				loggerKeyControllerGroup, "controllerKind", loggerKeyControllerKind)
		}).
		WithEventFilter(p).
		Complete(r)
}

func (r *WBschedulerReconciler) getDataFlow(ctx context.Context, req ctrl.Request) (*ntthpcv1.DataFlow, error) {

	l := log.FromContext(ctx)

	var df *ntthpcv1.DataFlow = &ntthpcv1.DataFlow{}
	if err := r.Get(ctx, req.NamespacedName, df); err != nil {
		l.Error(err, "unable to fetch dataflow")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return nil, client.IgnoreNotFound(err)
	}

	// check dataflow status
	if df.Status.Status != "Scheduling in progress" {
		l.Info("Dataflow Status is not 'Scheduling in progress'")
		return nil, errStatusIsNotScheduling
	}

	// If not yet scheduled, schedule
	// shcedule info init
	if df.Status.ScheduledFunctions == nil {
		df.Status.ScheduledFunctions = make(map[string]ntthpcv1.FunctionScheduleInfo)
	}
	if df.Status.ScheduledConnections == nil {
		df.Status.ScheduledConnections = []ntthpcv1.ConnectionScheduleInfo{}
	}

	return df, nil
}

func (r *WBschedulerReconciler) getOrCreateSchedulingData(
	ctx context.Context,
	df *ntthpcv1.DataFlow,
) (*ntthpcv1.SchedulingData, error) {
	var sd *ntthpcv1.SchedulingData = &ntthpcv1.SchedulingData{}

	l := log.FromContext(ctx)

	// Recurse until SchedulingData is generated,
	// This will cause a SEGV if the DataFlow is deleted midway.
	if df == nil {
		return nil, errDataFlowIsDeleted
	}

	sdName := df.ObjectMeta.Name
	sdNameSpace := df.ObjectMeta.Namespace

	if err := r.Get(ctx, client.ObjectKey{Namespace: sdNameSpace, Name: sdName}, sd); err != nil {

		// fetch SchedulingDataList
		var schedulingDataList ntthpcv1.SchedulingDataList
		if err := r.List(ctx, &schedulingDataList, &client.ListOptions{}); err != nil {
			l.Error(err, "Failed to fetch SchedulingDataList")
			return nil, err
		}

		// Create a new SchedulingData only if SchedulingData for another dataflow does not exist
		if len(schedulingDataList.Items) == 0 {
			if err := r.createSchedulingData(ctx, df, sdName, sdNameSpace); err != nil {
				l.Error(err, fmt.Sprintf("create SchedulingData %v failed", sdName))
				return nil, err
			}
			l.Info(fmt.Sprintf("SchedulingData %v is generated successfully", sdName))
			return r.getOrCreateSchedulingData(ctx, df)
		} else {
			l.Error(errSchedulingDataIsProcessing, fmt.Sprintf("create SchedulingData %v failed", sdName))
			return nil, errSchedulingDataIsProcessing
		}
	}

	if sd.Status.Status == "" {
		sd.Status.Status = filter_template.FilteringStatus
		sd.Status.CurrentFilterIndex = func() *int32 { var v int32 = initialFilterIndex; return &v }()
		if err := r.Status().Update(ctx, sd); err != nil {
			return nil, err
		}
	}

	// l.Info(fmt.Sprintf("Spec : %v", sd.Spec))
	// l.Info(fmt.Sprintf("Status %v", sd.Status))

	return sd, nil
}

func (r *WBschedulerReconciler) isSchedulingFinish(sd *ntthpcv1.SchedulingData) bool {
	return sd.Status.Status == filter_template.FinishStatus
}

func (r *WBschedulerReconciler) checkAndGetResult(
	ctx context.Context,
	sd *ntthpcv1.SchedulingData,
	df *ntthpcv1.DataFlow,
) (*ntthpcv1.TargetCombinationStruct, error) {

	l := log.FromContext(ctx)

	sdName := sd.ObjectMeta.Name

	if len(sd.Status.TargetCombinations) == 0 {
		l.Error(errValidCombinationDoesNotExist, "sd.Status.TargetCombinations's length is 0")
		if err := r.Delete(ctx, sd); err != nil {
			l.Error(err, "unable to delete SchedulingData")
			return nil, err
		}
		l.Info(fmt.Sprintf("SchedulingData %v is deleted successfully", sdName))
		return nil, errValidCombinationDoesNotExist
	}

	// Get the highest Score
	best := &sd.Status.TargetCombinations[0]
	for i := 1; i < len(sd.Status.TargetCombinations); i++ {
		comb_i := i
		if *best.Score < *sd.Status.TargetCombinations[comb_i].Score {
			best = &sd.Status.TargetCombinations[comb_i]
		}
	}

	return best, nil
}

//nolint:unused // FIXME: remove this function
func (r *WBschedulerReconciler) generateConnectionScheduleInfoList(
	resComb *ntthpcv1.TargetCombinationStruct,
	df *ntthpcv1.DataFlow,
) ([]ntthpcv1.ConnectionScheduleInfo, map[string]ntthpcv1.ConnectionIfStruct) {

	conScheInfoList := make([]ntthpcv1.ConnectionScheduleInfo, 0)
	connectionInterfaceListMap := make(map[string]ntthpcv1.ConnectionIfStruct, 1)

	fc := df.Status.FunctionChain

	cks := make([]string, len(resComb.ScheduledConnections))
	for i, c := range resComb.ScheduledConnections {
		cks[i] = c.ConnectionMethod
	}
	funcScheInfoMap := resComb.ScheduledFunctions

	type conTypeKey struct {
		From string
		To   string
	}
	conTypeMap := make(map[conTypeKey]string)

	for i, v := range fc.Spec.Connections {
		conTypeMap[conTypeKey{From: v.From.FunctionKey, To: v.To.FunctionKey}] = cks[i]
	}

	// for _, fcFuncKey := range sortFuncKeys {
	for fcFuncKey, _ := range fc.Spec.Functions {
		// connection schedule info create
		for _, cv := range fc.Spec.Connections {
			// connection information of the Scheduled Function
			if cv.From.FunctionKey == fcFuncKey || (strings.HasPrefix(cv.From.FunctionKey, StartPointKeyword) && cv.To.FunctionKey == fcFuncKey) {
				for _, conv := range df.Status.ConnectionType {
					// ConnectionType infomation
					if conv.Spec.ConnectionTypeName == cv.ConnectionTypeName || cv.ConnectionTypeName == "auto" {
						for availableType, _ := range conv.Status.AvailableInterfaces {
							// Get deployable DeviceType from ConnectionType
							if strings.Contains(availableType, funcScheInfoMap[fcFuncKey].DeviceType) {
								var csi ntthpcv1.ConnectionScheduleInfo

								if cv.From.FunctionKey == fcFuncKey {
									csi.From.FunctionKey = fcFuncKey
									csi.To.FunctionKey = cv.To.FunctionKey
								} else {
									csi.From.FunctionKey = cv.From.FunctionKey
									csi.To.FunctionKey = fcFuncKey
								}

								// csi.Type = cv.ConnectionType
								csi.ConnectionMethod = conTypeMap[conTypeKey{cv.From.FunctionKey, cv.To.FunctionKey}]

								var connectionInterface ntthpcv1.ConnectionIfStruct
								connectionInterface.InterfaceList = PutStrSliceNoDupe(connectionInterfaceListMap[fcFuncKey].InterfaceList, availableType)
								*connectionInterface.NodeName = resComb.ScheduledFunctions[fcFuncKey].NodeName
								connectionInterfaceListMap[fcFuncKey] = connectionInterface
								conScheInfoList = append(conScheInfoList, csi)
							}
						}
					}
				}
			}
		}
	}

	return conScheInfoList, connectionInterfaceListMap

}

// Get the Type of ScheduledConnections.
// Returns "host-mem" for FPGA-GPU Connections, and "host-100gether" for other Connections.
// The function name does not match "wb-start-of-chain" or "wb-end-of-chain" at the beginning, and
// If the DeviceType of the FunctionTarget where the Function is deployed is a string starting with "alveo" → FPGA
// If not applicable to the FPGA above → GPU
func (r *WBschedulerReconciler) getConnectionMethod(connection ntthpcv1.ConnectionStruct,
	fromFunctionScheduleInfo ntthpcv1.FunctionScheduleInfo, toFunctionScheduleInfo ntthpcv1.FunctionScheduleInfo) string {
	if strings.HasPrefix(connection.From.FunctionKey, StartPointKeyword) || strings.HasPrefix(connection.To.FunctionKey, EndPointKeyword) {
		return connMethodHost100g
	}

	fromDeviceType := fromFunctionScheduleInfo.DeviceType
	toDeviceType := toFunctionScheduleInfo.DeviceType
	if isFPGA(fromDeviceType) && !isFPGA(toDeviceType) || !isFPGA(fromDeviceType) && isFPGA(toDeviceType) {
		return connMethodHostMem
	}
	return connMethodHost100g
}

func isFPGA(regionKind string) bool {
	return strings.HasPrefix(regionKind, fpgaPrefix)
}

//nolint:unused // FIXME: remove this function
func (r *WBschedulerReconciler) getIoName(conTypeList []*ntthpcv1.ConnectionType, funcInterfaceListMap map[string]ntthpcv1.ConnectionIfStruct, conScheInfo ntthpcv1.ConnectionScheduleInfo) string {
	fromIfList := funcInterfaceListMap[conScheInfo.From.FunctionKey].InterfaceList
	toIfList := funcInterfaceListMap[conScheInfo.To.FunctionKey].InterfaceList
	var targetIoList []string
	for _, fromIf := range fromIfList {
		for _, toIf := range toIfList {
			sameNode_flg := false
			if funcInterfaceListMap[conScheInfo.From.FunctionKey].NodeName == funcInterfaceListMap[conScheInfo.To.FunctionKey].NodeName {
				sameNode_flg = true
			}
			for _, conType := range conTypeList {
				// TODO: want to delete connectiontype const (example "node" "rack")
				if conScheInfo.ConnectionMethod == "auto" {
					if (StringToConst(conType.Spec.ConnectionTypeName) <= StringToConst("node") && !sameNode_flg) ||
						(StringToConst(conType.Spec.ConnectionTypeName) >= StringToConst("rack") && sameNode_flg) {
						continue
					}
				} else if conScheInfo.ConnectionMethod == "node" {
					if (StringToConst(conType.Spec.ConnectionTypeName) <= StringToConst("node") && !sameNode_flg) ||
						StringToConst(conType.Spec.ConnectionTypeName) >= StringToConst("rack") {
						continue
					}
				} else if conType.Spec.ConnectionTypeName != conScheInfo.ConnectionMethod {
					continue
				}
				var candidateIoList []string
				for ifName, destinationsMap := range conType.Status.AvailableInterfaces {
					// Select common interface
					for _, checkIo := range [2]string{fromIf, toIf} {
						if checkIo == ifName {
							for distIo, _ := range destinationsMap.Destinations {
								appendFlg := false
								for _, candidateIo := range candidateIoList {
									if candidateIo == distIo {
										targetIoList = append(targetIoList, distIo)
										appendFlg = true
									}
								}
								if !appendFlg {
									candidateIoList = append(candidateIoList, distIo)
								}
							}
						}
					}
				}
			}
		}
	}

	if len(targetIoList) > 0 {
		ioPriorityList := []string{"direct", "host-mem", "pcie", "rdma", "host-100gether"}
		for _, ioPriority := range ioPriorityList {
			for _, targetIo := range targetIoList {
				if strings.Contains(targetIo, "-direct") {
					targetIo = "direct"
				}
				if ioPriority == targetIo {
					return targetIo
				}
			}
		}
	}
	return "not found"
}

func (r *WBschedulerReconciler) createSchedulingData(ctx context.Context, df *ntthpcv1.DataFlow, scheDataName string, scheDataNamespace string) error {

	var (
		filterPipeline []string
		err            error
	)

	filterPipeline, err = r.loadFilterPipeline(ctx, df)

	switch err {
	case errFailedToFetchConfigMap:
		filterPipeline = strings.Split(r.DefaultFilterPipeline, ",")
	case errInvalidConfigMapParamter:
		return errInvalidConfigMapParamter
	}

	dep := ntthpcv1.SchedulingData{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: scheDataNamespace,
			Name:      scheDataName,
		},
		Spec: ntthpcv1.SchedulingDataSpec{
			FilterPipeline: filterPipeline,
		},
	}

	ctrl.SetControllerReference(df, &dep, r.Scheme)

	if err := r.Create(ctx, &dep); err != nil {
		return err
	}

	return nil
}

func (r *WBschedulerReconciler) loadFilterPipeline(ctx context.Context, df *ntthpcv1.DataFlow) (filterPipeline []string, err error) {

	l := log.FromContext(ctx)
	var strategy string = ""

	defaultNamespaceCandidates := []string{df.ObjectMeta.Namespace, "default"}

	// Fetch userRequirement configmap from some namespace candidates. ["same-namespace-of-dataflow", "default"]
	cm, err := tryFetchConfigMapFromSeveralNameSpaceCandidates(
		r, ctx, GetDereferencedValueOrZeroValue(df.Spec.UserRequirement), defaultNamespaceCandidates)
	if err != nil {
		l.Info("Unable to fetch userRequirement configmap. Scheduler will use default FilterPipeline.")
		return filterPipeline, errFailedToFetchConfigMap
	}

	var foundStrategyFlag bool = false
	// setting dataflow
	for k, v := range cm.Data {
		if k == "strategy" {
			strategy = v
			foundStrategyFlag = true
		}
	}

	if !foundStrategyFlag {
		return filterPipeline, errFailedToFetchConfigMap
	}

	return r.fetchStrategy(ctx, strategy, defaultNamespaceCandidates)
}

func (r *WBschedulerReconciler) fetchStrategy(
	ctx context.Context,
	strategy string,
	nameSpaceCands []string,
) ([]string, error) {

	l := log.FromContext(ctx)
	cm, err := tryFetchConfigMapFromSeveralNameSpaceCandidates(r, ctx, strategy, nameSpaceCands)
	if err != nil {
		l.Info("Unable to fetch strategy configmap. Scheduler will use default FilterPipeline.")
		return nil, errFailedToFetchConfigMap
	}

	if fp, ok := cm.Data["filterPipeline"]; ok {
		var filterPipeline []string
		parseYAML(fp, &filterPipeline)
		return filterPipeline, nil
	}

	if rp, ok := cm.Data["referenceParameter"]; ok {
		return r.fetchStrategy(ctx, rp, nameSpaceCands)
	}

	return nil, errFailedToFetchConfigMap
}

// func (r *WBschedulerReconciler) checkScheduledDeviceTypesIsValid(
// 	ctx *context.Context,
// 	df *ntthpcv1.DataFlow,
// 	funcScheInfoMap map[string]ntthpcv1.FunctionScheduleInfo) bool {
//
// 	l := log.FromContext(*ctx)
//
// 	fkList := df.Status.FunctionType
// 	fMap := df.Status.FunctionChain.Spec.Functions
//
// 	// check device candidata able to apply target functions
// 	// for _, fcFuncKey := range sortFuncKeys {
// 	for fcFuncKey, v := range fMap {
// 		candDevKind := funcScheInfoMap[fcFuncKey].DeviceType
// 		targetFunc := v.FunctionName
//
// 		isDeviceSupported := false
// 		for _, funcKind := range fkList {
// 			if targetFunc == funcKind.Spec.Name {
// 				for _, availableDevKind := range funcKind.Status.DeviceTypeCandidates {
// 					if candDevKind == availableDevKind {
// 						isDeviceSupported = true
// 						break
// 					}
// 				}
// 			}
// 		}
//
// 		if !isDeviceSupported {
// 			l.Error(errInvalidDeviceTypeScheduled,
// 				fmt.Sprintf("Invalid Device Kind is Scheduled. FunctionType : %v, Scheduled DeviceType : %v",
// 					fcFuncKey, candDevKind))
// 			return false
// 		}
// 	}
//
// 	return true
// }

func (r *WBschedulerReconciler) allocate(df *ntthpcv1.DataFlow, resComb *ntthpcv1.TargetCombinationStruct) {

	df.Status.ScheduledFunctions = resComb.ScheduledFunctions

	if resComb.ScheduledConnections == nil {
		df.Status.ScheduledConnections = r.createScheduledConnectionsInfos(resComb.ScheduledFunctions, df.Status.FunctionChain.Spec.Connections)
	} else {
		conns := df.Status.FunctionChain.Spec.Connections
		for i, conn := range conns {
			csi := resComb.ScheduledConnections[i]
			var connType string
			if strings.HasPrefix(conn.From.FunctionKey, StartPointKeyword) || strings.HasPrefix(conn.To.FunctionKey, EndPointKeyword) {
				connType = "host-100gether"
			} else {
				connType = csi.ConnectionMethod
			}
			df.Status.ScheduledConnections = append(df.Status.ScheduledConnections, ntthpcv1.ConnectionScheduleInfo{
				From: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey:   conn.From.FunctionKey,
					Port:          csi.From.Port,
					InterfaceType: csi.From.InterfaceType,
				},
				To: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey:   conn.To.FunctionKey,
					Port:          csi.To.Port,
					InterfaceType: csi.To.InterfaceType,
				},
				ConnectionPath:   csi.ConnectionPath,
				ConnectionMethod: connType,
			})
		}
	}
}

func (r *WBschedulerReconciler) finalize(ctx context.Context, df *ntthpcv1.DataFlow, sd *ntthpcv1.SchedulingData) error {

	l := log.FromContext(ctx)

	sdName := sd.ObjectMeta.Name

	// status update
	df.Status.Status = dataflowstatus.WB_CreationInProgress

	if err := r.Status().Update(ctx, df); err != nil {
		l.Error(err, "unable to update df status")
		return err
	}
	l.Info("Updated DataFlow.Status successfully")

	if err := r.Delete(ctx, sd); err != nil {
		l.Error(err, "failed to delete SchedulingData")
		return err
	}
	l.Info(fmt.Sprintf("SchedulingData %v is deleted successfully", sdName))

	// Waiting to respond to continuous deployment requests (temporary solution)
	time.Sleep(time.Duration(r.WaitTimeSec) * time.Second)

	return nil
}

func (r *WBschedulerReconciler) createScheduledConnectionsInfos(functionScheduleInfoMap map[string]ntthpcv1.FunctionScheduleInfo,
	connections []ntthpcv1.ConnectionStruct) []ntthpcv1.ConnectionScheduleInfo {

	connectionScheduleInfos := []ntthpcv1.ConnectionScheduleInfo{}
	for _, connection := range connections {
		connType := r.getConnectionMethod(connection, functionScheduleInfoMap[connection.From.FunctionKey], functionScheduleInfoMap[connection.To.FunctionKey])
		connectionScheduleInfos = append(connectionScheduleInfos, ntthpcv1.ConnectionScheduleInfo{
			From: ntthpcv1.FromToFunctionScheduleInfo{
				FunctionKey: connection.From.FunctionKey,
			},
			To: ntthpcv1.FromToFunctionScheduleInfo{
				FunctionKey: connection.To.FunctionKey,
			},
			ConnectionMethod: connType,
		})
	}
	return connectionScheduleInfos
}
