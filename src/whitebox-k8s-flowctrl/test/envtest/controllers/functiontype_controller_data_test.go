/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controllers_test

import (
	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	corev1 "k8s.io/api/core/v1"                                 //nolint:stylecheck // ST1019: intentional import as another name
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"                   //nolint:stylecheck // ST1019: intentional import as another name
)

var functionType9 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest9",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest9",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}

var functionType10 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest10",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest10",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}

var functionType11 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest11",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest11",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}

var functionType12 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest12",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest12",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}

var functionType13 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest13",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest13",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}

var functionType14 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest14",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest14",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}

var functionType15 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest15",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest15",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}
var functionType16 = ntthpcv1.FunctionType{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest16",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTypeSpec{
		FunctionName: "test",
		FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontest16",
			Namespace: "default",
		},
		Version: "1.0.0",
	},
}
var functionInfo9 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest9",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-test1",
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

var functionInfo10 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest10",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-test2_1",
				"specName": "spec1"
			},
			{
				"name": "item2",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "mem",
				"configName": "cpufunc-config-test2_2",
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

var functionInfo11 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest11",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-test3_1",
				"specName": "spec1"
			},
			{
				"name": "item2",
				"regionType": "alveo",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "mem",
				"configName": "cpufunc-config-test3_2",
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

var functionInfo12 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest12",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-test4",
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

var functionInfo13 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest13",
		Namespace: "default",
	},
	Data: map[string]string{
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

var functionInfo14 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest14",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-test6",
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
		"recommend": `[
			{
			  "deployableItemName": "item1"
			}
		]`,
	},
}

var functionInfo15 = corev1.ConfigMap{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontest15",
		Namespace: "default",
	},
	Data: map[string]string{
		"deployableItems": `[
			{
				"name": "item1",
				"regionType": "cpu",
				"inputInterfaceType": "host100gether",
				"outputInterfaceType": "host100gether",
				"configName": "cpufunc-config-test7",
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
		"other": `[
			{
			  "otherName": "otherItem1"
			}
		]`,
	},
}
