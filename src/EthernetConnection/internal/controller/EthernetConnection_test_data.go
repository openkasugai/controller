/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "EthernetConnection/api/v1"
	controllertestcpu "EthernetConnection/internal/controller/test/type/CPU"
	controllertestfpga "EthernetConnection/internal/controller/test/type/FPGA"
	controllertestgpu "EthernetConnection/internal/controller/test/type/GPU"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var podGPU = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest3-wbfunction-high-infer-main-mps-dgpu-gpu--pod",
		Namespace: "default",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "testgpu",
				Image: "testgpu",
			},
		},
	},
}

var podCPU = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest3-mps-dgpu-node01-cpu0-pod",
		Namespace: "default",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "testcpu",
				Image: "testcpu",
			},
		},
	},
}

// indirect variables difinitions
var frameworkKernelID int32 = 0
var functionChannelID int32 = 0
var functionIndex int32 = 0
var functionKernelID int32 = 0
var ptuKernelID int32 = 0
var fdmaconnectorID int32 = 512
var dmachannelID int32 = 0
var port1 int32 = 1111
var port2 int32 = 2222
var ip string = "111.111.111.111"

var FPGAFunctiondecode = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest1-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T01J",
			},
		},
		ConfigName: "fpgafunc-config-decode",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
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
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
			Port:            &port2,
			IPAddress:       &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "RTP",
			Port:      &port1,
			IPAddress: &ip,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest1-wbfunction-decode-main",
			CommandQueueID:  "ethernetconnectiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},
		FunctionName:        "decode",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver2_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
			Port:            &port2,
			IPAddress:       &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "RTP",
			Port:      &port1,
			IPAddress: &ip,
		},
		Status: "Running",
	},
}

var port3 int32 = 3333
var port4 int32 = 4444
var ip2 string = "222.222.222.222"

var FPGAFunctionfilter = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
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
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
			Port:            &port4,
			IPAddress:       &ip2,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:        "DMA",
			FDMAConnectorID: &fdmaconnectorID,
			DMAChannelID:    &dmachannelID,
			Port:            &port4,
			IPAddress:       &ip2,
		},
		Status: "Running",
	},
}

var EthernetConnection1 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest1-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
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
		Name:      "ethernetconnectiontest3-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "ethernetconnectiontest3-wbfunction-decode-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "ethernetconnectiontest3-wbfunction-filter-resize-high-infer-main",
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
			FilePrefix:      "test01-ethernetconnectiontest3-wbfunction-decode-main",
			CommandQueueID:  "test01-ethernetconnectiontest3-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
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
		Name:      "ethernetconnectiontest2-wbfunction-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-low-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
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
			Protocol: "DMA",
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-low-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol: "TCP",
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol: "DMA",
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

var EthernetConnection2 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest2-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest2-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
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
		Name:      "ethernetconnectiontest3-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
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
			FilePrefix:      "ethernetconnectiontest3-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest3-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol: "TCP",
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol: "DMA",
		},
		Status: "Running",
	},
}

var GPUFunctionhighinfer = controllertestgpu.GPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest3-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				ID:            "GPU-",
				PartitionName: "ethernetconnectiontest3-wbfunction-high-infer-main",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		NodeName:     "node01",
		RegionName:   "a100",
		SharedMemory: &controllertestgpu.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest3-wbfunction-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest3-wbfunction-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestgpu.GPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
			Namespace: "default",
		},
		FunctionName: "high-infer",
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
		Name:      "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-high-infer",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "ethernetconnectiontest4-wbfunction-high-infer-main",
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
					Name:      "ethernetconnectiontest4-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var EthernetConnection3 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest3-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest3",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest3-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest3-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var EthernetConnection4 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest4-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest4-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var EthernetConnection5 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest5-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest5",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest5-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest5-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var EthernetConnection6 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest6-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest6",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest6-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest6-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var ChildBitstreamID string = "aaaaaa"
var ChildBitstreamCRName string = "cccccc"

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
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &ChildBitstreamID,
			ChildBitstreamCRName: &ChildBitstreamCRName,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:  &ChildBitstreamID,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}

var FPGA2 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test03",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &ChildBitstreamID,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &ChildBitstreamID,
			ChildBitstreamCRName: &ChildBitstreamCRName,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test04",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &ChildBitstreamID,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}

var FPGAFunctiondecode_tcp = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest2-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T01J",
			},
		},
		ConfigName: "fpgafunc-config-decode",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
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
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest2-wbfunction-decode-main",
			CommandQueueID:  "ethernetconnectiontest2-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
			Namespace: "default",
		},
		FunctionName:        "decode",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		Status: "Running",
	},
}

var FPGAFunctionfilter_tcp = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
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
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest2",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		Status: "Running",
	},
}

var FPGAFunctiondecode_tcp2 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest4-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T01J",
			},
		},
		ConfigName: "fpgafunc-config-decode",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
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
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest4-wbfunction-decode-main",
			CommandQueueID:  "ethernetconnectiontest4-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
			Namespace: "default",
		},
		FunctionName:        "decode",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		Status: "Running",
	},
}

