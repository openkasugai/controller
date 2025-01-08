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
func (in *WBConnection) DeepCopyInto(out *WBConnection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnection.
func (in *WBConnection) DeepCopy() *WBConnection {
	if in == nil {
		return nil
	}
	out := new(WBConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WBConnection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionEdge) DeepCopyInto(out *WBConnectionEdge) {
	*out = *in
	out.From = in.From
	out.To = in.To
	if in.IntParams != nil {
		in, out := &in.IntParams, &out.IntParams
		*out = make(map[string]int, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.StrParams != nil {
		in, out := &in.StrParams, &out.StrParams
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionEdge.
func (in *WBConnectionEdge) DeepCopy() *WBConnectionEdge {
	if in == nil {
		return nil
	}
	out := new(WBConnectionEdge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionIO) DeepCopyInto(out *WBConnectionIO) {
	*out = *in
	if in.IntParams != nil {
		in, out := &in.IntParams, &out.IntParams
		*out = make(map[string]int, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.StrParams != nil {
		in, out := &in.StrParams, &out.StrParams
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionIO.
func (in *WBConnectionIO) DeepCopy() *WBConnectionIO {
	if in == nil {
		return nil
	}
	out := new(WBConnectionIO)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionList) DeepCopyInto(out *WBConnectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WBConnection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionList.
func (in *WBConnectionList) DeepCopy() *WBConnectionList {
	if in == nil {
		return nil
	}
	out := new(WBConnectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WBConnectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionPath) DeepCopyInto(out *WBConnectionPath) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionPath.
func (in *WBConnectionPath) DeepCopy() *WBConnectionPath {
	if in == nil {
		return nil
	}
	out := new(WBConnectionPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionRequirementsInfo) DeepCopyInto(out *WBConnectionRequirementsInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionRequirementsInfo.
func (in *WBConnectionRequirementsInfo) DeepCopy() *WBConnectionRequirementsInfo {
	if in == nil {
		return nil
	}
	out := new(WBConnectionRequirementsInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionSpec) DeepCopyInto(out *WBConnectionSpec) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.ConnectionPath != nil {
		in, out := &in.ConnectionPath, &out.ConnectionPath
		*out = make([]WBConnectionPath, len(*in))
		copy(*out, *in)
	}
	out.From = in.From
	out.To = in.To
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]intstr.IntOrString, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Requirements != nil {
		in, out := &in.Requirements, &out.Requirements
		*out = new(WBConnectionRequirementsInfo)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionSpec.
func (in *WBConnectionSpec) DeepCopy() *WBConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(WBConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WBConnectionStatus) DeepCopyInto(out *WBConnectionStatus) {
	*out = *in
	out.DataFlowRef = in.DataFlowRef
	if in.ConnectionPath != nil {
		in, out := &in.ConnectionPath, &out.ConnectionPath
		*out = make([]WBConnectionPath, len(*in))
		copy(*out, *in)
	}
	out.From = in.From
	out.To = in.To
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]intstr.IntOrString, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SatisfiedRequirements != nil {
		in, out := &in.SatisfiedRequirements, &out.SatisfiedRequirements
		*out = new(WBConnectionRequirementsInfo)
		**out = **in
	}
	if in.IOs != nil {
		in, out := &in.IOs, &out.IOs
		*out = make(map[string]WBConnectionIO, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Connections != nil {
		in, out := &in.Connections, &out.Connections
		*out = make([]WBConnectionEdge, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WBConnectionStatus.
func (in *WBConnectionStatus) DeepCopy() *WBConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(WBConnectionStatus)
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
