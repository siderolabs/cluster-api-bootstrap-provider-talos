// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	talosclientconfig "github.com/talos-systems/talos/pkg/machinery/client/config"
	machineconfig "github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestIntegration(t *testing.T) {
	ctx, c := setupSuite(t)

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName)
		machine := createMachine(ctx, t, c, cluster)
		talosConfig := createTalosConfig(ctx, t, c, machine)
		waitForReady(ctx, t, c, talosConfig)

		// check talosConfig
		{
			assert.Equal(t, machine.Name+"-bootstrap-data", pointer.GetString(talosConfig.Status.DataSecretName), "%+v", talosConfig)
			clientConfig, err := talosclientconfig.FromString(talosConfig.Status.TalosConfig)
			require.NoError(t, err)
			creds := validateClientConfig(t, clientConfig)
			talosCA := parsePEMCertificate(t, creds.CA)
			_ = talosCA
			// t.Logf("Talos CA:\n%s", spew.Sdump(talosCA))
		}

		// get <cluster>-ca secret
		var caSecret corev1.Secret
		key := types.NamespacedName{
			Namespace: namespaceName,
			Name:      cluster.Name + "-ca",
		}
		require.NoError(t, c.Get(ctx, key, &caSecret))

		// check <cluster>-ca secret
		{
			assert.Len(t, caSecret.Data, 2)
			assert.Equal(t, corev1.SecretTypeOpaque, caSecret.Type) // TODO why not SecretTypeTLS?
			assert.NotEmpty(t, caSecret.Data[corev1.TLSCertKey])
			assert.NotEmpty(t, caSecret.Data[corev1.TLSPrivateKeyKey])
			kubeCA := parsePEMCertificate(t, caSecret.Data[corev1.TLSCertKey])
			_ = kubeCA
			// t.Logf("kubeCA:\n%s", spew.Sdump(kubeCA))
		}

		// get <cluster>-talos secret
		var talosSecret corev1.Secret
		key = types.NamespacedName{
			Namespace: namespaceName,
			Name:      cluster.Name + "-talos",
		}
		require.NoError(t, c.Get(ctx, key, &talosSecret))

		// check <cluster>-talos secret
		{
			assert.Len(t, talosSecret.Data, 3)
			assert.NotEmpty(t, talosSecret.Data["certs"])
			assert.NotEmpty(t, talosSecret.Data["kubeSecrets"])
			assert.NotEmpty(t, talosSecret.Data["trustdInfo"])
		}

		// get <machine>-bootstrap-data secret
		var bootstrapDataSecret corev1.Secret
		key = types.NamespacedName{
			Namespace: namespaceName,
			Name:      machine.Name + "-bootstrap-data",
		}
		require.NoError(t, c.Get(ctx, key, &bootstrapDataSecret))

		// check <machine>-bootstrap-data secret
		var provider machineconfig.Provider
		{
			assert.Len(t, bootstrapDataSecret.Data, 1)
			var err error
			provider, err = configloader.NewFromBytes(bootstrapDataSecret.Data["value"])
			require.NoError(t, err)
			_, err = provider.Validate(runtimeMode{false}, machineconfig.WithStrict())
			require.NoError(t, err)
		}

		// cross-checks
		{
			secretsBundle := generate.NewSecretsBundleFromConfig(generate.NewClock(), provider)

			var certs generate.Certs
			require.NoError(t, yaml.Unmarshal(talosSecret.Data["certs"], &certs))
			assert.NotEmpty(t, certs.Admin)
			certs.Admin = nil
			assert.Equal(t, secretsBundle.Certs, &certs)
			assert.Equal(t, caSecret.Data[corev1.TLSCertKey], certs.K8s.Crt)

			var kubeSecrets generate.Secrets
			require.NoError(t, yaml.Unmarshal(talosSecret.Data["kubeSecrets"], &kubeSecrets))
			assert.Equal(t, secretsBundle.Secrets, &kubeSecrets)

			var trustdInfo generate.TrustdInfo
			require.NoError(t, yaml.Unmarshal(talosSecret.Data["trustdInfo"], &trustdInfo))
			assert.Equal(t, secretsBundle.TrustdInfo, &trustdInfo)
		}

		// create the second machine
		machine2 := createMachine(ctx, t, c, cluster)
		talosConfig2 := createTalosConfig(ctx, t, c, machine2)
		waitForReady(ctx, t, c, talosConfig2)

		// get <machine>-bootstrap-data secret
		var bootstrapDataSecret2 corev1.Secret
		key = types.NamespacedName{
			Namespace: namespaceName,
			Name:      machine2.Name + "-bootstrap-data",
		}
		require.NoError(t, c.Get(ctx, key, &bootstrapDataSecret2))

		assert.Equal(t, bootstrapDataSecret.Data, bootstrapDataSecret2.Data) // ?!
	})
}
