/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "FPGAFunction/api/v1"

	/* Additional files */
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
	"unsafe"

	// #cgo pkg-config: libdpdk
	// #cgo CFLAGS:  -I/usr/local/include/fpgalib/
	// #cgo CFLAGS:  -I/usr/local/include/fpgalib/dpdk/include/
	// #cgo LDFLAGS: -L/usr/local/lib/fpgalib/ -lfpga
	// #cgo CXXFLAGS: -std=c++11
	// #cgo LDFLAGS: -L. -lstdc++
	// #cgo LDFLAGS: -L. -lpciaccess
	// #cgo LDFLAGS: -L/usr/local/lib/fpgalib/ -lfpgadb
	// #include <liblldma.h>
	// #include <libfpgactl.h>
	// #include <libchain.h>
	// #include <libdmacommon.h>
	// #include <liblogging.h>
	// #include <libshmem.h>
	// #include <libshmem_controller.h>
	// #include <libfpgabs.h>
	// #include <libfunction.h>
	// #include <libdirecttrans.h>
	// #include <libfpgadb.h>
	"C"
	/* Additional files end here */)

// FPGAFunctionReconciler reconciles a FPGAFunction object
type FPGAFunctionReconciler struct {
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

// Device deployment status status type
const (
	OK   = "OK"
	INIT = "INIT"
)

// Overall Status type
const (
	PENDING = "Pending"
	RUNNING = "Running"
)

// Kind type
const (
	KINDFPGAFUNCTION       = "FPGAFunction"
	KINDFPGARECONFIGRATION = "FPGAReconfiguration"
)

const (
	COMPUTERESOURCENAME = "compute-"
)

// Hold own node information
var myNodeName string
var myClusterName string

//+kubebuilder:rbac:groups=example.com,resources=fpgas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=fpgas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=fpgas/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=childbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=childbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=childbs/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=fpgareconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=fpgareconfigurations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=fpgareconfigurations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FPGAFunction object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *FPGAFunctionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var eventKind int32 // 0:Add, 1:Upd,  2:Del
	var retCInt C.int
	var requeueFlag bool
	var childBitstreamCRBase examplecomv1.ChildBs

	breakFlag := false

	reqKind, err := r.GetKind(ctx, req)
	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
			retCInt = 0
		} else {
			logger.Info("Kind Get Error")
		}
	} else if KINDFPGAFUNCTION == reqKind {

		var crData examplecomv1.FPGAFunction

		// Get CR information
		err := r.Get(ctx, req.NamespacedName, &crData)
		if errors.IsNotFound(err) {
			// If CR does not exist
			logger.Info("NotFound to fetch CR")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		} else if err != nil {
			logger.Error(err, "unable to fetch CR")
			return ctrl.Result{}, err
		}

		if myNodeName == crData.Spec.NodeName {

			// Get Event type
			eventKind = r.GetEventKind(&crData)
			if eventKind == CREATE {
				// For creation
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Create", "Create Start")

				var cfgData []byte
				var functionConfigData examplecomv1.FPGAFuncConfig
				var childBitstreamCRData examplecomv1.ChildBs
				var deviceUUID string
				var deviceID C.uint
				var fpgaCRName string
				var fpgaCRData examplecomv1.FPGA

				procStatus := false

				for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

					var functionConfigDataMap map[int32]examplecomv1.FPGAFuncConfig
					functionConfigDataMap = make(map[int32]examplecomv1.FPGAFuncConfig)

					// Get DeviceFilePath information
					if "" != crData.Spec.AcceleratorIDs[0].ID {
						deviceFilePath := crData.Spec.AcceleratorIDs[0].ID
						deviceUUID = strings.ReplaceAll(deviceFilePath, "/dev/xpcie_", "")
					} else {
						logger.Info("FPGAFunction.Spec.AcceleratorIDs[0].ID: " +
							crData.Spec.AcceleratorIDs[0].ID)
						break
					}

					// Get config information
					mainLane := int32(0)
					err = r.getConfigData(ctx, crData.Spec.ConfigName, &cfgData)
					if nil == err {
						err = FunctionConfigDataJsonUnmarshal(&cfgData, &functionConfigData)
						if nil != err {
							logger.Error(err, "unable to unmarshal. ConfigMap="+crData.Spec.ConfigName)
							break
						}
						functionConfigDataMap[mainLane] = functionConfigData
					} else {
						break
					}

					fpgaCRName = "fpga-" + strings.ToLower(deviceUUID) + "-" + crData.Spec.NodeName
					err = r.getFPGACRData(ctx,
						fpgaCRName,
						&fpgaCRData)
					if errors.IsNotFound(err) {
						// CR does not exist
						logger.Info("NotFound to fetch CR")
						break
					} else if err != nil {
						logger.Error(err, "unable to fetch CR")
						break
					}

					err = r.getConfigMapForWriteChildBs(ctx)
					if nil != err {
						logger.Error(err, "unable to unmarshal. getConfigMapForWriteChildBs()")
						break
					}

					if crData.Spec.FunctionIndex == nil &&
						fpgaCRData.Status.ChildBitstreamCRName == nil {

						retCInt, err = r.getChildBsConfig(ctx,
							functionConfigDataMap[mainLane].ParentBitstream.ID,
							functionConfigDataMap[mainLane].ChildBitstream.ID,
							&childBitstreamCRData)
						if nil != err || 0 > retCInt {
							break
						}

						// check ChildBitstream Region Name equal ComputeResource Region Namae.
						// check ComputeResource CurrentCapacity equal Zero.
						capacityIsZero := r.CurrentCapacityIsZero(ctx, childBitstreamCRData, deviceUUID)
						if false == capacityIsZero {
							break
						}

						// create ChildBitstream CR
						err = r.createChildBsCR(ctx,
							functionConfigDataMap,
							false,
							&childBitstreamCRData,
							&fpgaCRData,
							&childBitstreamCRBase)
						if nil != err {
							break
						}

						// update FPGACR for ChildBitstream Name
						err = r.updFPGACR(ctx,
							&fpgaCRData,
							childBitstreamCRData.Spec.ChildBitstreamID)
						if nil != err {
							break
						}

						retCInt = r.WriteFpgaBitstream(ctx,
							deviceUUID,
							nil,
							&deviceID,
							&childBitstreamCRData,
							&fpgaCRData)
						if 0 > retCInt {
							break
						}

						err = r.getChildBsData(ctx,
							fpgaCRName,
							*fpgaCRData.Status.ChildBitstreamID,
							&childBitstreamCRData,
							&childBitstreamCRBase)
						if errors.IsNotFound(err) {
							// CR does not exist
							logger.Info("NotFound to fetch CR")
							break
						} else if err != nil {
							logger.Error(err, "unable to fetch CR")
							break
						}

						err = r.updChildBsCR(ctx,
							&childBitstreamCRData,
							examplecomv1.ChildBsConfiguringParam,
							examplecomv1.ChildBsStatusPreparing,
							&childBitstreamCRBase)
						if nil != err {
							break
						}

						logger.Info("Bitstream file writing has completed successfully.")

						retCInt, breakFlag = r.SetFpgaInfo(ctx,
							deviceID,
							childBitstreamCRData)
						if 0 > retCInt || true == breakFlag {
							break
						}

						err = r.updDeployinfoCM(ctx,
							deviceUUID,
							myNodeName,
							&crData.Spec.FunctionName,
							false,
							&childBitstreamCRData.Spec.Regions)
						if nil != err {
							logger.Error(err, "DeployInfoCM Update Error")
							break
						}

						err = r.getChildBsData(ctx,
							fpgaCRName,
							*fpgaCRData.Status.ChildBitstreamID,
							&childBitstreamCRData,
							&childBitstreamCRBase)
						if errors.IsNotFound(err) {
							// CR does not exist
							logger.Info("NotFound to fetch CR")
							break
						} else if err != nil {
							logger.Error(err, "unable to fetch CR")
							break
						}

						nicNum := 0
						var upChildBitstreamCRData examplecomv1.ChildBs
						for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
							specRegion := childBitstreamCRData.Spec.Regions[regionIndex]
							statusRegion := childBitstreamCRData.Status.Regions[regionIndex]
							var upFunctionData []examplecomv1.ChildBsFunctions
							for functionIndex := 0; functionIndex < len(*specRegion.Modules.Functions); functionIndex++ {
								specFunctionsData := (*specRegion.Modules.Functions)[functionIndex]
								statusFunctionsData := (*statusRegion.Modules.Functions)[functionIndex]
								statusFunctionsData.Parameters = specFunctionsData.Parameters
								upFunctionData = append(upFunctionData, statusFunctionsData)
							}
							statusRegion.Modules.Functions = &upFunctionData
							upChildBitstreamCRData.Status.Regions = append(upChildBitstreamCRData.Status.Regions, statusRegion)
							if nil != specRegion.Modules.Ptu && nil != specRegion.Modules.Ptu.Parameters {
								nicNum++
							}
						}
						childBitstreamCRData.Status = upChildBitstreamCRData.Status

						if 0 == nicNum {
							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsReady,
								examplecomv1.ChildBsStatusReady,
								&childBitstreamCRBase)
							if nil != err {
								break
							} else {
								logger.Info("This Child-Bitstream has no NIC.")
							}
						} else {
							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsNoConfigureNetwork,
								examplecomv1.ChildBsStatusPreparing,
								&childBitstreamCRBase)
							if nil != err {
								break
							} else {
								logger.Info("Application parameters settings for all lanes has completed successfully.")
							}
						}
					}

					err = r.getChildBsData(ctx,
						fpgaCRName,
						*fpgaCRData.Status.ChildBitstreamID,
						&childBitstreamCRData,
						&childBitstreamCRBase)
					if errors.IsNotFound(err) {
						// CR does not exist
						logger.Info("NotFound to fetch CR")
						break
					} else if err != nil {
						logger.Error(err, "unable to fetch CR")
						break
					} else {
						logger.Info("Success to Get ChildBitstreamData.")
						logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State))
						logger.Info("ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
					}

					if examplecomv1.ChildBsReady == childBitstreamCRData.Status.State {

						var fpgafuncRxData examplecomv1.RxTxData
						var fpgafuncTxData examplecomv1.RxTxData
						var functionKernelID int32
						var chainID int32
						var functionChannelID int32
						err = r.allocateFPGAResource(ctx,
							req,
							&crData,
							&childBitstreamCRData,
							&fpgafuncRxData,
							&fpgafuncTxData,
							&functionKernelID,
							&chainID,
							&functionChannelID,
							&childBitstreamCRBase)
						if nil != err {
							logger.Info("wait because ConnectionCR is not found.")
							return ctrl.Result{Requeue: true}, nil
						}

						// Summary of information for updating the Status area
						var statused examplecomv1.AccStatuses
						var acceleratorStatuses examplecomv1.AccStatusesByDevice
						var crStatusString string
						crStatusString = OK

						partitionName := strconv.Itoa(int(functionKernelID))
						statused.AcceleratorID = &crData.Spec.AcceleratorIDs[0].ID
						statused.Status = &crStatusString
						if nil != crData.Spec.FunctionIndex {
							acceleratorStatuses.PartitionName = &partitionName
						} else {
							acceleratorStatuses.PartitionName = crData.Spec.AcceleratorIDs[0].PartitionName
						}
						acceleratorStatuses.Statused = append(acceleratorStatuses.Statused, statused)

						// Status column update
						crData.Status.DataFlowRef = crData.Spec.DataFlowRef
						crData.Status.FunctionName = crData.Spec.FunctionName
						crData.Status.SharedMemory = crData.Spec.SharedMemory
						crData.Status.ParentBitstreamName = functionConfigDataMap[mainLane].ParentBitstream.ID
						crData.Status.ChildBitstreamName = functionConfigDataMap[mainLane].ChildBitstream.ID
						crData.Status.FunctionChannelID = functionChannelID
						crData.Status.Rx = fpgafuncRxData
						crData.Status.Tx = fpgafuncTxData
						crData.Status.FrameworkKernelID = chainID
						crData.Status.FunctionKernelID = functionKernelID
						crData.Status.PtuKernelID = chainID

						crData.Status.AcceleratorStatuses = append(crData.Status.AcceleratorStatuses, acceleratorStatuses)

						r.UpdCustomResource(ctx, &crData, RUNNING)
						procStatus = true
					} else if examplecomv1.ChildBsError == childBitstreamCRData.Status.State {
						logger.Info("ChildBitstream State Is requeue end.")
						logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
							", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
						procStatus = true
					} else {
						logger.Info("ChildBitstream State Is wait.")
						logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
							", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
						return ctrl.Result{Requeue: true}, nil
					}
				}
				if procStatus == false && requeueFlag == false {
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Create", "Create Err")
				} else if true == procStatus && false == requeueFlag {
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Create", "Create End")
				} else {
					logger.Info("RequeueFlag is True.")
				}
			} else if eventKind == UPDATE {
				// In case of update
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Update", "Update Start")
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Update", "Update End")
			} else if eventKind == DELETE {
				// In case of deletion
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete Start")

				var childBitstreamCRData examplecomv1.ChildBs
				var deviceUUID string
				var fpgaCRName string
				var fpgaCRData examplecomv1.FPGA

				procStatus := false

				for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

					// Get DeviceFilePath information
					if "" != crData.Spec.AcceleratorIDs[0].ID {
						deviceFilePath := crData.Spec.AcceleratorIDs[0].ID
						deviceUUID = strings.ReplaceAll(deviceFilePath, "/dev/xpcie_", "")
					} else {
						logger.Info("FPGAFunction.Spec.AcceleratorIDs[0].ID: " +
							crData.Spec.AcceleratorIDs[0].ID)
						break
					}

					fpgaCRName = "fpga-" + strings.ToLower(deviceUUID) + "-" + crData.Spec.NodeName
					err = r.getFPGACRData(ctx,
						fpgaCRName,
						&fpgaCRData)
					if errors.IsNotFound(err) {
						// CR does not exist
						logger.Info("NotFound to fetch CR")
						break
					} else if err != nil {
						logger.Error(err, "unable to fetch CR")
						break
					}

					if nil != fpgaCRData.Status.ChildBitstreamCRName {
						err = r.getChildBsData(ctx,
							fpgaCRName,
							*fpgaCRData.Status.ChildBitstreamID,
							&childBitstreamCRData,
							&childBitstreamCRBase)
						if errors.IsNotFound(err) {
							// CR does not exist
							logger.Info("NotFound to fetch CR")
							break
						} else if err != nil {
							logger.Error(err, "unable to fetch CR")
							break
						}

						// FPGAResource free
						r.freeFPGAResource(ctx, req, &crData, &childBitstreamCRData, &childBitstreamCRBase)
					}

					// Delete the Finalizer statement.
					err = r.DelCustomResource(ctx, &crData)
					if err != nil {
						r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete Err")
						return ctrl.Result{}, client.IgnoreNotFound(err)
					}

					procStatus = true
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete End")
				}
				if false == procStatus {
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete Err")
				}
			}
		}
	} else if KINDFPGARECONFIGRATION == reqKind {

		var crData examplecomv1.FPGAReconfiguration

		// Get CR information
		err := r.Get(ctx, req.NamespacedName, &crData)
		if errors.IsNotFound(err) {
			// If CR does not exist
			logger.Info("NotFound to fetch CR")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		} else if err != nil {
			logger.Error(err, "unable to fetch CR")
			return ctrl.Result{}, err
		}
		logger.Info("Evevt Start Name=" + req.Name)

		if myNodeName == crData.Spec.NodeName {

			// Get Event type
			eventKind = r.GetEventKindFPGAReconfiguration(&crData)
			if eventKind == CREATE {
				// For creation
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Create", "Create Start")

				var cfgData []byte
				var functionConfigData examplecomv1.FPGAFuncConfig
				var childBitstreamCRData examplecomv1.ChildBs
				var deviceUUID string
				var deviceID C.uint
				var fpgaCRName string
				var fpgaCRData examplecomv1.FPGA

				procStatus := false

				for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

					var functionConfigDataMap map[int32]examplecomv1.FPGAFuncConfig
					functionConfigDataMap = make(map[int32]examplecomv1.FPGAFuncConfig)

					// Get DeviceFilePath information
					if "" != crData.Spec.DeviceFilePath {
						deviceFilePath := crData.Spec.DeviceFilePath
						deviceUUID = strings.ReplaceAll(deviceFilePath, "/dev/xpcie_", "")
					} else {
						logger.Info("FPGAReconfiguration.Spec.DeviceFilePath: " +
							crData.Spec.DeviceFilePath)
						break
					}

					mainLane := int32(-1)
					for CfgIndex := 0; len(crData.Spec.ConfigNames) > CfgIndex; CfgIndex++ {
						// Get config information
						err = r.getConfigData(ctx, crData.Spec.ConfigNames[CfgIndex].ConfigName, &cfgData)
						if nil == err {
							err = FunctionConfigDataJsonUnmarshal(&cfgData, &functionConfigData)
							if nil != err {
								logger.Error(err, "unable to unmarshal. ConfigMap="+
									crData.Spec.ConfigNames[CfgIndex].ConfigName)
								break
							}
							functionConfigDataMap[crData.Spec.ConfigNames[CfgIndex].LaneIndex] = functionConfigData
							if 0 > mainLane {
								mainLane = crData.Spec.ConfigNames[CfgIndex].LaneIndex
							}
						} else {
							break
						}
					}
					if nil != err {
						break
					}

					fpgaResetFlag := false
					childbsResetFlag := false
					if nil != crData.Spec.FPGAResetFlag {
						fpgaResetFlag = *crData.Spec.FPGAResetFlag
					}
					if nil != crData.Spec.ChildBsResetFlag {
						childbsResetFlag = *crData.Spec.ChildBsResetFlag
					}

					if false == childbsResetFlag && 0 > mainLane {
						logger.Error(err, "ConfigName parameter not found.")
						break
					}

					fpgaCRName = "fpga-" + strings.ToLower(deviceUUID) + "-" + crData.Spec.NodeName
					err = r.getFPGACRData(ctx,
						fpgaCRName,
						&fpgaCRData)
					if errors.IsNotFound(err) {
						// CR does not exist
						logger.Error(err, "not found to fetch FPGACR")
						break
					} else if err != nil {
						logger.Error(err, "unable to fetch FPGACR")
						break
					}

					if fpgaResetFlag == true && childbsResetFlag == true {
						logger.Error(err, "request reset flag exclusive error.")
						break

					} else if fpgaResetFlag == false &&
						childbsResetFlag == false &&
						fpgaCRData.Status.ChildBitstreamCRName != nil {
						logger.Info("already set FPGA.")

						err = r.getConfigMapForWriteChildBs(ctx)
						if nil != err {
							logger.Error(err, "unable to unmarshal. getConfigMapForWriteChildBs()")
							break
						}

						err = r.getChildBsData(ctx,
							fpgaCRName,
							*fpgaCRData.Status.ChildBitstreamID,
							&childBitstreamCRData,
							&childBitstreamCRBase)
						if errors.IsNotFound(err) {
							logger.Info("not found ChildBitstreamCR.")
							break
						} else if err != nil {
							logger.Error(err, "unable to fetch ChildBitstreamCR")
							break
						}

						if examplecomv1.ChildBsNoConfigureNetwork != childBitstreamCRData.Status.State &&
							examplecomv1.ChildBsConfiguringNetwork != childBitstreamCRData.Status.State &&
							examplecomv1.ChildBsReady != childBitstreamCRData.Status.State &&
							examplecomv1.ChildBsError != childBitstreamCRData.Status.State {
							logger.Info("ChildBitstream State Is Unmatch.")
							logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
								", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
							break
						}

					} else if fpgaResetFlag == false &&
						childbsResetFlag == false &&
						fpgaCRData.Status.ChildBitstreamCRName == nil {

						err = r.getConfigMapForWriteChildBs(ctx)
						if nil != err {
							logger.Error(err, "unable to unmarshal. getConfigMapForWriteChildBs()")
							break
						}

						regiontypeIsMatch := r.RegionTypeIsMatch(ctx,
							crData,
							functionConfigDataMap[mainLane].ChildBitstream.ID,
							deviceUUID)
						if false == regiontypeIsMatch {
							logger.Error(err, "regiontype unmatch.")
							break
						}

						retCInt, err = r.getChildBsConfig(ctx,
							functionConfigDataMap[mainLane].ParentBitstream.ID,
							functionConfigDataMap[mainLane].ChildBitstream.ID,
							&childBitstreamCRData)
						if nil != err || 0 > retCInt {
							break
						}

						// check FPGAReconfigurationCR ConfigNames num equal ChildBitstream Region num.
						// check FPGAReconfigurationCR ConfigNames.LaneIndex not duplicate.
						// check FPGAReconfigurationCR ConfigNames.LaneIndex equal ChildBitstream Region Name.
						confignamesIsMatch := r.RegionNamesIsMatch(ctx, crData, childBitstreamCRData)
						if false == confignamesIsMatch {
							break
						}

						// check ChildBitstream Region Name equal ComputeResource Region Namae.
						// check ComputeResource CurrentCapacity equal Zero.
						capacityIsZero := r.CurrentCapacityIsZero(ctx, childBitstreamCRData, deviceUUID)
						if false == capacityIsZero {
							break
						}

						// create ChildBitstream CR
						err = r.createChildBsCR(ctx,
							functionConfigDataMap,
							true,
							&childBitstreamCRData,
							&fpgaCRData,
							&childBitstreamCRBase)
						if nil != err {
							break
						}
						// update FPGACR for ChildBitstream Name
						err = r.updFPGACR(ctx,
							&fpgaCRData,
							childBitstreamCRData.Spec.ChildBitstreamID)
						if nil != err {
							break
						}

					} else if (true == fpgaResetFlag ||
						true == childbsResetFlag) &&
						nil == fpgaCRData.Status.ChildBitstreamCRName {
						// CR does not exist
						if true == fpgaResetFlag {
							logger.Info("already reset FPGA.")
						} else {
							logger.Info("not set FPGA.")
						}
						break

					} else if (true == fpgaResetFlag ||
						true == childbsResetFlag) &&
						nil != fpgaCRData.Status.ChildBitstreamCRName {

						err = r.getConfigMapForWriteChildBs(ctx)
						if nil != err {
							logger.Error(err, "unable to unmarshal. getConfigMapForWriteChildBs()")
							break
						}

						err = r.getChildBsData(ctx,
							fpgaCRName,
							*fpgaCRData.Status.ChildBitstreamID,
							&childBitstreamCRData,
							&childBitstreamCRBase)
						if errors.IsNotFound(err) {
							// CR does not exist
							if true == fpgaResetFlag {
								logger.Info("already reset FPGA.")
							} else {
								logger.Info("not set FPGA.")
							}
							break
						} else if err != nil {
							logger.Error(err, "unable to fetch ChildBitstreamCR")
							break
						}
						if true == fpgaResetFlag {
							if examplecomv1.ChildBsNotStopNetworkModule != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsStoppingNetworkModule != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsNotWriteBsfile != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsWritingBsfile != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsConfiguringParam != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsReady != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsError != childBitstreamCRData.Status.State {
								logger.Info("ChildBitstream State Is Unmatch.")
								logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
									", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
								break
							}
						} else {
							if examplecomv1.ChildBsNotStopNetworkModule != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsStoppingNetworkModule != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsNotWriteBsfile != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsWritingBsfile != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsConfiguringParam != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsNoConfigureNetwork != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsConfiguringNetwork != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsReady != childBitstreamCRData.Status.State &&
								examplecomv1.ChildBsError != childBitstreamCRData.Status.State {
								logger.Info("ChildBitstream State Is Unmatch.")
								logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
									", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
								break
							}
						}
						if examplecomv1.ChildBsReady == childBitstreamCRData.Status.State ||
							examplecomv1.ChildBsError == childBitstreamCRData.Status.State {

							// check ChildBitstream Region Name equal ComputeResource Region Namae.
							// check ComputeResource CurrentCapacity equal Zero.
							capacityIsZero := r.CurrentCapacityIsZero(ctx, childBitstreamCRData, deviceUUID)
							if false == capacityIsZero {
								break
							}

							err = r.updDeployinfoCM(ctx,
								deviceUUID,
								myNodeName,
								nil,
								true,
								&childBitstreamCRData.Spec.Regions)
							if nil != err {
								logger.Error(err, "DeployInfoCM Update Error")
								break
							}

							err = r.getChildBsData(ctx,
								fpgaCRName,
								*fpgaCRData.Status.ChildBitstreamID,
								&childBitstreamCRData,
								&childBitstreamCRBase)
							if errors.IsNotFound(err) {
								// CR does not exist
								logger.Info("NotFound to fetch CR")
								break
							} else if err != nil {
								logger.Error(err, "unable to fetch CR")
								break
							}
							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsStoppingModule,
								examplecomv1.ChildBsStatusPreparing,
								&childBitstreamCRBase)
							if nil != err {
								break
							}
						}

					} else {
						logger.Error(err, "request timming unmatch."+
							" FPGAResetFlag="+strconv.FormatBool(fpgaResetFlag)+
							" ChildBsResetFlag="+strconv.FormatBool(childbsResetFlag)+
							" ChildBitstreamCRName="+*fpgaCRData.Status.ChildBitstreamCRName)
						if nil != fpgaCRData.Status.ChildBitstreamCRName {
							logger.Error(err, "request timming unmatch."+
								" ChildBitstreamCRName="+*fpgaCRData.Status.ChildBitstreamCRName)
						}
						break
					}

					if examplecomv1.ChildBsStoppingModule == childBitstreamCRData.Status.State {

						_ = r.stopFpgaModule(ctx,
							deviceUUID,
							&childBitstreamCRData)

						logger.Info("Stop FPGA Module has completed successfully.")

						nicNum := 0
						for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
							region := childBitstreamCRData.Spec.Regions[regionIndex]
							if nil != region.Modules.Ptu && nil != region.Modules.Ptu.Parameters {
								nicNum++
							}
						}

						err = r.getChildBsData(ctx,
							fpgaCRName,
							*fpgaCRData.Status.ChildBitstreamID,
							&childBitstreamCRData,
							&childBitstreamCRBase)
						if errors.IsNotFound(err) {
							// CR does not exist
							logger.Info("NotFound to fetch CR")
							break
						} else if err != nil {
							logger.Error(err, "unable to fetch CR")
							break
						}
						if 0 == nicNum {
							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsNotWriteBsfile,
								examplecomv1.ChildBsStatusPreparing,
								&childBitstreamCRBase)
						} else {
							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsNotStopNetworkModule,
								examplecomv1.ChildBsStatusPreparing,
								&childBitstreamCRBase)
							procStatus = true
						}
						if nil != err {
							break
						}
					}

					if true == fpgaResetFlag {
						if examplecomv1.ChildBsNotWriteBsfile == childBitstreamCRData.Status.State {

							err = r.getChildBsData(ctx,
								fpgaCRName,
								*fpgaCRData.Status.ChildBitstreamID,
								&childBitstreamCRData,
								&childBitstreamCRBase)
							if errors.IsNotFound(err) {
								// CR does not exist
								logger.Info("NotFound to fetch CR")
								break
							} else if err != nil {
								logger.Error(err, "unable to fetch CR")
								break
							}
							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsReconfiguring,
								examplecomv1.ChildBsStatusPreparing,
								&childBitstreamCRBase)
							if nil != err {
								break
							}

							// update FPGACR for clear ChildBitstream Name
							err = r.updFPGACR(ctx,
								&fpgaCRData,
								nil)
							if nil != err {
								break
							}

							functionConfigData = functionConfigDataMap[mainLane]

							// FPGA write Initialize Only.
							retCInt = r.WriteFpgaBitstream(ctx,
								deviceUUID,
								&functionConfigData,
								&deviceID,
								&childBitstreamCRData,
								&fpgaCRData)
							if 0 > retCInt {
								break
							}
							logger.Info("Bitstream file initalize has completed successfully.")

							// delete ChildBitstream CR
							err = r.deleteChildBsCR(ctx, &childBitstreamCRData)
							if nil != err {
								break
							}
							r.UpdCRFPGAReconfiguration(ctx,
								&crData,
								examplecomv1.FPGARECONFSTATUSSUCCEEDED)
							procStatus = true
						} else if examplecomv1.ChildBsReconfiguring == childBitstreamCRData.Status.State {
							logger.Info("ChildBitstream State Is requeue end.")
							logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
								", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
							procStatus = true
						} else {
							logger.Info("ChildBitstream State Is wait.")
							logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
								", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
							return ctrl.Result{Requeue: true}, nil
						}
					} else {

						if examplecomv1.ChildBsNotWriteBsfile == childBitstreamCRData.Status.State ||
							examplecomv1.ChildBsWritingBsfile == childBitstreamCRData.Status.State {

							if true == childbsResetFlag {
								err = r.getChildBsData(ctx,
									fpgaCRName,
									*fpgaCRData.Status.ChildBitstreamID,
									&childBitstreamCRData,
									&childBitstreamCRBase)
								if errors.IsNotFound(err) {
									// CR does not exist
									logger.Info("NotFound to fetch CR")
									break
								} else if err != nil {
									logger.Error(err, "unable to fetch CR")
									break
								}

								err = r.updChildBsCR(ctx,
									&childBitstreamCRData,
									examplecomv1.ChildBsWritingBsfile,
									examplecomv1.ChildBsStatusPreparing,
									&childBitstreamCRBase)
								if nil != err {
									break
								}
							}

							// FPGA write Bitstream & update ChildBitstream CR
							retCInt = r.WriteFpgaBitstream(ctx,
								deviceUUID,
								nil,
								&deviceID,
								&childBitstreamCRData,
								&fpgaCRData)
							if 0 > retCInt {
								break
							}
							err = r.getChildBsData(ctx,
								fpgaCRName,
								*fpgaCRData.Status.ChildBitstreamID,
								&childBitstreamCRData,
								&childBitstreamCRBase)
							if errors.IsNotFound(err) {
								// CR does not exist
								logger.Info("NotFound to fetch CR")
								break
							} else if err != nil {
								logger.Error(err, "unable to fetch CR")
								break
							}

							err = r.updChildBsCR(ctx,
								&childBitstreamCRData,
								examplecomv1.ChildBsConfiguringParam,
								examplecomv1.ChildBsStatusPreparing,
								&childBitstreamCRBase)
							if nil != err {
								break
							}
							logger.Info("Bitstream file writing has completed successfully.")

							retCInt, breakFlag = r.SetFpgaInfo(ctx,
								deviceID,
								childBitstreamCRData)
							if 0 > retCInt || true == breakFlag {
								break
							}
							logger.Info("fpga setting has completed successfully.")

							err = r.getChildBsData(ctx,
								fpgaCRName,
								*fpgaCRData.Status.ChildBitstreamID,
								&childBitstreamCRData,
								&childBitstreamCRBase)
							if errors.IsNotFound(err) {
								// CR does not exist
								logger.Info("NotFound to fetch CR")
								break
							} else if err != nil {
								logger.Error(err, "unable to fetch CR")
								break
							}

							err = r.updDeployinfoCM(ctx,
								deviceUUID,
								myNodeName,
								nil,
								false,
								&childBitstreamCRData.Spec.Regions)
							if nil != err {
								logger.Error(err, "DeployInfoCM Update Error")
								break
							}
							logger.Info("deployinfo configmap update completed successfully.")

							r.reallocateFPGAResource(ctx, &childBitstreamCRData)
							logger.Info("fpga resource reallocation has completed successfully.")

							nicNum := 0
							var upChildBitstreamCRData examplecomv1.ChildBs
							for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
								specRegion := childBitstreamCRData.Spec.Regions[regionIndex]
								statusRegion := childBitstreamCRData.Status.Regions[regionIndex]
								var upFunctionData []examplecomv1.ChildBsFunctions
								for functionIndex := 0; functionIndex < len(*specRegion.Modules.Functions); functionIndex++ {
									specFunctionsData := (*specRegion.Modules.Functions)[functionIndex]
									statusFunctionsData := (*statusRegion.Modules.Functions)[functionIndex]
									statusFunctionsData.Parameters = specFunctionsData.Parameters
									upFunctionData = append(upFunctionData, statusFunctionsData)
								}
								statusRegion.Modules.Functions = &upFunctionData
								upChildBitstreamCRData.Status.Regions = append(upChildBitstreamCRData.Status.Regions, statusRegion)
								if nil != specRegion.Modules.Ptu && nil != specRegion.Modules.Ptu.Parameters {
									nicNum++
								}
							}
							childBitstreamCRData.Status = upChildBitstreamCRData.Status

							if 0 == nicNum {
								err = r.updChildBsCR(ctx,
									&childBitstreamCRData,
									examplecomv1.ChildBsReady,
									examplecomv1.ChildBsStatusReady,
									&childBitstreamCRBase)
								if nil != err {
									break
								} else {
									logger.Info("This Child-Bitstream has no NIC.")
								}
							} else {
								err = r.updChildBsCR(ctx,
									&childBitstreamCRData,
									examplecomv1.ChildBsNoConfigureNetwork,
									examplecomv1.ChildBsStatusPreparing,
									&childBitstreamCRBase)
								if nil != err {
									break
								} else {
									logger.Info("Application parameters settings for all lanes has completed successfully.")
								}
							}
							r.UpdCRFPGAReconfiguration(ctx,
								&crData,
								examplecomv1.FPGARECONFSTATUSSUCCEEDED)
							procStatus = true
						} else if examplecomv1.ChildBsNoConfigureNetwork == childBitstreamCRData.Status.State ||
							examplecomv1.ChildBsConfiguringNetwork == childBitstreamCRData.Status.State {
							logger.Info("ChildBitstream State Is requeue end.")
							logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
								", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
							procStatus = true
						} else if examplecomv1.ChildBsReady == childBitstreamCRData.Status.State ||
							examplecomv1.ChildBsError == childBitstreamCRData.Status.State {
							logger.Info("ChildBitstream State Is requeue end.")
							logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
								", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
						} else {
							logger.Info("ChildBitstream State Is wait.")
							logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
								", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
							return ctrl.Result{Requeue: true}, nil
						}
					}
				}
				if false == procStatus {
					r.UpdCRFPGAReconfiguration(ctx,
						&crData,
						examplecomv1.FPGARECONFSTATUSFAILED)
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Create", "Create Err")
				} else {
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Create", "Create End")
				}

			} else if eventKind == UPDATE {
				// In case of update
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Update", "Update Start")
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Update", "Update End")

			} else if eventKind == DELETE {
				// In case of deletion
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete Start")

				// Delete the Finalizer statement.
				err = r.DelCRFPGAReconfiguration(ctx, &crData)
				if err != nil {
					r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete Err")
					return ctrl.Result{}, client.IgnoreNotFound(err)
				}
				r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete End")
			}
		}
		logger.Info("Evevt End Name=" + req.Name)
	} else {
		logger.Info("Kind unknown. reqKind=" + reqKind)
	}

	if nil != err {
		return ctrl.Result{}, err
	} else if requeueFlag == true {
		return ctrl.Result{Requeue: requeueFlag}, nil
	} else if 0 > retCInt {
		return ctrl.Result{}, nil
	} else if false == breakFlag {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *FPGAFunctionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.FPGAFunction{}).
		Watches(&examplecomv1.FPGAReconfiguration{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

func (r *FPGAFunctionReconciler) GetfinalizerName(pCRData *examplecomv1.FPGAFunction) string {
	// Value to set in the finalizer
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
}

func (r *FPGAFunctionReconciler) GetfinalizerNameFPGAReconfiguration(pCRData *examplecomv1.FPGAReconfiguration) string {
	// Value to set in the finalizer
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
}

func (r *FPGAFunctionReconciler) GetKind(ctx context.Context, req ctrl.Request) (string, error) {
	logger := log.FromContext(ctx)
	var err error
	retKind := ""

	var KindList []string = []string{KINDFPGAFUNCTION, KINDFPGARECONFIGRATION}

	kcrData := &unstructured.Unstructured{}

	for n := 0; n < len(KindList); n++ {
		kcrData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    KindList[n],
		})
		err = r.Get(ctx, client.ObjectKey{
			Namespace: req.Namespace,
			Name:      req.Name}, kcrData)
		if errors.IsNotFound(err) {
			continue
		} else if err != nil {
			logger.Info("unable to fetch CR")
			break
		} else {
			retKind = KindList[n]
			break
		}

	}
	return retKind, err
}

