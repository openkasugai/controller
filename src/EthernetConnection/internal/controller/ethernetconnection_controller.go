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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "EthernetConnection/api/v1"

	/* Additional files */
	"encoding/json"
	"fmt"
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
	// #cgo CFLAGS: -mrtm
	// #cgo CFLAGS:  -I/usr/local/include/fpgalib/
	// #cgo CFLAGS: -I/usr/include/
	// #cgo CFLAGS:  -I/usr/local/include/fpgalib/dpdk/include/
	// #cgo LDFLAGS: -L/usr/local/lib/fpgalib/ -lfpga
	// #cgo LDFLAGS: -L/usr/lib/gcc/x86_64-linux-gnu/9/ -lstdc++
	// #cgo CXXFLAGS: -std=c++11
	// #cgo LDFLAGS: -L. -lstdc++
	// #cgo LDFLAGS: -L. -lpciaccess
	// #include <liblldma.h>
	// #include <libptu.h>
	// #include <libfpgactl.h>
	// #include <libchain.h>
	// #include <libdmacommon.h>
	// #include <liblogging.h>
	// #include <libshmem.h>
	// #include <libshmem_controller.h>
	// #include <arpa/inet.h>
	"C"
)

// EthernetConnectionReconciler reconciles a EthernetConnection object
type EthernetConnectionReconciler struct {
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

// DeviceID type
const (
	SRC_OR_DST = iota
	SRC_AND_DST
)

// Status type
const (
	STATUS_OK   = "OK"
	STATUS_INIT = "INIT"
	STATUS_NG   = "NG"
)

// Overall Status type
const (
	PENDING = "Pending"
	RUNNING = "Running"
)

// Load FPGA file
const (
	CONFIG_FILE_OK = iota
	CONFIG_FILE_NG
)

// FunctionType type
const (
	FUNCTYPE_FPGA = "FPGAFunction"
	FUNCTYPE_GPU  = "GPUFunction"
	FUNCTYPE_GATE = "GateFunction"
	FUNCTYPE_CPU  = "CPUFunction"
)

// Connection method
const (
	PROT_RTP = "RTP"
	PROT_TCP = "TCP"
	PROT_DMA = "DMA"
)

// Timeout value
type timeout_s struct {
	tv_sec  int32
	tv_usec int32
}

// Structure for obtaining Function information
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   examplecomv1.FunctionData       `json:"spec,omitempty"`
	Status examplecomv1.FunctionStatusData `json:"status,omitempty"`
}

var Timeout = timeout_s{3, 0}
var FpgaDevList = []string{} // Supports multiple FPGAs

