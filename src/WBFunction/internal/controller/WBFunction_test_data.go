package controller

import (
	examplecomv1 "WBFunction/api/v1"
	controllertestcpu "WBFunction/internal/controller/test/type/CPU"
	controllertestfpga "WBFunction/internal/controller/test/type/FPGA"
	controllertestgpu "WBFunction/internal/controller/test/type/GPU"
	// k8scnicncfio "github.com/k8snetworkplumbingwg/network-attachment-definition-client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// runtime "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	// "sigs.k8s.io/controller-runtime/pkg/scheme"
)

/*
var status_def_data int32 = 0
var status_def_data1 int32 = 0

	var status_def_data2 = examplecomv1.WBFunctionRequirementsInfo{
		Capacity: 15,
	}
*/
var MaxCapacity_cpu int32 = 1
var MaxDataFlows_cpu int32 = 20
var FunctionIndex_cpu *int32 = nil
var FunctionIndex_cpu1 int32 = 12
var Requirements_cpu = examplecomv1.WBFunctionRequirementsInfo{
	Capacity: 15,
}

var WBFunction_cpu_decode = examplecomv1.WBFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBFunctionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		NodeName:      "node01",
		NodeSelector:  map[string]string{},
		DeviceType:    "cpu",
		DeviceIndex:   4,
		RegionName:    "test-cpu",
		FunctionIndex: nil,
		FunctionName:  "decode",
		ConfigName:    "cpufunc-config-decode",
		Params: map[string]intstr.IntOrString{
			"ipAddress": {
				StrVal: "10.1.1.10/24",
			},
			"inputIPAddress": {
				StrVal: "10.38.119.157",
			},
			"inputPort": {
				IntVal: 5004,
			},
		},
		NextWBFunctions: map[string]examplecomv1.FromToWBFunction{
			"1234": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctiontest1-wbfunction-decode-main",
					Namespace: "default",
				},
				Port: 8080,
			},
		},
		MaxDataFlows: &MaxDataFlows_cpu,
		MaxCapacity:  &MaxCapacity_cpu,
		Requirements: &Requirements_cpu,
	},
}

var CPUFunction1 = controllertestcpu.CPUFunctionStatus{
	StartTime: metav1.Now(),
	DataFlowRef: controllertestcpu.WBNamespacedName{
		Name:      "df-night01",
		Namespace: "test01",
	},
	FunctionName:  "cpu-decode",
	FunctionIndex: &FunctionIndex_cpu1,
	ImageURI:      "localhost/host_decode:3.1.0",
	ConfigName:    "cpufunc-config-decode",
	Status:        "Running",
}

/**************************************************************************
TEST Data GPUFunctionCR Create
**************************************************************************/
/*
var status_def_data_gpu int32 = 0
var status_def_data1_gpu int32 = 0
var status_def_data2_gpu = examplecomv1.WBFunctionRequirementsInfo{
	Capacity: 15,
}
*/
var MaxCapacity_gpu int32 = 1
var MaxDataFlows_gpu int32 = 20
var FunctionIndex_gpu *int32 = nil
var FunctionIndex_gpu1 int32 = 12
var Requirements_gpu = examplecomv1.WBFunctionRequirementsInfo{
	Capacity: 15,
}

var WBFunction_gpu = examplecomv1.WBFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "WBFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBFunctionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		NodeName:     "node01",
		NodeSelector: map[string]string{},
		DeviceType:   "a100",
		DeviceIndex:  3,
		RegionName:   "gpu",
		FunctionName: "a100",
		ConfigName:   "gpufunc-config-high-infer",
		Params: map[string]intstr.IntOrString{
			"ipAddress": {
				StrVal: "10.1.1.14/24",
			},
			"outputIPAddress": {
				StrVal: "192.174.91.10",
			},
			"outputPort": {
				IntVal: 2001,
			},
		},

		PreviousWBFunctions: map[string]examplecomv1.FromToWBFunction{
			"1234": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer",
					Namespace: "default",
				},
				Port: 8081,
			},
		},
		MaxDataFlows: &MaxDataFlows_gpu,
		MaxCapacity:  &MaxCapacity_gpu,
		Requirements: &Requirements_gpu,
	},
}