func (r *FPGAFunctionReconciler) GetEventKind(pCRData *examplecomv1.FPGAFunction) int32 {
	var eventKind int32
	eventKind = UPDATE
	// Whether or not there is a deletion timestamp
	if !pCRData.ObjectMeta.DeletionTimestamp.IsZero() {
		eventKind = DELETE
	} else if !controllerutil.ContainsFinalizer(pCRData, r.GetfinalizerName(pCRData)) {
		// Whether or not Finalizer is written
		eventKind = CREATE
	}
	return eventKind
}

func (r *FPGAFunctionReconciler) GetEventKindFPGAReconfiguration(pCRData *examplecomv1.FPGAReconfiguration) int32 {
	var eventKind int32
	eventKind = UPDATE
	// Whether or not there is a deletion timestamp
	if !pCRData.ObjectMeta.DeletionTimestamp.IsZero() {
		eventKind = DELETE
	} else if !controllerutil.ContainsFinalizer(pCRData, r.GetfinalizerNameFPGAReconfiguration(pCRData)) {
		// Whether or not Finalizer is written
		eventKind = CREATE
	}
	return eventKind
}

func (r *FPGAFunctionReconciler) RegionTypeIsMatch(ctx context.Context,
	crData examplecomv1.FPGAReconfiguration,
	childBitstreamID string,
	deviceUUID string) bool {
	retBool := true
	logger := log.FromContext(ctx)

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		logger.Info("len(crData.Spec.ConfigNames)=" + strconv.Itoa(len(crData.Spec.ConfigNames)))
		for crConfigIndex := 0; len(crData.Spec.ConfigNames) > crConfigIndex; crConfigIndex++ {
			laneName := "lane" + strconv.Itoa(int(crData.Spec.ConfigNames[crConfigIndex].LaneIndex))
			laneIndex := strconv.Itoa(int(crData.Spec.ConfigNames[crConfigIndex].LaneIndex))
			regionuniType := ""
			for ruiIndex := 0; ruiIndex < len(gRegionUniqueInfo); ruiIndex++ {
				regionInfo := gRegionUniqueInfo[ruiIndex]
				if regionInfo.SubDeviceSpecRef != childBitstreamID {
					continue
				}
				for targetIndex := 0; targetIndex < len(regionInfo.FunctionTargets); targetIndex++ {
					target := regionInfo.FunctionTargets[targetIndex]
					if laneName == target.RegionName {
						logger.Info("target.RegionType=" + target.RegionType)
						regionuniType = target.RegionType
						break
					}
				}
				if "" == regionuniType {
					break
				}
			}
			if "" == regionuniType {
				retBool = false
				break
			}
			predeterminedregionType := ""
			for pdriIndex := 0; pdriIndex < len(gPreDeterminedRegionInfo); pdriIndex++ {
				predeterminedregionInfo := gPreDeterminedRegionInfo[pdriIndex]
				if predeterminedregionInfo.NodeName != crData.Spec.NodeName {
					continue
				}
				if predeterminedregionInfo.DeviceUUID != deviceUUID {
					continue
				}
				if predeterminedregionInfo.SubDeviceSpecRef == laneIndex {
					predeterminedregionType = predeterminedregionInfo.RegionType
					break
				}
			}
			if "" == predeterminedregionType {
				retBool = false
				break
			}
			if predeterminedregionType != regionuniType {
				retBool = false
				break
			}
		}
	}
	return retBool
}

