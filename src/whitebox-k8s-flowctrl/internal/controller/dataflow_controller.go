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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
)

// DataFlowReconciler reconciles a DataFlow object
type DataFlowReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var endPointKeyword = "wb-end-of-chain"
var startPointKeyword = "wb-start-of-chain"

func strSliceCheck(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

//nolint:unused // FIXME: remove this function
func putStrSliceNoDupe(slice []string, target string) []string {
	if !strSliceCheck(slice, target) {
		slice = append(slice, target)
	}
	return slice
}

//+kubebuilder:rbac:groups=example.com,resources=dataflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=dataflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=dataflows/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=wbfunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=wbfunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=wbfunctions/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=wbconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=wbconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=wbconnections/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=functionchains,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=functionchains/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=functionchains/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=functiontypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=functiontypes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=functiontypes/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=connectiontypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=connectiontypes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=connectiontypes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

func (r *DataFlowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var dataflow ntthpcv1.DataFlow

	l.Info("fetching DataFlow Resource")
	if err := r.Get(ctx, req.NamespacedName, &dataflow); err != nil {
		l.Error(err, "unable to fetch DataFlow")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var functionChain ntthpcv1.FunctionChain
	referFunctionChain := dataflow.Spec.FunctionChainRef

	// Branching by Dataflow status
	switch dataflow.Status.Status {
	case "":
		// FunctionChain parsing logic
		// Get FunctionChain
		err := r.Get(ctx, client.ObjectKey{Namespace: referFunctionChain.Namespace,
			Name: referFunctionChain.Name},
			&functionChain)
		l.Info("fetching FunctionChain Resource:" + functionChain.Name)
		if err != nil {
			l.Error(err, "unable to fetch FunctionChain")
			r.Recorder.Eventf(&dataflow,
				corev1.EventTypeNormal,
				"Reconcile",
				"unable to fetch FunctionChain")
			dataflow.Status.Status = "Error"
			r.patchApply1(ctx, req, dataflow)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		// Put FunctionChain in Status of DataFlow
		if functionChain.Status.Status == "Ready" {
			dataflow.Status.FunctionChain = &functionChain
			l.Info("Put FunctionChain: " + functionChain.Name + "  in Status of DataFlow")

			// TODO : Check Status of FunctionType and ConnectionType and add processing if Ready
			// Get ConnectionTypeList
			var connectiontypes ntthpcv1.ConnectionTypeList
			l.Info("fetching ConnectionTypeList")
			if err := r.List(ctx, &connectiontypes,
				&client.ListOptions{Namespace: functionChain.Spec.ConnectionTypeNamespace}); err != nil {
				return ctrl.Result{}, err
			}
			// Get and put ConnectionType in Status of DataFlow
			var alreadyConKindName []string
			for _, fcconnections := range functionChain.Spec.Connections {
				for _, cks := range connectiontypes.Items {
					if (fcconnections.ConnectionTypeName == cks.Name ||
						fcconnections.ConnectionTypeName == "auto") && !strSliceCheck(alreadyConKindName, cks.Name) {
						getcks := cks
						dataflow.Status.ConnectionType = append(dataflow.Status.ConnectionType, &getcks)
						alreadyConKindName = append(alreadyConKindName, cks.Name)
						l.Info("Put ConnectionType: " + cks.Name + " in Status of DataFlow")
					}
				}
			}
			// Put an empty slice in DataFlow.Status.ConnectionType when all ConnectionType parameters are "auto".
			if dataflow.Status.ConnectionType == nil {
				l.Info("ConnectionType parameter is all 'auto' and not saved in DataFlow")
				dataflow.Status.ConnectionType = []*ntthpcv1.ConnectionType{}
			}

			// Get FunctionTypeList
			var functiontypes ntthpcv1.FunctionTypeList
			l.Info("fetching FunctionTypeList")
			if err := r.List(ctx, &functiontypes,
				&client.ListOptions{Namespace: functionChain.Spec.FunctionTypeNamespace}); err != nil {
				return ctrl.Result{}, err
			}

			// Get and put FunctionType in Status of DataFlow
			for _, fcfunctions := range functionChain.Spec.Functions {
				for _, fks := range functiontypes.Items {
					if fcfunctions.FunctionName == fks.Spec.FunctionName &&
						fcfunctions.Version == fks.Spec.Version {
						getfks := fks
						dataflow.Status.FunctionType = append(dataflow.Status.FunctionType, &getfks)
						l.Info("Put FunctionType: " + fks.Name + "  in Status of DataFlow")
					}
				}
			}

			dataflow.Status.Status = "Scheduling in progress"
			l.Info("Update DataFlow status to Scheduling in progress")
			if _, err := r.patchApply2(ctx, req, dataflow); err != nil {
				l.Error(err, "unable to update df status")
				return ctrl.Result{}, err
			}
		} else {
			// l.Error(err, "FunctionChain STATUS is NotReady")
			l.Error(err, fmt.Sprintf("FunctionChain was found but not Ready state: %s(%s)",
				functionChain.Name, functionChain.Status.Status))
			r.Recorder.Eventf(&dataflow,
				corev1.EventTypeNormal,
				"Reconcile",
				"FunctionChain was found but not Ready state: %s(%s)",
				functionChain.Name, functionChain.Status.Status)
			dataflow.Status.Status = "Error"
			r.patchApply1(ctx, req, dataflow)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

	case "WBFunction/WBConnection creation in progress":
		// Resource creation logic

		var functiontargetList ntthpcv1.FunctionTargetList
		l.Info("fetching FunctionTargetList")
		if err := r.List(ctx, &functiontargetList, &client.ListOptions{}); err != nil {
			return ctrl.Result{}, err
		}
		ftmap := map[string]ntthpcv1.FunctionTarget{}
		for _, functiontarget := range functiontargetList.Items {
			ftv := functiontarget
			ftmap[functiontarget.ObjectMeta.Name] = ftv
		}

		wbfmap := map[string]*ntthpcv1.WBFunction{}
		wbcmap := map[string]*ntthpcv1.WBConnection{}

		// create wbfunction data
		for dfFcFuncKey, dfFcFuncValue := range dataflow.Status.FunctionChain.Spec.Functions {
			wbfunctionName := dataflow.Name + "-wbfunction-" + dfFcFuncKey
			wbfunction := &ntthpcv1.WBFunction{
				ObjectMeta: metav1.ObjectMeta{
					Name:      wbfunctionName,
					Namespace: req.Namespace,
				},
			}

			wbfmap[dfFcFuncKey] = wbfunction

			// get DataFlowName Info
			if wbfunction.Spec.DataFlowRef.Name != dataflow.Name {
				wbfunction.Spec.DataFlowRef.Name = dataflow.Name
			}
			if wbfunction.Spec.DataFlowRef.Namespace != dataflow.Namespace {
				wbfunction.Spec.DataFlowRef.Namespace = dataflow.Namespace
			}

			// FunctionScheduleInfo
			for dfFuncScheInfoKey, dfFuncScheInfoValue := range dataflow.Status.ScheduledFunctions {
				if dfFuncScheInfoKey == dfFcFuncKey {
					// get DeviceKind Info
					if wbfunction.Spec.DeviceType != dfFuncScheInfoValue.DeviceType {
						wbfunction.Spec.DeviceType = dfFuncScheInfoValue.DeviceType
					}

					// get DeviceIndex Info
					if wbfunction.Spec.DeviceIndex != int32(dfFuncScheInfoValue.DeviceIndex) {
						wbfunction.Spec.DeviceIndex = int32(dfFuncScheInfoValue.DeviceIndex)
					}

					// get Region Info
					if wbfunction.Spec.RegionName != dfFuncScheInfoValue.RegionName {
						wbfunction.Spec.RegionName = dfFuncScheInfoValue.RegionName
					}

					// get FunctionIndex Info
					if wbfunction.Spec.FunctionIndex != dfFuncScheInfoValue.FunctionIndex {
						wbfunction.Spec.FunctionIndex = dfFuncScheInfoValue.FunctionIndex
					}

					// get Node Info
					if wbfunction.Spec.NodeName != dfFuncScheInfoValue.NodeName {
						wbfunction.Spec.NodeName = dfFuncScheInfoValue.NodeName
					}
				}
			}

			// get Function Info
			if wbfunction.Spec.FunctionName != dfFcFuncValue.FunctionName {
				wbfunction.Spec.FunctionName = dfFcFuncValue.FunctionName
			}

			// set ConfigName
			var FunctionInfoName string
			var FunctionInfoNamespace string

			for _, dfFuncKindValue := range dataflow.Status.FunctionType {
				if dfFuncKindValue.Spec.FunctionName == dfFcFuncValue.FunctionName &&
					dfFuncKindValue.Spec.Version == dfFcFuncValue.Version {
					FunctionInfoName = dfFuncKindValue.Spec.FunctionInfoCMRef.Name
					FunctionInfoNamespace = dfFuncKindValue.Spec.FunctionInfoCMRef.Namespace
				}
			}

			// get Functioninfo
			funcInfo := &corev1.ConfigMap{}
			if err := r.Get(ctx, client.ObjectKey{
				Name:      FunctionInfoName,
				Namespace: FunctionInfoNamespace,
			}, funcInfo); err != nil {
				l.Error(err, "unable to fetch configmap")
				r.Recorder.Eventf(&dataflow,
					corev1.EventTypeNormal,
					"Reconcile",
					"unable to fetch configmap")
				dataflow.Status.Status = "Error"
				r.patchApply1(ctx, req, dataflow)
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			ftk := wbfunction.Spec.NodeName + "." + wbfunction.Spec.DeviceType + "-" +
				strconv.FormatInt(int64(wbfunction.Spec.DeviceIndex), 10) + "." + wbfunction.Spec.RegionName
			rk := ftmap[ftk].Status.RegionType
			funcInfoJson := funcInfo.Data

			// Get the type (connection method) value of ScheduledConnections for which the function is From or To
			var inputInterfaceType string
			var outputInterfaceType string
			for _, dfConScheValue := range dataflow.Status.ScheduledConnections {
				if strings.HasPrefix(dfConScheValue.To.FunctionKey, endPointKeyword) ||
					strings.HasPrefix(dfConScheValue.From.FunctionKey, startPointKeyword) {
					continue
				}
				if dfConScheValue.To.FunctionKey == dfFcFuncKey && dfConScheValue.To.InterfaceType != nil {
					inputInterfaceType = *dfConScheValue.To.InterfaceType
				}
				if dfConScheValue.From.FunctionKey == dfFcFuncKey && dfConScheValue.From.InterfaceType != nil {
					outputInterfaceType = *dfConScheValue.From.InterfaceType
				}
				if inputInterfaceType != "" && outputInterfaceType != "" {
					l.Info("funckey:" + dfFcFuncKey + " in:" + inputInterfaceType + " out:" + outputInterfaceType)
					break
				}
			}

			var result []map[string]string
			var result2 []map[string]interface{}
			json.Unmarshal([]byte(funcInfoJson["deployableItems"]), &result)
			json.Unmarshal([]byte(funcInfoJson["spec"]), &result2)
			for _, jsonv := range result {
				if rk == jsonv["regionType"] &&
					(inputInterfaceType == jsonv["inputInterfaceType"] || inputInterfaceType == "") &&
					(outputInterfaceType == jsonv["outputInterfaceType"] || outputInterfaceType == "") {
					configName := jsonv["configName"]
					wbfunction.Spec.ConfigName = strings.ReplaceAll(configName, "_", "-")
					specName := jsonv["specName"]
					if wbfunction.Spec.FunctionIndex == nil {
						for _, jsonv2 := range result2 {
							if specName == jsonv2["name"] {
								if _, exist := jsonv2["maxDataFlowsBase"]; exist {
									maxDataFlowsBase := int32(jsonv2["maxDataFlowsBase"].(float64))
									wbfunction.Spec.MaxDataFlows = &maxDataFlowsBase
								}
								if _, exist := jsonv2["maxCapacityBase"]; exist {
									maxCapacityBase := int32(jsonv2["maxCapacityBase"].(float64))
									wbfunction.Spec.MaxCapacity = &maxCapacityBase
								}
								break
							}
						}
					}
					break
				}
			}

			// Get the Port and InterfaceType values of the ScheduledConnections for which the current function is From or To
			inputInterface := make(map[string]string)
			outputInterface := make(map[string]string)
			for _, dfConScheValue := range dataflow.Status.ScheduledConnections {
				if dfConScheValue.To.FunctionKey == dfFcFuncKey &&
					dfConScheValue.To.Port != nil &&
					dfConScheValue.To.InterfaceType != nil {
					inputInterface[strconv.Itoa(int(*dfConScheValue.To.Port))] = *dfConScheValue.To.InterfaceType
				}
				if dfConScheValue.From.FunctionKey == dfFcFuncKey &&
					dfConScheValue.From.Port != nil &&
					dfConScheValue.From.InterfaceType != nil {
					outputInterface[strconv.Itoa(int(*dfConScheValue.From.Port))] = *dfConScheValue.From.InterfaceType
				}
			}
			// Set it to WBFunction
			if len(inputInterface) != 0 {
				wbfunction.Spec.InputInterface = inputInterface
			}
			if len(outputInterface) != 0 {
				wbfunction.Spec.OutputInterface = outputInterface
			}

			// TODO : The following parameters are optional and have low priority
			// if ioName == "host-mem" {
			//		m := map[string]string{"0": "host-mem"}
			//		wbfunction.Spec.InputInterface = m
			//		wbfunction.Spec.OutputInterface = m
			//	}
			// get InputInterface
			// get OutputInterface

			// set custom parameter
			for fcCustomKey, fcCustomValue := range dfFcFuncValue.CustomParameter {
				if wbfunction.Spec.Params == nil {
					wbfunction.Spec.Params = make(map[string]intstr.IntOrString)
				}
				wbfunction.Spec.Params[fcCustomKey] = fcCustomValue
			}
			for _, userParamStruct := range dataflow.Spec.FunctionUserParameter {
				if userParamStruct.FunctionKey == dfFcFuncKey {
					if wbfunction.Spec.Params == nil {
						wbfunction.Spec.Params = make(map[string]intstr.IntOrString)
					}
					for k, v := range userParamStruct.UserParams {
						wbfunction.Spec.Params[k] = v
					}
					break
				}
			}

			// get Requirements
			if dataflow.Spec.Requirements != nil && dataflow.Spec.Requirements.All != nil {
				wbfunction.Spec.Requirements = &ntthpcv1.WBFunctionRequirementsInfo{}
				wbfunction.Spec.Requirements.Capacity = dataflow.Spec.Requirements.All.Capacity
			}

			// create wbconnection data
			for _, dfFcConValue := range dataflow.Status.FunctionChain.Spec.Connections {

				// set connection direction Input or Output
				var direction string
				if dfFcConValue.From.FunctionKey == dfFcFuncKey {
					direction = "Output"
				} else if dfFcConValue.To.FunctionKey == dfFcFuncKey {
					direction = "Input"
				} else {
					continue
				}

				wbconnectionName := dataflow.Name + "-wbconnection-" + dfFcConValue.From.FunctionKey +
					"-" + dfFcConValue.To.FunctionKey
				wbconnection, exists := wbcmap[wbconnectionName]
				if !exists {
					wbconnection = &ntthpcv1.WBConnection{
						ObjectMeta: metav1.ObjectMeta{
							Name:      wbconnectionName,
							Namespace: req.Namespace,
						},
					}
					wbcmap[wbconnectionName] = wbconnection
				}

				l.Info(fmt.Sprintf("Create for WBConnection Object"))

				// get dataflowname Info
				if wbconnection.Spec.DataFlowRef.Name != dataflow.Name {
					wbconnection.Spec.DataFlowRef.Name = dataflow.Name
				}
				if wbconnection.Spec.DataFlowRef.Namespace != dataflow.Namespace {
					wbconnection.Spec.DataFlowRef.Namespace = dataflow.Namespace
				}

				// ScheduledConnections
				for _, dfConScheInfoValue := range dataflow.Status.ScheduledConnections {
					if (direction == "Output" && dfConScheInfoValue.From.FunctionKey == dfFcFuncKey) ||
						(direction == "Input" && dfConScheInfoValue.To.FunctionKey == dfFcFuncKey) {
						// get Type Info
						if wbconnection.Spec.ConnectionMethod != dfConScheInfoValue.ConnectionMethod {
							wbconnection.Spec.ConnectionMethod = dfConScheInfoValue.ConnectionMethod
						}
					}
				}

				// get FromPort
				if dfFcConValue.From.Port != 0 {
					if (wbconnection.Spec.From.Port) != dfFcConValue.From.Port {
						wbconnection.Spec.From.Port = dfFcConValue.From.Port
					}
				}
				// get ToPort
				if dfFcConValue.To.Port != 0 {
					if wbconnection.Spec.To.Port != dfFcConValue.To.Port {
						wbconnection.Spec.To.Port = dfFcConValue.To.Port
					}
				}

				// get FromFunction
				if dfFcFuncKey == dfFcConValue.From.FunctionKey {
					if wbconnection.Spec.From.WBFunctionRef.Name != wbfunctionName {
						wbconnection.Spec.From.WBFunctionRef.Name = wbfunctionName
					}
					if wbconnection.Spec.From.WBFunctionRef.Namespace != wbfunction.Namespace {
						wbconnection.Spec.From.WBFunctionRef.Namespace = wbfunction.Namespace
					}
				}
				// get ToFunction
				if dfFcFuncKey == dfFcConValue.To.FunctionKey {
					if wbconnection.Spec.To.WBFunctionRef.Name != wbfunctionName {
						wbconnection.Spec.To.WBFunctionRef.Name = wbfunctionName
					}
					if wbconnection.Spec.To.WBFunctionRef.Namespace != wbfunction.Namespace {
						wbconnection.Spec.To.WBFunctionRef.Namespace = wbfunction.Namespace
					}
				}

				// set custom parameter
				if dfFcFuncKey == dfFcConValue.From.FunctionKey {
					if dfFcConValue.CustomParameter != nil {
						if wbconnection.Spec.Params == nil {
							wbconnection.Spec.Params = make(map[string]intstr.IntOrString)
						}
						for fcCustomKey, fcCustomValue := range dfFcConValue.CustomParameter {
							wbconnection.Spec.Params[fcCustomKey] = fcCustomValue
						}
					}
					if dataflow.Spec.ConnectionUserParameter != nil {
						for _, userParamStruct := range dataflow.Spec.ConnectionUserParameter {
							if dfFcFuncKey == userParamStruct.From.FunctionKey {
								if wbconnection.Spec.Params == nil {
									wbconnection.Spec.Params = make(map[string]intstr.IntOrString)
								}
								for k, v := range userParamStruct.UserParams {
									wbconnection.Spec.Params[k] = v
								}
								break
							}
						}
					}
				}
				if dfFcFuncKey == dfFcConValue.To.FunctionKey {
					if dfFcConValue.CustomParameter != nil {
						if wbconnection.Spec.Params == nil {
							wbconnection.Spec.Params = make(map[string]intstr.IntOrString)
						}
						for fcCustomKey, fcCustomValue := range dfFcConValue.CustomParameter {
							wbconnection.Spec.Params[fcCustomKey] = fcCustomValue
						}
					}
					if dataflow.Spec.ConnectionUserParameter != nil {
						for _, userParamStruct := range dataflow.Spec.ConnectionUserParameter {
							if dfFcFuncKey == userParamStruct.To.FunctionKey {
								if wbconnection.Spec.Params == nil {
									wbconnection.Spec.Params = make(map[string]intstr.IntOrString)
								}
								for k, v := range userParamStruct.UserParams {
									wbconnection.Spec.Params[k] = v
								}
								break
							}
						}
					}
				}

				// For start and end points, set network information
				if strings.HasPrefix(dfFcConValue.From.FunctionKey, startPointKeyword) {
					wbconnection.Spec.From.WBFunctionRef.Name = dfFcConValue.From.FunctionKey
					wbconnection.Spec.From.WBFunctionRef.Namespace = dataflow.Namespace
					if dataflow.Spec.StartPoint != nil && dataflow.Spec.StartPoint.IP != "" {
						if wbconnection.Spec.Params == nil {
							wbconnection.Spec.Params = make(map[string]intstr.IntOrString)
						}
						wbconnection.Spec.Params["TargetPort"] = intstr.IntOrString{Type: intstr.Int,
							IntVal: int32(dataflow.Spec.StartPoint.Port)}
						wbconnection.Spec.Params["TargetIP"] = intstr.IntOrString{Type: intstr.String,
							StrVal: dataflow.Spec.StartPoint.IP}
						wbconnection.Spec.Params["Protocol"] = intstr.IntOrString{Type: intstr.String,
							StrVal: string(dataflow.Spec.StartPoint.Protocol)}
					}
				} else if strings.HasPrefix(dfFcConValue.To.FunctionKey, endPointKeyword) {
					wbconnection.Spec.To.WBFunctionRef.Name = dfFcConValue.To.FunctionKey
					wbconnection.Spec.To.WBFunctionRef.Namespace = dataflow.Namespace
					if dataflow.Spec.EndPoint != nil && dataflow.Spec.EndPoint.IP != "" {
						if wbconnection.Spec.Params == nil {
							wbconnection.Spec.Params = make(map[string]intstr.IntOrString)
						}
						wbconnection.Spec.Params["TargetPort"] = intstr.IntOrString{Type: intstr.Int,
							IntVal: int32(dataflow.Spec.EndPoint.Port)}
						wbconnection.Spec.Params["TargetIP"] = intstr.IntOrString{Type: intstr.String,
							StrVal: dataflow.Spec.EndPoint.IP}
						wbconnection.Spec.Params["Protocol"] = intstr.IntOrString{Type: intstr.String,
							StrVal: string(dataflow.Spec.EndPoint.Protocol)}
					}
				}

				// get ScheduledConnections
				for _, dfConScheValue := range dataflow.Status.ScheduledConnections {
					// if dataflow has coressponding route, set it and CapacityUsed to WBConnection
					if dfConScheValue.ConnectionPath != nil {
						// get dfFcFuncKey from FromFunction.Name
						var wbConFromFcFuncKey string
						if !strings.HasPrefix(wbconnection.Spec.From.WBFunctionRef.Name, startPointKeyword) &&
							len(wbconnection.Spec.From.WBFunctionRef.Name) != 0 {
							wbConFromFcFuncKey =
								strings.Split(wbconnection.Spec.From.WBFunctionRef.Name, "-wbfunction-")[1]
						}
						// get dfFcFuncKey from ToFunction.Name
						var wbConToFcFuncKey string
						if !strings.HasPrefix(wbconnection.Spec.To.WBFunctionRef.Name, endPointKeyword) &&
							len(wbconnection.Spec.To.WBFunctionRef.Name) != 0 {
							wbConToFcFuncKey =
								strings.Split(wbconnection.Spec.To.WBFunctionRef.Name, "-wbfunction-")[1]
						}
						if (strings.HasPrefix(wbconnection.Spec.From.WBFunctionRef.Name, startPointKeyword) &&
							wbConToFcFuncKey == dfConScheValue.To.FunctionKey) ||
							(wbConFromFcFuncKey == dfConScheValue.From.FunctionKey &&
								wbConToFcFuncKey == dfConScheValue.To.FunctionKey) ||
							(wbConFromFcFuncKey == dfConScheValue.From.FunctionKey &&
								strings.HasPrefix(wbconnection.Spec.To.WBFunctionRef.Name, endPointKeyword)) {
							// get Route
							wbconnection.Spec.ConnectionPath = dfConScheValue.ConnectionPath
							// get CapacityUsed
							if dataflow.Spec.Requirements != nil && dataflow.Spec.Requirements.All != nil {
								wbconnection.Spec.Requirements = &ntthpcv1.WBConnectionRequirementsInfo{}
								wbconnection.Spec.Requirements.Capacity = dataflow.Spec.Requirements.All.Capacity
							}
						}
					}
				}
			}
		}

		// set previous and next
		for _, conn := range dataflow.Status.FunctionChain.Spec.Connections {
			fromWBFunc := wbfmap[conn.From.FunctionKey]
			toWBFunc := wbfmap[conn.To.FunctionKey]
			if fromWBFunc == nil || toWBFunc == nil {
				continue
			}

			if fromWBFunc.Spec.NextWBFunctions == nil {
				fromWBFunc.Spec.NextWBFunctions = map[string]ntthpcv1.FromToWBFunction{}
			}
			fromPort := int32(0)
			if conn.From.Port != 0 {
				fromPort = conn.From.Port
			}
			fromWBFunc.Spec.NextWBFunctions[strconv.Itoa(int(fromPort))] = ntthpcv1.FromToWBFunction{
				WBFunctionRef: ntthpcv1.WBNamespacedName{Namespace: toWBFunc.Namespace, Name: toWBFunc.Name},
				Port:          conn.To.Port,
			}

			if toWBFunc.Spec.PreviousWBFunctions == nil {
				toWBFunc.Spec.PreviousWBFunctions = map[string]ntthpcv1.FromToWBFunction{}
			}
			toPort := int32(0)
			if conn.To.Port != 0 {
				toPort = conn.To.Port
			}
			toWBFunc.Spec.PreviousWBFunctions[strconv.Itoa(int(toPort))] = ntthpcv1.FromToWBFunction{
				WBFunctionRef: ntthpcv1.WBNamespacedName{Namespace: fromWBFunc.Namespace, Name: fromWBFunc.Name},
				Port:          conn.From.Port,
			}
		}

		// Create wbfunction resources
		for _, wbfunction := range wbfmap {
			if _, err := ctrl.CreateOrUpdate(ctx, r.Client, wbfunction, func() error {
				l.Info(fmt.Sprintf("CreateOrUpdate for WBFunction Resource"))
				// set the owner so that garbage collection can kicks in
				if err := ctrl.SetControllerReference(&dataflow, wbfunction, r.Scheme); err != nil {
					l.Error(err, "unable to set ownerReference from DataFlow to WBFunction")
					return err
				}

				// end of ctrl.CreateOrUpdate
				return nil

			}); err != nil {
				// error handling of ctrl.CreateOrUpdate
				l.Error(err, "unable to ensure wbfunction is correct")
				return ctrl.Result{Requeue: true}, nil
			}
		}

		// Create wbconnection resources
		for _, wbconnection := range wbcmap {
			if _, err := ctrl.CreateOrUpdate(ctx, r.Client, wbconnection, func() error {
				l.Info(fmt.Sprintf("CreateOrUpdate for WBConnection Resource"))
				// set the owner so that garbage collection can kicks in
				if err := ctrl.SetControllerReference(&dataflow, wbconnection, r.Scheme); err != nil {
					l.Error(err, "unable to set ownerReference from DataFlow to WBConnection")
					return err
				}
				return nil
			}); err != nil {
				// error handling of ctrl.CreateOrUpdate
				l.Error(err, "unable to ensure wbconnection is correct")
				return ctrl.Result{Requeue: true}, nil
			}
		}

		dataflow.Status.Status = "WBFunction/WBConnection created"
		l.Info("Update DataFlow status to WBFunction/WBConnection created")
		if _, err := r.patchApply1(ctx, req, dataflow); err != nil {
			l.Error(err, "unable to update df status")
			return ctrl.Result{}, err
		}

	case "WBFunction/WBConnection created":
		l.Info("check WBfunction/WBConnection")
		var wbfunctionList ntthpcv1.WBFunctionList
		l.Info("fetching wbfunctionList")
		if err := r.List(ctx, &wbfunctionList, &client.ListOptions{}); err != nil {
			return ctrl.Result{}, err
		}
		for _, wbfunction := range wbfunctionList.Items {
			if wbfunction.Spec.DataFlowRef.Name == dataflow.Name {
				if wbfunction.Status.Status != "Deployed" {
					return ctrl.Result{Requeue: true}, nil
				}
			}
		}
		var wbconnectionList ntthpcv1.WBConnectionList
		l.Info("fetching wbconnectionList")
		if err := r.List(ctx, &wbconnectionList, &client.ListOptions{}); err != nil {
			return ctrl.Result{}, err
		}
		for _, wbconnection := range wbconnectionList.Items {
			if wbconnection.Spec.DataFlowRef.Name == dataflow.Name {
				if wbconnection.Status.Status != "Deployed" {
					return ctrl.Result{Requeue: true}, nil
				}
			}
		}

		dataflow.Status.Status = "Deployed"
		l.Info("Update DataFlow status to Deployed")
		if _, err := r.patchApply1(ctx, req, dataflow); err != nil {
			l.Error(err, "unable to update df status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func fcToJsonTagMap(fc *ntthpcv1.FunctionChain) map[string]interface{} {

	result := map[string]interface{}{}

	b, _ := json.Marshal(fc)
	json.Unmarshal(b, &result)

	return result
}

func ftListToJsonTagSlice(ftList []*ntthpcv1.FunctionType) []map[string]interface{} {
	result := []map[string]interface{}{}

	b, _ := json.Marshal(ftList)
	json.Unmarshal(b, &result)

	return result
}

func ctListToJsonTagSlice(ctList []*ntthpcv1.ConnectionType) []map[string]interface{} {
	result := []map[string]interface{}{}

	b, _ := json.Marshal(ctList)
	json.Unmarshal(b, &result)

	return result
}

func (r *DataFlowReconciler) patchApply1(ctx context.Context, req ctrl.Request,
	df ntthpcv1.DataFlow) (ctrl.Result, error) {
	patch := &unstructured.Unstructured{}
	patch.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   ntthpcv1.GroupVersion.Group,
		Version: ntthpcv1.GroupVersion.Version,
		Kind:    "DataFlow",
	})
	patch.SetNamespace(df.Namespace)
	patch.SetName(df.Name)
	patch.UnstructuredContent()["status"] = map[string]interface{}{
		"status": df.Status.Status,
	}
	srpOpts := &client.SubResourcePatchOptions{}
	srpOpts.FieldManager = "dataflow-controller"
	srpOpts.Force = pointer.Bool(true)
	err := r.Status().Patch(ctx, patch, client.Apply, srpOpts)

	return ctrl.Result{}, err
}

func (r *DataFlowReconciler) patchApply2(ctx context.Context, req ctrl.Request,
	df ntthpcv1.DataFlow) (ctrl.Result, error) {
	patch := &unstructured.Unstructured{}
	patch.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   ntthpcv1.GroupVersion.Group,
		Version: ntthpcv1.GroupVersion.Version,
		Kind:    "DataFlow",
	})
	patch.SetNamespace(df.Namespace)
	patch.SetName(df.Name)

	fc := fcToJsonTagMap(df.Status.FunctionChain)
	ft := ftListToJsonTagSlice(df.Status.FunctionType)
	ct := ctListToJsonTagSlice(df.Status.ConnectionType)

	patch.UnstructuredContent()["status"] = map[string]interface{}{
		"status":         df.Status.Status,
		"functionChain":  fc,
		"functionType":   ft,
		"connectionType": ct,
	}
	srpOpts := &client.SubResourcePatchOptions{}
	srpOpts.FieldManager = "dataflow-controller"
	srpOpts.Force = pointer.Bool(true)
	err := r.Status().Patch(ctx, patch, client.Apply, srpOpts)

	return ctrl.Result{}, err
}

func (r *DataFlowReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.DataFlow{}).
		Owns(&ntthpcv1.WBFunction{}).
		Owns(&ntthpcv1.WBConnection{}).
		Complete(r)
}
