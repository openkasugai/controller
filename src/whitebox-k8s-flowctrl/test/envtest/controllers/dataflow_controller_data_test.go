/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package controllers_test

import (
	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getIntPtr(x int32) *int32 {
	return &x
}

var df5 = ntthpcv1.DataFlow{
	ObjectMeta: v1.ObjectMeta{
		Name:      "test",
		Namespace: "default",
	},
	Spec: ntthpcv1.DataFlowSpec{
		Requirements: &ntthpcv1.DataFlowRequirementsStruct{
			All: &ntthpcv1.AllRequirementsInfo{
				Capacity: 15,
			},
		},
	},
	Status: ntthpcv1.DataFlowStatus{
		Status: "WBFunction/WBConnection creation in progress",
		FunctionChain: &ntthpcv1.FunctionChain{
			Spec: ntthpcv1.FunctionChainSpec{
				Functions: map[string]ntthpcv1.FunctionStruct{
					"decode-main": {
						FunctionName: "decode",
						Version:      "1.0.0",
					},
					"filter-resize-high-infer-main": {
						FunctionName: "filter-resize-high-infer",
						Version:      "1.0.0",
					},
					"high-infer-main": {
						FunctionName: "high-infer",
						Version:      "1.0.0",
					},
				},
				Connections: []ntthpcv1.ConnectionStruct{
					{
						From: ntthpcv1.FromToFunction{
							FunctionKey: "wb-start-of-chain",
						},
						To: ntthpcv1.FromToFunction{
							FunctionKey: "decode-main",
						},
						ConnectionTypeName: "auto",
					},
					{
						From: ntthpcv1.FromToFunction{
							FunctionKey: "decode-main",
						},
						To: ntthpcv1.FromToFunction{
							FunctionKey: "filter-resize-high-infer-main",
						},
						ConnectionTypeName: "auto",
					},
					{
						From: ntthpcv1.FromToFunction{
							FunctionKey: "filter-resize-high-infer-main",
						},
						To: ntthpcv1.FromToFunction{
							FunctionKey: "high-infer-main",
						},
						ConnectionTypeName: "auto",
					},
					{
						From: ntthpcv1.FromToFunction{
							FunctionKey: "high-infer-main",
						},
						To: ntthpcv1.FromToFunction{
							FunctionKey: "wb-end-of-chain",
						},
						ConnectionTypeName: "auto",
					},
				},
			},
		},
		FunctionType: []*ntthpcv1.FunctionType{
			{
				Spec: ntthpcv1.FunctionTypeSpec{
					FunctionName: "decode",
					FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
						Name:      "funcinfo-decode",
						Namespace: "default",
					},
					Version: "1.0.0",
				},
				Status: ntthpcv1.FunctionTypeStatus{
					Status:               "Ready",
					RegionTypeCandidates: []string{"alveo"},
				},
			},
			{
				Spec: ntthpcv1.FunctionTypeSpec{
					FunctionName: "filter-resize-high-infer",
					FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
						Name:      "funcinfo-filter-resize-high-infer",
						Namespace: "default",
					},
					Version: "1.0.0",
				},
				Status: ntthpcv1.FunctionTypeStatus{
					Status:               "Ready",
					RegionTypeCandidates: []string{"alveo"},
				},
			},
			{
				Spec: ntthpcv1.FunctionTypeSpec{
					FunctionName: "high-infer",
					FunctionInfoCMRef: ntthpcv1.WBNamespacedName{
						Name:      "funcinfo-high-infer",
						Namespace: "default",
					},
					Version: "1.0.0",
				},
				Status: ntthpcv1.FunctionTypeStatus{
					Status:               "Ready",
					RegionTypeCandidates: []string{"a100"},
				},
			},
		},
		ScheduledFunctions: map[string]ntthpcv1.FunctionScheduleInfo{
			"decode-main": {
				NodeName:      "node1",
				DeviceType:    "alveo",
				DeviceIndex:   1,
				RegionName:    "lane0",
				FunctionIndex: getIntPtr(1),
			},
			"filter-resize-high-infer-main": {
				NodeName:      "node1",
				DeviceType:    "alveo",
				DeviceIndex:   0,
				RegionName:    "lane1",
				FunctionIndex: getIntPtr(5),
			},
			"high-infer-main": {
				NodeName:      "node1",
				DeviceType:    "a100",
				DeviceIndex:   1,
				RegionName:    "gpu-5b771964-ab74-a674-15d7-8f0d2cee4ef8",
				FunctionIndex: nil,
			},
		},
		ScheduledConnections: []ntthpcv1.ConnectionScheduleInfo{
			{
				From: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "wb-start-of-chain",
				},
				To: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "decode-main",
				},
				ConnectionMethod: "host-100gether",
			},
			{
				From: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "decode-main",
				},
				To: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "filter-resize-high-infer-main",
				},
				ConnectionMethod: "host-100gether",
			},
			{
				From: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "filter-resize-high-infer-main",
				},
				To: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "high-infer-main",
				},
				ConnectionMethod: "host-mem",
			},
			{
				From: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "high-infer-main",
				},
				To: ntthpcv1.FromToFunctionScheduleInfo{
					FunctionKey: "wb-end-of-chain",
				},
				ConnectionMethod: "host-100gether",
			},
		},
	},
}