func (r *FPGAFunctionReconciler) RegionNamesIsMatch(ctx context.Context,
	crData examplecomv1.FPGAReconfiguration,
	childBitstreamCRData examplecomv1.ChildBs) bool {
	logger := log.FromContext(ctx)
	retBool := true
	// check FPGAReconfigurationCR ConfigNames num equal ChildBitstream Region num.
	// check FPGAReconfigurationCR ConfigNames.LaneIndex not duplicate.
	// check FPGAReconfigurationCR ConfigNames.LaneIndex equal ChildBitstream Region Name.
	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		if len(crData.Spec.ConfigNames) != len(childBitstreamCRData.Spec.Regions) {
			logger.Info("ChildBitstreamCR Regions num not equal ConfigNames num")
			retBool = false
			break
		}
		regionsCheckBit := 0
		regionsNGFlag := false
		for crConfigIndex := 0; len(crData.Spec.ConfigNames) > crConfigIndex; crConfigIndex++ {
			if 0 != (regionsCheckBit & (1 << crData.Spec.ConfigNames[crConfigIndex].LaneIndex)) {
				logger.Info("ConfigNames LaneIndex is duplicate.")
				regionsNGFlag = true
				break
			}
			regionsCheckBit |= (1 << int(crData.Spec.ConfigNames[crConfigIndex].LaneIndex))
			laneName := "lane" + strconv.Itoa(int(crData.Spec.ConfigNames[crConfigIndex].LaneIndex))
			laneCheck := false
			for childregionIndex := 0; len(childBitstreamCRData.Spec.Regions) > childregionIndex; childregionIndex++ {
				if *childBitstreamCRData.Spec.Regions[childregionIndex].Name == laneName {
					laneCheck = true
					break
				}
			}
			if false == laneCheck {
				logger.Info("ConfigNames LaneIndex not found in CheildBitstreamCR Regions.")
				regionsNGFlag = true
				break

			}
		}
		if true == regionsNGFlag {
			retBool = false
			break
		}
	}
	return retBool
}

