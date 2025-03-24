/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "WBConnection/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	// Additional imports
	"encoding/json"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// WBConnectionReconciler reconciles a WBConnection object
type WBConnectionReconciler struct {
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
	PENDING = "Pending" // WBConnection.status.status(Pending)
	RUNNING = "Running" // WBConnection.status.status(Running)
)

type ConnectionFuncSpec struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type ConnectionSpec struct {
	DataFlowRef ConnectionFuncSpec `json:"dataFlowRef,omitempty"`
	From        ConnectionFuncSpec `json:"srcFunc,omitempty"`
	To          ConnectionFuncSpec `json:"dstFunc,omitempty"`
}

type ConnectionStatus struct {
	Status string `json:"status,omitempty"`
}

//+kubebuilder:rbac:groups=example.com,resources=wbconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=wbconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=wbconnections/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=ethernetconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=ethernetconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=ethernetconnections/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=pcieconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=pcieconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=pcieconnections/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WBConnection object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *WBConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get WBConnection CR
	var crData examplecomv1.WBConnection
	var eventKind int           // 0:Add, 1:Upd,  2:Del
	var eventConnectionKind int // 0:Add, 1:Upd,  2:Del
	var connectionKind string
	var connectionAPIVersion string
	var crConnectionMetadata metav1.ObjectMeta
	var crConnectionSpec ConnectionSpec
	var crConnectionStatus ConnectionStatus
	var requeueFlag bool
	var err error

	connectionKind = ""
	connectionAPIVersion = ""
	requeueFlag = false

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger.Info("Reconcile start.")
		err = r.Get(ctx, req.NamespacedName, &crData)
		if errors.IsNotFound(err) {
			// If WBConnection CR does not exist
			logger.Info("WBConnection does not exist.")
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err, "unable to get WBConnection.")
			break
		}

		// Create a resource
		ret := convTypeMeta(ctx, &crData, &connectionKind, &connectionAPIVersion)
		if false == ret {
			break
		}

		// Get Event type
		eventKind = getEventKind(ctx, &crData)

		eventConnectionKind = UPDATE
		if eventKind == CREATE || eventKind == UPDATE {
			if strings.HasPrefix(crData.Spec.From.WBFunctionRef.Name, "wb-start-of-chain") == true ||
				strings.HasPrefix(crData.Spec.To.WBFunctionRef.Name, "wb-end-of-chain") == true {
				eventConnectionKind = CREATE
			} else {
				err = r.getCustomResourceData(ctx, req,
					connectionKind, connectionAPIVersion,
					&crConnectionSpec, &crConnectionStatus, &crConnectionMetadata)
				if errors.IsNotFound(err) {
					// If FromFunction CR does not exist
					logger.Info("Maked Connection does not exist.")
					eventConnectionKind = CREATE
				} else if nil != err {
					logger.Error(err, "Maked Connection unable to fetch CR.")
					requeueFlag = true
					break
				}
			}
		}

		if eventKind == CREATE {

			for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Create", "Create Start")

				if CREATE == eventConnectionKind {
					chgStatus := examplecomv1.WBDeployStatusWaiting
					if strings.HasPrefix(crData.Spec.From.WBFunctionRef.Name, "wb-start-of-chain") == true ||
						strings.HasPrefix(crData.Spec.To.WBFunctionRef.Name, "wb-end-of-chain") == true {
						chgStatus = examplecomv1.WBDeployStatusDeployed
					} else if 0 == len(crData.Spec.From.WBFunctionRef.Name) ||
						0 == len(crData.Spec.From.WBFunctionRef.Namespace) ||
						0 == len(crData.Spec.To.WBFunctionRef.Name) ||
						0 == len(crData.Spec.To.WBFunctionRef.Namespace) {
						logger.Info("FromFunction/ToFunction len is zero.")
						break
					} else {
						// Create a resource
						err = r.createCustomResource(ctx, req,
							connectionKind, connectionAPIVersion, &crData)
						if nil != err {
							break
						}
					}

					// Update the spec information in the status field.
					setSpecToStatusData(ctx, &crData)

					logger.Info("Status Information Change start.")
					err = r.updCustomResource(ctx, &crData, eventConnectionKind,
						chgStatus)
					if nil != err {
						break
					}
					logger.Info("Status Information Change end.")
				} else if examplecomv1.WBDeployStatusDeployed !=
					crData.Status.Status {
					requeueFlag = true
				}

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Create", "Create End")
				break //nolint:staticcheck // SA4004: Intentional break
			}

		} else if eventKind == UPDATE {

			for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

				// In case of update
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Update", "Update Start")

				logger.Info("eventConnectionKind=" + strconv.Itoa(eventConnectionKind))
				logger.Info("crData.Status.Status=" +
					string(crData.Status.Status))
				logger.Info("crConnectionStatus.Status=" + crConnectionStatus.Status)

				if UPDATE == eventConnectionKind {
					if examplecomv1.WBDeployStatusDeployed !=
						crData.Status.Status {
						if RUNNING == crConnectionStatus.Status {
							// Processing is complete, so update the status
							logger.Info("Status Running Change start.")
							err = r.updCustomResource(ctx, &crData,
								eventConnectionKind,
								examplecomv1.WBDeployStatusDeployed)
							if nil != err {
								break
							}
							logger.Info("Status Running Change end.")
						}
					}
				}
				if examplecomv1.WBDeployStatusDeployed !=
					crData.Status.Status {
					requeueFlag = true
				}

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Update", "Update End")
				break //nolint:staticcheck // SA4004: Intentional break
			}

		} else if eventKind == DELETE {
			// In case of deletion
			r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
				"Delete", "Delete Start")

			// check topologyinfo usage flag
			if crData.Status.ConnectionPath != nil {
				// check if capacity substract process has finished
				if controllerutil.ContainsFinalizer(&crData, "topologyinfo.finalizers.example.com.v1") {
					requeueFlag = true
					break
				}
			}

			var deleteFlag bool
			deleteFlag = false

			err = r.getCustomResourceData(ctx, req,
				connectionKind, connectionAPIVersion,
				&crConnectionSpec, &crConnectionStatus, &crConnectionMetadata)
			if errors.IsNotFound(err) {
				// If FromFunction CR does not exist
				logger.Info("Deleted Connection does not exist.")
				deleteFlag = true
			} else if nil != err {
				logger.Error(err, "Deleted Connection unable to fetch CR.")
				requeueFlag = true
				break
			}

			if deleteFlag == false {
				if crConnectionMetadata.DeletionTimestamp.IsZero() {
					// Delete ConnectionCR
					err = r.deleteConnectionCR(ctx, req, connectionKind, connectionAPIVersion)
					if err != nil {
						logger.Error(err, "Delete ConnectionCR Request Failed.")
					} else {
						logger.Info("Delete ConnectionCR Request Successful.")
					}
				}
				if crData.Status.Status == examplecomv1.WBDeployStatusDeployed {
					// Changed Status statement.
					err = r.updCustomResource(ctx, &crData,
						eventConnectionKind,
						examplecomv1.WBDeployStatusTerminating)
					if err != nil {
						logger.Error(err, "WBConnectionCR state transition failed.")
					} else {
						logger.Info("WBConnectionCR state transition success.")
					}
				} else {
					requeueFlag = true
				}
				break
			}
			// Changed Status.
			crData.Status.Status = examplecomv1.WBDeployStatusReleased

			// Delete the Finalizer statement.
			err = r.deleteWBConnectionCR(ctx, &crData)
			if err != nil {
				break
			}

			r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
				"Delete", "Delete End")
		}
		break //nolint:staticcheck // SA4004: Intentional break
	}
	if requeueFlag == true {
		return ctrl.Result{Requeue: requeueFlag}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WBConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.WBConnection{}).
		Complete(r)
}

