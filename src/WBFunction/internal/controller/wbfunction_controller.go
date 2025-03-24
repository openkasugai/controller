/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "WBFunction/api/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	// Additional imports
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// WBFunctionReconciler reconciles a WBFunction object
type WBFunctionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// Event type
const (
	CREATE = iota
	UPDATE
	DELETE
)

// constructor
const (
	PENDING = "Pending" // WBFunction.status.status(Pending)
	RUNNING = "Running" // WBFunction.status.status(Running)
)

const (
	DEFAULTSECOND      = 10  // Default requeue interval (s)
	DEFAULTMICROSECOND = 500 // Default requeue interval (ms)
)

const (
	MYCRKIND         = "WBFunction"
	DEVICEINFOCRKIND = "DeviceInfo"
	DEVICEINFOCRNAME = "deviceinfo-"
)

type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   examplecomv1.FunctionData       `json:"spec,omitempty"`
	Status examplecomv1.FunctionStatusData `json:"status,omitempty"`
}

//+kubebuilder:rbac:groups=example.com,resources=wbfunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=wbfunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=wbfunctions/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=deviceinfoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=deviceinfoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=deviceinfoes/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=cpufunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=cpufunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=cpufunctions/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=gpufunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=gpufunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=gpufunctions/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WBFunction object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *WBFunctionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// WBFunction CR Get
	var functionKind string
	var functionAPIVersion string
	var requeueFlag bool
	var err error

	functionKind = ""
	functionAPIVersion = ""
	requeueFlag = false

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger.Info("Reconcile start.")

		// WBFunction CR Get
		var crData examplecomv1.WBFunction
		err = r.getWBFunctionCR(ctx, req, &crData)
		if errors.IsNotFound(err) {
			// If WBFunction CR does not exist
			break
		} else if nil != err {
			// If an error occurs in the Get function
			break
		}
		// Search for the kind to create a CR for
		ret := convTypeMeta(ctx, &crData, &functionKind, &functionAPIVersion)
		if false == ret {
			break
		}

		// Get Event type
		var eventKind int32 // 0:Add, 1:Upd,  2:Del
		eventKind = getEventKind(ctx, &crData)

		// Get the kind and apiVersion of the created CR
		var eventFunctionKind int32 // 0:Add, 1:Upd,  2:Del
		var crFunc Function
		eventFunctionKind = UPDATE

		// Function CR acquisition
		err = r.getFunctionCR(ctx,
			req,
			functionKind,
			functionAPIVersion,
			&crFunc)
		if errors.IsNotFound(err) {
			// If Function CR does not exist
			eventFunctionKind = CREATE
		} else if nil != err {
			// If an error occurs in the Get function
			requeueFlag = true
			break
		}

		var crDeviceInfoData examplecomv1.DeviceInfo
		var reqKind string
		reqKind = examplecomv1.RequestDeploy

		if crData.Status.Status == examplecomv1.WBDeployStatusDeployed &&
			UPDATE == eventKind {
			// Measures for reconcile operation when updating status
			logger.Info("Reconcile end.")
			break
		}
		if CREATE == eventKind {
			for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

				// In case of creation
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Create", "Create Start")
				if CREATE == eventFunctionKind {
					// Create a resource
					err = r.createCustomResource(ctx, req,
						functionKind, functionAPIVersion, &crData)
					if nil != err {
						break
					}
					logger.Info("Status Information Change start.")
					err = r.updCustomResource(ctx,
						&crData,
						examplecomv1.WBDeployStatusWaiting,
						eventKind)
					if nil != err {
						logger.Error(err,
							"failed to update WBFunction status.")
						fmt.Printf("%#v\n", crData)
						break
					}
					logger.Info("Status Information Change end.")
				}
				if examplecomv1.WBDeployStatusDeployed != crData.Status.Status &&
					examplecomv1.WBDeployStatusWaiting != crData.Status.Status {
					requeueFlag = true
				}

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Create", "Create End")
				break //nolint:staticcheck // SA4004: Intentional break
			}
		} else if UPDATE == eventKind {

			for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

				// In case of update
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Update", "Update Start")

				logger.Info("eventFunctionKind=" +
					strconv.Itoa(int(eventFunctionKind)))
				logger.Info("crData.Status.Status=" +
					string(crData.Status.Status))
				logger.Info("crFunc.status.Status=" + crFunc.Status.Status)

				if UPDATE == eventFunctionKind {
					if examplecomv1.WBDeployStatusDeployed != crData.Status.Status {
						if RUNNING == crFunc.Status.Status {
							// Store spec information in status
							err = setSpecToStatusData(ctx,
								&crData,
								crFunc.Status.FunctionIndex)
							if nil != err {
								break
							}
							// Get DeviceInfo CR
							err = r.getDevMngCR(ctx, &crData, &crDeviceInfoData, req)
							if errors.IsNotFound(err) {
								_ = r.createDevMngCR(ctx, // FIXME: need to handle return value (error)
									&crData,
									&crDeviceInfoData,
									req,
									reqKind)
							} else if nil != err {
								// If an error occurs in the Get function
								break
							}
							if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceError {
								err = r.updCustomResource(ctx, &crData,
									examplecomv1.WBDeployStatusFailed,
									eventKind)
								if nil != err {
									logger.Error(err, "failed to update WBFunction status.")
									fmt.Printf("%#v\n", crData)
									break
								}
								logger.Error(fmt.Errorf("DeviceInfo process Status error."),
									"Reconcile process abnormal end.")
								break
							}

							if crDeviceInfoData.Status.Response.Status != examplecomv1.ResponceDeployed {
								requeueFlag = true
								break
							}

							// Processing is complete, so update the status
							logger.Info("Status Running Change start.")
							err = r.updCustomResource(ctx, &crData,
								examplecomv1.WBDeployStatusDeployed,
								eventKind)
							if nil != err {
								break
							}
							err = r.deleteDevMngCR(ctx, &crDeviceInfoData)
							if nil != err {
								break
							}
							logger.Info("Status Running Change end.")
						}
					}
				}
				if examplecomv1.WBDeployStatusDeployed != crData.Status.Status {
					requeueFlag = true
				}

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Update", "Update End")
				break //nolint:staticcheck // SA4004: Intentional break
			}

		} else if DELETE == eventKind {
			for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Delete", "Delete Start")
				if crData.Status.Status == examplecomv1.WBDeployStatusDeployed {
					if eventFunctionKind != CREATE {
						// Delete FunctionCR
						if crFunc.DeletionTimestamp.IsZero() {
							logger.Info("Delete FunctionCR start.")
							err = r.delFunctionCR(ctx, functionKind, functionAPIVersion, req)
							if err != nil {
								logger.Info("Failed to delete FunctionCR.")
								break
							}
						}
					}
					logger.Info("Update WBFunction to Terminating.")
					err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusTerminating, eventKind)
					if err != nil {
						logger.Info("Failed to update WBFunction.")
						break
					}
					logger.Info("Delete FunctionCR was end.")
					logger.Info("WBFunction change to Terminating was end.")
					break
				} else if crData.Status.Status == examplecomv1.WBDeployStatusTerminating {
					if eventFunctionKind != CREATE {
						requeueFlag = true
						break
					}
				}

				if crData.Status.Status == examplecomv1.WBDeployStatusTerminating ||
					crData.Status.Status == examplecomv1.WBDeployStatusReleased {
					// Get DeviceInfo CR
					err = r.getDevMngCR(ctx, &crData, &crDeviceInfoData, req)
					if errors.IsNotFound(err) {

						if crData.Status.Status == examplecomv1.WBDeployStatusTerminating {
							logger.Info("Create DeviceInfoCR.")
							reqKind = examplecomv1.RequestUndeploy
							_ = r.createDevMngCR(ctx,
								&crData,
								&crDeviceInfoData,
								req,
								reqKind)
							requeueFlag = true

						} else if crData.Status.Status == examplecomv1.WBDeployStatusReleased {
							logger.Info("Delete WBFunction.")
							err = r.delCustomResource(ctx, &crData)
							if err != nil {
								logger.Info("Failed to delete WBFunction")
							} else {
								logger.Info("Delete WBFunctionCR was end.")
							}

						}
						break
					}
				}

				if crData.Status.Status == examplecomv1.WBDeployStatusTerminating {

					if crDeviceInfoData.DeletionTimestamp.IsZero() {

						if crDeviceInfoData.Status.Response.Status == "" &&
							crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestDeploy {
							requeueFlag = true
							break

						} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceError &&
							crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestDeploy {
							logger.Info("Delete DeviceInfo start.")
							err = r.deleteDevMngCR(ctx, &crDeviceInfoData)
							if err != nil {
								logger.Info("Failed to delete DeviceInfo.")
							} else {
								logger.Info("Delete DeviceInfo was end.")
							}
							logger.Info("Update WBFunctions to Released.")
							err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusReleased, eventKind)
							if err != nil {
								logger.Info("Failed to update WBFunction.")
							} else {
								logger.Info("WBFunction change to Released was end.")
							}

						} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceDeployed {
							logger.Info("Delete DeviceInfo start.")
							err = r.deleteDevMngCR(ctx, &crDeviceInfoData)
							if err != nil {
								logger.Info("Failed to delete DeviceInfo.")
							} else {
								logger.Info("Delete DeviceInfo was end.")
							}
							requeueFlag = true

						} else if crDeviceInfoData.Status.Response.Status == "" &&
							crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestUndeploy {
							requeueFlag = true
							break

						} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceError &&
							crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestUndeploy {
							logger.Info("Update WBFunctions to Failed.")
							err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusFailed, eventKind)
							if err != nil {
								logger.Info("Failed to update WBFunction.")
							}
							break

						} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceUndeployed {
							logger.Info("Delete DeviceInfo start.")
							err = r.deleteDevMngCR(ctx, &crDeviceInfoData)
							if err != nil {
								logger.Info("Failed to delete DeviceInfo.")
							} else {
								logger.Info("Delete DeviceInfo was end.")
							}
							logger.Info("Update WBFunctions to Released.")
							err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusReleased, eventKind)
							if err != nil {
								logger.Info("Failed to update WBFunction.")
							} else {
								logger.Info("WBFunction change to Released was end.")
							}

						} else {
							logger.Info("Unexpected route: WBFunctionCRStatus: " +
								string(crData.Status.Status) +
								", DeviceInfoCRRequestType: " +
								string(crDeviceInfoData.Spec.Request.RequestType) +
								", DeletionTimeStamp: DeletionTimeStamp is nil,  DeviceInfoCRResponceStatus: " +
								string(crDeviceInfoData.Status.Response.Status))
							break
						}

					} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceDeployed {
						requeueFlag = true
						break

					} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceUndeployed {
						logger.Info("Update WBFunctions to Released.")
						err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusReleased, eventKind)
						if err != nil {
							logger.Info("Failed to update WBFunction.")
						} else {
							logger.Info("WBFunction change to Released was end.")
						}

					} else if crDeviceInfoData.Status.Response.Status == "" &&
						crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestDeploy {
						requeueFlag = true
						break

					} else if crDeviceInfoData.Status.Response.Status == "" &&
						crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestUndeploy {
						logger.Info("Update WBFunctions to Released.")
						err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusReleased, eventKind)
						if err != nil {
							logger.Info("Failed to update WBFunction.")
						} else {
							logger.Info("WBFunction change to Released was end.")
						}

					} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceError &&
						crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestDeploy {
						logger.Info("Update WBFunctions to Released.")
						err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusReleased, eventKind)
						if err != nil {
							logger.Info("Failed to update WBFunction.")
						} else {
							logger.Info("WBFunction change to Released was end.")
						}

					} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceError &&
						crDeviceInfoData.Spec.Request.RequestType == examplecomv1.RequestUndeploy {
						logger.Info("Update WBFunctions to Failed.")
						err = r.updCustomResource(ctx, &crData, examplecomv1.WBDeployStatusFailed, eventKind)
						if err != nil {
							logger.Info("Failed to update WBFunction.")
						}
						break
					}

				} else if crData.Status.Status == examplecomv1.WBDeployStatusReleased {
					if crDeviceInfoData.DeletionTimestamp.IsZero() {
						logger.Info("Unexpected route: WBFunctionCRStatus: " +
							string(crData.Status.Status) +
							", DeviceInfoCRRequestType: " +
							string(crDeviceInfoData.Spec.Request.RequestType) +
							", DeletionTimeStamp: DeletionTimeStamp is nil,  DeviceInfoCRResponceStatus: " +
							string(crDeviceInfoData.Status.Response.Status))
						break
					} else if crDeviceInfoData.Status.Response.Status == examplecomv1.ResponceUndeployed {
						requeueFlag = true
						break
					} else {
						logger.Info("Unexpected route: WBFunctionCRStatus: " +
							string(crData.Status.Status) +
							", DeletionTimeStamp: " +
							crDeviceInfoData.DeletionTimestamp.Format("2006-01-02T15:04:05Z07:00") +
							", DeviceInfoCRRequestType: " +
							string(crDeviceInfoData.Spec.Request.RequestType) +
							", DeviceInfoCRResponceStatus: " +
							string(crDeviceInfoData.Status.Response.Status))
						break
					}
				} else {
					logger.Info("Unexpected route: WBFunctionCRStatus: " +
						string(crData.Status.Status) +
						", DeletionTimeStamp: " +
						crDeviceInfoData.DeletionTimestamp.Format("2006-01-02T15:04:05Z07:00") +
						", DeviceInfoCRRequestType: " +
						string(crDeviceInfoData.Spec.Request.RequestType) +
						", DeviceInfoCRResponceStatus: " +
						string(crDeviceInfoData.Status.Response.Status))
					break
				}
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Delete", "Delete End")
			}
		}
	}
	if requeueFlag == true {
		return ctrl.Result{Requeue: requeueFlag}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WBFunctionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.WBFunction{}).
		Complete(r)
}

