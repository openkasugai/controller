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

type FPGAstatus string

const (
	FPGAStatusNotReady  FPGAstatus = "NotReady"
	FPGAStatusPreparing FPGAstatus = "Preparing"
	FPGAStatusReady     FPGAstatus = "Ready"
	FPGAStatusError     FPGAstatus = "Error"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FPGASpec defines the desired state of FPGA
type FPGASpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ChildBitstreamID  *string `json:"childBitstreamID,omitempty"`
	DeviceIndex       int32   `json:"deviceIndex"`
	DeviceFilePath    string  `json:"deviceFilePath"`
	DeviceUUID        string  `json:"deviceUUID"`
	NodeName          string  `json:"nodeName"`
	ParentBitstreamID string  `json:"parentBitstreamID"`
	PCIDomain         int32   `json:"pciDomain"`
	PCIBus            int32   `json:"pciBus"`
	PCIDevice         int32   `json:"pciDevice"`
	PCIFunction       int32   `json:"pciFunction"`
	Vendor            string  `json:"vendor"`
}

// FPGAStatus defines the observed state of FPGA
type FPGAStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ChildBitstreamID     *string    `json:"childBitstreamID,omitempty"`
	ChildBitstreamCRName *string    `json:"childBitstreamCRName,omitempty"`
	DeviceFilePath       string     `json:"deviceFilePath"`
	DeviceIndex          int32      `json:"deviceIndex"`
	DeviceUUID           string     `json:"deviceUUID"`
	NodeName             string     `json:"nodeName"`
	ParentBitstreamID    string     `json:"parentBitstreamID"`
	PCIDomain            int32      `json:"pciDomain"`
	PCIBus               int32      `json:"pciBus"`
	PCIDevice            int32      `json:"pciDevice"`
	PCIFunction          int32      `json:"pciFunction"`
	Status               FPGAstatus `json:"status"`
	Vendor               string     `json:"vendor"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status

// FPGA is the Schema for the fpgas API
type FPGA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FPGASpec   `json:"spec,omitempty"`
	Status FPGAStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FPGAList contains a list of FPGA
type FPGAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FPGA `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FPGA{}, &FPGAList{})
}
