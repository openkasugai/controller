/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package v1

/*
import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
*/

const (
	RequestMemory = "32GiB"
)

const (
	FunctionCRKindFPGA = "FPGAFunction"
	FunctionCRKindGPU  = "GPUFunction"
	FunctionCRKindGATE = "GateFunction"
	FunctionCRKindCPU  = "CPUFunction"
)

const (
	ConnectionCRKindPCIe = "PCIeConnection"
	ConnectionCRKindEth  = "EthernetConnection"
)

const (
	InputIP    = "inputIPAddress"
	InputPort  = "inputPort"
	InputMAC   = "inputMACAddress"
	OutputIP   = "outputIPAddress"
	OutputPort = "outputPort"
	OutputMAC  = "outputMACAddress"
)

const (
	Width            = "WIDTH"
	Height           = "HEIGHT"
	ArgsWidth        = "%" + Width + "%"
	ArgsHeight       = "%" + Height + "%"
	ArgsArpIP        = "%ARPIP%"
	ArgsIP           = "%IP%"
	ArgsPort         = "%PORT%"
	ArgsArpMAC       = "%MAC%"
	ChangeArgsWidth  = "width="
	ChangeArgsHeight = "height="
	ChangeArgsIP     = "host="
	ChangeArgsPort   = "port="
)

const (
	CMServicerMgmtInfo   = "servicer-mgmt-info"
	CMDeployInfo         = "deployinfo"
	CMRegionUniqueInfo   = "region-unique-info"
	CMFunctionUniqueInfo = "function-unique-info"
	CMFilterResizeInfo   = "filter-resize-ch"
)

const (
	CMNameSpace = "default"
)

type FPGAFuncConfig struct {
	ParentBitstream BitstreamData     `json:"parentBitstream,omitempty"`
	ChildBitstream  BitstreamData     `json:"childBitstream,omitempty"`
	SharedMemoryGiB int32             `json:"sharedMemoryGiB,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	Commands        []string          `json:"commands,omitempty"`
	Args            []string          `json:"args,omitempty"`
	Parameters      AnyData           `json:"parameters,omitempty"`
}

type BitstreamData struct {
	File string `json:"file,omitempty"`
	ID   string `json:"id,omitempty"`
}

type AnyData struct {
	Functions FrameSizeData `json:"functions,omitempty"`
}

type FrameSizeData struct {
	InputWidth   int32 `json:"i_width"`
	InputHeight  int32 `json:"i_height"`
	OutputWidth  int32 `json:"o_width"`
	OutputHeight int32 `json:"o_height"`
}

// Function-based CR structure
type FunctionData struct {
	DataFlowRef       WBNamespacedName  `json:"dataFlowRef"`
	FunctionName      string            `json:"functionName"`
	NodeName          string            `json:"nodeName"`
	DeviceType        string            `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo       `json:"acceleratorIDs"`
	RegionName        string            `json:"regionName"`
	FunctionIndex     *int32            `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo        `json:"envs,omitempty"`
	RequestMemorySize *int32            `json:"requestMemorySize,omitempty"`
	SharedMemory      *SharedMemorySpec `json:"sharedMemory,omitempty"`
	Protocol          *string           `json:"protocol,omitempty"`
	ConfigName        string            `json:"configName"`
	FunctionKernelID  int32             `json:"functionKernelID,omitempty"`
	FunctionChannelID int32             `json:"functionChannelID,omitempty"`
	PtuKernelID       int32             `json:"ptuKernelID,omitempty"`
	FrameworkKernelID int32             `json:"frameworkKernelID,omitempty"`
	Rx                RxTxSpecFunc      `json:"rx,omitempty"`
	Tx                RxTxSpecFunc      `json:"tx,omitempty"`
}

type RxTxSpecFunc struct {
	Protocol        string  `json:"protocol"`
	IPAddress       *string `json:"iPAddress,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	SubnetAddress   *string `json:"subnetAddress,omitempty"`
	GatewayAddress  *string `json:"gatewayAddress,omitempty"`
	DMAChannelID    *int32  `json:"dmaChannelID,omitempty"`
	FDMAConnectorID *int32  `json:"fdmaConnectorID,omitempty"`
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
	MaxFunctions int32                       `json:"maxFunctions"`
	MaxCapacity  int32                       `json:"maxCapacity"`
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

// Detailed information about each function
type FunctionDetail struct {
	FunctionChannelID int32              `json:"functionChannelID"`
	Rx                FPGACatalogmapRxTx `json:"rx"`
	Tx                FPGACatalogmapRxTx `json:"tx"`
}

// FPGACatalogMapDetails
type FPGACatalogmapRxTx struct {
	Protocol *map[string]Details `json:"protocol"`
}

/*
type Details struct {
	Port             *int32 `json:"port,omitempty"`
	DMAChannelID     *int32 `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32 `json:"lldmaConnectorID,omitempty"`
}
*/
