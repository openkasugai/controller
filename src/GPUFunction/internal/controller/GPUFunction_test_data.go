/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "GPUFunction/api/v1"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var partitionName1 string = "df-night01-wbfunction-high-infer-main"

var GPUFunction1 = examplecomv1.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"outputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"outputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "df-night01-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night01-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "a100",
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctiontest",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var PCIeConnection1 = PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: WBNamespacedName{
			Name:      "pciefunctiontest",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "fpgafunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "pciefunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var frameworkKernelID int32 = 0
var functionChannelID int32 = 0
var functionIndex int32 = 0
var functionKernelID int32 = 0
var ptuKernelID int32 = 0
var partitionNamefpga1 string = "0"

var FPGAFunction1 = FPGAFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: &partitionNamefpga1,
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "worker1",
		PtuKernelID:       &ptuKernelID,
		RegionName:        "lane0",
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
	},
	Status: FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: WBNamespacedName{
			Name:      "fpgafunctiontest",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
		Status: "pending",
	},
}

var NetworkAttachmentDefinition1 = NetworkAttachmentDefinition{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "k8s.cni.cncf.io",
		Kind:       "NetworkAttachmentDefinition",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "worker1-config-net-sriov",
		Namespace: "default",
		Annotations: map[string]string{
			"k8s.v1.cni.cncf.io/resourceName": "intel.com/intel_sriov_netdevice",
		},
	},
	Spec: NetworkAttachmentDefinitionSpec{
		Config: `{
			"type": "sriov",
			"cniVersion": "0.3.1",
			"name": "worker1-config-net-sriov",
			"ipam": {
				"type": "static"
				}
			}`,
	},
}

var functionIndexg int32 = 99

var GPUFunction2 = examplecomv1.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night02",
			Namespace: "default",
		},
		DeviceType:    "a100",
		FunctionIndex: &functionIndexg,
		FunctionName:  "high-infer",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"outputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"outputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "df-night02-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night02-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "a100",
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctiontest",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var PCIeConnection2 = PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "df-night02",
			Namespace: "default",
		},
		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: WBNamespacedName{
			Name:      "pciefunctiontest",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "fpgafunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "pciefunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var GPUFunction314 = examplecomv1.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night314-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer314",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"outputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"outputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "df-night01-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night01-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "a100",
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctiontest",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}
var gpuconfigdecode = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunc-config-high-infer",
		Namespace: "default",
	},
	Data: map[string]string{
		"gpufunc-config-high-infer.json": ` 
		[{
			"rxProtocol": "DMA",
			"txProtocol":"RTP",
			"sharedMemoryMiB": 256,
			"imageURI": "localhost/gpu-deepstream-app:3.1.0",
			"additionalNetwork": true,
			"virtualNetworkDeviceDriverType": "sriov",
			"envs":{
			  "CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
			  "CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
			  "SHMEM_SECONDARY": "1",
			  "HEIGHT": "1280",
			  "WIDTH": "1280"
			},
			"template":{
			  "apiVersion": "v1",
			  "kind": "Pod",
			  "spec":{
			    "containers":[{
			      "name": "gfunc-hi-1",
			      "workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
			      "command": ["sh", "-c"],
			      "args":["cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgasrc !",
			         "'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
			         "! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA'",
			         "! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1",
			         "%WIDTH%",
			         "%HEIGHT%",
			         "! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1",
			         "model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert !",
			         "'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
			         "%OUTPUTIP%",
			         "%OUTPUTPORT%",
			         "sync=true"],
			      "securityContext":{
			        "privileged": true
			      },
			      "volumeMounts":[{
			        "name": "hugepage-1gi",
			        "mountPath": "/dev/hugepages"
			      },{
			        "name": "host-nvidia-mps",
			        "mountPath": "/tmp/nvidia-mps"
			      },{
			        "name": "dpdk",
			        "mountPath": "/var/run/dpdk"
			      }],
			      "resources":{
			        "requests":{
			          "memory": "32Gi"
			        },
			        "limits":{
			          "hugepages-1Gi": "1Gi"
			        }
			      }
			    }],
			    "volumes":[{
			      "name": "hugepage-1gi",
			      "hostPath":
			       {"path": "/dev/hugepages"}
			    },{
			      "name": "host-nvidia-mps",
			      "hostPath":
			       {"path": "/tmp/nvidia-mps"}
			    },{
			      "name": "dpdk",
			      "hostPath":
			       {"path": "/var/run/dpdk"}
			    }],
			    "hostNetwork": false,
			    "hostIPC": true,
			    "restartPolicy": "Always"
				}
			}
		}]`,
	},
}

var gpuconfigdecode314 = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunc-config-high-infer314",
		Namespace: "default",
	},
	Data: map[string]string{
		"gpufunc-config-high-infer.json": ` 
		[{
			"rxProtocol": "DMA",
			"txProtocol":"RTP",
			"sharedMemoryMiB": 256,
			"imageURI": "localhost/gpu-deepstream-app:3.1.0",
			"additionalNetwork": true,
			"virtualNetworkDeviceDriverType": "sriov",
			"envs":{
			  "CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
			  "CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
			  "SHMEM_SECONDARY": "1",
			  "HEIGHT": "1280",
			  "WIDTH": "1280"
			},
			"template":{
			  "apiVersion": "v1",
			  "kind": "Pod",
			  "spec":{
			    "containers":[{
			      "name": "gfunc-hi-1",
			      "workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
			      "command": ["sh", "-c"],
			      "args":["cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgasrc !",
			         "'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
			         "! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA'",
			         "! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1",
			         "%WIDTH%",
			         "%HEIGHT%",
			         "! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1",
			         "model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert !",
			         "'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
			         "%OUTPUTIP%",
			         "%OUTPUTPORT%",
			         "sync=true"],
			      "securityContext":{
			        "privileged": true
			      },
			      "volumeMounts":[{
			        "name": "hugepage-1gi",
			        "mountPath": "/dev/hugepages"
			      },{
			        "name": "host-nvidia-mps",
			        "mountPath": "/tmp/nvidia-mps"
			      },{
			        "name": "dpdk",
			        "mountPath": "/var/run/dpdk"
			      }],
			      "resources":{
			        "requests":{
			          "memory": "32Gi"
			        },
			        "limits":{
			          "hugepages-1Gi": "1Gi"
			        }
			      },
                  "lifecycle":{
                      "Prestop": {
						"exec": {
							"command": ["sh","-c", "kill -KILL $(pidof gst-launch-1.0)"]
						}
					  }
				  }
			    }],
			    "volumes":[{
			      "name": "hugepage-1gi",
			      "hostPath":
			       {"path": "/dev/hugepages"}
			    },{
			      "name": "host-nvidia-mps",
			      "hostPath":
			       {"path": "/tmp/nvidia-mps"}
			    },{
			      "name": "dpdk",
			      "hostPath":
			       {"path": "/var/run/dpdk"}
			    }],
			    "hostNetwork": false,
			    "hostIPC": true,
			    "restartPolicy": "Always"
				}
			}
		}]`,
	},
}

