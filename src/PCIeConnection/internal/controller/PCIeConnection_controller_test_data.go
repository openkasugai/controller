/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "PCIeConnection/api/v1"
	controllertestcpu "PCIeConnection/internal/controller/test/type/CPU"
	controllertestfpga "PCIeConnection/internal/controller/test/type/FPGA"
	controllertestgpu "PCIeConnection/internal/controller/test/type/GPU"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
)

// indirect variables difinitions
var frameworkKernelID int32 = 0
var functionChannelID int32 = 0
var functionIndex int32 = 0
var functionKernelID int32 = 0
var partitionName string = "0"
var ptuKernelID int32 = 0

var lldmaconnectorID int32 = 512
var dmachannelID int32 = 0
var lldmaconnectorIDNil int32
var dmachannelIDNil int32

var protocolRTP = "RTP"
var protocolTCP = "TCP"
var protocolDMA = "DMA"

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
				PartitionName: &partitionName,
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
		Rx: &controllertestfpga.RxTxData{
			Protocol:         protocolRTP,
			DMAChannelID:     &dmachannelIDNil,
			LLDMAConnectorID: &lldmaconnectorIDNil,
		},
		Tx: &controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			DMAChannelID:     &dmachannelID,
			LLDMAConnectorID: &lldmaconnectorID,
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
		FunctionIndex:       0,
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver2_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxData{
			Protocol:         protocolRTP,
			DMAChannelID:     &dmachannelIDNil,
			LLDMAConnectorID: &lldmaconnectorIDNil,
		},
		Tx: controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			DMAChannelID:     &dmachannelID,
			LLDMAConnectorID: &lldmaconnectorID,
		},
		Status: "Running",
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest1-wbfunction-decode-main",
			CommandQueueID:  "pcieconnectiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
}

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
				PartitionName: &partitionName,
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
		Rx: &controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
		},
		Tx: &controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
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
		FunctionIndex:       0,
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
		},
		Tx: controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
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

var partitionNameCPUDecode2 string = "pcieconnectiontest2-wbfunction-decode-main"
var podNameCPUDecode2 string = "pcieconnectiontest1-wbfunction-decode-main-cpu-pod"

var CPUFunctionDecode2 = controllertestcpu.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest2-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: &partitionNameCPUDecode2,
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
		PodName:      &podNameCPUDecode2,
		ConfigName:   "configname",
		Status:       "Running",
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "test01-pcieconnectiontest2-wbfunction-decode-main",
			CommandQueueID:  "test01-pcieconnectiontest2-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
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
				PartitionName: &partitionName,
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
		Rx: &controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
		},
		Tx: &controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
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
		FunctionIndex:       0,
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
		},
		Tx: controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
		},
		Status: "Running",
	},
}

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
				PartitionName: &partitionName,
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
		Rx: &controllertestfpga.RxTxData{
			Protocol:         protocolTCP,
			LLDMAConnectorID: &lldmaconnectorIDNil,
			DMAChannelID:     &dmachannelIDNil,
		},
		Tx: &controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
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
		FunctionIndex:       0,
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxData{
			Protocol:         protocolTCP,
			LLDMAConnectorID: &lldmaconnectorIDNil,
			DMAChannelID:     &dmachannelIDNil,
		},
		Tx: controllertestfpga.RxTxData{
			Protocol:         protocolDMA,
			LLDMAConnectorID: &lldmaconnectorID,
			DMAChannelID:     &dmachannelID,
		},
		Status: "Running",
	},
}

var partitionNameGPUHigh string = "pcieconnectiontest3-wbfunction-high-infer-main"
var podNameGPUhighinfer string = "pcieconnectiontest1-wbfunction-high-infer-main-gpu-pod"

var GPUFunctionhighinfer = controllertestgpu.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest3-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				ID:            "GPU-",
				PartitionName: &partitionNameGPUHigh,
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
		PodName:      &podNameGPUhighinfer,
		ConfigName:   "configname",
		Status:       "Running",
		SharedMemory: &controllertestgpu.SharedMemorySpec{
			FilePrefix:      "pcieconnectiontest3-wbfunction-high-infer-main",
			CommandQueueID:  "pcieconnectiontest3-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
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

var partitionNameCPUDecode4 string = "pcieconnectiontest4-wbfunction-decode-main"
var podNameCPUDecode4 string = "pcieconnectiontest1-wbfunction-decode-main-cpu-pod"

var CPUFunctionDecode4 = controllertestcpu.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest4-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: &partitionNameCPUDecode4,
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
		PodName:      &podNameCPUDecode4,
		ConfigName:   "configname",
		Status:       "Running",
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "test01-pcieconnectiontest4-wbfunction-decode-main",
			CommandQueueID:  "test01-pcieconnectiontest4-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
	},
}

var partitionNameCPUFR string = "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main"
var podNameCPUFilterResize string = "pcieconnectiontest2-wbfunction-decode-main-cpu-pod"

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
				PartitionName: &partitionNameCPUFR,
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
		PodName:      &podNameCPUFilterResize,
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

var PCIeConnection5 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
		Name:      "pcieconnectiontest5-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest5",
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
			Name:      "pcieconnectiontest5",
			Namespace: "default",
		},
		Status: "Running",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest2-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest2-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var PCIeConnection6 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
		Name:      "pcieconnectiontest6-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest6",
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
			Name:      "pcieconnectiontest6",
			Namespace: "default",
		},
		Status: "Terminating",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var PCIeConnection7 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
		Name:      "pcieconnectiontest7-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest7",
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
			Name:      "pcieconnectiontest7",
			Namespace: "default",
		},
		Status: "Terminating",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest2-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest2-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var PCIeConnection8 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
		Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest8",
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
			Name:      "pcieconnectiontest8",
			Namespace: "default",
		},
		Status: "Terminating",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest3-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var PCIeConnection9 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"pcieconnection.finalizers.example.com.v1",
		},
		Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.PCIeConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "pcieconnectiontest9",
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
			Name:      "pcieconnectiontest9",
			Namespace: "default",
		},
		Status: "Terminating",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest4-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var CPUPod1 = corev1.Pod{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"kubernetes",
		},
		Name:      "pcieconnectiontest1-wbfunction-decode-main-cpu-pod",
		Namespace: "default",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			0: corev1.Container{
				Image: "pcieconnectiontest1-cpu_decode",
				Name:  "pcieconnectiontest1-cpu-pod",
			},
		},
	},
}

var CPUPod2 = corev1.Pod{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"kubernetes",
		},
		Name:      "pcieconnectiontest2-wbfunction-decode-main-cpu-pod",
		Namespace: "default",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			0: corev1.Container{
				Image: "pcieconnectiontest2-cpu_decode",
				Name:  "pcieconnectiontest2-cpu-pod",
			},
		},
	},
}

var GPUPod1 = corev1.Pod{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
	ObjectMeta: metav1.ObjectMeta{
		Finalizers: []string{
			"kubernetes",
		},
		Name:      "pcieconnectiontest1-wbfunction-high-infer-main-gpu-pod",
		Namespace: "default",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			0: corev1.Container{
				Image: "pcieconnectiontest1-gpu_infer_dma",
				Name:  "pcieconnectiontest1-gpu-pod",
			},
		},
	},
}

var PCIeConnection723 = examplecomv1.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pcieconnectiontest723-wbconnection-decode-main-filter-resize-high-infer-main",
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
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},
		Status: "Running",
		From: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest4-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.PCIeFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}
