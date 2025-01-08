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

	examplecomv1 "PCIeConnection/api/v1"

	/* Additional imports */
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
	"unsafe"
	// #cgo LDFLAGS: -lpciaccess
	// #cgo pkg-config: libdpdk
	// #cgo CFLAGS:  -I/usr/local/include/fpgalib/
	// #cgo CFLAGS:  -I/usr/local/include/fpgalib/dpdk/include/
	// #cgo LDFLAGS: -L/usr/local/lib/fpgalib/ -lfpga
	// #cgo CXXFLAGS: -std=c++11
	// #cgo LDFLAGS: -L. -lstdc++
	// #cgo LDFLAGS: -L. -lpciaccess
	// #include <liblldma.h>
	// #include <libfpgactl.h>
	// #include <libchain.h>
	// #include <libdmacommon.h>
	// #include <liblogging.h>
	// #include <libshmem.h>
	// #include <libshmem_controller.h>
	"C"
)

// PCIeConnectionReconciler reconciles a PCIeConnection object
type PCIeConnectionReconciler struct {
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

// Status type
const (
	STATUS_OK   = "OK"
	STATUS_INIT = "INIT"
	STATUS_NG   = "NG"
)

// FunctionType type
const (
	FUNCTYPE_FPGA = "FPGAFunction"
	FUNCTYPE_GPU  = "GPUFunction"
	FUNCTYPE_GATE = "GateFunction"
	FUNCTYPE_CPU  = "CPUFunction"
)

// Overall Status type
const (
	PENDING = "Pending"
	RUNNING = "Running"
)

var FpgaDevList = []string{}
var ShmemEnable = map[string]bool{}
var ShmemInit = map[string]bool{}

const Hugepagesz int32 = 0x100000 // MB

// Structure for obtaining Function information
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   examplecomv1.FunctionData       `json:"spec,omitempty"`
	Status examplecomv1.FunctionStatusData `json:"status,omitempty"`
}

