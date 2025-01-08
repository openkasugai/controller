/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package v1

type WBDeployStatus string
type WBIOType string

type WBNamespacedName struct {

	//+kubebuilder:validation:Required

	Namespace string `json:"namespace"`

	//+kubebuilder:validation:Required

	Name string `json:"name"`
}

type FromToWBFunction struct {

	//+kubebuilder:validation:Required

	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`

	//+kubebuilder:validation:Required

	Port int32 `json:"port"`
}

const (
	WBDeployStatusNone        WBDeployStatus = ""
	WBDeployStatusFailed      WBDeployStatus = "Failed"
	WBDeployStatusAllocated   WBDeployStatus = "Allocated"
	WBDeployStatusDeployed    WBDeployStatus = "Deployed"
	WBDeployStatusWaiting     WBDeployStatus = "Waiting"
	WBDeployStatusReleased    WBDeployStatus = "Released"
	WBDeployStatusTerminating WBDeployStatus = "Terminating"
)

const (
	WBIOTypeIncoming WBIOType = "Incoming"
	WBIOTypeOutgoing WBIOType = "Outgoing"
)

const (
	WBTargetIP     = "TargetIP"
	WBTargetPort   = "TargetPort"
	WBProtocol     = "Protocol"
	WBFlowID       = "FlowID"
	WBFunctionID   = "FunctionID"
	WBIOID         = "IOID"
	WBConnectionID = "ConnectionID"
	WBPhysicalPort = "PhysPort"
)

const (
	WBEndOfDataFlowName = "wb-end-of-data-flow"
)

type WBIOUsedType string

const (
	WBIOUsedTypeNone                WBIOUsedType = ""
	WBIOUsedTypeIncoming            WBIOUsedType = "Incoming"
	WBIOUsedTypeOutgoing            WBIOUsedType = "Outgoing"
	WBIOUsedTypeIncomingAndOutgoing WBIOUsedType = "IncomingAndOutgoing"
)

type WBRegionStatus string

const (
	WBRegionStatusNotReady  WBRegionStatus = "NotReady"
	WBRegionStatusPreparing WBRegionStatus = "Preparing"
	WBRegionStatusReady     WBRegionStatus = "Ready"
	WBRegionStatusError     WBRegionStatus = "Error"
)
