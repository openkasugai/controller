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

const (
	FPGARECONFSTATUSSUCCEEDED = "Succeeded"
	FPGARECONFSTATUSFAILED    = "Failed"
)

// FPGAReconfigurationSpec defines the desired state of FPGAFunction
type FPGAReconfigurationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NodeName         string            `json:"nodeName"`
	DeviceFilePath   string            `json:"deviceFilePath"`
	FPGAResetFlag    *bool             `json:"fpgaResetFlag,omitempty"`
	ChildBsResetFlag *bool             `json:"childBsResetFlag,omitempty"`
	ConfigNames      []FPGAConfigNames `json:"configNames,omitempty"`
}

type FPGAConfigNames struct {
	LaneIndex  int32  `json:"laneIndex"`
	ConfigName string `json:"configName"`
}

// FPGAReconfigrationStatus defines the observed state of FPGAFunction
type FPGAReconfigurationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+kubebuilder:default=Pending
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FPGAReconfigration is the Schema for the fpgafunctions API
type FPGAReconfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FPGAReconfigurationSpec   `json:"spec,omitempty"`
	Status FPGAReconfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FPGAReconfigrationList contains a list of FPGAFunction
type FPGAReconfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FPGAReconfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FPGAReconfiguration{}, &FPGAReconfigurationList{})
}