//+kubebuilder:rbac:groups=example.com,resources=childbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=childbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=childbs/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com.example.com,resources=ethernetconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com.example.com,resources=ethernetconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com.example.com,resources=ethernetconnections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EthernetConnection object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *EthernetConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var crData examplecomv1.EthernetConnection
	var eventKind int32 // 0:Add, 1:Upd,  2:Del
	var srcStatus string
	var dstStatus string
	var srcDeviceID int32
	var dstDeviceID int32
	var crStatus string
	var srcFunctionData Function
	var dstFunctionData Function
	var srcFunctionKind string
	var dstFunctionKind string

	crStatus = PENDING
	srcStatus = STATUS_OK
	dstStatus = STATUS_OK
	srcDeviceID = -1
	dstDeviceID = -1
	srcFunctionKind = ""
	dstFunctionKind = ""

	// Collect information about your own node
	myNodeName := os.Getenv("K8S_NODENAME")

	if true == strings.Contains(req.NamespacedName.Name, "-wbconnection-") {

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
		}

		if eventKind == CREATE {
			// In case of creation
			r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
				"Create", "Create Start")

			// Check whether to process on the current node
			if (myNodeName == srcFunctionData.Spec.NodeName) ||
				(myNodeName == dstFunctionData.Spec.NodeName) {
				var srcAccID string
				var dstAccID string
				if FUNCTYPE_FPGA == srcFunctionKind {
					// Get deviceID from AccId
					for count := 0; count < len(srcFunctionData.Spec.AcceleratorIDs); count++ {
						if *srcFunctionData.Spec.AcceleratorIDs[count].PartitionName ==
							strconv.Itoa(int(*srcFunctionData.Status.FunctionKernelID)) {
							srcAccID = srcFunctionData.Spec.AcceleratorIDs[count].ID
							break
						}
					}
					srcDeviceID = EthernetConnection_AccIdToDevId(
						srcAccID)
				}
				if FUNCTYPE_FPGA == dstFunctionKind {
					// Get deviceID from AccId
					for count := 0; count < len(dstFunctionData.Spec.AcceleratorIDs); count++ {
						if *dstFunctionData.Spec.AcceleratorIDs[count].PartitionName ==
							strconv.Itoa(int(*dstFunctionData.Status.FunctionKernelID)) {
							dstAccID = dstFunctionData.Spec.AcceleratorIDs[count].ID
							break
						}
					}
					dstDeviceID = EthernetConnection_AccIdToDevId(
						dstAccID)
				}

				// Is the output side deployment node the current node?
				if dstFunctionData.Spec.NodeName == myNodeName &&
					crData.Status.To.Status != STATUS_OK {
					// Get the Status of the output FunctionType
					if dstFunctionData.Status.Status == RUNNING {
						switch dstFunctionKind {
						case FUNCTYPE_FPGA:
							if PROT_TCP == dstFunctionData.Status.Rx.Protocol {
								dstStatus = EthernetConnectionFPGAListen(ctx,
									dstDeviceID,
									&dstFunctionData.Status)
							} else if PROT_RTP == dstFunctionData.Status.Rx.Protocol {
								// No processing
							} else {
								dstStatus = STATUS_NG // Status update
								logger.Info("dstFunctionData.Rx.Protocol NG")
							}
						case FUNCTYPE_CPU, FUNCTYPE_GPU:
							dstStatus = r.handleCreateExternalNetworkFunc(ctx, &crData, &crData.Spec.To,
								&dstFunctionData, &srcFunctionData, dstFunctionKind, srcFunctionKind)
						default:
							logger.Info(dstFunctionKind + " nothing to do")
						}
					} else {
						dstStatus = STATUS_NG
						logger.Info("dstStatusCheck NG")
					}
				}

				// The input side is deployed to the current Node.
				if srcFunctionData.Spec.NodeName == myNodeName &&
					crData.Status.From.Status != STATUS_OK {
					if srcFunctionData.Status.Status == RUNNING {
						switch srcFunctionKind {
						case FUNCTYPE_FPGA:
							if PROT_TCP == srcFunctionData.Status.Tx.Protocol {
								srcStatus = EthernetConnectionFPGAConnect(
									ctx,
									srcDeviceID,
									&srcFunctionData.Status,
									&dstFunctionData.Status)
							} else if PROT_RTP == srcFunctionData.Status.Tx.Protocol {
								// No processing
							} else {
								srcStatus = STATUS_NG // Status update
								logger.Info("srcFunctionData.Tx.Protocol NG")
							}
						case FUNCTYPE_CPU, FUNCTYPE_GPU:
							srcStatus = r.handleCreateExternalNetworkFunc(ctx, &crData,
								&crData.Spec.From, &srcFunctionData, &dstFunctionData, srcFunctionKind, dstFunctionKind)
						default:
							logger.Info(srcFunctionKind + " nothing to do")
						}
					} else {
						srcStatus = STATUS_NG // Status update
						logger.Info("srcStatusCheck NG")
					}
				}
				if srcFunctionData.Spec.NodeName != myNodeName &&
					srcFunctionKind == FUNCTYPE_GATE {
					if srcFunctionData.Status.Status != RUNNING {
						srcStatus = STATUS_NG // Status update
						logger.Info("gatefunction srcStatusCheck NG")
					}
				}

				// Is the output side deployment node the current node?
				if dstFunctionData.Spec.NodeName == myNodeName &&
					crData.Status.To.Status != STATUS_OK {
					if dstStatus == STATUS_OK {
						switch dstFunctionKind {
						case FUNCTYPE_FPGA:
							if PROT_TCP == dstFunctionData.Status.Rx.Protocol {
								dstStatus = EthernetConnectionFPGAAccept(ctx,
									dstDeviceID,
									&srcFunctionData.Status,
									&dstFunctionData.Status)
								if dstStatus == STATUS_OK {
									if dstFunctionKind == FUNCTYPE_FPGA {
										logger.Info(
											"EthernetConnectionFPGAAccept() OK" +
												" chenge srcStatus = STATUS_OK")
										srcStatus = STATUS_OK
									}
								}
							} else if PROT_RTP == dstFunctionData.Status.Rx.Protocol {
								// No processing
							} else {
								dstStatus = STATUS_NG // Status update
								logger.Info("dstFunctionData.Rx.Protocol NG")
							}
						default:
							logger.Info(dstFunctionKind + " nothing to do")
						}
					}
				}
				if dstFunctionData.Spec.NodeName != myNodeName &&
					dstFunctionKind == FUNCTYPE_GATE {
					if dstFunctionData.Status.Status != RUNNING {
						dstStatus = STATUS_NG // Status update
						logger.Info("gatefunction dstStatusCheck NG")
					}
				}

				logger.Info("myNodeName " + myNodeName)
				logger.Info("src_node " + srcFunctionData.Spec.NodeName)
				logger.Info("srcStatus " + srcStatus)
				logger.Info("dst_node " + dstFunctionData.Spec.NodeName)
				logger.Info("dstStatus " + dstStatus)

				// Status description
				if srcFunctionData.Spec.NodeName == myNodeName {
					crData.Status.From.WBFunctionRef = crData.Spec.From.WBFunctionRef
					crData.Status.From.Status = srcStatus
				} else if srcFunctionKind == FUNCTYPE_GATE {
					crData.Status.From.WBFunctionRef = crData.Spec.From.WBFunctionRef
					crData.Status.From.Status = srcStatus
				}
				if dstFunctionData.Spec.NodeName == myNodeName {
					crData.Status.To.WBFunctionRef = crData.Spec.To.WBFunctionRef
					crData.Status.To.Status = dstStatus
				} else if dstFunctionKind == FUNCTYPE_GATE {
					crData.Status.To.WBFunctionRef = crData.Spec.To.WBFunctionRef
					crData.Status.To.Status = dstStatus
				}

				var srcLLDMAret bool
				var dstLLDMAret bool

				srcLLDMAret = true
				dstLLDMAret = true

				if srcFunctionKind == FUNCTYPE_FPGA {
					srcLLDMAret = lldmaInit(ctx, srcDeviceID, &srcFunctionData.Status)
				}

				if dstFunctionKind == FUNCTYPE_FPGA {
					dstLLDMAret = lldmaInit(ctx, dstDeviceID, &dstFunctionData.Status)
				}

				if crData.Status.From.Status == STATUS_OK &&
					crData.Status.To.Status == STATUS_OK &&
					srcLLDMAret == true && dstLLDMAret == true {
					crStatus = RUNNING
				}
				r.UpdCustomResource(ctx, &crData, crStatus)

				r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
					"Create", "Create End")

				if crStatus != RUNNING {
					// Requeue
					logger.Error(nil, "CR Status is not Running.")
					return ctrl.Result{Requeue: true}, nil
				}
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
			err = r.DelCustomResource(ctx, &crData)
			if err != nil {
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			r.Recorder.Eventf(&crData, corev1.EventTypeNormal,
				"Delete", "Delete End")
		}
	} else {
		_ = r.ChildBSCR(ctx, req, myNodeName)
		/*
			if nil != err {
				requeueFlag = true
			}
		*/
	}
	return ctrl.Result{}, nil
}

