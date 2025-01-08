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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BsConfigInfo) DeepCopyInto(out *BsConfigInfo) {
	*out = *in
	if in.ChildBitstreamIDs != nil {
		in, out := &in.ChildBitstreamIDs, &out.ChildBitstreamIDs
		*out = make([]ChildBsSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BsConfigInfo.
func (in *BsConfigInfo) DeepCopy() *BsConfigInfo {
	if in == nil {
		return nil
	}
	out := new(BsConfigInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBs) DeepCopyInto(out *ChildBs) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBs.
func (in *ChildBs) DeepCopy() *ChildBs {
	if in == nil {
		return nil
	}
	out := new(ChildBs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChildBs) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsChain) DeepCopyInto(out *ChildBsChain) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(int32)
		**out = **in
	}
	if in.Identifier != nil {
		in, out := &in.Identifier, &out.Identifier
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsChain.
func (in *ChildBsChain) DeepCopy() *ChildBsChain {
	if in == nil {
		return nil
	}
	out := new(ChildBsChain)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsConversion) DeepCopyInto(out *ChildBsConversion) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(int32)
		**out = **in
	}
	if in.Module != nil {
		in, out := &in.Module, &out.Module
		*out = new([]ConversionModule)
		if **in != nil {
			in, out := *in, *out
			*out = make([]ConversionModule, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsConversion.
func (in *ChildBsConversion) DeepCopy() *ChildBsConversion {
	if in == nil {
		return nil
	}
	out := new(ChildBsConversion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsDirecttrans) DeepCopyInto(out *ChildBsDirecttrans) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(int32)
		**out = **in
	}
	if in.Identifier != nil {
		in, out := &in.Identifier, &out.Identifier
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsDirecttrans.
func (in *ChildBsDirecttrans) DeepCopy() *ChildBsDirecttrans {
	if in == nil {
		return nil
	}
	out := new(ChildBsDirecttrans)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsFunctions) DeepCopyInto(out *ChildBsFunctions) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(int32)
		**out = **in
	}
	if in.Module != nil {
		in, out := &in.Module, &out.Module
		*out = new([]FunctionsModule)
		if **in != nil {
			in, out := *in, *out
			*out = make([]FunctionsModule, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = new(map[string]intstr.IntOrString)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]intstr.IntOrString, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.IntraResourceMgmtMap != nil {
		in, out := &in.IntraResourceMgmtMap, &out.IntraResourceMgmtMap
		*out = new(map[string]FunctionsIntraResourceMgmtMap)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]FunctionsIntraResourceMgmtMap, len(*in))
			for key, val := range *in {
				(*out)[key] = *val.DeepCopy()
			}
		}
	}
	in.DeploySpec.DeepCopyInto(&out.DeploySpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsFunctions.
func (in *ChildBsFunctions) DeepCopy() *ChildBsFunctions {
	if in == nil {
		return nil
	}
	out := new(ChildBsFunctions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsLLDMA) DeepCopyInto(out *ChildBsLLDMA) {
	*out = *in
	if in.Cids != nil {
		in, out := &in.Cids, &out.Cids
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsLLDMA.
func (in *ChildBsLLDMA) DeepCopy() *ChildBsLLDMA {
	if in == nil {
		return nil
	}
	out := new(ChildBsLLDMA)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsList) DeepCopyInto(out *ChildBsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChildBs, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsList.
func (in *ChildBsList) DeepCopy() *ChildBsList {
	if in == nil {
		return nil
	}
	out := new(ChildBsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChildBsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsModule) DeepCopyInto(out *ChildBsModule) {
	*out = *in
	if in.Ptu != nil {
		in, out := &in.Ptu, &out.Ptu
		*out = new(ChildBsPtu)
		(*in).DeepCopyInto(*out)
	}
	if in.LLDMA != nil {
		in, out := &in.LLDMA, &out.LLDMA
		*out = new(ChildBsLLDMA)
		(*in).DeepCopyInto(*out)
	}
	if in.Chain != nil {
		in, out := &in.Chain, &out.Chain
		*out = new(ChildBsChain)
		(*in).DeepCopyInto(*out)
	}
	if in.Directtrans != nil {
		in, out := &in.Directtrans, &out.Directtrans
		*out = new(ChildBsDirecttrans)
		(*in).DeepCopyInto(*out)
	}
	if in.Conversion != nil {
		in, out := &in.Conversion, &out.Conversion
		*out = new(ChildBsConversion)
		(*in).DeepCopyInto(*out)
	}
	if in.Functions != nil {
		in, out := &in.Functions, &out.Functions
		*out = new([]ChildBsFunctions)
		if **in != nil {
			in, out := *in, *out
			*out = make([]ChildBsFunctions, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsModule.
func (in *ChildBsModule) DeepCopy() *ChildBsModule {
	if in == nil {
		return nil
	}
	out := new(ChildBsModule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsPtu) DeepCopyInto(out *ChildBsPtu) {
	*out = *in
	if in.Cids != nil {
		in, out := &in.Cids, &out.Cids
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(int32)
		**out = **in
	}
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = new(map[string]intstr.IntOrString)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]intstr.IntOrString, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsPtu.
func (in *ChildBsPtu) DeepCopy() *ChildBsPtu {
	if in == nil {
		return nil
	}
	out := new(ChildBsPtu)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsRegion) DeepCopyInto(out *ChildBsRegion) {
	*out = *in
	if in.Modules != nil {
		in, out := &in.Modules, &out.Modules
		*out = new(ChildBsModule)
		(*in).DeepCopyInto(*out)
	}
	if in.MaxFunctions != nil {
		in, out := &in.MaxFunctions, &out.MaxFunctions
		*out = new(int32)
		**out = **in
	}
	if in.MaxCapacity != nil {
		in, out := &in.MaxCapacity, &out.MaxCapacity
		*out = new(int32)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsRegion.
func (in *ChildBsRegion) DeepCopy() *ChildBsRegion {
	if in == nil {
		return nil
	}
	out := new(ChildBsRegion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsSpec) DeepCopyInto(out *ChildBsSpec) {
	*out = *in
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]ChildBsRegion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ChildBitstreamID != nil {
		in, out := &in.ChildBitstreamID, &out.ChildBitstreamID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsSpec.
func (in *ChildBsSpec) DeepCopy() *ChildBsSpec {
	if in == nil {
		return nil
	}
	out := new(ChildBsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildBsStatus) DeepCopyInto(out *ChildBsStatus) {
	*out = *in
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]ChildBsRegion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ChildBitstreamID != nil {
		in, out := &in.ChildBitstreamID, &out.ChildBitstreamID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildBsStatus.
func (in *ChildBsStatus) DeepCopy() *ChildBsStatus {
	if in == nil {
		return nil
	}
	out := new(ChildBsStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComputeResource) DeepCopyInto(out *ComputeResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComputeResource.
func (in *ComputeResource) DeepCopy() *ComputeResource {
	if in == nil {
		return nil
	}
	out := new(ComputeResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ComputeResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComputeResourceList) DeepCopyInto(out *ComputeResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ComputeResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComputeResourceList.
func (in *ComputeResourceList) DeepCopy() *ComputeResourceList {
	if in == nil {
		return nil
	}
	out := new(ComputeResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ComputeResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComputeResourceSpec) DeepCopyInto(out *ComputeResourceSpec) {
	*out = *in
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]RegionInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComputeResourceSpec.
func (in *ComputeResourceSpec) DeepCopy() *ComputeResourceSpec {
	if in == nil {
		return nil
	}
	out := new(ComputeResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComputeResourceStatus) DeepCopyInto(out *ComputeResourceStatus) {
	*out = *in
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]RegionInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComputeResourceStatus.
func (in *ComputeResourceStatus) DeepCopy() *ComputeResourceStatus {
	if in == nil {
		return nil
	}
	out := new(ComputeResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConversionModule) DeepCopyInto(out *ConversionModule) {
	*out = *in
	if in.Identifier != nil {
		in, out := &in.Identifier, &out.Identifier
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConversionModule.
func (in *ConversionModule) DeepCopy() *ConversionModule {
	if in == nil {
		return nil
	}
	out := new(ConversionModule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Details) DeepCopyInto(out *Details) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.DMAChannelID != nil {
		in, out := &in.DMAChannelID, &out.DMAChannelID
		*out = new(int32)
		**out = **in
	}
	if in.LLDMAConnectorID != nil {
		in, out := &in.LLDMAConnectorID, &out.LLDMAConnectorID
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Details.
func (in *Details) DeepCopy() *Details {
	if in == nil {
		return nil
	}
	out := new(Details)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceInfo) DeepCopyInto(out *DeviceInfo) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceInfo.
func (in *DeviceInfo) DeepCopy() *DeviceInfo {
	if in == nil {
		return nil
	}
	out := new(DeviceInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeviceInfo) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceInfoList) DeepCopyInto(out *DeviceInfoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DeviceInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceInfoList.
func (in *DeviceInfoList) DeepCopy() *DeviceInfoList {
	if in == nil {
		return nil
	}
	out := new(DeviceInfoList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeviceInfoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceInfoSpec) DeepCopyInto(out *DeviceInfoSpec) {
	*out = *in
	in.Request.DeepCopyInto(&out.Request)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceInfoSpec.
func (in *DeviceInfoSpec) DeepCopy() *DeviceInfoSpec {
	if in == nil {
		return nil
	}
	out := new(DeviceInfoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceInfoStatus) DeepCopyInto(out *DeviceInfoStatus) {
	*out = *in
	in.Response.DeepCopyInto(&out.Response)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceInfoStatus.
func (in *DeviceInfoStatus) DeepCopy() *DeviceInfoStatus {
	if in == nil {
		return nil
	}
	out := new(DeviceInfoStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGA) DeepCopyInto(out *FPGA) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGA.
func (in *FPGA) DeepCopy() *FPGA {
	if in == nil {
		return nil
	}
	out := new(FPGA)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FPGA) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAList) DeepCopyInto(out *FPGAList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FPGA, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAList.
func (in *FPGAList) DeepCopy() *FPGAList {
	if in == nil {
		return nil
	}
	out := new(FPGAList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FPGAList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGASpec) DeepCopyInto(out *FPGASpec) {
	*out = *in
	if in.ChildBitstreamID != nil {
		in, out := &in.ChildBitstreamID, &out.ChildBitstreamID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGASpec.
func (in *FPGASpec) DeepCopy() *FPGASpec {
	if in == nil {
		return nil
	}
	out := new(FPGASpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGAStatus) DeepCopyInto(out *FPGAStatus) {
	*out = *in
	if in.ChildBitstreamID != nil {
		in, out := &in.ChildBitstreamID, &out.ChildBitstreamID
		*out = new(string)
		**out = **in
	}
	if in.ChildBitstreamCRName != nil {
		in, out := &in.ChildBitstreamCRName, &out.ChildBitstreamCRName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FPGAStatus.
func (in *FPGAStatus) DeepCopy() *FPGAStatus {
	if in == nil {
		return nil
	}
	out := new(FPGAStatus)
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
func (in *FunctionInfrastruct) DeepCopyInto(out *FunctionInfrastruct) {
	*out = *in
	if in.MaxDataFlows != nil {
		in, out := &in.MaxDataFlows, &out.MaxDataFlows
		*out = new(int32)
		**out = **in
	}
	if in.CurrentDataFlows != nil {
		in, out := &in.CurrentDataFlows, &out.CurrentDataFlows
		*out = new(int32)
		**out = **in
	}
	if in.MaxCapacity != nil {
		in, out := &in.MaxCapacity, &out.MaxCapacity
		*out = new(int32)
		**out = **in
	}
	if in.CurrentCapacity != nil {
		in, out := &in.CurrentCapacity, &out.CurrentCapacity
		*out = new(int32)
		**out = **in
	}
	if in.MaxTimeSlicingSeconds != nil {
		in, out := &in.MaxTimeSlicingSeconds, &out.MaxTimeSlicingSeconds
		*out = new(int32)
		**out = **in
	}
	if in.CurrentTimeSlicingSeconds != nil {
		in, out := &in.CurrentTimeSlicingSeconds, &out.CurrentTimeSlicingSeconds
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionInfrastruct.
func (in *FunctionInfrastruct) DeepCopy() *FunctionInfrastruct {
	if in == nil {
		return nil
	}
	out := new(FunctionInfrastruct)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionsDeploySpec) DeepCopyInto(out *FunctionsDeploySpec) {
	*out = *in
	if in.MaxCapacity != nil {
		in, out := &in.MaxCapacity, &out.MaxCapacity
		*out = new(int32)
		**out = **in
	}
	if in.MaxDataFlows != nil {
		in, out := &in.MaxDataFlows, &out.MaxDataFlows
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionsDeploySpec.
func (in *FunctionsDeploySpec) DeepCopy() *FunctionsDeploySpec {
	if in == nil {
		return nil
	}
	out := new(FunctionsDeploySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionsIntraResourceMgmtMap) DeepCopyInto(out *FunctionsIntraResourceMgmtMap) {
	*out = *in
	if in.Available != nil {
		in, out := &in.Available, &out.Available
		*out = new(bool)
		**out = **in
	}
	if in.FunctionCRName != nil {
		in, out := &in.FunctionCRName, &out.FunctionCRName
		*out = new(string)
		**out = **in
	}
	if in.Rx != nil {
		in, out := &in.Rx, &out.Rx
		*out = new(RxTxSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Tx != nil {
		in, out := &in.Tx, &out.Tx
		*out = new(RxTxSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionsIntraResourceMgmtMap.
func (in *FunctionsIntraResourceMgmtMap) DeepCopy() *FunctionsIntraResourceMgmtMap {
	if in == nil {
		return nil
	}
	out := new(FunctionsIntraResourceMgmtMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionsModule) DeepCopyInto(out *FunctionsModule) {
	*out = *in
	if in.FunctionChannelIDs != nil {
		in, out := &in.FunctionChannelIDs, &out.FunctionChannelIDs
		*out = new(string)
		**out = **in
	}
	if in.Identifier != nil {
		in, out := &in.Identifier, &out.Identifier
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionsModule.
func (in *FunctionsModule) DeepCopy() *FunctionsModule {
	if in == nil {
		return nil
	}
	out := new(FunctionsModule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegionInfo) DeepCopyInto(out *RegionInfo) {
	*out = *in
	if in.DeviceUUID != nil {
		in, out := &in.DeviceUUID, &out.DeviceUUID
		*out = new(string)
		**out = **in
	}
	if in.MaxFunctions != nil {
		in, out := &in.MaxFunctions, &out.MaxFunctions
		*out = new(int32)
		**out = **in
	}
	if in.CurrentFunctions != nil {
		in, out := &in.CurrentFunctions, &out.CurrentFunctions
		*out = new(int32)
		**out = **in
	}
	if in.MaxCapacity != nil {
		in, out := &in.MaxCapacity, &out.MaxCapacity
		*out = new(int32)
		**out = **in
	}
	if in.CurrentCapacity != nil {
		in, out := &in.CurrentCapacity, &out.CurrentCapacity
		*out = new(int32)
		**out = **in
	}
	if in.MaxTimeSlicingSeconds != nil {
		in, out := &in.MaxTimeSlicingSeconds, &out.MaxTimeSlicingSeconds
		*out = new(int32)
		**out = **in
	}
	if in.CurrentTimeSlicingSeconds != nil {
		in, out := &in.CurrentTimeSlicingSeconds, &out.CurrentTimeSlicingSeconds
		*out = new(int32)
		**out = **in
	}
	if in.Functions != nil {
		in, out := &in.Functions, &out.Functions
		*out = make([]FunctionInfrastruct, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegionInfo.
func (in *RegionInfo) DeepCopy() *RegionInfo {
	if in == nil {
		return nil
	}
	out := new(RegionInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RxTxSpec) DeepCopyInto(out *RxTxSpec) {
	*out = *in
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(map[string]Details)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]Details, len(*in))
			for key, val := range *in {
				(*out)[key] = *val.DeepCopy()
			}
		}
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
func (in *WBFuncRequest) DeepCopyInto(out *WBFuncRequest) {
	*out = *in
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
	if in.MaxDataFlows != nil {
		in, out := &in.MaxDataFlows, &out.MaxDataFlows
		*out = new(int32)
		**out = **in
	}
	if in.MaxCapacity != nil {
		in, out := &in.MaxCapacity, &out.MaxCapacity
		*out = new(int32)
		**out = **in
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBFuncRequest.
func (in *WBFuncRequest) DeepCopy() *WBFuncRequest {
	if in == nil {
		return nil
	}
	out := new(WBFuncRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBFuncResponse) DeepCopyInto(out *WBFuncResponse) {
	*out = *in
	if in.FunctionIndex != nil {
		in, out := &in.FunctionIndex, &out.FunctionIndex
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBFuncResponse.
func (in *WBFuncResponse) DeepCopy() *WBFuncResponse {
	if in == nil {
		return nil
	}
	out := new(WBFuncResponse)
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