/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "DeviceInfo/api/v1"
	// Additional imports
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	// "strconv"
	"strings"
)

// DeviceInfoReconciler reconciles a DeviceInfo object
type DeviceInfoReconciler struct {
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

const (
	DEFAULTSECOND      = 10  // Default requeue interval (s)
	DEFAULTMICROSECOND = 500 // Default requeue interval (ms)
)

const (
	COMPUTERESOURCEAPIVERSION = "example.com/v1"
	COMPUTERESOUCEKIND        = "ComputeResource"
	COMPUTERESOURCENAME       = "compute-"
	DEVICEINFOAPIVERSION      = "example.com/v1"
	DEVICEINFOKIND            = "DeviceInfo"
)

var gMyNodeName string    // Store the node name
var gMyClusterName string // Store the node name

//+kubebuilder:rbac:groups=example.com,resources=fpgas,verbs=get;list;watch
//+kubebuilder:rbac:groups=example.com,resources=fpgas/status,verbs=get
//+kubebuilder:rbac:groups=example.com,resources=fpgas/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=childbs,verbs=get;list;watch
//+kubebuilder:rbac:groups=example.com,resources=childbs/status,verbs=get
//+kubebuilder:rbac:groups=example.com,resources=childbs/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=deviceinfos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=deviceinfos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=deviceinfos/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=computeresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=computeresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=computeresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeviceInfo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *DeviceInfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get DeviceInfo CR
	var crDeviceInfoData examplecomv1.DeviceInfo
	var crComputeResourceData examplecomv1.ComputeResource
	var requeueFlag bool
	var err error

	requeueFlag = false

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		logger.Info("Reconcile start.")

		if true == strings.Contains(req.NamespacedName.Name, "deviceinfo-") {

			err = r.Get(ctx, req.NamespacedName, &crDeviceInfoData)
			if errors.IsNotFound(err) {
				// If DeviceInfo CR does not exist
				logger.Info("DeviceInfo does not exist.")
				break
			} else if nil != err {
				// If an error occurs in the Get function
				logger.Error(err, "unable to get DeviceInfo.")
				break
			}

			// Node check
			if gMyNodeName != crDeviceInfoData.Spec.Request.NodeName {
				// nothing to do
				break
			}

			// Get Event type
			var eventKind int // 0:Add, 1:Upd,  2:Del
			eventKind = getEventKind(ctx, &crDeviceInfoData)

			if CREATE == eventKind {
				controllerutil.AddFinalizer(&crDeviceInfoData, getFinalizerName(ctx, &crDeviceInfoData))

				crRequestData := crDeviceInfoData.Spec.Request

				err := r.Get(ctx, client.ObjectKey{
					Namespace: gMyClusterName,
					Name:      COMPUTERESOURCENAME + gMyNodeName,
				}, &crComputeResourceData)
				if errors.IsNotFound(err) {
					// If ComputeResource CR does not exist
					logger.Info("ComputeResource does not exist.")
					break
				} else if nil != err {
					// If an error occurs in the Get function
					logger.Error(err, "unable to get ComputeResource.")
					break
				}

				if crRequestData.RequestType == examplecomv1.RequestDeploy {
					// When requesting allocation of deployment area
					err = r.makeDeploySpace(ctx, &req, &crDeviceInfoData, &crComputeResourceData)
					if errors.IsConflict(err) {
						requeueFlag = true
						break
					} else if nil != err {
						logger.Error(err, "makeDeploySpace() Error.")
						requeueFlag = true
						break
					}
				} else if crRequestData.RequestType == examplecomv1.RequestUndeploy {
					// In case of a deployment area release request
					err = r.freeDeploySpace(ctx, &req, &crDeviceInfoData, &crComputeResourceData)
					if errors.IsConflict(err) {
						requeueFlag = true
						break
					} else if nil != err {
						logger.Error(err, "freeDeploySpace() Error.")
						requeueFlag = true
						break
					}
				} else {
					// does not process
					break
				}
			} else if UPDATE == eventKind {
				// No processing
				logger.Info("Update process start.")
				logger.Info("Update process end.")
				break
			} else if DELETE == eventKind {
				logger.Info("Delete process start.")
				_ = r.delCustomResource(ctx, &crDeviceInfoData)
				logger.Info("Delete process end.")
				break
			}
			break //nolint:staticcheck // SA4004: Intentional break
		} else {
			err = r.ChildBSCR(ctx, req)
			if nil != err {
				requeueFlag = true
			}
			break
		}
	}
	logger.Info("Reconcile end.")
	if requeueFlag == true {
		return ctrl.Result{Requeue: requeueFlag}, nil
	}
	return ctrl.Result{}, nil
}

