// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

const (
	// ConfigFinalizer allows us to clean up resources before deletion
	ConfigFinalizer = "talosconfig.bootstrap.cluster.x-k8s.io"
)

// TalosConfigSpec defines the desired state of TalosConfig
type TalosConfigSpec struct {
	TalosVersion  string          `json:"talosVersion,omitempty"` //talos version formatted like v0.8. used for backwards compatibility
	GenerateType  string          `json:"generateType"`           //none,init,controlplane,worker mutually exclusive w/ data
	Data          string          `json:"data,omitempty"`
	ConfigPatches []ConfigPatches `json:"configPatches,omitempty"`
	// Talos Linux machine configuration strategic merge patch list.
	StrategicPatches []string `json:"strategicPatches,omitempty"`
	// Set hostname in the machine configuration to some value.
	Hostname HostnameSpec `json:"hostname,omitempty"`
	// Important: Run "make" to regenerate code after modifying this file
}

// HostnameSource is the definition of hostname source.
type HostnameSource string

// HostnameSourceMachineName sets the hostname in the generated configuration to the machine name.
const HostnameSourceMachineName HostnameSource = "MachineName"

// HostnameSourceInfrastructureName sets the hostname in the generated configuration to the name of the machine's infrastructure.
const HostnameSourceInfrastructureName HostnameSource = "InfrastructureName"

// HostnameSpec defines the hostname source.
type HostnameSpec struct {
	// Source of the hostname.
	//
	// Allowed values:
	// "MachineName" (use linked Machine's Name).
	// "InfrastructureName" (use linked Machine's infrastructure's name).
	Source HostnameSource `json:"source,omitempty"`
}

// TalosConfigStatus defines the observed state of TalosConfig
type TalosConfigStatus struct {
	// Ready indicates the BootstrapData field is ready to be consumed
	Ready bool `json:"ready,omitempty"`

	// DataSecretName is the name of the secret that stores the bootstrap data script.
	// +optional
	DataSecretName *string `json:"dataSecretName,omitempty"`

	// Talos config will be a string containing the config for download.
	//
	// Deprecated: please use `<cluster>-talosconfig` secret.
	//
	// +optional
	TalosConfig string `json:"talosConfig,omitempty"`

	// FailureReason will be set on non-retryable errors
	//
	// Deprecated: this field will be removed in the next apiVersion.
	//
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage will be set on non-retryable errors
	//
	// Deprecated: this field will be removed in the next apiVersion.
	//
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// ObservedGeneration is the latest generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions defines current service state of the TalosConfig.
	// +optional
	Conditions capiv1.Conditions `json:"conditions,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in TalosConfig's status with the V1Beta2 version.
	// +optional
	V1Beta2 *TalosConfigV1Beta2Status `json:"v1beta2Status,omitempty"`
}

type TalosConfigV1Beta2Status struct {
	// Conditions represents the observations of a TalosConfig's current state.
	// Known condition types are Ready, DataSecretAvailable, ClientConfigAvailable.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=talosconfigs,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status

// TalosConfig is the Schema for the talosconfigs API
type TalosConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TalosConfigSpec   `json:"spec,omitempty"`
	Status TalosConfigStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (c *TalosConfig) GetConditions() capiv1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (c *TalosConfig) SetConditions(conditions capiv1.Conditions) {
	c.Status.Conditions = conditions
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (c *TalosConfig) GetV1Beta2Conditions() []metav1.Condition {
	if c.Status.V1Beta2 == nil {
		return nil
	}
	return c.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (c *TalosConfig) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if c.Status.V1Beta2 == nil {
		c.Status.V1Beta2 = &TalosConfigV1Beta2Status{}
	}
	c.Status.V1Beta2.Conditions = conditions
}

// +kubebuilder:object:root=true

// TalosConfigList contains a list of TalosConfig
type TalosConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TalosConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TalosConfig{}, &TalosConfigList{})
}
