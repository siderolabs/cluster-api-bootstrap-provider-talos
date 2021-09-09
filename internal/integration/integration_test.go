// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talos-systems/talos/pkg/machinery/client"
	clientconfig "github.com/talos-systems/talos/pkg/machinery/client/config"
	"github.com/talos-systems/talos/pkg/machinery/config"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
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

		// wait for TalosConfig to be reconciled
		for ctx.Err() == nil {
			key := types.NamespacedName{
				Namespace: namespaceName,
				Name:      talosConfig.Name,
			}

			err := c.Get(ctx, key, talosConfig)
			require.NoError(t, err)

			if talosConfig.Status.Ready {
				break
			}

			t.Log("Waiting ...")
			sleepCtx(ctx, 5*time.Second)
		}

		assert.Equal(t, machine.Name+"-bootstrap-data", pointer.GetString(talosConfig.Status.DataSecretName), "%+v", talosConfig)

		clientConfig, err := clientconfig.FromString(talosConfig.Status.TalosConfig)
		require.NoError(t, err)
		assert.Len(t, clientConfig.Contexts, 1)
		assert.NotEmpty(t, clientConfig.Context)
		context := clientConfig.Contexts[clientConfig.Context]
		require.NotNil(t, context)

		assert.Empty(t, context.Endpoints)
		assert.Empty(t, context.Nodes)
		creds, err := client.CredentialsFromConfigContext(context)
		require.NoError(t, err)
		assert.NotEmpty(t, creds.CA)

		var caSecret corev1.Secret
		key := types.NamespacedName{
			Namespace: namespaceName,
			Name:      cluster.Name + "-ca",
		}
		require.NoError(t, c.Get(ctx, key, &caSecret))
		assert.Len(t, caSecret.Data, 2)
		assert.Equal(t, corev1.SecretTypeOpaque, caSecret.Type)                     // TODO why not SecretTypeTLS?
		assert.NotEmpty(t, creds.Crt.Certificate, caSecret.Data[corev1.TLSCertKey]) // TODO decode and load
		assert.NotEmpty(t, caSecret.Data[corev1.TLSPrivateKeyKey])

		var talosSecret corev1.Secret
		key = types.NamespacedName{
			Namespace: namespaceName,
			Name:      cluster.Name + "-talos",
		}
		require.NoError(t, c.Get(ctx, key, &talosSecret))
		assert.Len(t, talosSecret.Data, 3)
		assert.NotEmpty(t, talosSecret.Data["certs"]) // TODO more tests
		assert.NotEmpty(t, talosSecret.Data["kubeSecrets"])
		assert.NotEmpty(t, talosSecret.Data["trustdInfo"])

		var bootstrapDataSecret corev1.Secret
		key = types.NamespacedName{
			Namespace: namespaceName,
			Name:      machine.Name + "-bootstrap-data",
		}
		require.NoError(t, c.Get(ctx, key, &bootstrapDataSecret))
		assert.Len(t, bootstrapDataSecret.Data, 1)
		provider, err := configloader.NewFromBytes(bootstrapDataSecret.Data["value"])
		require.NoError(t, err)

		provider.(*v1alpha1.Config).ClusterConfig.ControlPlane.Endpoint.Host = "FIXME"

		// TODO more tests
		_, err = provider.Validate(runtimeMode{false}, config.WithStrict())
		require.NoError(t, err)
	})
}