//+kubebuilder:rbac:groups=example.com,resources=pcieconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=pcieconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=pcieconnections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PCIeConnection object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *PCIeConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var crData examplecomv1.PCIeConnection /* Structure for holding CR information */
	var eventKind int32                    /* 0:Add, 1:Upd,  2:Del */
	var srcDeviceID int32                  /* Src side FPGA device number */
	var dstDeviceID int32                  /* Device number of the FPGA on the destination side */
	var crStatus string
	var srcFunctionData Function
	var dstFunctionData Function
	var srcFunctionKind string
	var dstFunctionKind string

	crStatus = PENDING
	srcDeviceID = -1
	dstDeviceID = -1
	srcFunctionKind = ""
	dstFunctionKind = ""

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

	// Get Event type
	eventKind = r.GetEventKind(&crData)
	if eventKind == CREATE {
		// FunctionData acquisition (Src side)
		err = r.GetFunctionData(ctx, crData.Spec.From.WBFunctionRef, &srcFunctionData, &srcFunctionKind)
		if err != nil {
			if errors.IsNotFound(err) {
				// do nothing
			} else {
				logger.Error(err, "GetFunctionData() error")
			}
			return ctrl.Result{Requeue: true}, nil
		}

		// FunctionData acquisition (Dst side)
		err = r.GetFunctionData(ctx, crData.Spec.To.WBFunctionRef, &dstFunctionData, &dstFunctionKind)
		if err != nil {
			if errors.IsNotFound(err) {
				// do nothing
			} else {
				logger.Error(err, "GetFunctionData() error")
			}
			return ctrl.Result{Requeue: true}, nil
		}

		// Collect information about your own node
		myNodeName := os.Getenv("K8S_NODENAME")

		if (myNodeName == srcFunctionData.Spec.NodeName) ||
			(myNodeName == dstFunctionData.Spec.NodeName) {
			// do nothing
		} else {
			// Do nothing except the target worker node
			return ctrl.Result{}, nil
		}

		// For creation
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Create", "Create Start")

		// Setting process
		var srcAccID string
		if FUNCTYPE_FPGA == srcFunctionKind {
			if nil == srcFunctionData.Status.FunctionKernelID {
				return ctrl.Result{Requeue: true}, nil
			}
			for count := 0; count < len(srcFunctionData.Spec.AcceleratorIDs); count++ {
				if nil == srcFunctionData.Spec.AcceleratorIDs[count].PartitionName {
					srcAccID = srcFunctionData.Spec.AcceleratorIDs[count].ID
					break
				} else if *srcFunctionData.Spec.AcceleratorIDs[count].PartitionName ==
					strconv.Itoa(int(*srcFunctionData.Status.FunctionKernelID)) {
					srcAccID = srcFunctionData.Spec.AcceleratorIDs[count].ID
					break
				}
			}
			srcDeviceID = PCIeConnectionAccIDToDeviceID(srcAccID)
		}

		// Setting process
		var dstAccID string
		if FUNCTYPE_FPGA == dstFunctionKind {
			if nil == dstFunctionData.Status.FunctionKernelID {
				return ctrl.Result{Requeue: true}, nil
			}
			for count := 0; count < len(dstFunctionData.Spec.AcceleratorIDs); count++ {
				if nil == dstFunctionData.Spec.AcceleratorIDs[count].PartitionName {
					dstAccID = dstFunctionData.Spec.AcceleratorIDs[count].ID
					break
				} else if *dstFunctionData.Spec.AcceleratorIDs[count].PartitionName ==
					strconv.Itoa(int(*dstFunctionData.Status.FunctionKernelID)) {
					dstAccID = dstFunctionData.Spec.AcceleratorIDs[count].ID
					break
				}
			}
			dstDeviceID = PCIeConnectionAccIDToDeviceID(dstAccID)
		}

		if myNodeName == srcFunctionData.Spec.NodeName {
			if srcFunctionData.Status.Status == RUNNING &&
				dstFunctionData.Status.Status == RUNNING {
				switch srcFunctionKind {
				case FUNCTYPE_FPGA:
					if FUNCTYPE_FPGA == dstFunctionKind {
						// Call D2D processing function
						ret := PCIeConnectionD2D(ctx, &crData,
							srcDeviceID, &srcFunctionData.Status,
							dstDeviceID, &dstFunctionData.Status)
						if STATUS_OK == ret {
							crData.Status.From.Status = STATUS_OK
							crData.Status.To.Status = STATUS_OK
						} else {
							crData.Status.From.Status = STATUS_NG
							crData.Status.To.Status = STATUS_NG
						}
					} else {
						if done, exist := ShmemEnable[dstFunctionData.Spec.SharedMemory.FilePrefix]; !exist || (exist && !done) {
							// Call the prefix enable function
							var size []C.uint
							size = append(size, C.uint(dstFunctionData.Spec.SharedMemory.SharedMemoryMiB)) //nolint:staticcheck // 'size' is never used. FIXME: remove variable
							ret := C.fpga_shmem_enable(
								C.CString(dstFunctionData.Spec.SharedMemory.FilePrefix), nil)
							logger.Info("debug prefix = " +
								dstFunctionData.Spec.SharedMemory.FilePrefix)
							if ret == 0 {
								ShmemEnable[dstFunctionData.Spec.SharedMemory.FilePrefix] = true
								logger.Info("fpga_shmem_enable() OK ret = " +
									strconv.Itoa(int(ret)))
								crData.Status.From.Status = STATUS_OK
							} else {
								ShmemEnable[dstFunctionData.Spec.SharedMemory.FilePrefix] = false
								logger.Info("fpga_shmem_enable() NG ret = " +
									strconv.Itoa(int(ret)))
								crData.Status.From.Status = STATUS_NG
							}
						}
						// Register connection ID and call FDMA setting function
						if crData.Status.From.Status != STATUS_NG {
							crData.Status.From.Status =
								PCIeConnectionSrcFPGA(ctx,
									srcDeviceID, &srcFunctionData.Status, &dstFunctionData.Status)
						}
					}
				case FUNCTYPE_CPU, FUNCTYPE_GPU:
					if FUNCTYPE_FPGA == dstFunctionKind {
						logger.Info("nothing to do")
						crData.Status.From.Status = STATUS_OK
					} else if FUNCTYPE_CPU == dstFunctionKind || FUNCTYPE_GPU == dstFunctionKind {
						if done, exist := ShmemEnable[srcFunctionData.Spec.SharedMemory.FilePrefix]; !exist || (exist && !done) {
							// Call the prefix enable function
							ret := C.fpga_shmem_enable(
								C.CString(srcFunctionData.Spec.SharedMemory.FilePrefix), nil)
							logger.Info("debug prefix = " +
								srcFunctionData.Spec.SharedMemory.FilePrefix)
							if ret == 0 {
								ShmemEnable[srcFunctionData.Spec.SharedMemory.FilePrefix] = true
								logger.Info("fpga_shmem_enable() OK ret = " +
									strconv.Itoa(int(ret)))
								crData.Status.From.Status = STATUS_OK
							} else {
								ShmemEnable[srcFunctionData.Spec.SharedMemory.FilePrefix] = false
								logger.Info("fpga_shmem_enable() NG ret = " +
									strconv.Itoa(int(ret)))
								crData.Status.From.Status = STATUS_NG
							}
						}
					} else {
						logger.Info("nothing to do")
					}
				default:
					logger.Info("nothing to do")
					crData.Status.From.Status = STATUS_OK
				}
			}
		}

		if myNodeName == dstFunctionData.Spec.NodeName {
			if srcFunctionData.Status.Status == RUNNING &&
				dstFunctionData.Status.Status == RUNNING {
				switch dstFunctionKind {
				case FUNCTYPE_FPGA:
					if FUNCTYPE_FPGA != srcFunctionKind {
						if done, exist := ShmemEnable[srcFunctionData.Spec.SharedMemory.FilePrefix]; !exist || (exist && !done) {
							// Call the prefix enable function
							var size []C.uint
							size = append(size, C.uint(srcFunctionData.Spec.SharedMemory.SharedMemoryMiB)) //nolint:staticcheck // 'size' is never used. FIXME: remove variable
							ret := C.fpga_shmem_enable(
								C.CString(srcFunctionData.Spec.SharedMemory.FilePrefix), nil)
							logger.Info("debug prefix = " +
								srcFunctionData.Spec.SharedMemory.FilePrefix)
							if ret == 0 {
								ShmemEnable[srcFunctionData.Spec.SharedMemory.FilePrefix] = true
								logger.Info("fpga_shmem_enable() OK ret = " +
									strconv.Itoa(int(ret)))
								crData.Status.To.Status = STATUS_OK
							} else {
								ShmemEnable[srcFunctionData.Spec.SharedMemory.FilePrefix] = false
								logger.Info("fpga_shmem_enable() NG ret = " +
									strconv.Itoa(int(ret)))
								crData.Status.To.Status = STATUS_NG
							}
						}
						// Register connection ID and call FDMA setting function
						if crData.Status.To.Status != STATUS_NG {
							crData.Status.To.Status =
								PCIeConnectionDstFPGA(ctx,
									dstDeviceID, &srcFunctionData.Status, &dstFunctionData.Status)
						}
					}
				case FUNCTYPE_CPU, FUNCTYPE_GPU:
					if FUNCTYPE_FPGA == srcFunctionKind {
						logger.Info("nothing to do")
						crData.Status.To.Status = STATUS_OK
					} else if FUNCTYPE_CPU == srcFunctionKind || FUNCTYPE_GPU == srcFunctionKind {
						logger.Info("nothing to do")
						crData.Status.To.Status = STATUS_OK
					} else {
						logger.Info("nothing to do")
					}
				default:
					logger.Info("nothing to do")
					crData.Status.To.Status = STATUS_OK
				}
			}
		}

		// PCIeConnection is created
		if crData.Status.From.Status == STATUS_OK &&
			crData.Status.To.Status == STATUS_OK {
			crStatus = RUNNING
		}
		r.UpdCustomResource(ctx, r, &crData, crStatus)

		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Create", "Create End")

		if crStatus != RUNNING {
			// Requeue
			logger.Error(nil, "CR Status is not Running")
			return ctrl.Result{Requeue: true}, nil
		}

	} else if eventKind == UPDATE {
		// In case of update
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Update", "Update Start")
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Update", "Update End")
	} else if eventKind == DELETE {
		// In case of deletion
		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Delete", "Delete Start")

		// Delete the Finalizer statement.
		err = r.DelCustomResource(ctx, r, &crData)
		if err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
			"Delete", "Delete End")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PCIeConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.PCIeConnection{}).
		Complete(r)
}