// Copy from spec area to status area
func setSpecToStatusData(
	ctx context.Context,
	pCRData *examplecomv1.WBFunction,
	functionindex *int32) error {

	logger := log.FromContext(ctx)
	var err error = nil // To determine the return value

	for doWhile := 0; doWhile < 1; doWhile++ {
		jsonstr, err := json.Marshal(pCRData.Spec)
		if nil != err {
			logger.Error(err, "WBFunction.Spec unable to Marshal.")
			break
		}
		err = json.Unmarshal(jsonstr, &pCRData.Status)
		if nil != err {
			logger.Error(err, DEVICEINFOCRKIND+" unable to Unmarshal.")
			break
		}
		pCRData.Status.FunctionIndex = *functionindex
	}
	return err
}

const (
	APIVERSION = "example.com/v1"
)

// Create CR attribute information and obtain it
func convTypeMeta(
	ctx context.Context,
	pCRData *examplecomv1.WBFunction,
	pFunctionKind *string,
	pFunctionAPIVersion *string) bool {

	logger := log.FromContext(ctx)

	var ret bool = false

	for count := 0; count < len(gFunctionKindMap); count++ {
		if gFunctionKindMap[count].DeviceType == pCRData.Spec.DeviceType {
			*pFunctionKind = gFunctionKindMap[count].FunctionCRKind
			*pFunctionAPIVersion = APIVERSION
		}
	}
	if *pFunctionAPIVersion == "" || *pFunctionKind == "" {
		logger.Info("Convert Resource Kind Error." +
			" apivesion=" + *pFunctionAPIVersion +
			" kind= " + *pFunctionKind)
	} else {
		ret = true
	}
	return ret
}