var GPUFunction1 = controllertestgpu.GPUFunctionStatus{
	StartTime: metav1.Now(),
	DataFlowRef: controllertestgpu.WBNamespacedName{
		Name:      "df-night01",
		Namespace: "test01",
	},
	FunctionName:  "cpu-decode",
	FunctionIndex: &FunctionIndex_gpu1,
	ImageURI:      "localhost/gpu-deepstream-app:3.1.0",
	ConfigName:    "gpufunc-config-high-infer",
	Status:        "Running",
}

/**************************************************************************
TEST Data GPUFunctionCR Create
**************************************************************************/
/*
var status_def_data_fpga int32 = 0
var status_def_data1_fpga int32 = 0
var status_def_data2_fpga = examplecomv1.WBFunctionRequirementsInfo{
	Capacity: 15,
}
*/
var MaxCapacity_fpga int32 = 1
var MaxDataFlows_fpga int32 = 20
var FunctionIndex_fpga *int32 = nil
var FunctionIndex_fpga1 int32 = 12
var Requirements_fpga = examplecomv1.WBFunctionRequirementsInfo{
	Capacity: 15,
}

var WBFunction_fpga_fr = examplecomv1.WBFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "WBFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-filter-resize-high-infer",
		Namespace: "default",
	},
	Spec: examplecomv1.WBFunctionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		NodeName:     "k8s-worker",
		DeviceType:   "alveo",
		DeviceIndex:  0,
		RegionName:   "lane0",
		FunctionName: "filter-resize-high-infer",
		ConfigName:   "fpgafunc-config-filter-resize-high-infer",
		NextWBFunctions: map[string]examplecomv1.FromToWBFunction{
			"1234": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-high-infer-main",
					Namespace: "default",
				},
				Port: 8081,
			},
		},
		PreviousWBFunctions: map[string]examplecomv1.FromToWBFunction{
			"1234": {
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-decode-main",
					Namespace: "default",
				},
				Port: 8081,
			},
		},
		MaxDataFlows: &MaxDataFlows_fpga,
		MaxCapacity:  &MaxCapacity_fpga,
		Requirements: &Requirements_fpga,
	},
}

var FPGAFunction1 = controllertestfpga.FPGAFunctionStatus{
	StartTime: metav1.Now(),
	DataFlowRef: controllertestfpga.WBNamespacedName{
		Name:      "df-night01",
		Namespace: "test01",
	},
	FunctionName:  "filter-resize-high-infer",
	FunctionIndex: FunctionIndex_fpga1,
	// ConfigName:    "fpgafunc-config-filter-resize-high-infer",
	Status: "Running",
}

var configdata = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "functionkindmap",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"functionkindmap.json": `
		[{
			"deviceType": "alveo",
			"functionCRKind":"FPGAFunction"
		},{
			"deviceType": "t4",
			"functionCRKind":"GPUFunction"
		},{
			"deviceType": "a100",
			"functionCRKind":"GPUFunction"
		},{
			"deviceType": "cpu",
			"functionCRKind":"CPUFunction"
		}]`,
	},
}

var cpuconfig = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunc-config-decode",
		Namespace: "default",
	},
	Data: map[string]string{
		"cpufunc-config-decode.json": `
    [{
      "rxProtocol":"RTP",
      "txProtocol":"DMA",
      "sharedMemoryMiB": 256,
      "imageURI": "localhost/host_decode:3.1.0",
      "additionalNetwork": true,
      "virtualNetworkDeviceDriverType": "sriov",
      "envs":{
        "DECENV_APPLOG_LEVEL": "6",
        "DECENV_FRAME_WIDTH": "3840",
        "DECENV_FRAME_HEIGHT": "2160",
        "DECENV_VIDEO_CONNECT_LIMIT": "0",
        "DECENV_VIDEOSRC_PROTOCOL": "RTP",
        "DECENV_OUTDST_PROTOCOL": "DMA"
      },
      "template":{
        "apiVersion": "v1",
        "kind": "Pod",
        "spec":{
          "containers":[{
            "name": "cfunc-1",
            "command": ["sh","-c"],
            "args":["./tools/host_decode/build/host_decode-shared"],
            "securityContext":{
              "privileged": true
            },
            "volumeMounts":[{
              "name": "hugepage-1gi",
              "mountPath": "/dev/hugepages"
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
      "rxProtocol":"RTP",
      "txProtocol":"TCP",
      "imageURI": "localhost/host_decode:3.1.0",
      "additionalNetwork": true,
      "virtualNetworkDeviceDriverType": "",
      "envs":{
        "DECENV_APPLOG_LEVEL": "6",
        "DECENV_FRAME_WIDTH": "3840",
        "DECENV_FRAME_HEIGHT": "2160",
        "DECENV_VIDEO_CONNECT_LIMIT": "0",
        "DECENV_VIDEOSRC_PROTOCOL": "RTP",
        "DECENV_OUTDST_PROTOCOL": "TCP"
      },
      "template":{
        "apiVersion": "v1",
        "kind": "Pod",
        "spec":{
          "containers":[{
            "name": "cfunc-1",
            "command": ["sh","-c"],
            "args":["./tools/host_decode/build/host_decode-shared"],
            "securityContext":{
              "privileged": true
            }
          }],
          "hostNetwork": false,
          "hostIPC": true,
          "restartPolicy": "Always"
        }
      }
    }]`,
	},
}

