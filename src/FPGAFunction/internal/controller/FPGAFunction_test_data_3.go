package controller

import (
	examplecomv1 "FPGAFunction/api/v1"
	controllertestcpu "FPGAFunction/internal/controller/test/type/CPU"
	controllertestethernet "FPGAFunction/internal/controller/test/type/Ethernet"
	controllertestgpu "FPGAFunction/internal/controller/test/type/GPU"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var frameworkKernelID2 int32 = 1
var functionChannelID2 int32 = 1
var functionIndex2 int32 = 1
var functionKernelID2 int32 = 1
var ptuKernelID2 int32 = 1
var partitionName2 string = "0"

var FPGAFunction2 = examplecomv1.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.FPGAFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName2,
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night02",
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
		FrameworkKernelID: &frameworkKernelID2,
		FunctionChannelID: &functionChannelID2,
		FunctionIndex:     &functionIndex2,
		FunctionKernelID:  &functionKernelID2,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "test01",
		PtuKernelID:       &ptuKernelID2,
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

var CPUFunction2 = controllertestcpu.CPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "CPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cpufunctiontest2-wbfunction-decode-main",
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
			Name:      "df-night02",
			Namespace: "default",
		},
		DeviceType:   "cpu",
		FunctionName: "cpu-decode",
		NextFunctions: map[string]controllertestcpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestcpu.WBNamespacedName{
					Name:      "cpufunctiontest2-wbfunction-filter-resize-low-infer-main",
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
			FilePrefix:      "test02-cpufunctiontest1-wbfunction-decode-main",
			CommandQueueID:  "test02-cpufunctiontest1-wbfunction-decode-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "cpu",
	},
	Status: controllertestcpu.CPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: controllertestcpu.WBNamespacedName{
			Name:      "df-night02",
			Namespace: "default",
		},
		FunctionName: "cpu-decode",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var functionIndexg2 int32 = 0