// Infrastructure information ConfigMap
type InfrastructureInfo struct {
	DeviceFilePath *string `json:"deviceFilePath,omitempty"`
	NodeName       string  `json:"nodeName"`
	DeviceUUID     *string `json:"deviceUUID,omitempty"`
	DeviceType     string  `json:"deviceType"`
	DeviceIndex    int32   `json:"deviceIndex"`
}

// Config information storage area
var gInfraStructureInfo map[string][]InfrastructureInfo

var config_load_tbl = []ConfigTable{
	{"infrastructureinfo"},
}

// Load ConfigMap (startup function)
func loadConfigMap(ctx context.Context, r *WBFunctionReconciler) error {

	logger := log.FromContext(ctx)
	var cfgdata []byte
	var err error

	for _, record := range config_load_tbl {
		err = r.getConfigMap(ctx, record.name, &cfgdata)
		if nil != err {
			break
		}
		if "infrastructureinfo" == record.name {
			err = json.Unmarshal(cfgdata, &gInfraStructureInfo)
		}
		if nil != err {
			logger.Error(err, "unable to unmarshal. ConfigMap="+record.name)
			break
		}
	}
	return err
}

// Each CRC is created using Resource
func (r *WBFunctionReconciler) createCustomResource(
	ctx context.Context,
	req ctrl.Request,
	functionKind string,
	functionAPIVersion string,
	pCRData *examplecomv1.WBFunction) error {

	logger := log.FromContext(ctx)

	var err error = nil // To determine the return value
	var funccfg examplecomv1.FunctionConfigMap
	var configname string
	var cfgdata []byte
	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		_ = loadConfigMap(ctx, r)

		infraDevice := gInfraStructureInfo["devices"]

		configname = pCRData.Spec.ConfigName
		if functionKind == examplecomv1.FunctionCRKindFPGA {
			err = r.getConfigMap(ctx, configname, &cfgdata)
			if nil != err {
				break
			}
			err = json.Unmarshal(cfgdata, &funccfg)
		} else {
			err = r.getConfigData(ctx, configname, &funccfg)
		}
		if nil != err {
			logger.Error(err, "FunctionConfigMap json.Unmarshal Error.")
			break
		}

		// Organize the cr information to be created
		crData := &unstructured.Unstructured{}
		crData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: functionAPIVersion,
			Kind:    functionKind,
		})
		crData.SetName(req.Name)
		crData.SetNamespace(req.Namespace)
		memorySize := int32(examplecomv1.RequestMemory)

		var crSpec examplecomv1.FunctionData
		crSpec.DataFlowRef = pCRData.Spec.DataFlowRef
		crSpec.FunctionName = pCRData.Spec.FunctionName
		crSpec.NodeName = pCRData.Spec.NodeName
		crSpec.DeviceType = pCRData.Spec.DeviceType
		crSpec.RegionName = pCRData.Spec.RegionName
		if nil != pCRData.Spec.FunctionIndex {
			crSpec.FunctionIndex = pCRData.Spec.FunctionIndex
		}
		crSpec.RequestMemorySize = &memorySize
		crSpec.ConfigName = &pCRData.Spec.ConfigName
		if 0 != funccfg.SharedMemoryMiB {
			crSpec.SharedMemory.FilePrefix = req.Namespace + "-" + req.Name
			crSpec.SharedMemory.CommandQueueID = req.Namespace + "-" + req.Name
			crSpec.SharedMemory.SharedMemoryMiB = funccfg.SharedMemoryMiB
		}

		if functionKind == examplecomv1.FunctionCRKindGPU ||
			functionKind == examplecomv1.FunctionCRKindCPU {
			crSpec.Params = pCRData.Spec.Params
		}
		crSpec.PreviousFunctions = pCRData.Spec.PreviousWBFunctions
		crSpec.NextFunctions = pCRData.Spec.NextWBFunctions

		// Params[examplecomv1.FPS] in CPUFunction CR is passed from Capacity in WBFunction
		if functionKind == examplecomv1.FunctionCRKindCPU &&
			crSpec.Params != nil {
			crSpec.Params[examplecomv1.FPS] = intstr.FromInt32(pCRData.Spec.Requirements.Capacity)
		}

		var deviceFilePath string
		var deviceUUID string

		// Repeat for the number of devices in the Infrastructure information ConfigMap
		for infraIndex := 0; infraIndex < len(infraDevice); infraIndex++ {
			infraData := infraDevice[infraIndex]
			if infraData.NodeName != pCRData.Spec.NodeName {
				continue
			}
			if infraData.DeviceType != pCRData.Spec.DeviceType {
				continue
			}
			if infraData.DeviceIndex != pCRData.Spec.DeviceIndex {
				continue
			}
			if nil != infraData.DeviceFilePath {
				deviceFilePath = *infraData.DeviceFilePath
			}
			deviceUUID = *infraData.DeviceUUID
		}
		var acceleratorInfo examplecomv1.AccIDInfo
		var defaultPartitionName string = "0"
		if functionKind == examplecomv1.FunctionCRKindFPGA {
			if nil == pCRData.Spec.FunctionIndex {
				acceleratorInfo.PartitionName = &defaultPartitionName
			}
			acceleratorInfo.ID = deviceFilePath
		} else {
			acceleratorInfo.PartitionName = &req.Name
			acceleratorInfo.ID = deviceUUID
		}
		crSpec.AcceleratorIDs =
			append(crSpec.AcceleratorIDs, acceleratorInfo)
		crData.UnstructuredContent()["spec"] = crSpec

		// Set the parent-child relationship between WBFunctionCR and the created CR
		err = ctrl.SetControllerReference(pCRData, crData, r.Scheme)
		if nil != err {
			logger.Error(err, "ctrl.SetControllerReference Error.")
		} else {
			// Create CR
			err = r.Create(ctx, crData)
			if err != nil {
				logger.Error(err, "CustomResource Create Error.")
			} else {
				logger.Info("CustomResource Create.")
			}
		}
		logger.Info("kind :" + functionKind)
		logger.Info("apiVersion :" + functionAPIVersion)
		logger.Info("name :" + req.Name)
		logger.Info("namespace :" + req.Namespace)
		break //nolint:staticcheck // SA4004: Intentional break
	}
	return err
}