func lldmaInit(ctx context.Context,
	deviceID int32,
	pFunctionData *examplecomv1.FunctionStatusData) bool {
	logger := log.FromContext(ctx)
	var result bool
	result = true

	if pFunctionData.Rx.Protocol == "DMA" &&
		(pFunctionData.Tx.Protocol == "TCP" || pFunctionData.Tx.Protocol == "RTP") {
		// Call FDMA setting function
		connectionID := C.CString(pFunctionData.SharedMemory.CommandQueueID)
		var dmaInfo C.dma_info_t
		ret := C.fpga_lldma_init(C.uint(deviceID),
			C.dma_dir_t(C.DMA_NW_TO_DEV),
			C.uint(*pFunctionData.Rx.DMAChannelID),
			connectionID,
			&dmaInfo) //nolint:gocritic //suspicious identical LHS and RHS for `==` operator
		if 0 != ret {
			logger.Info("fpga_lldma_init() DMA_NW_TO_DEV NG ret = " +
				strconv.Itoa(int(ret)))
			result = false
		}
	}
	if pFunctionData.Tx.Protocol == "DMA" &&
		(pFunctionData.Rx.Protocol == "TCP" || pFunctionData.Rx.Protocol == "RTP") {
		// Call FDMA setting function
		connectionID := C.CString(pFunctionData.SharedMemory.CommandQueueID)
		var dmaInfo C.dma_info_t
		ret := C.fpga_lldma_init(C.uint(deviceID),
			C.dma_dir_t(C.DMA_DEV_TO_NW),
			C.uint(*pFunctionData.Tx.DMAChannelID),
			connectionID,
			&dmaInfo) //nolint:gocritic //suspicious identical LHS and RHS for `==` operator
		if 0 != ret {
			logger.Info("fpga_lldma_init() DMA_DEV_TO_NW NG ret = " +
				strconv.Itoa(int(ret)))
			result = false
		}
	}
	return result

}

