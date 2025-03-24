/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "FPGAFunction/api/v1"
	controllertestcpu "FPGAFunction/internal/controller/test/type/CPU"
	controllertestethernet "FPGAFunction/internal/controller/test/type/Ethernet"
	controllertestgpu "FPGAFunction/internal/controller/test/type/GPU"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var t = metav1.Time{
	Time: time.Now(),
}
var testTime = metav1.Time{
	Time: t.Time.AddDate(0, 0, -1),
}

var frameworkKernelID1 int32 = 0
var functionChannelID1 int32 = 0

var functionKernelID1 int32 = 0
var ptuKernelID1 int32 = 0
var partitionName1 string = "0"

var FPGAFunction1 = examplecomv1.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.FPGAFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType: "alveo",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctiontest1-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-high-infer-main",
					Namespace: "default",
				},
			},
		},
		FrameworkKernelID: &frameworkKernelID1,
		FunctionChannelID: &functionChannelID1,
		FunctionKernelID:  &functionKernelID1,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "test01",
		PtuKernelID:       &ptuKernelID1,
		RegionName:        "lane0",
	},
	Status: examplecomv1.FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "fpgafunctiontest",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: examplecomv1.RxTxData{
			Protocol: "TCP",
		},
		Tx: examplecomv1.RxTxData{
			Protocol: "DMA",
		},
		Status: "pending",
	},
}

var CPUFunction1 = controllertestcpu.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctiontest1-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "cpufunctiontest1-wbfunction-decode-main",
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "cpufunctiontest1-wbfunction-filter-resize-low-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "test01",
		Params: map[string]intstr.IntOrString{
			"decEnvFrameFPS": {
				IntVal: 5,
			},
			"inputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"inputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &controllertestcpu.SharedMemorySpec{
			FilePrefix:      "test01-cpufunctiontest1-wbfunction-decode-main",
			CommandQueueID:  "test01-cpufunctiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var functionIndexg1 int32 = 0

var GPUFunction1 = controllertestgpu.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				PartitionName: "df-night01-wbfunction-high-infer-main",
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:    "a100",
		FunctionIndex: &functionIndexg1,
		FunctionName:  "high-infer",
		PreviousFunctions: map[string]controllertestgpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestgpu.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "test01",
		Params: map[string]intstr.IntOrString{
			"outputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"outputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &controllertestgpu.SharedMemorySpec{
			FilePrefix:      "df-night01-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night01-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "a100",
	},
	Status: controllertestgpu.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "gpufunctiontest",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var childbsid1 string = "00000001"

var name1 string = "child1"

var FPGA1 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpga-21320621v00d-test01",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid1,
			DeviceIndex:       0,
			DeviceFilePath:    "/dev/xethernet_21330621T01J",
			DeviceUUID:        "21330621T01J",
			NodeName:          "test01",
			ParentBitstreamID: "bbbbbbbbb",
			PCIDomain:         3,
			PCIBus:            4,
			PCIDevice:         5,
			PCIFunction:       6,
			Vendor:            "zzzzzzvendor",
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
			ChildBitstreamID:  &childbsid1,
			DeviceIndex:       1,
			DeviceFilePath:    "/dev/xethernet_21330621T00D",
			DeviceUUID:        "21330621T00D",
			NodeName:          "test01",
			ParentBitstreamID: "bbbbbbbbb",
			PCIDomain:         3,
			PCIBus:            4,
			PCIDevice:         5,
			PCIFunction:       6,
			Vendor:            "zzzzzzvendor",
		},
	},
}

