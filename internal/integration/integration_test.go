// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	talosmachine "github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	"inet.af/netaddr"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

func TestIntegration(t *testing.T) {
	require.NotEmpty(t, TalosVersion)

	ctx, c := setupSuite(t)

	t.Run("SingleNode", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, true)
		talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
		})
		createMachine(ctx, t, c, cluster, talosConfig, true)
		waitForReady(ctx, t, c, talosConfig)

		assertClientConfig(t, talosConfig)
		assertClusterClientConfig(ctx, t, c, cluster)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, talosmachine.TypeInit, provider.Machine().Type())

		assertClusterCA(ctx, t, c, cluster, provider)

		assertControllerSecret(ctx, t, c, cluster, provider)
	})

	t.Run("Cluster", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, true)

		controlplanes := []*bootstrapv1alpha3.TalosConfig{}
		controlplaneMachines := []*capiv1.Machine{}

		for i := 0; i < 3; i++ {
			machineType := talosmachine.TypeInit

			if i > 0 {
				machineType = talosmachine.TypeControlPlane
			}

			talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: machineType.String(),
				TalosVersion: TalosVersion,
			})
			controlplaneMachines = append(controlplaneMachines, createMachine(ctx, t, c, cluster, talosConfig, true))

			controlplanes = append(controlplanes, talosConfig)
		}

		workers := []*bootstrapv1alpha3.TalosConfig{}
		workerMachines := []*capiv1.Machine{}

		for i := 0; i < 4; i++ {
			talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: talosmachine.TypeWorker.String(),
				TalosVersion: TalosVersion,
			})
			workerMachines = append(workerMachines, createMachine(ctx, t, c, cluster, talosConfig, false))

			workers = append(workers, talosConfig)
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
				assert.Equal(t, talosmachine.TypeWorker, provider.Machine().Type())
			}
		}

		assertClusterClientConfig(ctx, t, c, cluster)
		assertClusterCA(ctx, t, c, cluster, assertMachineConfiguration(ctx, t, c, controlplanes[0]))
		assertControllerSecret(ctx, t, c, cluster, assertMachineConfiguration(ctx, t, c, controlplanes[0]))

		// compare control plane secrets completely
		assertSameMachineConfigSecrets(ctx, t, c, controlplanes...)

		// compare all configs in more relaxed mode
		assertCompatibleMachineConfigs(ctx, t, c, append(append([]*bootstrapv1alpha3.TalosConfig{}, controlplanes...), workers...)...)

		// attach addresses to machines
		ip := netaddr.MustParseIP("10.5.0.2")
		expectedEndpoints := []string{}

		for _, cpMachine := range controlplaneMachines {
			expectedEndpoints = append(expectedEndpoints, ip.String())
			patchMachineAddress(ctx, t, c, cpMachine, ip.String())

			ip = ip.Next()
		}

		for _, wMachine := range workerMachines {
			patchMachineAddress(ctx, t, c, wMachine, ip.String())

			ip = ip.Next()
		}

		waitForEndpointsClusterClientConfig(ctx, t, c, cluster, len(expectedEndpoints))

		assertClusterClientConfig(ctx, t, c, cluster, expectedEndpoints...)
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
		}, true)
		talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
			TalosVersion: TalosVersion,
		})
		createMachine(ctx, t, c, cluster, talosConfig, true)
		waitForReady(ctx, t, c, talosConfig)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, "https://example.com:443", provider.Cluster().Endpoint().String())
		assert.Equal(t, "mycluster.local", provider.Cluster().Network().DNSDomain())
		assert.Equal(t, "10.0.0.0/16,fdbb:bbbb:cccc:15::/64", strings.Join(provider.Cluster().Network().PodCIDRs(), ","))
		assert.Equal(t, "192.168.0.0/16,fdaa:bbbb:cccc:15::/64", strings.Join(provider.Cluster().Network().ServiceCIDRs(), ","))
	})

	t.Run("ConfigPatches", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, true)
		talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
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
		createMachine(ctx, t, c, cluster, talosConfig, true)
		waitForReady(ctx, t, c, talosConfig)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, []string{"time.cloudflare.com"}, provider.Machine().Time().Servers())
		assert.Equal(t, []string{"myserver.com"}, provider.Machine().Security().CertSANs())
	})

	t.Run("LegacyClusterSecret", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, true)

		// create a secret which imitates legacy secret format.
		clusterSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespaceName,
				Name:      cluster.Name + "-talos",
				Labels: map[string]string{
					capiv1.ClusterLabelName: cluster.Name,
				},
			},
		}

		require.NoError(t, json.Unmarshal([]byte(legacySecretData), &clusterSecret.Data))
		require.NoError(t, c.Create(ctx, &clusterSecret))

		talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeControlPlane.String(),
			TalosVersion: TalosVersion,
		})
		createMachine(ctx, t, c, cluster, talosConfig, true)
		waitForReady(ctx, t, c, talosConfig)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assertClusterCA(ctx, t, c, cluster, provider)
		assertClusterClientConfig(ctx, t, c, cluster)

		assert.Equal(t, "o19zh7.yv7rxce3lsptnme9", provider.Machine().Security().Token())
		assert.Equal(t, "5dwzrh", provider.Cluster().Token().ID())
		assert.Equal(t, "5ms9d5eke1muskrg", provider.Cluster().Token().Secret())
		assert.Equal(t, "-----BEGIN CERTIFICATE-----\nMIIBiTCCAS+gAwIBAgIQM4a04RExgV7BBZ2qmazx3TAKBggqhkjOPQQDBDAVMRMw\nEQYDVQQKEwprdWJlcm5ldGVzMB4XDTIxMDkyMDE4NDE0OVoXDTMxMDkxODE4NDE0\nOVowFTETMBEGA1UEChMKa3ViZXJuZXRlczBZMBMGByqGSM49AgEGCCqGSM49AwEH\nA0IABLezryg3QXmplOVP7+ap/ZTQCSlL3qiOeV7m3G8w8rvRaf+La9D0fCVJ9Rj/\nTyuuQFxQ203oeXPIfmE9HqtdjwqjYTBfMA4GA1UdDwEB/wQEAwIChDAdBgNVHSUE\nFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4E\nFgQUW0vg9AdP/ZK5+yR/73BpfvPRHMkwCgYIKoZIzj0EAwQDSAAwRQIgdvTMbjH+\n4XOMZzFIDjnq42I/suDw4cnGXcrlWdJ+aZYCIQDurrEAKmPrMgNqT2wP6JWYylla\n3l7yV8hS5CgCpJTaEg==\n-----END CERTIFICATE-----\n", string(provider.Cluster().CA().Crt)) //nolint:lll
		assert.Equal(t, "-----BEGIN CERTIFICATE-----\nMIIBPzCB8qADAgECAhEArv8iYjWXC8Mataa8e2pezDAFBgMrZXAwEDEOMAwGA1UE\nChMFdGFsb3MwHhcNMjEwOTIwMTg0MTQ5WhcNMzEwOTE4MTg0MTQ5WjAQMQ4wDAYD\nVQQKEwV0YWxvczAqMAUGAytlcAMhAOCRMlGNjsdQmgls2PCSgMdMeAIB8fAKsnCp\naXX3rfUKo2EwXzAOBgNVHQ8BAf8EBAMCAoQwHQYDVR0lBBYwFAYIKwYBBQUHAwEG\nCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFIDgT1HeMDtWHHXl\nmVhYqUPDU0JoMAUGAytlcANBAD2GLO2vG9MHGxt9658X4xZLSYNldAgDy2tHmZ7l\nnAjAR0npZoQXBVhorrQEcea7g6To9BDmtzrF0StW895d0Ak=\n-----END CERTIFICATE-----\n", string(provider.Machine().Security().CA().Crt))
	})

	t.Run("ConfigTypeNone", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, true)

		secretsBundle, err := generate.NewSecretsBundle(generate.NewClock())
		require.NoError(t, err)

		input, err := generate.NewInput(cluster.Name, "https://example.com:6443/", "v1.22.2", secretsBundle)
		require.NoError(t, err)

		workers := []*bootstrapv1alpha3.TalosConfig{}

		for i := 0; i < 4; i++ {
			machineconfig, err := generate.Config(talosmachine.TypeWorker, input)
			require.NoError(t, err)

			configdata, err := machineconfig.Bytes()
			require.NoError(t, err)

			talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: "none",
				Data:         string(configdata),
			})
			createMachine(ctx, t, c, cluster, talosConfig, true)

			workers = append(workers, talosConfig)
		}

		controlplanes := []*bootstrapv1alpha3.TalosConfig{}

		for i := 0; i < 3; i++ {
			machineType := talosmachine.TypeInit

			if i > 0 {
				machineType = talosmachine.TypeControlPlane
			}

			machineconfig, err := generate.Config(machineType, input)
			require.NoError(t, err)

			configdata, err := machineconfig.Bytes()
			require.NoError(t, err)

			talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: "none",
				Data:         string(configdata),
			})
			createMachine(ctx, t, c, cluster, talosConfig, false)

			controlplanes = append(controlplanes, talosConfig)
		}

		for i, talosConfig := range append(append([]*bootstrapv1alpha3.TalosConfig{}, controlplanes...), workers...) {
			waitForReady(ctx, t, c, talosConfig)

			provider := assertMachineConfiguration(ctx, t, c, talosConfig)

			switch {
			case i == 0:
				assert.Equal(t, talosmachine.TypeInit, provider.Machine().Type())
			case i < len(controlplanes):
				assert.Equal(t, talosmachine.TypeControlPlane, provider.Machine().Type())
			default:
				assert.Equal(t, talosmachine.TypeWorker, provider.Machine().Type())
			}

			if provider.Machine().Type() != talosmachine.TypeWorker {
				// with user config, can only generate config for controlplane nodes
				assertClientConfig(t, talosConfig)
			}
		}

		assertClusterCA(ctx, t, c, cluster, assertMachineConfiguration(ctx, t, c, controlplanes[0]))

		// compare control plane secrets completely
		assertSameMachineConfigSecrets(ctx, t, c, controlplanes...)

		// compare all configs in more relaxed mode
		assertCompatibleMachineConfigs(ctx, t, c, append(append([]*bootstrapv1alpha3.TalosConfig{}, controlplanes...), workers...)...)
	})

	t.Run("InfrastructureNotReady", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, false)
		talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
		})
		createMachine(ctx, t, c, cluster, talosConfig, true)

		// assert that controller reports condition

		for ctx.Err() == nil {
			key := types.NamespacedName{
				Namespace: talosConfig.Namespace,
				Name:      talosConfig.Name,
			}

			err := c.Get(ctx, key, talosConfig)
			require.NoError(t, err)

			if conditions.IsFalse(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition) &&
				conditions.GetReason(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition) == bootstrapv1alpha3.WaitingForClusterInfrastructureReason {
				break
			}

			t.Log("Waiting ...")
			sleepCtx(ctx, 3*time.Second)
		}

		require.NoError(t, ctx.Err())

		assert.Equal(t, capiv1.ConditionSeverityInfo, *conditions.GetSeverity(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition))

		// patch to make infrastructure ready
		patchHelper, err := patch.NewHelper(cluster, c)
		require.NoError(t, err)

		cluster.Status.InfrastructureReady = true
		require.NoError(t, patchHelper.Patch(ctx, cluster))

		waitForReady(ctx, t, c, talosConfig)

		assertClientConfig(t, talosConfig)
		assertClusterClientConfig(ctx, t, c, cluster)

		provider := assertMachineConfiguration(ctx, t, c, talosConfig)

		assert.Equal(t, talosmachine.TypeInit, provider.Machine().Type())

		assertClusterCA(ctx, t, c, cluster, provider)

		assertControllerSecret(ctx, t, c, cluster, provider)
	})

	t.Run("BadConfigPatch", func(t *testing.T) {
		t.Parallel()

		namespaceName := setupTest(ctx, t, c)
		cluster := createCluster(ctx, t, c, namespaceName, nil, true)
		talosConfig := createTalosConfig(ctx, t, c, namespaceName, bootstrapv1alpha3.TalosConfigSpec{
			GenerateType: talosmachine.TypeInit.String(),
			TalosVersion: TalosVersion,
			ConfigPatches: []bootstrapv1alpha3.ConfigPatches{
				{
					Op:   "add",
					Path: "/machine/time/servers",
					Value: apiextensions.JSON{
						Raw: []byte(`["time.cloudflare.com"]`),
					},
				},
			},
		})
		createMachine(ctx, t, c, cluster, talosConfig, true)

		// assert that controller reports failure condition
		for ctx.Err() == nil {
			key := types.NamespacedName{
				Namespace: talosConfig.Namespace,
				Name:      talosConfig.Name,
			}

			err := c.Get(ctx, key, talosConfig)
			require.NoError(t, err)

			if conditions.IsFalse(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition) &&
				conditions.GetReason(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition) == bootstrapv1alpha3.DataSecretGenerationFailedReason {
				break
			}

			t.Log("Waiting ...")
			sleepCtx(ctx, 3*time.Second)
		}

		require.NoError(t, ctx.Err())

		assert.Equal(t, capiv1.ConditionSeverityError, *conditions.GetSeverity(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition))
		assert.Equal(t,
			"failure applying rfc6902 patches to talos machine config: add operation does not apply: doc is missing path: \"/machine/time/servers\": missing value",
			conditions.GetMessage(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition))
	})
}