// SetupWithManager sets up the controller with the Manager.
func (r *EthernetConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.EthernetConnection{}).
		Watches(&examplecomv1.ChildBs{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

// ChildBitstream CR
func (r *EthernetConnectionReconciler) ChildBSCR(ctx context.Context,
	req ctrl.Request, myNodeName string) error {
	logger := log.FromContext(ctx)

	var err error
	var crChildBs examplecomv1.ChildBs
	var crFPGA examplecomv1.FPGA

	var result int32

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		result = -1
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

		if crFPGA.Status.NodeName != myNodeName {
			break
		}

		if crChildBs.Status.State != examplecomv1.ChildBsNoConfigureNetwork {
			break
		}

		crChildBs.Status.State = examplecomv1.ChildBsConfiguringNetwork
		err = r.Update(ctx, &crChildBs)
		if nil != err {
			logger.Error(err, "ChildBsCR Update Error")
			break
		} else {
			logger.Info("ChildBsCR Update Success")
		}

		for i := 0; i < len(crChildBs.Spec.Regions); i++ {
			/*
				ptuParams := crChildBs.Spec.Regions[i].Modules.Ptu.Parameters

				//			var IPAddress string
				//			var SubnetAddress string
				//			var GatewayAddress string

				var MACAddress string

				if "" != (*ptuParams)["IPAddress"].StrVal &&
					"" != (*ptuParams)["SubnetAddress"].StrVal &&
					"" != (*ptuParams)["GatewayAddress"].StrVal &&
					"" != (*ptuParams)["MacAddress"].StrVal {
					//				IPAddress = (*ptuParams)["IPAddress"].StrVal
					//				SubnetAddress = (*ptuParams)["SubnetAddress"].StrVal
					//				GatewayAddress = (*ptuParams)["GatewayAddress"].StrVal
					MACAddress = (*ptuParams)["MacAddress"].StrVal
				} else {
					break
				}

				//			Name := *crChildBs.Spec.Regions[i].Name

				//			DeviceIndex := EthernetConnection_AccIdToDevId(crFPGA.Spec.DeviceFilePath)
				//			LaneIndex, _ := strconv.Atoi(Name[len(Name)-1:])


				var macDataPh3 []uint8
				// MAC address conversion process
				strMACData := strings.Split(MACAddress, ":")
				for _, data := range strMACData {
					dataHex, _ := strconv.ParseUint(data, 16, 0)
					macDataPh3 = append(macDataPh3, uint8(dataHex))
				}
			*/
			//			macAddressPh3 := (*C.uchar)(unsafe.Pointer(&macDataPh3[0]))

			// Use Phase3 functions
			//			ret := C.fpga_ptu_init(
			//				C.uint(DeviceIndex),
			//				C.uint(int32(LaneIndex)),
			//				strconvInetAddressToUint32Swap(IPAddress),
			//				strconvInetAddressToUint32Swap(SubnetAddress),
			//				strconvInetAddressToUint32Swap(GatewayAddress),
			//				macAddressPh3)
			ret := 0
			logger.Info("fpga_ptu_init() ret = " + strconv.Itoa(int(ret)))
			if ret != 0 {
				result = 0
				break
			} else {
				result = 1
			}
		}

		if nil == err {
			if 1 == result {
				crChildBs.Status.State = examplecomv1.ChildBsReady
				err = r.Update(ctx, &crChildBs)
				if nil != err {
					logger.Error(err, "ChildBsCR Update Error")
					break
				} else {
					logger.Info("ChildBsCR Update Success")
					break
				}
			} else if 0 == result {
				crChildBs.Status.State = examplecomv1.ChildBsError
				err = r.Update(ctx, &crChildBs)
				if nil != err {
					logger.Error(err, "ChildBsCR Update Error")
					break
				} else {
					logger.Info("ChildBsCR Update Success")
					break
				}
			}
		}
	}
	return err
}

func strconvInetAddressToUint32Swap(ipAddress string) C.uint32_t {
	var inetAddress C.in_addr_t

	// Store in inetAddress_t type in network byte order.
	inetAddress = C.inet_addr(C.CString(ipAddress))

	// Due to FPGA specifications, the network byte order is reversed.
	return (((inetAddress & 0xff000000) >> 24) |
		((inetAddress & 0x00ff0000) >> 8) |
		((inetAddress & 0x0000ff00) << 8) |
		((inetAddress & 0x000000ff) << 24))
}

func EthernetConnectionFPGAListen(ctx context.Context,
	deviceID int32,
	pFunctionDstData *examplecomv1.FunctionStatusData) string {
	logger := log.FromContext(ctx)
	var status string
	status = STATUS_NG

	// Use Phase3 functions
	if 0 <= deviceID {
		// Call the listen function
		ret := C.fpga_ptu_listen(C.uint(deviceID),
			/*#if 0 * IT ph-2 temporary workaround (adjust fpga_ptu_listen argument) ****
			  C.uint(pCRData.Spec.To.FuncKernelId),
			**#else * IT ph-2 temporary workaround (adjust fpga_ptu_listen argument) ***/
			/*#endif* IT ph-2 temporary solution (fpga_ptu_listen argument adjustment) ***/
			C.uint(*pFunctionDstData.PtuKernelID),
			C.ushort(*pFunctionDstData.Rx.Port))
		// Execution result of the listen function
		if ret == 0 {
			status = STATUS_OK // Status update
			logger.Info("fpga_ptu_listen() OK ret = " +
				strconv.Itoa(int(ret)))
		} else if ret == -(1) {
			status = STATUS_OK // Status update
			logger.Info("fpga_ptu_listen() NG ret = " +
				strconv.Itoa(int(ret)) + " but temporary OK")
		} else {
			logger.Info("fpga_ptu_listen() NG ret = " +
				strconv.Itoa(int(ret)))
		}
	} else {
		logger.Info(
			"EthernetConnectionFPGAListen() not deviceID" +
				" Setting deviceID = " + strconv.Itoa(int(deviceID)))
	}
	return status
}

func EthernetConnectionFPGAAccept(ctx context.Context,
	deviceID int32,
	pFunctionSrcData *examplecomv1.FunctionStatusData,
	pFunctionDstData *examplecomv1.FunctionStatusData) string {
	logger := log.FromContext(ctx)
	var status string
	status = STATUS_NG

	if 0 <= deviceID {
		// Use Phase3 functions
		// Call the accept function/connection ID registration function
		var timeout C.struct_timeval
		timeout.tv_sec = C.long(Timeout.tv_sec)
		timeout.tv_usec = C.long(Timeout.tv_usec)
		cid := C.uint(0)
		ret := C.fpga_ptu_accept(C.uint(deviceID),
			C.uint(*pFunctionDstData.PtuKernelID),
			C.in_port_t(*pFunctionDstData.Rx.Port),
			strconvInetAddressToUint32Swap(*pFunctionSrcData.Tx.IPAddress),
			C.in_port_t(*pFunctionSrcData.Tx.Port),
			&timeout,
			&cid)
		// Execution result of the accept function
		if ret == 0 {
			logger.Info("fpga_ptu_accept() OK ret = " +
				strconv.Itoa(int(ret)))
			// Connection ID registration function call
			ret := C.fpga_chain_connect_ingress(C.uint(deviceID),
				/** if 0 (FPGAlibrary update) **
				C.uint(*pFunctionDstData.FrameworkKernelID),
				C.uint(*pFunctionDstData.FunctionChannelID),
				C.uint(cid))
				*** else if **/
				C.uint(*pFunctionDstData.FrameworkKernelID),
				C.uint(*pFunctionDstData.FunctionChannelID),
				C.uint(1),
				C.uint(cid),
				C.uint8_t(1),
				C.uint8_t(0))
			/** end if  **/
			// Execution result of the connection ID registration function
			if ret == 0 {
				status = STATUS_OK // Status update
				logger.Info("fpga_chain_connect_ingress() OK ret = " +
					strconv.Itoa(int(ret)))
			} else {
				logger.Info("fpga_chain_connect_ingress() NG ret = " +
					strconv.Itoa(int(ret)))
			}
		} else {
			logger.Info("fpga_ptu_accept() NG ret = " +
				strconv.Itoa(int(ret)))
		}
	} else {
		logger.Info(
			"EthernetConnectionFPGAAccept() not deviceID" +
				" Setting deviceID = " + strconv.Itoa(int(deviceID)))
	}
	return status
}

func EthernetConnectionFPGAConnect(ctx context.Context,
	deviceID int32,
	pFunctionSrcData *examplecomv1.FunctionStatusData,
	pFunctionDstData *examplecomv1.FunctionStatusData) string {
	logger := log.FromContext(ctx)
	var status string
	status = STATUS_NG

	if 0 <= deviceID {
		// Use Phase3 functions
		var timeout C.struct_timeval
		timeout.tv_sec = C.long(Timeout.tv_sec)
		timeout.tv_usec = C.long(Timeout.tv_usec)
		cid := C.uint(0)
		ret := C.fpga_ptu_connect(C.uint(deviceID),
			C.uint(*pFunctionSrcData.PtuKernelID),
			C.in_port_t(*pFunctionSrcData.Tx.Port),
			strconvInetAddressToUint32Swap(*pFunctionDstData.Rx.IPAddress),
			C.in_port_t(*pFunctionDstData.Rx.Port),
			&timeout,
			&cid)
		if ret == 0 {
			logger.Info("fpga_ptu_connect() OK ret = " +
				strconv.Itoa(int(ret)))
			ret = C.fpga_chain_connect_egress(C.uint(deviceID),
				/** if 0 (FPGAlibrary update) **
					C.uint(*pFunctionSrcData.FrameworkKernelID),
					C.uint(*pFunctionSrcData.FunctionChannelID),
					C.uint(cid))
				*** else if **/
				C.uint(*pFunctionSrcData.FrameworkKernelID),
				C.uint(*pFunctionSrcData.FunctionChannelID),
				C.uint(1),
				C.uint(cid),
				C.uint8_t(1),
				C.uint8_t(0),
				C.uint8_t(1))
			/** end if */
			if ret == 0 {
				status = STATUS_OK // Status update
				logger.Info("fpga_chain_connect_egress() OK ret = " +
					strconv.Itoa(int(ret)))
			} else {
				logger.Info("fpga_chain_connect_egress() NG ret = " +
					strconv.Itoa(int(ret)))
			}
		} else if ret == -(3) {
			status = STATUS_OK // Status update
			logger.Info("fpga_ptu_connect() OK ret = " +
				strconv.Itoa(int(ret)))
		} else {
			logger.Info("fpga_ptu_connect() NG ret = " +
				strconv.Itoa(int(ret)))
		}
	} else {
		logger.Info(
			"EthernetConnectionFPGAConnect() not deviceID" +
				" Setting deviceID = " + strconv.Itoa(int(deviceID)))
	}
	return status
}

func EthernetConnection_AccIdToDevId(accId string) int32 {
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

func (r *EthernetConnectionReconciler) GetfinalizerName(pCRData *examplecomv1.EthernetConnection) string {
	// Value to set in the finalizer
	gvks, _, _ := r.Client.Scheme().ObjectKinds(pCRData)
	return strings.ToLower(gvks[0].Kind) + ".finalizers." +
		strings.ToLower(gvks[0].Group+"."+gvks[0].Version)
	/*
	   // Value to set in the finalizer
	   return strings.ToLower(pCRData.Kind) + ".finalizers." +

	   	strings.ReplaceAll(pCRData.APIVersion, "/", ".")
	*/
}

func (r *EthernetConnectionReconciler) GetEventKind(pCRData *examplecomv1.EthernetConnection) int32 {
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

func (r *EthernetConnectionReconciler) UpdCustomResource(ctx context.Context,
	pCRData *examplecomv1.EthernetConnection,
	status string) error {
	logger := log.FromContext(ctx)
	var err error
	err = nil

	pCRData.Status.StartTime = metav1.Now()

	if status == RUNNING {
		// Write a Finalizer
		controllerutil.AddFinalizer(pCRData, r.GetfinalizerName(pCRData))
		// status update
		pCRData.Status.Status = status
		pCRData.Status.DataFlowRef = pCRData.Spec.DataFlowRef
	}
	err = r.Update(ctx, pCRData)
	if err != nil {
		logger.Error(err, "Status Update Error.")
	} else {
		logger.Info("Status Update.")
	}
	return err
}

func (r *EthernetConnectionReconciler) DelCustomResource(ctx context.Context,
	pCRData *examplecomv1.EthernetConnection) error {
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

// Function Data Get
func (r *EthernetConnectionReconciler) GetFunctionData(
	ctx context.Context,
	findNamespacedName examplecomv1.WBNamespacedName,
	functionData *Function,
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
			Namespace: findNamespacedName.Namespace,
			Name:      findNamespacedName.Name}, fcrData)
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
		err = json.Unmarshal(bytes, &functionData.Spec)
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
		err = json.Unmarshal(bytes, &functionData.Status)
		if err != nil {
			logger.Error(err, "unable to json.unmarshal")
			return err
		}
	}
	return nil
}

// Phase1/Phase3FPGA setting function
func CRCStartPtuInitSet(mng ctrl.Manager) {
	ctx := context.Background()

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
		for deviceID, devPath := range fpgaList {
			C.fpga_enable_regrw(C.uint(deviceID))
			FpgaDevList = append(FpgaDevList, devPath)
		}
	}
	/* Since the process stops when the controller is stopped, do not call exit for each init.
	defer C.fpga_finish() // Close FPGA (deferred execution)
	C.fpga_ptu_exit(
		C.uint(0),
		C.uint(cr_data.Spec.DstFunc.PtuKernelId))
	C.fpga_ptu_exit_ph1(
		C.uint(0),
		C.uint(cr_data.Spec.SrcFunc.PtuKernelId)) // End event monitoring thread (deferred execution)
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
		//		if myNodeName != fpgaCRData.Items[i].Status.NodeName {
		if myNodeName != fpgaCRData.Items[i].Status.NodeName {
			continue
		}
		//		fpgaList = append(fpgaList, fpgaCRData.Items[i].Status.DeviceFilePath)
		fpgaList = append(fpgaList, fpgaCRData.Items[i].Status.DeviceFilePath)
	}
	return fpgaList
}

func (r *EthernetConnectionReconciler) handleCreateExternalNetworkFunc(
	ctx context.Context,
	ethConn *examplecomv1.EthernetConnection,
	fnSpec *examplecomv1.EthernetFunctionSpec,
	fnData *Function,
	othrSideFnData *Function,
	fnKind string,
	othrSideFnKind string) string {

	logger := log.FromContext(ctx)

	type directiontype string
	const (
		DirectionUnknown directiontype = "Unknown"
		DirectionFrom    directiontype = "From"
		DirectionTo      directiontype = "To"
	)

	direction := DirectionUnknown
	var othrSideFnSpec *examplecomv1.EthernetFunctionSpec
	switch {
	case &ethConn.Spec.From == fnSpec:
		othrSideFnSpec = &ethConn.Spec.To
		direction = DirectionFrom
	case &ethConn.Spec.To == fnSpec:
		othrSideFnSpec = &ethConn.Spec.From
		direction = DirectionTo
	}

	switch direction {
	case DirectionTo:

		wbf := &unstructured.Unstructured{}
		wbf.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "WBFunction",
		})
		err := r.Get(ctx, client.ObjectKey{
			Namespace: fnSpec.WBFunctionRef.Namespace,
			Name:      fnSpec.WBFunctionRef.Name,
		}, wbf)
		if errors.IsNotFound(err) {
			logger.Error(err, "WBFunctionCR does not exist")
			return STATUS_NG
		}

		podName := ""
		switch fnKind {
		case FUNCTYPE_CPU:
			podName = fmt.Sprintf("%s-cpu-pod", fnSpec.WBFunctionRef.Name)
		case FUNCTYPE_GPU:
			if 0 < len(fnData.Spec.AcceleratorIDs) {
				podName = fmt.Sprintf("%s-mps-dgpu-%s-pod", fnSpec.WBFunctionRef.Name, strings.ToLower(fnData.Spec.AcceleratorIDs[0].ID))
			}
		default:
			return STATUS_OK
		}

		pod := corev1.Pod{}
		pod.Name = podName
		pod.Namespace = fnSpec.WBFunctionRef.Namespace
		err = r.Get(ctx, client.ObjectKeyFromObject(&pod), &pod)
		if err != nil {
			logger.Error(err, "POD not found ("+string(direction)+")",
				"name", fnSpec.WBFunctionRef.Name, "namespace", fnSpec.WBFunctionRef.Namespace)
			return STATUS_NG
		}
		if _, ok := pod.Annotations["ethernet.swb.example.com/network"]; !ok {
			return STATUS_OK
		}

		err = r.createService(ctx, ethConn, fnSpec, othrSideFnKind, pod)
		if err != nil {
			logger.Error(err, "service create fail ("+string(direction)+")",
				"name", ethConn.Spec.To.WBFunctionRef.Name, "namespace", ethConn.Spec.To.WBFunctionRef.Namespace)
			return STATUS_NG
		}

		_, found, _ := unstructured.NestedMap(wbf.Object, "spec", "nextWBFunctions")
		if !found {
			configmap := &corev1.ConfigMap{}
			configmap.Name = fmt.Sprintf("ethcrl.%s", podName)
			configmap.Namespace = fnSpec.WBFunctionRef.Namespace
			_, err := controllerutil.CreateOrUpdate(ctx, r.Client, configmap, func() error {
				raw, _ := json.MarshalIndent(struct {
				}{}, "", "  ")
				configmap.Data = map[string]string{
					"config": string(raw),
				}
				return controllerutil.SetControllerReference(ethConn, configmap, r.Scheme)
			})
			if err != nil {
				logger.Error(err, "configmap create fail ("+string(direction)+")",
					"name", ethConn.Spec.To.WBFunctionRef.Name, "namespace", ethConn.Spec.To.WBFunctionRef.Namespace)
				return STATUS_NG
			}
		}
	case DirectionFrom:

		owbf := &unstructured.Unstructured{}
		owbf.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "WBFunction",
		})
		err := r.Get(ctx, client.ObjectKey{
			Namespace: othrSideFnSpec.WBFunctionRef.Namespace,
			Name:      othrSideFnSpec.WBFunctionRef.Name,
		}, owbf)
		if errors.IsNotFound(err) {
			logger.Error(err, "WBFunctionCR does not exist")
			return STATUS_NG
		}

		othrSidePodName := ""
		switch othrSideFnKind {
		case FUNCTYPE_CPU:
			othrSidePodName = fmt.Sprintf("%s-cpu-pod", othrSideFnSpec.WBFunctionRef.Name)
		case FUNCTYPE_GPU:
			if 0 < len(othrSideFnData.Spec.AcceleratorIDs) {
				othrSidePodName = fmt.Sprintf("%s-mps-dgpu-%s-pod", othrSideFnSpec.WBFunctionRef.Name, strings.ToLower(othrSideFnData.Spec.AcceleratorIDs[0].ID))
			}
		default:
			return STATUS_OK
		}

		pod := corev1.Pod{}
		pod.Name = othrSidePodName
		pod.Namespace = othrSideFnSpec.WBFunctionRef.Namespace
		err = r.Get(ctx, client.ObjectKeyFromObject(&pod), &pod)
		if err != nil {
			logger.Error(err, "POD not found ("+string(direction)+")",
				"name", othrSideFnSpec.WBFunctionRef.Name, "namespace", othrSideFnSpec.WBFunctionRef.Namespace)
			return STATUS_NG
		}
		if _, ok := pod.Annotations["ethernet.swb.example.com/network"]; !ok {
			return STATUS_OK
		}

		extnw := ""

		if v, ok := pod.Annotations["ethernet.swb.example.com/network"]; ok {
			extnw = v
		}

		destIPs := []string{}
		destPorts := []struct {
			Port     int32  `json:"port"`
			Protocol string `json:"protocol"`
		}{}

		switch extnw {
		case "sriov":
			if v, ok := pod.Annotations["k8s.v1.cni.cncf.io/network-status"]; !ok {
				logger.Info("k8s.v1.cni.cncf.io/network-status Annotations not found ("+string(direction)+")",
					"name", othrSideFnSpec.WBFunctionRef.Name, "namespace", othrSideFnSpec.WBFunctionRef.Namespace)
				return STATUS_NG
			} else {
				nwss := []struct {
					Name       string   `json:"name"`
					Interface  *string  `json:"interface"`
					IPs        []string `json:"ips"`
					MAC        *string  `json:"mac"`
					DeviceInfo *struct {
						Type string `json:"type"`
						PCI  *struct {
							PCIAddress string `json:"pci-address"`
						} `json:"pci"`
					} `json:"device-info"`
				}{}
				err := json.Unmarshal([]byte(v), &nwss)
				if err != nil {
					logger.Error(err, "k8s.v1.cni.cncf.io/network-status Annotations fail ("+string(direction)+")",
						"name", othrSideFnSpec.WBFunctionRef.Name, "namespace", othrSideFnSpec.WBFunctionRef.Namespace)
					return STATUS_NG
				}
				for _, nws := range nwss {
					if nws.Interface != nil && nws.MAC != nil && nws.DeviceInfo != nil && nws.DeviceInfo.Type == "pci" {
						destIPs = append(destIPs, nws.IPs...)
					}
				}
				if 0 == len(destIPs) {
					logger.Info("destIP not found 1 ("+string(direction)+")",
						"name", fnSpec.WBFunctionRef.Name, "namespace", fnSpec.WBFunctionRef.Namespace)
					return STATUS_NG
				}
				for _, c := range pod.Spec.Containers {
					for _, p := range c.Ports {
						destPorts = append(destPorts,
							struct {
								Port     int32  `json:"port"`
								Protocol string `json:"protocol"`
							}{
								p.ContainerPort,
								string(p.Protocol),
							})
					}
				}
			}
		case "loadbalancer":
			fallthrough
		case "clusterip":
			service := corev1.Service{}
			service.Name = othrSideFnSpec.WBFunctionRef.Name + "-service"
			service.Namespace = othrSideFnSpec.WBFunctionRef.Namespace
			err := r.Get(ctx, client.ObjectKeyFromObject(&service), &service)
			if err != nil {
				logger.Error(err, "Service not found ("+string(direction)+")",
					"name", fnSpec.WBFunctionRef.Name, "namespace", fnSpec.WBFunctionRef.Namespace)
				return STATUS_NG
			}
			switch service.Spec.Type {
			case corev1.ServiceTypeClusterIP:
				destIPs = append(destIPs, service.Spec.ClusterIP)
			case corev1.ServiceTypeLoadBalancer:
				for _, v := range service.Status.LoadBalancer.Ingress {
					destIPs = append(destIPs, v.IP)
				}
			default:
				logger.Info("destIP not found ("+string(direction)+")",
					"name", fnSpec.WBFunctionRef.Name, "namespace", fnSpec.WBFunctionRef.Namespace)
				return STATUS_NG
			}

			if 0 == len(destIPs) {
				logger.Info("destIP not found 1 ("+string(direction)+")",
					"name", fnSpec.WBFunctionRef.Name, "namespace", fnSpec.WBFunctionRef.Namespace)
				return STATUS_NG
			}

			for _, p := range service.Spec.Ports {
				destPorts = append(destPorts,
					struct {
						Port     int32  `json:"port"`
						Protocol string `json:"protocol"`
					}{
						p.Port,
						string(p.Protocol),
					})
			}
		}

		wbf := &unstructured.Unstructured{}
		wbf.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "WBFunction",
		})
		err = r.Get(ctx, client.ObjectKey{
			Namespace: fnSpec.WBFunctionRef.Namespace,
			Name:      fnSpec.WBFunctionRef.Name,
		}, wbf)
		if errors.IsNotFound(err) {
			logger.Info("WBFunctionCR does not exist")
			return STATUS_NG
		}

		podName := ""
		switch fnKind {
		case FUNCTYPE_CPU:
			podName = fmt.Sprintf("%s-cpu-pod", fnSpec.WBFunctionRef.Name)
		case FUNCTYPE_GPU:
			if 0 < len(fnData.Spec.AcceleratorIDs) {
				podName = fmt.Sprintf("%s-mps-dgpu-%s-pod", fnSpec.WBFunctionRef.Name, strings.ToLower(fnData.Spec.AcceleratorIDs[0].ID))
			}
		default:
			return STATUS_OK
		}

		configmap := &corev1.ConfigMap{}
		configmap.Name = fmt.Sprintf("ethcrl.%s", podName)
		configmap.Namespace = fnSpec.WBFunctionRef.Namespace
		_, err = controllerutil.CreateOrUpdate(ctx, r.Client, configmap, func() error {
			raw, _ := json.MarshalIndent(struct {
				IPs   []string `json:"ips"`
				Ports []struct {
					Port     int32  `json:"port"`
					Protocol string `json:"protocol"`
				} `json:"ports"`
			}{
				destIPs,
				destPorts,
			}, "", "  ")
			configmap.Data = map[string]string{
				"config": string(raw),
			}
			return controllerutil.SetControllerReference(ethConn, configmap, r.Scheme)
		})
		if err != nil {
			logger.Error(err, "configmap create fail ("+string(direction)+")",
				"name", ethConn.Spec.To.WBFunctionRef.Name, "namespace", ethConn.Spec.To.WBFunctionRef.Namespace)
			return STATUS_NG
		}

		_, found, _ := unstructured.NestedMap(wbf.Object, "spec", "previousWBFunctions")
		if !found {
			pod := corev1.Pod{}
			pod.Name = podName
			pod.Namespace = fnSpec.WBFunctionRef.Namespace
			err := r.Get(ctx, client.ObjectKeyFromObject(&pod), &pod)
			if err != nil {
				logger.Info("POD not found ("+string(direction)+")",
					"name", fnSpec.WBFunctionRef.Name, "namespace", fnSpec.WBFunctionRef.Namespace, "err", err)
				return STATUS_NG
			}
			err = r.createService(ctx, ethConn, fnSpec, "external", pod)
			if err != nil {
				logger.Error(err, "service create fail ("+string(direction)+")",
					"name", ethConn.Spec.To.WBFunctionRef.Name, "namespace", ethConn.Spec.To.WBFunctionRef.Namespace)
				return STATUS_NG
			}
		}

	}
	return STATUS_OK
}

func (r *EthernetConnectionReconciler) createService(ctx context.Context,
	ethConn *examplecomv1.EthernetConnection,
	fnSpec *examplecomv1.EthernetFunctionSpec,
	othrSideFnKind string, pod corev1.Pod) error {

	const MetalLBCPUFuncAddressPoolName = "cpufunc-pool"

	extnw := ""
	if v, ok := pod.Annotations["ethernet.swb.example.com/network"]; ok {
		extnw = v
	}
	switch extnw {
	case "sriov":
		return nil
	case "loadbalancer":
		fallthrough
	case "clusterip":
		servicePorts := []corev1.ServicePort{}
		for _, c := range pod.Spec.Containers {
			for _, p := range c.Ports {
				servicePorts = append(servicePorts, corev1.ServicePort{
					Protocol: p.Protocol,
					Port:     p.ContainerPort,
				})
			}
		}
		service := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fnSpec.WBFunctionRef.Name + "-service",
				Namespace: fnSpec.WBFunctionRef.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(&ethConn.ObjectMeta, ethConn.GroupVersionKind()),
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"swb/func-type": pod.Labels["swb/func-type"],
					"swb/func-name": pod.Labels["swb/func-name"],
				},
				Ports: servicePorts,
			},
		}
		switch othrSideFnKind {
		case FUNCTYPE_CPU, FUNCTYPE_GPU:
			switch extnw {
			case "loadbalancer":
				service.Spec.Type = corev1.ServiceTypeLoadBalancer
				service.Annotations = map[string]string{
					"metallb.universe.tf/address-pool": MetalLBCPUFuncAddressPoolName,
				}
			default:
				service.Spec.Type = corev1.ServiceTypeClusterIP
			}
		default:
			service.Spec.Type = corev1.ServiceTypeLoadBalancer
			service.Annotations = map[string]string{
				"metallb.universe.tf/address-pool": MetalLBCPUFuncAddressPoolName,
			}
		}

		return r.Update(ctx, &service)
	default:
		return fmt.Errorf("Error: Annotations[ethernet.swb.example.com/network] == %s", extnw)
	}
}
