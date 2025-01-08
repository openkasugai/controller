/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FPGAFunctionSpec defines the desired state of FPGAFunction
type FPGAFunctionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DataFlowRef       WBNamespacedName            `json:"dataFlowRef"`
	FunctionName      string                      `json:"functionName"`
	NodeName          string                      `json:"nodeName"`
	DeviceType        string                      `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                 `json:"acceleratorIDs"`
	RegionName        string                      `json:"regionName"`
	FunctionIndex     *int32                      `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo                  `json:"envs,omitempty"`
	ConfigName        string                      `json:"configName"`
	SharedMemory      *SharedMemorySpec           `json:"sharedMemory,omitempty"`
	FunctionKernelID  *int32                      `json:"functionKernelID,omitempty"`
	FunctionChannelID *int32                      `json:"functionChannelID,omitempty"`
	PtuKernelID       *int32                      `json:"ptuKernelID,omitempty"`
	FrameworkKernelID *int32                      `json:"frameworkKernelID,omitempty"`
	Rx                RxTxData                    `json:"rx,omitempty"`
	Tx                RxTxData                    `json:"tx,omitempty"`
	PreviousFunctions map[string]FromToWBFunction `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction `json:"nextFunctions,omitempty"`
}

type AccIDInfo struct {
	PartitionName string `json:"partitionName"`
	ID            string `json:"id"`
}

type EnvsInfo struct {
	PartitionName string     `json:"partitionName"`
	EachEnv       []EnvsData `json:"eachEnv"`
}

type EnvsData struct {
	EnvKey   string `json:"envKey"`
	EnvValue string `json:"envValue"`
}

type RxTxData struct {
	Protocol         string  `json:"protocol"`
	IPAddress        *string `json:"ipAddress,omitempty"`
	Port             *int32  `json:"port,omitempty"`
	SubnetAddress    *string `json:"subnetAddress,omitempty"`
	GatewayAddress   *string `json:"gatewayAddress,omitempty"`
	DMAChannelID     *int32  `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32  `json:"lldmaConnectorID,omitempty"`
}

type SharedMemorySpec struct {
	FilePrefix      string `json:"filePrefix"`
	CommandQueueID  string `json:"commandQueueID"`
	SharedMemoryMiB int32  `json:"sharedMemoryMiB"`
}

// FPGAFunctionStatus defines the observed state of FPGAFunction
type FPGAFunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	StartTime metav1.Time `json:"startTime"`
	//+kubebuilder:default=Pending
	Status              string                `json:"status"`
	DataFlowRef         WBNamespacedName      `json:"dataFlowRef"`
	FunctionName        string                `json:"functionName"`
	FunctionIndex       int32                 `json:"functionIndex"`
	ParentBitstreamName string                `json:"parentBitstreamName"`
	ChildBitstreamName  string                `json:"childBitstreamName"`
	SharedMemory        *SharedMemorySpec     `json:"sharedMemory,omitempty"`
	FunctionKernelID    int32                 `json:"functionKernelID"`
	FunctionChannelID   int32                 `json:"functionChannelID"`
	PtuKernelID         int32                 `json:"ptuKernelID"`
	FrameworkKernelID   int32                 `json:"frameworkKernelID"`
	Rx                  RxTxData              `json:"rx"`
	Tx                  RxTxData              `json:"tx"`
	AcceleratorStatuses []AccStatusesByDevice `json:"acceleratorStatuses,omitempty"`
}

type AccStatusesByDevice struct {
	PartitionName *string       `json:"partitionName,omitempty"`
	Statused      []AccStatuses `json:"statuses,omitempty"`
}

type AccStatuses struct {
	AcceleratorID *string `json:"acceleratorID,omitempty"`
	Status        *string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FPGAFunction is the Schema for the fpgafunctions API
type FPGAFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FPGAFunctionSpec   `json:"spec,omitempty"`
	Status FPGAFunctionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FPGAFunctionList contains a list of FPGAFunction
type FPGAFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FPGAFunction `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FPGAFunction{}, &FPGAFunctionList{})
}
