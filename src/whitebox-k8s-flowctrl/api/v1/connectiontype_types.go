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

// ConnectionTypeSpec defines the desired state of ConnectionType
type ConnectionTypeSpec struct {
	//+kubebuilder:validation:Required
	ConnectionInfoNameSpaces []string `json:"connectionInfoNamespaces"`

	//+kubebuilder:validation:Optional
	ConnectionTypeName string `json:"connectionTypeName"`
}

// ConnectionTypeStatus defines the observed state of ConnectionType
type ConnectionTypeStatus struct {
	//+kubebuilder:validation:Required
	Status string `json:"status"`

	//+kubebuilder:validation:Optional
	AvailableInterfaces map[string]AvailableInterfaceStruct `json:"availableInterfaces"`

	//+kubebuilder:validation:Optional
	Interfaces []string `json:"interfaces"`
}

type AvailableInterfaceStruct struct {
	//+kubebuilder:validation:Optional
	Destinations map[string]DestinationStruct `json:"destinations"`
}

type DestinationStruct struct {
	//+kubebuilder:validation:Optional
	Route []RouteStruct `json:"route"`
}

type RouteStruct struct {
	//+kubebuilder:validation:Optional
	Type string `json:"type"`
	//+kubebuilder:validation:Optional
	From string `json:"from"`
	//+kubebuilder:validation:Optional
	To string `json:"to"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="INTERFACES",type="string",JSONPath=".status.interfaces"

// ConnectionType is the Schema for the connectiontypes API
type ConnectionType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectionTypeSpec   `json:"spec,omitempty"`
	Status ConnectionTypeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConnectionTypeList contains a list of ConnectionType
type ConnectionTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectionType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConnectionType{}, &ConnectionTypeList{})
}
