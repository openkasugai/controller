/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Width            = "WIDTH"
	Height           = "HEIGHT"
	ArgsWidth        = "%" + Width + "%"
	ArgsHeight       = "%" + Height + "%"
	ArgsArpIP        = "%ARPIP%"
	ArgsInputIP      = "%INPUTIP%"
	ArgsInputPort    = "%INPUTPORT%"
	ArgsOutputIP     = "%OUTPUTIP%"
	ArgsOutputPort   = "%OUTPUTPORT%"
	ArgsArpMAC       = "%MAC%"
	ChangeArgsWidth  = "width="
	ChangeArgsHeight = "height="
	ChangeArgsIP     = "host="
	ChangeArgsPort   = "port="
)

type GPUFuncConfig struct {
	RxProtocol                     *string           `json:"rxProtocol,omitempty"`
	TxProtocol                     *string           `json:"txProtocol,omitempty"`
	SharedMemoryGiB                *int32            `json:"sharedMemoryGiB,omitempty"`
	VirtualNetworkDeviceDriverType string            `json:"virtualNetworkDeviceDriverType,omitempty"`
	AdditionalNetwork              bool              `json:"additionalNetwork,omitempty"`
	ImageURI                       string            `json:"imageURI"`
	Envs                           map[string]string `json:"envs"`
	Template                       PodTemplate       `json:"template"`
}

type PodTemplate struct {
	metav1.TypeMeta `json:",inline"`
	Spec            GPUPodSpec `json:"spec"`
	// Spec corev1.PodSpec `json:"spec,omitempty"`
}

type GPUPodSpec struct {
	Volumes       []corev1.Volume      `json:"volumes,omitempty" `
	Containers    []GPUContainer       `json:"containers"`
	RestartPolicy corev1.RestartPolicy `json:"restartPolicy,omitempty"`
	HostNetwork   bool                 `json:"hostNetwork,omitempty"`
	HostIPC       bool                 `json:"hostIPC,omitempty"`
}

type GPUContainer struct {
	Name            *string                     `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Command         []string                    `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	Args            []string                    `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	WorkingDir      string                      `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	SecurityContext *corev1.SecurityContext     `json:"securityContext,omitempty" protobuf:"bytes,15,opt,name=securityContext"`
	VolumeMounts    []corev1.VolumeMount        `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
}
