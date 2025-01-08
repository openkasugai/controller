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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WBFunctionSpec defines the desired state of WBFunction
type WBFunctionSpec struct {

	//+kubebuilder:validation:Required

	DataFlowRef WBNamespacedName `json:"dataFlowRef"`

	//+kubebuilder:validation:Required

	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	NodeSelector map[string]string `json:"nodeSelector"`

	//+kubebuilder:validation:Optional

	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	//+kubebuilder:validation:Required

	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required

	DeviceIndex int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Required

	RegionName string `json:"regionName"`

	//+kubebuilder:validation:Optional

	FunctionIndex *int32 `json:"functionIndex"`

	//+kubebuilder:validation:Required

	FunctionName string `json:"functionName"`

	//+kubebuilder:validation:Required

	ConfigName string `json:"configName"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	InputInterface map[string]string `json:"inputInterface"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	OutputInterface map[string]string `json:"outputInterface"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Params map[string]intstr.IntOrString `json:"params"`

	//+kubebuilder:validation:Optional

	PreviousWBFunctions map[string]FromToWBFunction `json:"previousWBFunctions"`

	//+kubebuilder:validation:Optional

	NextWBFunctions map[string]FromToWBFunction `json:"nextWBFunctions"`

	//+kubebuilder:validation:Optional

	MaxDataFlows *int32 `json:"maxDataFlows"`

	//+kubebuilder:validation:Optional

	MaxCapacity *int32 `json:"maxCapacity"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Requirements *WBFunctionRequirementsInfo `json:"requirements"`
}

type WBFunctionRequirementsInfo struct {

	//+kubebuilder:validation:Required

	Capacity int32 `json:"capacity"`
}

// WBFunctionStatus defines the observed state of WBFunction
type WBFunctionStatus struct {

	//+kubebuilder:validation:Required

	DataFlowRef WBNamespacedName `json:"dataFlowRef"`

	//+kubebuilder:validation:Required

	Status WBDeployStatus `json:"status"`

	//+kubebuilder:validation:Required

	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Required

	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required

	DeviceIndex int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Required

	RegionName string `json:"regionName"`

	//+kubebuilder:validation:Required

	FunctionIndex int32 `json:"functionIndex"`

	//+kubebuilder:validation:Required

	FunctionName string `json:"functionName"`

	//+kubebuilder:validation:Required

	ConfigName string `json:"configName"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	InputInterface map[string]string `json:"inputInterface"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	OutputInterface map[string]string `json:"outputInterface"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Params map[string]intstr.IntOrString `json:"params"`

	//+kubebuilder:validation:Optional

	PreviousWBFunctions map[string]FromToWBFunction `json:"previousWBFunctions"`

	//+kubebuilder:validation:Optional

	NextWBFunctions map[string]FromToWBFunction `json:"nextWBFunctions"`

	//+kubebuilder:validation:Optional

	MaxDataFlows *int32 `json:"maxDataFlows"`

	//+kubebuilder:validation:Optional

	MaxCapacity *int32 `json:"maxCapacity"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	SatisfiedRequirements WBFunctionRequirementsInfo `json:"satisfiedRequirements"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="FunctionName",type="string",JSONPath=".status.functionName"
//+kubebuilder:printcolumn:name="NodeName",type="string",JSONPath=".status.nodeName"
//+kubebuilder:printcolumn:name="DeviceType",type="string",JSONPath=".status.deviceType"
//+kubebuilder:printcolumn:name="DeviceIndex",type="integer",JSONPath=".status.deviceIndex"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// WBFunction is the Schema for the wbfunctions API
type WBFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WBFunctionSpec   `json:"spec,omitempty"`
	Status WBFunctionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WBFunctionList contains a list of WBFunction
type WBFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WBFunction `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WBFunction{}, &WBFunctionList{})
}
