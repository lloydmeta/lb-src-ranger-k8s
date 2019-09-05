// +build !ignore_autogenerated

/*

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

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LbSrcRanger) DeepCopyInto(out *LbSrcRanger) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LbSrcRanger.
func (in *LbSrcRanger) DeepCopy() *LbSrcRanger {
	if in == nil {
		return nil
	}
	out := new(LbSrcRanger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LbSrcRanger) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LbSrcRangerList) DeepCopyInto(out *LbSrcRangerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LbSrcRanger, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LbSrcRangerList.
func (in *LbSrcRangerList) DeepCopy() *LbSrcRangerList {
	if in == nil {
		return nil
	}
	out := new(LbSrcRangerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LbSrcRangerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LbSrcRangerSpec) DeepCopyInto(out *LbSrcRangerSpec) {
	*out = *in
	if in.TargetLabels != nil {
		in, out := &in.TargetLabels, &out.TargetLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.UpdateEvery = in.UpdateEvery
	if in.SrcIPUrls != nil {
		in, out := &in.SrcIPUrls, &out.SrcIPUrls
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LbSrcRangerSpec.
func (in *LbSrcRangerSpec) DeepCopy() *LbSrcRangerSpec {
	if in == nil {
		return nil
	}
	out := new(LbSrcRangerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LbSrcRangerStatus) DeepCopyInto(out *LbSrcRangerStatus) {
	*out = *in
	in.LastRunAt.DeepCopyInto(&out.LastRunAt)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LbSrcRangerStatus.
func (in *LbSrcRangerStatus) DeepCopy() *LbSrcRangerStatus {
	if in == nil {
		return nil
	}
	out := new(LbSrcRangerStatus)
	in.DeepCopyInto(out)
	return out
}