// FINALIZER name generation
func getFinalizerName(ctx context.Context,
	pCRData *examplecomv1.WBFunction) string {

	var finalizername string

	// Value to set in the finalizer
	if len(pCRData.Kind) == 0 {
		finalizername = MYCRKIND + ".finalizers." +
			examplecomv1.GroupVersion.Group + "." +
			examplecomv1.GroupVersion.Version
	} else {
		finalizername = strings.ToLower(pCRData.Kind) + ".finalizers." +
			examplecomv1.GroupVersion.Group + "." +
			examplecomv1.GroupVersion.Version
	}
	return finalizername
}

// Event type determines
func getEventKind(ctx context.Context,
	pCRData *examplecomv1.WBFunction) int32 {

	var eventKind int32
	eventKind = UPDATE

	// Whether or not there is a deletion timestamp
	if !pCRData.ObjectMeta.DeletionTimestamp.IsZero() {
		eventKind = DELETE
	} else if !controllerutil.ContainsFinalizer(pCRData, getFinalizerName(ctx, pCRData)) {
		// Whether or not Finalizer is written
		eventKind = CREATE
	}
	return eventKind
}

func (r *WBFunctionReconciler) getWBFunctionCR(ctx context.Context,
	req ctrl.Request,
	pCRData *examplecomv1.WBFunction) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		// WBFunction CR Get
		err = r.Get(ctx, req.NamespacedName, pCRData)
		if errors.IsNotFound(err) {
			// If WBFunction CR does not exist
			logger.Info(MYCRKIND + " does not exist.")
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err, MYCRKIND+" unable to get.")
			break
		}
		break //nolint:staticcheck // SA4004: Intentional break
	}
	return err
}

