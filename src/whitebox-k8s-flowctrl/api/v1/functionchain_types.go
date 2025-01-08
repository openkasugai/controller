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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FunctionChainSpec defines the desired state of FunctionChain
type FunctionChainSpec struct {

	//+kubebuilder:validation:Required

	FunctionTypeNamespace string `json:"functionTypeNamespace"`

	//+kubebuilder:validation:Required

	ConnectionTypeNamespace string `json:"connectionTypeNamespace"`

	//+kubebuilder:validation:Required

	Functions map[string]FunctionStruct `json:"functions"`

	//+kubebuilder:validation:Required

	Connections []ConnectionStruct `json:"connections"`
}

type FunctionStruct struct {

	//+kubebuilder:validation:Required

	FunctionName string `json:"functionName"`

	//+kubebuilder:validation:Required

	Version string `json:"version"`

	//+kubebuilder:validation:Optional

	CustomParameter map[string]intstr.IntOrString `json:"customParameter,omitempty"`
}

type ConnectionStruct struct {

	//+kubebuilder:validation:Required

	From FromToFunction `json:"from"`

	//+kubebuilder:validation:Required

	To FromToFunction `json:"to"`

	//+kubebuilder:validation:Required
	// +kubebuilder:default:="auto"

	ConnectionTypeName string `json:"connectionTypeName"`

	//+kubebuilder:validation:Optional

	CustomParameter map[string]intstr.IntOrString `json:"customParameter,omitempty"`
}

type FromToFunction struct {

	//+kubebuilder:validation:Required

	FunctionKey string `json:"functionKey"`

	//+kubebuilder:validation:Required

	Port int32 `json:"port"`
}

type FromToFunctionInfo struct {

	//+kubebuilder:validation:Required

	FunctionKey string `json:"functionKey"`
}

type FromToFunctionScheduleInfo struct {

	//+kubebuilder:validation:Required

	FunctionKey string `json:"functionKey"`

	//+kubebuilder:validation:Optional

	Port *int32 `json:"port"`

	//+kubebuilder:validation:Optional

	InterfaceType *string `json:"interfaceType"`
}

// FunctionChainStatus defines the observed state of FunctionChain
type FunctionChainStatus struct {
	//+kubebuilder:validation:Required
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"

// FunctionChain is the Schema for the functionchains API
type FunctionChain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionChainSpec   `json:"spec,omitempty"`
	Status FunctionChainStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FunctionChainList contains a list of FunctionChain
type FunctionChainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FunctionChain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FunctionChain{}, &FunctionChainList{})
}
