/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controllers_test

import (
	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	corev1 "k8s.io/api/core/v1"                                 //nolint:stylecheck // ST1019: intentional import as another name
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"                   //nolint:stylecheck // ST1019: intentional import as another name
)

var functionChain1 = ntthpcv1.FunctionChain{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functionchaintest1",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionChainSpec{
		FunctionTypeNamespace: "default",
		Functions: map[string]ntthpcv1.FunctionStruct{
			"decode-main": {
				FunctionName: "cpu-decode",
				Version:      "1.0.0",
			},
			"filter-resize-main": {
				FunctionName: "cpu-filter-resize",
				Version:      "1.0.0",
			},
			"copy-branch-main": {
				FunctionName: "copy-branch",
				Version:      "1.0.0",
			},
			"infer-1": {
				FunctionName: "person-infer",
				Version:      "1.0.0",
			},
			"infer-2": {
				FunctionName: "vehicle-infer",
				Version:      "1.0.0",
			},
		},
		Connections: []ntthpcv1.ConnectionStruct{
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "wb-start-of-chain",
					Port:        0,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "decode-main",
					Port:        0,
				},
			},
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "decode-main",
					Port:        0,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "filter-resize-main",
					Port:        0,
				},
			},
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "filter-resize-main",
					Port:        0,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "copy-branch-main",
					Port:        0,
				},
			},
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "copy-branch-main",
					Port:        0,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "infer-1",
					Port:        0,
				},
			},
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "copy-branch-main",
					Port:        1,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "infer-2",
					Port:        0,
				},
			},
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "infer-1",
					Port:        0,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "wb-end-of-chain-1",
					Port:        0,
				},
			},
			{
				ConnectionTypeName: "auto",
				From: ntthpcv1.FromToFunction{
					FunctionKey: "infer-2",
					Port:        0,
				},
				To: ntthpcv1.FromToFunction{
					FunctionKey: "wb-end-of-chain-2",
					Port:        0,
				},
			},
		},
	},
	Status: ntthpcv1.FunctionChainStatus{
		Status: "Not Ready",
	},
}

var functionType1 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype1",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "wb-start-of-chain",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo1",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo1 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo1",
		Namespace: "default",
	},
	Data: map[string]string{
		"spec": `[
			{
				"name": "spec1",
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType2 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype2",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "cpu-decode",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo2",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo2 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo2",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
        		"regionType": "cpu",
        		"inputInterfaceType": "host100gether",
        		"outputInterfaceType": "host100gether",
        		"configName": "cpufunc-config-decode",
        		"specName": "spec1"
			},
			{
				"name": "item2",
        		"regionType": "cpu",
        		"inputInterfaceType": "host100gether",
        		"outputInterfaceType": "host100gether",
        		"configName": "cpufunc-config-decode2",
        		"specName": "spec1"
			}
		]`,
		"spec": `[
			{
				"name": "spec1",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 20,
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType3 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype3",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "cpu-filter-resize",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo3",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo3 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo3",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-filter-resize",
				"specName": "spec1"
			}
		]`,
		"spec": `[
			{
				"name": "spec1",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 20,
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType4 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype4",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "copy-branch",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo4",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo4 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo4",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-copy-branch",
				"specName": "spec1"
			}
		]`,
		"spec": `[
			{
				"name": "spec1",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 15,
				"maxInputNum": 1,
				"maxOutputNum":1
			},
			{
				"name": "spec2",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 15,
				"maxInputNum": 1,
				"maxOutputNum":10
			},
			{
				"name": "spec3",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 15,
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType5 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype5",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "person-infer",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo5",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo5 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo5",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
        		"regionType": "cpu",
        		"inputInterfaceType": "host100gether",
        		"outputInterfaceType": "host100gether",
        		"configName": "gpufunc-config-person-infer",
        		"specName": "spec1"
			},
			{
				"name": "item2",
        		"regionType": "cpu",
        		"inputInterfaceType": "host100gether",
        		"outputInterfaceType": "host100gether",
        		"configName": "gpufunc-config-person-infer2",
        		"specName": "spec1"
			}
		]`,
		"spec": `[
			{
				"name": "spec1",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 15,
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType6 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype6",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "vehicle-infer",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo6",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo6 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo6",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
        		"regionType": "cpu",
        		"inputInterfaceType": "host100gether",
        		"outputInterfaceType": "host100gether",
        		"configName": "gpufunc-config-vehicle-infer",
        		"specName": "spec1"
			},
			{
				"name": "item2",
        		"regionType": "cpu",
        		"inputInterfaceType": "host100gether",
        		"outputInterfaceType": "host100gether",
        		"configName": "gpufunc-config-vehicle-infer2",
        		"specName": "spec1"
			}
		]`,
		"spec": `[
			{
				"name": "spec1",
				"minCore": 1,
				"maxCore": 1,
				"maxDataFlowsBase": 1,
				"maxCapacityBase": 15,
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType7 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype7",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "wb-end-of-chain-1",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo1",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo7 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo7",
		Namespace: "default",
	},
	Data: map[string]string{
		"spec": `[
			{
				"name": "spec1",
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}

var functionType8 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontype8",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "wb-end-of-chain-2",
		Version:      "1.0.0",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functioninfo8",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTypeStatus{
		Status: "Ready",
	},
}

var functionInfo8 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functioninfo8",
		Namespace: "default",
	},
	Data: map[string]string{
		"spec": `[
			{
				"name": "spec1",
				"maxInputNum": 1,
				"maxOutputNum":1
			}
		]`,
	},
}
