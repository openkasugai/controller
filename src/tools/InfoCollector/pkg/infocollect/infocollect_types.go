/*
Copyright 2024 NTT Corporation, FUJITSU LIMITED

Structures to be stored internally
*/

package infocollect

import ()

// Node & device information
type DeviceInfo struct {
	DeviceFilePath *string `json:"deviceFilePath,omitempty"`
	NodeName       string  `json:"nodeName"`
	DeviceUUID     *string `json:"deviceUUID,omitempty"`
	DeviceType     string  `json:"deviceType"`
	DeviceIndex    int32   `json:"deviceIndex"`
}

// Structure for storing device information
type DeviceBasicInfo struct {
	NodeName         string  `json:"nodeName"`
	DeviceFilePath   *string `json:"deviceFilePath,omitempty"`
	DeviceUUID       *string `json:"deviceUUID,omitempty"`
	DeviceType       string  `json:"deviceType"`
	DeviceIndex      int32   `json:"deviceIndex"`
	ParentID         string  `json:"parentID"`
	SubDeviceSpecRef string  `json:"subDeviceSpecRef"`
	DeviceVendor     *string `json:"deviceVendor,omitempty"`
	PCIDomain        *int32  `json:"pciDomain,omitempty"`
	PCIBus           *int32  `json:"pciBus,omitempty"`
	PCIDevice        *int32  `json:"pciDevice,omitempty"`
	PCIFunction      *int32  `json:"pciFunction,omitempty"`
	DeviceCategory   string  `json:"deviceCategory"`
}

// Device deployment information
type DeviceRegionInfo struct {
	NodeName         string           `json:"nodeName"`
	DeviceFilePath   *string          `json:"deviceFilePath,omitempty"`
	DeviceUUID       *string          `json:"deviceUUID,omitempty"`
	SubDeviceSpecRef string           `json:"subDeviceSpecRef"`
	FunctionTargets  []RegionInDevice `json:"functionTargets"`
}

// Device deployment information, domain-specific information - FunctionTarget
type RegionInDevice struct {
	RegionName   string                       `json:"regionName"`
	RegionType   string                       `json:"regionType"`
	MaxFunctions *int32                       `json:"maxFunctions,omitempty"`
	MaxCapacity  *int32                       `json:"maxCapacity,omitempty"`
	Functions    *[]SimpleFunctionInfraStruct `json:"functions,omitempty"`
}

// Device deployment information - Function
type SimpleFunctionInfraStruct struct {
	FunctionName      string `json:"functionName"`
	FunctionIndex     int32  `json:"functionIndex"`
	FrameworkKernelID int32  `json:"frameworkKernelID"`
	PartitionName     string `json:"partitionName"`
}

// Domain specific information
type RegionSpecificInfo struct {
	SubDeviceSpecRef string           `json:"subDeviceSpecRef"`
	FunctionTargets  []RegionInDevice `json:"functionTargets"`
}

// Function specific information
type FPGACatalog struct {
	FunctionID   int32  `json:"functionID"`
	FunctionName string `json:"functionName"`
	MaxDataFlows int32  `json:"maxDataFlows"`
	MaxCapacity  int32  `json:"maxCapacity"`
}

// Function-specific information
type FunctionDedicatedInfo struct {
	PartitionName         string           `json:"partitionName"`
	FunctionChannelIDList []int32          `json:"functionChannelIDList"`
	FunctionChannelIDs    []FunctionDetail `json:"functionChannelIDs"`
}

// Function-specific information - FuncCHId
type FunctionDetail struct {
	FunctionChannelID int32                 `json:"functionChannelID"`
	Rx                FPGAConnectionCatalog `json:"rx"`
	Tx                FPGAConnectionCatalog `json:"tx"`
}

// Function-specific information - RxTx
type FPGAConnectionCatalog struct {
	Protocol map[string]FPGAConnectionCatalogDetails `json:"protocol"`
}