// Copy the difference from the spec area to the status area
func setSpecToStatusData(
	ctx context.Context,
	pCRData *examplecomv1.WBConnection) {

	pCRData.Status.DataFlowRef = pCRData.Spec.DataFlowRef
	pCRData.Status.ConnectionMethod = pCRData.Spec.ConnectionMethod
	pCRData.Status.From = pCRData.Spec.From
	pCRData.Status.To = pCRData.Spec.To
	pCRData.Status.Params = pCRData.Spec.Params
	pCRData.Status.ConnectionPath = pCRData.Spec.ConnectionPath
	pCRData.Status.SatisfiedRequirements = pCRData.Spec.Requirements

	return
}

const (
	APIVERSION = "example.com/v1"
)

// Create CR attribute information and obtain it
func convTypeMeta(
	ctx context.Context,
	pCRData *examplecomv1.WBConnection,
	pConnectionKind *string,
	pConnectionAPIVersion *string) bool {

	logger := log.FromContext(ctx)

	var ret bool = false

	for count := 0; count < len(gConnectionKindmap); count++ {
		if gConnectionKindmap[count].ConnectionMethod == pCRData.Spec.ConnectionMethod {
			*pConnectionKind = gConnectionKindmap[count].ConnectionCRKind
			*pConnectionAPIVersion = APIVERSION
		}
	}

	if *pConnectionAPIVersion == "" || *pConnectionKind == "" {
		logger.Info("Convert Resource Kind Error." +
			" kind=" + *pConnectionKind +
			" apivesion=" + *pConnectionAPIVersion)
	} else {
		ret = true
	}
	return ret
}

