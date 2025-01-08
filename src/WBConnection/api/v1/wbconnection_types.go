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

// WBConnectionSpec defines the desired state of WBConnection
type WBConnectionSpec struct {

	//+kubebuilder:validation:Required

	DataFlowRef WBNamespacedName `json:"dataFlowRef"`

	//+kubebuilder:validation:Required

	ConnectionMethod string `json:"connectionMethod"`

	//+kubebuilder:validation:Optional
	ConnectionPath []WBConnectionPath `json:"connectionPath"`

	//+kubebuilder:validation:Required

	From FromToWBFunction `json:"from"`

	//+kubebuilder:validation:Required

	To FromToWBFunction `json:"to"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Params map[string]intstr.IntOrString `json:"params"`

	//+kubebuilder:validation:Optional
	Requirements *WBConnectionRequirementsInfo `json:"requirements"`
}

// WBConnectionStatus defines the observed state of WBConnection
type WBConnectionStatus struct {

	//+kubebuilder:validation:Required

	DataFlowRef WBNamespacedName `json:"dataFlowRef"`

	//+kubebuilder:validation:Required

	Status WBDeployStatus `json:"status"`

	//+kubebuilder:validation:Required

	ConnectionMethod string `json:"connectionMethod"`

	//+kubebuilder:validation:Optional
	ConnectionPath []WBConnectionPath `json:"connectionPath"`

	//+kubebuilder:validation:Required

	From FromToWBFunction `json:"from"`

	//+kubebuilder:validation:Required

	To FromToWBFunction `json:"to"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Params map[string]intstr.IntOrString `json:"params"`

	//+kubebuilder:validation:Optional
	SatisfiedRequirements *WBConnectionRequirementsInfo `json:"satisfiedRequirements"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	IOs map[string]WBConnectionIO `json:"ios"`

	//+kubebuilder:validation:Optional

	Connections []WBConnectionEdge `json:"connections,omitempty"`
}

type WBConnectionRequirementsInfo struct {
	//+kubebuilder:validation:Required
	Capacity int32 `json:"capacity"`
}

type WBConnectionPath struct {

	//+kubebuilder:validation:Required
	EntityID string `json:"entityID"`

	//+kubebuilder:validation:Required
	UsedType WBIOUsedType `json:"usedType"`
}

type WBConnectionIO struct {

	//+kubebuilder:validation:Required

	Status WBDeployStatus `json:"status"`

	//+kubebuilder:validation:Required

	IoType WBIOType `json:"ioType"`

	//+kubebuilder:validation:Required

	Node string `json:"node"`

	//+kubebuilder:validation:Required

	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required

	DeviceIndex int `json:"deviceIndex"`

	//+kubebuilder:validation:Required

	IoName string `json:"ioName"`

	//+kubebuilder:validation:Required

	Port int `json:"port"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	IntParams map[string]int `json:"intParams"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	StrParams map[string]string `json:"strParams"`
}

type WBConnectionEdge struct {

	//+kubebuilder:validation:Required

	Status WBDeployStatus `json:"status"`

	//+kubebuilder:validation:Required

	From WBNamespacedName `json:"from"`

	//+kubebuilder:validation:Required

	To WBNamespacedName `json:"to"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	IntParams map[string]int `json:"intParams"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	StrParams map[string]string `json:"strParams"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ConnectionMethod",type="string",JSONPath=".status.connectionMethod"
//+kubebuilder:printcolumn:name="From",type="string",JSONPath=".status.from.wbFunctionRef.name"
//+kubebuilder:printcolumn:name="To",type="string",JSONPath=".status.to.wbFunctionRef.name"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// WBConnection is the Schema for the wbconnections API
type WBConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WBConnectionSpec   `json:"spec,omitempty"`
	Status WBConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WBConnectionList contains a list of WBConnection
type WBConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WBConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WBConnection{}, &WBConnectionList{})
}
