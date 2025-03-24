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
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ChildBitstreamState string

const (
	ChildBsStoppingModule        ChildBitstreamState = "StoppingModule"
	ChildBsNotStopNetworkModule  ChildBitstreamState = "NotStopNetworkModule"
	ChildBsStoppingNetworkModule ChildBitstreamState = "StoppingNetworkModule"
	ChildBsNotWriteBsfile        ChildBitstreamState = "NotWriteBitstreamFile"
	ChildBsReconfiguring         ChildBitstreamState = "Reconfiguring"
	ChildBsWritingBsfile         ChildBitstreamState = "WritingBitstreamFile"
	ChildBsConfiguringParam      ChildBitstreamState = "ConfiguringParameters"
	ChildBsNoConfigureNetwork    ChildBitstreamState = "NoConfigureNetwork"
	ChildBsConfiguringNetwork    ChildBitstreamState = "ConfiguringNetwork"
	ChildBsReady                 ChildBitstreamState = "Ready"
	ChildBsError                 ChildBitstreamState = "Error"
)

type ChildBitstreamStatus string

const (
	ChildBsStatusNotReady  ChildBitstreamStatus = "NotReady"
	ChildBsStatusPreparing ChildBitstreamStatus = "Preparing"
	ChildBsStatusReady     ChildBitstreamStatus = "Ready"
	ChildBsStatusError     ChildBitstreamStatus = "Error"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type ChildBsDetails struct {
	Port             *int32 `json:"port,omitempty"`
	DMAChannelID     *int32 `json:"dmaChannelID,omitempty"`
	LLDMAConnectorID *int32 `json:"lldmaConnectorID,omitempty"`
}

type RxTxSpec struct {
	Protocol *map[string]ChildBsDetails `json:"protocol,omitempty"`
}

type FunctionsDeploySpec struct {
	MaxCapacity  *int32 `json:"maxCapacity,omitempty"`
	MaxDataFlows *int32 `json:"maxDataFlows,omitempty"`
}

type FunctionsIntraResourceMgmtMap struct {
	Available      *bool     `json:"available,omitempty"`
	FunctionCRName *string   `json:"functionCRName,omitempty"`
	Rx             *RxTxSpec `json:"rx,omitempty"`
	Tx             *RxTxSpec `json:"tx,omitempty"`
}

type FunctionsModule struct {
	FunctionChannelIDs *string `json:"function-channel-ids,omitempty"`
	Identifier         *string `json:"identifier,omitempty"`
	Type               *string `json:"type,omitempty"`
	Version            *string `json:"version,omitempty"`
}

type ChildBsFunctions struct {
	ID                   *int32                                    `json:"id,omitempty"`
	FunctionName         *string                                   `json:"functionname,omitempty"`
	Module               *[]FunctionsModule                        `json:"module,omitempty"`
	Parameters           *map[string]intstr.IntOrString            `json:"parameters,omitempty"`
	IntraResourceMgmtMap *map[string]FunctionsIntraResourceMgmtMap `json:"intraResourceMgmtMap,omitempty"`
	DeploySpec           FunctionsDeploySpec                       `json:"deploySpec"`
}

type ConversionModule struct {
	Identifier *string `json:"identifier,omitempty"`
	Type       *string `json:"type,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type ChildBsConversion struct {
	ID     *int32              `json:"id,omitempty"`
	Module *[]ConversionModule `json:"module,omitempty"`
}

type ChildBsDirecttrans struct {
	ID         *int32  `json:"id,omitempty"`
	Identifier *string `json:"identifier,omitempty"`
	Type       *string `json:"type,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type ChildBsChain struct {
	ID         *int32  `json:"id,omitempty"`
	Identifier *string `json:"identifier,omitempty"`
	Type       *string `json:"type,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type ChildBsLLDMA struct {
	Cids *string `json:"cids,omitempty"`
	ID   *int32  `json:"id,omitempty"`
	// ExtIfID *int32  `json:"extIfId,omitempty"`
}

type ChildBsPtu struct {
	Cids       *string                        `json:"cids,omitempty"`
	ID         *int32                         `json:"id,omitempty"`
	Parameters *map[string]intstr.IntOrString `json:"parameters,omitempty"`
	// ExtIfID    *int32                         `json:"extIfId,omitempty"`
}

type ChildBsModule struct {
	Ptu         *ChildBsPtu         `json:"ptu,omitempty"`
	LLDMA       *ChildBsLLDMA       `json:"lldma,omitempty"`
	Chain       *ChildBsChain       `json:"chain,omitempty"`
	Directtrans *ChildBsDirecttrans `json:"directtrans,omitempty"`
	Conversion  *ChildBsConversion  `json:"conversion,omitempty"`
	Functions   *[]ChildBsFunctions `json:"functions,omitempty"`
}

type ChildBsRegion struct {
	Modules      *ChildBsModule `json:"modules,omitempty"`
	MaxFunctions *int32         `json:"maxFunctions,omitempty"`
	MaxCapacity  *int32         `json:"maxCapacity,omitempty"`
	Name         *string        `json:"name,omitempty"`
}

type BsConfigInfo struct {
	ChildBitstreamIDs []ChildBsSpec `json:"child-bitstream-ids"`
	ParentBitstreamID string        `json:"parent-bitstream-id"`
}

// ChildBitstreamSpec defines the desired state of ChildBitstream
type ChildBsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Regions            []ChildBsRegion `json:"regions"`
	ChildBitstreamID   *string         `json:"child-bitstream-id,omitempty"`
	ChildBitstreamFile *string         `json:"child-bitstream-file,omitempty"`
}

// ChildBitstreamStatus defines the observed state of ChildBitstream
type ChildBsStatus struct {
	Regions []ChildBsRegion `json:"regions"`
	//+kubebuilder:default=NotReady
	Status             ChildBitstreamStatus `json:"status"`
	State              ChildBitstreamState  `json:"state"`
	ChildBitstreamID   *string              `json:"child-bitstream-id,omitempty"`
	ChildBitstreamFile *string              `json:"child-bitstream-file,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status

// ChildBs is the Schema for the childbs API
type ChildBs struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChildBsSpec   `json:"spec,omitempty"`
	Status ChildBsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChildBsList contains a list of ChildBs
type ChildBsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChildBs `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChildBs{}, &ChildBsList{})
}