//nolint:unused // ftList is unused. FIXME: remove this variable
var ftList = []ntthpcv1.FunctionTarget{
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.alveou250-0.lane0",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "lane0",
			RegionType:       "alveo",
			NodeName:         "node1",
			DeviceType:       "alveou250",
			DeviceIndex:      0,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(2),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(16),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "filter-resize",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(8),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "filter-resize",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(8),
				},
			},
		},
	},
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.alveou250-0.lane1",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "lane1",
			RegionType:       "alveo",
			NodeName:         "node1",
			DeviceType:       "alveou250",
			DeviceIndex:      0,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(2),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(16),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "filter-resize",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "filter-resize",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(1),
				},
			},
		},
	},
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.alveou250-1.lane0",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "lane0",
			RegionType:       "alveo",
			NodeName:         "node1",
			DeviceType:       "alveou250",
			DeviceIndex:      1,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(2),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(30),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "decode",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "decode",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(60),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
			},
		},
	},
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.alveou250-1.lane1",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "lane1",
			RegionType:       "alveo",
			NodeName:         "node1",
			DeviceType:       "alveou250",
			DeviceIndex:      1,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(2),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(30),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "decode",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "decode",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(8),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
			},
		},
	},
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.a100-0.gpu",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "gpu",
			RegionType:       "a100",
			NodeName:         "node1",
			DeviceType:       "a100",
			DeviceIndex:      0,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(120),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(30),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "high-infer",
					Available:        false,
					MaxDataFlows:     func(i int32) *int32 { return &i }(1),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "high-infer",
					Available:        false,
					MaxDataFlows:     func(i int32) *int32 { return &i }(1),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
			},
		},
	},
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.a100-1.gpu",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "gpu",
			RegionType:       "a100",
			NodeName:         "node1",
			DeviceType:       "a100",
			DeviceIndex:      1,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(120),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(45),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "high-infer",
					Available:        false,
					MaxDataFlows:     func(i int32) *int32 { return &i }(1),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(30),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "high-infer",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(1),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
			},
		},
	},
	{
		ObjectMeta: v1.ObjectMeta{
			Name:      "node1.a100-2.gpu",
			Namespace: "default",
		},
		Status: ntthpcv1.FunctionTargetStatus{
			RegionName:       "gpu",
			RegionType:       "a100",
			NodeName:         "node1",
			DeviceType:       "a100",
			DeviceIndex:      2,
			Available:        true,
			MaxFunctions:     func(i int32) *int32 { return &i }(120),
			CurrentFunctions: func(i int32) *int32 { return &i }(2),
			MaxCapacity:      func(i int32) *int32 { return &i }(100),
			CurrentCapacity:  func(i int32) *int32 { return &i }(30),
			Functions: []ntthpcv1.FunctionCapStruct{
				{
					FunctionIndex:    1,
					FunctionName:     "high-infer",
					Available:        false,
					MaxDataFlows:     func(i int32) *int32 { return &i }(1),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
				{
					FunctionIndex:    2,
					FunctionName:     "high-infer",
					Available:        true,
					MaxDataFlows:     func(i int32) *int32 { return &i }(1),
					CurrentDataFlows: func(i int32) *int32 { return &i }(1),
					MaxCapacity:      func(i int32) *int32 { return &i }(30),
					CurrentCapacity:  func(i int32) *int32 { return &i }(15),
				},
			},
		},
	},
}

