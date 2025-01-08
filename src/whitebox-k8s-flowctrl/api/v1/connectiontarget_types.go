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

// ConnectionTargetSpec defines the desired state of ConnectionTarget
type ConnectionTargetSpec struct {

	//+kubebuilder:validation:Required
	IOResourceRef WBNamespacedName `json:"ioResourceRef"`
}

// ConnectionTargetStatus defines the desired state of ConnectionTarget
type ConnectionTargetStatus struct {

	//+kubebuilder:validation:Optional
	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Optional
	DeviceType string `json:"DeviceType"`

	//+kubebuilder:validation:Optional
	DeviceIndex int `json:"deviceIndex"`

	//+kubebuilder:validation:Optional
	RegionName string `json:"region"`

	//+kubebuilder:validation:Optional
	InterfaceName string `json:"interfaceName"`

	//+kubebuilder:validation:Optional
	InterfaceType string `json:"interfaceType"`

	//+kubebuilder:validation:Optional
	InterfaceSideType string `json:"interfaceSideType"`

	//+kubebuilder:validation:Optional
	Available bool `json:"available"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="NodeName",type="string",JSONPath=".status.nodeName"
//+kubebuilder:printcolumn:name="DeviceType",type="string",JSONPath=".status.DeviceType"
//+kubebuilder:printcolumn:name="DeviceIndex",type="integer",JSONPath=".status.deviceIndex"
//+kubebuilder:printcolumn:name="RegionName",type="integer",JSONPath=".status.regionName"
//+kubebuilder:printcolumn:name="InterfaceName",type="string",JSONPath=".status.interfaceName"
//+kubebuilder:printcolumn:name="InterfaceType",type="string",JSONPath=".status.interfaceType"
//+kubebuilder:printcolumn:name="InterfaceSideType",type="string",JSONPath=".status.interfaceSideType"
//+kubebuilder:printcolumn:name="Available",type="boolean",JSONPath=".status.available"

// ConnectionTarget is the Schema for the connectiontargets API
type ConnectionTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectionTargetSpec   `json:"spec,omitempty"`
	Status ConnectionTargetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConnectionTargetList contains a list of ConnectionTarget
type ConnectionTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectionTarget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConnectionTarget{}, &ConnectionTargetList{})
}