func PCIeConnectionAccIDToDeviceID(accId string) int32 {
	var deviceID int32
	deviceID = -1
	for idx, devPath := range FpgaDevList {
		if accId == devPath {
			deviceID = int32(idx)
			break
		}
	}
	return deviceID
}

func PCIeConnectionD2D(ctx context.Context,
	pCRData *examplecomv1.PCIeConnection,
	srcDeviceID int32, pFunctionSrcData *examplecomv1.FunctionStatusData,
	dstDeviceID int32, pFunctionDstData *examplecomv1.FunctionStatusData) string {

	logger := log.FromContext(ctx)
	var status string
	status = STATUS_NG

	for i := 0; i < 1; i++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		if 0 > srcDeviceID || 0 > dstDeviceID {
			logger.Info("not deviceID setting from AccId NG srcDeviceID = " +
				strconv.Itoa(int(srcDeviceID)) + "dstDeviceID = " +
				strconv.Itoa(int(dstDeviceID)))
			break
		}

		if done, exist := ShmemEnable[pFunctionSrcData.SharedMemory.FilePrefix]; !exist || (exist && !done) {
			// Call the prefix enable function
			var size []C.uint
			size = append(size, C.uint(pFunctionSrcData.SharedMemory.SharedMemoryMiB)) //nolint:staticcheck // 'size' is never used. FIXME: remove variable
			ret := C.fpga_shmem_enable(
				C.CString(pFunctionSrcData.SharedMemory.FilePrefix), nil)
			if ret != 0 {
				ShmemEnable[pFunctionSrcData.SharedMemory.FilePrefix] = false
				logger.Info("fpga_shmem_enable() NG ret = " +
					strconv.Itoa(int(ret)))
				break
			}
			ShmemEnable[pFunctionSrcData.SharedMemory.FilePrefix] = true
			logger.Info("fpga_shmem_enable() OK ret = " +
				strconv.Itoa(int(ret)))
		}

		if done, exist := ShmemInit[pFunctionSrcData.SharedMemory.FilePrefix]; !exist || (exist && !done) {
			// Call the initialization function for shared memory
			dpdkLogFlag, _ := strconv.Atoi(os.Getenv("K8S_DPDK_LOG_FLAG"))
			ret := C.fpga_shmem_init(
				C.CString(pFunctionSrcData.SharedMemory.FilePrefix), nil,
				C.int(dpdkLogFlag))
			if ret < 0 {
				ShmemInit[pFunctionSrcData.SharedMemory.FilePrefix] = false
				logger.Info("fpga_shmem_init() NG ret = " +
					strconv.Itoa(int(ret)))
				break
			}
			ShmemInit[pFunctionSrcData.SharedMemory.FilePrefix] = true
			logger.Info("fpga_shmem_init() OK ret = " +
				strconv.Itoa(int(ret)))
		}

		// Call the shared memory area allocation function
		bufaddr := C.fpga_shmem_aligned_alloc(
			C.ulong(pFunctionSrcData.SharedMemory.SharedMemoryMiB * Hugepagesz))
		if nil == bufaddr {
			logger.Info("fpga_shmem_aligned_alloc() NG")
			break
		}
		logger.Info("fpga_shmem_aligned_alloc() OK")

		// Connection ID registration function call
		ret := C.fpga_chain_connect_egress(C.uint(srcDeviceID),
			/** if 0 (FPGAlibrary update) **
			C.uint(*pFunctionSrcData.FrameworkKernelID),
			C.uint(*pFunctionSrcData.FunctionChannelID),
			C.uint(*pFunctionSrcData.Tx.FDMAConnectorID))
			*** else if **/
			C.uint(*pFunctionSrcData.FrameworkKernelID),
			C.uint(*pFunctionSrcData.FunctionChannelID),
			C.uint(0),
			C.uint(*pFunctionSrcData.Tx.LLDMAConnectorID),
			C.uint8_t(1),
			C.uint8_t(0),
			C.uint8_t(1))
		/** end if  **/
		// Execution result of the connection ID registration function
		if ret != 0 {
			logger.Info("fpga_chain_connect_egress() NG ret = " +
				strconv.Itoa(int(ret)))
			break
		}
		logger.Info("fpga_chain_connect_egress() OK ret = " +
			strconv.Itoa(int(ret)))

		// Connection ID registration function call
		ret = C.fpga_chain_connect_ingress(C.uint(dstDeviceID),
			/** if 0 (FPGAlibrary update) **
			C.uint(*pFunctionSrcData.FrameworkKernelID),
			C.uint(*pFunctionDstData.FrameworkKernelID),
			C.uint(*pFunctionDstData.FunctionChannelID),
			C.uint(*pFunctionDstData.Rx.FDMAConnectorID))
			*** else if **/
			C.uint(*pFunctionDstData.FrameworkKernelID),
			C.uint(*pFunctionDstData.FunctionChannelID),
			C.uint(0),
			C.uint(*pFunctionDstData.Rx.LLDMAConnectorID),
			C.uint8_t(1),
			C.uint8_t(0))
		/** end if  **/
		// Execution result of the connection ID registration function
		if ret != 0 {
			logger.Info("fpga_chain_connect_ingress() NG ret = " +
				strconv.Itoa(int(ret)))
			break
		}
		logger.Info("fpga_chain_connect_ingress() OK ret = " +
			strconv.Itoa(int(ret)))

		// Call the buffer connection function
		bufferConnection := C.fpga_lldma_connect_t{}
		bufferConnection.tx_dev_id = C.uint(srcDeviceID)
		bufferConnection.tx_chid = C.uint(*pFunctionSrcData.Tx.DMAChannelID)
		bufferConnection.rx_dev_id = C.uint(dstDeviceID)
		bufferConnection.rx_chid = C.uint(*pFunctionDstData.Rx.DMAChannelID)
		bufferConnection.buf_size = C.uint(pFunctionSrcData.SharedMemory.SharedMemoryMiB * Hugepagesz)
		bufferConnection.buf_addr = bufaddr
		bufferConnection.connector_id = C.CString(
			pFunctionSrcData.SharedMemory.CommandQueueID)
		ret = C.fpga_lldma_buf_connect(&bufferConnection) //nolint:gocritic // suspicious indentical LHS and RHS for next line "!=". QUESTION: why?
		if ret != 0 {
			logger.Info("fpga_buf_connect() NG ret = " +
				strconv.Itoa(int(ret)))
			// Call the buffer release function
			C.fpga_shmem_free(bufaddr)
			break
		}
		logger.Info("fpga_buf_connect() OK ret = " +
			strconv.Itoa(int(ret)))

		status = STATUS_OK
		break //nolint:staticcheck // SA4004: Intentional break
	}
	return status
}