// ChildBitstream CR
func (r *DeviceInfoReconciler) ChildBSCR(ctx context.Context, req ctrl.Request) error {
	logger := log.FromContext(ctx)

	var err error
	var crChildBs examplecomv1.ChildBs
	var crComputeResourceData examplecomv1.ComputeResource
	var deviceRequest examplecomv1.WBFuncRequest
	var crFPGA examplecomv1.FPGA

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		err = r.Get(ctx, req.NamespacedName, &crChildBs)
		if errors.IsNotFound(err) {
			// If ChildBitstream CR does not exist
			logger.Info("ChildBitstream CR does not exist.")
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err, "unable to get ChildBitstream.")
			break
		}

		// Get FPGA CR
		namespacedName := client.ObjectKey{Name: crChildBs.OwnerReferences[0].Name,
			Namespace: "default"}
		err = r.Get(ctx, namespacedName, &crFPGA)
		if errors.IsNotFound(err) {
			// If FPGA CR does not exist
			logger.Info("FPGA CR does not exist.")
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err, "unable to get FPGA CR.")
			break
		}

		if crFPGA.Status.NodeName != gMyNodeName {
			break
		}

		err := r.Get(ctx, client.ObjectKey{
			Namespace: gMyClusterName,
			Name:      COMPUTERESOURCENAME + gMyNodeName,
		}, &crComputeResourceData)
		if errors.IsNotFound(err) {
			// If ComputeResource CR does not exist
			logger.Info("ComputeResource does not exist.")
			break
		} else if nil != err {
			// If an error occurs in the Get function
			logger.Error(err, "unable to get ComputeResource.")
			break
		}

		for index := 0; index < len(crChildBs.Spec.Regions); index++ {
			childData := crChildBs.Spec.Regions[index]
			childDataFunctions := childData.Modules.Functions
			deviceRequest.RequestType = "Update"
			deviceRequest.DeviceType = "alveou250"
			deviceRequest.DeviceIndex = crFPGA.Status.DeviceIndex
			deviceRequest.RegionName = *childData.Name
			deviceRequest.NodeName = crFPGA.Status.NodeName
			deviceRequest.FunctionName = *(*(*childDataFunctions)[0].Module)[0].Type

			// Get the target RegionInfos from ComputeResource
			err, crComputeResourceRegionInfo, regionIndex :=
				getComputeResourceTargetDeviceInfo(ctx,
					&deviceRequest,
					&crComputeResourceData.Spec.NodeName,
					&crComputeResourceData.Spec.Regions)
			if nil != err {
				break
			}

			switch crChildBs.Status.State {
			case examplecomv1.ChildBsWritingBsfile,
				examplecomv1.ChildBsConfiguringParam,
				examplecomv1.ChildBsNoConfigureNetwork,
				examplecomv1.ChildBsConfiguringNetwork:
				crComputeResourceRegionInfo.Available = false
				crComputeResourceRegionInfo.Status = examplecomv1.WBRegionStatusPreparing
			case examplecomv1.ChildBsReady:
				// Get information for initial ComputeResourceCR generation
				err = loadConfigMap(ctx, r)
				if nil != err {
					break
				}
				updateComputeResourceCR(ctx, crComputeResourceRegionInfo)

			case examplecomv1.ChildBsError:
				crComputeResourceRegionInfo.Available = false
				crComputeResourceRegionInfo.Status = examplecomv1.WBRegionStatusError
			}
			crComputeResourceData.Spec.Regions[regionIndex] = *crComputeResourceRegionInfo
			crComputeResourceData.Status.Regions[regionIndex] = *crComputeResourceRegionInfo
		}
		err = r.Update(ctx, &crComputeResourceData)
		if errors.IsConflict(err) {
			logger.Info("ComputeResource CR Update Conflict")
			break
		} else if nil != err {
			logger.Error(err, "unable to get ComputeResource.")
		} else {
			logger.Info("ComputeResource Update Success")
			break
		}
	}
	return err
}

