/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "PCIeConnection/api/v1"
	controllertestcpu "PCIeConnection/internal/controller/test/type/CPU"
	controllertestfpga "PCIeConnection/internal/controller/test/type/FPGA"
	controllertestgpu "PCIeConnection/internal/controller/test/type/GPU"
	"time"

	//	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	// "sigs.k8s.io/controller-runtime/pkg/scheme"
)

// fpgalist-ph3 config
/*
var fpgalist = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fpgalist-ph3",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"fpgalist-ph3.json": `
    [
      {
        "nodeName": "node01",
        "deviceFilePaths": [
        "/dev/xpcie_21330621T04L","/dev/xpcie_21330621T01J","/dev/xpcie_21330621T00Y","/dev/xpcie_21330621T00D"
        ],
        "networkInfo": [
              {
            "deviceIndex": 0,
            "laneIndex": 0,
            "functionType": "filter/resize",
            "ipAddress": "192.174.90.91",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5D:A0",
            "rx":{
              "protocol":"TCP"
            },
            "tx":{
              "protocol":"DMA"
            }
          },
          {
            "deviceIndex": 0,
            "laneIndex": 1,
            "functionType": "filter/resize",
            "ipAddress": "192.174.90.92",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5D:A1",
            "rx":{
              "protocol":"TCP"
            },
            "tx":{
              "protocol":"DMA"
            }
          },
                {
            "deviceIndex": 1,
            "laneIndex": 0,
            "functionType": "decode",
            "ipAddress": "192.174.90.81",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5C:A1",
            "rx":{
              "protocol":"RTP",
              "startPort":5004,
              "endPort":5027
            },
            "tx":{
              "protocol":"TCP"
            }
          },
          {
            "deviceIndex": 1,
            "laneIndex": 1,
            "functionType": "decode",
            "ipAddress": "192.174.90.82",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5C:A2",
            "rx":{
              "protocol":"RTP",
              "startPort":5004,
              "endPort":5027
            },
            "tx":{
              "protocol":"TCP"
            }
          },
          {
            "deviceIndex": 2,
            "laneIndex": 0,
            "functionType": "decode",
            "ipAddress": "192.174.90.83",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5B:A0",
            "rx":{
              "protocol":"RTP",
              "startPort":5004,
              "endPort":5027
            },
            "tx":{
              "protocol":"TCP"
            }
          },
          {
            "deviceIndex": 2,
            "laneIndex": 1,
            "functionType": "decode",
            "ipAddress": "192.174.90.84",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5B:A1",
            "rx":{
              "protocol":"RTP",
              "startPort":5004,
              "endPort":5027
            },
            "tx":{
              "protocol":"TCP"
            }
          },
          {
            "deviceIndex": 3,
            "laneIndex": 0,
            "functionType": "filter/resize",
            "ipAddress": "192.174.90.93",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5A:A1",
            "rx":{
              "protocol":"TCP"
            },
            "tx":{
              "protocol":"DMA"
            }
          },
          {
            "deviceIndex": 3,
            "laneIndex": 1,
            "functionType": "filter/resize",
            "ipAddress": "192.174.90.94",
            "subnetAddress": "255.255.255.0",
            "gatewayAddress": "192.174.90.1",
            "macAddress": "00:12:34:00:5A:A2",
            "rx":{
              "protocol":"TCP"
            },
            "tx":{
              "protocol":"DMA"
            }
          }
        ]
      }
    ]`,
	},
}
*/
// indirect variables difinitions
var frameworkKernelID int32 = 0
var functionChannelID int32 = 0
var functionIndex int32 = 0
var functionKernelID int32 = 0
var ptuKernelID int32 = 0

//	var sharedmem = controllertestfpga.SharedMemorySpec{
//		FilePrefix:      "pcieconnectiontest1-wbfunction-decode-main",
//		CommandQueueID:  "pcieconnectiontest1-wbfunction-decode-main",
//		SharedMemoryMiB: 1,
//	}
var fdmaconnectorID int32 = 512
var dmachannelID int32 = 0

