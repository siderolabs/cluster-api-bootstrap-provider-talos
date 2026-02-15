// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"context"
	"testing"
	"time"

	bootstrapv1beta1 "github.com/siderolabs/cluster-api-bootstrap-provider-talos/api/v1beta1"
	talosclientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	machineconfig "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
	"github.com/siderolabs/talos/pkg/machinery/config/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// assertClusterClientConfig checks that Talos client config as a cluster-wide secret is valid.
func assertClusterClientConfig(ctx context.Context, t *testing.T, c client.Client, cluster *capiv1.Cluster, endpoints ...string) {
	t.Helper()

	var secret corev1.Secret

	require.NoError(t, c.Get(ctx, client.ObjectKey{Namespace: cluster.Namespace, Name: cluster.Name + "-talosconfig"}, &secret))

	clientConfig, err := talosclientconfig.FromString(string(secret.Data["talosconfig"]))
	require.NoError(t, err)
	validateClientConfig(t, clientConfig, endpoints...)
}

// assertMachineConfiguration checks that generated bootstrap data is a valid Talos machine configuration.
func assertMachineConfiguration(ctx context.Context, t *testing.T, c client.Client, talosConfig *bootstrapv1beta1.TalosConfig) machineconfig.Provider {
	var bootstrapDataSecret corev1.Secret

	key := types.NamespacedName{
		Namespace: talosConfig.Namespace,
		Name:      talosConfig.Status.DataSecretName,
	}
	require.NoError(t, c.Get(ctx, key, &bootstrapDataSecret))

	assert.Len(t, bootstrapDataSecret.Data, 1)

	assert.Less(t, len(bootstrapDataSecret.Data["value"]), 32*1024) // 32KB limit on user-data size

	provider, err := configloader.NewFromBytes(bootstrapDataSecret.Data["value"])
	require.NoError(t, err)

	_, err = provider.Validate(runtimeMode{false}, validation.WithStrict())
	assert.NoError(t, err)

	return provider
}

// assertClusterCA checks that generated cluster CA secret matches secrets in machine config (machine config from controlplane node required).
func assertClusterCA(ctx context.Context, t *testing.T, c client.Client, cluster *capiv1.Cluster, provider machineconfig.Provider) {
	var caSecret corev1.Secret

	key := types.NamespacedName{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-ca",
	}
	require.NoError(t, c.Get(ctx, key, &caSecret))

	assert.Len(t, caSecret.Data, 2)
	assert.Equal(t, corev1.SecretTypeOpaque, caSecret.Type) // TODO why not SecretTypeTLS?

	assert.NotEmpty(t, caSecret.Data[corev1.TLSCertKey])
	assert.NotEmpty(t, caSecret.Data[corev1.TLSPrivateKeyKey])

	assert.Equal(t, provider.Cluster().IssuingCA().Crt, caSecret.Data[corev1.TLSCertKey])
	assert.Equal(t, provider.Cluster().IssuingCA().Key, caSecret.Data[corev1.TLSPrivateKeyKey])
}

// assertControllerSecret checks that persisted controller secret (used to bootstrap more machines with same secrets) maches generated controlplane config.
func assertControllerSecret(ctx context.Context, t *testing.T, c client.Client, cluster *capiv1.Cluster, provider machineconfig.Provider) {
	var talosSecret corev1.Secret
	key := types.NamespacedName{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-talos",
	}
	require.NoError(t, c.Get(ctx, key, &talosSecret))

	assert.Len(t, talosSecret.Data, 1)
	assert.NotEmpty(t, talosSecret.Data["bundle"])

	// cross-checks
	secretsBundle := secrets.NewBundleFromConfig(secrets.NewFixedClock(time.Now()), provider)
	secretsBundle.Clock = nil

	var savedBundle secrets.Bundle
	require.NoError(t, yaml.Unmarshal(talosSecret.Data["bundle"], &savedBundle))
	assert.Equal(t, *secretsBundle, savedBundle)
}

// assertSameMachineConfigSecrets checks that control plane configs share same set of secrets.
func assertSameMachineConfigSecrets(ctx context.Context, t *testing.T, c client.Client, talosConfigs ...*bootstrapv1beta1.TalosConfig) {
	providers := make([]machineconfig.Provider, len(talosConfigs))

	for i := range providers {
		providers[i] = assertMachineConfiguration(ctx, t, c, talosConfigs[i])
	}

	clock := secrets.NewFixedClock(time.Now())

	secretsBundle0 := secrets.NewBundleFromConfig(clock, providers[0])

	for _, provider := range providers[1:] {
		assert.Equal(t, secretsBundle0, secrets.NewBundleFromConfig(clock, provider))
	}
}

// assertCompatibleMachineConfigs checks that configs share same set of core secrets so that nodes can build a cluster.
func assertCompatibleMachineConfigs(ctx context.Context, t *testing.T, c client.Client, talosConfigs ...*bootstrapv1beta1.TalosConfig) {
	providers := make([]machineconfig.Provider, len(talosConfigs))

	for i := range providers {
		providers[i] = assertMachineConfiguration(ctx, t, c, talosConfigs[i])
	}

	checks := []func(p machineconfig.Provider) any{
		func(p machineconfig.Provider) any { return p.Machine().Security().Token() },
		func(p machineconfig.Provider) any { return p.Machine().Security().IssuingCA().Crt },
		func(p machineconfig.Provider) any { return p.Cluster().ID() },
		func(p machineconfig.Provider) any { return p.Cluster().Secret() },
		func(p machineconfig.Provider) any { return p.Cluster().Endpoint().String() },
		func(p machineconfig.Provider) any { return p.Cluster().Token().ID() },
		func(p machineconfig.Provider) any { return p.Cluster().Token().Secret() },
		func(p machineconfig.Provider) any { return p.Cluster().IssuingCA().Crt },
	}

	for _, check := range checks {
		value0 := check(providers[0])

		for _, provider := range providers[1:] {
			assert.Equal(t, value0, check(provider))
		}
	}
}