func updateComputeResourceCR(ctx context.Context, regioninfo *examplecomv1.RegionInfo) {

	deployDevice := gDeployInfo["devices"]

	// Repeat for the number of devices in Deployment Information ConfigMap
	for deployIndex := 0; deployIndex < len(deployDevice); deployIndex++ {
		deployDevice := deployDevice[deployIndex]

		if deployDevice.NodeName != gMyNodeName {
			// If the node in the deployment information does not match the current node
			continue
		}

		// @TODO Must be modified when converting FPGA UUID
		if deployDevice.DeviceUUID == nil && deployDevice.DeviceFilePath == nil {
			continue
		} else if deployDevice.DeviceUUID != nil {
			if *deployDevice.DeviceUUID != *regioninfo.DeviceUUID {
				// If the DeviceUUID in the ComputeResource and the deployment information do not match
				continue
			}
		} else {
			// @TODO Must be modified when converting FPGA UUID
			if *deployDevice.DeviceFilePath != regioninfo.DeviceFilePath {
				// If the DeviceFilePath of the ComputeResource and the deployment information do not match
				continue
			}
		}

		if 0 != len(regioninfo.Functions) {
			break
		}

		regionList, funcData := createFunctionData(*deployDevice.FunctionTargets)

		// Create RegionInfos for each FunctionTargets(Region)
		for regionIndex := 0; regionIndex < len(*deployDevice.FunctionTargets); regionIndex++ {
			deployFunctionTargets := (*deployDevice.FunctionTargets)[regionIndex]

			if deployFunctionTargets.RegionName != regioninfo.Name {
				// break
				continue
			}

			var defaultCurrentFunctions int32
			var defaultCurrentCapacity int32
			defaultCurrentFunctions = 0
			defaultCurrentCapacity = 0

			regioninfo.Available = true
			regioninfo.MaxFunctions = deployFunctionTargets.MaxFunctions
			regioninfo.CurrentFunctions = &defaultCurrentFunctions
			regioninfo.MaxCapacity = deployFunctionTargets.MaxCapacity
			regioninfo.CurrentCapacity = &defaultCurrentCapacity

			if 0 != len(regionList) &&
				0 != len(funcData[regionList[regionIndex]]) {
				regioninfo.Functions = funcData[regionList[regionIndex]]
				*regioninfo.CurrentFunctions = int32(len(funcData[regionList[regionIndex]]))
				regioninfo.Status = examplecomv1.WBRegionStatusReady
			} else {
				if "" != regioninfo.DeviceFilePath {
					regioninfo.Status = examplecomv1.WBRegionStatusNotReady
				}
				regioninfo.Status = examplecomv1.WBRegionStatusReady
			}

			// fmt.Println("regionIndex : ", regionIndex) //debug
			// fmt.Println("regionList : ", regionList[regionIndex]) //debug
			// fmt.Println("regioninfo.Functions : ", regioninfo.Functions) //debug
		}
	}
	// fmt.Println("data : ", data) //debug
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeviceInfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.DeviceInfo{}).
		Watches(&examplecomv1.ChildBs{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

// FINALIZER name generation
func getFinalizerName(ctx context.Context,
	pCRData *examplecomv1.DeviceInfo) string {

	var finalizername string

	// Value to set in the finalizer
	if len(pCRData.Kind) == 0 {
		finalizername = DEVICEINFOKIND + ".finalizers." +
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
	pCRData *examplecomv1.DeviceInfo) int {

	var eventKind int
	eventKind = UPDATE

	// Whether or not there is a deletion timestamp
	if !pCRData.ObjectMeta.DeletionTimestamp.IsZero() {
		eventKind = DELETE
	} else if !controllerutil.ContainsFinalizer(pCRData, getFinalizerName(ctx, pCRData)) {
		// Whether or not Finalizer is written
		eventKind = CREATE
	}
	//	fmt.Printf("getEventKind %#v\n", pCRData)
	return eventKind
}

// Delete a custom resource
func (r *DeviceInfoReconciler) delCustomResource(ctx context.Context,
	pCRData *examplecomv1.DeviceInfo) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	// Delete the Finalizer statement.
	if controllerutil.ContainsFinalizer(pCRData,
		getFinalizerName(ctx, pCRData)) {
		controllerutil.RemoveFinalizer(pCRData,
			getFinalizerName(ctx, pCRData))
		err := r.Update(ctx, pCRData)
		if nil != err {
			logger.Error(err, "RemoveFinalizer Update Error.")
		}
	}
	return err
}

// Allocate deployment area
func (r *DeviceInfoReconciler) makeDeploySpace(ctx context.Context,
	req *ctrl.Request,
	pCRDeviceInfo *examplecomv1.DeviceInfo,
	pCRComputeResource *examplecomv1.ComputeResource) error {

	logger := log.FromContext(ctx)
	var err error
	//	var specifyingFlag bool // Whether or not FunctionIndex is specified
	var ResponceData examplecomv1.WBFuncResponse
	var crComputeResourceRegionInfo *examplecomv1.RegionInfo
	var regionIndex int32

	err = nil
	//	specifyingFlag = true
	ResponceData.Status = examplecomv1.ResponceDeployed

	deviceRequest := pCRDeviceInfo.Spec.Request

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Get the target RegionInfos from ComputeResource
		err, crComputeResourceRegionInfo, regionIndex =
			getComputeResourceTargetDeviceInfo(ctx,
				&deviceRequest,
				&pCRComputeResource.Spec.NodeName,
				&pCRComputeResource.Spec.Regions)
		if nil != err {
			break
		}

		// Allocation of deployment area
		err = deployProccessing(ctx,
			&deviceRequest, crComputeResourceRegionInfo, &ResponceData)
		if nil != err {
			break
		}
		pCRComputeResource.Spec.Regions[regionIndex] = *crComputeResourceRegionInfo

	}
	if nil == err {
		// Update the ComputeResource
		pCRComputeResource.Status.Regions = pCRComputeResource.Spec.Regions
		if err == nil {
			err = r.Update(ctx, pCRComputeResource)
			if errors.IsConflict(err) {
				logger.Info("ComputeResource CR Update Conflict")
			} else if nil != err {
				logger.Error(err, "unable to get ComputeResource.")
			} else {
				logger.Info("ComputeResource Update Success")
				// Pass the results to WBFunction
				r.createResponseCR(ctx, pCRDeviceInfo, &ResponceData, err)
			}
		}
	}

	return err
}

// Get the target RegionInfos from ComputeResource
func getComputeResourceTargetDeviceInfo(ctx context.Context,
	pCRDeviceInfo *examplecomv1.WBFuncRequest,
	pCRComputeResourceNodeName *string,
	pCRComputeResource *[]examplecomv1.RegionInfo) (error, *examplecomv1.RegionInfo, int32) {

	var regionIndex int32 = -1
	var pCRComputeResourceRegionInfo *examplecomv1.RegionInfo

	for index := 0; index < len(*pCRComputeResource); index++ {
		regionData := (*pCRComputeResource)[index]
		switch pCRDeviceInfo.RequestType {
		case "Deploy":
			fallthrough
		case "Undeploy":
			if *pCRComputeResourceNodeName != pCRDeviceInfo.NodeName {
				continue
			}
			if regionData.Name != pCRDeviceInfo.RegionName {
				continue
			}
			if regionData.DeviceType != pCRDeviceInfo.DeviceType {
				continue
			}
			if regionData.DeviceIndex != pCRDeviceInfo.DeviceIndex {
				continue
			}
			// Temporarily store the corresponding RegionInfo.
			pCRComputeResourceRegionInfo = &regionData
			regionIndex = int32(index)
			break
		case "Update":
			if *pCRComputeResourceNodeName != pCRDeviceInfo.NodeName {
				continue
			}
			if regionData.Name != pCRDeviceInfo.RegionName {
				continue
			}
			if regionData.DeviceIndex != pCRDeviceInfo.DeviceIndex {
				continue
			}
			// Temporarily store the corresponding RegionInfo.
			pCRComputeResourceRegionInfo = &regionData
			regionIndex = int32(index)
			break
		}
	}
	if pCRComputeResourceRegionInfo == nil {
		// If there is no corresponding information in the region information
		return fmt.Errorf("compute resource get error."), pCRComputeResourceRegionInfo, -1
	}

	return nil, pCRComputeResourceRegionInfo, regionIndex
}

// Deployment process
func deployProccessing(ctx context.Context,
	pCRDeviceInfo *examplecomv1.WBFuncRequest,
	pRegionInfo *examplecomv1.RegionInfo,
	pResponceData *examplecomv1.WBFuncResponse) error {

	var err error
	var addFlag bool // Add to the list of Functions for judgment
	var targetFunction examplecomv1.FunctionInfrastruct
	var functionsIndex *int32

	err = nil
	addFlag = false
	pFunctions := pRegionInfo.Functions

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Additional deployment to existing deployment area

		if len(pFunctions) == 0 {
			addFlag = true
			targetFunction.FunctionIndex = *pCRDeviceInfo.FunctionIndex
			// Set PartitionName
			targetFunction.PartitionName = *pRegionInfo.DeviceUUID

			var defaultInt int32
			defaultInt = 1
			*pRegionInfo.CurrentFunctions = *pRegionInfo.CurrentFunctions + int32(1)
			targetFunction.FunctionName = pCRDeviceInfo.FunctionName
			targetFunction.MaxDataFlows = pCRDeviceInfo.MaxDataFlows
			targetFunction.CurrentDataFlows = &defaultInt
			targetFunction.MaxCapacity = pCRDeviceInfo.MaxCapacity
			targetFunction.CurrentCapacity = pCRDeviceInfo.Capacity
		} else {
			var functionIndex int32
			var functionflag bool
			functionflag = false
			for functionIndex = 0; int(functionIndex) < len(pFunctions); functionIndex++ {
				if pFunctions[functionIndex].FunctionIndex !=
					*pCRDeviceInfo.FunctionIndex {
					continue
				} else {
					// Allocate Functions that contain the specified FunctionIndex
					targetFunction = pRegionInfo.Functions[functionIndex]
					functionsIndex = &functionIndex
					functionflag = true
					break
				}
			}
			if functionflag == false {
				// If there is no free function information
				addFlag = true
				targetFunction.FunctionIndex = *pCRDeviceInfo.FunctionIndex
				var defaultInt int32
				defaultInt = 1
				// Set PartitionName
				targetFunction.PartitionName = *pRegionInfo.DeviceUUID
				*pRegionInfo.CurrentFunctions = *pRegionInfo.CurrentFunctions + int32(1)
				targetFunction.FunctionName = pCRDeviceInfo.FunctionName
				targetFunction.MaxDataFlows = pCRDeviceInfo.MaxDataFlows
				targetFunction.CurrentDataFlows = &defaultInt
				targetFunction.MaxCapacity = pCRDeviceInfo.MaxCapacity
				targetFunction.CurrentCapacity = pCRDeviceInfo.Capacity
				// break
			} else {
				currentDataFlows := *targetFunction.CurrentDataFlows
				*targetFunction.CurrentDataFlows =
					currentDataFlows + int32(1)
				currentCapacity := *targetFunction.CurrentCapacity
				*targetFunction.CurrentCapacity =
					currentCapacity + *pCRDeviceInfo.Capacity
			}
		}
		// Common settings for new/existing
		*pRegionInfo.CurrentCapacity =
			*pRegionInfo.CurrentCapacity + *pCRDeviceInfo.Capacity
		pResponceData.FunctionIndex = &targetFunction.FunctionIndex
		pResponceData.DeviceUUID = *pRegionInfo.DeviceUUID
		pResponceData.DeviceFilePath = pRegionInfo.DeviceFilePath

		maxFunctions := int32(0)
		currentFunctions := int32(0)
		regionMaxCapacity := int32(0)
		regionCurrentCapacity := int32(0)

		if nil != pRegionInfo.MaxFunctions {
			maxFunctions = *pRegionInfo.MaxFunctions
		}
		if nil != pRegionInfo.CurrentFunctions {
			currentFunctions = *pRegionInfo.CurrentFunctions
		}
		if nil != pRegionInfo.MaxCapacity {
			regionMaxCapacity = *pRegionInfo.MaxCapacity
		}
		if nil != pRegionInfo.CurrentCapacity {
			regionCurrentCapacity = *pRegionInfo.CurrentCapacity
		}

		// Set Available
		if maxFunctions-currentFunctions >= 0 &&
			regionMaxCapacity-regionCurrentCapacity > 0 {
			pRegionInfo.Available = true
		} else {
			pRegionInfo.Available = false
		}

		maxDataFlows := int32(0)
		currentDataFlows := int32(0)
		funcMaxCapacity := int32(0)
		funcCurrentCapacity := int32(0)

		if nil != targetFunction.MaxDataFlows {
			maxDataFlows = *targetFunction.MaxDataFlows
		}
		if nil != targetFunction.CurrentDataFlows {
			currentDataFlows = *targetFunction.CurrentDataFlows
		}
		if nil != targetFunction.MaxCapacity {
			funcMaxCapacity = *targetFunction.MaxCapacity
		}
		if nil != targetFunction.CurrentCapacity {
			funcCurrentCapacity = *targetFunction.CurrentCapacity
		}

		if maxDataFlows-currentDataFlows > 0 &&
			funcMaxCapacity-funcCurrentCapacity > 0 {
			targetFunction.Available = true
		} else {
			targetFunction.Available = false
		}
		if true == addFlag {
			pRegionInfo.Functions = append(pRegionInfo.Functions, targetFunction)
		} else {
			pRegionInfo.Functions[int(*functionsIndex)] = targetFunction
		}
	}
	return err
}

