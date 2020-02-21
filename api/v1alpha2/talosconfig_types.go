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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ConfigFinalizer allows us to clean up resources before deletion
	ConfigFinalizer = "talosconfig.bootstrap.cluster.x-k8s.io"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TalosConfigSpec defines the desired state of TalosConfig
type TalosConfigSpec struct {
	GenerateType string `json:"generateType"` //none,init,controlplane,worker mutually exclusive w/ data
	Data         string `json:"data,omitempty"`
	// Important: Run "make" to regenerate code after modifying this file
}

// TalosConfigStatus defines the observed state of TalosConfig
type TalosConfigStatus struct {
	// Ready indicates the BootstrapData field is ready to be consumed
	Ready bool `json:"ready,omitempty"`

	// BootstrapData will be a slice of bootstrap data
	// +optional
	BootstrapData []byte `json:"bootstrapData,omitempty"`

	// Talos config will be a string containing the config for download
	// +optional
	TalosConfig string `json:"talosConfig,omitempty"`

	// ErrorReason will be set on non-retryable errors
	// +optional
	ErrorReason string `json:"errorReason,omitempty"`

	// ErrorMessage will be set on non-retryable errors
	// +optional
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=talosconfigs,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status

// TalosConfig is the Schema for the talosconfigs API
type TalosConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TalosConfigSpec   `json:"spec,omitempty"`
	Status TalosConfigStatus `json:"status,omitempty"`
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

// TODO(rsmitty): this is disgusting, we should figure out how to do deepcopy
// and use the already existing talos pkg

// Device represents a network interface.
type Device struct {
	Interface string `json:"interface"`
	Ignore    bool   `json:"ignore"`
}
