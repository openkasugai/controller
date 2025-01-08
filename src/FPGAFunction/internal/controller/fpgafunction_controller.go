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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// Hold own node information
var myNodeName string

//+kubebuilder:rbac:groups=example.com,resources=fpgas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=fpgas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=fpgas/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=childbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=childbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=childbs/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=fpgafunctions/finalizers,verbs=update

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

	var crData examplecomv1.FPGAFunction
	var eventKind int32 // 0:Add, 1:Upd,  2:Del
	var retCInt C.int
	breakFlag := false

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
			var frameSizeConfigCChar *C.char

			for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
				// Get config information
				err = r.getConfigData(ctx, crData.Spec.ConfigName, &cfgData)
				if nil == err {
					err = FunctionConfigDataJsonUnmarshal(&cfgData, &functionConfigData)
					if nil != err {
						logger.Error(err, "unable to unmarshal. ConfigMap="+crData.Spec.ConfigName)
						break
					}
				} else {
					break
				}

				if "" != crData.Spec.AcceleratorIDs[0].ID {
					deviceFilePath := crData.Spec.AcceleratorIDs[0].ID
					deviceUUID = strings.ReplaceAll(deviceFilePath, "/dev/xpcie_", "")
					// deviceUUID = strings.Replace(deviceFilePath, "/dev/xpcie_", "", -1)
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
				} else {
					logger.Info("FPGAFunction.Spec.AcceleratorIDs[0].ID: " +
						crData.Spec.AcceleratorIDs[0].ID)
					break
				}

				if crData.Spec.FunctionIndex == nil &&
					fpgaCRData.Status.ChildBitstreamCRName == nil {

					err = r.createChildBsCR(ctx,
						functionConfigData,
						crData,
						&childBitstreamCRData,
						&fpgaCRData)
					if nil != err {
						break
					}

					err = updFPGACR(ctx, r,
						&fpgaCRData,
						childBitstreamCRData.Spec.ChildBitstreamID)
					if nil != err {
						break
					}

					if "" != crData.Spec.AcceleratorIDs[0].ID {
						deviceFilePath := crData.Spec.AcceleratorIDs[0].ID
						deviceUUID = strings.ReplaceAll(deviceFilePath, "/dev/xpcie_", "")
					} else {
						break
					}
					fpgaDeviceUUIDCString := C.CString(deviceUUID)
					childBsFile := C.CString(functionConfigData.ChildBitstream.File)
					retCInt = C.fpga_get_dev_id(fpgaDeviceUUIDCString, &deviceID)
					if 0 > retCInt {
						logger.Info("fpga_get_dev_id() err = " +
							strconv.Itoa(int(retCInt)))
						break
					} else {
						logger.Info("fpga_get_dev_id() ret = " +
							strconv.Itoa(int(retCInt)))
					}

					logger.Info("Start to write Child-Bitstream.")

					retCInt = C.fpga_write_bitstream(deviceID, 0, childBsFile)
					if 0 > retCInt {
						logger.Info("fpga_write_bitstream() err = " +
							strconv.Itoa(int(retCInt)))
						break
					} else {
						logger.Info("fpga_write_bitstream() ret = " +
							strconv.Itoa(int(retCInt)))
					}

					err = r.getChildBsData(ctx,
						fpgaCRName,
						*fpgaCRData.Status.ChildBitstreamID,
						&childBitstreamCRData)
					if errors.IsNotFound(err) {
						// CR does not exist
						logger.Info("NotFound to fetch CR")
						break
					} else if err != nil {
						logger.Error(err, "unable to fetch CR")
						break
					}

					err = updChildBsCR(ctx,
						r,
						&childBitstreamCRData,
						examplecomv1.ChildBsConfiguringParam,
						examplecomv1.ChildBsStatusPreparing)
					if nil != err {
						break
					}

					logger.Info("Bitstream file writing has completed successfully.")

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
						// lldmaExtIfID := *region.Modules.LLDMA.ExtIfID
						// ptuExtIfID := *region.Modules.Ptu.ExtIfID
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

							FunctionsData := (*region.Modules.Functions)[functionIndex]
							for funcmodIndex := 0; funcmodIndex < len(*FunctionsData.Module); funcmodIndex++ {
								funcNameCString := C.CString(*(*FunctionsData.Module)[funcmodIndex].Type)

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

								var frameSizeStruct examplecomv1.FrameSizeData
								frameSizeStruct.InputWidth = functionConfigData.Parameters.Functions.InputWidth
								frameSizeStruct.InputHeight = functionConfigData.Parameters.Functions.InputHeight
								frameSizeStruct.OutputWidth = functionConfigData.Parameters.Functions.OutputWidth
								frameSizeStruct.OutputHeight = functionConfigData.Parameters.Functions.OutputHeight
								frameSizeBytes, err := json.Marshal(frameSizeStruct)
								if nil != err {
									logger.Error(err, "FrameSize unable to Marshal.")
									breakFlag = true
									break
								}
								frameSizeCString := C.CString(string(frameSizeBytes))

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
									breakFlag = true
									break
								}
								if true == breakFlag {
									break
								}
							}
						}
						if true == breakFlag {
							break
						}
					}

					var nicFlag bool
					err = updDeployinfoCM(ctx,
						r,
						deviceUUID,
						crData.Spec.NodeName,
						crData.Spec.FunctionName,
						&nicFlag,
						&childBitstreamCRData.Spec.Regions)
					if nil != err {
						logger.Error(err, "DeployInfoCM Update Error")
						break
					}

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
					}
					childBitstreamCRData.Status = upChildBitstreamCRData.Status

					if false == nicFlag {
						updChildBsCR(ctx,
							r,
							&childBitstreamCRData,
							examplecomv1.ChildBsReady,
							examplecomv1.ChildBsStatusReady)
						if nil != err {
							breakFlag = true
							break
						} else {
							logger.Info("This Child-Bitstream has no NIC.")
						}
					} else {
						updChildBsCR(ctx,
							r,
							&childBitstreamCRData,
							examplecomv1.ChildBsNoConfigureNetwork,
							examplecomv1.ChildBsStatusPreparing)
						if nil != err {
							breakFlag = true
							break
						} else {
							logger.Info("Application parameters settings for all lanes has completed successfully.")
						}
					}
				}

				err = r.getChildBsData(ctx,
					fpgaCRName,
					*fpgaCRData.Status.ChildBitstreamID,
					&childBitstreamCRData)
				if errors.IsNotFound(err) {
					// CR does not exist
					logger.Info("NotFound to fetch CR")
					break
				} else if err != nil {
					logger.Error(err, "unable to fetch CR")
					break
				}

				// logger.Info("debug fpgaCRData : ", fpgaCRData)
				// logger.Info("debug fpgaCRData : ")
				// fmt.Println("debug childBitstreamCRData : ", childBitstreamCRData.Status.State)
				// logger.Info("debug childBitstreamCRData : "+ examplecomv1.ChildBsReady)
				if examplecomv1.ChildBsReady == childBitstreamCRData.Status.State &&
					examplecomv1.ChildBsStatusReady == childBitstreamCRData.Status.Status {

					var fpgafuncRxData examplecomv1.RxTxData
					var fpgafuncTxData examplecomv1.RxTxData
					var functionKernelID int32
					var chainID int32
					var functionChannelID int32
					err = allocateFPGAResource(ctx,
						r,
						req,
						&crData,
						&childBitstreamCRData,
						&fpgafuncRxData,
						&fpgafuncTxData,
						&functionKernelID,
						&chainID,
						&functionChannelID)
					if nil != err {
						break
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
					crData.Status.ParentBitstreamName = functionConfigData.ParentBitstream.ID
					crData.Status.ChildBitstreamName = functionConfigData.ChildBitstream.ID
					crData.Status.FunctionChannelID = functionChannelID
					crData.Status.Rx = fpgafuncRxData
					crData.Status.Tx = fpgafuncTxData
					crData.Status.FrameworkKernelID = chainID
					crData.Status.FunctionKernelID = functionKernelID
					crData.Status.PtuKernelID = chainID

					crData.Status.AcceleratorStatuses = append(crData.Status.AcceleratorStatuses, acceleratorStatuses)

					r.UpdCustomResource(ctx, r, &crData, RUNNING)
				} else {

					logger.Info("ChildBs.Status.State: " + string(childBitstreamCRData.Status.State) +
						", ChildBs.Status: " + string(childBitstreamCRData.Status.Status))
					return ctrl.Result{Requeue: true}, nil
				}

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
			err = r.DelCustomResource(ctx, r, &crData)
			if err != nil {
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			r.Recorder.Eventf(&crData, corev1.EventTypeNormal, "Delete", "Delete End")
		}
	}

	if nil != err {
		return ctrl.Result{}, err
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
		Complete(r)
}

