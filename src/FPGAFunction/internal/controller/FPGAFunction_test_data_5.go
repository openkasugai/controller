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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var frameworkKernelID4 int32 = 3
var functionChannelID4 int32 = 3
var functionIndex4 int32 = 3
var functionKernelID4 int32 = 3
var ptuKernelID4 int32 = 3
var partitionName4 string = "0"

var FPGAFunction4 = examplecomv1.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night04-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.FPGAFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName4,
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night04",
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
		FrameworkKernelID: &frameworkKernelID4,
		FunctionChannelID: &functionChannelID4,
		FunctionIndex:     &functionIndex4,
		FunctionKernelID:  &functionKernelID4,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "test01",
		PtuKernelID:       &ptuKernelID4,
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
			Protocol: "DMA",
		},
		Tx: examplecomv1.RxTxData{
			Protocol: "DMA",
		},
		Status: "pending",
	},
}

var CPUFunction4 = controllertestcpu.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctiontest4-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: controllertestcpu.CPUFunctionSpec{
		AcceleratorIDs: []controllertestcpu.AccIDInfo{
			{
				PartitionName: "cpufunctiontest4-wbfunction-decode-main",
				ID:            "",
			},
		},
		ConfigName: "cpufunc-config-decode",
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night04",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "cpufunctiontest4-wbfunction-filter-resize-low-infer-main",
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
			FilePrefix:      "test04-cpufunctiontest1-wbfunction-decode-main",
			CommandQueueID:  "test04-cpufunctiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night04",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname1",
		Status:       "pending",
	},
}

var functionIndexg4 int32 = 3