/*
	var fpgaconfig = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpgafunc-config-filter-resize-high-infer",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
		},
		Data: map[string]string{
			"fpgafunc-config-filter-resize-high-infer.json": `
			{
				"parentBitstream": {
					"file": "/home/ubuntu/mcap-lib/OpenKasugai-fpga-example-design-1.0.0-1.mcs",
					"id": "0100001c"
				},
				"childBitstream": {
					"file": "/home/ubuntu/mcap-lib/OpenKasugai-fpga-example-design-1.0.0-2.bit",
					"id": "0100001c"
				},
				"sharedMemoryMiB": 256
			}`,
		},
	}
*/
var infrainfo_configdata = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "infrastructureinfo",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"infrastructureinfo.json": `
		{"devices": [{
			"deviceFilePath": "/dev/xpcie_21330621T04L",
			"nodeName": "node01",
			"deviceUUID": "21330621T04L",
			"deviceIndex": 0,
			"deviceType": "alveo"
		},{
			"deviceFilePath": "/dev/xpcie_21330621T01J",
			"nodeName": "node01",
			"deviceUUID": "21330621T01J",
			"deviceIndex": 1,
			"deviceType": "alveo"
		},{
			"deviceFilePath": "",
			"nodeName": "node01",
			"deviceUUID": "gpu-123456789t4",
			"deviceIndex": 0,
			"deviceType": "t4"
		},{
			"deviceFilePath": "",
			"nodeName": "node01",
			"deviceUUID": "gpu-123456789a100",
			"deviceIndex": 1,
			"deviceType": "a100"
		},{
			"deviceFilePath": "",
			"nodeName": "node01",
			"deviceUUID": "test01-cpu",
			"deviceIndex": 0,
			"deviceType": "cpu"
		}]}`,
	},
}

var gpuconfig = corev1.ConfigMap{
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

var fpgafuncconfig_fr_high_infer = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fpgafunc-config-filter-resize-high-infer",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"fpgafunc-config-filter-resize-high-infer.json": ` 
		{
			"parentBitstream": {
				"file": "/home/ubuntu/mcap-lib/OpenKasugai-fpga-example-design-1.0.0-1.mcs",
				"id": "0100001c"
			},
			"childBitstream": {
				"file": "/home/ubuntu/mcap-lib/OpenKasugai-fpga-example-design-1.0.0-2.bit",
				"id": "0100001c"
			},
			"parameters": {
				"functions": {
					"i_width": 3840,
					"i_height": 2160,
					"o_width": 1280,
					"o_height": 1280
		  		}
			},
			"sharedMemoryMiB": 256,
			"functionDedicatedInfo": "filter-resize-ch"
		}`,
	},
}

var DeviceInfoRetCPU = examplecomv1.DeviceInfoStatus{
	Response: examplecomv1.WBFuncResponse{
		Status:        "Deployed",
		FunctionIndex: &FunctionIndex_cpu1,
		DeviceUUID:    "CPU",
		//			},
	},
}

var DeviceInfoRetGPU = examplecomv1.DeviceInfoStatus{
	Response: examplecomv1.WBFuncResponse{
		Status:        "Deployed",
		FunctionIndex: &FunctionIndex_cpu1,
		DeviceUUID:    "GPU",
		//			},
	},
}

var DeviceInfoRetFPGA = examplecomv1.DeviceInfoStatus{
	Response: examplecomv1.WBFuncResponse{
		Status:        "Deployed",
		FunctionIndex: &FunctionIndex_cpu1,
		DeviceUUID:    "FPGA",
		//			},
	},
}