var FPGA1Status = []examplecomv1.FPGAStatus{
	{
		ChildBitstreamID:     &childbsid1,
		ChildBitstreamCRName: &name1,
		DeviceIndex:          0,
		DeviceFilePath:       "/dev/xethernet_21330621T01J",
		DeviceUUID:           "21330621T01J",
		NodeName:             "test01",
		ParentBitstreamID:    "bbbbbbbbb",
		PCIDomain:            3,
		PCIBus:               4,
		PCIDevice:            5,
		PCIFunction:          6,
		Vendor:               "zzzzzzvendor",
		Status:               examplecomv1.FPGAStatusPreparing,
	}, {
		ChildBitstreamID:     &childbsid1,
		ChildBitstreamCRName: &name1,
		DeviceIndex:          1,
		DeviceFilePath:       "/dev/xethernet_21330621T00D",
		DeviceUUID:           "21330621T00D",
		NodeName:             "test01",
		ParentBitstreamID:    "bbbbbbbbb",
		PCIDomain:            3,
		PCIBus:               4,
		PCIDevice:            5,
		PCIFunction:          6,
		Vendor:               "zzzzzzvendor",
		Status:               examplecomv1.FPGAStatusPreparing,
	},
}

var fpgafuncconfig_fr_high_infer_1 = corev1.ConfigMap{
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
                "id": "00000001"
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
            "functionDedicatedInfo": "filter-resize-ch",
            "functionName": "filter-resize-high-infer"
        }`,
	},
}

var servicerMgmtConfig1 = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "servicer-mgmt-info",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"servicer-mgmt-info.json": `
		[{
			"nodeName":"test01",
			"networkInfo":[
				{
					"deviceFilePath": "/dev/xpcie_21330621T04L",
					"laneIndex":0,
					"ipAddress":"192.174.90.91",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5D:A0"
				},
				{
					"deviceFilePath": "/dev/ xpcie_21330621T04L",
					"laneIndex":1,
					"ipAddress":"192.174.90.92",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5D:A1"
				},
				{
					"deviceFilePath": "/dev/xpcie_21330621T01J",
					"laneIndex":0,
					"ipAddress":"192.174.90.81",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5C:A1"
				},
				{
					"deviceFilePath": "/dev/ xpcie_21330621T01J",
					"laneIndex":1,
					"ipAddress":"192.174.90.82",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5C:A2"
				},
				{
					"deviceFilePath": "/dev/xethernet_21330621T01J",
					"laneIndex":0,
					"ipAddress":"192.174.90.81",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5C:A1"
				},
				{
					"deviceFilePath": "/dev/xethernet_21330621T01J",
					"laneIndex":1,
					"ipAddress":"192.174.90.82",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5C:A2"
				},
				{
					"deviceFilePath": "/dev/xpcie_21330621T00Y",
					"laneIndex":0,
					"ipAddress":"192.174.90.83",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5B:A0"
				},
				{
					"deviceFilePath": "/dev/ xpcie_21330621T00Y",
					"laneIndex":1,
					"ipAddress":"192.174.90.84",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5B:A1"
				},
				{
					"deviceFilePath": "/dev/xpcie_21330621T00D",
					"laneIndex":0,
					"ipAddress":"192.174.90.93",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5A:A1"
				},
				{
					"deviceFilePath": "/dev/ xpcie_21330621T00D",
					"laneIndex":1,
					"ipAddress":"192.174.90.94",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5A:A2"
				},
				{
					"deviceFilePath": "/dev/xpcie_21320621V00D",
					"laneIndex":0,
					"ipAddress":"192.174.90.93",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5A:A1"
				},
				{
					"deviceFilePath": "/dev/xpcie_21320621V00D",
					"laneIndex":1,
					"ipAddress":"192.174.90.94",
					"subnetAddress":"255.255.255.0",
					"gatewayAddress":"192.174.90.1",
					"macAddress":"00:12:34:00:5A:A2"
				}
			]
		}]`,
	},
}

var deployinfo_configdata1 = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deployinfo",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"deployinfo.json": `
		{"devices": [{
			"nodeName": "test01",
			"deviceFilePath": "/dev/xpcie_21330621T04L",
			"deviceUUID": "21330621T04L",
			"functionTargets": [{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane0",
				"maxFunctions": 8,
				"maxCapacity": 40
				},{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane1",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		},{
			"nodeName": "test01",
			"deviceFilePath": "/dev/xpcie_21330621T01J",
			"deviceUUID": "21330621T01J",
			"functionTargets": [{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane0",
				"maxFunctions": 8,
				"maxCapacity": 40
				},{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane1",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		},{
			"nodeName": "test01",
			"deviceFilePath": "/dev/xpcie_21320621V00D",
			"deviceUUID": "21320621V00D",
			"functionTargets": [{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane0",
				"maxFunctions": 4,
				"maxCapacity": 20
				},{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane1",
				"maxFunctions": 3,
				"maxCapacity": 30
			}]
		},{
			"nodeName": "test01",
			"deviceFilePath": "/dev/xpcie_21320621V01D",
			"deviceUUID": "21320621V01D",
			"functionTargets": [{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane0",
				"maxFunctions": 4,
				"maxCapacity": 20
				},{
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"regionName": "lane1",
				"maxFunctions": 3,
				"maxCapacity": 30
			}]
		},{
			"nodeName": "test01",
			"deviceUUID": "gpu-123456789t4",
			"functionTargets": [{
				"regionType": "gpu",
				"regionName": "t4",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		},{
			"nodeName": "test01",
			"deviceUUID": "gpu-123456789a100",
			"functionTargets": [{
				"regionType": "gpu",
				"regionName": "a100",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		},{
			"nodeName": "test01",
			"deviceUUID": "test01-cpu",
			"functionTargets": [{
				"regionType": "cpu",
				"regionName": "cpu",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		}]}`,
	},
}

var regionUniqueInfoConfig1 = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "region-unique-info",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"region-unique-info.json": `
		[{
			"subDeviceSpecRef": "00000001",
			"functionTargets":[{
				"regionName": "lane0",
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"maxFunctions": 2,
				"maxCapacity": 40,
				"funcitons":[{
					"functionIndex": 0,
					"partitionName": "0",
					"functionName": "filter-resize-high-infer",
					"maxDataFlows": 8,
					"maxCapacity": 40
				}]
			},{
				"regionName": "lane1",
				"regionType": "alveou250-0100001c-2lanes-1nics",
				"maxFunctions": 2,
				"maxCapacity": 40
			}]
		},{
			"subDeviceSpecRef": "Tesla T4",
			"functionTargets":[{
				"regionName": "t4",
				"regionType": "t4",
				"maxFunctions": 110,
				"maxCapacity": 40
			}]
		},{
			"subDeviceSpecRef": "NVIDIA A100 80GB PCIe",
			"functionTargets":[{
				"regionName": "a100",
				"regionType": "a100",
				"maxFunctions": 110,
				"maxCapacity": 120
			}]
		},{
			"subDeviceSpecRef": "Intel(R) Xeon(R) Gold 6346 CPU @ 3.10GHz",
			"functionTargets":[{
				"regionName": "cpu",
				"regionType": "cpu",
				"maxFunctions": 110,
				"maxCapacity": 120
			}]
		},{
			"subDeviceSpecRef": "Intel(R) Xeon(R) Gold 6348 CPU @ 2.60GHz",
			"functionTargets":[{
				"regionName": "cpu",
				"regionType": "cpu",
				"maxFunctions": 110,
				"maxCapacity": 120
			}]
		},{
			"subDeviceSpecRef": "Intel(R) Xeon(R) Gold 6330 CPU @ 2.00GHz",
			"functionTargets":[{
				"regionName": "cpu",
				"regionType": "cpu",
				"maxFunctions": 110,
				"maxCapacity": 120
			}]
		}]`,
	},
}

var functionUniqueInfoConfig1 = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "function-unique-info",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"function-unique-info.json": `
		[{
			"functionID" : 0,
			"functionName" : "decode",
			"maxDataFlows": 6,
			"maxCapacity": 20
		},{
			"functionID" : 0,
			"functionName" : "filter-resize",
			"maxDataFlows": 8,
			"maxCapacity": 40
		}]`,
	},
}

var fr_childbs_Config1 = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "filter-resize-ch",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"filter-resize-childbs.json": `
		{
			"functionKernels":[{
				"partitionName": "0",
				"functionChannelIDList": [0,1,2,3,4,5,6,7],
				"functionChannelIDs":[{
					"functionChannelID": 0,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12300
							},
							"DMA":{
								"port": 12300,
								"lldmaConnectorID": 512,
								"dmaChannelID": 0
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12300
							},
							"DMA":{
								"port": 12300,
								"lldmaConnectorID": 512,
								"dmaChannelID": 0
							}
						}
					}
				},{
					"functionChannelID": 1,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12301
							},
							"DMA":{
								"port": 12301,
								"lldmaConnectorID": 513,
								"dmaChannelID": 0
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12301
							},
							"DMA":{
								"port": 12301,
								"lldmaConnectorID": 513,
								"dmaChannelID": 1
							}
						}
					}
				},{
					"functionChannelID": 2,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12302
							},
							"DMA":{
								"port": 12302,
								"lldmaConnectorID": 514,
								"dmaChannelID": 2
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12302
							},
							"DMA":{
								"port": 12302,
								"lldmaConnectorID": 514,
								"dmaChannelID": 2
							}
						}
					}
				},{
					"functionChannelID": 3,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12303
							},
							"DMA":{
								"port": 12303,
								"lldmaConnectorID": 515,
								"dmaChannelID": 3
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12303
							},
							"DMA":{
								"port": 12303,
								"lldmaConnectorID": 515,
								"dmaChannelID": 3
							}
						}
					}
				},{
					"functionChannelID": 4,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12304
							},
							"DMA":{
								"port": 12304,
								"lldmaConnectorID": 516,
								"dmaChannelID": 4
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12304
							},
							"DMA":{
								"port": 12304,
								"lldmaConnectorID": 516,
								"dmaChannelID": 4
							}
						}
					}
				},{
					"functionChannelID": 5,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12305
							},
							"DMA":{
								"port": 12305,
								"lldmaConnectorID": 517,
								"dmaChannelID": 5
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12305
							},
							"DMA":{
								"port": 12305,
								"lldmaConnectorID": 517,
								"dmaChannelID": 5
							}
						}
					}
				},{
					"functionChannelID": 6,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12306
							},
							"DMA":{
								"port": 12306,
								"lldmaConnectorID": 518,
								"dmaChannelID": 6
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12306
							},
							"DMA":{
								"port": 12306,
								"lldmaConnectorID": 518,
								"dmaChannelID": 6
							}
						}
					}
				},{
					"functionChannelID": 7,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12307
							},
							"DMA":{
								"port": 12307,
								"lldmaConnectorID": 519,
								"dmaChannelID": 7
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12307
							},
							"DMA":{
								"port": 12307,
								"lldmaConnectorID": 519,
								"dmaChannelID": 7
							}
						}
					}
				}]
			},{
				"partitionName": "1",
				"functionChannelIDList": [8,9,10,11,12,13,14,15],
				"functionChannelIDs":[{
					"functionChannelID": 8,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12308
							},
							"DMA":{
								"port": 12308,
								"lldmaConnectorID": 520,
								"dmaChannelID": 8
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12308
							},
							"DMA":{
								"port": 12308,
								"lldmaConnectorID": 520,
								"dmaChannelID": 8
							}
						}
					}
				},{
					"functionChannelID": 9,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12309
							},
							"DMA":{
								"port": 12309,
								"lldmaConnectorID": 521,
								"dmaChannelID": 9
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12309
							},
							"DMA":{
								"port": 12309,
								"lldmaConnectorID": 521,
								"dmaChannelID": 9
							}
						}
					}
				},{
					"functionChannelID": 10,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12310
							},
							"DMA":{
								"port": 12310,
								"lldmaConnectorID": 522,
								"dmaChannelID": 10
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12310
							},
							"DMA":{
								"port": 12310,
								"lldmaConnectorID": 522,
								"dmaChannelID": 10
							}
						}
					}
				},{
					"functionChannelID": 11,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12311
							},
							"DMA":{
								"port": 12311,
								"lldmaConnectorID": 523,
								"dmaChannelID": 11
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12311
							},
							"DMA":{
								"port": 12311,
								"lldmaConnectorID": 523,
								"dmaChannelID": 11
							}
						}
					}
				},{
					"functionChannelID": 12,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12312
							},
							"DMA":{
								"port": 12312,
								"lldmaConnectorID": 524,
								"dmaChannelID": 12
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12312
							},
							"DMA":{
								"port": 12312,
								"lldmaConnectorID": 524,
								"dmaChannelID": 12
							}
						}
					}
				},{
					"functionChannelID": 13,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12313
							},
							"DMA":{
								"port": 12313,
								"lldmaConnectorID": 525,
								"dmaChannelID": 13
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12313
							},
							"DMA":{
								"port": 12313,
								"lldmaConnectorID": 525,
								"dmaChannelID": 13
							}
						}
					}
				},{
					"functionChannelID": 14,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12314
							},
							"DMA":{
								"port": 12314,
								"lldmaConnectorID": 526,
								"dmaChannelID": 14
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12314
							},
							"DMA":{
								"port": 12314,
								"lldmaConnectorID": 526,
								"dmaChannelID": 14
							}
						}
					}
				},{
					"functionChannelID": 15,
					"rx":{
						"protocol":{
							"TCP":{
								"port": 12315
							},
							"DMA":{
								"port": 12315,
								"lldmaConnectorID": 527,
								"dmaChannelID": 15
							}
						}
					},
					"tx":{
						"protocol":{
							"TCP":{
								"port": 12315
							},
							"DMA":{
								"port": 12315,
								"lldmaConnectorID": 527,
								"dmaChannelID": 15
							}
						}
					}
				}
			]}]
		}`,
	},
}

var EthernetConnection1 = controllertestethernet.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestethernet.EthernetConnectionSpec{
		DataFlowRef: controllertestethernet.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},

		From: controllertestethernet.EthernetFunctionSpec{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: controllertestethernet.EthernetFunctionSpec{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestethernet.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestethernet.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: controllertestethernet.EthernetFunctionStatus{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: controllertestethernet.EthernetFunctionStatus{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var EthernetConnection2 = controllertestethernet.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestethernet.EthernetConnectionSpec{
		DataFlowRef: controllertestethernet.WBNamespacedName{
			Name:      "ethernetconnectiontest1",
			Namespace: "default",
		},

		From: controllertestethernet.EthernetFunctionSpec{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: controllertestethernet.EthernetFunctionSpec{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestethernet.EthernetConnectionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestethernet.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: controllertestethernet.EthernetFunctionStatus{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: controllertestethernet.EthernetFunctionStatus{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var CurrentCapacity int32 = 0
var CurrentCapacity2 int32 = 0
var CurrentCapacity3 int32 = 0
var CurrentCapacity4 int32 = 0
var CurrentFunctions int32 = 2
var CurrentFunctions2 int32 = 1
var CurrentFunctions3 int32 = 0
var CurrentFunctions4 int32 = 8
var CurrentDataFlows int32 = 1
var CurrentDataFlows2 int32 = 2
var MaxCapacity int32 = 40
var MaxCapacity2 int32 = 40
var MaxCapacity3 int32 = 120
var MaxCapacity4 int32 = 240
var MaxDataFlows1 int32 = 6
var MaxDataFlows2 int32 = 8
var DeviceUUID string = "21320621V00D"
var DeviceUUID2 string = "21330621T04L"
var DeviceUUID3 string = "GPU-b8b4f1f5-bf51-eaa3-6ec4-97190b7f6c98"
var DeviceUUID4 string = "GPU-5b771964-ab74-a674-15d7-8f0d2cee4ef8"
var DeviceUUID5 string = "swb-sm7-cpu0"
var MaxFunctions int32 = 2
var MaxFunctions2 int32 = 1
var MaxFunctions3 int32 = 110
var MaxFunctions4 int32 = 220

var comres1 = examplecomv1.ComputeResource{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "compute-test01",
		Namespace: "default",
	},
	Spec: examplecomv1.ComputeResourceSpec{
		NodeName: "test01",
		Regions: []examplecomv1.RegionInfo{
			{
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "0",
				}, {
					Available:        true,
					CurrentCapacity:  &CurrentCapacity2,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "1",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxCapacity2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity2,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "2",
				}, {
					Available:        true,
					CurrentCapacity:  &CurrentCapacity2,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "3",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxFunctions,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity,
					CurrentDataFlows: &CurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &MaxCapacity2,
					MaxDataFlows:     &MaxDataFlows2,
					PartitionName:    "0",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxFunctions2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity,
					CurrentDataFlows: &CurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &MaxCapacity2,
					MaxDataFlows:     &MaxDataFlows2,
					PartitionName:    "1",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxDataFlows2,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity3,
				CurrentFunctions: &CurrentFunctions3,
				DeviceFilePath:   "",
				DeviceIndex:      0,
				DeviceType:       "t4",
				DeviceUUID:       &DeviceUUID3,
				MaxCapacity:      &MaxCapacity2,
				MaxFunctions:     &MaxFunctions3,
				Name:             "t4",
				Type:             "t4",
			}, {
				Available:        false,
				CurrentCapacity:  &CurrentCapacity4,
				CurrentFunctions: &CurrentFunctions4,
				DeviceFilePath:   "",
				DeviceIndex:      1,
				DeviceType:       "a100",
				DeviceUUID:       &DeviceUUID4,
				MaxCapacity:      &MaxCapacity3,
				MaxFunctions:     &MaxFunctions3,
				Name:             "a100",
				Type:             "a100",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity3,
				CurrentFunctions: &CurrentFunctions3,
				DeviceFilePath:   "",
				DeviceIndex:      0,
				DeviceType:       "cpu",
				DeviceUUID:       &DeviceUUID5,
				MaxCapacity:      &MaxCapacity4,
				MaxFunctions:     &MaxFunctions4,
				Name:             "cpu",
				Type:             "cpu",
			},
		},
	},
	Status: examplecomv1.ComputeResourceStatus{
		NodeName: "test01",
		Regions: []examplecomv1.RegionInfo{
			{
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "0",
				}, {
					Available:        true,
					CurrentCapacity:  &CurrentCapacity2,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "1",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxCapacity2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity2,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "2",
				}, {
					Available:        true,
					CurrentCapacity:  &CurrentCapacity2,
					CurrentDataFlows: &CurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &MaxCapacity,
					MaxDataFlows:     &MaxDataFlows1,
					PartitionName:    "3",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxFunctions,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity,
					CurrentDataFlows: &CurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &MaxCapacity2,
					MaxDataFlows:     &MaxDataFlows2,
					PartitionName:    "0",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxFunctions2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity,
				CurrentFunctions: &CurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &DeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &CurrentCapacity,
					CurrentDataFlows: &CurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &MaxCapacity2,
					MaxDataFlows:     &MaxDataFlows2,
					PartitionName:    "1",
				}},
				MaxCapacity:  &MaxCapacity2,
				MaxFunctions: &MaxDataFlows2,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity3,
				CurrentFunctions: &CurrentFunctions3,
				DeviceFilePath:   "",
				DeviceIndex:      0,
				DeviceType:       "t4",
				DeviceUUID:       &DeviceUUID3,
				MaxCapacity:      &MaxCapacity2,
				MaxFunctions:     &MaxFunctions3,
				Name:             "t4",
				Type:             "t4",
			}, {
				Available:        false,
				CurrentCapacity:  &CurrentCapacity4,
				CurrentFunctions: &CurrentFunctions4,
				DeviceFilePath:   "",
				DeviceIndex:      1,
				DeviceType:       "a100",
				DeviceUUID:       &DeviceUUID4,
				MaxCapacity:      &MaxCapacity3,
				MaxFunctions:     &MaxFunctions3,
				Name:             "a100",
				Type:             "a100",
			}, {
				Available:        true,
				CurrentCapacity:  &CurrentCapacity3,
				CurrentFunctions: &CurrentFunctions3,
				DeviceFilePath:   "",
				DeviceIndex:      0,
				DeviceType:       "cpu",
				DeviceUUID:       &DeviceUUID5,
				MaxCapacity:      &MaxCapacity4,
				MaxFunctions:     &MaxFunctions4,
				Name:             "cpu",
				Type:             "cpu",
			},
		},
	},
}