// + kubectl get functiontargets.example.com -A -o yaml
// apiVersion: v1
// items:
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.a100-4.gpu-5b771964-ab74-a674-15d7-8f0d2cee4ef8
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549741"
//     uid: c1a8bc48-e8c8-444d-90eb-4af58f96e074
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: false
//     capacityTotal: 120
//     capacityUsed: 120
//     deviceIndex: 4
//     deviceKind: a100
//     functions:
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 201
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 202
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 203
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 204
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 205
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 206
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 207
//       name: high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 15
//       capacityUsed: 15
//       functionIndex: 208
//       name: high-infer
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: gpu-5b771964-ab74-a674-15d7-8f0d2cee4ef8
//     regionType: a100
//     total: 110
//     used: 8
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-0.lane0
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2548723"
//     uid: f174c31d-4815-49fb-a031-c31ddebe9f88
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 42
//     capacityUsed: 30
//     deviceIndex: 0
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 5
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 6
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane0
//     regionType: alveo
//     total: 8
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-0.lane1
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549122"
//     uid: 6c7d37e7-3d63-4a53-8161-397c1c94bf56
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 42
//     capacityUsed: 30
//     deviceIndex: 0
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 7
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 8
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane1
//     regionType: alveo
//     total: 8
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-1.lane0
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549104"
//     uid: 780c9ef6-8e40-4b7b-b4dd-064f255e84aa
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 42
//     capacityUsed: 30
//     deviceIndex: 1
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 1
//       name: decode
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 2
//       name: decode
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane0
//     regionType: alveo
//     total: 12
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-1.lane1
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2548788"
//     uid: 57325164-473e-4ec1-b21f-4ce18eb9c586
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 42
//     capacityUsed: 30
//     deviceIndex: 1
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 3
//       name: decode
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 21
//       capacityUsed: 15
//       functionIndex: 4
//       name: decode
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane1
//     regionType: alveo
//     total: 12
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-2.lane0
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549413"
//     uid: d2211467-d88e-47c7-b06b-e3b5e7047123
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 40
//     capacityUsed: 30
//     deviceIndex: 2
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 9
//       name: decode
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 10
//       name: decode
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane0
//     regionType: alveo
//     total: 12
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-2.lane1
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549695"
//     uid: 45c04a1c-01b7-423a-a864-0e11e01da0a0
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 40
//     capacityUsed: 30
//     deviceIndex: 2
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 11
//       name: decode
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 12
//       name: decode
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane1
//     regionType: alveo
//     total: 12
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-3.lane0
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549734"
//     uid: b08b9523-68c7-4778-aead-e22f66141230
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 40
//     capacityUsed: 30
//     deviceIndex: 3
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 13
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 14
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane0
//     regionType: alveo
//     total: 8
//     used: 2
// - apiVersion: example.com/v1
//   kind: FunctionTarget
//   metadata:
//     creationTimestamp: "2023-03-22T06:20:24Z"
//     generation: 1
//     name: swb-sm7.alveo-3.lane1
//     namespace: default
//     ownerReferences:
//     - apiVersion: example.com/v1
//       blockOwnerDeletion: true
//       controller: true
//       kind: ComputeResource
//       name: compute-swb-sm7
//       uid: 4784b334-45fb-491f-ba03-b18db6a4d973
//     resourceVersion: "2549423"
//     uid: 865e461a-d0c8-42b6-868e-640348a40c41
//   spec:
//     computeResourceName: compute-swb-sm7
//     namespace: default
//   status:
//     available: true
//     capacityTotal: 40
//     capacityUsed: 30
//     deviceIndex: 3
//     deviceKind: alveo
//     functions:
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 15
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     - available: false
//       capacityTotal: 20
//       capacityUsed: 15
//       functionIndex: 16
//       name: filter-resize-high-infer
//       total: 1
//       used: 1
//     node: swb-sm7
//     region: lane1
//     regionType: alveo
//     total: 8
//     used: 2
// kind: List
// metadata:
//   resourceVersion: ""
//   selfLink: ""

// FunctionInfo