var FPGAFunctiondecode = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest1-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T01J",
			},
		},
		ConfigName: "fpgafunc-config-decode",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest1",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "decode",
		NodeName:          "node01",
		PtuKernelID:       &ptuKernelID,
		RegionName:        "lane0",
		Rx: controllertestfpga.RxTxSpec{
			Protocol: "RTP",
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest1-wbfunction-decode-main",
			CommandQueueID:  "pcieconnectiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest1",
			Namespace: "default",
		},
		FunctionName:        "decode",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver2_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol: "RTP",
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Status: "Running",
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest1-wbfunction-decode-main",
			CommandQueueID:  "pcieconnectiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
}

//	var sharedmemfilter = controllertestfpga.SharedMemorySpec{
//		FilePrefix:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
//		CommandQueueID:  "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
//		SharedMemoryMiB: 1,
//	}
var FPGAFunctionfilter = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest1",
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
		Rx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest1",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:        "TCP",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Status: "Running",
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
}

var PCIeConnection1 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest1-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest1",
			Namespace: "default",
		},

		From: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var CPUFunctiondecode = controllertestcpu.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest2-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "pcieconnectiontest2-wbfunction-decode-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "pcieconnectiontest2",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "pcieconnectiontest2-wbfunction-filter-resize-high-infer-main",
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
				IntVal: 5004,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.111",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
		},
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "test01-pcieconnectiontest2-wbfunction-decode-main",
			CommandQueueID:  "test01-pcieconnectiontest2-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "pcieconnectiontest2",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var FPGAFunctionfilter2 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest2-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-low-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest2",
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
		Rx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest2-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "pcieconnectiontest2-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest2",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-low-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Status: "Running",
	},
}

// var t, _ = time.Parse("2006-01-02T15:04:05Z", "2023-12-01T10:00:00Z")
var t = metav1.Time{
	Time: time.Now(),
}
var testTime = metav1.Time{
	Time: t.Time.AddDate(0, 0, -1),
}

var PCIeConnection2 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest2",
			Namespace: "default",
		},

		From: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest2-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest2-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var FPGAFunctionfilter3 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest3-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest3",
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
		Rx: controllertestfpga.RxTxSpec{
			Protocol: "TCP",
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest3-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "pcieconnectiontest3-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "pcieconnectiontest3",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol: "TCP",
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
		},
		Status: "Running",
	},
}

var GPUFunctionhighinfer = controllertestgpu.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest3-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				ID:            "GPU-",
				PartitionName: "pcieconnectiontest3-wbfunction-high-infer-main",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "pcieconnectiontest3",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		NodeName:     "node01",
		RegionName:   "a100",
		SharedMemory: &controllertestgpu.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest3-wbfunction-high-infer-main",
			CommandQueueID:  "pcieconnectiontest3-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestgpu.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "pcieconnectiontest3",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}
var PCIeConnection3 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest3",
			Namespace: "default",
		},

		From: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var CPUFunctionDecode = controllertestcpu.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest4-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "pcieconnectiontest4-wbfunction-decode-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
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
				IntVal: 5004,
			},
			"outputIPAddress": {
				StrVal: "192.168.90.111",
				Type:   1,
			},
			"outputPort": {
				IntVal: 15000,
			},
		},
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "test01-pcieconnectiontest4-wbfunction-decode-main",
			CommandQueueID:  "test01-pcieconnectiontest4-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var CPUFunctionFilterResize = controllertestcpu.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-high-infer",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "pcieconnectiontest4-wbfunction-high-infer-main",
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
		PreviousFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "pcieconnectiontest4-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var PCIeConnection4 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest4-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},

		From: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest4-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var pcieconnectiontestUPDATE = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontestupdate-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontestupdate",
			Namespace: "default",
		},

		From: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestupdate-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestupdate-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.PCIeConnectionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontestupdate",
			Namespace: "default",
		},
		Status: "Pending",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestupdate-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "Pending",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestupdate-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "Pending",
		},
	},
}