var FPGAFunctionfilter_tcp2 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
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
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest4-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest4",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		Status: "Running",
	},
}

var FPGAFunctiondecode_tcp3 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest5-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T01J",
			},
		},
		ConfigName: "fpgafunc-config-decode",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest5",
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
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest5-wbfunction-decode-main",
			CommandQueueID:  "ethernetconnectiontest5-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest5",
			Namespace: "default",
		},
		FunctionName:        "decode",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		Status: "Running",
	},
}

var FPGAFunctiondecode_tcp4 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest6-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T01J",
			},
		},
		ConfigName: "fpgafunc-config-decode",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest6",
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
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest6-wbfunction-decode-main",
			CommandQueueID:  "ethernetconnectiontest6-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest6",
			Namespace: "default",
		},
		FunctionName:        "decode",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port2,
			IPAddress: &ip,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port1,
			IPAddress: &ip,
		},
		Status: "Running",
	},
}

var FPGAFunctionfilter_tcp3 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest5-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest5",
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
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest5-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest5-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest5",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		Status: "Running",
	},
}

var FPGAFunctionfilter_tcp4 = controllertestfpga.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest6-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestfpga.FPGAFunctionSpec{
		AcceleratorIDs: []controllertestfpga.AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xethernet_21330621T00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest6",
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
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		SharedMemory: &controllertestfpga.SharedMemorySpec{
			FilePrefix:      "ethernetconnectiontest6-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetconnectiontest6-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: controllertestfpga.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestfpga.WBNamespacedName{
			Name:      "ethernetconnectiontest6",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tethernet_tandem1.mcs",
		ChildBitstreamName:  "ver1_tethernet_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port3,
			IPAddress: &ip2,
		},
		Tx: controllertestfpga.RxTxSpec{
			Protocol:  "TCP",
			Port:      &port4,
			IPAddress: &ip2,
		},
		Status: "Running",
	},
}

var childbsid string = "111111111"
var maxFunctions int32 = 1
var maxCapacity int32 = 2
var name string = "child1"
var cids string = "111"
var id int32 = 3
var identifier string = "child1_identifier"
var typ string = "childbs_chaintype"
var varsion string = "childbs_varsion1.1.3"
var maxDataflows int32 = 4
var available bool = true
var funcCRName string = "funcCRName"
var port5 int32 = 5
var dmaChannel int32 = 6
var lldmaConnector int32 = 7
var port6 int32 = 8
var dmaChannel2 int32 = 9
var lldmaConnector2 int32 = 10
var uid types.UID = "aaaaaaa"

var Childbs = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "Childbs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testchildbs",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga3",
				UID:        uid,
			},
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "111.111.111",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "222.222.222",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "333.333.333",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "444.444.444",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		ChildBitstreamID: &childbsid,
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "555.555.555",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "666.666.666",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							Parameters: &map[string]intstr.IntOrString{
								"5": {
									StrVal: "param01",
									IntVal: 12345,
									Type:   1,
								},
							},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		Status:           examplecomv1.ChildBsStatusPreparing,
		State:            examplecomv1.ChildBsNoConfigureNetwork,
		ChildBitstreamID: &childbsid,
	},
}

var Childbs2 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "Childbs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testchildbs2",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga4",
				UID:        uid,
			},
		},
		Finalizers: []string{
			"deviceinfo.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "111.111.111",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "222.222.222",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "333.333.333",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "444.444.444",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		ChildBitstreamID: &childbsid,
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "555.555.555",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "666.666.666",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							Parameters: &map[string]intstr.IntOrString{
								"5": {
									StrVal: "param01",
									IntVal: 12345,
									Type:   1,
								},
							},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		Status:           examplecomv1.ChildBsStatusPreparing,
		State:            examplecomv1.ChildBsNoConfigureNetwork,
		ChildBitstreamID: &childbsid,
	},
}

var Childbs3 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "Childbs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testchildbs3",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga5",
				UID:        uid,
			},
		},
		Finalizers: []string{
			"deviceinfo.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "111.111.111",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "222.222.222",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "333.333.333",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "444.444.444",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		ChildBitstreamID: &childbsid,
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "555.555.555",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "666.666.666",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							Parameters: &map[string]intstr.IntOrString{
								"5": {
									StrVal: "param01",
									IntVal: 12345,
									Type:   1,
								},
							},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		Status:           examplecomv1.ChildBsStatusPreparing,
		State:            examplecomv1.ChildBsNoConfigureNetwork,
		ChildBitstreamID: &childbsid,
	},
}