func PCIeConnectionSrcFPGA(ctx context.Context,
	deviceID int32, pFunctionSrcData *examplecomv1.FunctionStatusData,
	pFunctionDstData *examplecomv1.FunctionStatusData) string {
	logger := log.FromContext(ctx)
	var status string
	status = STATUS_NG

	if 0 <= deviceID {
		// Connection ID registration function call
		ret := C.fpga_chain_connect_egress(C.uint(deviceID),
			/** if 0 (FPGAlibrary update) **
			C.uint(*pFunctionSrcData.FrameworkKernelID),
			C.uint(*pFunctionSrcData.FunctionChannelID),
			C.uint(*pFunctionSrcData.Tx.FDMAConnectorID))
			*** else if **/
			C.uint(*pFunctionSrcData.FrameworkKernelID),
			C.uint(*pFunctionSrcData.FunctionChannelID),
			C.uint(0),
			C.uint(*pFunctionSrcData.Tx.LLDMAConnectorID),
			C.uint8_t(1),
			C.uint8_t(0),
			C.uint8_t(1))
		/** end if */
		// Execution result of the connection ID registration function
		if ret != 0 {
			logger.Info("fpga_chain_connect_egress() NG ret = " +
				strconv.Itoa(int(ret)))
		} else {
			logger.Info("fpga_chain_connect_egress() OK ret = " +
				strconv.Itoa(int(ret)))

			// Call FDMA setting function
			connectionID := C.CString(pFunctionDstData.SharedMemory.CommandQueueID)
			var dmaInfo C.dma_info_t
			ret = C.fpga_lldma_init(C.uint(deviceID),
				C.dma_dir_t(C.DMA_DEV_TO_HOST),
				C.uint(*pFunctionSrcData.Tx.DMAChannelID),
				connectionID,
				&dmaInfo) //nolint:gocritic // suspicious indentical LHS and RHS for next block "==". QUESTION: why?
			// Execution result of FDMA setting function
			if ret == 0 {
				status = STATUS_OK
				logger.Info("fpga_lldma_init() OK ret = " +
					strconv.Itoa(int(ret)))
			} else if ret == (-C.ALREADY_ACTIVE_CHID) {
				status = STATUS_OK
				logger.Info("fpga_lldma_init() OK ret = " +
					strconv.Itoa(int(ret)))
			} else {
				logger.Info("fpga_lldma_init() NG ret = " +
					strconv.Itoa(int(ret)))
			}
		}
	} else {
		logger.Info("not deviceID setting from AccId NG deviceID = " +
			strconv.Itoa(int(deviceID)))
	}

	return status
}

