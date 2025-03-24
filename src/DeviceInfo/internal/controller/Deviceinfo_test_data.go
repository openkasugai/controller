/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	examplecomv1 "DeviceInfo/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var MaxDataFlows int32 = 1
var MaxCapacity int32 = 20
var Capacity int32 = 15
var FunctionIndex int32 = 0
var FunctionIndex2 int32 = 1

var DeviceInfo1 = examplecomv1.DeviceInfo{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "DeviceInfo",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deviceinfo-df-night01-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.DeviceInfoSpec{
		Request: examplecomv1.WBFuncRequest{
			RequestType:   "Deploy",
			DeviceType:    "cpu",
			DeviceIndex:   0,
			RegionName:    "cpu",
			NodeName:      "test01",
			FunctionName:  "cpu-decode",
			FunctionIndex: &FunctionIndex,
			MaxDataFlows:  &MaxDataFlows,
			MaxCapacity:   &MaxCapacity,
			Capacity:      &Capacity,
		},
	},
}

var DeviceInfo2 = examplecomv1.DeviceInfo{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "DeviceInfo",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deviceinfo-df-night01-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.DeviceInfoSpec{
		Request: examplecomv1.WBFuncRequest{
			RequestType:   "Undeploy",
			DeviceType:    "cpu",
			DeviceIndex:   0,
			RegionName:    "cpu",
			NodeName:      "test01",
			FunctionName:  "cpu-decode",
			FunctionIndex: &FunctionIndex,
			MaxDataFlows:  &MaxDataFlows,
			MaxCapacity:   &MaxCapacity,
			Capacity:      &Capacity,
		},
	},
}

var DeviceInfo3 = examplecomv1.DeviceInfo{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "DeviceInfo",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deviceinfo-df-night02-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.DeviceInfoSpec{
		Request: examplecomv1.WBFuncRequest{
			RequestType:   "Deploy",
			DeviceType:    "cpu",
			DeviceIndex:   0,
			RegionName:    "cpu",
			NodeName:      "test01",
			FunctionName:  "cpu-decode",
			FunctionIndex: &FunctionIndex2,
			MaxDataFlows:  &MaxDataFlows,
			MaxCapacity:   &MaxCapacity,
			Capacity:      &Capacity,
		},
	},
}

var MaxDataFlowsFPGA int32 = 1
var MaxCapacityFPGA int32 = 20
var CapacityFPGA int32 = 15
var FunctionIndexFPGA int32 = 0

var DeviceInfo4 = examplecomv1.DeviceInfo{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "DeviceInfo",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deviceinfo-df-night03-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.DeviceInfoSpec{
		Request: examplecomv1.WBFuncRequest{
			RequestType:   "Deploy",
			DeviceType:    "alveo",
			DeviceIndex:   1,
			RegionName:    "lane0",
			NodeName:      "test01",
			FunctionName:  "decode",
			FunctionIndex: &FunctionIndexFPGA,
			MaxDataFlows:  &MaxDataFlowsFPGA,
			MaxCapacity:   &MaxCapacityFPGA,
			Capacity:      &CapacityFPGA,
		},
	},
}