// Update the deployment result in the CR Status.
func (r *DeviceInfoReconciler) createResponseCR(ctx context.Context,
	pCRDeviceInfo *examplecomv1.DeviceInfo,
	pResponceData *examplecomv1.WBFuncResponse,
	err error) error {

	logger := log.FromContext(ctx)

	for doWhile := 0; doWhile < 1; doWhile++ {
		if nil != err {
			// If there is a problem collecting information
			pCRDeviceInfo.Status.Response.Status = examplecomv1.ResponceError
			break
		}
		pCRDeviceInfo.Status.Response = *pResponceData
	}

	// Update the Status area
	err = r.Update(ctx, pCRDeviceInfo)
	if nil != err {
		logger.Error(err, "DeviceInfoCR Status Update Error")
	} else {
		logger.Info("DeviceInfoCR Update Success")
	}

	return err
}

// Deployment area release process
func (r *DeviceInfoReconciler) freeDeploySpace(ctx context.Context,
	req *ctrl.Request,
	pCRDeviceInfo *examplecomv1.DeviceInfo,
	pCRComputeResource *examplecomv1.ComputeResource) error {

	logger := log.FromContext(ctx)
	var err error
	var ResponceData examplecomv1.WBFuncResponse
	var crComputeResourceRegionInfo *examplecomv1.RegionInfo
	var regionIndex int32

	err = nil
	ResponceData.Status = examplecomv1.ResponceUndeployed

	deviceRequest := pCRDeviceInfo.Spec.Request

	for doWhile := 0; doWhile < 1; doWhile++ {

		// Get the target RegionInfos from ComputeResource
		err, crComputeResourceRegionInfo, regionIndex =
			getComputeResourceTargetDeviceInfo(ctx,
				&deviceRequest,
				&pCRComputeResource.Spec.NodeName,
				&pCRComputeResource.Spec.Regions)
		if nil != err {
			break
		}

		// Deployment area release process
		err = undeployProccessing(ctx, &deviceRequest, crComputeResourceRegionInfo)
		if nil != err {
			break
		}
		pCRComputeResource.Spec.Regions[regionIndex] = *crComputeResourceRegionInfo
	}
	if nil == err {
		// Update the ComputeResource
		pCRComputeResource.Status.Regions = pCRComputeResource.Spec.Regions
		err = r.Update(ctx, pCRComputeResource)
		if errors.IsConflict(err) {
			logger.Info("ComputeResource CR Update Conflict")
		} else if nil != err {
			logger.Error(err, "unable to get ComputeResource.")
		} else {
			logger.Info("ComputeResource Update Success")
			// Pass the results to WBFunction
			r.createResponseCR(ctx, pCRDeviceInfo, &ResponceData, err)
		}
	}
	return err
}

