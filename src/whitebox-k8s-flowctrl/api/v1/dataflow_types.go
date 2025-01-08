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

// DataFlowSpec defines the desired state of DataFlow
type DataFlowSpec struct {

	//+kubebuilder:validation:Required

	FunctionChainRef WBNamespacedName `json:"functionChainRef"`

	//+kubebuilder:validation:Optional

	DryRun *bool `json:"dryrun"`

	//+kubebuilder:validation:Optional

	StartPoint *StartEndPoint `json:"startPoint"`

	//+kubebuilder:validation:Optional

	EndPoint *StartEndPoint `json:"endPoint"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	FunctionUserParameter []FunctionParamStruct `json:"functionUserParameter"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	ConnectionUserParameter []ConnectionParamStruct `json:"connectionUserParameter"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Requirements *DataFlowRequirementsStruct `json:"requirements"`

	//+kubebuilder:validation:Optional
	UserRequirement *string `json:"userRequirement"`
}

type StartEndPoint struct {

	//+kubebuilder:validation:Required

	IP string `json:"ip"`

	//+kubebuilder:validation:Required

	Port int32 `json:"port"`

	//+kubebuilder:validation:Required

	Protocol corev1.Protocol `json:"protocol"`
}

type FunctionParamStruct struct {

	//+kubebuilder:validation:Required

	FunctionKey string `json:"functionKey"`

	//+kubebuilder:validation:Required

	UserParams map[string]intstr.IntOrString `json:"userParams"`
}

type ConnectionParamStruct struct {

	//+kubebuilder:validation:Required

	From FromToFunctionInfo `json:"from"`

	//+kubebuilder:validation:Required

	To FromToFunctionInfo `json:"to"`

	//+kubebuilder:validation:Required

	UserParams map[string]intstr.IntOrString `json:"userParams"`
}

type DataFlowRequirementsStruct struct {

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	All *AllRequirementsInfo `json:"all"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Functions []FunctionRequirementsInfo `json:"functions"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	Connections []ConnectionRequirementsInfo `json:"connections"`
}

type AllRequirementsInfo struct {

	//+kubebuilder:validation:Required

	Capacity int32 `json:"capacity"`
}

type FunctionRequirementsInfo struct {

	//+kubebuilder:validation:Required

	FunctionKey string `json:"functionKey"`

	//+kubebuilder:validation:Required

	Capacity int32 `json:"capacity"`
}

type ConnectionRequirementsInfo struct {

	//+kubebuilder:validation:Required

	From FromToFunctionInfo `json:"from"`

	//+kubebuilder:validation:Required

	To FromToFunctionInfo `json:"to"`

	//+kubebuilder:validation:Required

	Capacity int32 `json:"capacity"`
}

// DataFlowStatus defines the observed state of DataFlow
type DataFlowStatus struct {

	//+kubebuilder:validation:Required
	Status string `json:"status"`

	//+kubebuilder:validation:Optional
	FunctionChain *FunctionChain `json:"functionChain"`

	//+kubebuilder:validation:Optional
	FunctionType []*FunctionType `json:"functionType"`

	//+kubebuilder:validation:Optional
	ConnectionType []*ConnectionType `json:"connectionType"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	ScheduledFunctions map[string]FunctionScheduleInfo `json:"scheduledFunctions"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:pruning:PreserveUnknownFields

	ScheduledConnections []ConnectionScheduleInfo `json:"scheduledConnections"`

	//+kubebuilder:validation:Optional
	StartPoint string `json:"startPoint"`

	//+kubebuilder:validation:Optional
	EndPoint string `json:"endPoint"`
}

type FunctionScheduleInfo struct {

	//+kubebuilder:validation:Required

	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Required

	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required

	DeviceIndex int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Required

	RegionName string `json:"regionName"`

	//+kubebuilder:validation:Optional

	FunctionIndex *int32 `json:"functionIndex"`
}

type ConnectionScheduleInfo struct {

	//+kubebuilder:validation:Required

	From FromToFunctionScheduleInfo `json:"from"`

	//+kubebuilder:validation:Required

	To FromToFunctionScheduleInfo `json:"to"`

	//+kubebuilder:validation:Required

	ConnectionMethod string `json:"connectionMethod"`

	//+kubebuilder:validation:Optional
	ConnectionPath []WBConnectionPath `json:"connectionPath"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="FunctionChain",type="string",JSONPath=".spec.functionChainRef.name"

// DataFlow is the Schema for the dataflows API
type DataFlow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataFlowSpec   `json:"spec,omitempty"`
	Status DataFlowStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataFlowList contains a list of DataFlow
type DataFlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataFlow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataFlow{}, &DataFlowList{})
}
