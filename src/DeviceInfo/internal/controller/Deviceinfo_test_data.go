package controller

import (
	examplecomv1 "DeviceInfo/api/v1"
	// k8scnicncfio "github.com/k8snetworkplumbingwg/network-attachment-definition-client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// runtime "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/util/intstr"
	// "sigs.k8s.io/controller-runtime/pkg/scheme"
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
	/*
		Status: examplecomv1.DeviceInfoStatus{
			Response: examplecomv1.WBFuncResponse{
				Status:        "Deployed",
				FunctionIndex: &FunctionIndex,
				DeviceUUID:    "k8s-worker-cpu0",
			},
		},
	*/
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
	/*
	   Status: examplecomv1.DeviceInfoStatus{
	       Response: examplecomv1.WBFuncResponse{
	           Status:        "Deployed",
	           FunctionIndex: &FunctionIndex,
	           DeviceUUID:    "k8s-worker-cpu0",
	       },
	   },
	*/
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
	/*
	   Status: examplecomv1.DeviceInfoStatus{
	       Response: examplecomv1.WBFuncResponse{
	           Status:        "Deployed",
	           FunctionIndex: &FunctionIndex,
	           DeviceUUID:    "k8s-worker-cpu0",
	       },
	   },
	*/
}

var CurrentCapacity int32 = 30
var CurrentCapacity2 int32 = 15
var CurrentCapacity3 int32 = 0
var CurrentCapacity4 int32 = 120
var CurrentFunctions int32 = 2
var CurrentFunctions2 int32 = 1
var CurrentFunctions3 int32 = 0
var CurrentFunctions4 int32 = 8
var CurrentDataFlows int32 = 1
var CurrentDataFlows2 int32 = 2
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
	/*
		TypeMeta: metav1.TypeMeta{
			APIVersion: "example.com/v1",
			Kind:       "ComputeResource",
		},
	*/
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

/*
var chkComRes1 = examplecomv1.ComputeResource{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "compute-test01",
		Namespace: "default",
	},
	Spec: examplecomv1.ComputeResourceSpec{
		NodeName: "test01",
		Regions: []examplecomv1.RegionInfo{
			{
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "0",
				}, {
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity2,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "1",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxCapacity2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity2,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "2",
				}, {
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity2,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "3",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxFunctions,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity,
					CurrentDataFlows: &chkCurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &chkMaxCapacity2,
					MaxDataFlows:     &chkMaxDataFlows2,
					PartitionName:    "0",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxFunctions2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity,
					CurrentDataFlows: &chkCurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &chkMaxCapacity2,
					MaxDataFlows:     &chkMaxDataFlows2,
					PartitionName:    "1",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxDataFlows2,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity3,
				CurrentFunctions: &chkCurrentFunctions3,
				DeviceFilePath:   "",
				DeviceIndex:      0,
				DeviceType:       "t4",
				DeviceUUID:       &chkDeviceUUID3,
				MaxCapacity:      &chkMaxCapacity2,
				MaxFunctions:     &chkMaxFunctions3,
				Name:             "t4",
				Type:             "t4",
			}, {
				Available:        false,
				CurrentCapacity:  &chkCurrentCapacity4,
				CurrentFunctions: &chkCurrentFunctions4,
				DeviceFilePath:   "",
				DeviceIndex:      1,
				DeviceType:       "a100",
				DeviceUUID:       &chkDeviceUUID4,
				MaxCapacity:      &chkMaxCapacity3,
				MaxFunctions:     &chkMaxFunctions3,
				Name:             "a100",
				Type:             "a100",
			}, {
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
					PartitionName:    "0",
				}},
				MaxCapacity:      &chkMaxCapacity4,
				MaxFunctions:     &chkMaxFunctions4,
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
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "0",
				}, {
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity2,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "1",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxCapacity2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions,
				DeviceFilePath:   "/dev/xpcie_21330621T01J",
				DeviceIndex:      1,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity2,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    0,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "2",
				}, {
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity2,
					CurrentDataFlows: &chkCurrentDataFlows,
					FunctionIndex:    1,
					FunctionName:     "decode",
					MaxCapacity:      &chkMaxCapacity,
					MaxDataFlows:     &chkMaxDataFlows1,
					PartitionName:    "3",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxFunctions,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity,
					CurrentDataFlows: &chkCurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &chkMaxCapacity2,
					MaxDataFlows:     &chkMaxDataFlows2,
					PartitionName:    "0",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxFunctions2,
				Name:         "lane0",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity,
				CurrentFunctions: &chkCurrentFunctions2,
				DeviceFilePath:   "/dev/xpcie_21330621T04L",
				DeviceIndex:      0,
				DeviceType:       "alveo",
				DeviceUUID:       &chkDeviceUUID2,
				Functions: []examplecomv1.FunctionInfrastruct{{
					Available:        true,
					CurrentCapacity:  &chkCurrentCapacity,
					CurrentDataFlows: &chkCurrentDataFlows2,
					FunctionIndex:    0,
					FunctionName:     "filter-resize-high-infer",
					MaxCapacity:      &chkMaxCapacity2,
					MaxDataFlows:     &chkMaxDataFlows2,
					PartitionName:    "1",
				}},
				MaxCapacity:  &chkMaxCapacity2,
				MaxFunctions: &chkMaxDataFlows2,
				Name:         "lane1",
				Type:         "alveo",
			}, {
				Available:        true,
				CurrentCapacity:  &chkCurrentCapacity3,
				CurrentFunctions: &chkCurrentFunctions3,
				DeviceFilePath:   "",
				DeviceIndex:      0,
				DeviceType:       "t4",
				DeviceUUID:       &chkDeviceUUID3,
				MaxCapacity:      &chkMaxCapacity2,
				MaxFunctions:     &chkMaxFunctions3,
				Name:             "t4",
				Type:             "t4",
			}, {
				Available:        false,
				CurrentCapacity:  &chkCurrentCapacity4,
				CurrentFunctions: &chkCurrentFunctions4,
				DeviceFilePath:   "",
				DeviceIndex:      1,
				DeviceType:       "a100",
				DeviceUUID:       &chkDeviceUUID4,
				MaxCapacity:      &chkMaxCapacity3,
				MaxFunctions:     &chkMaxFunctions3,
				Name:             "a100",
				Type:             "a100",
			}, {
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
					PartitionName:    "0",
				}},
				MaxCapacity:      &chkMaxCapacity4,
				MaxFunctions:     &chkMaxFunctions4,
				Name:             "cpu",
				Type:             "cpu",
			},
		},
	},
}
*/

var uid types.UID = "aaaaaaa"
var regionname string = "lane0"
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
	//	apiVersion: "example.com/v1",
	//	Kind:       "ChildBs",
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
	//	apiVersion: "example.com/v1",
	//	Kind:       "ChildBs",
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
		Status: examplecomv1.ChildBsStatusReady,
		State:  examplecomv1.ChildBsReady,
	},
}

var childBsID string = "aaaaaaa"
var childBsCRName string = ""
var fpgaCRdata = examplecomv1.FPGA{
	//	apiVersion: "example.com/v1",
	//	Kind:       "ChildBs",
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
	/*
		Status: examplecomv1.FPGAStatus{
			Regions: []examplecomv1.ChildBsRegion{
				{
					Name: &regionname,
				},
			},
			Status: examplecomv1.ChildBsStatusNotReady,
			State:  examplecomv1.ChildBsWritingBsfile,
		},
	*/
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

/*
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
*/

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