// Check the resource generation status for each CRC
func (r *WBConnectionReconciler) getCustomResourceData(
	ctx context.Context,
	req ctrl.Request,
	connectionKind string,
	connectionAPIVersion string,
	pConnectionSpec *ConnectionSpec,
	pConnectionStatus *ConnectionStatus,
	pConnectionMetadata *metav1.ObjectMeta) error {

	logger := log.FromContext(ctx)

	var err error

	var crDataStringMap map[string]interface{}

	crData := &unstructured.Unstructured{}

	crData.SetGroupVersionKind(schema.GroupVersionKind{
		Version: connectionAPIVersion,
		Kind:    connectionKind,
	})
	err = r.Get(ctx,
		client.ObjectKey{
			Namespace: req.Namespace,
			Name:      req.Name,
		},
		crData)

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		if err != nil {
			break
		}

		if len(crData.Object) != 0 {
			// Store metadata information
			crDataStringMap, _, _ = unstructured.NestedMap(crData.Object, "metadata")

			// Convert the obtained mapdata to byte type
			metadatabytes, err := json.Marshal(crDataStringMap)
			if err != nil {
				logger.Error(err, "unable to json.marshal.")
				break
			}
			// Replace with a struct
			err = json.Unmarshal(metadatabytes, pConnectionMetadata)
			if err != nil {
				logger.Error(err, "unable to json.unmarshal.")
				break
			}

			// Store spec information
			crDataStringMap, _, _ = unstructured.NestedMap(crData.Object, "spec")

			// Convert the obtained mapdata to byte type
			specbytes, err := json.Marshal(crDataStringMap)
			if err != nil {
				logger.Error(err, "unable to json.marshal.")
				break
			}
			// Replace with a struct
			err = json.Unmarshal(specbytes, pConnectionSpec)
			if err != nil {
				logger.Error(err, "unable to json.unmarshal.")
				break
			}

			// Store the status information
			crDataStringMap, _, _ = unstructured.NestedMap(crData.Object, "status")

			// Convert the obtained mapdata to byte type
			statusbytes, err := json.Marshal(crDataStringMap)
			if err != nil {
				logger.Error(err, "unable to json.marshal.")
				break
			}
			// Replace with a struct
			err = json.Unmarshal(statusbytes, pConnectionStatus)
			if err != nil {
				logger.Error(err, "unable to json.unmarshal.")
				break
			}
		}
		break //nolint:staticcheck // SA4004: Intentional break
	}
	return err
}