var pcieconnectiontestDELETE = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontestdelete-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontestdelete",
			Namespace: "default",
		},

		From: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestdelete-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.PCIeFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestdelete-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.PCIeConnectionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontestdelete",
			Namespace: "default",
		},
		Status: "OK",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestdelete-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontestdelete-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var ChildBitstreamID string = "aaaaaa"
var ChildBitstreamID2 string = "cccccc"
var ChildBitstreamCRName string = "ddddddd"
var ChildBitstreamCRName2 string = "ddddddd"

var FPGA1 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test01",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &ChildBitstreamID,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xpcie_21330621T00D",
			DeviceUUID:        "21330621T00D",
			NodeName:          "node01",
			ParentBitstreamID: "bbbbbbbbb",
			PCIDomain:         3,
			PCIBus:            4,
			PCIDevice:         5,
			PCIFunction:       6,
			Vendor:            "zzzzzzvendor",
		},
		Status: examplecomv1.FPGAStatus{
			ChildBitstreamID:     &ChildBitstreamID,
			ChildBitstreamCRName: &ChildBitstreamCRName,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xpcie_21330621T00D",
			DeviceUUID:           "21330621T00D",
			NodeName:             "node01",
			ParentBitstreamID:    "bbbbbbbbb",
			PCIDomain:            3,
			PCIBus:               4,
			PCIDevice:            5,
			PCIFunction:          6,
			Vendor:               "zzzzzzvendor",
			Status:               examplecomv1.FPGAStatusPreparing,
		},
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test02",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &ChildBitstreamID2,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xpcie_21330621T01J",
			DeviceUUID:        "21330621T01J",
			NodeName:          "node01",
			ParentBitstreamID: "bbbbbbbbb",
			PCIDomain:         3,
			PCIBus:            4,
			PCIDevice:         5,
			PCIFunction:       6,
			Vendor:            "zzzzzzvendor",
		},
		Status: examplecomv1.FPGAStatus{
			ChildBitstreamID:     &ChildBitstreamID2,
			ChildBitstreamCRName: &ChildBitstreamCRName2,
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xpcie_21330621T01J",
			DeviceUUID:           "21330621T01J",
			NodeName:             "node01",
			ParentBitstreamID:    "bbbbbbbbb",
			PCIDomain:            3,
			PCIBus:               4,
			PCIDevice:            5,
			PCIFunction:          6,
			Vendor:               "zzzzzzvendor",
			Status:               examplecomv1.FPGAStatusPreparing,
		},
	},
}

