// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TalosConfigTemplateSpec defines the desired state of TalosConfigTemplate
type TalosConfigTemplateSpec struct {
	Template TalosConfigTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=talosconfigtemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ClusterClass",type="string",JSONPath=`.metadata.ownerReferences[?(@.kind=="ClusterClass")].name`,description="Name of the ClusterClass owning this template"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=`.metadata.ownerReferences[?(@.kind=="Cluster")].name`,description="Name of the Cluster owning this template"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of TalosConfigTemplate"

// TalosConfigTemplate is the Schema for the talosconfigtemplates API
type TalosConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TalosConfigTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// TalosConfigTemplateList contains a list of TalosConfigTemplate
type TalosConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TalosConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TalosConfigTemplate{}, &TalosConfigTemplateList{})
}
