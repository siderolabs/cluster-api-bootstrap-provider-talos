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

package v1alpha3

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfig) DeepCopyInto(out *TalosConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfig.
func (in *TalosConfig) DeepCopy() *TalosConfig {
	if in == nil {
		return nil
	}
	out := new(TalosConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TalosConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigList) DeepCopyInto(out *TalosConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TalosConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigList.
func (in *TalosConfigList) DeepCopy() *TalosConfigList {
	if in == nil {
		return nil
	}
	out := new(TalosConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TalosConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigSpec) DeepCopyInto(out *TalosConfigSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigSpec.
func (in *TalosConfigSpec) DeepCopy() *TalosConfigSpec {
	if in == nil {
		return nil
	}
	out := new(TalosConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigStatus) DeepCopyInto(out *TalosConfigStatus) {
	*out = *in
	if in.DataSecretName != nil {
		in, out := &in.DataSecretName, &out.DataSecretName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigStatus.
func (in *TalosConfigStatus) DeepCopy() *TalosConfigStatus {
	if in == nil {
		return nil
	}
	out := new(TalosConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigTemplate) DeepCopyInto(out *TalosConfigTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigTemplate.
func (in *TalosConfigTemplate) DeepCopy() *TalosConfigTemplate {
	if in == nil {
		return nil
	}
	out := new(TalosConfigTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TalosConfigTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigTemplateList) DeepCopyInto(out *TalosConfigTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TalosConfigTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigTemplateList.
func (in *TalosConfigTemplateList) DeepCopy() *TalosConfigTemplateList {
	if in == nil {
		return nil
	}
	out := new(TalosConfigTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TalosConfigTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigTemplateResource) DeepCopyInto(out *TalosConfigTemplateResource) {
	*out = *in
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigTemplateResource.
func (in *TalosConfigTemplateResource) DeepCopy() *TalosConfigTemplateResource {
	if in == nil {
		return nil
	}
	out := new(TalosConfigTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TalosConfigTemplateSpec) DeepCopyInto(out *TalosConfigTemplateSpec) {
	*out = *in
	out.Template = in.Template
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TalosConfigTemplateSpec.
func (in *TalosConfigTemplateSpec) DeepCopy() *TalosConfigTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(TalosConfigTemplateSpec)
	in.DeepCopyInto(out)
	return out
}