// Deployment area release process
func undeployProccessing(ctx context.Context,
	pCRDeviceInfo *examplecomv1.WBFuncRequest,
	pRegionInfo *examplecomv1.RegionInfo) error {

	var err error
	var targetRegion examplecomv1.RegionInfo
	var targetFunction examplecomv1.FunctionInfrastruct
	var functionsIndex int

	err = nil
	pFunctions := pRegionInfo.Functions

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Release the deployment area
		for functionIndex := 0; functionIndex < len(pFunctions); functionIndex++ {
			if pFunctions[functionIndex].FunctionIndex !=
				*pCRDeviceInfo.FunctionIndex {
				continue
			} else {
				// Store the Functions containing the FunctionIndex to be deleted.
				targetFunction = pFunctions[functionIndex]
				functionsIndex = functionIndex
				break
			}
		}
		if targetFunction.PartitionName == "" {
			err = fmt.Errorf("function info get error.")
			break
		}

		// Usage settings
		*pRegionInfo.CurrentCapacity =
			*pRegionInfo.CurrentCapacity - *pCRDeviceInfo.Capacity

		fpgaJudge := acceleratorJudge(ctx, pRegionInfo.Name, pRegionInfo.DeviceFilePath)
		if false == fpgaJudge && *targetFunction.CurrentDataFlows-int32(1) == 0 {
			*pRegionInfo.CurrentFunctions = *pRegionInfo.CurrentFunctions - int32(1)
		}

		// Set Available
		if *pRegionInfo.MaxFunctions-*pRegionInfo.CurrentFunctions >= 0 &&
			*pRegionInfo.MaxCapacity-*pRegionInfo.CurrentCapacity > 0 {
			pRegionInfo.Available = true
		} else {
			pRegionInfo.Available = false
		}

		if false == fpgaJudge {
			targetRegion.Available = pRegionInfo.Available
			targetRegion.CurrentCapacity = pRegionInfo.CurrentCapacity
			targetRegion.CurrentFunctions = pRegionInfo.CurrentFunctions
			targetRegion.DeviceFilePath = pRegionInfo.DeviceFilePath
			targetRegion.DeviceIndex = pRegionInfo.DeviceIndex
			targetRegion.DeviceType = pRegionInfo.DeviceType
			targetRegion.DeviceUUID = pRegionInfo.DeviceUUID
			targetRegion.MaxCapacity = pRegionInfo.MaxCapacity
			targetRegion.MaxFunctions = pRegionInfo.MaxFunctions
			targetRegion.Name = pRegionInfo.Name
			targetRegion.Status = pRegionInfo.Status
			targetRegion.Type = pRegionInfo.Type

			pRegionInfo.Functions = pRegionInfo.Functions[:functionsIndex+
				copy(pRegionInfo.Functions[functionsIndex:],
					pRegionInfo.Functions[functionsIndex+1:])]
			if 0 != len(pRegionInfo.Functions) {
				targetRegion.Functions = pRegionInfo.Functions
			}

			*pRegionInfo = targetRegion
			break
		}

		*targetFunction.CurrentDataFlows =
			*targetFunction.CurrentDataFlows - int32(1)
		*targetFunction.CurrentCapacity =
			*targetFunction.CurrentCapacity - *pCRDeviceInfo.Capacity

		if *targetFunction.MaxDataFlows-*targetFunction.CurrentDataFlows > 0 &&
			*targetFunction.MaxCapacity-*targetFunction.CurrentCapacity > 0 {
			targetFunction.Available = true
		} else {
			targetFunction.Available = false
		}
		pRegionInfo.Functions[functionsIndex] = targetFunction
	}
	return err
}