func PCIeConnectionDstFPGA(ctx context.Context,
	deviceID int32, pFunctionSrcData *examplecomv1.FunctionStatusData,
	pFunctionDstData *examplecomv1.FunctionStatusData) string {
	logger := log.FromContext(ctx)
	var status string
	status = STATUS_NG

	if 0 <= deviceID {
		// Connection ID registration function call
		ret := C.fpga_chain_connect_ingress(C.uint(deviceID),
			/** if 0 (FPGAlibrary update) **
			C.uint(*pFunctionDstData.FrameworkKernelID),
			C.uint(*pFunctionDstData.FunctionChannelID),
			C.uint(*pFunctionDstData.Rx.FDMAConnectorID))
			*** else if **/
			C.uint(*pFunctionDstData.FrameworkKernelID),
			C.uint(*pFunctionDstData.FunctionChannelID),
			C.uint(0),
			C.uint(*pFunctionDstData.Rx.LLDMAConnectorID),
			C.uint8_t(1),
			C.uint8_t(0))
		/** end if  **/
		// Execution result of the connection ID registration function
		if ret != 0 {
			logger.Info("fpga_chain_connect_ingress() NG ret = " +
				strconv.Itoa(int(ret)))
		} else {
			logger.Info("fpga_chain_connect_ingress() OK ret = " +
				strconv.Itoa(int(ret)))

			// Call FDMA setting function
			connectionID := C.CString(pFunctionSrcData.SharedMemory.CommandQueueID)
			var dmaInfo C.dma_info_t
			ret = C.fpga_lldma_init(C.uint(deviceID),
				C.dma_dir_t(C.DMA_HOST_TO_DEV),
				C.uint(*pFunctionDstData.Rx.DMAChannelID), // @TODO
				connectionID,
				&dmaInfo) //nolint:gocritic // suspicious indentical LHS and RHS for next block "==". QUESTION: why?
			// Execution result of FDMA setting function
			if ret == 0 {
				status = STATUS_OK
				logger.Info("fpga_lldma_init() OK ret = " +
					strconv.Itoa(int(ret)))
			} else if ret == (-C.ALREADY_ACTIVE_CHID) {
				status = STATUS_OK
				logger.Info("fpga_lldma_init() OK ret = " +
					strconv.Itoa(int(ret)))
			} else {
				logger.Info("fpga_ldma_init() NG ret = " +
					strconv.Itoa(int(ret)))
			}
		}
	} else {
		logger.Info("not deviceID setting from AccId NG deviceID = " +
			strconv.Itoa(int(deviceID)))
	}

	return status
}

