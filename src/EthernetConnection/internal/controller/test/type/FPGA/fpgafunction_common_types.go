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
	FPGADeviceConfigName = "fpgalist-ph3"
)

type FPGAFuncConfig struct {
	ParentBitstream string            `json:"parentBitstream,omitempty"`
	ChildBitstream  string            `json:"childBitstream,omitempty"`
	SharedMemoryGiB int32             `json:"sharedMemoryGiB,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	Commands        []string          `json:"commands,omitempty"`
	Args            []string          `json:"args,omitempty"`
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