func (r *FPGAFunctionReconciler) CurrentCapacityIsZero(ctx context.Context,
	childBitstreamCRData examplecomv1.ChildBs,
	deviceUUID string) bool {
	logger := log.FromContext(ctx)
	retBool := true
	var ComputeResourceData examplecomv1.ComputeResource
	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		err = r.getComputeResourceCR(ctx,
			&ComputeResourceData)
		if nil != err {
			retBool = false
			break
		}

		// check ChildBitstream Region Name equal ComputeResource Region Namae.
		// check ComputeResource CurrentCapacity equal Zero.
		for childregionIndex := 0; len(childBitstreamCRData.Spec.Regions) > childregionIndex; childregionIndex++ {
			checkFlag := false
			for regionIndex := 0; len(ComputeResourceData.Spec.Regions) > regionIndex; regionIndex++ {
				if deviceUUID !=
					*ComputeResourceData.Spec.Regions[regionIndex].DeviceUUID {
					continue
				}
				if ComputeResourceData.Spec.Regions[regionIndex].Name !=
					*childBitstreamCRData.Spec.Regions[childregionIndex].Name {
					continue
				}
				checkFlag = true
				if nil != ComputeResourceData.Spec.Regions[regionIndex].CurrentCapacity {
					if 0 != int(*ComputeResourceData.Spec.Regions[regionIndex].CurrentCapacity) {
						logger.Info("ComputeResource CurrentCapacity not equal Zero")
						retBool = false
						break
					}
				}
			}
			if false == checkFlag {
				logger.Info("ComputeResource Region unmatch. CurrentCapacity Check Error")
				retBool = false
				break
			}
			if false == retBool {
				break
			}
		}
	}
	return retBool
}

func (r *FPGAFunctionReconciler) stopFpgaModule(ctx context.Context,
	deviceUUID string,
	childBitstreamCRData *examplecomv1.ChildBs) C.int {
	logger := log.FromContext(ctx)

	var deviceID C.uint
	var retCInt C.int
	var funcRetCInt C.int

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger.Info("Start to stop FPGA Module.")

		fpgaDeviceUUIDCString := C.CString(deviceUUID)
		defer C.free(unsafe.Pointer(fpgaDeviceUUIDCString))
		retCInt = C.fpga_get_dev_id(fpgaDeviceUUIDCString, &deviceID)
		if 0 > retCInt {
			logger.Info("fpga_get_dev_id() err = " +
				strconv.Itoa(int(retCInt)))
			break
		} else {
			logger.Info("fpga_get_dev_id() ret = " +
				strconv.Itoa(int(retCInt)))
		}

		for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {

			region := childBitstreamCRData.Spec.Regions[regionIndex]
			nameSlice := []rune(*region.Name)
			lane, _ := strconv.Atoi(string(nameSlice[len(nameSlice)-1:]))

			funcRetCInt = C.fpga_chain_stop(deviceID, C.uint(lane))
			if 0 > funcRetCInt {
				logger.Info("fpga_chain_stop() err = " +
					strconv.Itoa(int(funcRetCInt)))
				if 0 == retCInt {
					retCInt = funcRetCInt
				}
			}
			funcRetCInt = C.fpga_direct_stop(deviceID, C.uint(lane))
			if 0 > funcRetCInt {
				logger.Info("fpga_direct_stop() err = " +
					strconv.Itoa(int(funcRetCInt)))
				if 0 == retCInt {
					retCInt = funcRetCInt
				}
			}
			funcRetCInt = C.fpga_function_finish(deviceID, C.uint(lane), nil)
			if 0 > funcRetCInt {
				logger.Info("fpga_function_finish() err = " +
					strconv.Itoa(int(funcRetCInt)))
				if 0 == retCInt {
					retCInt = funcRetCInt
				}
			}
		}
		funcRetCInt = C.fpga_ref_cleanup(deviceID)
		if 0 > funcRetCInt {
			logger.Info("fpga_ref_cleanup() ret = " +
				strconv.Itoa(int(funcRetCInt)))
			if 0 == retCInt {
				retCInt = funcRetCInt
			}
		}
	}
	return retCInt
}

func (r *FPGAFunctionReconciler) WriteFpgaBitstream(ctx context.Context,
	deviceUUID string,
	functionConfigData *examplecomv1.FPGAFuncConfig,
	deviceID *C.uint,
	childBitstreamCRData *examplecomv1.ChildBs,
	fpgaCRData *examplecomv1.FPGA) C.int {
	logger := log.FromContext(ctx)

	var retCInt C.int

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger.Info("Start to write Child-Bitstream.")

		fpgaDeviceUUIDCString := C.CString(deviceUUID)
		defer C.free(unsafe.Pointer(fpgaDeviceUUIDCString))
		retCInt = C.fpga_get_dev_id(fpgaDeviceUUIDCString, deviceID)
		if 0 > retCInt {
			logger.Info("fpga_get_dev_id() err = " +
				strconv.Itoa(int(retCInt)))
			break
		} else {
			logger.Info("fpga_get_dev_id() ret = " +
				strconv.Itoa(int(retCInt)))
		}

		var childBsFile string
		if nil == functionConfigData {
			if nil != childBitstreamCRData.Spec.ChildBitstreamFile {
				childBsFile = *childBitstreamCRData.Spec.ChildBitstreamFile
			} else {
				childBsFile = ""
			}
		} else {
			childBsFile = functionConfigData.ChildBitstream.File
		}

		if "" != childBsFile {
			childBsFileCString := C.CString(childBsFile)
			defer C.free(unsafe.Pointer(childBsFileCString))
			retCInt = C.fpga_write_bitstream(*deviceID, 0, childBsFileCString)
			if 0 > retCInt {
				logger.Info("fpga_write_bitstream() err = " +
					strconv.Itoa(int(retCInt)))
				break
			} else {
				logger.Info("fpga_write_bitstream() ret = " +
					strconv.Itoa(int(retCInt)))
			}
		} else {
			logger.Info("fpga_write_bitstream() write skip. not saved ChildBitstreamFile")
		}
	}
	return retCInt
}

func (r *FPGAFunctionReconciler) SetFpgaInfo(ctx context.Context,
	deviceID C.uint,
	childBitstreamCRData examplecomv1.ChildBs) (C.int, bool) {
	logger := log.FromContext(ctx)

	var retCInt C.int
	breakFlag := false
	var frameSizeConfigCChar *C.char

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		retCInt = C.fpga_update_info(deviceID)
		if 0 > retCInt {
			logger.Info("fpga_update_info() err = " +
				strconv.Itoa(int(retCInt)))
			break
		} else {
			logger.Info("fpga_update_info() ret = " +
				strconv.Itoa(int(retCInt)))
		}

		retCInt = C.fpga_lldma_setup_buffer(deviceID)
		if 0 > retCInt {
			logger.Info("fpga_lldma_setup_buffer() err = " +
				strconv.Itoa(int(retCInt)))
			break
		} else {
			logger.Info("fpga_lldma_setup_buffer() ret = " +
				strconv.Itoa(int(retCInt)))
		}

		for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
			region := childBitstreamCRData.Spec.Regions[regionIndex]
			nameSlice := []rune(*region.Name)
			lane, _ := strconv.Atoi(string(nameSlice[len(nameSlice)-1:]))
			lldmaExtIfID := 0
			ptuExtIfID := 1

			retCInt = C.fpga_chain_set_ddr(deviceID, C.uint(lane), C.uint(lldmaExtIfID))
			if 0 > retCInt {
				logger.Info("fpga_chain_set_ddr() err = " +
					strconv.Itoa(int(retCInt)))
				breakFlag = true
				break
			} else {
				logger.Info("fpga_chain_set_ddr() ret = " +
					strconv.Itoa(int(retCInt)))
			}

			retCInt = C.fpga_chain_set_ddr(deviceID, C.uint(lane), C.uint(ptuExtIfID))
			if 0 > retCInt {
				logger.Info("fpga_chain_set_ddr() err = " +
					strconv.Itoa(int(retCInt)))
				breakFlag = true
				break
			} else {
				logger.Info("fpga_chain_set_ddr() ret = " +
					strconv.Itoa(int(retCInt)))
			}

			for functionIndex := 0; functionIndex < len(*region.Modules.Functions); functionIndex++ {

				var frameSizeBytes []byte
				var frameSizeStruct examplecomv1.FrameSizeData
				FunctionsData := (*region.Modules.Functions)[functionIndex]
				if nil == FunctionsData.Parameters {
					logger.Info("FrameSize not found in ChildBitstreamsCR.")
					breakFlag = true
					break
				}
				frameSizeStruct.InputWidth = (*FunctionsData.Parameters)["InputWidth"].IntVal
				frameSizeStruct.InputHeight = (*FunctionsData.Parameters)["InputHeight"].IntVal
				frameSizeStruct.OutputWidth = (*FunctionsData.Parameters)["OutputWidth"].IntVal
				frameSizeStruct.OutputHeight = (*FunctionsData.Parameters)["OutputHeight"].IntVal
				frameSizeBytes, err := json.Marshal(frameSizeStruct)
				if nil != err {
					logger.Error(err, "FrameSize unable to Marshal. lane="+strconv.Itoa(lane))
					breakFlag = true
					break
				}

				for funcmodIndex := 0; funcmodIndex < len(*FunctionsData.Module); funcmodIndex++ {
					funcNameCString := C.CString(*(*FunctionsData.Module)[funcmodIndex].Type)
					defer C.free(unsafe.Pointer(funcNameCString))

					retCInt = C.fpga_chain_start(deviceID, C.uint(*FunctionsData.ID))
					if 0 > retCInt {
						logger.Info("fpga_chain_start() err = " +
							strconv.Itoa(int(retCInt)))
						breakFlag = true
						break
					} else {
						logger.Info("fpga_chain_start() ret = " +
							strconv.Itoa(int(retCInt)))
					}
					retCInt = C.fpga_direct_start(deviceID, C.uint(*FunctionsData.ID))
					if 0 > retCInt {
						logger.Info("fpga_direct_start() err = " +
							strconv.Itoa(int(retCInt)))
						breakFlag = true
						break
					} else {
						logger.Info("fpga_direct_start() ret = " +
							strconv.Itoa(int(retCInt)))
					}

					retCInt = C.fpga_function_config(deviceID, C.uint(*FunctionsData.ID), funcNameCString)
					if 0 > retCInt {
						logger.Info("fpga_function_config() err = " +
							strconv.Itoa(int(retCInt)))
						breakFlag = true
						break
					} else {
						logger.Info("fpga_function_config() ret = " +
							strconv.Itoa(int(retCInt)))
					}

					retCInt = C.fpga_function_init(deviceID, C.uint(*FunctionsData.ID), nil)
					if 0 > retCInt {
						logger.Info("fpga_function_init() err = " +
							strconv.Itoa(int(retCInt)))
						breakFlag = true
						break
					} else {
						logger.Info("fpga_function_init() ret = " +
							strconv.Itoa(int(retCInt)))
					}

					if 0 != len(string(frameSizeBytes)) {
						frameSizeCString := C.CString(string(frameSizeBytes))
						defer C.free(unsafe.Pointer(frameSizeCString))

						retCInt = C.fpga_function_set(deviceID, C.uint(*FunctionsData.ID), frameSizeCString)
						if 0 > retCInt {
							logger.Info("fpga_function_set() err = " +
								strconv.Itoa(int(retCInt)))
							breakFlag = true
							break
						} else {
							logger.Info("fpga_function_set() ret = " +
								strconv.Itoa(int(retCInt)))
						}

						// Get frame size information
						var frameSizeConfig map[string]examplecomv1.FrameSizeData
						retCInt = C.fpga_function_get(deviceID, C.uint(*FunctionsData.ID),
							(**C.char)(unsafe.Pointer(&frameSizeConfigCChar))) //nolint:gocritic // suspicious identical LHS and RHS for `==` operator
						if 0 > retCInt {
							logger.Info("fpga_function_get() err = " +
								strconv.Itoa(int(retCInt)))
							breakFlag = true
							break
						} else {
							logger.Info("fpga_function_get() ret = " +
								strconv.Itoa(int(retCInt)))
						}

						n := 0
						head := (*byte)(unsafe.Pointer(frameSizeConfigCChar))
						for ptr := head; *ptr != 0; n++ {
							ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
						}

						frameSizeConfigByte := C.GoBytes(unsafe.Pointer(
							frameSizeConfigCChar), C.int(n))

						err = json.Unmarshal(frameSizeConfigByte,
							&frameSizeConfig)
						if nil != err {
							logger.Info("unable to unmarshal fpga_function_get() data.")
							breakFlag = true
							break
						}

						if frameSizeConfig["fr"].InputWidth != frameSizeStruct.InputWidth ||
							frameSizeConfig["fr"].InputHeight != frameSizeStruct.InputHeight ||
							frameSizeConfig["fr"].OutputWidth != frameSizeStruct.OutputWidth ||
							frameSizeConfig["fr"].OutputHeight != frameSizeStruct.OutputHeight {
							logger.Info("fpga_function_get() data is wrong." +
								"deviceID: " + strconv.Itoa(int(deviceID)) +
								", FunctionID: " + strconv.Itoa(int(*FunctionsData.ID)) +
								", InputWidth: " + strconv.Itoa(int(frameSizeConfig["fr"].InputWidth)) +
								", InputHeight: " + strconv.Itoa(int(frameSizeConfig["fr"].InputHeight)) +
								", OutputWidth: " + strconv.Itoa(int(frameSizeConfig["fr"].OutputWidth)) +
								", OutputHeight: " + strconv.Itoa(int(frameSizeConfig["fr"].OutputHeight)))
							logger.Info("fpga_function_set() data is set." +
								", InputWidth: " + strconv.Itoa(int(frameSizeStruct.InputWidth)) +
								", InputHeight: " + strconv.Itoa(int(frameSizeStruct.InputHeight)) +
								", OutputWidth: " + strconv.Itoa(int(frameSizeStruct.OutputWidth)) +
								", OutputHeight: " + strconv.Itoa(int(frameSizeStruct.OutputHeight)))
							breakFlag = true
							break
						}
					}
				}
				if true == breakFlag {
					break
				}
			}
			if true == breakFlag {
				break
			}
		}
	}
	return retCInt, breakFlag
}