type FPGAConnectionCatalogDetails struct {
	Port             *int32 `json:"port,omitempty"`
	DMAChannelID     *int32 `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32 `json:"lldmaConnectorID,omitempty"`
}

// Bitstream information
type BsConfigInfo struct {
	ChildBitstreamIDs []ChildBsConfigInfo `json:"child-bitstream-ids"`
	ParentBitstreamID string              `json:"parent-bitstream-id"`
}

// ChildBitstream information
type ChildBsConfigInfo struct {
	Regions          []ChildBsRegion `json:"regions"`
	ChildBitstreamID *string         `json:"child-bitstream-id,omitempty"`
}

type ChildBsRegion struct {
	Modules *ChildBsModule `json:"modules,omitempty"`
	Name    *string        `json:"name,omitempty"`
}

type ChildBsModule struct {
	Ptu         *ChildBsPtu         `json:"ptu,omitempty"`
	LLDMA       *ChildBsLLDMA       `json:"lldma,omitempty"`
	Chain       *ChildBsChain       `json:"chain,omitempty"`
	Directtrans *ChildBsDirecttrans `json:"directtrans,omitempty"`
	Conversion  *ChildBsConversion  `json:"conversion,omitempty"`
	Functions   *[]ChildBsFunctions `json:"functions,omitempty"`
}

type ChildBsPtu struct {
	Cids *string `json:"cids,omitempty"`
	ID   *int32  `json:"id,omitempty"`
	// ExtIfID *int32  `json:"extIfId,omitempty"`
}

type ChildBsLLDMA struct {
	Cids *string `json:"cids,omitempty"`
	ID   *int32  `json:"id,omitempty"`
	// ExtIfID *int32  `json:"extIfId,omitempty"`
}

type ChildBsChain struct {
	ID         *int32  `json:"id,omitempty"`
	Identifier *string `json:"identifier,omitempty"`
	Type       *string `json:"type,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type ChildBsDirecttrans struct {
	ID         *int32  `json:"id,omitempty"`
	Identifier *string `json:"identifier,omitempty"`
	Type       *string `json:"type,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type ChildBsConversion struct {
	ID     *int32              `json:"id,omitempty"`
	Module *[]ConversionModule `json:"module,omitempty"`
}

type ChildBsFunctions struct {
	ID     *int32             `json:"id,omitempty"`
	Module *[]FunctionsModule `json:"module,omitempty"`
}

type ConversionModule struct {
	Identifier *string `json:"identifier,omitempty"`
	Type       *string `json:"type,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type FunctionsModule struct {
	FunctionChannelIDs *string `json:"function-channel-ids,omitempty"`
	Identifier         *string `json:"identifier,omitempty"`
	Type               *string `json:"type,omitempty"`
	Version            *string `json:"version,omitempty"`
}

/* Provisional support (dynamic reconfiguration not supported) */
// FunctionName conversion information
type FunctionNameMap struct {
	Height       uint32 `json:"height"`
	Width        uint32 `json:"width"`
	FunctionName string `json:"functionName"`
}

/* Provisional support (dynamic reconfiguration not supported) */
// DeviceType conversion information
type DeviceTypeMap struct {
	InputDeviceType  string `json:"inputDeviceType"`
	OutputDeviceType string `json:"outputDeviceType"`
}

// Frame size setting information
type FrameSizeConfig struct {
	InputWidth   uint32 `json:"i_width"`
	InputHeight  uint32 `json:"i_height"`
	OutputWidth  uint32 `json:"o_width"`
	OutputHeight uint32 `json:"o_height"`
}

// Structure for ServicerMgmtInfo information acquisition
type ServicerMgmtInfo struct {
	NodeName    string        `json:"nodeName"`
	NetworkInfo []NetworkData `json:"networkInfo"`
}

type NetworkData struct {
	DeviceFilePath string `json:"deviceFilePath"`
	LaneIndex      int32  `json:"laneIndex"`
	IPAddress      string `json:"ipAddress"`
	SubnetAddress  string `json:"subnetAddress"`
	GatewayAddress string `json:"gatewayAddress"`
	MACAddress     string `json:"macAddress"`
}
