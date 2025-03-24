/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "CPUFunction/api/v1"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var t = metav1.Time{
	Time: time.Now(),
}
var testTime = metav1.Time{
	Time: t.Time.AddDate(0, 0, -1),
}

var partitionName1 string = "cpufunctest1-wbfunction-decode-main"

var CPUFunction1 = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest1-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest1",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest1-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.112",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
			"inputPort": {
				IntVal: 8556,
			},
			"ipAddress": {
				StrVal: "192.174.90.102/24",
				Type:   1,
			},
			"inputIPAddress": {
				StrVal: "192.174.90.102",
				Type:   1,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest1-wbfunction-decode-main",
			CommandQueueID:  "test01-cpufunctest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
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

var PCIeConnection1 = PCIeConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest1-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest1",
			Namespace: "default",
		},

		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest1-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest1",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest1-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest1-wbfunction-filter-resize-low-infer-main",
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

var FPGAFunction1 = FPGAFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest1-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-low-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest1",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-low-infer",
		NodeName:          "node01",
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
			Name:      "cpufunctest1",
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
		Status: "pending",
	},
}

var cpuconfigdecode = corev1.ConfigMap{
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
      "virtualNetworkDeviceDriverType": "sriov",
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

var cpuconfigfrhigh = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunc-config-filter-resize-high-infer",
		Namespace: "default",
	},
	Data: map[string]string{
		"cpufunc-config-filter-resize-high-infer.json": `
    [{
      "rxProtocol":"TCP",
      "txProtocol":"TCP",
      "additionalNetwork": true,
      "virtualNetworkDeviceDriverType": "sriov",
      "imageURI": "localhost/cpu-filterresize-app:3.1.0",
      "envs":{
        "FRENV_APPLOG_LEVEL": "INFO",
        "FRENV_INPUT_WIDTH": "3840",
        "FRENV_INPUT_HEIGHT": "2160",
        "FRENV_OUTPUT_WIDTH": "1280",
        "FRENV_OUTPUT_HEIGHT": "1280"
      },
      "template":{
        "apiVersion": "v1",
        "kind": "Pod",
        "spec":{
          "containers":[{
            "name": "fr",
            "command": ["python",
               "fr.py",
               "--in_port=$(FRENV_INPUT_PORT)",
               "--out_addr=$(FRENV_OUTPUT_IP)",
               "--out_port=$(FRENV_OUTPUT_PORT)",
               "--in_width=$(FRENV_INPUT_WIDTH)",
               "--in_height=$(FRENV_INPUT_HEIGHT)",
               "--out_width=$(FRENV_OUTPUT_WIDTH)",
               "--out_height=$(FRENV_OUTPUT_HEIGHT)",
               "--loglevel=$(FRENV_APPLOG_LEVEL)"],
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

var cpuconfigfrlow = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunc-config-filter-resize-low-infer",
		Namespace: "default",
	},
	Data: map[string]string{
		"cpufunc-config-filter-resize-low-infer.json": `
    [{
      "rxProtocol":"TCP",
      "txProtocol":"TCP",
      "additionalNetwork": true,
      "virtualNetworkDeviceDriverType": "sriov",
      "imageURI": "localhost/cpu-filterresize-app:3.1.0",
      "envs":{
        "FRENV_APPLOG_LEVEL": "INFO",
        "FRENV_INPUT_WIDTH": "3840",
        "FRENV_INPUT_HEIGHT": "2160",
        "FRENV_OUTPUT_WIDTH": "416",
        "FRENV_OUTPUT_HEIGHT": "416"
      },
      "template":{
        "apiVersion": "v1",
        "kind": "Pod",
        "spec":{
          "containers":[{
            "name": "fr",
            "command": ["python",
               "fr.py",
               "--in_port=$(FRENV_INPUT_PORT)",
               "--out_addr=$(FRENV_OUTPUT_IP)",
               "--out_port=$(FRENV_OUTPUT_PORT)",
               "--in_width=$(FRENV_INPUT_WIDTH)",
               "--in_height=$(FRENV_INPUT_HEIGHT)",
               "--out_width=$(FRENV_OUTPUT_WIDTH)",
               "--out_height=$(FRENV_OUTPUT_HEIGHT)",
               "--loglevel=$(FRENV_APPLOG_LEVEL)"],
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

var cpuconfigcopybranch = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunc-config-copy-branch",
		Namespace: "default",
	},
	Data: map[string]string{
		"cpufunc-config-copy-branch.json": `
    [{
      "rxProtocol":"TCP",
      "txProtocol":"TCP",
      "additionalNetwork": true,
      "virtualNetworkDeviceDriverType": "sriov",
      "copyMemorySize": "1024",
      "imageURI": "localhost/cpu-copybranch-app:3.1.0",
      "template":{
        "apiVersion": "v1",
        "kind": "Pod",
        "spec":{
          "containers":[{
            "name": "cfunc-copy-branch-1",
            "workingDir": "/opt/fpga-software/tools/copy_branch",
            "command": ["sh","-c"],
            "args":["./copy_branch",
               "%RECEIVING%",
               "%NUM%",
               "%FORWARDING%",
               "%MEMSIZE%"],
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

var cpuconfigglue = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunc-config-glue-fdma-to-tcp",
		Namespace: "default",
	},
	Data: map[string]string{
		"cpufunc-config-glue-fdma-to-tcp.json": `
    [{
      "rxProtocol":"DMA",
      "txProtocol":"TCP",
      "sharedMemoryMiB": 256,
      "imageURI": "localhost/cpu-glue-app:3.1.0",
      "additionalNetwork": true,
      "virtualNetworkDeviceDriverType": "sriov",
      "template":{
        "apiVersion": "v1",
        "kind": "Pod",
        "spec":{
          "containers":[{
            "name": "cfunc-glue-fdma-to-tcp-1",
            "workingDir": "/opt/fpga-software/tools/glue_fdma_tcp",
            "command": ["sh","-c"],
            "args":["./build/glue",
               "%FORWARDING%",
               "%WIDTH%",
               "%HEIGHT%"],
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
    }]`,
	},
}

var partitionName2 string = "cpufunctest2-wbfunction-decode-main"
var CPUFunction2 = examplecomv1.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest2-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName2,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest2",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest2-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 15,
			},
			"ipAddress": {
				StrVal: "192.174.90.111/24",
				Type:   1,
			},
			"inputPort": {
				IntVal: 5004,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.131",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
			"inputIPAddress": {
				StrVal: "192.174.90.111",
				Type:   1,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest2-wbfunction-decode-main",
			CommandQueueID:  "test01-cpufunctest2-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
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

var EthernetConnection2 = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest2-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest2",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest2-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest2-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest2",
			Namespace: "default",
		},
		Status: "pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest2-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest2-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var partitionName3 string = "cpufunctest3-wbfunction-filter-resize-high-infer-main"

var CPUFunction3 = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName3,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-high-infer",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest3-wbfunction-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 15,
			},
			"ipAddress": {
				StrVal: "192.168.122.50/24",
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
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest3-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest3-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "test01-cpufunctest3-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
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

var PCIeConnection3 = PCIeConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest3-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},

		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var EthernetConnectionfrhigh = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest3-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},
		Status: "pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var EthernetConnection3frhigh = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest3",
			Namespace: "default",
		},
		Status: "pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

/*
testcase 1-1-4: cpu-filter-resize-low-infer
*/

var EthernetConnectionfrlow = EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest4-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest4",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest4",
			Namespace: "default",
		},
		Status: "Pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

var EthernetConnection4frlow = EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest4-wbconnection-filter-resize-low-infer-main-low-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest4",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest4",
			Namespace: "default",
		},
		Status: "Pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest4-wbfunction-low-infer-main",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

var partitionName4 string = "cpufunctest4-wbfunction-filter-resize-low-infer-main"

var CPUFunction4frlow = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName4,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-low-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest4",
			Namespace: "default",
		},
		DeviceType: "cpu",
		Envs: []examplecomv1.EnvsInfo{
			{
				PartitionName: "cpufunction4",
				EachEnv: []examplecomv1.EnvsData{
					{
						EnvKey:   "test",
						EnvValue: "testvalue",
					},
				},
			},
		},
		FunctionName: "cpu-filter-resize-low-infer",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest4-wbfunction-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"ipAddress": {
				StrVal: "192.168.122.150/24",
				Type:   1,
			},
			"inputPort": {
				IntVal: 15000,
			},
			"outputIPAddress": {
				StrVal: "192.168.122.121",
				Type:   1,
			},
			"outputPort": {
				IntVal: 16000,
			},
			"inputIPAddress": {
				StrVal: "192.168.122.150",
				Type:   1,
			},
		},
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest4-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest4-wbfunction-filter-resize-low-infer-main",
			CommandQueueID:  "test01-cpufunctest4-wbfunction-filter-resize-low-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
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
testcase 1-1-5: copybranch
*/
var reqMemSize int32 = 32
var partitionName5 string = "cpufunctest5-wbfunction-copy-branch-main"

var CPUFunction5copy = examplecomv1.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest5-wbfunction-copy-branch-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName5,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-copy-branch",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "copy-branch",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest5-wbfunction-infer-1",
					Namespace: "default",
				},
			},
			"1": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest5-wbfunction-infer-2",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"inputIPAddress": {
				StrVal: "192.168.122.121",
				Type:   1,
			},
			"inputPort": {
				IntVal: 16000,
			},
			"branchOutputIPAddress": {
				StrVal: "192.168.90.141,192.168.90.142",
				Type:   1,
			},
			"branchOutputPort": {
				StrVal: "17000,18000",
				Type:   1,
			},
			"ipAddress": {
				StrVal: "192.168.122.121/24",
				Type:   1,
			},
		},
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest5-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest5-wbfunction-copy-branch-main",
			CommandQueueID:  "test01-cpufunctest5-wbfunction-copy-branch-main",
			SharedMemoryMiB: 0,
		},
		RegionName:        "cpu",
		RequestMemorySize: &reqMemSize,
	},
	Status: examplecomv1.CPUFunctionStatus{
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

var EthernetConnection4 = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest5-wbconnection-filter-resize-low-infer-main-copy-branch-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		Status: "Pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

var EthernetConnection5 = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest5-wbconnection-copy-branch-main-infer-1",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-infer-1",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		Status: "Pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-infer-1",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

var EthernetConnection6 = EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest5-wbconnection-copy-branch-main-infer-2",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-infer-2",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest5",
			Namespace: "default",
		},
		Status: "Pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest5-wbfunction-infer-2",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

/*
testcase 1-1-6: glue
*/
var partitionName6 string = "cpufunctest6-wbfunction-glue-fdma-to-tcp-main"

var CPUFunction6glue = examplecomv1.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName6,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-glue-fdma-to-tcp",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest6",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "glue-fdma-to-tcp",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest6-wbfunction-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 15,
			},
			"inputIPAddress": {
				StrVal: "192.168.122.131",
				Type:   1,
			},
			"inputPort": {
				IntVal: 16000,
			},
			"outputIPAddress": {
				StrVal: "192.168.122.100",
				Type:   1,
			},
			"outputPort": {
				IntVal: 16000,
			},
			"glueOutputIPAddress": {
				StrVal: "192.174.90.141",
				Type:   1,
			},
			"glueOutputPort": {
				StrVal: "16000",
				Type:   1,
			},
			"ipAddress": {
				StrVal: "192.174.122.131/24",
				Type:   1,
			},
		},
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest6-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
			CommandQueueID:  "test01-cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
			SharedMemoryMiB: 256,
		},
		RegionName:        "cpu",
		RequestMemorySize: &reqMemSize,
	},
	Status: examplecomv1.CPUFunctionStatus{
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

var PCIeConnection5 = PCIeConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest6-wbconnection-filter-resize-high-infer-main-glue-fdma-to-tcp-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest6",
			Namespace: "default",
		},

		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest6",
			Namespace: "default",
		},
		Status: "Pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

var FPGAFunction5 = FPGAFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest6-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest6",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "node01",
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
			Name:      "cpufunctest6",
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

var EthernetConnection6glue = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest6-wbconnection-glue-fdma-to-tcp-main-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest6",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest6",
			Namespace: "default",
		},
		Status: "",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				Namespace: "default",
			},
			Status: "",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest6-wbfunction-high-infer-main",
				Namespace: "default",
			},
			Status: "",
		},
	},
}

var NetworkAttachmentDefinition1 = NetworkAttachmentDefinition{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "k8s.cni.cncf.io",
		Kind:       "NetworkAttachmentDefinition",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node01-config-net-sriov",
		Namespace: "default",
		Annotations: map[string]string{
			"k8s.v1.cni.cncf.io/resourceName": "intel.com/intel_sriov_netdevice",
		},
	},
	Spec: NetworkAttachmentDefinitionSpec{
		Config: `{
			"type": "sriov",
			"cniVersion": "0.3.1",
			"name": "node01-config-net-sriov",
			"ipam": {
				"type": "static"
				}
			}`,
	},
}

/*
testcase 1-2
*/
var EthernetConnection1 = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctiontest3-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctiontest3-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		Status: "pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctiontest3-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var EthernetConnection3 = EthernetConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: EthernetConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		From: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: EthernetFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctiontest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		Status: "pending",
		From: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: EthernetFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctiontest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var partitionName12 string = "dftest-wbfunction-filter-resize-high-infer-main"

var CPUFunction12 = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "dftest-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName12,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-high-infer",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctiontest3-wbfunction-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
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
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctiontest3-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-dftest-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "test01-dftest-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var functionIndex122 int32 = 99
var partitionName122 string = "dftest-wbfunction-filter-resize-high-infer-main"
var CPUFunction122 = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "dftest-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName122,
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctiontest3",
			Namespace: "default",
		},
		DeviceType:    "cpu",
		FunctionName:  "cpu-filter-resize-high-infer",
		FunctionIndex: &functionIndex122,
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctiontest3-wbfunction-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
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
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctiontest3-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-dftest122-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "test01-dftest122-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctiontest122",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

/*
testcase 2-1 Update
*/
var partitionNameUpdate string = "cpufunctestupdate-wbfunction-decode-main"

var cpufunctestUPDATE = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctestupdate-wbfunction-decode-main",
		Namespace: "default",
		Finalizers: []string{
			"cpufunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionNameUpdate,
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctestupdate",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctestupdate-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.112",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
			"inputPort": {
				IntVal: 8556,
			},
			"ipAddress": {
				StrVal: "192.174.90.102/24",
				Type:   1,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctestupdate-wbfunction-decode-main",
			CommandQueueID:  "test01-cpufunctestupdate-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctestupdate",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "localhost/host_decode:3.1.0",
		ConfigName:   "cpufunc-config-decode",
		Status:       "Running",
	},
}

/*
testcase 2-2 Delete
*/
var partitionNameDelete string = "cpufunctestupdate-wbfunction-decode-main"
var functionIndexDelete int32 = 6

var cpufunctestDELETE = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctestdelete-wbfunction-decode-main",
		Namespace: "default",
		Finalizers: []string{
			"cpufunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionNameDelete,
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctestupdate",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctestupdate-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.112",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
			"inputPort": {
				IntVal: 8556,
			},
			"ipAddress": {
				StrVal: "192.174.90.102/24",
				Type:   1,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctestupdate-wbfunction-decode-main",
			CommandQueueID:  "test01-cpufunctestupdate-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctestupdate",
			Namespace: "default",
		},
		FunctionName:  "cpu-decode",
		FunctionIndex: &functionIndexDelete,
		ImageURI:      "localhost/host_decode:3.1.0",
		ConfigName:    "cpufunc-config-decode",
		Status:        "Running",
	},
}

/*
testcase 2-1-4
*/

var cpuconfigdecode214 = corev1.ConfigMap{
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
          "lifecycle":{
            "preStop":{
              "exec":{
                "command": ["sh","-c", "kill -TERM $(pidof cpu_decode-shared)"]}}},
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
      "virtualNetworkDeviceDriverType": "sriov",
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

var PCIeConnection214 = PCIeConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest214-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest214",
			Namespace: "default",
		},

		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest214-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest214-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest214",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest214-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "cpufunctest214-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var FPGAFunction214 = FPGAFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest214-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-low-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "cpufunctest214",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-low-infer",
		NodeName:          "node01",
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
			Name:      "cpufunctest214",
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
		Status: "pending",
	},
}

var NetworkAttachmentDefinition214 = NetworkAttachmentDefinition{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "k8s.cni.cncf.io",
		Kind:       "NetworkAttachmentDefinition",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node01-config-net-sriov",
		Namespace: "default",
		Annotations: map[string]string{
			"k8s.v1.cni.cncf.io/resourceName": "intel.com/intel_sriov_netdevice",
		},
	},
	Spec: NetworkAttachmentDefinitionSpec{
		Config: `{
			"type": "sriov",
			"cniVersion": "0.3.1",
			"name": "node01-config-net-sriov",
			"ipam": {
				"type": "static"
				}
			}`,
	},
}

var CPUFunction214 = examplecomv1.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctest214-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.CPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "cpufunctest214",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest214-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "node01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.112",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
			"inputPort": {
				IntVal: 8556,
			},
			"ipAddress": {
				StrVal: "192.174.90.102/24",
				Type:   1,
			},
			"inputIPAddress": {
				StrVal: "192.174.90.102",
				Type:   1,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctest1-wbfunction-decode-main",
			CommandQueueID:  "test01-cpufunctest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: examplecomv1.CPUFunctionStatus{
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

/*
	registration of PCIeConnection
*/

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

func init() {
	SchemeBuilderpcie.Register(&PCIeConnection{}, &PCIeConnectionList{})
}

// groupversion_info.go
var (
	// GroupVersion is group version used to register these objects
	GroupVersionpcie = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilderpcie = &scheme.Builder{GroupVersion: GroupVersionpcie}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemepcie = SchemeBuilderpcie.AddToScheme
)

// zz_generated.deepcopy.go
// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccIDInfo) DeepCopyInto(out *AccIDInfo) {
	*out = *in
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
func (in *FunctionStatusData) DeepCopyInto(out *FunctionStatusData) {
	*out = *in
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
func (in *RxTxSpecFunc) DeepCopyInto(out *RxTxSpecFunc) {
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

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RxTxSpecFunc.
func (in *RxTxSpecFunc) DeepCopy() *RxTxSpecFunc {
	if in == nil {
		return nil
	}
	out := new(RxTxSpecFunc)
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
func (in *SharedMemoryStatus) DeepCopyInto(out *SharedMemoryStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SharedMemoryStatus.
func (in *SharedMemoryStatus) DeepCopy() *SharedMemoryStatus {
	if in == nil {
		return nil
	}
	out := new(SharedMemoryStatus)
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

// function_common_types.go
// Function CR structure
type FunctionData struct {
	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     int32                         `json:"functionIndex"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      SharedMemorySpec              `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        *string                       `json:"configName,omitempty"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
	FunctionKernelID  *int32                        `json:"functionKernelID,omitempty"`
	FunctionChannelID *int32                        `json:"functionChannelID,omitempty"`
	PtuKernelID       *int32                        `json:"ptuKernelID,omitempty"`
	FrameworkKernelID *int32                        `json:"frameworkKernelID,omitempty"`
	Rx                RxTxSpecFunc                  `json:"rx,omitempty"`
	Tx                RxTxSpecFunc                  `json:"tx,omitempty"`
}

// Function CR structure
type FunctionStatusData struct {
	Status string `json:"status"`
}
type AccIDInfo struct {
	PartitionName string `json:"partitionName"`
	ID            string `json:"id"`
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

type FPGARxTx struct {
	Protocol  string `json:"protocol"`
	StartPort int32  `json:"startPort,omitempty"`
	EndPort   int32  `json:"endPort,omitempty"`
}

type FunctionConfigMap struct {
	RxProtocol      string            `json:"rxProtocol,omitempty"`
	TxProtocol      string            `json:"txProtocol,omitempty"`
	SharedMemoryMiB int32             `json:"sharedMemoryMiB,omitempty"`
	ImageURI        string            `json:"imageURI,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	ParentBitStream string            `json:"parentBitStream,omitempty"`
	ChildBitStream  string            `json:"childBitStream,omitempty"`
}

// function_common_types.go
type WBNamespacedName struct {

	//+kubebuilder:validation:Required

	Namespace string `json:"namespace"`

	//+kubebuilder:validation:Required

	Name string `json:"name"`
}

type FromToWBFunction struct {

	//+kubebuilder:validation:Required

	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`

	//+kubebuilder:validation:Required

	Port int32 `json:"port"`
}

/*
	registration of EthernetConnection
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
	SchemeBuilderpcie.Register(&EthernetConnection{}, &EthernetConnectionList{})
}

// groupversion_info.go
var (
	// GroupVersion is group version used to register these objects
	GroupVersioneth = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuildereth = &scheme.Builder{GroupVersion: GroupVersioneth}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemeeth = SchemeBuildereth.AddToScheme
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

/*
	registration of FPGAFucntion
*/

// fpgafunction_types.go
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