func (r *FPGAFunctionReconciler) UpdCustomResource(ctx context.Context,
	pCRData *examplecomv1.FPGAFunction, status string) error {
	logger := log.FromContext(ctx)
	var err error

	if status == RUNNING {
		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, r.GetfinalizerName(pCRData))
		// status update
		pCRData.Status.StartTime = metav1.Now()
		pCRData.Status.Status = status

		if nil == pCRData.Spec.FunctionIndex {
			pCRData.Status.FunctionIndex = 0
		} else {
			pCRData.Status.FunctionIndex = *pCRData.Spec.FunctionIndex
		}
	}
	err = r.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "FPGAFunctionCR Status Update Error.")
	} else {
		logger.Info("FPGAFunctionCR Status Update.")
	}
	return err
}

func (r *FPGAFunctionReconciler) UpdCRFPGAReconfiguration(ctx context.Context,
	pCRData *examplecomv1.FPGAReconfiguration, status string) error {
	logger := log.FromContext(ctx)
	var err error

	if status == examplecomv1.FPGARECONFSTATUSSUCCEEDED {
		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, r.GetfinalizerNameFPGAReconfiguration(pCRData))
	}
	// status update
	pCRData.Status.Status = status

	err = r.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "FPGAReconfigurationCR Status Update Error.")
	} else {
		logger.Info("FPGAReconfigurationCR Status Update.")
	}
	return err
}

func (r *FPGAFunctionReconciler) DelCustomResource(ctx context.Context,
	pCRData *examplecomv1.FPGAFunction) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	// Delete the Finalizer statement.
	if controllerutil.ContainsFinalizer(pCRData,
		r.GetfinalizerName(pCRData)) {
		controllerutil.RemoveFinalizer(pCRData, r.GetfinalizerName(pCRData))
		err := r.Update(ctx, pCRData)
		if err != nil {
			logger.Error(err, "RemoveFinalizer Update Error.")
		}
	}
	return err
}

func (r *FPGAFunctionReconciler) DelCRFPGAReconfiguration(ctx context.Context,
	pCRData *examplecomv1.FPGAReconfiguration) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	// Delete the Finalizer statement.
	if controllerutil.ContainsFinalizer(pCRData,
		r.GetfinalizerNameFPGAReconfiguration(pCRData)) {
		controllerutil.RemoveFinalizer(pCRData, r.GetfinalizerNameFPGAReconfiguration(pCRData))
		err := r.Update(ctx, pCRData)
		if err != nil {
			logger.Error(err, "RemoveFinalizer Update Error.")
		}
	}
	return err
}

// Get Config information for FPGAFunc
func (r *FPGAFunctionReconciler) getConfigData(
	ctx context.Context,
	configName string,
	configData *[]byte) error {

	var err error
	var mapdata map[string]string

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
			mapdata, _, _ = unstructured.NestedStringMap(tmpData.Object, "data")
			for _, jsonrecord := range mapdata {
				*configData = []byte(jsonrecord)
			}
		}
	}
	return err
}

// Convert Config data for FPGAFunction
func FunctionConfigDataJsonUnmarshal(
	configData *[]byte,
	pConfigData *examplecomv1.FPGAFuncConfig) error {

	err := json.Unmarshal(*configData, &pConfigData)
	return err
}

// Get CR information for FPGA
func GetFPGACR(ctx context.Context,
	mngr ctrl.Manager,
	myNodeName string) (fpgaList []string) {
	logger := log.FromContext(ctx)
	r := mngr.GetAPIReader() // Get the manager to use the Get/List functions

	fpgaCRData := &examplecomv1.FPGAList{}
	// Get a ConfigMap by namespace/name
	err := r.List(context.Background(), fpgaCRData)
	if errors.IsNotFound(err) {
		logger.Error(err, "fpgaCR does not exist")
	}
	if err != nil {
		logger.Error(err, "unable to fetch fpgaCR")
	}

	for i := 0; i < len(fpgaCRData.Items); i++ {
		if myNodeName != fpgaCRData.Items[i].Status.NodeName {
			continue
		}
		fpgaList = append(fpgaList, fpgaCRData.Items[i].Status.DeviceFilePath)
	}
	return fpgaList
}

// FPGA controller startup process
func StartupProccessing(mng ctrl.Manager) error {

	ctx := context.Background()
	logger := log.FromContext(ctx)

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		myNodeName = os.Getenv("K8S_NODENAME")
		myClusterName = os.Getenv("K8S_CLUSTERNAME")

		if myNodeName == "" || myClusterName == "" {
			logger.Info("The node name or cluster name could not be obtained.")
			err = fmt.Errorf("startup error.")
			break
		}

		fpgaList := GetFPGACR(ctx, mng, myNodeName)

		if 0 != len(fpgaList) {
			// FPGAPhase3 initial setting
			var argv []*C.char
			// Phase 3 FPGA initialization variables
			argv0 := C.CString("proc")
			argv1 := C.CString("-d")
			argv2 := C.CString(strings.Join(fpgaList, ","))
			defer C.free(unsafe.Pointer(argv0))
			defer C.free(unsafe.Pointer(argv1))
			defer C.free(unsafe.Pointer(argv2))
			argv = []*C.char{argv0, argv1, argv2}
			argc := C.int(len(argv))
			C.libfpga_log_set_output_stdout()
			C.libfpga_log_set_level(C.LIBFPGA_LOG_ALL)
			C.fpga_init(argc, (**C.char)(unsafe.Pointer(&argv[0])))

			// Hold device information
			for deviceID, _ := range fpgaList {
				C.fpga_enable_regrw(C.uint(deviceID))
			}
		}
	}
	return err
}

// Config information storage area
var gServicerMgmtInfo []examplecomv1.ServicerMgmtInfo
var gDeployInfo map[string][]examplecomv1.DeviceRegionInfo
var gRegionUniqueInfo []examplecomv1.RegionSpecificInfo
var gFunctionUniqueInfo []examplecomv1.FPGACatalog
var gPreDeterminedRegionInfo []examplecomv1.PreDeterminedRegionInfo
var gFilterResizeInfo []examplecomv1.FunctionDetail

type ConfigTable struct {
	name string
}

var configLoadTableForWriteChildBs = []ConfigTable{
	{examplecomv1.CMServicerMgmtInfo},
	{examplecomv1.CMDeployInfo},
	{examplecomv1.CMRegionUniqueInfo},
	{examplecomv1.CMFunctionUniqueInfo},
	{examplecomv1.CMPreDeterminedRegionInfo},
	{examplecomv1.CMFilterResizeInfo},
}

// Get Config information for FPGA Write ChildBitstream
func (r *FPGAFunctionReconciler) getConfigMapForWriteChildBs(
	ctx context.Context) error {
	logger := log.FromContext(ctx)

	var cfgdata []byte
	var err error

	for _, record := range configLoadTableForWriteChildBs {
		err = r.getConfigData(ctx, record.name, &cfgdata)
		if nil != err {
			continue
		}
		if examplecomv1.CMServicerMgmtInfo == record.name {
			err = json.Unmarshal(cfgdata, &gServicerMgmtInfo)
		} else if examplecomv1.CMDeployInfo == record.name {
			err = json.Unmarshal(cfgdata, &gDeployInfo)
		} else if examplecomv1.CMRegionUniqueInfo == record.name {
			err = json.Unmarshal(cfgdata, &gRegionUniqueInfo)
		} else if examplecomv1.CMFunctionUniqueInfo == record.name {
			err = json.Unmarshal(cfgdata, &gFunctionUniqueInfo)
		} else if examplecomv1.CMPreDeterminedRegionInfo == record.name {
			err = json.Unmarshal(cfgdata, &gPreDeterminedRegionInfo)
		} else if examplecomv1.CMFilterResizeInfo == record.name {
			var filterResizeMap map[string][]examplecomv1.FunctionDetail
			err = json.Unmarshal(cfgdata, &filterResizeMap)
			gFilterResizeInfo = filterResizeMap["functionKernels"]
		}

		if nil != err {
			logger.Error(err,
				"unable to Unmarshal. ConfigMap="+record.name)
			break
		}
	}
	return err

}

// get ChildBs Config
func (r *FPGAFunctionReconciler) getChildBsConfig(ctx context.Context,
	parentBitstreamID string,
	childBitstreamID string,
	childBitstreamCRData *examplecomv1.ChildBs) (C.int, error) {
	logger := log.FromContext(ctx)

	var err error
	var retCInt C.int
	var bitstreamIDConfigCChar *C.char
	var bitstreamIDConfig examplecomv1.BsConfigInfo

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		parentBitstreamCString := C.CString(parentBitstreamID)
		defer C.free(unsafe.Pointer(parentBitstreamCString))
		childBitstreamCString := C.CString(childBitstreamID)
		defer C.free(unsafe.Pointer(childBitstreamCString))
		retCInt = C.fpga_db_get_device_config_by_bitstream_id(parentBitstreamCString,
			childBitstreamCString,
			(**C.char)(unsafe.Pointer(&bitstreamIDConfigCChar))) //nolint:gocritic // suspicious identical LHS and RHS for `==` operator
		if 0 > retCInt {
			logger.Info("fpga_db_get_device_config_by_bitstream_id()" +
				" err = " + strconv.Itoa(int(retCInt)))
			break
		} else {
			logger.Info("fpga_db_get_device_config_by_bitstream_id()" +
				" ret = " + strconv.Itoa(int(retCInt)))
		}

		if 0 == retCInt {

			// Get the number of characters in the JSON data
			n := 0
			head := (*byte)(unsafe.Pointer(bitstreamIDConfigCChar))
			for ptr := head; *ptr != 0; n++ {
				ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
			}

			// Convert to []bytes type
			bitstreamIDConfigBytes := C.GoBytes(unsafe.Pointer(bitstreamIDConfigCChar), C.int(n))

			// Get FPGA function information
			err = json.Unmarshal(bitstreamIDConfigBytes, &bitstreamIDConfig)
			if nil != err {
				logger.Info("unable to unmarshal fpga_db_get_device_config_by_bitstream_id() data.")
				break
			}

			childBitstreamCRData.Spec = bitstreamIDConfig.ChildBitstreamIDs[0]
		}
	}
	return retCInt, err
}