// Each CRC is created using Resource
func (r *WBConnectionReconciler) createCustomResource(
	ctx context.Context,
	req ctrl.Request,
	connectionKind string,
	connectionAPIVersion string,
	//	pConnectionSpec *ConnectionSpec,
	pCRData *examplecomv1.WBConnection) error {

	logger := log.FromContext(ctx)

	var err error // To determine the return value

	// Organize the cr information to be created
	crData := &unstructured.Unstructured{}
	crData.SetGroupVersionKind(schema.GroupVersionKind{
		Version: connectionAPIVersion,
		Kind:    connectionKind,
	})
	crData.SetName(req.Name)
	crData.SetNamespace(req.Namespace)
	crSpec := map[string]interface{}{
		"dataFlowRef": map[string]interface{}{
			"name":      pCRData.Spec.DataFlowRef.Name,
			"namespace": pCRData.Spec.DataFlowRef.Namespace,
		},
		"from": map[string]interface{}{
			"wbFunctionRef": map[string]interface{}{
				"name":      pCRData.Spec.From.WBFunctionRef.Name,
				"namespace": pCRData.Spec.From.WBFunctionRef.Namespace,
			},
		},
		"to": map[string]interface{}{
			"wbFunctionRef": map[string]interface{}{
				"name":      pCRData.Spec.To.WBFunctionRef.Name,
				"namespace": pCRData.Spec.To.WBFunctionRef.Namespace,
			},
		},
	}
	crData.UnstructuredContent()["spec"] = crSpec

	// Set the parent-child relationship between WBConnectionCR and the created CR
	err = ctrl.SetControllerReference(pCRData, crData, r.Scheme)
	if err != nil {
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
	logger.Info("kind :" + connectionKind)
	logger.Info("apiVersion :" + connectionAPIVersion)
	logger.Info("name :" + req.Name)
	logger.Info("namespace :" + req.Namespace)

	return err
}

func getFinalizerName(ctx context.Context,
	pCRData *examplecomv1.WBConnection) string {
	logger := log.FromContext(ctx)
	var finalizername string

	// Value to set in the finalizer
	if len(pCRData.Kind) == 0 ||
		len(pCRData.APIVersion) == 0 {
		finalizername = "WBConnection.finalizers.example.com.v1"
	} else {
		finalizername = strings.ToLower(pCRData.Kind) + ".finalizers." +
			strings.ReplaceAll(pCRData.APIVersion, "/", ".")
	}
	logger.Info("Finalizername=" + finalizername)
	return finalizername
}

func getEventKind(ctx context.Context,
	pCRData *examplecomv1.WBConnection) int {
	var eventKind int
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

func (r *WBConnectionReconciler) updCustomResource(ctx context.Context,
	pCRData *examplecomv1.WBConnection,
	eventConnectionKind int,
	status examplecomv1.WBDeployStatus) error {
	logger := log.FromContext(ctx)
	var err error

	// status update
	pCRData.Status.Status = status

	if CREATE == eventConnectionKind {
		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, getFinalizerName(ctx, pCRData))
	}
	err = r.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "Status Update Error.")
	} else {
		logger.Info("Status Update.")
	}
	return err
}

func (r *WBConnectionReconciler) deleteConnectionCR(ctx context.Context,
	req ctrl.Request,
	connectionKind string,
	connectionAPIVersion string) error {

	logger := log.FromContext(ctx)

	// Organize the cr information to be created
	crData := &unstructured.Unstructured{}
	crData.SetGroupVersionKind(schema.GroupVersionKind{
		Version: connectionAPIVersion,
		Kind:    connectionKind,
	})
	crData.SetName(req.Name)
	crData.SetNamespace(req.Namespace)

	err := r.Delete(ctx, crData)
	if err != nil {
		logger.Error(err, "Failed to delete ConnectionCR.")
	} else {
		logger.Info("Success to delete ConnectionCR.")
	}

	return err
}

func (r *WBConnectionReconciler) deleteWBConnectionCR(ctx context.Context,
	pCRData *examplecomv1.WBConnection) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	// Delete the Finalizer statement.
	if controllerutil.ContainsFinalizer(pCRData,
		getFinalizerName(ctx, pCRData)) {
		controllerutil.RemoveFinalizer(pCRData, getFinalizerName(ctx, pCRData))
		err := r.Update(ctx, pCRData)
		if err != nil {
			logger.Error(err, "RemoveFinalizer Update Error.")
		} else {
			logger.Info("RemoveFinalizer Update.")
		}
	}
	return err
}

type ConnectionKindMap struct {
	ConnectionMethod string `json:"connectionMethod"`
	ConnectionCRKind string `json:"connectionCRKind"`
}

var gConnectionKindmap []ConnectionKindMap

type ConfigTable struct {
	name string
}

var configLoadTable = []ConfigTable{
	{"connectionkindmap"},
}

// Load ConfigMap (for main)
func LoadConfigMap(r *WBConnectionReconciler) error {

	ctx := context.Background()
	logger := log.FromContext(ctx)

	var cfgdata []byte
	var err error

	for _, record := range configLoadTable {
		err = r.getConfigMap(ctx, record.name, &cfgdata)
		if nil != err {
			break
		}
		if "connectionkindmap" == record.name {
			err = json.Unmarshal(cfgdata, &gConnectionKindmap)
		}
		if nil != err {
			logger.Error(err, "unable to Unmarshal. ConfigMap="+record.name)
			break
		}
	}
	return err
}

// Get ConfigMap
func (r *WBConnectionReconciler) getConfigMap(ctx context.Context,
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