var gpuconfighigh = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunc-config-high-infer",
		Namespace: "default",
	},
	Data: map[string]string{
		"gpufunc-config-high-infer.json": `
	[{
		"rxProtocol": "DMA",
		"txProtocol":"RTP",
		"sharedMemoryMiB": 256,
		"imageURI": "localhost/gpu-deepstream-app:3.1.0",
		"additionalNetwork": true,
		"virtualNetworkDeviceDriverType": "sriov",
		"envs":{
			"CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
			"CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
			"SHMEM_SECONDARY": "1",
			"HEIGHT": "1280",
			"WIDTH": "1280"
		},
		"template":{
			"apiVersion": "v1",
			"kind": "Pod",
			"spec":{
			"containers":[{
				"name": "gfunc-hi-1",
				"workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
				"command": ["sh", "-c"],
				"args":["cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgasrc !",
				"'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
				"! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA'",
				"! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1",
				"%WIDTH%",
				"%HEIGHT%",
				"! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1",
				"model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert !",
				"'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
				"%OUTPUTIP%",
				"%OUTPUTPORT%",
				"sync=true"],
				"securityContext": {
					"privileged": true
					},
					"volumeMounts": [
					{
					"name": "hugepage-1gi",
					"mountPath": "/dev/hugepages"
					},
					{
					"name": "host-nvidia-mps",
					"mountPath": "/tmp/nvidia-mps"
					},
					{
					"name": "dpdk",
					"mountPath": "/var/run/dpdk"
					}],
					"resources":{
					"requests":{
						"memory": "32Gi"
					},
					"limits":{
						"hugepages-1Gi": "1Gi"
					}
				}
			}],
			"volumes":[{
				"name": "hugepage-1gi",
				"hostPath":
				{"path": "/dev/hugepages"}
			},
			{
				"name": "host-nvidia-mps",
				"hostPath":
				{"path": "/tmp/nvidia-mps"}
			},
			{
				"name": "dpdk",
				"hostPath":
				{"path": "/var/run/dpdk"}
			}],
			"hostNetwork": false,
			"hostIPC": true,
			"restartPolicy": "Always"
			}
		}
	},
	{
		"rxProtocol": "TCP",
		"txProtocol":"RTP",
		"imageURI": "localhost/gpu-deepstream-app_tcprcv:3.1.0",
		"additionalNetwork": true,
		"virtualNetworkDeviceDriverType": "sriov",
		"envs":{
			"CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
			"CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
			"GST_PLUGIN_PATH": "/opt/nvidia/deepstream/deepstream-6.3/fpga-software/tools/tcp_plugins/fpga_depayloader",
			"HEIGHT": "1280",
			"WIDTH": "1280"
		},
		"template":{
			"apiVersion": "v1",
			"kind": "Pod",
			"spec":{
			"containers":[{
				"name": "gfunc-hi-1",
				"workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
				"command": ["sh", "-c"],
				"args":["cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgadepay",
				"%INPUTIP%",
				"%INPUTPORT%",
				"! 'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
				"! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA'",
				"! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1",
				"%WIDTH%",
				"%HEIGHT%",
				"! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1",
				"model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert !",
				"'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
				"%OUTPUTIP%",
				"%OUTPUTPORT%",
				"sync=true"],
				"securityContext":{
				"privileged": true
				},
				"volumeMounts":[{
				"name": "host-nvidia-mps",
				"mountPath": "/tmp/nvidia-mps"
				}]
			}],
			"volumes":[{
				"name": "host-nvidia-mps",
				"hostPath":
				{"path": "/tmp/nvidia-mps"}
			}],
			"hostNetwork": false,
			"hostIPC": true,
			"restartPolicy": "Always"
			}
		}
	}]`,
	},
}

var gpuconfiglow = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunc-config-low-infer",
		Namespace: "default",
	},
	Data: map[string]string{
		"gpufunc-config-low-infer": `
	[{
	"rxProtocol": "DMA",
	"txProtocol":"RTP",
	"sharedMemoryMiB": 256,
	"imageURI": "localhost/gpu-deepstream-app:3.1.0",
	"additionalNetwork": true,
	"virtualNetworkDeviceDriverType": "sriov",
	"envs":{
		"CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
		"CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
		"SHMEM_SECONDARY": "1",
		"HEIGHT": "416",
		"WIDTH": "416"
	},
	"template":{
		"apiVersion": "v1",
		"kind": "Pod",
		"spec":{
		"containers":[{
			"name": "gfunc-n02-lo-1",
			"workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
			"command": ["sh", "-c"],
			"args":["cd /opt/nvidia/deepstream/deepstream-6.3/sources/objectDetector_Yolo/ && gst-launch-1.0 -ev fpgasrc !",
			"'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
			"! nvvideoconvert !",
			"'video/x-raw(memory:NVMM), format=(string)RGBA' !",
			"m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1 ",
			"%WIDTH%",
			"%HEIGHT%",
			"! queue ! nvinfer config-file-path=./config_infer_primary_yoloV3_tiny.txt",
			"batch-size=1 model-engine-file=./model_b1_gpu0_int8.engine ! queue ! nvvideoconvert !",
			"'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
			"%OUTPUTIP%",
			"%OUTPUTPORT%",
			"sync=true"],
			"securityContext":{
			"privileged": true
			},
			"volumeMounts":[{
			"name": "hugepage-1gi",
			"mountPath": "/dev/hugepages"
			},{
			"name": "host-nvidia-mps",
			"mountPath": "/tmp/nvidia-mps"
			},{
			"name": "dpdk",
			"mountPath": "/var/run/dpdk"
			}],
			"resources":{
			"requests":{
				"memory": "32Gi"
			},
			"limits":{
				"hugepages-1Gi": "1Gi"
			}
			}
		}],
		"volumes":[{
			"name": "hugepage-1gi",
			"hostPath":
			{"path": "/dev/hugepages"}
		},{
			"name": "host-nvidia-mps",
			"hostPath":
			{"path": "/tmp/nvidia-mps"}
		},{
			"name": "dpdk",
			"hostPath":
			{"path": "/var/run/dpdk"}
		}],
		"hostNetwork": false,
		"hostIPC": true,
		"restartPolicy": "Always"
		}
	}
	},
	{
	"rxProtocol": "TCP",
	"txProtocol":"RTP",
	"imageURI": "localhost/gpu-deepstream-app_tcprcv:3.1.0",
	"additionalNetwork": true,
	"virtualNetworkDeviceDriverType": "sriov",
	"envs":{
		"CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
		"CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
		"GST_PLUGIN_PATH": "/opt/nvidia/deepstream/deepstream-6.3/fpga-software/tools/tcp_plugins/fpga_depayloader",
		"HEIGHT": "416",
		"WIDTH": "416"
	},
	"template":{
		"apiVersion": "v1",
		"kind": "Pod",
		"spec":{
		"containers":[{
			"name": "gfunc-n02-lo-1",
			"workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
			"command": ["sh", "-c"],
			"args":["cd /opt/nvidia/deepstream/deepstream-6.3/sources/objectDetector_Yolo/ && gst-launch-1.0 -ev fpgadepay",
			"%INPUTIP%",
			"%INPUTPORT%",
			"! 'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
			"' ! nvvideoconvert !",
			"'video/x-raw(memory:NVMM), format=(string)RGBA' !",
			"m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1 ",
			"%WIDTH%",
			"%HEIGHT%",
			"! queue ! nvinfer config-file-path=./config_infer_primary_yoloV3_tiny.txt",
			"batch-size=1 model-engine-file=./model_b1_gpu0_int8.engine ! queue ! nvvideoconvert !",
			"'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
			"%OUTPUTIP%",
			"%OUTPUTPORT%",
			"sync=true"],
			"securityContext":{
			"privileged": true
			},
			"volumeMounts":[{
			"name": "host-nvidia-mps",
			"mountPath": "/tmp/nvidia-mps"
			}]
		}],
		"volumes":[{
			"name": "host-nvidia-mps",
			"hostPath":
			{"path": "/tmp/nvidia-mps"}
		}],
		"hostNetwork": false,
		"hostIPC": true,
		"restartPolicy": "Always"
		}
	}
	}]
	`,
	},
}
var t = metav1.Time{
	Time: time.Now(),
}
var testTime = metav1.Time{
	Time: t.Time.AddDate(0, 0, -1),
}

/*
1-1-1
*/
var partitionName111 string = "gpufunctest111-wbfunction-high-infer-main"

var GPUFunction111 = examplecomv1.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest111-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				ID:            "0",
				PartitionName: &partitionName111,
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctest111",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		NodeName:     "worker1",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Namespace: "default",
					Name:      "gpufunctest111-wbfunction-filter-resize-high-infer-main",
				},
				Port: 0,
			},
		},
		RegionName: "a100",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "gpufunctest111-wbfunction-high-infer-main",
			CommandQueueID:  "gpufunctest111-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
		Params: map[string]intstr.IntOrString{
			"inputIPAddress": {
				StrVal: "192.174.90.141",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"ipAddress": {
				StrVal: "192.174.90.141/24",
				Type:   1,
			},
			"outputIPAddress": {
				StrVal: "192.174.90.10",
				Type:   1,
			},
			"outputPort": {
				IntVal: 2001,
			},
		},
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		FunctionName: "",
		ImageURI:     "",
		ConfigName:   "",
		Status:       "",
	},
}