// Create ChildBsCR
func (r *FPGAFunctionReconciler) createChildBsCR(ctx context.Context,
	functionConfigDataMap map[int32]examplecomv1.FPGAFuncConfig,
	manualFlag bool,
	childBitstreamCRData *examplecomv1.ChildBs,
	fpgaCRData *examplecomv1.FPGA,
	childBitstreamCRBase *examplecomv1.ChildBs) error {
	logger := log.FromContext(ctx)

	var err error
	var bitstreamIDConfigCChar *C.char

	defer C.free(unsafe.Pointer(bitstreamIDConfigCChar))

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		var ptuData *examplecomv1.ChildBsPtu
		var functionData examplecomv1.ChildBsFunctions

		var networkInfo examplecomv1.NetworkData
		var setAvailable bool = true

		for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
			region := childBitstreamCRData.Spec.Regions[regionIndex]

			nameSlice := []rune(*region.Name)
			lane, _ := strconv.Atoi(string(nameSlice[len(nameSlice)-1:]))
			functionConfigData, functionConfigDataOk := functionConfigDataMap[int32(lane)]
			if false == functionConfigDataOk {
				if false == manualFlag {
					if 0 != lane {
						lane = 0
						functionConfigData, functionConfigDataOk = functionConfigDataMap[int32(lane)]
					}
				}
			}
			var inputHeight int32
			var inputWidth int32
			var outputHeight int32
			var outputWidth int32
			if true == functionConfigDataOk {
				inputHeight = functionConfigData.Parameters.Functions.InputHeight
				inputWidth = functionConfigData.Parameters.Functions.InputWidth
				outputHeight = functionConfigData.Parameters.Functions.OutputHeight
				outputWidth = functionConfigData.Parameters.Functions.OutputWidth
				if nil == childBitstreamCRData.Spec.ChildBitstreamFile {
					childBitstreamFile := functionConfigData.ChildBitstream.File
					childBitstreamCRData.Spec.ChildBitstreamFile = &childBitstreamFile
				}
			}

			for srvIndex := 0; srvIndex < len(gServicerMgmtInfo); srvIndex++ {
				if nil == region.Modules.Ptu {
					break
				}
				if gServicerMgmtInfo[srvIndex].NodeName != myNodeName {
					continue
				}
				networkInfos := gServicerMgmtInfo[srvIndex].NetworkInfo
				for networkIndex := 0; networkIndex < len(networkInfos); networkIndex++ {
					networkInfo = networkInfos[networkIndex]
					if networkInfo.DeviceFilePath != fpgaCRData.Spec.DeviceFilePath {
						continue
					}
					nameSlice := []rune(*region.Name)
					lane, _ := strconv.Atoi(string(nameSlice[len(nameSlice)-1:]))
					if networkInfo.LaneIndex == int32(lane) {
						var paramData map[string]intstr.IntOrString = make(map[string]intstr.IntOrString)

						paramData["GatewayAddress"] = intstr.IntOrString{Type: intstr.String,
							StrVal: networkInfo.GatewayAddress}
						paramData["IPAddress"] = intstr.IntOrString{Type: intstr.String,
							StrVal: networkInfo.IPAddress}
						paramData["MACAddress"] = intstr.IntOrString{Type: intstr.String,
							StrVal: networkInfo.MACAddress}
						paramData["SubnetAddress"] = intstr.IntOrString{Type: intstr.String,
							StrVal: networkInfo.SubnetAddress}
						ptuData = region.Modules.Ptu
						ptuData.Parameters = &paramData

						childBitstreamCRData.Spec.Regions[regionIndex].Modules.Ptu = ptuData
					}
				}
			}

			functions := len(*region.Modules.Functions)
			for functionIndex := 0; functionIndex < functions; functionIndex++ {

				functionData = (*region.Modules.Functions)[functionIndex]

				for fuiIndex := 0; fuiIndex < len(gFunctionUniqueInfo); fuiIndex++ {
					for fumIndex := 0; fumIndex < len(*functionData.Module); fumIndex++ {

						if gFunctionUniqueInfo[fuiIndex].FunctionName ==
							*(*functionData.Module)[fumIndex].Type {
							var maxCapacitydata int32 = gFunctionUniqueInfo[fuiIndex].MaxCapacity
							var maxDataFlowsdata int32 = gFunctionUniqueInfo[fuiIndex].MaxDataFlows
							functionData.DeploySpec.MaxCapacity = &maxCapacitydata
							functionData.DeploySpec.MaxDataFlows = &maxDataFlowsdata
						}
					}
				}

				for fumIndex := 0; fumIndex < len(*functionData.Module); fumIndex++ {
					functionChannelIDs := *(*functionData.Module)[fumIndex].FunctionChannelIDs
					idRange := strings.SplitN(functionChannelIDs, "-", 2)
					idRangeMin, _ := strconv.Atoi(idRange[0])
					idRangeMax, _ := strconv.Atoi(idRange[1])

					for idIndex := idRangeMin; idIndex <= idRangeMax; idIndex++ {
						if nil == functionData.IntraResourceMgmtMap {
							intraResourceMgmtMap := make(map[string]examplecomv1.FunctionsIntraResourceMgmtMap)
							functionData.IntraResourceMgmtMap = &intraResourceMgmtMap
						}
						available := setAvailable
						(*functionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)] =
							examplecomv1.FunctionsIntraResourceMgmtMap{Available: &available}
					}
				}

				if true == functionConfigDataOk {
					if nil == functionData.Parameters {
						parameters := make(map[string]intstr.IntOrString)
						functionData.Parameters = &parameters
					}
					(*functionData.Parameters)["InputHeight"] =
						intstr.IntOrString{Type: intstr.Int, IntVal: inputHeight}
					(*functionData.Parameters)["InputWidth"] =
						intstr.IntOrString{Type: intstr.Int, IntVal: inputWidth}
					(*functionData.Parameters)["OutputHeight"] =
						intstr.IntOrString{Type: intstr.Int, IntVal: outputHeight}
					(*functionData.Parameters)["OutputWidth"] =
						intstr.IntOrString{Type: intstr.Int, IntVal: outputWidth}
					functionName := functionConfigData.FunctionName
					functionData.FunctionName = &functionName

				}
				(*childBitstreamCRData.Spec.Regions[regionIndex].Modules.Functions)[functionIndex] = functionData
			}

			fmt.Println("debug region-unique-info: ", gRegionUniqueInfo)
			for ruiIndex := 0; ruiIndex < len(gRegionUniqueInfo); ruiIndex++ {
				regionInfo := gRegionUniqueInfo[ruiIndex]
				if regionInfo.SubDeviceSpecRef == *childBitstreamCRData.Spec.ChildBitstreamID {
					for targetIndex := 0; targetIndex < len(regionInfo.FunctionTargets); targetIndex++ {
						target := regionInfo.FunctionTargets[targetIndex]
						childBitstreamCRData.Spec.Regions[regionIndex].MaxCapacity =
							&target.MaxCapacity
						childBitstreamCRData.Spec.Regions[regionIndex].MaxFunctions =
							&target.MaxFunctions
					}
				}
			}
		}

		childBitstreamCRData.Status.Regions = childBitstreamCRData.Spec.Regions
		childBitstreamCRData.Status.State = examplecomv1.ChildBsWritingBsfile
		childBitstreamCRData.Status.Status = examplecomv1.ChildBsStatusPreparing

		if nil != childBitstreamCRData.Spec.ChildBitstreamID {
			logger.Info("debug fpgaCRData.ObjectMeta.Name :" + fpgaCRData.ObjectMeta.Name)
			logger.Info("debug childBitstreamCRData.Spec.ChildBitstreamID :" + *childBitstreamCRData.Spec.ChildBitstreamID)
		} else {
			logger.Info("debug childBitstreamCRData.Spec.ChildBitstreamID : nil ")
		}
		reqData := &examplecomv1.ChildBs{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fpgaCRData.ObjectMeta.Name + "-" + *childBitstreamCRData.Spec.ChildBitstreamID,
				Namespace: "default",
			},
			Spec:   childBitstreamCRData.Spec,
			Status: childBitstreamCRData.Status,
		}

		err = ctrl.SetControllerReference(fpgaCRData, reqData, r.Scheme)
		if nil != err {
			logger.Error(err, "ctrl.SetControllerReference Error.")
		} else {
			err = r.Create(ctx, reqData)
			if err != nil {
				logger.Error(err, "Failed to create ChildBitstreamCR.")
				break
			} else {
				logger.Info("Success to create ChildBitstreamCR.")
				if nil == childBitstreamCRBase {
					childBitstreamCRBase = childBitstreamCRData
					logger.Info("debug childBitstreamCRBase.Name : " + childBitstreamCRBase.Name)
				}
				break
			}
		}
	}

	return err
}

// FPGA Data Get
func (r *FPGAFunctionReconciler) getFPGACRData(ctx context.Context,
	crName string,
	crData *examplecomv1.FPGA) error {
	logger := log.FromContext(ctx)

	var err error

	err = r.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      crName}, crData)
	if errors.IsNotFound(err) {
		// If CR does not exist
		logger.Info("NotFound to fetch FPGACR")
	} else if err != nil {
		logger.Error(err, "unable to fetch FPGACR")
	}
	return err
}

// ChildBs Data Get
func (r *FPGAFunctionReconciler) getChildBsData(
	ctx context.Context,
	fpgaCRName string,
	childBsID string,
	crData *examplecomv1.ChildBs,
	childBitstreamCRBase *examplecomv1.ChildBs) error {
	logger := log.FromContext(ctx)

	var err error

	logger.Info("debug : " + fpgaCRName + "-" + childBsID)

	err = r.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      fpgaCRName + "-" + childBsID}, crData)
	if errors.IsNotFound(err) {
		// If CR does not exist
		logger.Info("NotFound to fetch ChiledBitstreaCR")
	} else if err != nil {
		logger.Error(err, "unable to fetch ChiledBitstreamCR")
	} else if err == nil {
		if nil == childBitstreamCRBase {
			childBitstreamCRBase = crData
			logger.Info("debug childBitstreamCRBase.Name : " + childBitstreamCRBase.Name)
		}
	}
	return err
}

// ComputeResource Data Get
func (r *FPGAFunctionReconciler) getComputeResourceCR(
	ctx context.Context,
	pComputeResourceData *examplecomv1.ComputeResource) error {
	logger := log.FromContext(ctx)

	var err error

	logger.Info("debug : " + COMPUTERESOURCENAME + myNodeName)

	err = r.Get(ctx, client.ObjectKey{
		Namespace: myClusterName,
		Name:      COMPUTERESOURCENAME + myNodeName}, pComputeResourceData)
	if errors.IsNotFound(err) {
		// If ComputeResource CR does not exist
		logger.Info("ComputeResource does not exist.")
	} else if nil != err {
		// If an error occurs in the Get function
		logger.Error(err, "unable to get ComputeResource.")
	}
	return err
}

// FPGA CR Update
func (r *FPGAFunctionReconciler) updFPGACR(
	ctx context.Context,
	pCRData *examplecomv1.FPGA,
	childBsID *string) error {
	logger := log.FromContext(ctx)
	var err error

	pCRData.Spec.ChildBitstreamID = childBsID
	pCRData.Status.ChildBitstreamID = childBsID
	if nil != childBsID {
		crName := pCRData.ObjectMeta.Name + "-" + *childBsID
		pCRData.Status.ChildBitstreamCRName = &crName
	} else {
		pCRData.Status.ChildBitstreamCRName = nil
	}

	err = r.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "FPGACR Status Update Error.")
	} else {
		logger.Info("FPGACR Status Update.")
	}
	return err
}

// ChildBs CR Update
func (r *FPGAFunctionReconciler) updChildBsCR(ctx context.Context,
	pCRData *examplecomv1.ChildBs,
	state examplecomv1.ChildBitstreamState,
	status examplecomv1.ChildBitstreamStatus,
	childBitstreamCRBase *examplecomv1.ChildBs) error {
	logger := log.FromContext(ctx)
	var err error

	pachData := childBitstreamCRBase.DeepCopy()
	pachData.Name = pCRData.Name
	pachData.Namespace = pCRData.Namespace
	pachData.Spec.Regions = pCRData.Spec.Regions
	pachData.Status.Regions = pCRData.Spec.Regions
	pachData.Status.Status = status
	pachData.Status.State = state

	patch := client.MergeFrom(childBitstreamCRBase)

	err = r.Patch(ctx, pachData, patch)
	if err != nil {
		logger.Error(err, "ChildBitstreamCR Status Update Error.")
	} else {
		logger.Info("ChildBitstreamCR Status Update.")
		pCRData.Status.Status = status
		pCRData.Status.State = state
	}
	return err
}

