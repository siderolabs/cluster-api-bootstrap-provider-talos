/*
Copyright 2019 The Kubernetes Authors.

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

// Modified by Pip Oomen: trimmed copy of sigs.k8s.io/cluster-api/bootstrap/util,
// which was moved into the kubeadm bootstrap provider in Cluster API v1.14.

package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigOwner provides a data interface for different config owner types.
type ConfigOwner struct {
	*unstructured.Unstructured
}

// ClusterName extracts spec.clusterName from the config owner.
func (co ConfigOwner) ClusterName() string {
	clusterName, _, err := unstructured.NestedString(co.Object, "spec", "clusterName")
	if err != nil {
		return ""
	}

	return clusterName
}

// DataSecretName extracts spec.bootstrap.dataSecretName from the config owner.
func (co ConfigOwner) DataSecretName() *string {
	dataSecretName, exist, err := unstructured.NestedString(co.Object, "spec", "bootstrap", "dataSecretName")
	if err != nil || !exist {
		return nil
	}

	return &dataSecretName
}

// IsMachinePool checks if the config owner is a MachinePool.
func (co ConfigOwner) IsMachinePool() bool {
	return co.GetKind() == "MachinePool"
}

// GetConfigOwner returns the Machine or MachinePool owning obj, or nil if there is none.
func GetConfigOwner(ctx context.Context, c client.Client, obj metav1.Object) (*ConfigOwner, error) {
	allowedGKs := []schema.GroupKind{
		{
			Group: clusterv1.GroupVersion.Group,
			Kind:  "Machine",
		},
	}

	if feature.Gates.Enabled(feature.MachinePool) {
		allowedGKs = append(allowedGKs, schema.GroupKind{
			Group: clusterv1.GroupVersion.Group,
			Kind:  "MachinePool",
		})
	}

	for _, ref := range obj.GetOwnerReferences() {
		refGV, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GroupVersion from %q: %w", ref.APIVersion, err)
		}

		for _, gk := range allowedGKs {
			if refGV.Group != gk.Group || ref.Kind != gk.Kind {
				continue
			}

			// Intentionally always using the latest apiVersion to get Machine or MachinePool,
			// even if the apiVersion in the ownerRef has not been bumped yet.
			owner, err := external.Get(ctx, c, &corev1.ObjectReference{
				APIVersion: clusterv1.GroupVersion.String(),
				Kind:       ref.Kind,
				Name:       ref.Name,
				Namespace:  obj.GetNamespace(),
			})
			if err != nil {
				return nil, err
			}

			return &ConfigOwner{owner}, nil
		}
	}

	return nil, nil
}
