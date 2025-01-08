package controller

import (
	examplecomv1 "GPUFunction/api/v1"
	// k8scnicncfio "github.com/k8snetworkplumbingwg/network-attachment-definition-client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var partitionName1 string = "df-night01-wbfunction-high-infer-main"

var GPUFunction1 = examplecomv1.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:   "a100",
		FunctionName: "high-infer",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"outputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"outputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "df-night01-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night01-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "a100",
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctiontest",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var PCIeConnection1 = PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: WBNamespacedName{
			Name:      "pciefunctiontest",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "fpgafunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "pciefunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var frameworkKernelID int32 = 0
var functionChannelID int32 = 0
var functionIndex int32 = 0
var functionKernelID int32 = 0
var ptuKernelID int32 = 0

var FPGAFunction1 = FPGAFunction{
	// TypeMeta: metav1.TypeMeta{
	// 	APIVersion: "example.com/v1",
	// 	Kind:       "FPGAFunction",
	// },
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
		Namespace: "default",
	},
	Spec: FPGAFunctionSpec{
		AcceleratorIDs: []AccIDInfo{
			{
				PartitionName: "0",
				ID:            "/dev/xpcie_21320621V00D",
			},
		},
		ConfigName: "fpgafunc-config-filter-resize-high-infer",
		DataFlowRef: WBNamespacedName{
			Name:      "df-night01",
			Namespace: "default",
		},
		DeviceType:        "alveo",
		FrameworkKernelID: &frameworkKernelID,
		FunctionChannelID: &functionChannelID,
		FunctionIndex:     &functionIndex,
		FunctionKernelID:  &functionKernelID,
		FunctionName:      "filter-resize-high-infer",
		NodeName:          "worker1",
		PtuKernelID:       &ptuKernelID,
		RegionName:        "lane0",
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
	},
	Status: FPGAFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: WBNamespacedName{
			Name:      "fpgafunctiontest",
			Namespace: "default",
		},
		FunctionName:        "filter-resize-high-infer",
		ParentBitstreamName: "ver2_tpcie_tandem1.mcs",
		ChildBitstreamName:  "ver1_tpcie_tandem2.bit",
		FrameworkKernelID:   0,
		FunctionChannelID:   0,
		PtuKernelID:         0,
		Rx: RxTxSpec{
			Protocol: "TCP",
		},
		Tx: RxTxSpec{
			Protocol: "DMA",
		},
		Status: "pending",
	},
}

var NetworkAttachmentDefinition1 = NetworkAttachmentDefinition{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "k8s.cni.cncf.io",
		Kind:       "NetworkAttachmentDefinition",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "worker1-config-net-sriov",
		Namespace: "default",
		Annotations: map[string]string{
			"k8s.v1.cni.cncf.io/resourceName": "intel.com/intel_sriov_netdevice",
		},
	},
	Spec: NetworkAttachmentDefinitionSpec{
		Config: `{
			"type": "sriov",
			"cniVersion": "0.3.1",
			"name": "worker1-config-net-sriov",
			"ipam": {
				"type": "static"
				}
			}`,
	},
}