var DeviceInfo5 = examplecomv1.DeviceInfo{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "DeviceInfo",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deviceinfo-df-night03-wbfunction-decode-main",
		Namespace: "default",
	},
	Spec: examplecomv1.DeviceInfoSpec{
		Request: examplecomv1.WBFuncRequest{
			RequestType:   "Undeploy",
			DeviceType:    "alveo",
			DeviceIndex:   1,
			RegionName:    "lane0",
			NodeName:      "test01",
			FunctionName:  "decode",
			FunctionIndex: &FunctionIndexFPGA,
			MaxDataFlows:  &MaxDataFlowsFPGA,
			MaxCapacity:   &MaxCapacityFPGA,
			Capacity:      &CapacityFPGA,
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
var CurrentFunctions4 int32 = 0
var CurrentDataFlows int32 = 0
var CurrentDataFlows2 int32 = 0
var MaxCapacity2 int32 = 40
var MaxCapacity3 int32 = 120
var MaxCapacity4 int32 = 240
var MaxDataFlows1 int32 = 6
var MaxDataFlows2 int32 = 8
var DeviceUUID string = "21330621T01J"
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

var chkCurrentCapacity3 int32 = 0
var chkCurrentCapacity99 int32 = 0
var chkCurrentFunctions3 int32 = 0
var chkCurrentDataFlows99 int32 = 0
var chkMaxCapacity4 int32 = 240
var chkMaxCapacity99 int32 = 20
var chkMaxDataFlows99 int32 = 1
var chkDeviceUUID5 string = "swb-sm7-cpu0"
var chkMaxFunctions4 int32 = 220

var chkComRes1 = examplecomv1.RegionInfo{
	Available:        true,
	CurrentCapacity:  &chkCurrentCapacity3,
	CurrentFunctions: &chkCurrentFunctions3,
	DeviceFilePath:   "",
	DeviceIndex:      0,
	DeviceType:       "cpu",
	DeviceUUID:       &chkDeviceUUID5,
	Functions: []examplecomv1.FunctionInfrastruct{{
		Available:        true,
		CurrentCapacity:  &chkCurrentCapacity99,
		CurrentDataFlows: &chkCurrentDataFlows99,
		FunctionIndex:    0,
		FunctionName:     "cpu-decode",
		MaxCapacity:      &chkMaxCapacity99,
		MaxDataFlows:     &chkMaxDataFlows99,
		PartitionName:    chkDeviceUUID5,
	}},
	MaxCapacity:  &chkMaxCapacity4,
	MaxFunctions: &chkMaxFunctions4,
	Name:         "cpu",
	Type:         "cpu",
}

var chkCurrentCapacity4 int32 = 15
var chkCurrentCapacity98 int32 = 15
var chkCurrentFunctions4 int32 = 1
var chkCurrentDataFlows98 int32 = 1
var chkMaxCapacity5 int32 = 240
var chkMaxCapacity98 int32 = 20
var chkMaxDataFlows98 int32 = 1
var chkDeviceUUID6 string = "swb-sm7-cpu0"
var chkMaxFunctions5 int32 = 220

var chkComRes2 = examplecomv1.RegionInfo{
	Available:        true,
	CurrentCapacity:  &chkCurrentCapacity4,
	CurrentFunctions: &chkCurrentFunctions4,
	DeviceFilePath:   "",
	DeviceIndex:      0,
	DeviceType:       "cpu",
	DeviceUUID:       &chkDeviceUUID6,
	Functions: []examplecomv1.FunctionInfrastruct{{
		Available:        false,
		CurrentCapacity:  &chkCurrentCapacity98,
		CurrentDataFlows: &chkCurrentDataFlows98,
		FunctionIndex:    1,
		FunctionName:     "cpu-decode",
		MaxCapacity:      &chkMaxCapacity98,
		MaxDataFlows:     &chkMaxDataFlows98,
		PartitionName:    chkDeviceUUID5,
	}},
	MaxCapacity:  &chkMaxCapacity5,
	MaxFunctions: &chkMaxFunctions5,
	Name:         "cpu",
	Type:         "cpu",
}

var chkCurrentCapacity5 int32 = 15
var chkCurrentCapacity97 int32 = 15
var chkCurrentFunctions5 int32 = 2
var chkCurrentDataFlows97 int32 = 1
var chkMaxCapacity6 int32 = 40
var chkMaxCapacity97 int32 = 20
var chkMaxDataFlows97 int32 = 6
var chkDeviceUUID7 string = "21330621T01J"
var chkMaxFunctions6 int32 = 40

var chkComRes3 = examplecomv1.RegionInfo{
	Available:        true,
	CurrentCapacity:  &chkCurrentCapacity5,
	CurrentFunctions: &chkCurrentFunctions5,
	DeviceFilePath:   "/dev/xpcie_21330621T01J",
	DeviceIndex:      1,
	DeviceType:       "alveo",
	DeviceUUID:       &chkDeviceUUID7,
	Functions: []examplecomv1.FunctionInfrastruct{{
		Available:        true,
		CurrentCapacity:  &chkCurrentCapacity97,
		CurrentDataFlows: &chkCurrentDataFlows97,
		FunctionIndex:    0,
		FunctionName:     "decode",
		MaxCapacity:      &chkMaxCapacity97,
		MaxDataFlows:     &chkMaxDataFlows97,
		PartitionName:    "0",
	}},
	MaxCapacity:  &chkMaxCapacity6,
	MaxFunctions: &chkMaxFunctions6,
	Name:         "lane0",
	Type:         "alveo",
}

var chkCurrentCapacity6 int32 = 0
var chkCurrentCapacity96 int32 = 0
var chkCurrentFunctions6 int32 = 2
var chkCurrentDataFlows96 int32 = 0
var chkMaxCapacity7 int32 = 40
var chkMaxCapacity96 int32 = 20
var chkMaxDataFlows96 int32 = 5
var chkMaxFunctions7 int32 = 40

var chkComRes4 = examplecomv1.RegionInfo{
	Available:        true,
	CurrentCapacity:  &chkCurrentCapacity6,
	CurrentFunctions: &chkCurrentFunctions6,
	DeviceFilePath:   "/dev/xpcie_21330621T01J",
	DeviceIndex:      1,
	DeviceType:       "alveo",
	DeviceUUID:       &chkDeviceUUID7,
	Functions: []examplecomv1.FunctionInfrastruct{{
		Available:        true,
		CurrentCapacity:  &chkCurrentCapacity96,
		CurrentDataFlows: &chkCurrentDataFlows96,
		FunctionIndex:    0,
		FunctionName:     "decode",
		MaxCapacity:      &chkMaxCapacity96,
		MaxDataFlows:     &chkMaxDataFlows96,
		PartitionName:    "0",
	}},
	MaxCapacity:  &chkMaxCapacity7,
	MaxFunctions: &chkMaxFunctions7,
	Name:         "lane0",
	Type:         "alveo",
}

var chkCurrentFunctions7 int32 = 1
var chkDeviceUUID8 string = "21330621T04L"
var chkMaxFunctions8 int32 = 1
var chkComRes5 = examplecomv1.RegionInfo{
	Available:        false,
	CurrentCapacity:  &chkCurrentCapacity6,
	CurrentFunctions: &chkCurrentFunctions7,
	DeviceFilePath:   "/dev/xpcie_21330621T04L",
	DeviceIndex:      0,
	DeviceType:       "alveo",
	DeviceUUID:       &chkDeviceUUID8,
	MaxCapacity:      &chkMaxCapacity7,
	MaxFunctions:     &chkMaxFunctions8,
	Name:             "lane0",
	Type:             "alveo",
}

var chkComRes6 = examplecomv1.RegionInfo{
	Available:        false,
	CurrentCapacity:  &chkCurrentCapacity6,
	CurrentFunctions: &chkCurrentFunctions7,
	DeviceFilePath:   "/dev/xpcie_21330621T04L",
	DeviceIndex:      0,
	DeviceType:       "alveo",
	DeviceUUID:       &chkDeviceUUID8,
	MaxCapacity:      &chkMaxCapacity7,
	MaxFunctions:     &chkMaxFunctions8,
	Name:             "lane1",
	Type:             "alveo",
}

var chkComRes7 = examplecomv1.RegionInfo{
	Available:        true,
	CurrentCapacity:  &chkCurrentCapacity6,
	CurrentFunctions: &chkCurrentFunctions7,
	DeviceFilePath:   "/dev/xpcie_21330621T04L",
	DeviceIndex:      0,
	DeviceType:       "alveo",
	DeviceUUID:       &chkDeviceUUID8,
	MaxCapacity:      &chkMaxCapacity7,
	MaxFunctions:     &chkMaxFunctions8,
	Name:             "lane0",
	Type:             "alveo",
}

var chkComRes8 = examplecomv1.RegionInfo{
	Available:        true,
	CurrentCapacity:  &chkCurrentCapacity6,
	CurrentFunctions: &chkCurrentFunctions7,
	DeviceFilePath:   "/dev/xpcie_21330621T04L",
	DeviceIndex:      0,
	DeviceType:       "alveo",
	DeviceUUID:       &chkDeviceUUID8,
	MaxCapacity:      &chkMaxCapacity7,
	MaxFunctions:     &chkMaxFunctions8,
	Name:             "lane1",
	Type:             "alveo",
}

var uid types.UID = "aaaaaaa"
var regionname string = "lane0"
var regionname1 string = "lane1"
var typestring string = "decode"
var maxcapachbs1 int32 = 40
var maxcapachbs2 int32 = 40
var maxdataflowschbs1 int32 = 8
var maxdataflowschbs2 int32 = 8
var functionIDchbs1 int32 = 0
var functionIDchbs2 int32 = 1

var modules = examplecomv1.ChildBsModule{
	Functions: &functions,
}
var functions = []examplecomv1.ChildBsFunctions{
	{
		Module:     &module1,
		DeploySpec: deployspec1,
		ID:         &functionIDchbs1,
	}, {
		Module:     &module2,
		DeploySpec: deployspec2,
		ID:         &functionIDchbs2,
	},
}

var module1 = []examplecomv1.FunctionsModule{{
	Type: &typestring,
}}
var module2 = []examplecomv1.FunctionsModule{{
	Type: &typestring,
}}

var deployspec1 = examplecomv1.FunctionsDeploySpec{
	MaxCapacity:  &maxcapachbs1,
	MaxDataFlows: &maxdataflowschbs1,
}
var deployspec2 = examplecomv1.FunctionsDeploySpec{
	MaxCapacity:  &maxcapachbs2,
	MaxDataFlows: &maxdataflowschbs2,
}

var childBsCRdata1 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "ChildBs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "childbs1",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga1",
				UID:        uid,
			},
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name:    &regionname,
				Modules: &modules,
			},
		},
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name: &regionname,
			},
		},
		Status: examplecomv1.ChildBsStatusNotReady,
		State:  examplecomv1.ChildBsWritingBsfile,
	},
}

