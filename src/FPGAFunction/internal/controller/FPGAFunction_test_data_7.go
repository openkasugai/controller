/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "FPGAFunction/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var fpgafunctionUPDATE = examplecomv1.FPGAFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGAFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fpgafunctestupdate-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
		Finalizers: []string{
			"fpgafunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.FPGAFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName5,
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "fpgafunctestupdate",
			Namespace: "default",
		},
		DeviceType: "alveo",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "fpgafunctestupdate-wbfunction-decode-main",
					Namespace: "default",
				},
			},
		},
		NextFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "fpgafunctestupdate-wbfunction-high-infer-main",
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
			Name:      "fpgafunctestupdate",
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