func (r *FPGAFunctionReconciler) GetfinalizerName(pCRData *examplecomv1.FPGAFunction) string {
	// Value to set in the finalizer
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
}

/*
func GetfinalizerName(pCRData *examplecomv1.FPGAFunction) string {
	// Value to set in the finalizer
	return strings.ToLower(pCRData.Kind) + ".finalizers." +
		strings.ReplaceAll(pCRData.APIVersion, "/", ".")
}
*/

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

func (r *FPGAFunctionReconciler) UpdCustomResource(ctx context.Context, pClient *FPGAFunctionReconciler,
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
	err = pClient.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "FPGAFunctionCR Status Update Error.")
	} else {
		logger.Info("FPGAFunctionCR Status Update.")
	}
	return err
}

func (r *FPGAFunctionReconciler) DelCustomResource(ctx context.Context, pClient *FPGAFunctionReconciler,
	pCRData *examplecomv1.FPGAFunction) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	// Delete the Finalizer statement.
	if controllerutil.ContainsFinalizer(pCRData,
		r.GetfinalizerName(pCRData)) {
		controllerutil.RemoveFinalizer(pCRData, r.GetfinalizerName(pCRData))
		err := pClient.Update(ctx, pCRData)
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
		logger.Error(err, "ConfigMap does not exist")
	}
	if err != nil {
		logger.Error(err, "unable to fetch ConfigMap")
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
func StartupProccessing(mng ctrl.Manager) {

	ctx := context.Background()

	myNodeName = os.Getenv("K8S_NODENAME")

	fpgaList := GetFPGACR(ctx, mng, myNodeName)

	if 0 != len(fpgaList) {
		// FPGAPhase3 initial setting
		var argv []*C.char
		// Phase 3 FPGA initialization variables
		argv = []*C.char{C.CString("proc"),
			C.CString("-d"),
			C.CString(strings.Join(fpgaList, ","))}
		argc := C.int(len(argv))
		C.libfpga_log_set_output_stdout()
		C.libfpga_log_set_level(C.LIBFPGA_LOG_ALL)
		C.fpga_init(argc, (**C.char)(unsafe.Pointer(&argv[0])))

		// Hold device information
		for deviceID, _ := range fpgaList {
			C.fpga_enable_regrw(C.uint(deviceID))
			// FpgaDevList = append(FpgaDevList, devPath)
		}
	}
}

// Config information storage area
var gServicerMgmtInfo []examplecomv1.ServicerMgmtInfo
var gDeployInfo map[string][]examplecomv1.DeviceRegionInfo
var gRegionUniqueInfo []examplecomv1.RegionSpecificInfo
var gFunctionUniqueInfo []examplecomv1.FPGACatalog
var gFilterResizeInfo []examplecomv1.FunctionDetail

type ConfigTable struct {
	name string
}

var configLoadTableForWriteChildBs = []ConfigTable{
	{examplecomv1.CMServicerMgmtInfo},
	{examplecomv1.CMDeployInfo},
	{examplecomv1.CMRegionUniqueInfo},
	{examplecomv1.CMFunctionUniqueInfo},
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

// Create ChildBsCR
func (r *FPGAFunctionReconciler) createChildBsCR(ctx context.Context,
	functionConfigData examplecomv1.FPGAFuncConfig,
	crData examplecomv1.FPGAFunction,
	childBitstreamCRData *examplecomv1.ChildBs,
	fpgaCRData *examplecomv1.FPGA) error {
	logger := log.FromContext(ctx)

	var err error
	var bitstreamIDConfigCChar *C.char
	var bitstreamIDConfig examplecomv1.BsConfigInfo

	defer C.free(unsafe.Pointer(bitstreamIDConfigCChar))

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		parentBitstreamCString := C.CString(functionConfigData.ParentBitstream.ID)
		childBitstreamCString := C.CString(functionConfigData.ChildBitstream.ID)
		ret := C.fpga_db_get_device_config_by_bitstream_id(parentBitstreamCString,
			childBitstreamCString,
			(**C.char)(unsafe.Pointer(&bitstreamIDConfigCChar))) //nolint:gocritic // suspicious identical LHS and RHS for `==` operator
		if 0 > ret {
			logger.Info("fpga_get_device_config_by_bitstream_id()" +
				" err = " + strconv.Itoa(int(ret)))
			break
		} else {
			logger.Info("fpga_get_device_config_by_bitstream_id()" +
				" ret = " + strconv.Itoa(int(ret)))
		}

		if 0 == ret {

			// Get the number of characters in the JSON data
			n := 0
			head := (*byte)(unsafe.Pointer(bitstreamIDConfigCChar))
			for ptr := head; *ptr != 0; n++ {
				ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
			}

			// Convert to []bytes type
			bitstreamIDConfigBytes := C.GoBytes(unsafe.Pointer(bitstreamIDConfigCChar), C.int(n))

			// Get FPGA function information
			functionError := json.Unmarshal(bitstreamIDConfigBytes, &bitstreamIDConfig)
			if nil != functionError {
				logger.Info("unable to unmarshal fpga_get_device_config_by_bitstream_id() data.")
			}

			childBitstreamCRData.Spec = bitstreamIDConfig.ChildBitstreamIDs[0]
		}

		err = r.getConfigMapForWriteChildBs(ctx)
		if nil != err {
			logger.Error(err, "unable to unmarshal. ConfigMap="+crData.Spec.ConfigName)
			break
		}

		var ptuData *examplecomv1.ChildBsPtu
		var functionData examplecomv1.ChildBsFunctions

		var networkInfo examplecomv1.NetworkData
		var ioFrameSize examplecomv1.FrameSizeData
		var setAvailable bool = true

		for regionIndex := 0; regionIndex < len(childBitstreamCRData.Spec.Regions); regionIndex++ {
			region := childBitstreamCRData.Spec.Regions[regionIndex]

			for srvIndex := 0; srvIndex < len(gServicerMgmtInfo); srvIndex++ {
				if nil == region.Modules.Ptu {
					continue
				}
				if gServicerMgmtInfo[srvIndex].NodeName != crData.Spec.NodeName {
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
					if networkInfo.LaneIndex != int32(lane) {
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
						(*functionData.IntraResourceMgmtMap)[strconv.Itoa(idIndex)] = examplecomv1.FunctionsIntraResourceMgmtMap{Available: &available}
					}
				}

				ioFrameSize = functionConfigData.Parameters.Functions
				if nil == functionData.Parameters {
					parameters := make(map[string]intstr.IntOrString)
					functionData.Parameters = &parameters
				}
				(*functionData.Parameters)["InputHeight"] = intstr.IntOrString{Type: intstr.Int,
					IntVal: ioFrameSize.InputHeight}
				(*functionData.Parameters)["InputWidth"] = intstr.IntOrString{Type: intstr.Int,
					IntVal: ioFrameSize.InputWidth}
				(*functionData.Parameters)["OutputHeight"] = intstr.IntOrString{Type: intstr.Int,
					IntVal: ioFrameSize.OutputHeight}
				(*functionData.Parameters)["OutputWidth"] = intstr.IntOrString{Type: intstr.Int,
					IntVal: ioFrameSize.OutputWidth}

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
	crData *examplecomv1.ChildBs) error {
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
	}
	return err
}

// ChildBs CR Update
func updChildBsCR(ctx context.Context, pClient *FPGAFunctionReconciler,
	pCRData *examplecomv1.ChildBs,
	state examplecomv1.ChildBitstreamState,
	status examplecomv1.ChildBitstreamStatus) error {
	logger := log.FromContext(ctx)
	var err error

	if examplecomv1.ChildBsReady == state {
		pCRData.Status.Status = status
	}
	pCRData.Status.State = state
	err = pClient.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "ChildBitstreamCR Status Update Error.")
	} else {
		logger.Info("ChildBitstreamCR Status Update.")
	}
	return err
}

// FPGA CR Update
func updFPGACR(ctx context.Context, pClient *FPGAFunctionReconciler,
	pCRData *examplecomv1.FPGA, childBsID *string) error {
	logger := log.FromContext(ctx)
	var err error

	crName := pCRData.ObjectMeta.Name + "-" + *childBsID

	pCRData.Spec.ChildBitstreamID = childBsID
	pCRData.Status.ChildBitstreamID = childBsID
	pCRData.Status.ChildBitstreamCRName = &crName

	err = pClient.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "FPGACR Status Update Error.")
	} else {
		logger.Info("FPGACR Status Update.")
	}
	return err
}