var functionIndexg int32 = 99
var GPUFunction2 = examplecomv1.GPUFunction{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "GPUFunction",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbfunction-high-infer-main",
		Namespace: "default",
	},
	Spec: examplecomv1.GPUFunctionSpec{
		AcceleratorIDs: []examplecomv1.AccIDInfo{
			{
				PartitionName: &partitionName1,
				ID:            "GPU-702fb653-43a4-732d-6bc4-7b3487696c90",
			},
		},
		ConfigName: "gpufunc-config-high-infer",
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "df-night02",
			Namespace: "default",
		},
		DeviceType:    "a100",
		FunctionIndex: &functionIndexg,
		FunctionName:  "high-infer",
		PreviousFunctions: map[string]examplecomv1.FromToWBFunction{
			"0": {
				Port: 0,
				WBFunctionRef: examplecomv1.WBNamespacedName{
					Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
					Namespace: "default",
				},
			},
		},
		NodeName: "worker1",
		Params: map[string]intstr.IntOrString{
			"outputIPAddress": {
				StrVal: "192.168.122.40",
				Type:   1,
			},
			"outputPort": {
				IntVal: 8556,
			},
		},
		SharedMemory: &examplecomv1.SharedMemorySpec{
			FilePrefix:      "df-night02-wbfunction-filter-resize-high-infer-main",
			CommandQueueID:  "df-night02-wbfunction-filter-resize-high-infer-main",
			SharedMemoryMiB: 1,
		},
		RegionName: "a100",
	},
	Status: examplecomv1.GPUFunctionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: examplecomv1.WBNamespacedName{
			Name:      "gpufunctiontest",
			Namespace: "default",
		},
		FunctionName: "high-infer",
		ImageURI:     "container",
		ConfigName:   "configname",
		Status:       "pending",
	},
}

var PCIeConnection2 = PCIeConnection{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "example.com/v1",
		Kind:       "PCIeConnection",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "df-night02-wbconnection-filter-resize-high-infer-main-high-infer-main",
		Namespace: "default",
	},
	Spec: PCIeConnectionSpec{
		DataFlowRef: WBNamespacedName{
			Name:      "df-night02",
			Namespace: "default",
		},
		From: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
				Namespace: "default",
			},
		},
		To: PCIeFunctionSpec{
			WBFunctionRef: WBNamespacedName{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			},
		},
	},
	Status: PCIeConnectionStatus{
		StartTime: metav1.Now(),
		DataFlowRef: WBNamespacedName{
			Name:      "pciefunctiontest",
			Namespace: "default",
		},
		Status: "pending",
		From: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "fpgafunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
		To: PCIeFunctionStatus{
			WBFunctionRef: WBNamespacedName{
				Name:      "pciefunctiontest",
				Namespace: "default",
			},
			Status: "pending",
		},
	},
}

