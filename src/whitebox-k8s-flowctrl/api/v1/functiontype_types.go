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

// FunctionTypeSpec defines the desired state of FunctionType
type FunctionTypeSpec struct {

	//+kubebuilder:validation:Required
	FunctionName string `json:"functionName,omitempty"`

	//+kubebuilder:validation:Required
	FunctionInfoCMRef WBNamespacedName `json:"functionInfoCMRef,omitempty"`

	//+kubebuilder:validation:Required
	Version string `json:"version,omitempty"`
}

// FunctionTypeStatus defines the observed state of FunctionType
type FunctionTypeStatus struct {
	//+kubebuilder:validation:Required
	Status string `json:"status"`

	//+kubebuilder:validation:Optional
	RegionTypeCandidates []string `json:"regionTypeCandidates,omitempty"`
	//+kubebuilder:validation:Optional
	RecommendConnection []string `json:"recommendConnection,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="FUNCTION NAME",type="string",JSONPath=".spec.functionName"
//+kubebuilder:printcolumn:name="VERSION",type="string",JSONPath=".spec.version"
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="REGION",type="string",JSONPath=".status.regionTypeCandidates"

// FunctionType is the Schema for the functiontypes API
type FunctionType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionTypeSpec   `json:"spec,omitempty"`
	Status FunctionTypeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FunctionTypeList contains a list of FunctionType
type FunctionTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FunctionType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FunctionType{}, &FunctionTypeList{})
}