var GPUFunction2 = controllertestgpu.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestgpu.GPUFunctionSpec{
		AcceleratorIDs: []controllertestgpu.AccIDInfo{
			{
				PartitionName: "df-night02-wbfunction-high-infer-main",
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: controllertestgpu.WBNamespacedName{
			Name:      "df-night02",
			Namespace: "default",
		},
		DeviceType:    "a100",
		FunctionIndex: &functionIndexg2,
		FunctionName:  "high-infer",
		PreviousFunctions: map[string]controllertestgpu.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: controllertestgpu.WBNamespacedName{
					Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
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
			FilePrefix:      "df-night02-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night02-wbfunction-filter-resize-high-infer-main",
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

var childbsid2 string = "2222222222"
var maxFunctions2 int32 = 1
var maxCapacity2 int32 = 2
var name2 string = "child1"
var cids2 string = "111"
var id2 int32 = 3
var identifier2 string = "child1_identifier"
var typ2 string = "childbs_chaintype"
var varsion2 string = "childbs_varsion1.1.3"
var maxDataflows2 int32 = 4
var available2 bool = true
var funcCRName2 string = "funcCRName"
var port3_1 int32 = 5
var dmaChannel3_1 int32 = 6
var lldmaConnector3_1 int32 = 7
var port3_2 int32 = 8
var dmaChannel3_2 int32 = 9
var lldmaConnector3_2 int32 = 10
var uid2 types.UID = "aaaaaaa"

var ChildBitstream1 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "Childbs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fpga-21320621v00dtest012222222222",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga3",
				UID:        uid2,
			},
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids2,
						ID:   &id2,
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
						Cids: &cids2,
						ID:   &id2,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id2,
						Identifier: &identifier2,
						Type:       &typ2,
						Version:    &varsion2,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id2,
						Identifier: &identifier2,
						Type:       &typ2,
						Version:    &varsion2,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id2,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier2,
							Type:       &typ2,
							Version:    &varsion2,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id2,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier2,
								Type:       &typ2,
								Version:    &varsion2,
							}},
							/*Parameters: &map[string]intstr.IntOrString{
							   	"0": {
							     	StrVal: "param01",
							         IntVal: 12345,
							         Type:   1,
							     },
							  },
							*/IntraResourceMgmtMap: &map[string]examplecomv1.FunctionsIntraResourceMgmtMap{
								"1": {
									Available:      &available2,
									FunctionCRName: &funcCRName2,
									Rx: &examplecomv1.RxTxSpec{

										Protocol: &map[string]examplecomv1.Details{
											"RTP": {
												Port:             &port3_1,
												DMAChannelID:     &dmaChannel3_1,
												LLDMAConnectorID: &lldmaConnector3_1,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.Details{
											"DMA": {
												Port:             &port3_2,
												DMAChannelID:     &dmaChannel3_2,
												LLDMAConnectorID: &lldmaConnector3_2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity2,
								MaxDataFlows: &maxDataflows2,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions2,
				MaxCapacity:  &maxCapacity2,
				Name:         &name2,
			},
		},
		ChildBitstreamID: &childbsid2,
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Modules: &examplecomv1.ChildBsModule{
					Ptu: &examplecomv1.ChildBsPtu{
						Cids: &cids2,
						ID:   &id2,
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
						Cids: &cids2,
						ID:   &id2,
					},
					Chain: &examplecomv1.ChildBsChain{
						ID:         &id2,
						Identifier: &identifier2,
						Type:       &typ2,
						Version:    &varsion2,
					},
					Directtrans: &examplecomv1.ChildBsDirecttrans{
						ID:         &id2,
						Identifier: &identifier2,
						Type:       &typ2,
						Version:    &varsion2,
					},
					Conversion: &examplecomv1.ChildBsConversion{
						ID: &id2,
						Module: &[]examplecomv1.ConversionModule{{
							Identifier: &identifier2,
							Type:       &typ2,
							Version:    &varsion2,
						}},
					},
					Functions: &[]examplecomv1.ChildBsFunctions{
						{
							ID: &id2,
							Module: &[]examplecomv1.FunctionsModule{{
								Identifier: &identifier2,
								Type:       &typ2,
								Version:    &varsion2,
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
									Available:      &available2,
									FunctionCRName: &funcCRName2,
									Rx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.Details{
											"RTP": {
												Port:             &port3_1,
												DMAChannelID:     &dmaChannel3_1,
												LLDMAConnectorID: &lldmaConnector3_1,
											},
										},
									},
									Tx: &examplecomv1.RxTxSpec{
										Protocol: &map[string]examplecomv1.Details{
											"DMA": {
												Port:             &port3_2,
												DMAChannelID:     &dmaChannel3_2,
												LLDMAConnectorID: &lldmaConnector3_2,
											},
										},
									},
								},
							},
							DeploySpec: examplecomv1.FunctionsDeploySpec{
								MaxCapacity:  &maxCapacity2,
								MaxDataFlows: &maxDataflows2,
							},
						},
					},
				},
				MaxFunctions: &maxFunctions2,
				MaxCapacity:  &maxCapacity2,
				Name:         &name2,
			},
		},
		Status:           examplecomv1.ChildBsStatusReady,
		State:            examplecomv1.ChildBsReady,
		ChildBitstreamID: &childbsid2,
	},
}

var FPGA2 = []examplecomv1.FPGA{
	{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "FPGA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fpga-21320621v00dtest01",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid2,
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
			ChildBitstreamID:     &childbsid2,
			ChildBitstreamCRName: &name2,
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
			Name:      "childbs2",
			Namespace: "default",
		},
		Spec: examplecomv1.FPGASpec{
			ChildBitstreamID:  &childbsid2,
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
			ChildBitstreamID:     &childbsid2,
			ChildBitstreamCRName: &name2,
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

var fpgafuncconfig_fr_high_infer_2 = corev1.ConfigMap{
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

var EthernetConnection3 = controllertestethernet.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbconnection-decode-main-filter-resize-high-infer-main",
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

var EthernetConnection4 = controllertestethernet.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "EthernetConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbconnection-filter-resize-high-infer-main-high-infer-main",
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
