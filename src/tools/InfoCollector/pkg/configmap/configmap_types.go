/* Copyright 2024 NTT Corporation , FUJITSU LIMITED */

package configmap

import ()

// Event type
const (
	CreateEvent = iota
	UpdateEvent
	DeleteEvent
)

// The type of ConfigMap to create.
const (
	CMInfraInfo        = "infrastructureinfo"
	CMDeployInfo       = "deployinfo"
	CMFPGACatalog      = "fpgacatalogmap"
	CMFPGADecode       = "decode-ch"
	CMFPGAFilterResize = "filter-resize-ch"
)

// Information for creating ConfigMapCR
const (
	CMNameSpace  = "default"
	CMKind       = "ConfigMap"
	CMAPIVersion = "v1"
)

// Store node name information
var GMyNodeName string

// Config information storage area
var GInfrastructureInfo map[string][]DeviceInfo
var GDeployInfo map[string][]DeviceRegionInfo
var GFPGACatalogMap map[string][]FPGACatalog
var GDecodeCH []FunctionDetail
var GFilterResizeCH []FunctionDetail

// Infrastructure information ConfigMap
type DeviceInfo struct {
	DeviceFilePath *string `json:"deviceFilePath,omitempty"`
	NodeName       string  `json:"nodeName"`
	DeviceUUID     *string `json:"deviceUUID,omitempty"`
	DeviceType     string  `json:"deviceType"`
	DeviceIndex    int32   `json:"deviceIndex"`
}

// Deployment information ConfigMap
type DeviceRegionInfo struct {
	NodeName        string           `json:"nodeName"`
	DeviceFilePath  *string          `json:"deviceFilePath,omitempty"`
	DeviceUUID      *string          `json:"deviceUUID,omitempty"`
	FunctionTargets []RegionInDevice `json:"functionTargets"`
}

// Deployment information FunctionTargets
type RegionInDevice struct {
	RegionType   string                      `json:"regionType"`
	RegionName   string                      `json:"regionName"`
	MaxFunctions *int32                      `json:"maxFunctions,omitempty"`
	MaxCapacity  *int32                      `json:"maxCapacity,omitempty"`
	Functions    []SimpleFunctionInfraStruct `json:"functions,omitempty"`
}

// Deployment informationFunctions
type SimpleFunctionInfraStruct struct {
	FunctionIndex *int32 `json:"functionIndex,omitempty"`
	PartitionName string `json:"partitionName"`
	FunctionName  string `json:"functionName"`
	MaxDataFlows  int32  `json:"maxDataFlows"`
	MaxCapacity   int32  `json:"maxCapacity"`
}

// FPGACatalogMap
type FPGACatalog struct {
	NodeName       string          `json:"nodeName"`
	DeviceFilePath string          `json:"deviceFilePath"`
	DeviceUUID     string          `json:"deviceUUID"`
	Details        []DeviceDetails `json:"details"`
}

type DeviceDetails struct {
	RegionName string           `json:"regionName"`
	IPAddress  string           `json:"ipAddress"`
	FuncData   []FuncDataStruct `json:"funcData"`
}

type FuncDataStruct struct {
	FunctionIndex      int32   `json:"functionIndex"`
	FunctionKernelID   int32   `json:"functionKernelID"`
	FrameworkKernelID  int32   `json:"frameworkKernelID"`
	FunctionChannelIDs []int32 `json:"functionChannelIDs"`
}

// Detailed information about each function
type FunctionDedicatedInfo struct {
	PartitionName      string           `json:"partitionName"`
	FunctionChannelIDs []FunctionDetail `json:"functionChannelIDs"`
}

type FunctionDetail struct {
	FunctionChannelID int32                 `json:"functionChannelID"`
	Rx                FPGAConnectionCatalog `json:"rx"`
	Tx                FPGAConnectionCatalog `json:"tx"`
}

type FPGAConnectionCatalog struct {
	Protocol map[string]FPGAConnectionCatalogDetails `json:"protocol"`
}

type FPGAConnectionCatalogDetails struct {
	Port             *int32 `json:"port,omitempty"`
	DMAChannelID     *int32 `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32 `json:"lldmaConnectorID,omitempty"`
}