// Get the DeviceInfo CR
func (r *WBFunctionReconciler) getDevMngCR(ctx context.Context,
	pCRData *examplecomv1.WBFunction,
	pCRDeviceInfoData *examplecomv1.DeviceInfo,
	req ctrl.Request) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Get DeviceInfo CR
		err = r.Get(ctx,
			client.ObjectKey{
				Name:      DEVICEINFOCRNAME + req.Name,
				Namespace: req.Namespace,
			},
			pCRDeviceInfoData)
		if errors.IsNotFound(err) {
			// If DeviceInfo CR does not exist
			logger.Info(DEVICEINFOCRKIND +
				" does not exist. Namespaces/Name=" +
				req.NamespacedName.Namespace + "/" +
				req.NamespacedName.Name)
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err,
				DEVICEINFOCRKIND+" unable to get. Name="+
					DEVICEINFOCRNAME+pCRData.Spec.NodeName)
			break
		}
	}
	return err
}

// Create a DeviceInfo CR
func (r *WBFunctionReconciler) createDevMngCR(ctx context.Context,
	pCRData *examplecomv1.WBFunction,
	pCRDeviceInfoData *examplecomv1.DeviceInfo,
	req ctrl.Request,
	reqType string) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ {
		jsonstr, err := json.Marshal(pCRData.Spec)
		if nil != err {
			logger.Error(err, "WBFunction.Spec unable to Marshal.")
			break
		}
		err = json.Unmarshal(jsonstr, &pCRDeviceInfoData.Spec.Request)
		if nil != err {
			logger.Error(err, DEVICEINFOCRKIND+" unable to Unmarshal.")
			break
		}
		pCRDeviceInfoData.Spec.Request.MaxDataFlows =
			pCRData.Spec.MaxDataFlows
		pCRDeviceInfoData.Spec.Request.MaxCapacity =
			pCRData.Spec.MaxCapacity
		pCRDeviceInfoData.Spec.Request.Capacity =
			&pCRData.Spec.Requirements.Capacity
		pCRDeviceInfoData.Spec.Request.RequestType = reqType

		pCRDeviceInfoData.Spec.Request.FunctionIndex =
			&pCRData.Status.FunctionIndex

		reqData := &examplecomv1.DeviceInfo{
			ObjectMeta: metav1.ObjectMeta{
				Name:      DEVICEINFOCRNAME + req.Name,
				Namespace: req.Namespace,
			},
			Spec: pCRDeviceInfoData.Spec,
		}

		// Set the parent-child relationship between WBFunctionCR and the created CR
		err = ctrl.SetControllerReference(pCRData, reqData, r.Scheme)
		if nil != err {
			logger.Error(err, "ctrl.SetControllerReference Error.")
		} else {
			err = r.Create(ctx, reqData)
			if err != nil {
				logger.Error(err, "Failed to create RequestCR.")
				break
			} else {
				logger.Info("Success to create RequestCR.")
				break
			}
		}
	}
	return err
}