func (r *PCIeConnectionReconciler) GetfinalizerName(pCRData *examplecomv1.PCIeConnection) string {
	// Value to set in the finalizer
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
}

func (r *PCIeConnectionReconciler) GetEventKind(pCRData *examplecomv1.PCIeConnection) int32 {
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

func (r *PCIeConnectionReconciler) UpdCustomResource(ctx context.Context, pClient *PCIeConnectionReconciler,
	pCRData *examplecomv1.PCIeConnection, status string) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	pCRData.Status.StartTime = metav1.Now()

	if status == RUNNING {
		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, r.GetfinalizerName(pCRData))
		// status update
		pCRData.Status.DataFlowRef = pCRData.Spec.DataFlowRef
		pCRData.Status.From.WBFunctionRef = pCRData.Spec.From.WBFunctionRef
		pCRData.Status.To.WBFunctionRef = pCRData.Spec.To.WBFunctionRef
		pCRData.Status.Status = status
	}
	err = pClient.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "Status Update Error.")
	} else {
		logger.Info("Status Update.")
	}
	return err
}

func (r *PCIeConnectionReconciler) DelCustomResource(ctx context.Context, pClient *PCIeConnectionReconciler,
	pCRData *examplecomv1.PCIeConnection) error {
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

// Function Data Get
func (r *PCIeConnectionReconciler) GetFunctionData(
	ctx context.Context,
	functionCR examplecomv1.WBNamespacedName,
	functionCRData *Function,
	functionKind *string) error {
	logger := log.FromContext(ctx)

	var err error

	var existCount int = 0
	var notFoundErr error
	var elseErr error
	var strmapFunctionCRData map[string]interface{}

	var functionKindList []string = []string{FUNCTYPE_FPGA,
		FUNCTYPE_GPU,
		FUNCTYPE_GATE,
		FUNCTYPE_CPU}

	fcrData := &unstructured.Unstructured{}

	for n := 0; n < len(functionKindList); n++ {
		fcrData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    functionKindList[n],
		})
		err = r.Get(ctx, client.ObjectKey{
			Namespace: functionCR.Namespace,
			Name:      functionCR.Name}, fcrData)
		if errors.IsNotFound(err) {
			existCount += 1
			notFoundErr = err
		} else if err != nil {
			if FUNCTYPE_GATE != functionKindList[n] {
				logger.Info("unable to fetch FunctionCR")
				elseErr = err
			} else {
				existCount += 1
			}
		} else {
			*functionKind = functionKindList[n]
			break
		}
		if existCount == len(functionKindList) {
			logger.Info("FunctionCR does not exist")
			err = notFoundErr
		} else {
			err = elseErr
		}
	}

	if err != nil {
		return err
	}

	if len(fcrData.Object) != 0 {
		// Store spec information
		strmapFunctionCRData, _, _ = unstructured.NestedMap(fcrData.Object, "spec")

		// Convert the obtained mapdata to byte type
		bytes, err := json.Marshal(strmapFunctionCRData)
		if err != nil {
			logger.Error(err, "unable to json.marshal")
			return err
		}
		// Replace with a struct
		err = json.Unmarshal(bytes, &functionCRData.Spec)
		if err != nil {
			logger.Error(err, "unable to json.unmarshal")
			return err
		}

		// Store the status information
		strmapFunctionCRData, _, _ = unstructured.NestedMap(fcrData.Object, "status")

		// Convert the obtained mapdata to byte type
		bytes, err = json.Marshal(strmapFunctionCRData)
		if err != nil {
			logger.Error(err, "unable to json.marshal")
			return err
		}
		// Replace with a struct
		err = json.Unmarshal(bytes, &functionCRData.Status)
		if err != nil {
			logger.Error(err, "unable to json.unmarshal")
			return err
		}
	}
	return nil
}

