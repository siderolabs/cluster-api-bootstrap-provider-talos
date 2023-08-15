// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"context"
	"testing"

	bootstrapv1alpha3 "github.com/siderolabs/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	"github.com/siderolabs/go-pointer"
	talosclientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	machineconfig "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// assertClientConfig checks that Talos client config as part of TalosConfig resource is valid.
func assertClientConfig(t *testing.T, talosConfig *bootstrapv1alpha3.TalosConfig) {
	t.Helper()

	clientConfig, err := talosclientconfig.FromString(talosConfig.Status.TalosConfig) //nolint:staticcheck
	require.NoError(t, err)
	validateClientConfig(t, clientConfig)
}

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
func assertMachineConfiguration(ctx context.Context, t *testing.T, c client.Client, talosConfig *bootstrapv1alpha3.TalosConfig) machineconfig.Provider {
	var bootstrapDataSecret corev1.Secret

	key := types.NamespacedName{
		Namespace: talosConfig.Namespace,
		Name:      pointer.SafeDeref(talosConfig.Status.DataSecretName),
	}
	require.NoError(t, c.Get(ctx, key, &bootstrapDataSecret))

	assert.Len(t, bootstrapDataSecret.Data, 1)

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

	assert.Equal(t, provider.Cluster().CA().Crt, caSecret.Data[corev1.TLSCertKey])
	assert.Equal(t, provider.Cluster().CA().Key, caSecret.Data[corev1.TLSPrivateKeyKey])
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
	secretsBundle := generate.NewSecretsBundleFromConfig(generate.NewClock(), provider)
	secretsBundle.Clock = nil

	var savedBundle generate.SecretsBundle
	require.NoError(t, yaml.Unmarshal(talosSecret.Data["bundle"], &savedBundle))
	assert.Equal(t, *secretsBundle, savedBundle)
}

// assertSameMachineConfigSecrets checks that control plane configs share same set of secrets.
func assertSameMachineConfigSecrets(ctx context.Context, t *testing.T, c client.Client, talosConfigs ...*bootstrapv1alpha3.TalosConfig) {
	providers := make([]machineconfig.Provider, len(talosConfigs))

	for i := range providers {
		providers[i] = assertMachineConfiguration(ctx, t, c, talosConfigs[i])
	}

	secretsBundle0 := generate.NewSecretsBundleFromConfig(generate.NewClock(), providers[0])

	for _, provider := range providers[1:] {
		assert.Equal(t, secretsBundle0, generate.NewSecretsBundleFromConfig(generate.NewClock(), provider))
	}
}

// assertCompatibleMachineConfigs checks that configs share same set of core secrets so that nodes can build a cluster.
func assertCompatibleMachineConfigs(ctx context.Context, t *testing.T, c client.Client, talosConfigs ...*bootstrapv1alpha3.TalosConfig) {
	providers := make([]machineconfig.Provider, len(talosConfigs))

	for i := range providers {
		providers[i] = assertMachineConfiguration(ctx, t, c, talosConfigs[i])
	}

	checks := []func(p machineconfig.Provider) interface{}{
		func(p machineconfig.Provider) interface{} { return p.Machine().Security().Token() },
		func(p machineconfig.Provider) interface{} { return p.Machine().Security().CA().Crt },
		func(p machineconfig.Provider) interface{} { return p.Cluster().ID() },
		func(p machineconfig.Provider) interface{} { return p.Cluster().Secret() },
		func(p machineconfig.Provider) interface{} { return p.Cluster().Endpoint().String() },
		func(p machineconfig.Provider) interface{} { return p.Cluster().Token().ID() },
		func(p machineconfig.Provider) interface{} { return p.Cluster().Token().Secret() },
		func(p machineconfig.Provider) interface{} { return p.Cluster().CA().Crt },
	}

	for _, check := range checks {
		value0 := check(providers[0])

		for _, provider := range providers[1:] {
			assert.Equal(t, value0, check(provider))
		}
	}
}