var childBsCRdata2 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "ChildBs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "childbs2",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga1",
				UID:        uid,
			},
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name:    &regionname,
				Modules: &modules,
			},
		},
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name: &regionname,
			},
		},
		Status: examplecomv1.ChildBsStatusReady,
		State:  examplecomv1.ChildBsReady,
	},
}

var childBsCRdata3 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "ChildBs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "childbs3",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga1",
				UID:        uid,
			},
		},
		Finalizers: []string{
			"fpgafunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name:    &regionname,
				Modules: &modules,
			},
			{
				Name:    &regionname1,
				Modules: &modules,
			},
		},
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name: &regionname,
			},
		},
		Status: examplecomv1.ChildBsStatusPreparing,
		State:  examplecomv1.ChildBsReconfiguring,
	},
}

var childBsCRdata4 = examplecomv1.ChildBs{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "ChildBs",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "childbs4",
		Namespace: "default",
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "example.com/v1",
				Kind:       "FPGA",
				Name:       "fpga1",
				UID:        uid,
			},
		},
		Finalizers: []string{
			"fpgafunction.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.ChildBsSpec{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name:    &regionname,
				Modules: &modules,
			},
			{
				Name:    &regionname1,
				Modules: &modules,
			},
		},
	},
	Status: examplecomv1.ChildBsStatus{
		Regions: []examplecomv1.ChildBsRegion{
			{
				Name: &regionname,
			},
		},
		Status: examplecomv1.ChildBsStatusPreparing,
		State:  examplecomv1.ChildBsReconfiguring,
	},
}
var childBsID string = "aaaaaaa"
var childBsCRName string = ""
var fpgaCRdata = examplecomv1.FPGA{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "FPGA",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fpga1",
		Namespace: "default",
	},
	Spec: examplecomv1.FPGASpec{
		ChildBitstreamID:  &childBsID,
		DeviceIndex:       0,
		DeviceFilePath:    "/dev/xpcie_21330621T04L",
		DeviceUUID:        "21330621T04L",
		NodeName:          "test01",
		ParentBitstreamID: "yyyyyyyyy.mcs",
		PCIDomain:         2,
		PCIBus:            3,
		PCIDevice:         4,
		PCIFunction:       5,
		Vendor:            "eeeeeeeeee",
	},
	Status: examplecomv1.FPGAStatus{
		ChildBitstreamID:     &childBsID,
		ChildBitstreamCRName: &childBsCRName,
		DeviceIndex:          0,
		DeviceFilePath:       "/dev/xpcie_21330621T04L",
		DeviceUUID:           "21330621T04L",
		NodeName:             "test01",
		ParentBitstreamID:    "yyyyyyyyy.mcs",
		PCIDomain:            2,
		PCIBus:               3,
		PCIDevice:            4,
		PCIFunction:          5,
		Vendor:               "eeeeeeeeee",
		Status:               examplecomv1.FPGAStatusPreparing,
	},
}