// var CPUFunction4 = examplecomv1.CPUFunction{
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:      "cpufunctiontest4-wbfunction-copy-branch-main",
// 		Namespace: "test01",
// 	},
// 	Spec: examplecomv1.CPUFunctionSpec{
// 		AcceleratorIDs: []examplecomv1.AccIDInfo{
// 			{
// 				PartitionName: "cpufunctiontest4-wbfunction-copy-branch-main",
// 				ID:            "node01-cpu0",
// 			},
// 		},
// 		ConfigName: "cpufunc-config-copy-branch",
// 		DataFlowRef: examplecomv1.WBNamespacedName{
// 			Name:      "cpufunctiontest4",
// 			Namespace: "test01",
// 		},
// 		DeviceType: "cpu",
// 		Envs: []examplecomv1.EnvsInfo{
// 			{
// 				PartitionName: "cpufunction4",
// 				EachEnv: []examplecomv1.EnvsData{
// 					{
// 						EnvKey:   "test",
// 						EnvValue: "testvalue",
// 					},
// 				},
// 			},
// 		},
// 		FunctionName: "copy-branch",
// 		NextFunctions: map[string]examplecomv1.FromToWBFunction{
// 			"0": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest4-wbfunction-infer-1",
// 					Namespace: "test01",
// 				},
// 			},
// 			"1": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest4-wbfunction-infer-2",
// 					Namespace: "test01",
// 				},
// 			},
// 		},
// 		NodeName: "node01",
// 		Params: map[string]intstr.IntOrString{
// 			"decEnvFrameFPS": {
// 				IntVal: 5,
// 			},
// 			"inputIPAddress": {
// 				StrVal: "192.168.122.121",
// 				Type:   1,
// 			},
// 			"inputPort": {
// 				IntVal: 16000,
// 			},
// 			"branchOutputIPAddress": {
// 				StrVal: "192.174.90.141,192.174.90.142",
// 				Type:   1,
// 			},
// 			"branchOutputPort": {
// 				StrVal: "17000,18000",
// 				Type:   1,
// 			},
// 			"ipAddress": {
// 				StrVal: "192.174.122.121/24",
// 				Type:   1,
// 			},
// 		},
// 		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
// 			"0": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest4-wbfunction-filter-resize-low-infer-main",
// 					Namespace: "test01",
// 				},
// 			},
// 		},
// 		SharedMemory: &examplecomv1.SharedMemorySpec{
// 			FilePrefix:      "test01-cpufunctiontest4-wbfunction-filter-resize-high-infer-main",
// 			CommandQueueID:  "test01-cpufunctiontest4-wbfunction-filter-resize-high-infer-main",
// 			SharedMemoryMiB: 0,
// 		},
// 		RegionName:        "cpu",
// 		RequestMemorySize: &reqMemSize,
// 	},
// 	Status: examplecomv1.CPUFunctionStatus{
// 		StartTime: metav1.Now(),
// 		DataFlowRef: examplecomv1.WBNamespacedName{
// 			Name:      "cpufunctiontest4",
// 			Namespace: "test01",
// 		},
// 		FunctionName: "copy-branch",
// 		ImageURI:     "container",
// 		ConfigName:   "configname",
// 		Status:       "pending",
// 	},
// }
// var CPUFunction4frlow = examplecomv1.CPUFunction{
// 	TypeMeta: metav1.TypeMeta{
// 		Kind:       "CPUFunction",
// 		APIVersion: "example.com/v1",
// 	},
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:      "cpufunctiontest4-wbfunction-filter-resize-low-infer-main",
// 		Namespace: "test01",
// 	},
// 	Spec: examplecomv1.CPUFunctionSpec{
// 		AcceleratorIDs: []examplecomv1.AccIDInfo{
// 			{
// 				PartitionName: "cpufunctiontest4-wbfunction-filter-resize-low-infer-main",
// 				ID:            "node01-cpu0",
// 			},
// 		},
// 		ConfigName: "cpufunc-config-filter-resize-low-infer",
// 		DataFlowRef: examplecomv1.WBNamespacedName{
// 			Name:      "cpufunctiontest4",
// 			Namespace: "test01",
// 		},
// 		DeviceType:   "cpu",
// 		FunctionName: "cpu-filter-resize-low-infer",
// 		NextFunctions: map[string]examplecomv1.FromToWBFunction{
// 			"0": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest4-wbfunction-low-infer-main",
// 					Namespace: "test01",
// 				},
// 			},
// 		},
// 		NodeName: "node01",
// 		Params: map[string]intstr.IntOrString{
// 			"decEnvFrameFPS": {
// 				IntVal: 5,
// 			},
// 			"inputIPAddress": {
// 				StrVal: "192.168.122.50",
// 				Type:   1,
// 			},
// 			"inputPort": {
// 				IntVal: 15000,
// 			},
// 			"outputIPAddress": {
// 				StrVal: "192.168.122.121",
// 				Type:   1,
// 			},
// 			"outputPort": {
// 				IntVal: 16000,
// 			},
// 		},
// 		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
// 			"0": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest4-wbfunction-decode-main",
// 					Namespace: "test01",
// 				},
// 			},
// 		},
// 		SharedMemory: &examplecomv1.SharedMemorySpec{
// 			FilePrefix:      "test01-cpufunctiontest4-wbfunction-filter-resize-low-infer-main",
// 			CommandQueueID:  "test01-cpufunctiontest4-wbfunction-filter-resize-low-infer-main",
// 			SharedMemoryMiB: 1,
// 		},
// 		RegionName: "cpu",
// 	},
// 	Status: examplecomv1.CPUFunctionStatus{
// 		StartTime: metav1.Now(),
// 		DataFlowRef: examplecomv1.WBNamespacedName{
// 			Name:      "cpufunctiontest4",
// 			Namespace: "test01",
// 		},
// 		FunctionName: "cpu-filter-resize",
// 		ImageURI:     "container",
// 		ConfigName:   "configname",
// 		Status:       "pending",
// 	},
// }

