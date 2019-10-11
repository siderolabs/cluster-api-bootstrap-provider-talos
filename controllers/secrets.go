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
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *TalosConfigReconciler) fetchCertSecret(ctx context.Context, config *bootstrapv1alpha2.TalosConfig, clusterName string) (*corev1.Secret, error) {

	certSecret := &corev1.Secret{}
	err := r.Client.Get(context.Background(), client.ObjectKey{
		Namespace: config.GetNamespace(),
		Name:      clusterName,
	}, certSecret)

	if err != nil && k8serrors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return certSecret, nil
}

func (r *TalosConfigReconciler) writeCertSecret(ctx context.Context, config *bootstrapv1alpha2.TalosConfig, clusterName string, certs *generate.Certs) error {

	certMarshal, err := yaml.Marshal(certs)
	if err != nil {
		return err
	}

	certSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: config.GetNamespace(),
			Name:      clusterName,
		},
		Data: map[string][]byte{"certs": certMarshal},
	}

	err = r.Client.Create(ctx, certSecret)
	if err != nil {
		return err
	}
	return nil
}
