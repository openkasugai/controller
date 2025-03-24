/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "FPGAFunction/api/v1"
	controllertestcpu "FPGAFunction/internal/controller/test/type/CPU"
	controllertestgpu "FPGAFunction/internal/controller/test/type/GPU"
	controllertestpcie "FPGAFunction/internal/controller/test/type/PCIe"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var frameworkKernelID3 int32 = 2
var functionChannelID3 int32 = 2

var functionKernelID3 int32 = 2
var ptuKernelID3 int32 = 2
var partitionName3 string = "0"

var FPGAFunction3 = examplecomv1.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.FPGAFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName3,
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night03",
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
		FrameworkKernelID: &frameworkKernelID3,
		FunctionChannelID: &functionChannelID3,
		FunctionKernelID:  &functionKernelID3,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "test01",
		PtuKernelID:       &ptuKernelID3,
		RegionName:        "lane0",
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "df-night03-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night03-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
	},
	Status: examplecomv1.FPGAFunctionStatus{
		StartTime: metav1.Now(),
	},
}

var CPUFunction3 = controllertestcpu.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctiontest3-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "cpufunctiontest3-wbfunction-decode-main",
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night03",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "cpufunctiontest3-wbfunction-filter-resize-low-infer-main",
					Namespace: "test01",
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
			FilePrefix:      "cpufunctiontest3-wbfunction-decode-main",
			CommandQueueID:  "cpufunctiontest3-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night03",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname1",
		Status:       "pending",
	},
}

var functionIndexg3 int32 = 3

var GPUFunction3 = controllertestgpu.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night03-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				PartitionName: "df-night03-wbfunction-high-infer-main",
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "df-night03",
			Namespace: "default",
		},
		DeviceType:    "a100",
		FunctionIndex: &functionIndexg3,
		FunctionName:  "high-infer",
		PreviousFunctions: map[string]controllertestgpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestgpu.WBNamespacedName{
					Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
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
			FilePrefix:      "df-night03-wbfunction-high-infer-main",
			CommandQueueID:  "df-night03-wbfunction-high-infer-main",
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
		ConfigName:   "configname1",
		Status:       "pending",
	},
}

var childbsid3 string = "333333333"

var name3 string = "child3"

var FPGA3 = []examplecomv1.FPGA{
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
			ChildBitstreamID:  &childbsid3,
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
			ChildBitstreamID:  &childbsid3,
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
		Status: examplecomv1.FPGAStatus{
			ChildBitstreamID:     &childbsid3,
			ChildBitstreamCRName: &name3,
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
	},
}

var fpgafuncconfig_fr_high_infer_3 = corev1.ConfigMap{
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
            "functionDedicatedInfo": "filter-resize-ch"
        }`,
	},
}

var servicerMgmtConfig2 = corev1.ConfigMap{
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
				}
			]
		}]`,
	},
}

var deployinfo_configdata2 = corev1.ConfigMap{
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
				"regionType": "alveou250-0100001c-2lanes-0nics",
				"regionName": "lane0",
				"maxFunctions": 8,
				"maxCapacity": 40
				},{
				"regionType": "alveo",
				"regionName": "lane1",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		},{
			"nodeName": "test01",
			"deviceFilePath": "/dev/xpcie_21330621T01J",
			"deviceUUID": "21330621T01J",
			"functionTargets": [{
				"regionType": "alveou250-0100001c-2lanes-0nics",
				"regionName": "lane0",
				"maxFunctions": 8,
				"maxCapacity": 40
				},{
				"regionType": "alveo",
				"regionName": "lane1",
				"maxFunctions": 8,
				"maxCapacity": 40
			}]
		},{
			"nodeName": "test01",
			"deviceFilePath": "/dev/xpcie_21320621V00D",
			"deviceUUID": "21320621V00D",
			"functionTargets": [{
				"regionType": "alveou250-0100001c-2lanes-0nics",
				"regionName": "lane0",
				"maxFunctions": 8,
				"maxCapacity": 40
				},{
				"regionType": "alveo",
				"regionName": "lane1",
				"maxFunctions": 8,
				"maxCapacity": 40
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

var functionUniqueInfoConfig2 = corev1.ConfigMap{
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

var fr_childbs_Config2 = corev1.ConfigMap{
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

var PCIeConnection3 = controllertestpcie.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night03-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestpcie.PCIeConnectionSpec{
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "pcieconnectiontest3",
			Namespace: "default",
		},

		From: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "df-night03-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestpcie.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "pcieconnectiontest3",
			Namespace: "",
		},
		Status: "",
		From: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}

var PCIeConnection4 = controllertestpcie.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night03-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestpcie.PCIeConnectionSpec{
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "pcieconnectiontest4",
			Namespace: "default",
		},

		From: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "cpufunctiontest3-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestpcie.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "",
		From: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
		To: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
			Status: "",
		},
	},
}