// Accelerator determination (currently FPGA or other)
func acceleratorJudge(ctx context.Context,
	region string,
	deviceFilePath string) bool {
	// old	deviceUUID string) bool {

	var fpgaJudge bool // For FPGA: true
	fpgaJudge = false

	if deviceFilePath != "" {
		fpgaJudge = true
	}
	return fpgaJudge
}

// Infrastructure information ConfigMap
type InfrastructureInfo struct {
	DeviceFilePath *string `json:"deviceFilePath,omitempty"`
	NodeName       string  `json:"nodeName"`
	DeviceUUID     *string `json:"deviceUUID,omitempty"`
	DeviceType     string  `json:"deviceType"`
	DeviceIndex    int32   `json:"deviceIndex"`
}

// Deployment Information ConfigMap
type DeployInfo struct {
	NodeName        string            `json:"nodeName"`
	DeviceFilePath  *string           `json:"deviceFilePath,omitempty"`
	DeviceUUID      *string           `json:"deviceUUID,omitempty"`
	FunctionTargets *[]regionIndevice `json:"functionTargets"`
}

// Deployment Information FunctionTargets
type regionIndevice struct {
	RegionType   string                       `json:"regionType"`
	RegionName   string                       `json:"regionName"`
	MaxFunctions *int32                       `json:"maxFunctions"`
	MaxCapacity  *int32                       `json:"maxCapacity"`
	Functions    *[]simplefunctioninfrastruct `json:"functions"`
}

