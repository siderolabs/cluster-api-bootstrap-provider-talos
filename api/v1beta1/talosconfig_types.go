// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// ConfigFinalizer allows us to clean up resources before deletion
	ConfigFinalizer = "talosconfig.bootstrap.cluster.x-k8s.io"
)

// TalosConfigSpec defines the desired state of TalosConfig
type TalosConfigSpec struct {
	// Version of Talos to generate machine configuration for, i.e. v1.0 (patch version may be omitted).
	// Defaults to the latest supported version if unspecified.
	// It is recommended to set this field explicitly to avoid unexpected issues during provider upgrades.
	TalosVersion string `json:"talosVersion,omitempty"`

	// Talos machine configuration type to generate: Supported values are:
	// controlplane, worker, init (deprecated), none (configuration must be then provided in data field).
	// +kubebuilder:validation:Enum=controlplane;worker;init;none
	GenerateType string `json:"generateType"`

	// Machine configuration in case generateType=none.
	Data string `json:"data,omitempty"`

	// RFC6902 JSON patches to apply to machine configuration.
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
	// initialization provides observations of the TalosConfig initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Machine provisioning.
	// +optional
	Initialization TalosConfigInitializationStatus `json:"initialization,omitempty,omitzero"`

	// DataSecretName is the name of the secret that stores the bootstrap data script.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	DataSecretName string `json:"dataSecretName,omitempty"`

	// ObservedGeneration is the latest generation observed by the controller.
	// +optional
	// kubebuilder:validation:Minimum=1
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represents the observations of a TalosConfig's current state.
	// Known condition types are Ready, DataSecretAvailable, ClientConfigAvailable.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *TalosConfigDeprecatedStatus `json:"deprecated,omitempty"`
}

// TalosConfigInitializationStatus provides observations of the TalosConfig initialization process.
// +kubebuilder:validation:MinProperties=1
type TalosConfigInitializationStatus struct {
	// dataSecretCreated is true when the Machine's bootstrap secret is created.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Machine provisioning.
	// +optional
	DataSecretCreated *bool `json:"dataSecretCreated,omitempty"`
}

// TalosConfigDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
type TalosConfigDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for CAPI v1beta1 contract will be dropped.
	// +optional
	V1Beta1 *TalosConfigV1Beta1DeprecatedStatus `json:"v1beta,omitempty"`
}

// TalosConfigV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for CAPI v1beta1 contract will be dropped.
type TalosConfigV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the TalosConfig.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for CAPI v1beta1 contract will be dropped.
	//
	// +optional
	Conditions capiv1.Conditions `json:"conditions,omitempty"`

	// failureReason will be set on non-retryable errors
	//
	// Deprecated: This field is deprecated and is going to be removed when support for CAPI v1beta1 contract will be dropped.
	//
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// failureMessage will be set on non-retryable errors
	//
	// Deprecated: This field is deprecated and is going to be removed when support for CAPI v1beta1 contract will be dropped.
	//
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=talosconfigs,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels['cluster\\.x-k8s\\.io/cluster-name']",description="Cluster"
// +kubebuilder:printcolumn:name="Initialized",type="string",JSONPath=`.status.initialization.dataSecretCreated`,description="Bootstrap secret is created"
// +kubebuilder:printcolumn:name="Bootstrap secret",type="string",JSONPath=`.status.dataSecretName`,description="Bootstrap secret name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of TalosConfig"

// TalosConfig is the Schema for the talosconfigs API
type TalosConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TalosConfigSpec   `json:"spec,omitempty"`
	Status TalosConfigStatus `json:"status,omitempty"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *TalosConfig) GetV1Beta1Conditions() capiv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *TalosConfig) SetV1Beta1Conditions(conditions capiv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &TalosConfigDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &TalosConfigV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *TalosConfig) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (c *TalosConfig) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
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