// deployinfo CM Update
func updDeployinfoCM(ctx context.Context, pClient *FPGAFunctionReconciler,
	deviceUUID string,
	nodeName string,
	functionName string,
	nicFlag *bool,
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

				if "lane0" == functionTarget.RegionName {
					typeList := strings.SplitN(functionTarget.RegionType, "-", 4)
					if "0nics" != typeList[3] {
						*nicFlag = true
					}
				}

				var childBsFunctions []examplecomv1.ChildBsFunctions
				var maxCapacity int32
				var maxFunctions int32
				// var functions *[]examplecomv1.ChildBsFunctions
				for bsIndex := 0; bsIndex < len(*childBsRegions); bsIndex++ {
					if *(*childBsRegions)[bsIndex].Name == functionTarget.RegionName {
						childBsFunctions =
							*(*childBsRegions)[bsIndex].Modules.Functions
						// childBsFunctions = *functions
						maxCapacity = *(*childBsRegions)[bsIndex].MaxCapacity
						maxFunctions = *(*childBsRegions)[bsIndex].MaxFunctions
					}
				}

				for funcIndex := 0; funcIndex < len(childBsFunctions); funcIndex++ {
					var functionData examplecomv1.SimpleFunctionInfraStruct

					function := (childBsFunctions)[funcIndex]

					partitionName := strconv.Itoa(int(*function.ID))
					funcIndexInt32 := int32(funcIndex)

					functionData.PartitionName = partitionName
					functionData.FunctionIndex = &funcIndexInt32
					functionData.FunctionName = functionName
					functionData.MaxCapacity = *function.DeploySpec.MaxCapacity
					functionData.MaxDataFlows = *function.DeploySpec.MaxDataFlows

					functionTargetData.Functions = append(functionTargetData.Functions, functionData)

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

		err = pClient.Update(ctx, createCM)
		if err != nil {
			logger.Error(err, examplecomv1.CMDeployInfo+" Update Error.")
		} else {
			logger.Info(examplecomv1.CMDeployInfo + " Update.")
		}
	}
	return err
}

