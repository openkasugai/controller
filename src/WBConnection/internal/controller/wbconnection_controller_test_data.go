/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "WBConnection/api/v1"
	controllertestethernet "WBConnection/internal/controller/test/type/ethernet"
	controllertestpcie "WBConnection/internal/controller/test/type/pcie"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// fpgalist-ph3 config
var connectionkindmap = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "connectionkindmap",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"connectionkindmap.json": `
		 [{
        "connectionMethod": "host-mem",
        "connectionCRKind": "PCIeConnection"
      },{
        "connectionMethod": "host-100gether",
        "connectionCRKind": "EthernetConnection"
      }]
		`,
	},
}
var wbconnection1Start = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest1-wbconnection-wb-start-of-chain-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-100gether",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest1",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wb-start-of-chain",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest1-wbfunction-decode-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection2Ether = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-100gether",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest2",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest2-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest2-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection3PCIe = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-mem",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest3",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest3-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest3-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection4End = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest4-wbconnection-high-infer-main-wb-end-of-chain",
		Namespace: "default",
	},
	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-100gether",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest4",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest4-wbfuncction-high-infer-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wb-end-of-chain",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection5PCIe = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-mem",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest5",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest5-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest5-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection6PCIe = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-mem",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest6",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest6-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest6-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection7PCIe = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest7-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},

	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-mem",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest7",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest7-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest7-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		Status: "",
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var wbconnection8PCIe = examplecomv1.WBConnection{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest8-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},

	Spec: examplecomv1.WBConnectionSpec{
		ConnectionMethod: "host-mem",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "wbconntest8",
			Namespace: "default",
		},
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest8-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "wbconntest8-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: examplecomv1.WBConnectionStatus{
		ConnectionMethod: "",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "",
			Namespace: "",
		},
		Status: "Deployed",
		From: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
		To: examplecomv1.FromToWBFunction{
			Port: 0,
			WBFunctionRef: examplecomv1.WBNamespacedName{
				Name:      "",
				Namespace: "",
			},
		},
	},
}

var pcieconnection = controllertestpcie.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest7-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: controllertestpcie.PCIeConnectionSpec{
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "wbconntest7",
			Namespace: "default",
		},

		From: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "wbconntest7-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "wbconntest7-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestpcie.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "wbconntest7",
			Namespace: "default",
		},
		Status: "Deployed",
		From: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "pcieconnectiontest1-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "Deployed",
		},
		To: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "pcieconnectiontest1-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
			Status: "Deployed",
		},
	},
}
var t = metav1.Time{
	Time: time.Now(),
}
var testTime = metav1.Time{
	Time: t.Time.AddDate(0, 0, -1),
}

var EthernetConnectionUpdate = controllertestethernet.EthernetConnection{
	TypeMeta: metav1.TypeMeta{
		Kind:       "EthernetConnection",
		APIVersion: "example.com/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest2upd-wbconnection-decode-main-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: controllertestethernet.EthernetConnectionSpec{
		DataFlowRef: controllertestethernet.WBNamespacedName{
			Name:      "wbconntest2upd",
			Namespace: "default",
		},
		From: controllertestethernet.EthernetFunctionSpec{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "wbconntest2upd-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: controllertestethernet.EthernetFunctionSpec{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "wbconntest2upd-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestethernet.EthernetConnectionStatus{
		DataFlowRef: controllertestethernet.WBNamespacedName{
			Name:      "wbconntest2upd",
			Namespace: "default",
		},
		From: controllertestethernet.EthernetFunctionStatus{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "wbconntest2upd-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		To: controllertestethernet.EthernetFunctionStatus{
			WBFunctionRef: controllertestethernet.WBNamespacedName{
				Name:      "wbconntest2upd-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
			Status: "OK",
		},
		Status:    "Running",
		StartTime: testTime,
	},
}

var PCIeConnectionUpdate = controllertestpcie.PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		Kind:       "",
		APIVersion: "",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "wbconntest3upd-wbconnection-decode-main-filter-resize-low-infer-main",
		Namespace: "default",
	},
	Spec: controllertestpcie.PCIeConnectionSpec{
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "wbconntest3upd",
			Namespace: "default",
		},
		From: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "wbconntest3upd-wbfunction-decode-main",
				Namespace: "default",
			},
		},
		To: controllertestpcie.PCIeFunctionSpec{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "wbconntest3upd-wbfunction-filter-resize-low-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: controllertestpcie.PCIeConnectionStatus{
		StartTime: testTime,
		DataFlowRef: controllertestpcie.WBNamespacedName{
			Name:      "wbconntest3upd",
			Namespace: "default",
		},
		Status: "Running",
		From: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "wbconntest3upd-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},

		To: controllertestpcie.PCIeFunctionStatus{
			WBFunctionRef: controllertestpcie.WBNamespacedName{
				Name:      "wbconntest3upd-wbfunction-decode-main",
				Namespace: "default",
			},
			Status: "OK",
		},
	},
}
