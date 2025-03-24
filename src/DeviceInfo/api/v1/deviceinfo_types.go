/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// constructor
const (
	RequestDeploy      = "Deploy"     // DeviceInfo.Spec.Request.RequestType(Deploy)
	RequestUpdate      = "Update"     // DeviceInfo.Spec.Request.RequestType(Update)
	RequestUndeploy    = "Undeploy"   // DeviceInfo.Spec.Request.RequestType(Undeploy)
	ResponceDeployed   = "Deployed"   // DeviceInfo.Status.Responce.Status(Deployed)
	ResponceUndeployed = "Undeployed" // DeviceInfo.Status.Responce.Status(Undeployed)
	ResponceInitial    = "Initial"    // DeviceInfo.Status.Responce.Status(Initial)
	ResponceError      = "Error"      // DeviceInfo.Status.Responce.Status(Error)
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeviceInfoSpec defines the desired state of DeviceInfo
type DeviceInfoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Request WBFuncRequest `json:"request"`
}

// DeviceInfoStatus defines the observed state of DeviceInfo
type DeviceInfoStatus struct {
	Response WBFuncResponse `json:"response,omitempty"`
}

type WBFuncRequest struct {
	RequestType   string `json:"requestType"`
	DeviceType    string `json:"deviceType"`
	DeviceIndex   int32  `json:"deviceIndex"`
	RegionName    string `json:"regionName"`
	NodeName      string `json:"nodeName"`
	FunctionIndex *int32 `json:"functionIndex,omitempty"`
	FunctionName  string `json:"functionName"`
	MaxDataFlows  *int32 `json:"maxDataFlows,omitempty"`
	MaxCapacity   *int32 `json:"maxCapacity,omitempty"`
	Capacity      *int32 `json:"capacity,omitempty"`
}

type WBFuncResponse struct {
	//+kubebuilder:default=Initial
	Status         string `json:"status"`
	FunctionIndex  *int32 `json:"functionIndex,omitempty"`
	DeviceUUID     string `json:"deviceUUID,omitempty"`
	DeviceFilePath string `json:"deviceFilePath,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:print32column:name="Status",type="string",JSONPath=".status.Status"
//+kubebuilder:print32column:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DeviceInfo is the Schema for the deviceinfos API
type DeviceInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceInfoSpec   `json:"spec,omitempty"`
	Status DeviceInfoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeviceInfoList contains a list of DeviceInfo
type DeviceInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeviceInfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeviceInfo{}, &DeviceInfoList{})
}