// deployinfo CM Update
func (r *FPGAFunctionReconciler) updDeployinfoCM(ctx context.Context,
	deviceUUID string,
	nodeName string,
	functionName *string,
	resetFlag bool,
	childBsRegions *[]examplecomv1.ChildBsRegion) error {
	logger := log.FromContext(ctx)

	var err error
	deployDevices := gDeployInfo["devices"]
	upData := make(map[string][]examplecomv1.DeviceRegionInfo)

	for deployIndex := 0; deployIndex < len(deployDevices); deployIndex++ {
		deployDevice := deployDevices[deployIndex]
		var deployDeviceData examplecomv1.DeviceRegionInfo

		if *deployDevice.DeviceUUID == deviceUUID &&
			deployDevice.NodeName == nodeName {
			functionTargets := deployDevice.FunctionTargets

			for ftIndex := 0; ftIndex < len(functionTargets); ftIndex++ {
				functionTarget := functionTargets[ftIndex]
				var functionTargetData examplecomv1.RegionInDevice
				var childBsFunctions []examplecomv1.ChildBsFunctions
				var maxCapacity int32
				var maxFunctions int32
				// var functions *[]examplecomv1.ChildBsFunctions
				for bsIndex := 0; bsIndex < len(*childBsRegions); bsIndex++ {
					if *(*childBsRegions)[bsIndex].Name == functionTarget.RegionName {
						childBsFunctions =
							*(*childBsRegions)[bsIndex].Modules.Functions
						// childBsFunctions = *functions
						if nil != (*childBsRegions)[bsIndex].MaxCapacity {
							maxCapacity = *(*childBsRegions)[bsIndex].MaxCapacity
						}
						if nil != (*childBsRegions)[bsIndex].MaxFunctions {
							maxFunctions = *(*childBsRegions)[bsIndex].MaxFunctions
						}
					}
				}

				if false == resetFlag {
					for funcIndex := 0; funcIndex < len(childBsFunctions); funcIndex++ {
						var functionData examplecomv1.SimpleFunctionInfraStruct

						function := (childBsFunctions)[funcIndex]

						partitionName := strconv.Itoa(int(*function.ID))
						funcIndexInt32 := int32(funcIndex)

						functionData.PartitionName = partitionName
						functionData.FunctionIndex = &funcIndexInt32
						if nil == functionName {
							if nil == function.FunctionName {
								functionData.FunctionName = ""
							} else {
								functionData.FunctionName = *function.FunctionName
							}
						} else {
							functionData.FunctionName = *functionName
						}
						functionData.MaxCapacity = *function.DeploySpec.MaxCapacity
						functionData.MaxDataFlows = *function.DeploySpec.MaxDataFlows

						functionTargetData.Functions = append(functionTargetData.Functions, functionData)
					}
				}

				functionTargetData.RegionType = functionTarget.RegionType
				functionTargetData.RegionName = functionTarget.RegionName
				functionTargetData.MaxCapacity = maxCapacity
				functionTargetData.MaxFunctions = maxFunctions

				deployDeviceData.FunctionTargets = append(deployDeviceData.FunctionTargets, functionTargetData)
			}
			deployDeviceData.NodeName = deployDevice.NodeName
			deployDeviceData.DeviceFilePath = deployDevice.DeviceFilePath
			deployDeviceData.DeviceUUID = deployDevice.DeviceUUID
			upData["devices"] = append(upData["devices"], deployDeviceData)
		} else {
			upData["devices"] = append(upData["devices"], deployDevice)
		}
	}
	createCM := &corev1.ConfigMap{}
	createCM.SetName(examplecomv1.CMDeployInfo)
	createCM.SetNamespace(examplecomv1.CMNameSpace)

	jsonData, err := json.Marshal(upData)
	if nil != err {
		logger.Error(err, "unable to marshal. ConfigMap="+
			examplecomv1.CMDeployInfo)
	} else {
		cmData := map[string]string{
			examplecomv1.CMDeployInfo + ".json": string(jsonData),
		}

		createCM.Data = cmData

		err = r.Update(ctx, createCM)
		if err != nil {
			logger.Error(err, examplecomv1.CMDeployInfo+" Update Error.")
		} else {
			logger.Info(examplecomv1.CMDeployInfo + " Update.")
		}
	}
	return err
}

// delete a ChildBs CR
func (r *FPGAFunctionReconciler) deleteChildBsCR(ctx context.Context,
	pCRData *examplecomv1.ChildBs) error {

	logger := log.FromContext(ctx)

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		err := r.Delete(ctx, pCRData)
		if err != nil {
			logger.Error(err, "Failed to delete ChildBitstreamCR.")
			break
		} else {
			logger.Info("Success to delete ChildBitstreamCR.")
			break
		}
	}
	return err
}

func (r *FPGAFunctionReconciler) getRxTxProtocol(ctx context.Context,
	req ctrl.Request,
	crData *examplecomv1.FPGAFunction,
	rxProtocol *string,
	txProtocol *string) error {
	logger := log.FromContext(ctx)

	var err error
	var previousListKey []string
	var nextListKey []string
	var previousConnectionCRName string
	var nextConnectionCRName string
	var fromConnectionKind string
	var toConnectionKind string

	dma := "DMA"
	tcp := "TCP"
	init := ""
	*rxProtocol = init
	*txProtocol = init

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		for key, _ := range crData.Spec.PreviousFunctions {
			previousListKey = append(previousListKey, key)
		}
		if 0 != len(previousListKey) {
			// From side ConnectionCR name generation
			MakeConnectionCRName(ctx,
				crData.Spec.DataFlowRef.Name,
				crData.Spec.PreviousFunctions[previousListKey[0]].WBFunctionRef.Name,
				req.Name,
				&previousConnectionCRName)
			logger.Info("previousConnectionCRName : " + previousConnectionCRName)

			// Get Kind of ConnectionCR on the From side
			err = r.GetConnectionData(ctx,
				req.Namespace,
				previousConnectionCRName,
				&fromConnectionKind)
			if nil != err {
				if errors.IsNotFound(err) {
					// do nothing
				} else {
					logger.Error(err, "GetConnectionData() error")
				}
				break
			}
		} else {
			fromConnectionKind = "nothing"
		}

		for key, _ := range crData.Spec.NextFunctions {
			nextListKey = append(nextListKey, key)
		}
		if 0 != len(previousListKey) {
			// Generate the To side ConnectionCR name
			MakeConnectionCRName(ctx,
				crData.Spec.DataFlowRef.Name,
				req.Name,
				crData.Spec.NextFunctions[nextListKey[0]].WBFunctionRef.Name,
				&nextConnectionCRName)
			logger.Info("nextConnectionCRName : " + nextConnectionCRName)

			// Get Kind of To side ConnectionCR
			err = r.GetConnectionData(ctx,
				req.Namespace,
				nextConnectionCRName,
				&toConnectionKind)
			if nil != err {
				if errors.IsNotFound(err) {
					// do nothing
				} else {
					logger.Error(err, "GetConnectionData() error")
				}
				break
			}
		} else {
			fromConnectionKind = "nothing"
		}

		switch fromConnectionKind {
		case examplecomv1.ConnectionCRKindPCIe:
			*rxProtocol = dma
		case examplecomv1.ConnectionCRKindEth:
			*rxProtocol = tcp
		default:
			*rxProtocol = init
		}

		switch toConnectionKind {
		case examplecomv1.ConnectionCRKindPCIe:
			*txProtocol = dma
		case examplecomv1.ConnectionCRKindEth:
			*txProtocol = tcp
		default:
			*txProtocol = init
		}

		logger.Info("rxProtocol : " + *rxProtocol)
		logger.Info("txProtocol : " + *txProtocol)
	}
	return err
}

// Allocate FPGA Resource
func (r *FPGAFunctionReconciler) allocateFPGAResource(ctx context.Context,
	req ctrl.Request,
	crData *examplecomv1.FPGAFunction,
	childBitstreamCRData *examplecomv1.ChildBs,
	fpgafuncRxData *examplecomv1.RxTxData,
	fpgafuncTxData *examplecomv1.RxTxData,
	functionKernelID *int32,
	chainID *int32,
	functionChannelID *int32,
	childBitstreamCRBase *examplecomv1.ChildBs) error {
	logger := log.FromContext(ctx)

	var err error
	var rxProtocol string
	var txProtocol string
	var upChildBitstreamCRData examplecomv1.ChildBs

	err = r.getRxTxProtocol(ctx,
		req,
		crData,
		&rxProtocol,
		&txProtocol)
	if nil != err {
		return err
	}

	for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
		specRegion := childBitstreamCRData.Spec.Regions[regionIndex]
		statusRegion := childBitstreamCRData.Status.Regions[regionIndex]

		if *specRegion.Name == crData.Spec.RegionName {

			var upSpecFunctionsData []examplecomv1.ChildBsFunctions
			var upStatusFunctionsData []examplecomv1.ChildBsFunctions
			for functionIndex := 0; functionIndex < len(*specRegion.Modules.Functions); functionIndex++ {

				specFunctionData := (*specRegion.Modules.Functions)[functionIndex]
				statusFunctionData := (*specRegion.Modules.Functions)[functionIndex]

				for funcmoduleIndex := 0; funcmoduleIndex < len(*specFunctionData.Module); funcmoduleIndex++ {
					functionChannelIDs := (*specFunctionData.Module)[funcmoduleIndex].FunctionChannelIDs

					if nil == functionChannelIDs {
						continue
					}

					var upIntraResourceMgmtMapData map[string]examplecomv1.FunctionsIntraResourceMgmtMap
					upIntraResourceMgmtMapData = make(map[string]examplecomv1.FunctionsIntraResourceMgmtMap)

					idRange := strings.SplitN(*functionChannelIDs, "-", 2)
					idRangeMin, _ := strconv.Atoi(idRange[0])
					idRangeMax, _ := strconv.Atoi(idRange[1])

					var intraResourceMgmtMapData examplecomv1.FunctionsIntraResourceMgmtMap
					var allocateFunctionChannelID int32 = -1
					var matchFunctionChannelID int32 = -1
					var functionChannelIDDetail examplecomv1.FPGAFunctionChannelIDs
					for idIndex := idRangeMin; idIndex <= idRangeMax; idIndex++ {
						if nil == (*specFunctionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)].Available {
							logger.Info("ChildBsCR IntraResourceMgmtMap Range Index unmatch." +
								" not found indexkey=" + strconv.Itoa(idIndex))
							continue
						}
						if -1 == allocateFunctionChannelID &&
							true == *(*specFunctionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)].Available {

							intraResourceMgmtMapData =
								(*specFunctionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)]
							allocateFunctionChannelID = int32(idIndex)
						} else {
							idIndexString := strconv.Itoa(idIndex)
							upIntraResourceMgmtMapData[idIndexString] =
								(*specFunctionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)]
						}
					}

					if -1 != allocateFunctionChannelID {
						for frIndex := 0; frIndex < len(gFilterResizeInfo); frIndex++ {
							functionDetailData := gFilterResizeInfo[frIndex]
							if strconv.Itoa(int(*specFunctionData.ID)) == functionDetailData.PartitionName {
								for i := range functionDetailData.FunctionChannelIDs {
									functionChannelIDsData := functionDetailData.FunctionChannelIDs[i]
									if functionChannelIDsData.FunctionChannelID == allocateFunctionChannelID {
										functionChannelIDDetail = functionChannelIDsData
										matchFunctionChannelID = allocateFunctionChannelID
										break
									}
								}
							}
						}
					}
					if -1 != matchFunctionChannelID {

						available := false

						var txData examplecomv1.RxTxSpec
						var rxData examplecomv1.RxTxSpec
						var txProtocolDatails map[string]examplecomv1.ChildBsDetails = make(map[string]examplecomv1.ChildBsDetails)
						var rxProtocolDatails map[string]examplecomv1.ChildBsDetails = make(map[string]examplecomv1.ChildBsDetails)
						var txDMAChannelID *int32
						var txLLDMAConnectorID *int32
						var txPort *int32
						var rxDMAChannelID *int32
						var rxLLDMAConnectorID *int32
						var rxPort *int32

						if nil != (*functionChannelIDDetail.Tx.Protocol)[txProtocol].DMAChannelID {
							txDMAChannelID = (*functionChannelIDDetail.Tx.Protocol)[txProtocol].DMAChannelID
						}
						if nil != (*functionChannelIDDetail.Tx.Protocol)[txProtocol].LLDMAConnectorID {
							txLLDMAConnectorID = (*functionChannelIDDetail.Tx.Protocol)[txProtocol].LLDMAConnectorID
						}
						if nil != (*functionChannelIDDetail.Tx.Protocol)[txProtocol].Port {
							txPort = (*functionChannelIDDetail.Tx.Protocol)[txProtocol].Port
						}
						if nil != (*functionChannelIDDetail.Rx.Protocol)[rxProtocol].DMAChannelID {
							rxDMAChannelID = (*functionChannelIDDetail.Rx.Protocol)[rxProtocol].DMAChannelID
						}
						if nil != (*functionChannelIDDetail.Rx.Protocol)[rxProtocol].LLDMAConnectorID {
							rxLLDMAConnectorID = (*functionChannelIDDetail.Rx.Protocol)[rxProtocol].LLDMAConnectorID
						}
						if nil != (*functionChannelIDDetail.Rx.Protocol)[rxProtocol].Port {
							rxPort = (*functionChannelIDDetail.Rx.Protocol)[rxProtocol].Port
						}

						txProtocolDatails[txProtocol] = examplecomv1.ChildBsDetails{
							DMAChannelID:     txDMAChannelID,
							LLDMAConnectorID: txLLDMAConnectorID,
							Port:             txPort}
						rxProtocolDatails[rxProtocol] = examplecomv1.ChildBsDetails{
							DMAChannelID:     rxDMAChannelID,
							LLDMAConnectorID: rxLLDMAConnectorID,
							Port:             rxPort}

						txData.Protocol = &txProtocolDatails
						rxData.Protocol = &rxProtocolDatails

						intraResourceMgmtMapData.Available = &available
						intraResourceMgmtMapData.FunctionCRName = &req.Name
						intraResourceMgmtMapData.Tx = &txData
						intraResourceMgmtMapData.Rx = &rxData

						*functionKernelID = *specFunctionData.ID
						*chainID = *specRegion.Modules.Chain.ID
						*functionChannelID = allocateFunctionChannelID

						if nil != specRegion.Modules.Ptu {
							if nil != specRegion.Modules.Ptu.Parameters {
								if "" != (*specRegion.Modules.Ptu.Parameters)["IPAddress"].StrVal {
									fpgafuncIPAddress := (*specRegion.Modules.Ptu.Parameters)["IPAddress"].StrVal
									fpgafuncRxData.IPAddress = &fpgafuncIPAddress
									fpgafuncTxData.IPAddress = &fpgafuncIPAddress
								}
								if "" != (*specRegion.Modules.Ptu.Parameters)["SubnetAddress"].StrVal {
									fpgafuncSubnetAddress := (*specRegion.Modules.Ptu.Parameters)["SubnetAddress"].StrVal
									fpgafuncRxData.SubnetAddress = &fpgafuncSubnetAddress
									fpgafuncTxData.SubnetAddress = &fpgafuncSubnetAddress
								}
								if "" != (*specRegion.Modules.Ptu.Parameters)["GatewayAddress"].StrVal {
									fpgafuncGatewayAddress := (*specRegion.Modules.Ptu.Parameters)["GatewayAddress"].StrVal
									fpgafuncRxData.GatewayAddress = &fpgafuncGatewayAddress
									fpgafuncTxData.GatewayAddress = &fpgafuncGatewayAddress
								}
							}
						}

						fpgafuncRxData.Protocol = rxProtocol
						fpgafuncRxData.DMAChannelID = rxDMAChannelID
						fpgafuncRxData.LLDMAConnectorID = rxLLDMAConnectorID
						fpgafuncRxData.Port = rxPort
						fpgafuncTxData.Protocol = txProtocol
						fpgafuncTxData.DMAChannelID = txDMAChannelID
						fpgafuncTxData.LLDMAConnectorID = txLLDMAConnectorID
						fpgafuncTxData.Port = txPort

						functionChannelIDString := strconv.Itoa(int(allocateFunctionChannelID))
						upIntraResourceMgmtMapData[functionChannelIDString] = intraResourceMgmtMapData
					}

					*specFunctionData.IntraResourceMgmtMap = upIntraResourceMgmtMapData
					*statusFunctionData.IntraResourceMgmtMap = upIntraResourceMgmtMapData
					upSpecFunctionsData = append(upSpecFunctionsData, specFunctionData)
					upStatusFunctionsData = append(upStatusFunctionsData, statusFunctionData)
				}
			}
			specRegion.Modules.Functions = &upSpecFunctionsData
			statusRegion.Modules.Functions = &upStatusFunctionsData
			upChildBitstreamCRData.Spec.Regions =
				append(upChildBitstreamCRData.Spec.Regions, specRegion)
			upChildBitstreamCRData.Status.Regions =
				append(upChildBitstreamCRData.Status.Regions, statusRegion)
		} else {
			upChildBitstreamCRData.Spec.Regions =
				append(upChildBitstreamCRData.Spec.Regions, specRegion)
			upChildBitstreamCRData.Status.Regions =
				append(upChildBitstreamCRData.Status.Regions, statusRegion)
		}
	}

	err = r.updChildBsCR(ctx,
		childBitstreamCRData,
		examplecomv1.ChildBsReady,
		examplecomv1.ChildBsStatusReady,
		childBitstreamCRBase)
	if err != nil {
		logger.Error(err, "ChildBsCR update failed.")
	}
	return err
}

