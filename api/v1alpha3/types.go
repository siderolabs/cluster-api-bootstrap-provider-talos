package v1alpha3

import apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

// TalosConfigTemplateResource defines the Template structure
type TalosConfigTemplateResource struct {
	Spec TalosConfigSpec `json:"spec,omitempty"`
}

// nb: we use apiextensions.JSON for the value below b/c we can't use interface{} with controller-gen.
// found this workaround here: https://github.com/kubernetes-sigs/controller-tools/pull/126#issuecomment-630769075

type ConfigPatches struct {
	Op    string             `json:"op"`
	Path  string             `json:"path"`
	Value apiextensions.JSON `json:"value,omitempty"`
}