// Delete DeviceInfo CR
func (r *WBFunctionReconciler) deleteDevMngCR(ctx context.Context,
	pCRDeviceInfoData *examplecomv1.DeviceInfo) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		err := r.Delete(ctx, pCRDeviceInfoData)
		if err != nil {
			logger.Error(err, "Failed to delete RequestCR.")
			break
		} else {
			logger.Info("Success to delete RequestCR.")
			break
		}
	}
	return err
}

// FunctionCR obtain
func (r *WBFunctionReconciler) getFunctionCR(ctx context.Context,
	req ctrl.Request,
	functionKind string,
	functionAPIVersion string,
	pCRFunction *Function) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		tmpData := &unstructured.Unstructured{}
		tmpData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: functionAPIVersion,
			Kind:    functionKind,
		})

		err = r.Get(ctx, req.NamespacedName, tmpData)
		if errors.IsNotFound(err) {
			// If Function CR does not exist
			logger.Info(functionKind + " does not exist.")
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err, functionKind+" unable to get.")
			break
		}
		var jsonstr []byte
		getstr, _, _ := unstructured.NestedString(tmpData.Object,
			"apiVersion")
		pCRFunction.TypeMeta.APIVersion = getstr

		getstr, _, _ = unstructured.NestedString(tmpData.Object, "kind")
		pCRFunction.TypeMeta.Kind = getstr

		getdata, _, _ := unstructured.NestedMap(tmpData.Object, "spec")
		jsonstr, err = json.Marshal(getdata)
		if nil != err {
			logger.Error(err, DEVICEINFOCRKIND+" unable to Marshal.")
			break
		}
		err = json.Unmarshal(jsonstr, &pCRFunction.Spec)
		if nil != err {
			logger.Error(err, DEVICEINFOCRKIND+" unable to Unmarshal.")
			break
		}
		getdata, _, _ = unstructured.NestedMap(tmpData.Object, "status")
		jsonstr, err = json.Marshal(getdata)
		if nil != err {
			logger.Error(err, DEVICEINFOCRKIND+" unable to Marshal.")
			break
		}
		err = json.Unmarshal(jsonstr, &pCRFunction.Status)
		if nil != err {
			logger.Error(err, DEVICEINFOCRKIND+" unable to Unmarshal.")
			break
		}
		break //nolint:staticcheck // SA4004: Intentional break
	}

	return err
}