var GPUFunction4 = controllertestgpu.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night04-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				PartitionName: "df-night04-wbfunction-high-infer-main",
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "df-night04",
			Namespace: "default",
		},
		DeviceType:    "a100",
		FunctionIndex: &functionIndexg4,
		FunctionName:  "high-infer",
		PreviousFunctions: map[string]controllertestgpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestgpu.WBNamespacedName{
					Name:      "df-night04-wbfunction-filter-resize-high-infer-main",
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
			FilePrefix:      "df-night04-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night04-wbfunction-filter-resize-high-infer-main",
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

var childbsid4 string = "4444444444"
var maxFunctions4 int32 = 3
var maxCapacity4 int32 = 4
var name4 string = "lane0"
var cids4 string = "444"
var id4 int32 = 0
var identifier4 string = "child4_identifier"
var typ4 string = "childbs_chaintype"
var varsion4 string = "childbs_varsion1.1.5"
var maxDataflows4 int32 = 6
var available4 bool = true
var available_false bool = false
var funcCRName4 string = "funcCRName"
var port5_1 int32 = 7
var dmaChannel5_1 int32 = 8
var lldmaConnector5_1 int32 = 9
var port5_2 int32 = 10
var dmaChannel5_2 int32 = 11
var lldmaConnector5_2 int32 = 12
var uid4 types.UID = "bbbbbbbbbbbbb"
var functionChannelIDs4_1 = "0-7"

var ChildBitstream2 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "Childbs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fpga-21320621v00d-test01-4444444444",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga3",
				UID:        uid4,
			},
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					/*
						Ptu: &examplecomv1.ChildBsPtu{
							Cids: &cids4,
							ID:   &id4,
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
					*/
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids4,
						ID:   &id4,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id4,
						Identifier: &identifier4,
						Type:       &typ4,
						Version:    &varsion4,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id4,
						Identifier: &identifier4,
						Type:       &typ4,
						Version:    &varsion4,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id4,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier4,
							Type:       &typ4,
							Version:    &varsion4,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id4,
							Module: &[]examplecomv1.FunctionsModule{{
								FunctionChannelIDs: &functionChannelIDs4_1,
								Identifier:         &identifier4,
								Type:               &typ4,
								Version:            &varsion4,
							}},
							/*Parameters: &map[string]intstr.IntOrString{
							   	"0": {
							     	StrVal: "param01",
							         IntVal: 12345,
							         Type:   1,
							     },
							  },
							*/IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"0": {
									Available:      &available_false,
									FunctionCRName: &funcCRName4,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5_1,
												DMAChannelID:     &dmaChannel5_1,
												LLDMAConnectorID: &lldmaConnector5_1,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port5_2,
												DMAChannelID:     &dmaChannel5_2,
												LLDMAConnectorID: &lldmaConnector5_2,
											},
										},
									},
								},
								"1": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5_1,
												DMAChannelID:     &dmaChannel5_1,
												LLDMAConnectorID: &lldmaConnector5_1,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port5_2,
												DMAChannelID:     &dmaChannel5_2,
												LLDMAConnectorID: &lldmaConnector5_2,
											},
										},
									},
								},
								"2": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"3": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"4": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"5": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"6": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"7": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity4,
								MaxDataFlows: &maxDataflows4,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions4,
				MaxCapacity:  &maxCapacity4,
				Name:         &name4,
			},
		},
		ChildBitstreamID: &childbsid4,
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					/*
						Ptu: &examplecomv1.ChildBsPtu{
							Cids: &cids4,
							ID:   &id4,
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
					*/
					LLDMA: &examplecomv1.ChildBsLLDMA{
						Cids: &cids4,
						ID:   &id4,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id4,
						Identifier: &identifier4,
						Type:       &typ4,
						Version:    &varsion4,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id4,
						Identifier: &identifier4,
						Type:       &typ4,
						Version:    &varsion4,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id4,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier4,
							Type:       &typ4,
							Version:    &varsion4,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id4,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier4,
								Type:       &typ4,
								Version:    &varsion4,
							}},
							Parameters: &map[string]intstr.IntOrString{
								"5": {
									StrVal: "param01",
									IntVal: 12345,
									Type:   1,
								},
							},
							IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"0": {
									Available:      &available_false,
									FunctionCRName: &funcCRName4,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5_1,
												DMAChannelID:     &dmaChannel5_1,
												LLDMAConnectorID: &lldmaConnector5_1,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port5_2,
												DMAChannelID:     &dmaChannel5_2,
												LLDMAConnectorID: &lldmaConnector5_2,
											},
										},
									},
								},
								"1": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
									Rx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"RTP": {
												Port:             &port5_1,
												DMAChannelID:     &dmaChannel5_1,
												LLDMAConnectorID: &lldmaConnector5_1,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.ChildBsDetails{
											"DMA": {
												Port:             &port5_2,
												DMAChannelID:     &dmaChannel5_2,
												LLDMAConnectorID: &lldmaConnector5_2,
											},
										},
									},
								},
								"2": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"3": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"4": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"5": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"6": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
								"7": {
									Available:      &available4,
									FunctionCRName: &funcCRName4,
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity4,
								MaxDataFlows: &maxDataflows4,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions4,
				MaxCapacity:  &maxCapacity4,
				Name:         &name4,
			},
		},
		Status:           examplecomv1.ChildBsStatusReady,
		State:            examplecomv1.ChildBsReady,
		ChildBitstreamID: &childbsid4,
	},
}

var FPGA4 = []examplecomv1.FPGA{
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
			ChildBitstreamID:  &childbsid4,
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
		Status: examplecomv1.FPGAStatus{
			ChildBitstreamID:     &childbsid4,
			ChildBitstreamCRName: &name4,
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
			ChildBitstreamID:  &childbsid4,
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
			ChildBitstreamID:     &childbsid4,
			ChildBitstreamCRName: &name4,
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

var fpgafuncconfig_fr_high_infer_4 = corev1.ConfigMap{
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

var PCIeConnection5 = controllertestpcie.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night04-wbconnection-filter-resize-high-infer-main-high-infer-main",
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

var PCIeConnection6 = controllertestpcie.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night04-wbconnection-decode-main-filter-resize-high-infer-main",
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
