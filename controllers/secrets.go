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

package controllers

import (
	"context"

	bootstrapv1alpha2 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha2"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *TalosConfigReconciler) fetchInputSecret(ctx context.Context, config *bootstrapv1alpha2.TalosConfig, clusterName string) (*corev1.Secret, error) {

	inputSecret := &corev1.Secret{}
	err := r.Client.Get(context.Background(), client.ObjectKey{
		Namespace: config.GetNamespace(),
		Name:      clusterName,
	}, inputSecret)

	if err != nil {
		return nil, err
	}

	return inputSecret, nil
}

func (r *TalosConfigReconciler) writeInputSecret(ctx context.Context, config *bootstrapv1alpha2.TalosConfig, clusterName string, input *generate.Input) (*corev1.Secret, error) {

	certMarshal, err := yaml.Marshal(input.Certs)
	if err != nil {
		return nil, err
	}

	kubeTokenMarshal, err := yaml.Marshal(input.KubeadmTokens)
	if err != nil {
		return nil, err
	}

	trustdInfoMarshal, err := yaml.Marshal(input.TrustdInfo)
	if err != nil {
		return nil, err
	}

	certSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: config.GetNamespace(),
			Name:      clusterName,
		},
		Data: map[string][]byte{
			"certs":      certMarshal,
			"kubeTokens": kubeTokenMarshal,
			"trustdInfo": trustdInfoMarshal,
		},
	}

	err = r.Client.Create(ctx, certSecret)
	if err != nil {
		return nil, err
	}
	return certSecret, nil
}
