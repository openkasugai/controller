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

type TypeCombinationStruct struct {

	//+kubebuilder:validation:Optional
	DeviceTypes map[string]string `json:"deviceTypes"`

	//+kubebuilder:validation:Optional
	ConnectionTypes []string `json:"connectionTypes"`

	//+kubebuilder:validation:Optional
	Score *int64 `json:"score"`
}

type TargetCombinationStruct struct {

	//+kubebuilder:validation:Optional
	ScheduledFunctions map[string]FunctionScheduleInfo `json:"scheduledFunctions"`

	//+kubebuilder:validation:Optional
	ScheduledConnections []ConnectionScheduleInfo `json:"scheduledConnections"`

	//+kubebuilder:validation:Optional
	Score *int64 ` json:"score"`
}

type ConnectionIfStruct struct {

	//+kubebuilder:validation:Optional
	NodeName *string `json:"nodeName"`

	//+kubebuilder:validation:Optional
	InterfaceList []string `json:"interfaceList"`
}

type SchedulingDataSpec struct {

	//+kubebuilder:validation:Required
	FilterPipeline []string `json:"filterPipeline"`
}

type SchedulingDataStatus struct {

	//+kubebuilder:validation:Required
	Status string `json:"status"`

	//+kubebuilder:validation:Optional
	CurrentFilterIndex *int32 `json:"currentFilterIndex"`

	//+kubebuilder:validation:Optional
	TypeCombinations []TypeCombinationStruct `json:"typeCombinations"`

	//+kubebuilder:validation:Optional
	TargetCombinations []TargetCombinationStruct `json:"targetCombinations"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
// +kubebuilder:printcolumn:name="Index",type="string",JSONPath=".status.currentFilterIndex"
// +kubebuilder:printcolumn:name="UserRequirement",type="string",JSONPath=".spec.userRequirement"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SchedulingData struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SchedulingDataSpec   `json:"spec,omitempty"`
	Status SchedulingDataStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type SchedulingDataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SchedulingData `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SchedulingData{}, &SchedulingDataList{})
}
