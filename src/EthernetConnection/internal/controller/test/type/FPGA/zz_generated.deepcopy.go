//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccIDInfo) DeepCopyInto(out *AccIDInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccIDInfo.
func (in *AccIDInfo) DeepCopy() *AccIDInfo {
	if in == nil {
		return nil
	}
	out := new(AccIDInfo)
	in.DeepCopyInto(out)
	return out
}

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
func (in *EnvsData) DeepCopyInto(out *EnvsData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvsData.
func (in *EnvsData) DeepCopy() *EnvsData {
	if in == nil {
		return nil
	}
	out := new(EnvsData)
	in.DeepCopyInto(out)
	return out
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAFuncConfig) DeepCopyInto(out *FPGAFuncConfig) {
	*out = *in
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Commands != nil {
		in, out := &in.Commands, &out.Commands
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAFuncConfig.
func (in *FPGAFuncConfig) DeepCopy() *FPGAFuncConfig {
	if in == nil {
		return nil
	}
	out := new(FPGAFuncConfig)
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
func (in *FPGARxTx) DeepCopyInto(out *FPGARxTx) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGARxTx.
func (in *FPGARxTx) DeepCopy() *FPGARxTx {
	if in == nil {
		return nil
	}
	out := new(FPGARxTx)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FromToWBFunction) DeepCopyInto(out *FromToWBFunction) {
	*out = *in
	out.WBFunctionRef = in.WBFunctionRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FromToWBFunction.
func (in *FromToWBFunction) DeepCopy() *FromToWBFunction {
	if in == nil {
		return nil
	}
	out := new(FromToWBFunction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionData) DeepCopyInto(out *FunctionData) {
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
	if in.RequestMemorySize != nil {
		in, out := &in.RequestMemorySize, &out.RequestMemorySize
		*out = new(int32)
		**out = **in
	}
	if in.SharedMemory != nil {
		in, out := &in.SharedMemory, &out.SharedMemory
		*out = new(SharedMemorySpec)
		**out = **in
	}
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(string)
		**out = **in
	}
	in.Rx.DeepCopyInto(&out.Rx)
	in.Tx.DeepCopyInto(&out.Tx)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionData.
func (in *FunctionData) DeepCopy() *FunctionData {
	if in == nil {
		return nil
	}
	out := new(FunctionData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkData) DeepCopyInto(out *NetworkData) {
	*out = *in
	out.Rx = in.Rx
	out.Tx = in.Tx
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkData.
func (in *NetworkData) DeepCopy() *NetworkData {
	if in == nil {
		return nil
	}
	out := new(NetworkData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Phase3Data) DeepCopyInto(out *Phase3Data) {
	*out = *in
	if in.DeviceFilePaths != nil {
		in, out := &in.DeviceFilePaths, &out.DeviceFilePaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NetworkInfo != nil {
		in, out := &in.NetworkInfo, &out.NetworkInfo
		*out = make([]NetworkData, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Phase3Data.
func (in *Phase3Data) DeepCopy() *Phase3Data {
	if in == nil {
		return nil
	}
	out := new(Phase3Data)
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RxTxSpecFunc) DeepCopyInto(out *RxTxSpecFunc) {
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

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RxTxSpecFunc.
func (in *RxTxSpecFunc) DeepCopy() *RxTxSpecFunc {
	if in == nil {
		return nil
	}
	out := new(RxTxSpecFunc)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SharedMemorySpec) DeepCopyInto(out *SharedMemorySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SharedMemorySpec.
func (in *SharedMemorySpec) DeepCopy() *SharedMemorySpec {
	if in == nil {
		return nil
	}
	out := new(SharedMemorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBNamespacedName) DeepCopyInto(out *WBNamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBNamespacedName.
func (in *WBNamespacedName) DeepCopy() *WBNamespacedName {
	if in == nil {
		return nil
	}
	out := new(WBNamespacedName)
	in.DeepCopyInto(out)
	return out
}