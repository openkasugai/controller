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

package controllertestethernet

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EthernetConnectionSpec defines the desired state of EthernetConnection
type EthernetConnectionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef WBNamespacedName     `json:"dataFlowRef"`
	From        EthernetFunctionSpec `json:"from"`
	To          EthernetFunctionSpec `json:"to"`
}

// EthernetConnectionStatus defines the observed state of EthernetConnection
type EthernetConnectionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef WBNamespacedName       `json:"dataFlowRef"`
	From        EthernetFunctionStatus `json:"from"`
	To          EthernetFunctionStatus `json:"to"`
	StartTime   metav1.Time            `json:"startTime"`
	//+kubebuilder:default=Pending
	Status string `json:"status"`
}

type EthernetFunctionSpec struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
}

type EthernetFunctionStatus struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
	//+kubebuilder:default=INIT
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="FROMFUNC_STATUS",type="string",JSONPath=".status.from.status"
//+kubebuilder:printcolumn:name="TOFUNC_STATUS",type="string",JSONPath=".status.to.status"
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// EthernetConnection is the Schema for the ethernetconnections API
type EthernetConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EthernetConnectionSpec   `json:"spec,omitempty"`
	Status EthernetConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EthernetConnectionList contains a list of EthernetConnection
type EthernetConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EthernetConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EthernetConnection{}, &EthernetConnectionList{})
}