// var cmList1 = []v1.ConfigMap{
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.alveou250-0.lane0",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "lane0",
// 			RegionType:    "alveo",
// 			Node:          "node1",
// 			DeviceKind:    "alveou250",
// 			DeviceIndex:   0,
// 			Available:     true,
// 			Total:         2,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  16,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "filter-resize",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  8,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "filter-resize",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  8,
// 				},
// 			},
// 		},
// 	},
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.alveou250-0.lane1",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "lane1",
// 			RegionType:    "alveo",
// 			Node:          "node1",
// 			DeviceKind:    "alveou250",
// 			DeviceIndex:   0,
// 			Available:     true,
// 			Total:         2,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  16,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "filter-resize",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "filter-resize",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  1,
// 				},
// 			},
// 		},
// 	},
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.alveou250-1.lane0",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "lane0",
// 			RegionType:    "alveo",
// 			Node:          "node1",
// 			DeviceKind:    "alveou250",
// 			DeviceIndex:   1,
// 			Available:     true,
// 			Total:         2,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  30,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "decode",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "decode",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 60,
// 					CapacityUsed:  15,
// 				},
// 			},
// 		},
// 	},
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.alveou250-1.lane1",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "lane1",
// 			RegionType:    "alveo",
// 			Node:          "node1",
// 			DeviceKind:    "alveou250",
// 			DeviceIndex:   1,
// 			Available:     true,
// 			Total:         2,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  30,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "decode",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "decode",
// 					Available:     true,
// 					Total:         8,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 			},
// 		},
// 	},
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.a100-0.gpu",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "gpu",
// 			RegionType:    "a100",
// 			Node:          "node1",
// 			DeviceKind:    "a100",
// 			DeviceIndex:   0,
// 			Available:     true,
// 			Total:         120,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  30,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "high-infer",
// 					Available:     false,
// 					Total:         1,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "high-infer",
// 					Available:     false,
// 					Total:         1,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 			},
// 		},
// 	},
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.a100-1.gpu",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "gpu",
// 			RegionType:    "a100",
// 			Node:          "node1",
// 			DeviceKind:    "a100",
// 			DeviceIndex:   1,
// 			Available:     true,
// 			Total:         120,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  45,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "high-infer",
// 					Available:     false,
// 					Total:         1,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  30,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "high-infer",
// 					Available:     true,
// 					Total:         1,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 			},
// 		},
// 	},
// 	{
// 		ObjectMeta: v1.ObjectMeta{
// 			Name:      "node1.a100-2.gpu",
// 			Namespace: "default",
// 		},
// 		Status: ntthpcv1.FunctionTargetStatus{
// 			Region:        "gpu",
// 			RegionType:    "a100",
// 			Node:          "node1",
// 			DeviceKind:    "a100",
// 			DeviceIndex:   2,
// 			Available:     true,
// 			Total:         120,
// 			Used:          2,
// 			CapacityTotal: 100,
// 			CapacityUsed:  30,
// 			Functions: []ntthpcv1.FunctionCapStruct{
// 				{
// 					FunctionIndex: 1,
// 					Name:          "high-infer",
// 					Available:     false,
// 					Total:         1,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 				{
// 					FunctionIndex: 2,
// 					Name:          "high-infer",
// 					Available:     true,
// 					Total:         1,
// 					Used:          1,
// 					CapacityTotal: 30,
// 					CapacityUsed:  15,
// 				},
// 			},
// 		},
// 		},
// }

// ubuntu@wb-m6:~/work_alpha/itac/yaml$ cat functioninfo.yaml
// apiVersion: v1
// items:
// - apiVersion: v1
//   kind: ConfigMap
//   metadata:
//     name: funcinfo-decode
//     namespace: wbfunc-imgproc
//   data:
//     alveo: '{
//       "items": {
//         "host-100gether": {
//           "configName": "fpgafunc-config-decode",
//           "minCore": 1,
//           "maxCore": 1,
//           "maxDataFlowsBase": 6,
//           "maxCapacityBase": 20
//         }
//       }
//     }'
// - apiVersion: v1
//   kind: ConfigMap
//   metadata:
//     name: funcinfo-filter-resize-high-infer
//     namespace: wbfunc-imgproc
//   data:
//     alveo: '{
//       "items": {
//         "host-100gether": {
//           "configName": "fpgafunc-config-filter-resize-high-infer",
//           "minCore": 1,
//           "maxCore": 1,
//           "maxDataFlowsBase": 8,
//           "maxCapacityBase": 40
//         },
//         "host-mem": {
//           "configName": "fpgafunc-config-filter-resize-high-infer",
//           "minCore": 1,
//           "maxCore": 1,
//           "maxDataFlowsBase": 8,
//           "maxCapacityBase": 40
//         }
//       }
//     }'
// - apiVersion: v1
//   kind: ConfigMap
//   metadata:
//     name: funcinfo-high-infer
//     namespace: wbfunc-imgproc
//   data:
//     a100: '{
//       "items": {
//         "host-mem": {
//           "configName": "gpufunc-config-high-infer",
//           "minCore": 1,
//           "maxCore": 1,
//           "maxDataFlowsBase": 110,
//           "maxCapacityBase": 120
//         },
//         "host-100gether": {
//           "configName": "gpufunc-config-high-infer",
//           "minCore": 1,
//           "maxCore": 1,
//           "maxDataFlowsBase": 110,
//           "maxCapacityBase": 120
//         }
//       }
//     }'
// kind: List
