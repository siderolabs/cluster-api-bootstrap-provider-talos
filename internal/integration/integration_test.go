// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	talosmachine "github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestIntegration(t *testing.T) {
	require.NotEmpty(t, TalosVersion)

	ctx, c := setupSuite(t)

	t.Run("SingleNode", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil)
		machine := createMachine(ctx, t, c, cluster)
		talosConfig := createTalosConfig(ctx, t, c, machine, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
		})
		waitForReady(ctx, t, c, talosConfig)

		assertClientConfig(t, talosConfig)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, talosmachine.TypeInit, provider.Machine().Type())

		assertClusterCA(ctx, t, c, cluster, provider)

		assertControllerSecret(ctx, t, c, cluster, provider)
	})

	t.Run("Cluster", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil)

		controlplanes := []*bootstrapv1alpha3.TalosConfig{}

		for i := 0; i < 3; i++ {
			machine := createMachine(ctx, t, c, cluster)

			machineType := talosmachine.TypeInit

			if i > 0 {
				machineType = talosmachine.TypeControlPlane
			}

			controlplanes = append(controlplanes, createTalosConfig(ctx, t, c, machine, bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: machineType.String(),
				TalosVersion: TalosVersion,
			}))
		}

		workers := []*bootstrapv1alpha3.TalosConfig{}

		for i := 0; i < 4; i++ {
			machine := createMachine(ctx, t, c, cluster)

			workers = append(workers, createTalosConfig(ctx, t, c, machine, bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: talosmachine.TypeJoin.String(),
				TalosVersion: TalosVersion,
			}))
		}

		for i, talosConfig := range append(append([]*bootstrapv1alpha3.TalosConfig{}, controlplanes...), workers...) {
			waitForReady(ctx, t, c, talosConfig)

			assertClientConfig(t, talosConfig)

			provider := assertMachineConfiguration(ctx, t, c, talosConfig)

			switch {
			case i == 0:
				assert.Equal(t, talosmachine.TypeInit, provider.Machine().Type())
			case i < len(controlplanes):
				assert.Equal(t, talosmachine.TypeControlPlane, provider.Machine().Type())
			default:
				assert.Equal(t, talosmachine.TypeJoin, provider.Machine().Type())
			}
		}

		assertClusterCA(ctx, t, c, cluster, assertMachineConfiguration(ctx, t, c, controlplanes[0]))
		assertControllerSecret(ctx, t, c, cluster, assertMachineConfiguration(ctx, t, c, controlplanes[0]))

		// compare control plane secrets completely
		assertSameMachineConfigSecrets(ctx, t, c, controlplanes...)

		// compare all configs in more relaxed mode
		assertCompatibleMachineConfigs(ctx, t, c, append(append([]*bootstrapv1alpha3.TalosConfig{}, controlplanes...), workers...)...)
	})

	t.Run("ClusterSpec", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, &capiv1.ClusterSpec{
			ClusterNetwork: &capiv1.ClusterNetwork{
				Services: &capiv1.NetworkRanges{
					CIDRBlocks: []string{
						"192.168.0.0/16",
						"fdaa:bbbb:cccc:15::/64",
					},
				},
				Pods: &capiv1.NetworkRanges{
					CIDRBlocks: []string{
						"10.0.0.0/16",
						"fdbb:bbbb:cccc:15::/64",
					},
				},
				ServiceDomain: "mycluster.local",
			},
			ControlPlaneEndpoint: capiv1.APIEndpoint{
				Host: "example.com",
				Port: 443,
			},
		})
		machine := createMachine(ctx, t, c, cluster)
		talosConfig := createTalosConfig(ctx, t, c, machine, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
			TalosVersion: TalosVersion,
		})
		waitForReady(ctx, t, c, talosConfig)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, "https://example.com:443", provider.Cluster().Endpoint().String())
		assert.Equal(t, "mycluster.local", provider.Cluster().Network().DNSDomain())
		assert.Equal(t, "10.0.0.0/16,fdbb:bbbb:cccc:15::/64", provider.Cluster().Network().PodCIDR())
		assert.Equal(t, "192.168.0.0/16,fdaa:bbbb:cccc:15::/64", provider.Cluster().Network().ServiceCIDR())
	})

	t.Run("ConfigPatches", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil)
		machine := createMachine(ctx, t, c, cluster)
		talosConfig := createTalosConfig(ctx, t, c, machine, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
			TalosVersion: TalosVersion,
			ConfigPatches: []bootstrapv1alpha3.ConfigPatches{
				{
					Op:   "add",
					Path: "/machine/time",
					Value: apiextensions.JSON{
						Raw: []byte(`{"servers": ["time.cloudflare.com"]}`),
					},
				},
				{
					Op:   "replace",
					Path: "/machine/certSANs",
					Value: apiextensions.JSON{
						Raw: []byte(`["myserver.com"]`),
					},
				},
			},
		})
		waitForReady(ctx, t, c, talosConfig)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, []string{"time.cloudflare.com"}, provider.Machine().Time().Servers())
		assert.Equal(t, []string{"myserver.com"}, provider.Machine().Security().CertSANs())
	})

}
