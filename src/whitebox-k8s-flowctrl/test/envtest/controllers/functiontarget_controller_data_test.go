/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controllers_test

import (
	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"                   //nolint:stylecheck // ST1019: intentional import as another name
)

var functionTarget1 = ntthpcv1.FunctionTarget{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontargettest1.alveou250-3.lane1",
		Namespace: "default",
	},
	Spec: ntthpcv1.FunctionTargetSpec{
		ComputeResourceRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontargettest1",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTargetStatus{
		RegionName:  "lane0",
		RegionType:  "alveo",
		NodeName:    "node1",
		DeviceType:  "alveou250",
		DeviceIndex: 0,
		Available:   true,
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
		},
	},
}

var functionTarget2 = ntthpcv1.FunctionTarget{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontargettest2.alveou250-3.lane1",
		Namespace: "default",
		OwnerReferences: []v1.OwnerReference{
			{
				APIVersion:         "v1",
				Kind:               "fuctiontarget",
				Name:               "anotherOwner",
				UID:                "anotherUID",
				Controller:         func(b bool) *bool { return &b }(true),
				BlockOwnerDeletion: func(b bool) *bool { return &b }(false),
			},
		},
	},
	Spec: ntthpcv1.FunctionTargetSpec{
		ComputeResourceRef: ntthpcv1.WBNamespacedName{
			Name:      "functiontargettest2",
			Namespace: "default",
		},
	},
	Status: ntthpcv1.FunctionTargetStatus{
		RegionName:  "lane0",
		RegionType:  "alveo",
		NodeName:    "node1",
		DeviceType:  "alveou250",
		DeviceIndex: 0,
		Available:   true,
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
		},
	},
}

var computeResource1 = ntthpcv1.ComputeResource{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontargettest1",
		Namespace: "default",
	},
	Spec: ntthpcv1.ComputeResourceSpec{
		NodeName: "functiontargettest1",
		Regions: []ntthpcv1.RegionInfo{
			{
				Name:           "lane1",
				Type:           "alveou250-0100001c-2lanes-0nics",
				DeviceFilePath: "/dev/xpcie_21330621T00D",
				DeviceUUID:     func(s string) *string { return &s }("21330621t00d"),
				DeviceType:     "alveou250",
				DeviceIndex:    3,
				Available:      true,
				Status:         "NotReady",
				Functions: []ntthpcv1.FunctionInfrastruct{
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
	},
	Status: ntthpcv1.ComputeResourceStatus{
		NodeName: "functiontargettest1",
		Regions: []ntthpcv1.RegionInfo{
			{
				Name:           "lane1",
				Type:           "alveou250-0100001c-2lanes-0nics",
				DeviceFilePath: "/dev/xpcie_21330621T00D",
				DeviceUUID:     func(s string) *string { return &s }("21330621t00d"),
				DeviceType:     "alveou250",
				DeviceIndex:    3,
				Available:      true,
				Status:         "NotReady",
				Functions: []ntthpcv1.FunctionInfrastruct{
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
	},
}

var computeResource2 = ntthpcv1.ComputeResource{
	ObjectMeta: v1.ObjectMeta{
		Name:      "functiontargettest2",
		Namespace: "default",
	},
	Spec: ntthpcv1.ComputeResourceSpec{
		NodeName: "functiontargettest2",
		Regions: []ntthpcv1.RegionInfo{
			{
				Name:           "lane1",
				Type:           "alveou250-0100001c-2lanes-0nics",
				DeviceFilePath: "/dev/xpcie_21330621T00D",
				DeviceUUID:     func(s string) *string { return &s }("21330621t00d"),
				DeviceType:     "alveou250",
				DeviceIndex:    3,
				Available:      true,
				Status:         "NotReady",
				Functions: []ntthpcv1.FunctionInfrastruct{
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
	},
	Status: ntthpcv1.ComputeResourceStatus{
		NodeName: "functiontargettest2",
		Regions: []ntthpcv1.RegionInfo{
			{
				Name:           "lane1",
				Type:           "alveou250-0100001c-2lanes-0nics",
				DeviceFilePath: "/dev/xpcie_21330621T00D",
				DeviceUUID:     func(s string) *string { return &s }("21330621t00d"),
				DeviceType:     "alveou250",
				DeviceIndex:    3,
				Available:      true,
				Status:         "NotReady",
				Functions: []ntthpcv1.FunctionInfrastruct{
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
	},
}