// legacy cluster secret format
const legacySecretData = `{
"certs": "YWRtaW46CiAgY3J0OiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VKTmFrTkNOV0ZCUkVGblJVTkJhRUpNZEV4U1R5OHlOMGQ0TkRsaVJXSnpVbVZYT1V0TlFWVkhRWGwwYkdORVFWRk5VVFIzUkVGWlJGWlJVVXNLUlhkV01GbFhlSFpqZWtGbFJuY3dlVTFVUVRWTmFrRjRUMFJSZUU1RWJHRkdkekI2VFZSQk5VMVVaM2hQUkZGNFRrUnNZVTFDVFhoRlZFRlFRbWRPVmdwQ1FXOVVRMGM1ZWs5dFJtdGlWMngxVFVOdmQwSlJXVVJMTWxaM1FYbEZRV0ZpWVZKU2MzTllaMDVhZUd3MVoyeG1abWx6V1ZkcFpIZG1OVUZHV1dwRENscEhheTlLT1Roek5FNXRhbFZxUWxGTlFUUkhRVEZWWkVSM1JVSXZkMUZGUVhkSlNHZEVRV1JDWjA1V1NGTlZSVVpxUVZWQ1oyZHlRbWRGUmtKUlkwUUtRVkZaU1V0M1dVSkNVVlZJUVhkSmQwaDNXVVJXVWpCcVFrSm5kMFp2UVZWblQwSlFWV1EwZDA4eFdXTmtaVmRhVjBacGNGRTRUbFJSYldkM1FsRlpSQXBMTWxaM1FUQkZRVGg0UlRNelRFeE9hbEpxVjFVdlJsZ3ZRVmRFZDFWc0t6RnNja3hRVkZRNU5UUXpXbHBtZDBncldVdDVWMkpqUmt4NlRFSnFaRXBKQ2pSSU5rZFFZekpqTVhwd2JqQlVZbHA0UzJJeWMxUmFjSGRZUkZKRVVUMDlDaTB0TFMwdFJVNUVJRU5GVWxSSlJrbERRVlJGTFMwdExTMEsKICBrZXk6IExTMHRMUzFDUlVkSlRpQkZSREkxTlRFNUlGQlNTVlpCVkVVZ1MwVlpMUzB0TFMwS1RVTTBRMEZSUVhkQ1VWbEVTekpXZDBKRFNVVkpUVTVQTjBkbVZrNTVlbU5XWmtoc1dtUXpaVEZSY0hOaU1WZGhRMDFMTkdWM2NYQmxVRTFwV1RsR2JBb3RMUzB0TFVWT1JDQkZSREkxTlRFNUlGQlNTVlpCVkVVZ1MwVlpMUzB0TFMwSwpldGNkOgogIGNydDogTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVSm1WRU5EUVZOVFowRjNTVUpCWjBsU1FVMXVPRGgyZVd0VE0wWnZTVnBFV0RZd04zSkhWamgzUTJkWlNVdHZXa2w2YWpCRlFYZFJkMFI2UlU0S1RVRnpSMEV4VlVWRGFFMUZXbGhTYWxwRVFXVkdkekI1VFZSQk5VMXFRWGhQUkZGNFRrUnNZVVozTUhwTlZFRTFUVlJuZUU5RVVYaE9SR3hoVFVFNGVBcEVWRUZNUW1kT1ZrSkJiMVJDUjFZd1dUSlJkMWRVUVZSQ1oyTnhhR3RxVDFCUlNVSkNaMmR4YUd0cVQxQlJUVUpDZDA1RFFVRlVVVlJLTDJGNVJucHhDa1JvTm5WYWN6SjNUbXBUZWt4emVHaGpVVTFwZWtaTFFWUktNWEkwVTA1bWFEVkRSM2RYU1V4elNtaFRZell4YUdocVZXZzNlRUZXWkRObGQxbG9iWFVLT0ZCUWN6TnNlRE5VUzJOd2J6SkZkMWg2UVU5Q1owNVdTRkU0UWtGbU9FVkNRVTFEUVc5UmQwaFJXVVJXVWpCc1FrSlpkMFpCV1VsTGQxbENRbEZWU0FwQmQwVkhRME56UjBGUlZVWkNkMDFEVFVFNFIwRXhWV1JGZDBWQ0wzZFJSazFCVFVKQlpqaDNTRkZaUkZaU01FOUNRbGxGUmtoVk1IWXJWR1l2U2pOVkNqSmlOM0l5TkdGemNqaGlja1ZaTUZaTlFXOUhRME54UjFOTk5EbENRVTFGUVRCalFVMUZVVU5KU0ZOdlkyaFdOMUZJY2psNVFXOXJiakJqYTNjNVdtY0tjemczYUU5eU1WRlZkbGRzUTI4d1ZUSTNSak5CYVVJMGJEVlJjREJ5U0VWR2NtWnhMeTlqZVVaTE9YaERVVzlFSzA5eFNVWnhWVEJ6YWsxVFYwUnZSZ3A2WnowOUNpMHRMUzB0UlU1RUlFTkZVbFJKUmtsRFFWUkZMUzB0TFMwSwogIGtleTogTFMwdExTMUNSVWRKVGlCRlF5QlFVa2xXUVZSRklFdEZXUzB0TFMwdENrMUlZME5CVVVWRlNVWlVZemxxVUU1QlZXbGtSbUV2WkUxR1NIcDFOV0ZuYm5jeFl6UkZkVnB4Vm1KTE1rbGFWa2xRUldSdlFXOUhRME54UjFOTk5Ea0tRWGRGU0c5VlVVUlJaMEZGTUVWNVpqSnphR00yWnpSbGNtMWlUbk5FV1RCemVUZE5XVmhGUkVsemVGTm5SWGxrWVN0RmFsZzBaVkZvYzBacFF6ZERXUXBWYms5MFdWbFpNVWxsT0ZGR1dHUXpjMGRKV25KMlJIbzNUalZqWkRCNWJrdFJQVDBLTFMwdExTMUZUa1FnUlVNZ1VGSkpWa0ZVUlNCTFJWa3RMUzB0TFFvPQprOHM6CiAgY3J0OiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VKcFZFTkRRVk1yWjBGM1NVSkJaMGxSVFRSaE1EUlNSWGhuVmpkQ1Fsb3ljVzFoZW5nelZFRkxRbWRuY1docmFrOVFVVkZFUWtSQlZrMVNUWGNLUlZGWlJGWlJVVXRGZDNCeVpGZEtiR050Tld4a1IxWjZUVUkwV0VSVVNYaE5SR3Q1VFVSRk5FNUVSVEJQVm05WVJGUk5lRTFFYTNoUFJFVTBUa1JGTUFwUFZtOTNSbFJGVkUxQ1JVZEJNVlZGUTJoTlMyRXpWbWxhV0VwMVdsaFNiR042UWxwTlFrMUhRbmx4UjFOTk5EbEJaMFZIUTBOeFIxTk5ORGxCZDBWSUNrRXdTVUZDVEdWNmNubG5NMUZZYlhCc1QxWlFOeXRoY0M5YVZGRkRVMnhNTTNGcFQyVldOMjB6UnpoM09ISjJVbUZtSzB4aE9VUXdaa05XU2psU2FpOEtWSGwxZFZGR2VGRXlNRE52WlZoUVNXWnRSVGxJY1hSa2FuZHhhbGxVUW1aTlFUUkhRVEZWWkVSM1JVSXZkMUZGUVhkSlEyaEVRV1JDWjA1V1NGTlZSUXBHYWtGVlFtZG5ja0puUlVaQ1VXTkVRVkZaU1V0M1dVSkNVVlZJUVhkSmQwUjNXVVJXVWpCVVFWRklMMEpCVlhkQmQwVkNMM3BCWkVKblRsWklVVFJGQ2tablVWVlhNSFpuT1VGa1VDOWFTelVyZVZJdk56TkNjR1oyVUZKSVRXdDNRMmRaU1V0dldrbDZhakJGUVhkUlJGTkJRWGRTVVVsblpIWlVUV0pxU0NzS05GaFBUVnA2UmtsRWFtNXhOREpKTDNOMVJIYzBZMjVIV0dOeWJGZGtTaXRoV2xsRFNWRkVkWEp5UlVGTGJWQnlUV2RPY1ZReWQxQTJTbGRaZVd4c1lRb3piRGQ1Vmpob1V6VkRaME53U2xSaFJXYzlQUW90TFMwdExVVk9SQ0JEUlZKVVNVWkpRMEZVUlMwdExTMHRDZz09CiAga2V5OiBMUzB0TFMxQ1JVZEpUaUJGUXlCUVVrbFdRVlJGSUV0RldTMHRMUzB0Q2sxSVkwTkJVVVZGU1VaTVRVdHJjMFExYlhOdFZqQTBLM2hPTTBGVFZtaHVOVmhWTW5WUGJFRkVabXRHWTNaRWFqQnlORGx2UVc5SFEwTnhSMU5OTkRrS1FYZEZTRzlWVVVSUlowRkZkRGRQZGt0RVpFSmxZVzFWTlZVdmRqVnhiamxzVGtGS1MxVjJaWEZKTlRWWWRXSmpZbnBFZVhVNVJuQXZOSFJ5TUZCU09BcEtWVzR4UjFBNVVFczJOVUZZUmtSaVZHVm9OV000YUN0WlZEQmxjVEV5VUVOblBUMEtMUzB0TFMxRlRrUWdSVU1nVUZKSlZrRlVSU0JMUlZrdExTMHRMUW89Cms4c2FnZ3JlZ2F0b3I6CiAgY3J0OiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VKWWFrTkRRVkZYWjBGM1NVSkJaMGxSWm5KTFV6RmthV001T0ZCc1ptMHZNMmw2WW1SV2FrRkxRbWRuY1docmFrOVFVVkZFUWtSQlFVMUNORmdLUkZSSmVFMUVhM2xOUkVVMFRrUkZNRTlXYjFoRVZFMTRUVVJyZUU5RVJUUk9SRVV3VDFadmQwRkVRbHBOUWsxSFFubHhSMU5OTkRsQlowVkhRME54UndwVFRUUTVRWGRGU0VFd1NVRkNRMUJqZFVJMkwycEVTa3RDT1RRck0xRlFZaTlQVWl0M1JqZEdaR2d3ZFhkb1UyOXlaak53ZUhSeFpqTkJZbGxOTkU4NUNteFlRek5yTVZwRE1tNUNkbEU1UXpac1ZGZEpWVVpYWkV0SFNtTkxhalpzZGpWNWFsbFVRbVpOUVRSSFFURlZaRVIzUlVJdmQxRkZRWGRKUTJoRVFXUUtRbWRPVmtoVFZVVkdha0ZWUW1kbmNrSm5SVVpDVVdORVFWRlpTVXQzV1VKQ1VWVklRWGRKZDBSM1dVUldVakJVUVZGSUwwSkJWWGRCZDBWQ0wzcEJaQXBDWjA1V1NGRTBSVVpuVVZVMVUwTlVUR054YVZwUGJ6WXhjVkJrWm5OdVZYZEtRV1ozVVZWM1EyZFpTVXR2V2tsNmFqQkZRWGRSUkZKM1FYZFNRVWxuQ2xNMGVVUlFRakZtZGxRNGRFbGxhR1ZWVEhkVFl6QXlWWFV2YkRWNVZHNUNVVTFUYlVZMldUazVVMDFEU1VWdlVUVnpMM0YwZEVKSldrWjNNMGd5Y1RBS1MzUTFOM1oxVjJKSUswOXdTMmhNVlV0MFZESkNjWFE1Q2kwdExTMHRSVTVFSUVORlVsUkpSa2xEUVZSRkxTMHRMUzBLCiAga2V5OiBMUzB0TFMxQ1JVZEpUaUJGUXlCUVVrbFdRVlJGSUV0RldTMHRMUzB0Q2sxSVkwTkJVVVZGU1VzNVN6VkpNVkI1TW1GU1RDdFdhU3R6ZFRVek5qUjBOMGxvTkVoS1MxTnRTVWMzY1RKNmRXMXdkVVZ2UVc5SFEwTnhSMU5OTkRrS1FYZEZTRzlWVVVSUlowRkZTVGw1TkVoeUswMU5hMjlJTTJvM1pFRTVkamcxU0RkQldITldNa2hUTjBOR1MybDBMMlZ1UnpKd0wyTkNkR2Q2WnpjeVZncGpUR1ZVVm10TVlXTkhPVVF3VEhGV1RsbG9VVlphTUc5WmJIZHhVSEZYTDI1QlBUMEtMUzB0TFMxRlRrUWdSVU1nVUZKSlZrRlVSU0JMUlZrdExTMHRMUW89Cms4c3NlcnZpY2VhY2NvdW50OgogIGtleTogTFMwdExTMUNSVWRKVGlCRlF5QlFVa2xXUVZSRklFdEZXUzB0TFMwdENrMUlZME5CVVVWRlNVbENTR3B1YzA5NWVGRTJNVkV5VFZGcVVrWTRSVUphVldaVVNHbElPVk41ZG1WSVIya3ZNVWhoYVc5dlFXOUhRME54UjFOTk5Ea0tRWGRGU0c5VlVVUlJaMEZGZFd4cWJuRTRTRGswZGxkSlZsbG9ZVEJEZGxkWVNFaHhNekJtUTFSbU1qaEZaVTB5UlUxU09XNVNTamxxZDBWVVVsWlNOUXBYSzFsT1ZITnVVRTR4ZFZOVU5rRlZRVGhZZGl0WVFqRlVkek5WYUVSdllsVkJQVDBLTFMwdExTMUZUa1FnUlVNZ1VGSkpWa0ZVUlNCTFJWa3RMUzB0TFFvPQpvczoKICBjcnQ6IExTMHRMUzFDUlVkSlRpQkRSVkpVU1VaSlEwRlVSUzB0TFMwdENrMUpTVUpRZWtOQ09IRkJSRUZuUlVOQmFFVkJjblk0YVZscVYxaERPRTFoZEdGaE9HVXljR1Y2UkVGR1FtZE5jbHBZUVhkRlJFVlBUVUYzUjBFeFZVVUtRMmhOUm1SSFJuTmlNMDEzU0doalRrMXFSWGRQVkVsM1RWUm5NRTFVVVRWWGFHTk9UWHBGZDA5VVJUUk5WR2N3VFZSUk5WZHFRVkZOVVRSM1JFRlpSQXBXVVZGTFJYZFdNRmxYZUhaamVrRnhUVUZWUjBGNWRHeGpRVTFvUVU5RFVrMXNSMDVxYzJSUmJXZHNjekpRUTFOblRXUk5aVUZKUWpobVFVdHpia053Q21GWVdETnlabFZMYnpKRmQxaDZRVTlDWjA1V1NGRTRRa0ZtT0VWQ1FVMURRVzlSZDBoUldVUldVakJzUWtKWmQwWkJXVWxMZDFsQ1FsRlZTRUYzUlVjS1EwTnpSMEZSVlVaQ2QwMURUVUU0UjBFeFZXUkZkMFZDTDNkUlJrMUJUVUpCWmpoM1NGRlpSRlpTTUU5Q1FsbEZSa2xFWjFReFNHVk5SSFJYU0VoWWJBcHRWbWhaY1ZWUVJGVXdTbTlOUVZWSFFYbDBiR05CVGtKQlJESkhURTh5ZGtjNVRVaEhlSFE1TmpVNFdEUjRXa3hUV1U1c1pFRm5SSGt5ZEVodFdqZHNDbTVCYWtGU01HNXdXbTlSV0VKV2FHOXljbEZGWTJWaE4yYzJWRzg1UWtSdGRIcHlSakJUZEZjNE9UVmtNRUZyUFFvdExTMHRMVVZPUkNCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2c9PQogIGtleTogTFMwdExTMUNSVWRKVGlCRlJESTFOVEU1SUZCU1NWWkJWRVVnUzBWWkxTMHRMUzBLVFVNMFEwRlJRWGRDVVZsRVN6SldkMEpEU1VWSlJtUkpZWE4xUzBNemVuZHBOREJFV2xNNVNUTldUSFJXUWtKVlNuaEdiVVYzYkdsdFVuWlZlSGhHWndvdExTMHRMVVZPUkNCRlJESTFOVEU1SUZCU1NWWkJWRVVnUzBWWkxTMHRMUzBLCg==",
"kubeSecrets": "Ym9vdHN0cmFwdG9rZW46IDVkd3pyaC41bXM5ZDVla2UxbXVza3JnCmFlc2NiY2VuY3J5cHRpb25zZWNyZXQ6IGp4aFBxMkM2TVJGYk5kQzdyZE5KU3dKbXNZM1lIMjNnUnpuYjdlZmhLTU09Cg==",
"trustdInfo": "dG9rZW46IG8xOXpoNy55djdyeGNlM2xzcHRubWU5Cg=="
}
`