// Delete to FunctionCR
func (r *WBFunctionReconciler) delFunctionCR(ctx context.Context,
	functionKind string,
	functionAPIVersion string,
	req ctrl.Request) error {
	var err error
	tmp := &unstructured.Unstructured{}
	tmp.SetGroupVersionKind(schema.GroupVersionKind{
		Version: functionAPIVersion,
		Kind:    functionKind,
	})
	tmp.SetName(req.Name)
	tmp.SetNamespace(req.Namespace)
	err = r.Delete(ctx, tmp)
	return err
}

// Update a custom resource
func (r *WBFunctionReconciler) updCustomResource(ctx context.Context,
	pCRData *examplecomv1.WBFunction,
	status examplecomv1.WBDeployStatus,
	evkind int32) error {

	logger := log.FromContext(ctx)

	var err error

	pCRData.Status.Status = status
	if status != examplecomv1.WBDeployStatusDeployed && evkind != DELETE {
		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, getFinalizerName(ctx, pCRData))
	}
	fmt.Printf("updCustomResource %#v\n", pCRData)
	err = r.Update(ctx, pCRData)
	if nil != err {
		logger.Error(err, "Status Update Error.")
	} else {
		logger.Info("Status Update.")
	}
	return err
}

// Delete a custom resource
func (r *WBFunctionReconciler) delCustomResource(ctx context.Context,
	pCRData *examplecomv1.WBFunction) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	// Delete the Finalizer statement.
	if controllerutil.ContainsFinalizer(pCRData,
		getFinalizerName(ctx, pCRData)) {
		controllerutil.RemoveFinalizer(pCRData,
			getFinalizerName(ctx, pCRData))
		err = r.Update(ctx, pCRData)

		if nil != err {
			logger.Error(err, "RemoveFinalizer Update Error.")
		} else {
			logger.Info("Finalizer Delete Success.")
		}
	}
	return err
}

