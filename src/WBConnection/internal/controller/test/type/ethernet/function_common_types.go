/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package controllertestethernet

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

// Function CR structure
type FunctionData struct {
	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     int32                         `json:"functionIndex"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      SharedMemorySpec              `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        *string                       `json:"configName,omitempty"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
	FunctionKernelID  *int32                        `json:"functionKernelID,omitempty"`
	FunctionChannelID *int32                        `json:"functionChannelID,omitempty"`
	PtuKernelID       *int32                        `json:"ptuKernelID,omitempty"`
	FrameworkKernelID *int32                        `json:"frameworkKernelID,omitempty"`
	Rx                RxTxSpecFunc                  `json:"rx,omitempty"`
	Tx                RxTxSpecFunc                  `json:"tx,omitempty"`
}

// Function CR structure
type FunctionStatusData struct {
	Status string `json:"status"`
}
type AccIDInfo struct {
	PartitionName string `json:"partitionName"`
	ID            string `json:"id"`
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

// FPGA Device information
type RxTxSpecFunc struct {
	Protocol        string  `json:"protocol,omitempty"`
	IPAddress       *string `json:"ipAddress,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	SubnetAddress   *string `json:"subnetAddress,omitempty"`
	GatewayAddress  *string `json:"gatewayAddress,omitempty"`
	DMAChannelID    *int32  `json:"dmaChannelID,omitempty"`
	FDMAConnectorID *int32  `json:"fdmaConnectorID,omitempty"`
}

// Shared memory information
type SharedMemorySpec struct {
	FilePrefix      string `json:"filePrefix"`
	CommandQueueID  string `json:"commandQueueID"`
	SharedMemoryMiB int32  `json:"sharedMemoryMiB"`
}

// Structure for acquiring Phase3FPGA information
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

type FunctionConfigMap struct {
	RxProtocol      string            `json:"rxProtocol,omitempty"`
	TxProtocol      string            `json:"txProtocol,omitempty"`
	SharedMemoryMiB int32             `json:"sharedMemoryMiB,omitempty"`
	ImageURI        string            `json:"imageURI,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	ParentBitStream string            `json:"parentBitStream,omitempty"`
	ChildBitStream  string            `json:"childBitStream,omitempty"`
}