var PCIeConnection111 = PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest111-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest111",
			Namespace: "default",
		},

		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest111-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest111-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var partitionNamefpga111 string = "0"

var FPGAFunction111 = FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest111-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: &partitionNamefpga111,
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest111",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "worker1",
		PtuKernelID:       &ptuKernelID,
		RegionName:        "lane0",
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
	},
	Status: FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest111",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
		Status: "",
	},
}

/*
1-1-2 cpu-filter-resize-tcp-high-infer
*/

var partitionNamecpufr string = "gpufunctest112-wbfunction-filter-resize-high-infer-main"

var CPUFunctionFilterResize112 = CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest112-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: CPUFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: &partitionNamecpufr,
				ID:            "worker1-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest112",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-high-infer",
		NextFunctions: map[string]FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: WBNamespacedName{
					Name:      "gpufunctest112-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 15,
			},
			"inputIPAddress": {
				StrVal: "192.168.122.50",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"outputIPAddress": {
				StrVal: "192.168.122.100",
				Type:   1,
			},
			"outputPort": {
				IntVal: 16000,
			},
		},
		PreviousFunctions: map[string]FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: WBNamespacedName{
					Name:      "gpufunctest112-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &SharedMemorySpec{
			FilePrefix:      "gpufunctest112-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "gpufunctest112-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest112",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var EthernetConnection112 = EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest112-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest112",
			Namespace: "default",
		},

		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest112-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest112-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}
var partitionName112 string = "gpufunctest112-wbfunction-high-infer-main"

var GPUFunction112 = examplecomv1.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest112-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				ID:            "0",
				PartitionName: &partitionName112,
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctest112",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		NodeName:     "worker1",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Namespace: "default",
					Name:      "gpufunctest112-wbfunction-filter-resize-high-infer-main",
				},
				Port: 0,
			},
		},
		RegionName: "a100",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "gpufunctest112-wbfunction-high-infer-main",
			CommandQueueID:  "gpufunctest112-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
		Params: map[string]intstr.IntOrString{
			"inputIPAddress": {
				StrVal: "192.174.90.142",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15001,
			},
			"ipAddress": {
				StrVal: "192.174.90.142/24",
				Type:   1,
			},
			"outputIPAddress": {
				StrVal: "192.174.90.10",
				Type:   1,
			},
			"outputPort": {
				IntVal: 2011,
			},
		},
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		FunctionName: "",
		ImageURI:     "",
		ConfigName:   "",
		Status:       "",
	},
}

/*
1-2-1
*/
var partitionName121 string = "gpufunctest121-wbfunction-low-infer-main"

var GPUFunction121 = examplecomv1.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest121-wbfunction-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				ID:            "0",
				PartitionName: &partitionName121,
			},
		},
		ConfigName: "gpufunc-config-low-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctest121",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "low-infer",
		NodeName:     "worker1",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Namespace: "default",
					Name:      "gpufunctest121-wbfunction-filter-resize-low-infer-main",
				},
				Port: 0,
			},
		},
		RegionName: "a100",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "gpufunctest121-wbfunction-low-infer-main",
			CommandQueueID:  "gpufunctest121-wbfunction-low-infer-main",
			SharedMemoryMiB: 1,
		},
		Params: map[string]intstr.IntOrString{
			"inputIPAddress": {
				StrVal: "192.174.90.141",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"ipAddress": {
				StrVal: "192.174.90.141/24",
				Type:   1,
			},
			"outputIPAddress": {
				StrVal: "192.174.90.10",
				Type:   1,
			},
			"outputPort": {
				IntVal: 2001,
			},
		},
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		FunctionName: "",
		ImageURI:     "",
		ConfigName:   "",
		Status:       "",
	},
}

var PCIeConnection121 = PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest121-wbconnection-filter-resize-low-infer-main-low-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest121",
			Namespace: "default",
		},

		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest121-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest121-wbfunction-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var partitionNamefpga121 string = "0"

var FPGAFunction121 = FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest121-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: &partitionNamefpga121,
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-low-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest121",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-low-infer",
		NodeName:          "worker1",
		PtuKernelID:       &ptuKernelID,
		RegionName:        "lane0",
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
	},
	Status: FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest121",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-low-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
		Status: "",
	},
}

/*
1-2-2
*/

var partitionNamecpufr122 string = "gpufunctest2-wbfunction-filter-resize-low-infer-main"

var CPUFunctionFilterResize122 = CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest122-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: CPUFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: &partitionNamecpufr122,
				ID:            "worker1-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-low-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest122",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-low-infer",
		NextFunctions: map[string]FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: WBNamespacedName{
					Name:      "gpufunctest122-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 15,
			},
			"inputIPAddress": {
				StrVal: "192.168.122.50",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"outputIPAddress": {
				StrVal: "192.168.122.100",
				Type:   1,
			},
			"outputPort": {
				IntVal: 16000,
			},
		},
		PreviousFunctions: map[string]FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: WBNamespacedName{
					Name:      "gpufunctest122-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &SharedMemorySpec{
			FilePrefix:      "gpufunctest122-wbfunction-filter-resize-low-infer-main",
			CommandQueueID:  "gpufunctest122-wbfunction-filter-resize-low-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest122",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var EthernetConnection122 = EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest122-wbconnection-filter-resize-low-infer-main-low-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "gpufunctest122",
			Namespace: "default",
		},

		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest122-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "gpufunctest122-wbfunction-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}
var partitionName122 string = "gpufunctest2-wbfunction-low-infer-main"

var GPUFunction122 = examplecomv1.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctest122-wbfunction-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				ID:            "0",
				PartitionName: &partitionName122,
			},
		},
		ConfigName: "gpufunc-config-low-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctest122",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "low-infer",
		NodeName:     "worker1",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Namespace: "default",
					Name:      "gpufunctest122-wbfunction-filter-resize-low-infer-main",
				},
				Port: 0,
			},
		},
		RegionName: "a100",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "gpufunctest122-wbfunction-low-infer-main",
			CommandQueueID:  "gpufunctest122-wbfunction-low-infer-main",
			SharedMemoryMiB: 1,
		},
		Params: map[string]intstr.IntOrString{
			"inputIPAddress": {
				StrVal: "192.174.90.142",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15001,
			},
			"ipAddress": {
				StrVal: "192.174.90.142/24",
				Type:   1,
			},
			"outputIPAddress": {
				StrVal: "192.174.90.10",
				Type:   1,
			},
			"outputPort": {
				IntVal: 2011,
			},
		},
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		FunctionName: "",
		ImageURI:     "",
		ConfigName:   "",
		Status:       "",
	},
}

/*
2-1
*/

var partitionNameUPDATE string = "gpufunctestupdate-wbfunction-high-infer-main"

var gpufunctestUPDATE = examplecomv1.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctestupdate-wbfunction-high-infer-main",
		Namespace: "default",
		Finalizers: []string{
			"gpufunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				ID:            "0",
				PartitionName: &partitionNameUPDATE,
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctestupdate",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		NodeName:     "worker1",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Namespace: "default",
					Name:      "gpufunctestupdate-wbfunction-filter-resize-high-infer-main",
				},
				Port: 0,
			},
		},
		RegionName: "a100",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "gpufunctestupdate-wbfunction-high-infer-main",
			CommandQueueID:  "gpufunctestupdate-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
		Params: map[string]intstr.IntOrString{
			"inputIPAddress": {
				StrVal: "192.174.90.141",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"ipAddress": {
				StrVal: "192.174.90.141/24",
				Type:   1,
			},
			"outputIPAddress": {
				StrVal: "192.174.90.10",
				Type:   1,
			},
			"outputPort": {
				IntVal: 2001,
			},
		},
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctestupdate",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "localhost/gpu-deepstream-app:3.1.0",
		ConfigName:   "gpufunc-config-high-infer",
		Status:       "Running",
	},
}