type FunctionKindMap struct {
	DeviceType     string `json:"deviceType"`
	FunctionCRKind string `json:"functionCRKind"`
}

var gFunctionKindMap []FunctionKindMap

type ConfigTable struct {
	name string
}

var configLoadTable = []ConfigTable{
	{"functionkindmap"},
}

// Load ConfigMap (for main)
func LoadConfigMap(r *WBFunctionReconciler) error {

	ctx := context.Background()
	logger := log.FromContext(ctx)

	var cfgdata []byte
	var err error

	for _, record := range configLoadTable {
		err = r.getConfigMap(ctx, record.name, &cfgdata)
		if nil != err {
			break
		}
		if "functionkindmap" == record.name {
			err = json.Unmarshal(cfgdata, &gFunctionKindMap)
		}
		if nil != err {
			logger.Error(err,
				"unable to Unmarshal. ConfigMap="+record.name)
			break
		}
	}
	return err
}

// Get ConfigMap
func (r *WBFunctionReconciler) getConfigMap(ctx context.Context,
	cfgname string, cfgdata *[]byte) error {

	logger := log.FromContext(ctx)
	var mapdata map[string]string

	tmpData := &unstructured.Unstructured{}
	tmpData.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "ConfigMap",
		Version: "v1",
	})

	// Get a ConfigMap by namespace/name
	err := r.Get(context.Background(),
		client.ObjectKey{
			Namespace: "default",
			Name:      cfgname,
		},
		tmpData)
	if errors.IsNotFound(err) {
		logger.Error(err, "ConfigMap does not exist. ConfigName="+cfgname)
	} else if nil != err {
		logger.Error(err, "ConfigMap unable to fetch. ConfigName="+cfgname)
	} else {
		mapdata, _, _ = unstructured.NestedStringMap(tmpData.Object, "data")
		for _, jsonrecord := range mapdata {
			*cfgdata = []byte(jsonrecord)
		}
	}
	return err
}

// Get Config information for Func
func (r *WBFunctionReconciler) getConfigData(
	ctx context.Context,
	configName string,
	pConfigData *examplecomv1.FunctionConfigMap) error {

	var err error
	var mapData map[string]interface{}
	var configSliceData []examplecomv1.FunctionConfigMap

	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ {
		tmpData := &unstructured.Unstructured{}
		tmpData.SetGroupVersionKind(schema.GroupVersionKind{
			Kind:    "ConfigMap",
			Version: "v1",
		})

		err = r.Get(ctx,
			client.ObjectKey{
				Namespace: "default",
				Name:      configName,
			}, tmpData)

		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			break
		} else {
			mapData, _, _ =
				unstructured.NestedMap(tmpData.Object, "data")
			for _, jsonRecord := range mapData {
				err = json.Unmarshal([]byte(jsonRecord.(string)), &configSliceData)
			}
			for _, configData := range configSliceData {
				cmRxProtocol := configData.RxProtocol
				cmTxProtocol := configData.TxProtocol
				if "DMA" == cmRxProtocol || "DMA" == cmTxProtocol {
					*pConfigData = configData
				}
			}
		}
	}
	return err
}