func PCIeConnectionFPGAInit(mng ctrl.Manager) {
	ctx := context.Background()
	logger := log.FromContext(ctx)

	myNodeName := os.Getenv("K8S_NODENAME")

	fpgaList := GetFPGACR(ctx, mng, myNodeName)

	if 0 != len(fpgaList) {
		// FPGAPhase3 initial setting
		var argv []*C.char
		// Phase 3 FPGA initialization variables
		argv = []*C.char{C.CString("proc"),
			C.CString("-d"),
			C.CString(strings.Join(fpgaList, ","))}
		argc := C.int(len(argv))
		/*#if 1 * IT ph-2 temporary solution (fpga log level change/standard output) ***/
		C.libfpga_log_set_output_stdout()
		C.libfpga_log_set_level(C.LIBFPGA_LOG_ALL)
		/*#else * IT ph-2 temporary solution (change fpga log level/standard output)
		**** **#endif* IT ph-2 temporary solution (change fpga log level/standard output) ***/
		C.fpga_init(argc, (**C.char)(unsafe.Pointer(&argv[0])))

		// Hold device information
		// FpgaDevList = []string{}
		for _, devPath := range fpgaList {
			FpgaDevList = append(FpgaDevList, devPath)
		}
	}
	// Start the shared memory controller
	port := 60000
	addr := "localhost"
	ret := C.fpga_shmem_controller_init(C.ushort(port),
		C.CString(addr))
	logger.Info("fpga_shmem_controller_init() ret = " +
		strconv.Itoa(int(ret)))
	/* Since the process stops when the controller is stopped, do not call exit for each init.
	   defer C.fpga_finish() // Close FPGA (deferred execution)
	*/
}

// Get FPGACR
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