// Deployment Information functions
type simplefunctioninfrastruct struct {
	FunctionIndex int32  `json:"functionIndex"`
	PartitionName string `json:"partitionName"`
	FunctionName  string `json:"functionName"`
	MaxDataFlows  int32  `json:"maxDataFlows"`
	MaxCapacity   int32  `json:"maxCapacity"`
}

// Config information storage area
var gInfraStructureInfo map[string][]InfrastructureInfo
var gDeployInfo map[string][]DeployInfo

var config_load_tbl = []ConfigTable{
	{"infrastructureinfo"},
	{"deployinfo"},
}

type ConfigTable struct {
	name string
}

// Create ComputeResourceCR (function for startup processing)
func StartupProccessing(r *DeviceInfoReconciler, mng ctrl.Manager) error {
	var crData examplecomv1.ComputeResource
	var err error
	ctx := context.Background()
	logger := log.FromContext(ctx)
	err = nil

	gMyNodeName = os.Getenv("K8S_NODENAME")
	gMyClusterName = os.Getenv("K8S_CLUSTERNAME")

	if gMyNodeName == "" || gMyClusterName == "" {
		logger.Info("The node name or cluster name could not be obtained.")
		return fmt.Errorf("startup error.")
	}

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Get information for initial ComputeResourceCR generation
		err = loadConfigMap(ctx, r)
		if nil != err {
			break
		}

		m := mng.GetAPIReader()
		err = m.Get(context.Background(),
			client.ObjectKey{
				Namespace: gMyClusterName,
				Name:      COMPUTERESOURCENAME + gMyNodeName,
			}, &crData)

		if errors.IsNotFound(err) {
			// If ComputeResourceCR does not exist
			deployData := createComputeResourceCR()
			crData.Spec.Regions = deployData
			crData.Spec.NodeName = gMyNodeName
			crData.Status.Regions = deployData
			crData.Status.NodeName = gMyNodeName
			// fmt.Println("spec data : ", crData.Spec.Regions) // debug

			// ComputeResource generation
			crData.SetName(COMPUTERESOURCENAME + gMyNodeName)
			crData.SetNamespace(gMyClusterName)
			crData.Kind = COMPUTERESOUCEKIND
			crData.APIVersion = COMPUTERESOURCEAPIVERSION

			// fmt.Println("crData : ", crData) //debug

			// Register ComputeResourceCR
			err = r.Create(ctx, &crData)
			if nil != err {
				logger.Error(err, "Startup Create Error")
				break
			}
			logger.Info("Startup Create Success")
		} else if nil == err {
			// If ComputeResourceCR exists
			logger.Info("ComputeResourceCR is exist")
			err = nil
		} else {
			// Failed to get ComputeResourceCR
			logger.Error(err, "Failed to get ComputeResourceCR at startup")
			break
		}
	}
	return err
}

// Load ConfigMap (startup function)
func loadConfigMap(ctx context.Context, r *DeviceInfoReconciler) error {

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
		} else if "deployinfo" == record.name {
			err = json.Unmarshal(cfgdata, &gDeployInfo)
		}
		if nil != err {
			logger.Error(err, "unable to unmarshal. ConfigMap="+record.name)
			break
		}
	}
	return err
}

// Get ConfigMap (startup function)
func (r *DeviceInfoReconciler) getConfigMap(ctx context.Context,
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
		// fmt.Printf("NestedMap %#v\n", mapdata) //debug
		for _, jsonrecord := range mapdata {
			*cfgdata = []byte(jsonrecord)
		}
	}
	return err
}

