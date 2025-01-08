//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package reference

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGA) DeepCopyInto(out *FPGA) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
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
	return
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FPGASpec) DeepCopyInto(out *FPGASpec) {
	*out = *in
	if in.ChildBitstreamID != nil {
		in, out := &in.ChildBitstreamID, &out.ChildBitstreamID
		*out = new(string)
		**out = **in
	}
	return
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
	if in.ChildBitStreamID != nil {
		in, out := &in.ChildBitStreamID, &out.ChildBitStreamID
		*out = new(string)
		**out = **in
	}
	if in.ChildBitStreamCRName != nil {
		in, out := &in.ChildBitStreamCRName, &out.ChildBitStreamCRName
		*out = new(string)
		**out = **in
	}
	return
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