var infrainfo_configdata = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "infrastructureinfo",
		Namespace: "default",
	},
	TypeMeta: metav1.TypeMeta{
		Kind: "ConfigMap",
	},
	Data: map[string]string{
		"infrastructureinfo.json": `
		{"devices": [{
			"deviceFilePath": "/dev/xpcie_21330621T04L",
			"nodeName": "test01",
			"deviceUUID": "21330621T04L",
			"deviceIndex": 0,
			"deviceType": "alveo"
		},{
			"deviceFilePath": "/dev/xpcie_21330621T01J",
			"nodeName": "test01",
			"deviceUUID": "21330621T01J",
			"deviceIndex": 1,
			"deviceType": "alveo"
		},{
			"deviceFilePath": "",
			"nodeName": "test01",
			"deviceUUID": "gpu-123456789t4",
			"deviceIndex": 0,
			"deviceType": "t4"
		},{
			"deviceFilePath": "",
			"nodeName": "test01",
			"deviceUUID": "gpu-123456789a100",
			"deviceIndex": 1,
			"deviceType": "a100"
		},{
			"deviceFilePath": "",
			"nodeName": "test01",
			"deviceUUID": "test01-cpu",
			"deviceIndex": 0,
			"deviceType": "cpu"
		}]}`,
	},
}

