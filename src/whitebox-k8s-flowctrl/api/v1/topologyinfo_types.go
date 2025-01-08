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

// TopologyInfoSpec defines the desired state of TopologyInfo
type TopologyInfoSpec struct {
	//+kubebuilder:validation:Required
	TopologyDataCMRef []WBNamespacedName `json:"topologyDataCMRef"`
}

// TopologyInfoStatus defines the observed state of TopologyInfo
type TopologyInfoStatus struct {
	//+kubebuilder:validation:Required
	Entities []EntityInfo `json:"entities"`

	//+kubebuilder:validation:Required
	Relations []RelationInfo `json:"relations"`
}

type EntityInfo struct {
	//+kubebuilder:validation:Required
	ID string `json:"id"`

	//+kubebuilder:validation:Required
	Type string `json:"type"`

	//+kubebuilder:validation:Optional
	LocationInfo *LocationInfo `json:"locationInfo"`

	//+kubebuilder:validation:Optional
	CapacityInfo *CapacityInfo `json:"capacityInfo"`

	//+kubebuilder:validation:Optional
	NodeInfo *NodeInfo `json:"nodeInfo"`

	//+kubebuilder:validation:Optional
	DeviceInfo *DeviceInfo `json:"deviceInfo"`

	//+kubebuilder:validation:Optional
	InterfaceInfo *InterfaceInfo `json:"interfaceInfo"`

	//+kubebuilder:validation:Optional
	NetworkInfo *NetworkInfo `json:"networkInfo"`

	//+kubebuilder:validation:Required
	Available bool `json:"available"`
}

type LocationInfo struct {
	//+kubebuilder:validation:Optional
	Rack string `json:"rack"`

	//+kubebuilder:validation:Optional
	DataCenter string `json:"dataCenter"`
}

type CapacityInfo struct {
	//+kubebuilder:validation:Optional
	MaxFunctions int32 `json:"maxFunctions"`

	//+kubebuilder:validation:Optional
	CurrentFunctions int32 `json:"currentFunctions"`

	//+kubebuilder:validation:Optional
	MaxIncomingCapacity int32 `json:"maxIncomingCapacity"`

	//+kubebuilder:validation:Optional
	CurrentIncomingCapacity int32 `json:"currentIncomingCapacity"`

	//+kubebuilder:validation:Optional
	MaxOutgoingCapacity int32 `json:"maxOutgoingCapacity"`

	//+kubebuilder:validation:Optional
	CurrentOutgoingCapacity int32 `json:"currentOutgoingCapacity"`
}

type NodeInfo struct {
	//+kubebuilder:validation:Required
	NodeName string `json:"nodeName"`
}

type DeviceInfo struct {
	//+kubebuilder:validation:Required
	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Required
	DeviceType string `json:"deviceType"`

	//+kubebuilder:validation:Required
	DeviceIndex int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Optional
	RegionName *string `json:"regionName"`
}

type InterfaceInfo struct {
	//+kubebuilder:validation:Required
	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Optional
	DeviceType *string `json:"deviceType"`

	//+kubebuilder:validation:Optional
	DeviceIndex *int32 `json:"deviceIndex"`

	//+kubebuilder:validation:Required
	InterfaceType string `json:"interfaceType"`

	//+kubebuilder:validation:Required
	InterfaceIndex int32 `json:"interfaceIndex"`

	//+kubebuilder:validation:Required
	InterfaceSideType string `json:"interfaceSideType"`
}

type NetworkInfo struct {
	//+kubebuilder:validation:Optional
	NodeName *string `json:"nodeName"`

	//+kubebuilder:validation:Required
	NetworkType string `json:"networkType"`

	//+kubebuilder:validation:Required
	NetworkIndex int32 `json:"networkIndex"`

	//+kubebuilder:validation:Required
	NetworkSideType string `json:"networkSideType"`
}

type RelationInfo struct {
	//+kubebuilder:validation:Required
	Type string `json:"type"`

	//+kubebuilder:validation:Required
	From string `json:"from"`

	//+kubebuilder:validation:Required
	To string `json:"to"`

	//+kubebuilder:validation:Required
	Available bool `json:"available"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=topologyinfos

// TopologyInfo is the Schema for the topologyinfos API
type TopologyInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologyInfoSpec   `json:"spec,omitempty"`
	Status TopologyInfoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TopologyInfoList contains a list of TopologyInfo
type TopologyInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopologyInfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopologyInfo{}, &TopologyInfoList{})
}
