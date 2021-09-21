// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"

	"github.com/talos-systems/crypto/x509"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
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

// getSecretsBundle either generates or loads existing secret.
func (r *TalosConfigReconciler) getSecretsBundle(ctx context.Context, scope *TalosConfigScope, secretName string, opts ...generate.GenOption) (*generate.SecretsBundle, error) {
	var secretsBundle *generate.SecretsBundle

retry:
	secret, err := r.fetchSecret(ctx, scope.Config, secretName)

	switch {
	case err != nil && k8serrors.IsNotFound(err):
		// no cluster secret yet, generate new one
		secretsBundle, err = generate.NewSecretsBundle(generate.NewClock(), opts...)
		if err != nil {
			return nil, fmt.Errorf("error generating new secrets bundle: %w", err)
		}

		if err = r.writeSecretsBundleSecret(ctx, scope, secretName, secretsBundle); err != nil {
			if k8serrors.IsAlreadyExists(err) {
				// conflict on creation, retry loading
				goto retry
			}

			return nil, fmt.Errorf("error writing secrets bundle: %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("error reading secrets bundle: %w", err)
	default:
		// successfully loaded secret, initialize secretsBundle from it
		secretsBundle = &generate.SecretsBundle{
			Clock: generate.NewClock(),
		}

		if _, ok := secret.Data["bundle"]; ok {
			// new format
			if err = yaml.Unmarshal(secret.Data["bundle"], secretsBundle); err != nil {
				return nil, fmt.Errorf("error unmarshaling secrets bundle: %w", err)
			}
		} else {
			// legacy format
			if err = yaml.Unmarshal(secret.Data["certs"], &secretsBundle.Certs); err != nil {
				return nil, fmt.Errorf("error unmarshaling certs: %w", err)
			}

			if err = yaml.Unmarshal(secret.Data["kubeSecrets"], &secretsBundle.Secrets); err != nil {
				return nil, fmt.Errorf("error unmarshaling secrets: %w", err)
			}

			if err = yaml.Unmarshal(secret.Data["trustdInfo"], &secretsBundle.TrustdInfo); err != nil {
				return nil, fmt.Errorf("error unmarshaling trustd info: %w", err)
			}

			// not stored in legacy format, use empty values
			secretsBundle.Cluster = &generate.Cluster{}
		}
	}

	return secretsBundle, nil
}

func (r *TalosConfigReconciler) writeSecretsBundleSecret(ctx context.Context, scope *TalosConfigScope, secretName string, secretsBundle *generate.SecretsBundle) error {
	bundle, err := yaml.Marshal(secretsBundle)
	if err != nil {
		return fmt.Errorf("error marshaling secrets bundle: %w", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: scope.Config.Namespace,
			Name:      secretName,
			Labels: map[string]string{
				capiv1.ClusterLabelName: scope.Cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(scope.Cluster, capiv1.GroupVersion.WithKind("Cluster")),
			},
		},
		Data: map[string][]byte{
			"bundle": bundle,
		},
	}

	return r.Client.Create(ctx, secret)
}

func (r *TalosConfigReconciler) writeK8sCASecret(ctx context.Context, scope *TalosConfigScope, certs *x509.PEMEncodedCertificateAndKey) error {
	// Create ca secret only if it doesn't already exist
	_, err := r.fetchSecret(ctx, scope.Config, scope.Cluster.Name+"-ca")
	if err == nil {
		return nil
	}

	if !k8serrors.IsNotFound(err) {
		return err
	}

	certSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: scope.Config.Namespace,
			Name:      scope.Cluster.Name + "-ca",
			Labels: map[string]string{
				capiv1.ClusterLabelName: scope.Cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(scope.Cluster, capiv1.GroupVersion.WithKind("Cluster")),
			},
		},
		Data: map[string][]byte{
			"tls.crt": certs.Crt,
			"tls.key": certs.Key,
		},
	}

	err = r.Client.Create(ctx, certSecret)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// writeBootstrapData creates a new secret with the data passed in as input
func (r *TalosConfigReconciler) writeBootstrapData(ctx context.Context, scope *TalosConfigScope, data []byte) (string, error) {
	// Create bootstrap secret only if it doesn't already exist
	ownerName := scope.ConfigOwner.GetName()
	dataSecretName := ownerName + "-bootstrap-data"

	r.Log.Info("handling bootstrap data for ", "owner", ownerName)

	_, err := r.fetchSecret(ctx, scope.Config, dataSecretName)
	if err == nil {
		return dataSecretName, nil
	}

	if err != nil && !k8serrors.IsNotFound(err) {
		return dataSecretName, err
	}

	certSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: scope.Config.Namespace,
			Name:      dataSecretName,
			Labels: map[string]string{
				capiv1.ClusterLabelName: scope.Cluster.Name,
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

	return dataSecretName, err
}