func (r *FPGAFunctionReconciler) freeFPGAResource(ctx context.Context,
	req ctrl.Request,
	crData *examplecomv1.FPGAFunction,
	childBitstreamCRData *examplecomv1.ChildBs,
	childBitstreamCRBase *examplecomv1.ChildBs) error {
	logger := log.FromContext(ctx)

	var err error

	// clear the region match and function modual info FunctionsIntraResourceMgmtMap FunctionCRName match
	for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
		specRegion := childBitstreamCRData.Spec.Regions[regionIndex]
		statusRegion := childBitstreamCRData.Status.Regions[regionIndex]

		if *specRegion.Name == crData.Spec.RegionName {

			var upSpecFunctionsData []examplecomv1.ChildBsFunctions
			var upStatusFunctionsData []examplecomv1.ChildBsFunctions
			for functionIndex := 0; functionIndex < len(*specRegion.Modules.Functions); functionIndex++ {

				specFunctionData := (*specRegion.Modules.Functions)[functionIndex]
				statusFunctionData := (*specRegion.Modules.Functions)[functionIndex]

				for funcmoduleIndex := 0; funcmoduleIndex < len(*specFunctionData.Module); funcmoduleIndex++ {
					functionChannelIDs := (*specFunctionData.Module)[funcmoduleIndex].FunctionChannelIDs

					if nil == functionChannelIDs {
						continue
					}

					var upIntraResourceMgmtMapData map[string]examplecomv1.FunctionsIntraResourceMgmtMap
					upIntraResourceMgmtMapData = make(map[string]examplecomv1.FunctionsIntraResourceMgmtMap)

					idRange := strings.SplitN(*functionChannelIDs, "-", 2)
					idRangeMin, _ := strconv.Atoi(idRange[0])
					idRangeMax, _ := strconv.Atoi(idRange[1])

					var intraResourceMgmtMapData examplecomv1.FunctionsIntraResourceMgmtMap
					for idIndex := idRangeMin; idIndex <= idRangeMax; idIndex++ {
						idIndexString := strconv.Itoa(idIndex)
						if nil == (*specFunctionData.IntraResourceMgmtMap)[idIndexString].FunctionCRName {
							upIntraResourceMgmtMapData[idIndexString] =
								(*specFunctionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)]
							continue
						}
						if false == *(*specFunctionData.IntraResourceMgmtMap)[idIndexString].Available &&
							*(*specFunctionData.IntraResourceMgmtMap)[idIndexString].FunctionCRName == req.Name {
							available := false
							intraResourceMgmtMapData.Available = &available
							intraResourceMgmtMapData.FunctionCRName = nil
							upIntraResourceMgmtMapData[idIndexString] = intraResourceMgmtMapData
						} else {
							upIntraResourceMgmtMapData[idIndexString] =
								(*specFunctionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)]
						}
					}

					*specFunctionData.IntraResourceMgmtMap = upIntraResourceMgmtMapData
					*statusFunctionData.IntraResourceMgmtMap = upIntraResourceMgmtMapData
					upSpecFunctionsData = append(upSpecFunctionsData, specFunctionData)
					upStatusFunctionsData = append(upStatusFunctionsData, statusFunctionData)
				}
			}
			specRegion.Modules.Functions = &upSpecFunctionsData
			statusRegion.Modules.Functions = &upStatusFunctionsData
		}
	}

	err = r.updChildBsCR(ctx,
		childBitstreamCRData,
		examplecomv1.ChildBsReady,
		examplecomv1.ChildBsStatusReady,
		childBitstreamCRBase)
	if err != nil {
		logger.Error(err, "ChildBsCR update failed.")
	}
	return err
}

// reallocate FPGA Resource
func (r *FPGAFunctionReconciler) reallocateFPGAResource(ctx context.Context,
	childBitstreamCRData *examplecomv1.ChildBs) {

	// clear the all region info and all function modual info FunctionsIntraResourceMgmtMap
	for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
		specRegion := childBitstreamCRData.Spec.Regions[regionIndex]
		statusRegion := childBitstreamCRData.Status.Regions[regionIndex]

		var upSpecFunctionsData []examplecomv1.ChildBsFunctions
		var upStatusFunctionsData []examplecomv1.ChildBsFunctions
		for functionIndex := 0; functionIndex < len(*specRegion.Modules.Functions); functionIndex++ {

			specFunctionData := (*specRegion.Modules.Functions)[functionIndex]
			statusFunctionData := (*specRegion.Modules.Functions)[functionIndex]

			for funcmoduleIndex := 0; funcmoduleIndex < len(*specFunctionData.Module); funcmoduleIndex++ {
				functionChannelIDs := (*specFunctionData.Module)[funcmoduleIndex].FunctionChannelIDs

				if nil == functionChannelIDs {
					continue
				}

				var upIntraResourceMgmtMapData map[string]examplecomv1.FunctionsIntraResourceMgmtMap
				upIntraResourceMgmtMapData = make(map[string]examplecomv1.FunctionsIntraResourceMgmtMap)

				idRange := strings.SplitN(*functionChannelIDs, "-", 2)
				idRangeMin, _ := strconv.Atoi(idRange[0])
				idRangeMax, _ := strconv.Atoi(idRange[1])

				var intraResourceMgmtMapData examplecomv1.FunctionsIntraResourceMgmtMap
				for idIndex := idRangeMin; idIndex <= idRangeMax; idIndex++ {
					idIndexString := strconv.Itoa(idIndex)
					available := true
					intraResourceMgmtMapData.Available = &available
					intraResourceMgmtMapData.FunctionCRName = nil
					upIntraResourceMgmtMapData[idIndexString] = intraResourceMgmtMapData
				}

				*specFunctionData.IntraResourceMgmtMap = upIntraResourceMgmtMapData
				*statusFunctionData.IntraResourceMgmtMap = upIntraResourceMgmtMapData
				upSpecFunctionsData = append(upSpecFunctionsData, specFunctionData)
				upStatusFunctionsData = append(upStatusFunctionsData, statusFunctionData)
			}
		}
		specRegion.Modules.Functions = &upSpecFunctionsData
		statusRegion.Modules.Functions = &upStatusFunctionsData
	}

	return
}

// Get Connection Data
func (r *FPGAFunctionReconciler) GetConnectionData(ctx context.Context,
	connectionCRNamespace string,
	connectionCRName string,
	connectionKind *string) error {
	logger := log.FromContext(ctx)

	var err error

	var existCount int = 0
	var notFoundErr error
	var elseErr error
	var connectionKindList []string = []string{examplecomv1.ConnectionCRKindPCIe,
		examplecomv1.ConnectionCRKindEth}

	connectionCRData := &unstructured.Unstructured{}

	for n := 0; n < len(connectionKindList); n++ {
		connectionCRData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    connectionKindList[n],
		})
		err = r.Get(ctx, client.ObjectKey{
			Namespace: connectionCRNamespace,
			Name:      connectionCRName}, connectionCRData)
		if errors.IsNotFound(err) {
			existCount += 1
			notFoundErr = err
		} else if err != nil {
			logger.Info("unable to fetch ConnectionCR")
			elseErr = err
		} else {
			*connectionKind = connectionKindList[n]
			break
		}
		if existCount == len(connectionKindList) {
			logger.Info("ConnectionCR does not exist")
			err = notFoundErr
		} else {
			err = elseErr
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// Generate ConnectionCR name
func MakeConnectionCRName(
	ctx context.Context,
	dataFlowName string,
	fromCRName string,
	toCRName string,
	connectionCRName *string) {

	var concatenationFlag bool = false
	var fromFunctionName string
	var toFunctionName string

	fromCRNameSplitList := strings.Split(fromCRName, "-")
	for splitIndex := 0; splitIndex < len(fromCRNameSplitList); splitIndex++ {
		if true == concatenationFlag {
			fromFunctionName = fromFunctionName +
				fromCRNameSplitList[splitIndex] + "-"
		}

		if "wbfunction" == fromCRNameSplitList[splitIndex] {
			concatenationFlag = true
		}

	}
	concatenationFlag = false
	toCRNameSplitList := strings.Split(toCRName, "-")
	for splitIndex := 0; splitIndex < len(toCRNameSplitList); splitIndex++ {
		if "" != toFunctionName {
			toFunctionName = toFunctionName + "-"
		}

		if true == concatenationFlag {
			toFunctionName = toFunctionName +
				toCRNameSplitList[splitIndex]
		}

		if "wbfunction" == toCRNameSplitList[splitIndex] {
			concatenationFlag = true
		}

	}
	concatenationFlag = false

	*connectionCRName = dataFlowName + "-wbconnection-" + fromFunctionName + toFunctionName
}