var gpuconfigdecode = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "gpufunc-config-high-infer",
		Namespace: "default",
	},
	Data: map[string]string{
		"gpufunc-config-high-infer.json": ` 
		[{
			"rxProtocol": "DMA",
			"txProtocol":"RTP",
			"sharedMemoryMiB": 256,
			"imageURI": "localhost/gpu-deepstream-app:3.1.0",
			"additionalNetwork": true,
			"virtualNetworkDeviceDriverType": "sriov",
			"envs":{
			  "CUDA_MPS_PIPE_DIRECTORY": "/tmp/nvidia-mps",
			  "CUDA_MPS_LOG_DIRECTORY": "/tmp/nvidia-mps",
			  "SHMEM_SECONDARY": "1",
			  "HEIGHT": "1280",
			  "WIDTH": "1280"
			},
			"template":{
			  "apiVersion": "v1",
			  "kind": "Pod",
			  "spec":{
			    "containers":[{
			      "name": "gfunc-hi-1",
			      "workingDir": "/opt/nvidia/deepstream/deepstream-6.3",
			      "command": ["sh", "-c"],
			      "args":["cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgasrc !",
			         "'video/x-raw,format=(string)BGR,%WIDTH%,%HEIGHT%'",
			         "! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA'",
			         "! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1",
			         "%WIDTH%",
			         "%HEIGHT%",
			         "! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1",
			         "model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert !",
			         "'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink",
			         "%OUTPUTIP%",
			         "%OUTPUTPORT%",
			         "sync=true"],
			      "securityContext":{
			        "privileged": true
			      },
			      "volumeMounts":[{
			        "name": "hugepage-1gi",
			        "mountPath": "/dev/hugepages"
			      },{
			        "name": "host-nvidia-mps",
			        "mountPath": "/tmp/nvidia-mps"
			      },{
			        "name": "dpdk",
			        "mountPath": "/var/run/dpdk"
			      }],
			      "resources":{
			        "requests":{
			          "memory": "32Gi"
			        },
			        "limits":{
			          "hugepages-1Gi": "1Gi"
			        }
			      }
			    }],
			    "volumes":[{
			      "name": "hugepage-1gi",
			      "hostPath":
			       {"path": "/dev/hugepages"}
			    },{
			      "name": "host-nvidia-mps",
			      "hostPath":
			       {"path": "/tmp/nvidia-mps"}
			    },{
			      "name": "dpdk",
			      "hostPath":
			       {"path": "/var/run/dpdk"}
			    }],
			    "hostNetwork": false,
			    "hostIPC": true,
			    "restartPolicy": "Always"
				}
			}
		}]`,
	},
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvsInfo) DeepCopyInto(out *EnvsInfo) {
	*out = *in
	if in.EachEnv != nil {
		in, out := &in.EachEnv, &out.EachEnv
		*out = make([]EnvsData, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvsInfo.
func (in *EnvsInfo) DeepCopy() *EnvsInfo {
	if in == nil {
		return nil
	}
	out := new(EnvsInfo)
	in.DeepCopyInto(out)
	return out
}

var (
	// GroupVersion is group version used to register these objects
	GroupVersionpcie = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilderpcie = &scheme.Builder{GroupVersion: GroupVersionpcie}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemepcie = SchemeBuilderpcie.AddToScheme
)

// PCIeConnection difines the PCIeConnection CR
type PCIeConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PCIeConnectionSpec   `json:"spec,omitempty"`
	Status PCIeConnectionStatus `json:"status,omitempty"`
}

// PCIeConnectionSpec defines the desired state of PCIeConnection
type PCIeConnectionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef WBNamespacedName `json:"dataFlowRef"`
	From        PCIeFunctionSpec `json:"from"`
	To          PCIeFunctionSpec `json:"to"`
}

// PCIeConnectionStatus defines the observed state of PCIeConnection
type PCIeConnectionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DataFlowRef  WBNamespacedName   `json:"dataFlowRef"`
	From         PCIeFunctionStatus `json:"from"`
	To           PCIeFunctionStatus `json:"to"`
	SharedMemory SharedMemoryStatus `json:"sharedMemory,omitempty"`
	StartTime    metav1.Time        `json:"startTime"`
	//+kubebuilder:default=Pending
	Status string `json:"status"`
}

type PCIeFunctionSpec struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
}

type PCIeFunctionStatus struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
	//+kubebuilder:default=INIT
	Status string `json:"status"`
}

type SharedMemoryStatus struct {
	// +optional
	Status string `json:"status,omitempty"`
}

// PCIeConnectionList contains a list of PCIeConnection
type PCIeConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PCIeConnection `json:"items"`
}