/*
2-2
*/

var partitionNameDELETE string = "gpufunctestdelete-wbfunction-high-infer-main"

var gpufunctestDELETE = examplecomv1.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunctestdelete-wbfunction-high-infer-main",
		Namespace: "default",
		Finalizers: []string{
			"gpufunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				ID:            "0",
				PartitionName: &partitionNameDELETE,
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctestdelete",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		NodeName:     "worker1",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Namespace: "default",
					Name:      "gpufunctestdelete-wbfunction-filter-resize-high-infer-main",
				},
				Port: 0,
			},
		},
		RegionName: "a100",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "gpufunctestdelete-wbfunction-high-infer-main",
			CommandQueueID:  "gpufunctestdelete-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
		Params: map[string]intstr.IntOrString{
			"inputIPAddress": {
				StrVal: "192.174.90.141",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"ipAddress": {
				StrVal: "192.174.90.141/24",
				Type:   1,
			},
			"outputIPAddress": {
				StrVal: "192.174.90.10",
				Type:   1,
			},
			"outputPort": {
				IntVal: 2001,
			},
		},
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctestdelete",
			Namespace: "default",
		},
		FunctionName:  "high-infer",
		ImageURI:      "localhost/gpu-deepstream-app:3.1.0",
		ConfigName:    "gpufunc-config-high-infer",
		Status:        "Running",
		FunctionIndex: &functionIndex,
	},
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvsInfo) DeepCopyInto(out *EnvsInfo) {
	*out = *in
	if in.EachEnv != nil {
		in, out := &in.EachEnv, &out.EachEnv
		*out = make([]EnvsData, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvsInfo.
func (in *EnvsInfo) DeepCopy() *EnvsInfo {
	if in == nil {
		return nil
	}
	out := new(EnvsInfo)
	in.DeepCopyInto(out)
	return out
}

/*
	registration of PCIeConnection
*/

// groupversion_info.go
var (
	// GroupVersion is group version used to register these objects
	GroupVersionpcie = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilderpcie = &scheme.Builder{GroupVersion: GroupVersionpcie}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemepcie = SchemeBuilderpcie.AddToScheme
)

// pcieconnection_types.go
// PCIeConnection difines the PCIeConnection CR
type PCIeConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PCIeConnectionSpec   `json:"spec,omitempty"`
	Status PCIeConnectionStatus `json:"status,omitempty"`
}

// PCIeConnectionSpec defines the desired state of PCIeConnection
type PCIeConnectionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef WBNamespacedName `json:"dataFlowRef"`
	From        PCIeFunctionSpec `json:"from"`
	To          PCIeFunctionSpec `json:"to"`
}

// PCIeConnectionStatus defines the observed state of PCIeConnection
type PCIeConnectionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef  WBNamespacedName   `json:"dataFlowRef"`
	From         PCIeFunctionStatus `json:"from"`
	To           PCIeFunctionStatus `json:"to"`
	SharedMemory SharedMemoryStatus `json:"sharedMemory,omitempty"`
	StartTime    metav1.Time        `json:"startTime"`
	//+kubebuilder:default=Pending
	Status string `json:"status"`
}

type PCIeFunctionSpec struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
}

type PCIeFunctionStatus struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
	//+kubebuilder:default=INIT
	Status string `json:"status"`
}

type SharedMemoryStatus struct {
	// +optional
	Status string `json:"status,omitempty"`
}

// PCIeConnectionList contains a list of PCIeConnection
type PCIeConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PCIeConnection `json:"items"`
}

type WBNamespacedName struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func init() {
	SchemeBuilderpcie.Register(&PCIeConnection{}, &PCIeConnectionList{})
}

