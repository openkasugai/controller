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

// ComputeResourceSpec defines the desired state of ComputeResource
type ComputeResourceSpec struct {

	//+kubebuilder:validation:Optional
	Regions []RegionInfo `json:"regions"`

	//+kubebuilder:validation:Required
	NodeName string `json:"nodeName"`
}

// ComputeResourceStatus defines the desired state of ComputeResource
type ComputeResourceStatus struct {

	//+kubebuilder:validation:Optional
	Regions []RegionInfo `json:"regions"`

	//+kubebuilder:validation:Required
	NodeName string `json:"nodeName"`
}

type RegionInfo struct {
	//+kubebuilder:validation:Required
	Name string `json:"name"`

	//+kubebuilder:validation:Required
	Type string `json:"type"`

	//+kubebuilder:validation:Required
	DeviceFilePath string `json:"deviceFilePath"`

	//+kubebuilder:validation:Optional
	DeviceUUID *string `json:"deviceUUID"`

	//+kubebuilder:validation:Required
	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required
	DeviceIndex int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Required
	Available bool `json:"available"`

	//+kubebuilder:validation:Required
	Status WBRegionStatus `json:"status"`

	//+kubebuilder:validation:Optional
	MaxFunctions *int32 `json:"maxFunctions"`

	//+kubebuilder:validation:Optional
	CurrentFunctions *int32 `json:"currentFunctions"`

	//+kubebuilder:validation:Optional
	MaxCapacity *int32 `json:"maxCapacity"`

	//+kubebuilder:validation:Optional
	CurrentCapacity *int32 `json:"currentCapacity"`

	//+kubebuilder:validation:Optional
	MaxTimeSlicingSeconds *int32 `json:"maxTimeSlicingSeconds"`

	//+kubebuilder:validation:Optional
	CurrentTimeSlicingSeconds *int32 `json:"currentTimeSlicingSeconds"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields
	Functions []FunctionInfrastruct `json:"functions"`
}

type FunctionInfrastruct struct {

	//+kubebuilder:validation:Required
	FunctionIndex int32 `json:"functionIndex"`

	//+kubebuilder:validation:Required
	PartitionName string `json:"partitionName"`

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
// kubebuilder:subresource:status

// ComputeResource is the Schema for the computeresources API
type ComputeResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComputeResourceSpec   `json:"spec,omitempty"`
	Status ComputeResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComputeResourceList contains a list of ComputeResource
type ComputeResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComputeResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComputeResource{}, &ComputeResourceList{})
}