type WBNamespacedName struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func init() {
	SchemeBuilderpcie.Register(&PCIeConnection{}, &PCIeConnectionList{})
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnection) DeepCopyInto(out *PCIeConnection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnection.
func (in *PCIeConnection) DeepCopy() *PCIeConnection {
	if in == nil {
		return nil
	}
	out := new(PCIeConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PCIeConnection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnectionList) DeepCopyInto(out *PCIeConnectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PCIeConnection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnectionList.
func (in *PCIeConnectionList) DeepCopy() *PCIeConnectionList {
	if in == nil {
		return nil
	}
	out := new(PCIeConnectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PCIeConnectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnectionSpec) DeepCopyInto(out *PCIeConnectionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	out.From = in.From
	out.To = in.To
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnectionSpec.
func (in *PCIeConnectionSpec) DeepCopy() *PCIeConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(PCIeConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeConnectionStatus) DeepCopyInto(out *PCIeConnectionStatus) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	out.From = in.From
	out.To = in.To
	out.SharedMemory = in.SharedMemory
	in.StartTime.DeepCopyInto(&out.StartTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeConnectionStatus.
func (in *PCIeConnectionStatus) DeepCopy() *PCIeConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(PCIeConnectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeFunctionSpec) DeepCopyInto(out *PCIeFunctionSpec) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeFunctionSpec.
func (in *PCIeFunctionSpec) DeepCopy() *PCIeFunctionSpec {
	if in == nil {
		return nil
	}
	out := new(PCIeFunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PCIeFunctionStatus) DeepCopyInto(out *PCIeFunctionStatus) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PCIeFunctionStatus.
func (in *PCIeFunctionStatus) DeepCopy() *PCIeFunctionStatus {
	if in == nil {
		return nil
	}
	out := new(PCIeFunctionStatus)
	in.DeepCopyInto(out)
	return out
}

// FPGAFunctionSpec defines the desired state of FPGAFunction
type FPGAFunctionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DataFlowRef       WBNamespacedName  `json:"dataFlowRef"`
	FunctionName      string            `json:"functionName"`
	NodeName          string            `json:"nodeName"`
	DeviceType        string            `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo       `json:"acceleratorIDs"`
	RegionName        string            `json:"regionName"`
	FunctionIndex     *int32            `json:"functionIndex,omitempty"`
	Envs              []EnvsInfo        `json:"envs,omitempty"`
	ConfigName        string            `json:"configName"`
	SharedMemory      *SharedMemorySpec `json:"sharedMemory,omitempty"`
	FunctionKernelID  *int32            `json:"functionKernelID"`
	FunctionChannelID *int32            `json:"functionChannelID"`
	PtuKernelID       *int32            `json:"ptuKernelID"`
	FrameworkKernelID *int32            `json:"frameworkKernelID"`
	Rx                RxTxSpec          `json:"rx"`
	Tx                RxTxSpec          `json:"tx"`
}

type RxTxSpec struct {
	Protocol        string  `json:"protocol"`
	IPAddress       *string `json:"ipAddress,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	SubnetAddress   *string `json:"subnetAddress,omitempty"`
	GatewayAddress  *string `json:"gatewayAddress,omitempty"`
	DMAChannelID    *int32  `json:"dmaChannelID,omitempty"`
	FDMAConnectorID *int32  `json:"fdmaConnectorID,omitempty"`
}

// FPGAFunctionStatus defines the observed state of FPGAFunction
type FPGAFunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	StartTime metav1.Time `json:"startTime"`
	//+kubebuilder:default=Pending
	Status              string                `json:"status"`
	DataFlowRef         WBNamespacedName      `json:"dataFlowRef"`
	FunctionName        string                `json:"functionName"`
	ParentBitstreamName string                `json:"parentBitstreamName"`
	ChildBitstreamName  string                `json:"childBitstreamName"`
	SharedMemory        *SharedMemorySpec     `json:"sharedMemory,omitempty"`
	FunctionKernelID    int32                 `json:"functionKernelID"`
	FunctionChannelID   int32                 `json:"functionChannelID"`
	PtuKernelID         int32                 `json:"ptuKernelID"`
	FrameworkKernelID   int32                 `json:"frameworkKernelID"`
	Rx                  RxTxSpec              `json:"rx"`
	Tx                  RxTxSpec              `json:"tx"`
	AcceleratorStatuses []AccStatusesByDevice `json:"acceleratorStatuses,omitempty"`
}

type AccStatusesByDevice struct {
	PartitionName *string       `json:"partitionName,omitempty"`
	Statused      []AccStatuses `json:"statuses,omitempty"`
}

type AccStatuses struct {
	AcceleratorID *string `json:"acceleratorID,omitempty"`
	Status        *string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FPGAFunction is the Schema for the fpgafunctions API
type FPGAFunction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FPGAFunctionSpec   `json:"spec,omitempty"`
	Status FPGAFunctionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FPGAFunctionList contains a list of FPGAFunction
type FPGAFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FPGAFunction `json:"items"`
}

func init() {
	func() *scheme.Builder {
		SchemeBuilderfpga.SchemeBuilder.Register(func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(SchemeBuilderfpga.GroupVersion, []runtime.Object{&FPGAFunction{}, &FPGAFunctionList{}}...)
			metav1.AddToGroupVersion(scheme, SchemeBuilderfpga.GroupVersion)
			return nil
		})
		return SchemeBuilderfpga
	}()
	// SchemeBuilderfpga.Register(&FPGAFunction{}, &FPGAFunctionList{})
}

var (
	// GroupVersion is group version used to register these objects
	GroupVersionfpga = schema.GroupVersion{Group: "example.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilderfpga = &scheme.Builder{GroupVersion: GroupVersionfpga}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToSchemefpga = SchemeBuilderfpga.AddToScheme
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccStatuses) DeepCopyInto(out *AccStatuses) {
	*out = *in
	if in.AcceleratorID != nil {
		in, out := &in.AcceleratorID, &out.AcceleratorID
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccStatuses.
func (in *AccStatuses) DeepCopy() *AccStatuses {
	if in == nil {
		return nil
	}
	out := new(AccStatuses)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccStatusesByDevice) DeepCopyInto(out *AccStatusesByDevice) {
	*out = *in
	if in.PartitionName != nil {
		in, out := &in.PartitionName, &out.PartitionName
		*out = new(string)
		**out = **in
	}
	if in.Statused != nil {
		in, out := &in.Statused, &out.Statused
		*out = make([]AccStatuses, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccStatusesByDevice.
func (in *AccStatusesByDevice) DeepCopy() *AccStatusesByDevice {
	if in == nil {
		return nil
	}
	out := new(AccStatusesByDevice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunction) DeepCopyInto(out *FPGAFunction) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunction.
func (in *FPGAFunction) DeepCopy() *FPGAFunction {
	if in == nil {
		return nil
	}
	out := new(FPGAFunction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FPGAFunction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunctionList) DeepCopyInto(out *FPGAFunctionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FPGAFunction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunctionList.
func (in *FPGAFunctionList) DeepCopy() *FPGAFunctionList {
	if in == nil {
		return nil
	}
	out := new(FPGAFunctionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FPGAFunctionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunctionSpec) DeepCopyInto(out *FPGAFunctionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.AcceleratorIDs != nil {
		in, out := &in.AcceleratorIDs, &out.AcceleratorIDs
		*out = make([]AccIDInfo, len(*in))
		copy(*out, *in)
	}
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]EnvsInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	if in.FunctionKernelID != nil {
		in, out := &in.FunctionKernelID, &out.FunctionKernelID
		*out = new(int32)
		**out = **in
	}
	if in.FunctionChannelID != nil {
		in, out := &in.FunctionChannelID, &out.FunctionChannelID
		*out = new(int32)
		**out = **in
	}
	if in.PtuKernelID != nil {
		in, out := &in.PtuKernelID, &out.PtuKernelID
		*out = new(int32)
		**out = **in
	}
	if in.FrameworkKernelID != nil {
		in, out := &in.FrameworkKernelID, &out.FrameworkKernelID
		*out = new(int32)
		**out = **in
	}
	in.Rx.DeepCopyInto(&out.Rx)
	in.Tx.DeepCopyInto(&out.Tx)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunctionSpec.
func (in *FPGAFunctionSpec) DeepCopy() *FPGAFunctionSpec {
	if in == nil {
		return nil
	}
	out := new(FPGAFunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFunctionStatus) DeepCopyInto(out *FPGAFunctionStatus) {
	*out = *in
	in.StartTime.DeepCopyInto(&out.StartTime)
	out.DataFlowRef = in.DataFlowRef
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	in.Rx.DeepCopyInto(&out.Rx)
	in.Tx.DeepCopyInto(&out.Tx)
	if in.AcceleratorStatuses != nil {
		in, out := &in.AcceleratorStatuses, &out.AcceleratorStatuses
		*out = make([]AccStatusesByDevice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFunctionStatus.
func (in *FPGAFunctionStatus) DeepCopy() *FPGAFunctionStatus {
	if in == nil {
		return nil
	}
	out := new(FPGAFunctionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RxTxSpec) DeepCopyInto(out *RxTxSpec) {
	*out = *in
	if in.IPAddress != nil {
		in, out := &in.IPAddress, &out.IPAddress
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.SubnetAddress != nil {
		in, out := &in.SubnetAddress, &out.SubnetAddress
		*out = new(string)
		**out = **in
	}
	if in.GatewayAddress != nil {
		in, out := &in.GatewayAddress, &out.GatewayAddress
		*out = new(string)
		**out = **in
	}
	if in.DMAChannelID != nil {
		in, out := &in.DMAChannelID, &out.DMAChannelID
		*out = new(int32)
		**out = **in
	}
	if in.FDMAConnectorID != nil {
		in, out := &in.FDMAConnectorID, &out.FDMAConnectorID
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RxTxSpec.
func (in *RxTxSpec) DeepCopy() *RxTxSpec {
	if in == nil {
		return nil
	}
	out := new(RxTxSpec)
	in.DeepCopyInto(out)
	return out
}

// Function CR structure
type FunctionData struct {
	DataFlowRef       WBNamespacedName              `json:"dataFlowRef"`
	FunctionName      string                        `json:"functionName"`
	NodeName          string                        `json:"nodeName"`
	DeviceType        string                        `json:"deviceType"`
	AcceleratorIDs    []AccIDInfo                   `json:"acceleratorIDs"`
	RegionName        string                        `json:"regionName"`
	FunctionIndex     int32                         `json:"functionIndex"`
	Envs              []EnvsInfo                    `json:"envs,omitempty"`
	RequestMemorySize *int32                        `json:"requestMemorySize,omitempty"`
	SharedMemory      SharedMemorySpec              `json:"sharedMemory,omitempty"`
	Protocol          *string                       `json:"protocol,omitempty"`
	ConfigName        *string                       `json:"configName,omitempty"`
	PreviousFunctions map[string]FromToWBFunction   `json:"previousFunctions,omitempty"`
	NextFunctions     map[string]FromToWBFunction   `json:"nextFunctions,omitempty"`
	Params            map[string]intstr.IntOrString `json:"params,omitempty"`
	FunctionKernelID  *int32                        `json:"functionKernelID,omitempty"`
	FunctionChannelID *int32                        `json:"functionChannelID,omitempty"`
	PtuKernelID       *int32                        `json:"ptuKernelID,omitempty"`
	FrameworkKernelID *int32                        `json:"frameworkKernelID,omitempty"`
	Rx                RxTxSpecFunc                  `json:"rx,omitempty"`
	Tx                RxTxSpecFunc                  `json:"tx,omitempty"`
}

// Function CR structure
type FunctionStatusData struct {
	Status string `json:"status"`
}
type AccIDInfo struct {
	PartitionName string `json:"partitionName"`
	ID            string `json:"id"`
}

type EnvsInfo struct {
	PartitionName string     `json:"partitionName"`
	EachEnv       []EnvsData `json:"eachEnv"`
}

// Environmental information
type EnvsData struct {
	EnvKey   string `json:"envKey"`
	EnvValue string `json:"envValue"`
}

// FPGA Device information
type RxTxSpecFunc struct {
	Protocol        string  `json:"protocol,omitempty"`
	IPAddress       *string `json:"ipAddress,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	SubnetAddress   *string `json:"subnetAddress,omitempty"`
	GatewayAddress  *string `json:"gatewayAddress,omitempty"`
	DMAChannelID    *int32  `json:"dmaChannelID,omitempty"`
	FDMAConnectorID *int32  `json:"fdmaConnectorID,omitempty"`
}

// Shared memory information
type SharedMemorySpec struct {
	FilePrefix      string `json:"filePrefix"`
	CommandQueueID  string `json:"commandQueueID"`
	SharedMemoryMiB int32  `json:"sharedMemoryMiB"`
}

// Structure for acquiring Phase3FPGA information
type Phase3Data struct {
	NodeName        string        `json:"nodeName"`
	DeviceFilePaths []string      `json:"deviceFilePaths"`
	NetworkInfo     []NetworkData `json:"networkInfo"`
}

type NetworkData struct {
	DeviceIndex    int32    `json:"deviceIndex"`
	LaneIndex      int32    `json:"laneIndex"`
	IPAddress      string   `json:"ipAddress"`
	SubnetAddress  string   `json:"subnetAddress"`
	GatewayAddress string   `json:"gatewayAddress"`
	MACAddress     string   `json:"macAddress"`
	Rx             FPGARxTx `json:"rx"`
	Tx             FPGARxTx `json:"tx"`
}

type FromToWBFunction struct {
	WBFunctionRef WBNamespacedName `json:"wbFunctionRef"`
	Port          int32            `json:"port"`
}

type FPGARxTx struct {
	Protocol  string `json:"protocol"`
	StartPort int32  `json:"startPort,omitempty"`
	EndPort   int32  `json:"endPort,omitempty"`
}

// This defines NetworkAttachmentDefinition CR
type NetworkAttachmentDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkAttachmentDefinitionSpec `json:"spec,omitempty"`
}

// PCIeConnectionSpec defines the desired state of PCIeConnection
type NetworkAttachmentDefinitionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Config string `json:"config"`
}

// EthernetConnectionList contains a list of EthernetConnection
type NetworkAttachmentDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkAttachmentDefinition `json:"items"`
}

func init() {
	SchemeBuilder1.Register(&NetworkAttachmentDefinition{}, &NetworkAttachmentDefinitionList{})
}

var (
	// GroupVersion is group version used to register these objects
	GroupVersion1 = schema.GroupVersion{Group: "k8s.cni.cncf.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder1 = &scheme.Builder{GroupVersion: GroupVersion1}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme1 = SchemeBuilder1.AddToScheme
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinition) DeepCopyInto(out *NetworkAttachmentDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Spec.DeepCopyInto(&out.Spec)
	//in.Status.DeepCopyInto(&out.Status)

}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnection.
func (in *NetworkAttachmentDefinition) DeepCopy() *NetworkAttachmentDefinition {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkAttachmentDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinitionList) DeepCopyInto(out *NetworkAttachmentDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NetworkAttachmentDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionList.
func (in *NetworkAttachmentDefinitionList) DeepCopy() *NetworkAttachmentDefinitionList {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkAttachmentDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinitionSpec) DeepCopyInto(out *NetworkAttachmentDefinitionSpec) {
	*out = *in
	out.Config = in.Config
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EthernetConnectionSpec.
func (in *NetworkAttachmentDefinitionSpec) DeepCopy() *NetworkAttachmentDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionConfigMap) DeepCopyInto(out *FunctionConfigMap) {
	*out = *in
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionConfigMap.
func (in *FunctionConfigMap) DeepCopy() *FunctionConfigMap {
	if in == nil {
		return nil
	}
	out := new(FunctionConfigMap)
	in.DeepCopyInto(out)
	return out
}

type FunctionConfigMap struct {
	RxProtocol      string            `json:"rxProtocol,omitempty"`
	TxProtocol      string            `json:"txProtocol,omitempty"`
	SharedMemoryMiB int32             `json:"sharedMemoryMiB,omitempty"`
	ImageURI        string            `json:"imageURI,omitempty"`
	Envs            map[string]string `json:"envs,omitempty"`
	ParentBitStream string            `json:"parentBitStream,omitempty"`
	ChildBitStream  string            `json:"childBitStream,omitempty"`
}
