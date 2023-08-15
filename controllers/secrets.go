// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-logr/logr"
	"github.com/siderolabs/crypto/x509"
	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
	talosmachine "github.com/siderolabs/talos/pkg/machinery/config/machine"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/collections"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bootstrapv1alpha3 "github.com/siderolabs/cluster-api-bootstrap-provider-talos/api/v1alpha3"
)

func (r *TalosConfigReconciler) fetchSecret(ctx context.Context, config *bootstrapv1alpha3.TalosConfig, secretName string) (*corev1.Secret, error) {
	retSecret := &corev1.Secret{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Namespace: config.GetNamespace(),
		Name:      secretName,
	}, retSecret)

	if err != nil {
		return nil, err
	}

	return retSecret, nil
}

// getSecretsBundle either generates or loads existing secret.
func (r *TalosConfigReconciler) getSecretsBundle(ctx context.Context, scope *TalosConfigScope, allowGenerate bool, versionContract *config.VersionContract) (*secrets.Bundle, error) {
	var secretsBundle *secrets.Bundle

	secretName := scope.Cluster.Name + "-talos"

retry:
	secret, err := r.fetchSecret(ctx, scope.Config, secretName)

	switch {
	case err != nil && k8serrors.IsNotFound(err):
		if !allowGenerate {
			return nil, fmt.Errorf("secrets bundle is missing")
		}

		// no cluster secret yet, generate new one
		secretsBundle, err = secrets.NewBundle(secrets.NewFixedClock(time.Now()), versionContract)
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
		secretsBundle = &secrets.Bundle{
			Clock: secrets.NewFixedClock(time.Now()),
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
			secretsBundle.Cluster = &secrets.Cluster{}
		}
	}

	return secretsBundle, nil
}

func (r *TalosConfigReconciler) writeSecretsBundleSecret(ctx context.Context, scope *TalosConfigScope, secretName string, secretsBundle *secrets.Bundle) error {
	bundle, err := yaml.Marshal(secretsBundle)
	if err != nil {
		return fmt.Errorf("error marshaling secrets bundle: %w", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: scope.Config.Namespace,
			Name:      secretName,
			Labels: map[string]string{
				capiv1.ClusterNameLabel: scope.Cluster.Name,
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
				capiv1.ClusterNameLabel: scope.Cluster.Name,
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
				capiv1.ClusterNameLabel: scope.Cluster.Name,
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

// reconcileClientConfig creates/updates a TalosConfig for the cluster.
func (r *TalosConfigReconciler) reconcileClientConfig(ctx context.Context, log logr.Logger, scope *TalosConfigScope) error {
	if !(scope.Config.Spec.GenerateType == talosmachine.TypeControlPlane.String() || scope.Config.Spec.GenerateType == talosmachine.TypeInit.String()) {
		// can only reconcile for control plane machines
		return nil
	}

	machines, err := collections.GetFilteredMachinesForCluster(ctx, r.Client, scope.Cluster, collections.ControlPlaneMachines(scope.Cluster.Name))
	if err != nil {
		return fmt.Errorf("failed getting control plane machines: %w", err)
	}

	var endpoints []string

	for _, machine := range machines {
		for _, addr := range machine.Status.Addresses {
			if addr.Type == capiv1.MachineExternalIP || addr.Type == capiv1.MachineInternalIP {
				endpoints = append(endpoints, addr.Address)
			}
		}
	}

	sort.Strings(endpoints)

	secretBundle, err := r.getSecretsBundle(ctx, scope, false, defaultVersionContract) // version contract doesn't matter, as we're getting the secrets
	if err != nil {
		return err
	}

	talosConfig, err := genTalosConfigFile(scope.Cluster.Name, secretBundle, endpoints)
	if err != nil {
		return err
	}

	// Create or update talosconfig secret
	dataSecretName := scope.Cluster.GetName() + "-talosconfig"

	log.Info("updating talosconfig", "endpoints", endpoints, "secret", dataSecretName)

	configSecret := &corev1.Secret{}

	err = r.Client.Get(ctx, client.ObjectKey{Namespace: scope.Cluster.Namespace, Name: dataSecretName}, configSecret)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("error fetching secret: %w", err)
		}

		configSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: scope.Cluster.Namespace,
				Name:      dataSecretName,
				Labels: map[string]string{
					capiv1.ClusterNameLabel: scope.Cluster.Name,
				},
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(scope.Cluster, capiv1.GroupVersion.WithKind("Cluster")),
				},
			},
			Data: map[string][]byte{
				"talosconfig": []byte(talosConfig),
			},
		}

		return r.Client.Create(ctx, configSecret)
	}

	configSecret.Data["talosconfig"] = []byte(talosConfig)

	err = r.Client.Update(ctx, configSecret)
	if k8serrors.IsConflict(err) {
		// ignore conflict errors, probably another reconcile fixed up the endpoints
		err = nil
	}

	return err
}