// Generate CR information for ComputeResource from infrastructure information and deployment area information
func createComputeResourceCR() (regioninfos []examplecomv1.RegionInfo) {

	infraDevice := gInfraStructureInfo["devices"]
	deployDevice := gDeployInfo["devices"]

	// Repeat for the number of devices in the Infrastructure information ConfigMap
	for infraIndex := 0; infraIndex < len(infraDevice); infraIndex++ {
		// fmt.Println("infra_node : ", infraDevice[infraIndex].Node) //debug
		// fmt.Println("infra_node : ", gMyNodeName) //debug
		infraData := infraDevice[infraIndex]

		if infraData.NodeName != gMyNodeName {
			continue
		}

		// For the second cycle: exclude memory information
		if infraData.DeviceType == "memory" {
			continue
		}

		// Repeat for the number of devices in Deployment Information ConfigMap
		for deployIndex := 0; deployIndex < len(deployDevice); deployIndex++ {
			deployDevice := deployDevice[deployIndex]

			if deployDevice.NodeName != gMyNodeName {
				// If the node in the deployment information does not match the current node
				continue
			}

			// @TODO Must be modified when converting FPGA UUID
			if deployDevice.DeviceUUID == nil && deployDevice.DeviceFilePath == nil {
				continue
			} else if deployDevice.DeviceUUID != nil {
				if *deployDevice.DeviceUUID != *infraData.DeviceUUID {
					// If the DeviceUUID in the infrastructure information and the deployment information do not match
					continue
				}
			} else {
				// @TODO Must be modified when converting FPGA UUID
				if *deployDevice.DeviceFilePath != *infraData.DeviceFilePath {
					// If the DeviceFilePath of the infrastructure information and the deployment information do not match
					continue
				}
			}

			regionList, funcData := createFunctionData(*deployDevice.FunctionTargets)

			// Create RegionInfos for each FunctionTargets(Region)
			for regionIndex := 0; regionIndex < len(*deployDevice.FunctionTargets); regionIndex++ {
				deployFunctionTargets := (*deployDevice.FunctionTargets)[regionIndex]

				var regions examplecomv1.RegionInfo
				var defaultCurrentFunctions int32
				var defaultCurrentCapacity int32
				defaultCurrentFunctions = 0
				defaultCurrentCapacity = 0
				var deviceFilePathString string
				if nil != infraData.DeviceFilePath {
					deviceFilePathString = *infraData.DeviceFilePath
				}
				regions.Name = deployFunctionTargets.RegionName
				regions.Type = deployFunctionTargets.RegionType
				regions.DeviceFilePath = deviceFilePathString
				regions.DeviceUUID = infraData.DeviceUUID
				regions.DeviceType = infraData.DeviceType
				regions.DeviceIndex = infraData.DeviceIndex
				regions.Available = true

				if 0 != len(regionList) &&
					0 != len(funcData[regionList[regionIndex]]) {
					regions.Functions = funcData[regionList[regionIndex]]
					var currentFunction int32 = int32(len(funcData[regionList[regionIndex]]))
					var maxFunction int32 = int32(len(funcData[regionList[regionIndex]]))
					regions.CurrentFunctions = &currentFunction
					regions.Status = examplecomv1.WBRegionStatusReady
					regions.MaxFunctions = &maxFunction
					regions.MaxCapacity = deployFunctionTargets.MaxCapacity
					regions.CurrentCapacity = &defaultCurrentCapacity
				} else {
					if "" != regions.DeviceFilePath {
						regions.Status = examplecomv1.WBRegionStatusNotReady
					} else {
						regions.Status = examplecomv1.WBRegionStatusReady
						regions.CurrentFunctions = &defaultCurrentFunctions
						regions.CurrentCapacity = &defaultCurrentCapacity
					}
					if nil != deployFunctionTargets.MaxFunctions {
						regions.MaxFunctions = deployFunctionTargets.MaxFunctions
					}
					if nil != deployFunctionTargets.MaxCapacity {
						regions.MaxCapacity = deployFunctionTargets.MaxCapacity
					}
				}

				// fmt.Println("regionIndex : ", regionIndex) //debug
				// fmt.Println("regionList : ", regionList[regionIndex]) //debug
				// fmt.Println("regions.Functions : ", regions.Functions) //debug

				// Stored for DeviceInfoCR
				regioninfos = append(regioninfos, regions)
			}
			break
		}
	}
	// fmt.Println("data : ", data) //debug
	return regioninfos
}

// Create Functions and Region information based on the deployment area information
func createFunctionData(
	deployFunctionTargets []regionIndevice) (
	[]string,
	map[string][]examplecomv1.FunctionInfrastruct) {

	var regionList []string
	functions := make(map[string][]examplecomv1.FunctionInfrastruct)

	for targetsIndex := 0; targetsIndex < len(deployFunctionTargets); targetsIndex++ {
		region := deployFunctionTargets[targetsIndex].RegionName
		if nil == deployFunctionTargets[targetsIndex].Functions {
			continue
		}
		deployFunctions := deployFunctionTargets[targetsIndex].Functions

		for deployFunctionIndex := 0; deployFunctionIndex < len(*deployFunctions); deployFunctionIndex++ {

			var listFlag bool

			depFunctions := (*deployFunctions)[deployFunctionIndex]

			listFlag = false

			for _, strRegion := range regionList {
				if strRegion == region {
					listFlag = true
					break
				}
			}

			// fmt.Println("listFlag:", listFlag) //debug
			if listFlag == false {
				regionList = append(regionList, region)
			}

			var funcData examplecomv1.FunctionInfrastruct
			var defaultCurrentDataFlows int32
			var defaultCurrentCapacity int32
			defaultCurrentDataFlows = 0
			defaultCurrentCapacity = 0
			funcData.FunctionIndex = depFunctions.FunctionIndex
			funcData.PartitionName = depFunctions.PartitionName
			funcData.FunctionName = depFunctions.FunctionName
			funcData.Available = true
			funcData.MaxDataFlows = &depFunctions.MaxDataFlows
			funcData.CurrentDataFlows = &defaultCurrentDataFlows
			funcData.MaxCapacity = &depFunctions.MaxCapacity
			funcData.CurrentCapacity = &defaultCurrentCapacity

			functions[region] =
				append(functions[region], funcData)
		}
	}
	// fmt.Println("functions : ", functions) //debug
	return regionList, functions
}