func getRxTxProtocol(ctx context.Context,
	pClient *FPGAFunctionReconciler,
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
			err = GetConnectionData(ctx,
				pClient,
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
			err = GetConnectionData(ctx,
				pClient,
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
func allocateFPGAResource(ctx context.Context,
	pClient *FPGAFunctionReconciler,
	req ctrl.Request,
	crData *examplecomv1.FPGAFunction,
	childBitstreamCRData *examplecomv1.ChildBs,
	fpgafuncRxData *examplecomv1.RxTxData,
	fpgafuncTxData *examplecomv1.RxTxData,
	functionKernelID *int32,
	chainID *int32,
	functionChannelID *int32) error {
	logger := log.FromContext(ctx)

	var err error
	var rxProtocol string
	var txProtocol string
	var upChildBitstreamCRData examplecomv1.ChildBs

	err = getRxTxProtocol(ctx,
		pClient,
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

					var upIntraResourceMgmtMapData map[string]examplecomv1.FunctionsIntraResourceMgmtMap
					upIntraResourceMgmtMapData = make(map[string]examplecomv1.FunctionsIntraResourceMgmtMap)

					idRange := strings.SplitN(*functionChannelIDs, "-", 2)
					idRangeMin, _ := strconv.Atoi(idRange[0])
					idRangeMax, _ := strconv.Atoi(idRange[1])

					var intraResourceMgmtMapData examplecomv1.FunctionsIntraResourceMgmtMap
					var allocateFunctionChannelID int32 = -1
					var functionChannelIDDetail examplecomv1.FPGAFunctionChannelIDs
					for idIndex := idRangeMin; idIndex <= idRangeMax; idIndex++ {
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

					for frIndex := 0; frIndex < len(gFilterResizeInfo); frIndex++ {
						functionDetailData := gFilterResizeInfo[frIndex]
						if strconv.Itoa(int(*specFunctionData.ID)) == functionDetailData.PartitionName {
							for i := range functionDetailData.FunctionChannelIDs {
								functionChannelIDsData := functionDetailData.FunctionChannelIDs[i]
								if functionChannelIDsData.FunctionChannelID == allocateFunctionChannelID {
									functionChannelIDDetail = functionChannelIDsData
								}
							}
						}
					}

					if -1 != allocateFunctionChannelID &&
						allocateFunctionChannelID == functionChannelIDDetail.FunctionChannelID {

						available := false

						var txData examplecomv1.RxTxSpec
						var rxData examplecomv1.RxTxSpec
						var txProtocolDatails map[string]examplecomv1.Details = make(map[string]examplecomv1.Details)
						var rxProtocolDatails map[string]examplecomv1.Details = make(map[string]examplecomv1.Details)
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

						txProtocolDatails[txProtocol] = examplecomv1.Details{
							DMAChannelID:     txDMAChannelID,
							LLDMAConnectorID: txLLDMAConnectorID,
							Port:             txPort}
						rxProtocolDatails[rxProtocol] = examplecomv1.Details{
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

	err = updChildBsCR(ctx,
		pClient,
		childBitstreamCRData,
		examplecomv1.ChildBsReady,
		examplecomv1.ChildBsStatusReady)
	if err != nil {
		logger.Error(err, "ChildBsCR update failed.")
		// do nothing
	}
	return err
}

// Get Connection Data
func GetConnectionData(ctx context.Context,
	pClient *FPGAFunctionReconciler,
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
		err = pClient.Get(ctx, client.ObjectKey{
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

	// logger := log.FromContext(ctx)

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