var deployinfo_configdata = corev1.ConfigMap{
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
				"regionType": "alveo",
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
				"regionType": "alveo",
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
		{"devices":[{
			"nodeName":"test01",
			"deviceFilePath": "/dev/xpcie_21330621T04L",
			"deviceUUID": "21330621T04L",
				"functionTargets":[{
					"regionType":"alveo",
					"regionName":"lane0",
					"maxFunctions":1,
					"maxCapacity":40,
					"functions":[{
						"functionIndex":0,
						"partitionName":"0",
						"functionName":"filter-resize-high-high-infer",
						"maxDataFlows":8,
						"maxCapacity":40
					}]
				},{
					"regionType":"alveo",
					"regionName":"lane1",
					"maxFunctions":1,
					"maxCapacity":40,
					"functions":[{
						"functionIndex":0,
						"partitionName":"1",
						"functionName":"filter-resize-low-low-infer",
						"maxDataFlows":8,
						"maxCapacity":40
					}]
				}]
			},{
			"nodeName":"test01",
			"deviceFilePath": "/dev/xpcie_21330621T01J",
			"deviceUUID": "21330621T01J",
			"functionTargets":[{
				"regionType":"alveo",
				"regionName":"lane0",
				"maxFunctions":2,
				"maxCapacity":40,
				"functions":[{
					"functionIndex":0,
					"partitionName":"0",
					"functionName":"decode",
					"maxDataFlows":6,
					"maxCapacity":20
				},{
					"functionIndex":1,
					"partitionName":"1",
					"functionName":"decode",
					"maxDataFlows":6,
					"maxCapacity":20
				}]
			},{
				"regionType":"alveo",
				"regionName":"lane1",
				"maxFunctions":2,
				"maxCapacity":40,
					"functions":[{
						"functionIndex":0,
						"partitionName":"2",
						"functionName":"decode",
						"maxDataFlows":6,
						"maxCapacity":20
					},{
						"functionIndex":1,
						"partitionName":"3",
						"functionName":"decode",
						"maxDataFlows":6,
						"maxCapacity":20
					}]
				}]
			},{
			"nodeName":"test01",
			"deviceUUID":"GPU-1234567890ab0",
			"functionTargets":[{
				"regionType":"t4",
				"regionName":"t4",
				"maxFunctions":110,
				"maxCapacity":40
			}]
		},{
			"nodeName":"test01",
				"deviceUUID":"GPU-1234567890ab1",
				"functionTargets":[{
					"regionType":"a100",
					"regionName":"a100",
					"maxFunctions":110,
					"maxCapacity":120
				}]
			},{
				"nodeName":"test01",
				"deviceUUID":"test01-cpu0",
				"functionTargets":[{"regionType":"cpu",
					"regionName":"cpu",
					"maxFunctions":220,
					"maxCapacity":240
				}]
		}]}`,
	},
}

var deployinfo_configdata3 = corev1.ConfigMap{
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
			"functionTargets":[{
					"regionType":"alveo",
					"regionName":"lane0",
					"maxFunctions":1,
					"maxCapacity":40
				},{
					"regionType":"alveo",
					"regionName":"lane1",
					"maxFunctions":1,
					"maxCapacity":40
				}]
			},{
			"nodeName":"test01",
			"deviceFilePath": "/dev/xpcie_21330621T01J",
			"deviceUUID": "21330621T01J",
			"functionTargets":[{
				"regionType":"alveo",
				"regionName":"lane0",
				"maxFunctions":2,
				"maxCapacity":40,
				"functions":[{
					"functionIndex":0,
					"partitionName":"0",
					"functionName":"decode",
					"maxDataFlows":6,
					"maxCapacity":20
				},{
					"functionIndex":1,
					"partitionName":"1",
					"functionName":"decode",
					"maxDataFlows":6,
					"maxCapacity":20
				}]
			},{
				"regionType":"alveo",
				"regionName":"lane1",
				"maxFunctions":2,
				"maxCapacity":40,
					"functions":[{
						"functionIndex":0,
						"partitionName":"2",
						"functionName":"decode",
						"maxDataFlows":6,
						"maxCapacity":20
					},{
						"functionIndex":1,
						"partitionName":"3",
						"functionName":"decode",
						"maxDataFlows":6,
						"maxCapacity":20
					}]
				}]
			},{
			"nodeName":"test01",
			"deviceUUID":"GPU-1234567890ab0",
			"functionTargets":[{
				"regionType":"t4",
				"regionName":"t4",
				"maxFunctions":110,
				"maxCapacity":40
			}]
		},{
			"nodeName":"test01",
				"deviceUUID":"GPU-1234567890ab1",
				"functionTargets":[{
					"regionType":"a100",
					"regionName":"a100",
					"maxFunctions":110,
					"maxCapacity":120
				}]
			},{
				"nodeName":"test01",
				"deviceUUID":"test01-cpu0",
				"functionTargets":[{"regionType":"cpu",
					"regionName":"cpu",
					"maxFunctions":220,
					"maxCapacity":240
				}]
		}]}`,
	},
}

/*
=====================================================
7-9-1 UPDATE
=====================================================
*/

var deviceInfoUpdate = examplecomv1.DeviceInfo{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "DeviceInfo",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "deviceinfo-update-wbfunction-decode-main",
		Namespace: "default",
		Finalizers: []string{
			"deviceinfo.finalizers.example.com.v1",
		},
	},
	Spec: examplecomv1.DeviceInfoSpec{
		Request: examplecomv1.WBFuncRequest{
			RequestType:   "Deploy",
			DeviceType:    "cpu",
			DeviceIndex:   0,
			RegionName:    "cpu",
			NodeName:      "test01",
			FunctionName:  "cpu-decode",
			FunctionIndex: &FunctionIndex,
			MaxDataFlows:  &MaxDataFlows,
			MaxCapacity:   &MaxCapacity,
			Capacity:      &Capacity,
		},
	},

	Status: examplecomv1.DeviceInfoStatus{
		Response: examplecomv1.WBFuncResponse{
			Status:        "Deployed",
			FunctionIndex: &FunctionIndex,
			DeviceUUID:    "k8s-worker-cpu0",
		},
	},
}