var Childbs4 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "Childbs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testchildbs4",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga6",
				UID:        uid,
			},
		},
		Finalizers: []string{
			"deviceinfo.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "111.111.111",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "222.222.222",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "333.333.333",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "444.444.444",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		ChildBitstreamID: &childbsid,
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids,
						ID:   &id,
						Parameters: &map[string]intstr.IntOrString{
							"IPAddress": {
								StrVal: "555.555.555",
								Type:   1,
							},
							"SubnetAddress": {
								StrVal: "666.666.666",
								Type:   1,
							},
							"GatewayAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
							"MacAddress": {
								StrVal: "777.777.777",
								Type:   1,
							},
						},
					},
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids,
						ID:   &id,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id,
						Identifier: &identifier,
						Type:       &typ,
						Version:    &varsion,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier,
							Type:       &typ,
							Version:    &varsion,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier,
								Type:       &typ,
								Version:    &varsion,
							}},
							Parameters: &map[string]intstr.IntOrString{
								"5": {
									StrVal: "param01",
									IntVal: 12345,
									Type:   1,
								},
							},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available,
									FunctionCRName: &funcCRName,
									Rx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5,
												DMAChannelID:     &dmaChannel,
												LLDMAConnectorID: &lldmaConnector,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port6,
												DMAChannelID:     &dmaChannel2,
												LLDMAConnectorID: &lldmaConnector2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity,
								MaxDataFlows: &maxDataflows,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions,
				MaxCapacity:  &maxCapacity,
				Name:         &name,
			},
		},
		Status:           examplecomv1.ChildBsStatusPreparing,
		State:            examplecomv1.ChildBsNoConfigureNetwork,
		ChildBitstreamID: &childbsid,
	},
}

var FPGA3 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpga3",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "childbs03",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}

var FPGA4 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpga4",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "childbs04",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}

var FPGA5 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpga5",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testchildbs03",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}

var FPGA6 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpga6",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testchildbs04",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			ChildBitstreamID:     &childbsid,
			ChildBitstreamCRName: &name,
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}

/////////////////////////////////////////////////////////

var EthernetConnection118 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetcontest118-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetcontest118",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontest118-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontest118-filter-resize-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var CPUFunctionDecode118 = controllertestcpu.CPUFunction{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetcontest118-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "ethernetcontest118-decode-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetcontest118",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "ethernetcontest118-filter-resize-main",
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
			FilePrefix:      "test01-ethernetcontest118-wbfunction-decode-main",
			CommandQueueID:  "test01-ethernetcontest118-wbfunction-decode-main",
			SharedMemoryMiB: 256,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetcontest118",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var CPUFunctionFilterResize118 = controllertestcpu.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		Kind:       "CPUFunction",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetcontest118-filter-resize-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "ethernetcontest118-filter-resize-main",
				ID:            "node01-cpu0",
			},
		},
		ConfigName: "cpufunc-config-filter-resize-high-infer",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetcontest118",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-filter-resize-high-infer",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "ethernetcontest118-high-infer-main",
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
					Name:      "ethernetcontest118-decode-main",
					Namespace: "default",
				},
			},
		},
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "ethernetcontest118-filter-resize-high-infer-main",
			CommandQueueID:  "ethernetcontest118-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "ethernetcontest118",
			Namespace: "default",
		},
		FunctionName: "cpu-filter-resize",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "Running",
	},
}

var CPUfilterresizePOD = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetcontest118-filter-resize-main-cpu-pod",
		Namespace: "default",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Env: []corev1.EnvVar{
					{
						Name:  "testEnv",
						Value: "test",
					},
				},
				Image: "image: localhost/cpu_filterresize:4.0.2",
				Name:  "cpu-filter-resize0",
			},
		},
	},
}

var ethernetConnectionUPDATE = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetcontestupdate-wbconnection-main",
		Namespace: "default",
		Finalizers: []string{
			"ethernetconnection.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetcontestupdate",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestupdate-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestupdate-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetcontestupdate",
			Namespace: "default",
		},
		Status: "Pending",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestupdate-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "Pending",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestupdate-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "Pending",
		},
	},
}

var ethernetConnectionDELETE = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetcontestdelete-wbconnection-main",
		Namespace: "default",
		Finalizers: []string{
			"ethernetconnection.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetcontestdelete",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestdelete-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestdelete-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetcontestdelete",
			Namespace: "default",
		},
		Status: "OK",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestdelete-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetcontestdelete-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var EthernetConnection621 = examplecomv1.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "ethernetconnectiontest621-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.EthernetConnectionSpec{
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},

		From: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.EthernetFunctionSpec{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},
		Status: "Running",
		From: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: examplecomv1.EthernetFunctionStatus{
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}

var FPGA621 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test621",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &ChildBitstreamID,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
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
			ChildBitstreamID:     &ChildBitstreamID,
			ChildBitstreamCRName: &ChildBitstreamCRName,
			DeviceIndex:          0,
			DeviceFilePath:       "/dev/xethernet_21330621T01J",
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
	}, {
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test0621",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &ChildBitstreamID,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
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
			DeviceIndex:          1,
			DeviceFilePath:       "/dev/xethernet_21330621T00D",
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
	},
}
