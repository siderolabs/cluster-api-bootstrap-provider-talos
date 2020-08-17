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
	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	"github.com/talos-systems/crypto/x509"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *TalosConfigReconciler) fetchSecret(ctx context.Context, config *bootstrapv1alpha3.TalosConfig, secretName string) (*corev1.Secret, error) {
	retSecret := &corev1.Secret{}
	err := r.Client.Get(context.Background(), client.ObjectKey{
		Namespace: config.GetNamespace(),
		Name:      secretName,
	}, retSecret)

	if err != nil {
		return nil, err
	}

	return retSecret, nil
}

func (r *TalosConfigReconciler) writeInputSecret(ctx context.Context, scope *TalosConfigScope, input *generate.Input) (*corev1.Secret, error) {

	certMarshal, err := yaml.Marshal(input.Certs)
	if err != nil {
		return nil, err
	}

	kubeSecretsMarshal, err := yaml.Marshal(input.Secrets)
	if err != nil {
		return nil, err
	}

	trustdInfoMarshal, err := yaml.Marshal(input.TrustdInfo)
	if err != nil {
		return nil, err
	}

	certSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: scope.Config.Namespace,
			Name:      scope.Cluster.Name + "-talos",
			Labels: map[string]string{
				clusterv1.ClusterLabelName: scope.Cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(scope.Config, bootstrapv1alpha3.GroupVersion.WithKind("TalosConfig")),
			},
		},
		Data: map[string][]byte{
			"certs":       certMarshal,
			"kubeSecrets": kubeSecretsMarshal,
			"trustdInfo":  trustdInfoMarshal,
		},
	}

	err = r.Client.Create(ctx, certSecret)
	if err != nil {
		return nil, err
	}
	return certSecret, nil
}

func (r *TalosConfigReconciler) writeK8sCASecret(ctx context.Context, scope *TalosConfigScope, certs *x509.PEMEncodedCertificateAndKey) error {
	// Create ca secret only if it doesn't already exist
	_, err := r.fetchSecret(ctx, scope.Config, scope.Cluster.Name+"-ca")
	if k8serrors.IsNotFound(err) {
		certSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: scope.Config.Namespace,
				Name:      scope.Cluster.Name + "-ca",
				Labels: map[string]string{
					clusterv1.ClusterLabelName: scope.Cluster.Name,
				},
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(scope.Config, bootstrapv1alpha3.GroupVersion.WithKind("TalosConfig")),
				},
			},
			Data: map[string][]byte{
				"tls.crt": certs.Crt,
				"tls.key": certs.Key,
			},
		}

		err = r.Client.Create(ctx, certSecret)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

// writeBootstrapData creates a new secret with the data passed in as input
func (r *TalosConfigReconciler) writeBootstrapData(ctx context.Context, scope *TalosConfigScope, data []byte) error {
	// Create ca secret only if it doesn't already exist
	_, err := r.fetchSecret(ctx, scope.Config, scope.Machine.Name+"-bootstrap-data")
	if k8serrors.IsNotFound(err) {
		certSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: scope.Config.Namespace,
				Name:      scope.Machine.Name + "-bootstrap-data",
				Labels: map[string]string{
					clusterv1.ClusterLabelName: scope.Cluster.Name,
				},
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(scope.Config, bootstrapv1alpha3.GroupVersion.WithKind("TalosConfig")),
				},
			},
			Data: map[string][]byte{
				"value": data,
			},
		}

		err = r.Client.Create(ctx, certSecret)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (r *TalosConfigReconciler) deleteSecret(ctx context.Context, config *bootstrapv1alpha2.TalosConfig, secretName string) error {
	err := r.Client.Delete(ctx,
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: config.GetNamespace(),
				Name:      secretName,
			},
		},
	)

	if err != nil && !k8serrors.IsNotFound(err) {
		return err
	}

	return nil
}
