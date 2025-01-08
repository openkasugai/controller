/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package v1

import (
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	RequestMemory = 32
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
	MyIP             = "ipAddress"
	InputIP          = "inputIPAddress"
	InputPort        = "inputPort"
	InputMAC         = "inputMACAddress"
	OutputIP         = "outputIPAddress"
	OutputPort       = "outputPort"
	OutputMAC        = "outputMACAddress"
	BranchOutputIP   = "branchOutputIPAddress"
	BranchOutputPort = "branchOutputPort"
	GlueOutputIP     = "glueOutputIPAddress"
	GlueOutputPort   = "glueOutputPort"
	FPS              = "decEnvFrameFPS"
)

const (
	FPGADeviceConfigName = "fpgalist-ph3"
)

// Function-based CR structure
type FunctionData struct {
	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     *int32                        `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      SharedMemorySpec              `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        *string                       `json:"configName,omitempty"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
}

// Function-based CR structure
type FunctionStatusData struct {
	Status            string           `json:"status"`
	FunctionIndex     *int32           `json:"functionIndex"`
	FunctionKernelID  *int32           `json:"functionKernelID,omitempty"`
	FunctionChannelID *int32           `json:"functionChannelID,omitempty"`
	PtuKernelID       *int32           `json:"ptuKernelID,omitempty"`
	FrameworkKernelID *int32           `json:"frameworkKernelID,omitempty"`
	Rx                RxTxData         `json:"rx,omitempty"`
	Tx                RxTxData         `json:"tx,omitempty"`
	SharedMemory      SharedMemorySpec `json:"sharedMemory,omitempty"`
}

// FPGADevice Connection Info
type RxTxData struct {
	Protocol         string  `json:"protocol,omitempty"`
	IPAddress        *string `json:"ipAddress,omitempty"`
	Port             *int32  `json:"port,omitempty"`
	SubnetAddress    *string `json:"subnetAddress,omitempty"`
	GatewayAddress   *string `json:"gatewayAddress,omitempty"`
	DMAChannelID     *int32  `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32  `json:"lldmaConnectorID,omitempty"`
}

type AccIDInfo struct {
	PartitionName *string `json:"partitionName,omitempty"`
	ID            string  `json:"id"`
}

type EnvsInfo struct {
	PartitionName string     `json:"partitionName"`
	EachEnv       []EnvsData `json:"eachEnv"`
}

// Environmental information
type EnvsData struct {
	EnvKey   string `json:"envKey"`
	EnvValue string `json:"envValue"`
}

// Shared memory information
type SharedMemorySpec struct {
	FilePrefix      string `json:"filePrefix"`
	CommandQueueID  string `json:"commandQueueID"`
	SharedMemoryMiB int32  `json:"sharedMemoryMiB"`
}

// Structure for Phase3FPGA information acquisition
type Phase3Data struct {
	NodeName        string        `json:"nodeName"`
	DeviceFilePaths []string      `json:"deviceFilePaths"`
	NetworkInfo     []NetworkData `json:"networkInfo"`
}

type NetworkData struct {
	DeviceIndex    int32    `json:"deviceIndex"`
	LaneIndex      int32    `json:"laneIndex"`
	IPAddress      string   `json:"ipAddress"`
	SubnetAddress  string   `json:"subnetAddress"`
	GatewayAddress string   `json:"gatewayAddress"`
	MACAddress     string   `json:"macAddress"`
	Rx             FPGARxTx `json:"rx"`
	Tx             FPGARxTx `json:"tx"`
}

type FPGARxTx struct {
	Protocol  string `json:"protocol"`
	StartPort int32  `json:"startPort,omitempty"`
	EndPort   int32  `json:"endPort,omitempty"`
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

type BitstreamData struct {
	File string `json:"file,omitempty"`
	ID   string `json:"id,omitempty"`
}

type FunctionConfigMap struct {
	RxProtocol      string            `json:"rxProtocol,omitempty"`
	TxProtocol      string            `json:"txProtocol,omitempty"`
	SharedMemoryMiB int32             `json:"sharedMemoryMiB,omitempty"`
	ImageURI        string            `json:"imageURI,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	//	ParentBitstream string            `json:"parentBitstream,omitempty"`
	//	ChildBitstream  string            `json:"childBitstream,omitempty"`
	ParentBitstream       BitstreamData `json:"parentBitstream,omitempty"`
	ChildBitstream        BitstreamData `json:"childBitstream,omitempty"`
	Commands              []string      `json:"commands,omitempty"`
	Args                  []string      `json:"args,omitempty"`
	Parameters            AnyData       `json:"parameters,omitempty"`
	FunctionDedicatedInfo string        `json:"functionDedicatedInfo,omitempty"`
}
