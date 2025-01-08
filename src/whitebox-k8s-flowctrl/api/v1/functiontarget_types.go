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

// FunctionTargetSpec defines the desired state of FunctionTarget
type FunctionTargetSpec struct {

	//+kubebuilder:validation:Required
	ComputeResourceRef WBNamespacedName `json:"computeResourceRef"`
}

// FunctionTargetStatus defines the desired state of FunctionTarget
type FunctionTargetStatus struct {
	//+kubebuilder:validation:Required
	Status WBRegionStatus `json:"status"`

	//+kubebuilder:validation:Required
	RegionName string `json:"regionName"`

	//+kubebuilder:validation:Required
	RegionType string `json:"regionType"`

	//+kubebuilder:validation:Required
	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Required
	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required
	DeviceIndex int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Required
	Available bool `json:"available"`

	//+kubebuilder:validation:Optional
	MaxFunctions *int32 `json:"maxFunctions"`

	//+kubebuilder:validation:Optional
	CurrentFunctions *int32 `json:"currentFunctions"`

	//+kubebuilder:validation:Optional
	MaxCapacity *int32 `json:"maxCapacity"`

	//+kubebuilder:validation:Optional
	CurrentCapacity *int32 `json:"currentCapacity"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields
	Functions []FunctionCapStruct `json:"functions"`
}

type FunctionCapStruct struct {

	//+kubebuilder:validation:Required
	FunctionIndex int32 `json:"functionIndex"`

	//+kubebuilder:validation:Required
	FunctionName string `json:"functionName"`

	//+kubebuilder:validation:Required
	Available bool `json:"available"`

	//+kubebuilder:validation:Optional
	MaxDataFlows *int32 `json:"maxDataFlows"`

	//+kubebuilder:validation:Optional
	CurrentDataFlows *int32 `json:"currentDataFlows"`

	//+kubebuilder:validation:Optional
	MaxCapacity *int32 `json:"maxCapacity"`

	//+kubebuilder:validation:Optional
	CurrentCapacity *int32 `json:"currentCapacity"`

	//+kubebuilder:validation:Optional
	MaxTimeSlicingSeconds *int32 `json:"maxTimeSlicingSeconds"`

	//+kubebuilder:validation:Optional
	CurrentTimeSlicingSeconds *int32 `json:"currentTimeSlicingSeconds"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="NodeName",type="string",JSONPath=".status.nodeName"
//+kubebuilder:printcolumn:name="DeviceType",type="string",JSONPath=".status.deviceType"
//+kubebuilder:printcolumn:name="DeviceIndex",type="integer",JSONPath=".status.deviceIndex"
//+kubebuilder:printcolumn:name="RegionType",type="string",JSONPath=".status.regionType"
//+kubebuilder:printcolumn:name="Available",type="boolean",JSONPath=".status.available"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="MaxFunctions",type="string",JSONPath=".status.maxFunctions"
//+kubebuilder:printcolumn:name="CurrentFunctions",type="string",JSONPath=".status.currentFunctions"
//+kubebuilder:printcolumn:name="MaxCapacity",type="string",JSONPath=".status.maxCapacity"
//+kubebuilder:printcolumn:name="CurrentCapacity",type="string",JSONPath=".status.currentCapacity"

// FunctionTarget is the Schema for the functiontargets API
type FunctionTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionTargetSpec   `json:"spec,omitempty"`
	Status FunctionTargetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FunctionTargetList contains a list of FunctionTarget
type FunctionTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FunctionTarget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FunctionTarget{}, &FunctionTargetList{})
}