// zz_generated.deepcopy.go
// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnection) DeepCopyInto(out *PCIeConnection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnection.
func (in *PCIeConnection) DeepCopy() *PCIeConnection {
	if in == nil {
		return nil
	}
	out := new(PCIeConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PCIeConnection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnectionList) DeepCopyInto(out *PCIeConnectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PCIeConnection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnectionList.
func (in *PCIeConnectionList) DeepCopy() *PCIeConnectionList {
	if in == nil {
		return nil
	}
	out := new(PCIeConnectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PCIeConnectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnectionSpec) DeepCopyInto(out *PCIeConnectionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	out.From = in.From
	out.To = in.To
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnectionSpec.
func (in *PCIeConnectionSpec) DeepCopy() *PCIeConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(PCIeConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnectionStatus) DeepCopyInto(out *PCIeConnectionStatus) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	out.From = in.From
	out.To = in.To
	out.SharedMemory = in.SharedMemory
	in.StartTime.DeepCopyInto(&out.StartTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnectionStatus.
func (in *PCIeConnectionStatus) DeepCopy() *PCIeConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(PCIeConnectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeFunctionSpec) DeepCopyInto(out *PCIeFunctionSpec) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeFunctionSpec.
func (in *PCIeFunctionSpec) DeepCopy() *PCIeFunctionSpec {
	if in == nil {
		return nil
	}
	out := new(PCIeFunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeFunctionStatus) DeepCopyInto(out *PCIeFunctionStatus) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeFunctionStatus.
func (in *PCIeFunctionStatus) DeepCopy() *PCIeFunctionStatus {
	if in == nil {
		return nil
	}
	out := new(PCIeFunctionStatus)
	in.DeepCopyInto(out)
	return out
}

/*
	registration of FPGAFunction
*/

// fpgafnction_types.go
// FPGAFunctionSpec defines the desired state of FPGAFunction
type FPGAFunctionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DataFlowRef       WBNamespacedName  `json:"dataFlowRef"`
	FunctionName      string            `json:"functionName"`
	NodeName          string            `json:"nodeName"`
	DeviceType        string            `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo       `json:"acceleratorIDs"`
	RegionName        string            `json:"regionName"`
	FunctionIndex     *int32            `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo        `json:"envs,omitempty"`
	ConfigName        string            `json:"configName"`
	SharedMemory      *SharedMemorySpec `json:"sharedMemory,omitempty"`
	FunctionKernelID  *int32            `json:"functionKernelID"`
	FunctionChannelID *int32            `json:"functionChannelID"`
	PtuKernelID       *int32            `json:"ptuKernelID"`
	FrameworkKernelID *int32            `json:"frameworkKernelID"`
	Rx                RxTxSpec          `json:"rx"`
	Tx                RxTxSpec          `json:"tx"`
}

type RxTxSpec struct {
	Protocol        string  `json:"protocol"`
	IPAddress       *string `json:"ipAddress,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	SubnetAddress   *string `json:"subnetAddress,omitempty"`
	GatewayAddress  *string `json:"gatewayAddress,omitempty"`
	DMAChannelID    *int32  `json:"dmaChannelID,omitempty"`
	FDMAConnectorID *int32  `json:"fdmaConnectorID,omitempty"`
}

// FPGAFunctionStatus defines the observed state of FPGAFunction
type FPGAFunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	StartTime metav1.Time `json:"startTime"`
	//+kubebuilder:default=Pending
	Status              string                `json:"status"`
	DataFlowRef         WBNamespacedName      `json:"dataFlowRef"`
	FunctionName        string                `json:"functionName"`
	ParentBitstreamName string                `json:"parentBitstreamName"`
	ChildBitstreamName  string                `json:"childBitstreamName"`
	SharedMemory        *SharedMemorySpec     `json:"sharedMemory,omitempty"`
	FunctionKernelID    int32                 `json:"functionKernelID"`
	FunctionChannelID   int32                 `json:"functionChannelID"`
	PtuKernelID         int32                 `json:"ptuKernelID"`
	FrameworkKernelID   int32                 `json:"frameworkKernelID"`
	Rx                  RxTxSpec              `json:"rx"`
	Tx                  RxTxSpec              `json:"tx"`
	AcceleratorStatuses []AccStatusesByDevice `json:"acceleratorStatuses,omitempty"`
}

type AccStatusesByDevice struct {
	PartitionName *string       `json:"partitionName,omitempty"`
	Statused      []AccStatuses `json:"statuses,omitempty"`
}

type AccStatuses struct {
	AcceleratorID *string `json:"acceleratorID,omitempty"`
	Status        *string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FPGAFunction is the Schema for the fpgafunctions API
type FPGAFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FPGAFunctionSpec   `json:"spec,omitempty"`
	Status FPGAFunctionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FPGAFunctionList contains a list of FPGAFunction
type FPGAFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FPGAFunction `json:"items"`
}

func init() {
	func() *scheme.Builder {
		SchemeBuilderfpga.SchemeBuilder.Register(func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(SchemeBuilderfpga.GroupVersion, []runtime.Object{&FPGAFunction{}, &FPGAFunctionList{}}...)
			metav1.AddToGroupVersion(scheme, SchemeBuilderfpga.GroupVersion)
			return nil
		})
		return SchemeBuilderfpga
	}()
}

// groupversion_info.go
var (
	// GroupVersion is group version used to register these objects
	GroupVersionfpga = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilderfpga = &scheme.Builder{GroupVersion: GroupVersionfpga}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemefpga = SchemeBuilderfpga.AddToScheme
)

// zz_generated.deepcopy.go
// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccStatuses) DeepCopyInto(out *AccStatuses) {
	*out = *in
	if in.AcceleratorID != nil {
		in, out := &in.AcceleratorID, &out.AcceleratorID
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccStatuses.
func (in *AccStatuses) DeepCopy() *AccStatuses {
	if in == nil {
		return nil
	}
	out := new(AccStatuses)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccStatusesByDevice) DeepCopyInto(out *AccStatusesByDevice) {
	*out = *in
	if in.PartitionName != nil {
		in, out := &in.PartitionName, &out.PartitionName
		*out = new(string)
		**out = **in
	}
	if in.Statused != nil {
		in, out := &in.Statused, &out.Statused
		*out = make([]AccStatuses, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccStatusesByDevice.
func (in *AccStatusesByDevice) DeepCopy() *AccStatusesByDevice {
	if in == nil {
		return nil
	}
	out := new(AccStatusesByDevice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunction) DeepCopyInto(out *FPGAFunction) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunction.
func (in *FPGAFunction) DeepCopy() *FPGAFunction {
	if in == nil {
		return nil
	}
	out := new(FPGAFunction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FPGAFunction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunctionList) DeepCopyInto(out *FPGAFunctionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FPGAFunction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunctionList.
func (in *FPGAFunctionList) DeepCopy() *FPGAFunctionList {
	if in == nil {
		return nil
	}
	out := new(FPGAFunctionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FPGAFunctionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunctionSpec) DeepCopyInto(out *FPGAFunctionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.AcceleratorIDs != nil {
		in, out := &in.AcceleratorIDs, &out.AcceleratorIDs
		*out = make([]AccIDInfo, len(*in))
		copy(*out, *in)
	}
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]EnvsInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	if in.FunctionKernelID != nil {
		in, out := &in.FunctionKernelID, &out.FunctionKernelID
		*out = new(int32)
		**out = **in
	}
	if in.FunctionChannelID != nil {
		in, out := &in.FunctionChannelID, &out.FunctionChannelID
		*out = new(int32)
		**out = **in
	}
	if in.PtuKernelID != nil {
		in, out := &in.PtuKernelID, &out.PtuKernelID
		*out = new(int32)
		**out = **in
	}
	if in.FrameworkKernelID != nil {
		in, out := &in.FrameworkKernelID, &out.FrameworkKernelID
		*out = new(int32)
		**out = **in
	}
	in.Rx.DeepCopyInto(&out.Rx)
	in.Tx.DeepCopyInto(&out.Tx)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunctionSpec.
func (in *FPGAFunctionSpec) DeepCopy() *FPGAFunctionSpec {
	if in == nil {
		return nil
	}
	out := new(FPGAFunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunctionStatus) DeepCopyInto(out *FPGAFunctionStatus) {
	*out = *in
	in.StartTime.DeepCopyInto(&out.StartTime)
	out.DataFlowRef = in.DataFlowRef
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	in.Rx.DeepCopyInto(&out.Rx)
	in.Tx.DeepCopyInto(&out.Tx)
	if in.AcceleratorStatuses != nil {
		in, out := &in.AcceleratorStatuses, &out.AcceleratorStatuses
		*out = make([]AccStatusesByDevice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunctionStatus.
func (in *FPGAFunctionStatus) DeepCopy() *FPGAFunctionStatus {
	if in == nil {
		return nil
	}
	out := new(FPGAFunctionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RxTxSpec) DeepCopyInto(out *RxTxSpec) {
	*out = *in
	if in.IPAddress != nil {
		in, out := &in.IPAddress, &out.IPAddress
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.SubnetAddress != nil {
		in, out := &in.SubnetAddress, &out.SubnetAddress
		*out = new(string)
		**out = **in
	}
	if in.GatewayAddress != nil {
		in, out := &in.GatewayAddress, &out.GatewayAddress
		*out = new(string)
		**out = **in
	}
	if in.DMAChannelID != nil {
		in, out := &in.DMAChannelID, &out.DMAChannelID
		*out = new(int32)
		**out = **in
	}
	if in.FDMAConnectorID != nil {
		in, out := &in.FDMAConnectorID, &out.FDMAConnectorID
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RxTxSpec.
func (in *RxTxSpec) DeepCopy() *RxTxSpec {
	if in == nil {
		return nil
	}
	out := new(RxTxSpec)
	in.DeepCopyInto(out)
	return out
}

// Function CR structure
type FunctionData struct {
	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     *int32                        `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      SharedMemorySpec              `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        *string                       `json:"configName,omitempty"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
}

type AccIDInfo struct {
	PartitionName *string `json:"partitionName,omitempty"`
	ID            string  `json:"id"`
}

type EnvsInfo struct {
	PartitionName string     `json:"partitionName"`
	EachEnv       []EnvsData `json:"eachEnv"`
}

// Environmental information
type EnvsData struct {
	EnvKey   string `json:"envKey"`
	EnvValue string `json:"envValue"`
}

// FPGA Device information
type RxTxSpecFunc struct {
	Protocol        string  `json:"protocol,omitempty"`
	IPAddress       *string `json:"ipAddress,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	SubnetAddress   *string `json:"subnetAddress,omitempty"`
	GatewayAddress  *string `json:"gatewayAddress,omitempty"`
	DMAChannelID    *int32  `json:"dmaChannelID,omitempty"`
	FDMAConnectorID *int32  `json:"fdmaConnectorID,omitempty"`
}

// Shared memory information
type SharedMemorySpec struct {
	FilePrefix      string `json:"filePrefix"`
	CommandQueueID  string `json:"commandQueueID"`
	SharedMemoryMiB int32  `json:"sharedMemoryMiB"`
}

// Structure for acquiring Phase3FPGA information
type Phase3Data struct {
	NodeName        string        `json:"nodeName"`
	DeviceFilePaths []string      `json:"deviceFilePaths"`
	NetworkInfo     []NetworkData `json:"networkInfo"`
}

type NetworkData struct {
	DeviceIndex    int32    `json:"deviceIndex"`
	LaneIndex      int32    `json:"laneIndex"`
	IPAddress      string   `json:"ipAddress"`
	SubnetAddress  string   `json:"subnetAddress"`
	GatewayAddress string   `json:"gatewayAddress"`
	MACAddress     string   `json:"macAddress"`
	Rx             FPGARxTx `json:"rx"`
	Tx             FPGARxTx `json:"tx"`
}

type FromToWBFunction struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
	Port          int32            `json:"port"`
}

type FPGARxTx struct {
	Protocol  string `json:"protocol"`
	StartPort int32  `json:"startPort,omitempty"`
	EndPort   int32  `json:"endPort,omitempty"`
}

type AnyData struct {
	Functions FrameSizeData `json:"functions,omitempty"`
}

type FrameSizeData struct {
	InputWidth   int32 `json:"i_width"`
	InputHeight  int32 `json:"i_height"`
	OutputWidth  int32 `json:"o_width"`
	OutputHeight int32 `json:"o_height"`
}

/*
	registration of NetworkAttachmentDefinition
*/

// This defines NetworkAttachmentDefinition CR
type NetworkAttachmentDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkAttachmentDefinitionSpec `json:"spec,omitempty"`
}

// PCIeConnectionSpec defines the desired state of PCIeConnection
type NetworkAttachmentDefinitionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Config string `json:"config"`
}

// EthernetConnectionList contains a list of EthernetConnection
type NetworkAttachmentDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkAttachmentDefinition `json:"items"`
}

func init() {
	SchemeBuilder1.Register(&NetworkAttachmentDefinition{}, &NetworkAttachmentDefinitionList{})
}

var (
	// GroupVersion is group version used to register these objects
	GroupVersion1 = schema.GroupVersion{Group: "k8s.cni.cncf.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder1 = &scheme.Builder{GroupVersion: GroupVersion1}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme1 = SchemeBuilder1.AddToScheme
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinition) DeepCopyInto(out *NetworkAttachmentDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnection.
func (in *NetworkAttachmentDefinition) DeepCopy() *NetworkAttachmentDefinition {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkAttachmentDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinitionList) DeepCopyInto(out *NetworkAttachmentDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NetworkAttachmentDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionList.
func (in *NetworkAttachmentDefinitionList) DeepCopy() *NetworkAttachmentDefinitionList {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkAttachmentDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinitionSpec) DeepCopyInto(out *NetworkAttachmentDefinitionSpec) {
	*out = *in
	out.Config = in.Config
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionSpec.
func (in *NetworkAttachmentDefinitionSpec) DeepCopy() *NetworkAttachmentDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionConfigMap) DeepCopyInto(out *FunctionConfigMap) {
	*out = *in
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionConfigMap.
func (in *FunctionConfigMap) DeepCopy() *FunctionConfigMap {
	if in == nil {
		return nil
	}
	out := new(FunctionConfigMap)
	in.DeepCopyInto(out)
	return out
}

type BitstreamData struct {
	File string `json:"file,omitempty"`
	ID   string `json:"id,omitempty"`
}

type FunctionConfigMap struct {
	RxProtocol      string            `json:"rxProtocol,omitempty"`
	TxProtocol      string            `json:"txProtocol,omitempty"`
	SharedMemoryMiB int32             `json:"sharedMemoryMiB,omitempty"`
	ImageURI        string            `json:"imageURI,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	//	ParentBitstream string            `json:"parentBitstream,omitempty"`
	//	ChildBitstream  string            `json:"childBitstream,omitempty"`
	ParentBitstream       BitstreamData `json:"parentBitstream,omitempty"`
	ChildBitstream        BitstreamData `json:"childBitstream,omitempty"`
	Commands              []string      `json:"commands,omitempty"`
	Args                  []string      `json:"args,omitempty"`
	Parameters            AnyData       `json:"parameters,omitempty"`
	FunctionDedicatedInfo string        `json:"functionDedicatedInfo,omitempty"`
}

/*
	registration of EternetConnection
*/

// ethernetconnection_types.go
// EthernetConnectionSpec defines the desired state of EthernetConnection
type EthernetConnectionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef WBNamespacedName     `json:"dataFlowRef"`
	From        EthernetFunctionSpec `json:"from"`
	To          EthernetFunctionSpec `json:"to"`
}

// EthernetConnectionStatus defines the observed state of EthernetConnection
type EthernetConnectionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef WBNamespacedName       `json:"dataFlowRef"`
	From        EthernetFunctionStatus `json:"from"`
	To          EthernetFunctionStatus `json:"to"`
	StartTime   metav1.Time            `json:"startTime"`
	//+kubebuilder:default=Pending
	Status string `json:"status"`
}

type EthernetFunctionSpec struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
}

type EthernetFunctionStatus struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
	//+kubebuilder:default=INIT
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="FROMFUNC_STATUS",type="string",JSONPath=".status.from.status"
//+kubebuilder:printcolumn:name="TOFUNC_STATUS",type="string",JSONPath=".status.to.status"
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// EthernetConnection is the Schema for the ethernetconnections API
type EthernetConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EthernetConnectionSpec   `json:"spec,omitempty"`
	Status EthernetConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EthernetConnectionList contains a list of EthernetConnection
type EthernetConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EthernetConnection `json:"items"`
}

func init() {
	SchemeBuilderether.Register(&EthernetConnection{}, &EthernetConnectionList{})
}

// groupversion_info.go
var (
	// GroupVersion is group version used to register these objects
	GroupVersionether = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilderether = &scheme.Builder{GroupVersion: GroupVersionether}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemeether = SchemeBuilderether.AddToScheme
)

// zz_generated.deepcopy.go
// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EthernetConnection) DeepCopyInto(out *EthernetConnection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnection.
func (in *EthernetConnection) DeepCopy() *EthernetConnection {
	if in == nil {
		return nil
	}
	out := new(EthernetConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EthernetConnection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EthernetConnectionList) DeepCopyInto(out *EthernetConnectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EthernetConnection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionList.
func (in *EthernetConnectionList) DeepCopy() *EthernetConnectionList {
	if in == nil {
		return nil
	}
	out := new(EthernetConnectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EthernetConnectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EthernetConnectionSpec) DeepCopyInto(out *EthernetConnectionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	out.From = in.From
	out.To = in.To
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionSpec.
func (in *EthernetConnectionSpec) DeepCopy() *EthernetConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(EthernetConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EthernetConnectionStatus) DeepCopyInto(out *EthernetConnectionStatus) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	out.From = in.From
	out.To = in.To
	in.StartTime.DeepCopyInto(&out.StartTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionStatus.
func (in *EthernetConnectionStatus) DeepCopy() *EthernetConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(EthernetConnectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EthernetFunctionSpec) DeepCopyInto(out *EthernetFunctionSpec) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetFunctionSpec.
func (in *EthernetFunctionSpec) DeepCopy() *EthernetFunctionSpec {
	if in == nil {
		return nil
	}
	out := new(EthernetFunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EthernetFunctionStatus) DeepCopyInto(out *EthernetFunctionStatus) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetFunctionStatus.
func (in *EthernetFunctionStatus) DeepCopy() *EthernetFunctionStatus {
	if in == nil {
		return nil
	}
	out := new(EthernetFunctionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGARxTx) DeepCopyInto(out *FPGARxTx) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGARxTx.
func (in *FPGARxTx) DeepCopy() *FPGARxTx {
	if in == nil {
		return nil
	}
	out := new(FPGARxTx)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FromToWBFunction) DeepCopyInto(out *FromToWBFunction) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FromToWBFunction.
func (in *FromToWBFunction) DeepCopy() *FromToWBFunction {
	if in == nil {
		return nil
	}
	out := new(FromToWBFunction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionData) DeepCopyInto(out *FunctionData) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.AcceleratorIDs != nil {
		in, out := &in.AcceleratorIDs, &out.AcceleratorIDs
		*out = make([]AccIDInfo, len(*in))
		copy(*out, *in)
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]EnvsInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RequestMemorySize != nil {
		in, out := &in.RequestMemorySize, &out.RequestMemorySize
		*out = new(int32)
		**out = **in
	}
	out.SharedMemory = in.SharedMemory
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(string)
		**out = **in
	}
	if in.ConfigName != nil {
		in, out := &in.ConfigName, &out.ConfigName
		*out = new(string)
		**out = **in
	}
	if in.PreviousFunctions != nil {
		in, out := &in.PreviousFunctions, &out.PreviousFunctions
		*out = make(map[string]FromToWBFunction, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NextFunctions != nil {
		in, out := &in.NextFunctions, &out.NextFunctions
		*out = make(map[string]FromToWBFunction, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]intstr.IntOrString, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionData.
func (in *FunctionData) DeepCopy() *FunctionData {
	if in == nil {
		return nil
	}
	out := new(FunctionData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkData) DeepCopyInto(out *NetworkData) {
	*out = *in
	out.Rx = in.Rx
	out.Tx = in.Tx
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkData.
func (in *NetworkData) DeepCopy() *NetworkData {
	if in == nil {
		return nil
	}
	out := new(NetworkData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Phase3Data) DeepCopyInto(out *Phase3Data) {
	*out = *in
	if in.DeviceFilePaths != nil {
		in, out := &in.DeviceFilePaths, &out.DeviceFilePaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NetworkInfo != nil {
		in, out := &in.NetworkInfo, &out.NetworkInfo
		*out = make([]NetworkData, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Phase3Data.
func (in *Phase3Data) DeepCopy() *Phase3Data {
	if in == nil {
		return nil
	}
	out := new(Phase3Data)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SharedMemorySpec) DeepCopyInto(out *SharedMemorySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SharedMemorySpec.
func (in *SharedMemorySpec) DeepCopy() *SharedMemorySpec {
	if in == nil {
		return nil
	}
	out := new(SharedMemorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBNamespacedName) DeepCopyInto(out *WBNamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBNamespacedName.
func (in *WBNamespacedName) DeepCopy() *WBNamespacedName {
	if in == nil {
		return nil
	}
	out := new(WBNamespacedName)
	in.DeepCopyInto(out)
	return out
}

/*
	registration of CPUFucntion
*/

// cpufunction_common_types.go
const (
	Width            = "WIDTH"
	Height           = "HEIGHT"
	ArgsWidth        = "%" + Width + "%"
	ArgsHeight       = "%" + Height + "%"
	ArgsArpIP        = "%ARPIP%"
	ArgsIP           = "%IP%"
	ArgsPort         = "%PORT%"
	ArgsArpMAC       = "%MAC%"
	ChangeArgsWidth  = "width="
	ChangeArgsHeight = "height="
	ChangeArgsIP     = "host="
	ChangeArgsPort   = "port="
	ArgsReceiving    = "%RECEIVING%"
	ArgsNum          = "%NUM%"
	ArgsForwarding   = "%FORWARDING%"
	ArgsMemSize      = "%MEMSIZE%"
)

type CPUFuncConfig struct {
	RxProtocol                     *string           `json:"rxProtocol,omitempty"`
	TxProtocol                     *string           `json:"txProtocol,omitempty"`
	SharedMemoryGiB                *int32            `json:"sharedMemoryGiB,omitempty"`
	VirtualNetworkDeviceDriverType string            `json:"virtualNetworkDeviceDriverType,omitempty"`
	AdditionalNetwork              bool              `json:"additionalNetwork,omitempty"`
	CopyMemorySize                 *string           `json:"copyMemorySize,omitempty"`
	ImageURI                       string            `json:"imageURI"`
	Envs                           map[string]string `json:"envs"`
	Template                       PodTemplate       `json:"template"`
	Annotations                    map[string]string `json:"annotations,omitempty"`
	Labels                         map[string]string `json:"labels,omitempty"`
	IPAM                           []string          `json:"ipam,omitempty"`
}

type PodTemplate struct {
	metav1.TypeMeta `json:",inline"`
	Spec            CPUPodSpec `json:"spec"`
	// Spec corev1.PodSpec `json:"spec,omitempty"` // TODO
}

type CPUPodSpec struct {
	Volumes               []corev1.Volume      `json:"volumes,omitempty" `
	Containers            []CPUContainer       `json:"containers"`
	RestartPolicy         corev1.RestartPolicy `json:"restartPolicy,omitempty"`
	HostNetwork           bool                 `json:"hostNetwork,omitempty"`
	HostIPC               bool                 `json:"hostIPC,omitempty"`
	ShareProcessNamespace bool                 `json:"shareProcessNamespace,omitempty"`
}

type CPUContainer struct {
	Name            *string                     `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Command         []string                    `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	Args            []string                    `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	WorkingDir      string                      `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	SecurityContext *corev1.SecurityContext     `json:"securityContext,omitempty" protobuf:"bytes,15,opt,name=securityContext"`
	VolumeMounts    []corev1.VolumeMount        `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	Image           string                      `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	Ports           []corev1.ContainerPort      `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
}

// cpufunction_types.go
// CPUFunctionSpec defines the desired state of CPUFunction
type CPUFunctionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     *int32                        `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      *SharedMemorySpec             `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        string                        `json:"configName"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
}

// CPUFunctionStatus defines the observed state of CPUFunction
type CPUFunctionStatus struct {
	DataFlowRef                    WBNamespacedName  `json:"dataFlowRef"`
	FunctionName                   string            `json:"functionName"`
	ImageURI                       string            `json:"imageURI"`
	SharedMemory                   *SharedMemorySpec `json:"sharedMemory,omitempty"`
	RxProtocol                     *string           `json:"rxProtocol,omitempty"`
	TxProtocol                     *string           `json:"txProtocol,omitempty"`
	ConfigName                     string            `json:"configName"`
	VirtualNetworkDeviceDriverType string            `json:"virtualNetworkDeviceDriverType,omitempty"`
	AdditionalNetwork              *bool             `json:"additionalNetwork,omitempty"`
	FunctionIndex                  *int32            `json:"functionIndex,omitempty"`
	StartTime                      metav1.Time       `json:"startTime"`
	//+kubebuilder:default=Pending
	Status              string                   `json:"status"`
	IPAddress           *string                  `json:"Ip,omitempty"`
	AcceleratorStatuses []AccStatusesByContainer `json:"acceleratorStatuses,omitempty"`
}

type AccStatusesByContainer struct {
	PartitionName *string       `json:"partitionName,omitempty"`
	Statuses      []AccStatuses `json:"statuses,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// CPUFunction is the Schema for the cpufunctions API
type CPUFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CPUFunctionSpec   `json:"spec,omitempty"`
	Status CPUFunctionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CPUFunctionList contains a list of CPUFunction
type CPUFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CPUFunction `json:"items"`
}

func init() {
	SchemeBuildercpu.Register(&CPUFunction{}, &CPUFunctionList{})
}

// Function-based CR structure
type FunctionStatusData struct {
	Status            string           `json:"status"`
	FunctionIndex     *int32           `json:"functionIndex"`
	FunctionKernelID  *int32           `json:"functionKernelID,omitempty"`
	FunctionChannelID *int32           `json:"functionChannelID,omitempty"`
	PtuKernelID       *int32           `json:"ptuKernelID,omitempty"`
	FrameworkKernelID *int32           `json:"frameworkKernelID,omitempty"`
	Rx                RxTxData         `json:"rx,omitempty"`
	Tx                RxTxData         `json:"tx,omitempty"`
	SharedMemory      SharedMemorySpec `json:"sharedMemory,omitempty"`
}

// FPGADevice Connection Info
type RxTxData struct {
	Protocol         string  `json:"protocol,omitempty"`
	IPAddress        *string `json:"ipAddress,omitempty"`
	Port             *int32  `json:"port,omitempty"`
	SubnetAddress    *string `json:"subnetAddress,omitempty"`
	GatewayAddress   *string `json:"gatewayAddress,omitempty"`
	DMAChannelID     *int32  `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32  `json:"lldmaConnectorID,omitempty"`
}

// groupversion_info.go
var (
	// GroupVersion is group version used to register these objects
	GroupVersioncpu = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuildercpu = &scheme.Builder{GroupVersion: GroupVersioncpu}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemecpu = SchemeBuildercpu.AddToScheme
)

// zz_generated.deepcopy.go
// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccIDInfo) DeepCopyInto(out *AccIDInfo) {
	*out = *in
	if in.PartitionName != nil {
		in, out := &in.PartitionName, &out.PartitionName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccIDInfo.
func (in *AccIDInfo) DeepCopy() *AccIDInfo {
	if in == nil {
		return nil
	}
	out := new(AccIDInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccStatusesByContainer) DeepCopyInto(out *AccStatusesByContainer) {
	*out = *in
	if in.PartitionName != nil {
		in, out := &in.PartitionName, &out.PartitionName
		*out = new(string)
		**out = **in
	}
	if in.Statuses != nil {
		in, out := &in.Statuses, &out.Statuses
		*out = make([]AccStatuses, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccStatusesByContainer.
func (in *AccStatusesByContainer) DeepCopy() *AccStatusesByContainer {
	if in == nil {
		return nil
	}
	out := new(AccStatusesByContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnyData) DeepCopyInto(out *AnyData) {
	*out = *in
	out.Functions = in.Functions
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnyData.
func (in *AnyData) DeepCopy() *AnyData {
	if in == nil {
		return nil
	}
	out := new(AnyData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BitstreamData) DeepCopyInto(out *BitstreamData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BitstreamData.
func (in *BitstreamData) DeepCopy() *BitstreamData {
	if in == nil {
		return nil
	}
	out := new(BitstreamData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUContainer) DeepCopyInto(out *CPUContainer) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]corev1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]corev1.ContainerPort, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUContainer.
func (in *CPUContainer) DeepCopy() *CPUContainer {
	if in == nil {
		return nil
	}
	out := new(CPUContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUFuncConfig) DeepCopyInto(out *CPUFuncConfig) {
	*out = *in
	if in.RxProtocol != nil {
		in, out := &in.RxProtocol, &out.RxProtocol
		*out = new(string)
		**out = **in
	}
	if in.TxProtocol != nil {
		in, out := &in.TxProtocol, &out.TxProtocol
		*out = new(string)
		**out = **in
	}
	if in.SharedMemoryGiB != nil {
		in, out := &in.SharedMemoryGiB, &out.SharedMemoryGiB
		*out = new(int32)
		**out = **in
	}
	if in.CopyMemorySize != nil {
		in, out := &in.CopyMemorySize, &out.CopyMemorySize
		*out = new(string)
		**out = **in
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Template.DeepCopyInto(&out.Template)
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.IPAM != nil {
		in, out := &in.IPAM, &out.IPAM
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUFuncConfig.
func (in *CPUFuncConfig) DeepCopy() *CPUFuncConfig {
	if in == nil {
		return nil
	}
	out := new(CPUFuncConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUFunction) DeepCopyInto(out *CPUFunction) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUFunction.
func (in *CPUFunction) DeepCopy() *CPUFunction {
	if in == nil {
		return nil
	}
	out := new(CPUFunction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CPUFunction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUFunctionList) DeepCopyInto(out *CPUFunctionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CPUFunction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUFunctionList.
func (in *CPUFunctionList) DeepCopy() *CPUFunctionList {
	if in == nil {
		return nil
	}
	out := new(CPUFunctionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CPUFunctionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUFunctionSpec) DeepCopyInto(out *CPUFunctionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.AcceleratorIDs != nil {
		in, out := &in.AcceleratorIDs, &out.AcceleratorIDs
		*out = make([]AccIDInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]EnvsInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RequestMemorySize != nil {
		in, out := &in.RequestMemorySize, &out.RequestMemorySize
		*out = new(int32)
		**out = **in
	}
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(string)
		**out = **in
	}
	if in.PreviousFunctions != nil {
		in, out := &in.PreviousFunctions, &out.PreviousFunctions
		*out = make(map[string]FromToWBFunction, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NextFunctions != nil {
		in, out := &in.NextFunctions, &out.NextFunctions
		*out = make(map[string]FromToWBFunction, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]intstr.IntOrString, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUFunctionSpec.
func (in *CPUFunctionSpec) DeepCopy() *CPUFunctionSpec {
	if in == nil {
		return nil
	}
	out := new(CPUFunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUFunctionStatus) DeepCopyInto(out *CPUFunctionStatus) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	if in.RxProtocol != nil {
		in, out := &in.RxProtocol, &out.RxProtocol
		*out = new(string)
		**out = **in
	}
	if in.TxProtocol != nil {
		in, out := &in.TxProtocol, &out.TxProtocol
		*out = new(string)
		**out = **in
	}
	if in.AdditionalNetwork != nil {
		in, out := &in.AdditionalNetwork, &out.AdditionalNetwork
		*out = new(bool)
		**out = **in
	}
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
	in.StartTime.DeepCopyInto(&out.StartTime)
	if in.IPAddress != nil {
		in, out := &in.IPAddress, &out.IPAddress
		*out = new(string)
		**out = **in
	}
	if in.AcceleratorStatuses != nil {
		in, out := &in.AcceleratorStatuses, &out.AcceleratorStatuses
		*out = make([]AccStatusesByContainer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUFunctionStatus.
func (in *CPUFunctionStatus) DeepCopy() *CPUFunctionStatus {
	if in == nil {
		return nil
	}
	out := new(CPUFunctionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUPodSpec) DeepCopyInto(out *CPUPodSpec) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]corev1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]CPUContainer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUPodSpec.
func (in *CPUPodSpec) DeepCopy() *CPUPodSpec {
	if in == nil {
		return nil
	}
	out := new(CPUPodSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvsData) DeepCopyInto(out *EnvsData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvsData.
func (in *EnvsData) DeepCopy() *EnvsData {
	if in == nil {
		return nil
	}
	out := new(EnvsData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrameSizeData) DeepCopyInto(out *FrameSizeData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrameSizeData.
func (in *FrameSizeData) DeepCopy() *FrameSizeData {
	if in == nil {
		return nil
	}
	out := new(FrameSizeData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionStatusData) DeepCopyInto(out *FunctionStatusData) {
	*out = *in
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
	if in.FunctionKernelID != nil {
		in, out := &in.FunctionKernelID, &out.FunctionKernelID
		*out = new(int32)
		**out = **in
	}
	if in.FunctionChannelID != nil {
		in, out := &in.FunctionChannelID, &out.FunctionChannelID
		*out = new(int32)
		**out = **in
	}
	if in.PtuKernelID != nil {
		in, out := &in.PtuKernelID, &out.PtuKernelID
		*out = new(int32)
		**out = **in
	}
	if in.FrameworkKernelID != nil {
		in, out := &in.FrameworkKernelID, &out.FrameworkKernelID
		*out = new(int32)
		**out = **in
	}
	in.Rx.DeepCopyInto(&out.Rx)
	in.Tx.DeepCopyInto(&out.Tx)
	out.SharedMemory = in.SharedMemory
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionStatusData.
func (in *FunctionStatusData) DeepCopy() *FunctionStatusData {
	if in == nil {
		return nil
	}
	out := new(FunctionStatusData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodTemplate) DeepCopyInto(out *PodTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodTemplate.
func (in *PodTemplate) DeepCopy() *PodTemplate {
	if in == nil {
		return nil
	}
	out := new(PodTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RxTxData) DeepCopyInto(out *RxTxData) {
	*out = *in
	if in.IPAddress != nil {
		in, out := &in.IPAddress, &out.IPAddress
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.SubnetAddress != nil {
		in, out := &in.SubnetAddress, &out.SubnetAddress
		*out = new(string)
		**out = **in
	}
	if in.GatewayAddress != nil {
		in, out := &in.GatewayAddress, &out.GatewayAddress
		*out = new(string)
		**out = **in
	}
	if in.DMAChannelID != nil {
		in, out := &in.DMAChannelID, &out.DMAChannelID
		*out = new(int32)
		**out = **in
	}
	if in.LLDMAConnectorID != nil {
		in, out := &in.LLDMAConnectorID, &out.LLDMAConnectorID
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RxTxData.
func (in *RxTxData) DeepCopy() *RxTxData {
	if in == nil {
		return nil
	}
	out := new(RxTxData)
	in.DeepCopyInto(out)
	return out
}