// var CPUFunction5 = examplecomv1.CPUFunction{
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:      "cpufunctiontest5-wbfunction-glue-fdma-to-tcp-main",
// 		Namespace: "test01",
// 	},
// 	Spec: examplecomv1.CPUFunctionSpec{
// 		AcceleratorIDs: []examplecomv1.AccIDInfo{
// 			{
// 				PartitionName: "cpufunctiontest5-wbfunction-glue-fdma-to-tcp-main",
// 				ID:            "node01-cpu0",
// 			},
// 		},
// 		ConfigName: "cpufunc-config-glue-fdma-to-tcp",
// 		DataFlowRef: examplecomv1.WBNamespacedName{
// 			Name:      "cpufunctiontest5",
// 			Namespace: "test01",
// 		},
// 		DeviceType:   "cpu",
// 		FunctionName: "glue-fdma-to-tcp",
// 		NextFunctions: map[string]examplecomv1.FromToWBFunction{
// 			"0": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest5-wbfunction-high-infer-main",
// 					Namespace: "test01",
// 				},
// 			},
// 		},
// 		NodeName: "node01",
// 		Params: map[string]intstr.IntOrString{
// 			"decEnvFrameFPS": {
// 				IntVal: 15,
// 			},
// 			"inputIPAddress": {
// 				StrVal: "192.168.122.121",
// 				Type:   1,
// 			},
// 			"inputPort": {
// 				IntVal: 16000,
// 			},
// 			"outputIPAddress": {
// 				StrVal: "192.168.122.100",
// 				Type:   1,
// 			},
// 			"outputPort": {
// 				IntVal: 16000,
// 			},
// 			"glueOutputIPAddress": {
// 				StrVal: "192.174.90.141",
// 				Type:   1,
// 			},
// 			"glueOutputPort": {
// 				StrVal: "16000",
// 				Type:   1,
// 			},
// 			"ipAddress": {
// 				StrVal: "192.174.122.131/24",
// 				Type:   1,
// 			},
// 		},
// 		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
// 			"0": {
// 				Port: 0,
// 				WBFunctionRef: examplecomv1.WBNamespacedName{
// 					Name:      "cpufunctiontest5-wbfunction-filter-resize-high-infer-main",
// 					Namespace: "test01",
// 				},
// 			},
// 		},
// 		SharedMemory: &examplecomv1.SharedMemorySpec{
// 			FilePrefix:      "test01-cpufunctiontest5-wbfunction-glue-fdma-to-tcp-main",
// 			CommandQueueID:  "test01-cpufunctiontest5-wbfunction-glue-fdma-to-tcp-main",
// 			SharedMemoryMiB: 256,
// 		},
// 		RegionName:        "cpu",
// 		RequestMemorySize: &reqMemSize,
// 	},
// 	Status: examplecomv1.CPUFunctionStatus{
// 		StartTime: metav1.Now(),
// 		DataFlowRef: examplecomv1.WBNamespacedName{
// 			Name:      "cpufunctiontest5",
// 			Namespace: "test01",
// 		},
// 		FunctionName: "glue-fdma-to-tcp",
// 		ImageURI:     "container",
// 		ConfigName:   "configname",
// 		Status:       "pending",
// 	},
// }
