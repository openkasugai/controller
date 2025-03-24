/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED

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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GPUFunctionSpec defines the desired state of GPUFunction
type GPUFunctionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     *int32                        `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      *SharedMemorySpec             `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        string                        `json:"configName"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
}

// GPUFunctionStatus defines the observed state of GPUFunction
type GPUFunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DataFlowRef                    WBNamespacedName  `json:"dataFlowRef"`
	FunctionName                   string            `json:"functionName"`
	ImageURI                       string            `json:"imageURI"`
	SharedMemory                   *SharedMemorySpec `json:"sharedMemory,omitempty"`
	RxProtocol                     *string           `json:"rxProtocol,omitempty"`
	TxProtocol                     *string           `json:"txProtocol,omitempty"`
	ConfigName                     string            `json:"configName"`
	VirtualNetworkDeviceDriverType string            `json:"virtualNetworkDeviceDriverType,omitempty"`
	AdditionalNetwork              *bool             `json:"additionalNetwork,omitempty"`
	FunctionIndex                  *int32            `json:"functionIndex,omitempty"`
	PodName                        *string           `json:"podName,omitempty"`
	StartTime                      metav1.Time       `json:"startTime"`
	//+kubebuilder:default=Pending
	Status              string                              `json:"status"`
	IPAddress           *string                             `json:"Ip,omitempty"`
	AcceleratorStatuses []GPUFunctionAccStatusesByContainer `json:"acceleratorStatuses,omitempty"`
}

type GPUFunctionAccStatusesByContainer struct {
	PartitionName *string                  `json:"partitionName,omitempty"`
	Statuses      []GPUFunctionAccStatuses `json:"statuses,omitempty"`
}

type GPUFunctionAccStatuses struct {
	AcceleratorID *string `json:"acceleratorID,omitempty"`
	Status        *string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// GPUFunction is the Schema for the gpufunctions API
type GPUFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GPUFunctionSpec   `json:"spec,omitempty"`
	Status GPUFunctionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GPUFunctionList contains a list of GPUFunction
type GPUFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GPUFunction `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GPUFunction{}, &GPUFunctionList{})
}
