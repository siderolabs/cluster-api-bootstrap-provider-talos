// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	talosclientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	talosmachine "github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/config/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bootstrapv1alpha3 "github.com/siderolabs/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	// +kubebuilder:scaffold:imports
)

var skipCleanup bool

func init() {
	const env = "INTEGRATION_SKIP_CLEANUP"
	def, _ := strconv.ParseBool(os.Getenv(env))
	flag.BoolVar(&skipCleanup, "skip-cleanup", def, fmt.Sprintf("Cleanup after tests [%s]", env))
}

// sleepCtx blocks until ctx is canceled or timeout passed.
func sleepCtx(ctx context.Context, timeout time.Duration) {
	sCtx, sCancel := context.WithTimeout(ctx, timeout)
	defer sCancel()
	<-sCtx.Done()
}

// generateName generates a unique name.
func generateName(t *testing.T, kind string) string {
	// use milliseconds since UTC midnight: unique enough, short enough, ordered
	now := time.Now().UTC()
	clock := time.Duration(now.Hour())*time.Hour +
		time.Duration(now.Minute())*time.Minute +
		time.Duration(now.Second())*time.Second +
		time.Duration(now.Nanosecond())
	n := clock / time.Microsecond

	return fmt.Sprintf("%s-%s-%d", strings.ReplaceAll(strings.ToLower(t.Name()), "/", "-"), kind, n)
}

// createCluster creates a Cluster with "ready" infrastructure.
func createCluster(ctx context.Context, t *testing.T, c client.Client, namespaceName string, spec *capiv1.ClusterSpec) *capiv1.Cluster {
	t.Helper()

	clusterName := generateName(t, "cluster")

	if spec == nil {
		spec = &capiv1.ClusterSpec{
			ClusterNetwork: capiv1.ClusterNetwork{},
			ControlPlaneEndpoint: capiv1.APIEndpoint{
				Host: clusterName + ".host",
				Port: 12345,
			},
		}
	}

	cluster := &capiv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespaceName,
			Name:      clusterName,
		},
		Spec: *spec,
	}

	require.NoError(t, c.Create(ctx, cluster), "can't create a cluster")

	return cluster
}

// createMachine creates a Machine owned by the Cluster.
func createMachine(ctx context.Context, t *testing.T, c client.Client, cluster *capiv1.Cluster, talosconfig *bootstrapv1alpha3.TalosConfig, controlplane bool) *capiv1.Machine {
	t.Helper()

	labels := map[string]string{
		capiv1.ClusterNameLabel: cluster.Name,
	}

	if controlplane {
		labels[capiv1.MachineControlPlaneLabel] = ""
	}

	machineName := generateName(t, "machine")
	machine := &capiv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      machineName,
			Labels:    labels,
		},
		Spec: capiv1.MachineSpec{
			ClusterName: cluster.Name,
			Bootstrap: capiv1.Bootstrap{
				ConfigRef: capiv1.ContractVersionedObjectReference{
					Kind:     "TalosConfig",
					APIGroup: bootstrapv1alpha3.GroupVersion.Group,
					Name:     talosconfig.GetName(),
				},
			},
			InfrastructureRef: capiv1.ContractVersionedObjectReference{
				Name: generateName(t, "infrastructure"),
			},
		},
	}

	require.NoError(t, c.Create(ctx, machine))

	return machine
}

// patchMachineAddress adds an address to the machine.
func patchMachineAddress(ctx context.Context, t *testing.T, c client.Client, machine *capiv1.Machine, addr string) {
	t.Helper()

	require.NoError(t, c.Get(ctx, client.ObjectKeyFromObject(machine), machine))

	patchHelper, err := patch.NewHelper(machine, c)
	require.NoError(t, err)

	machine.Status.Addresses = append(machine.Status.Addresses, capiv1.MachineAddress{
		Type:    capiv1.MachineExternalIP,
		Address: addr,
	})

	require.NoError(t, patchHelper.Patch(ctx, machine))
}

// createTalosConfig creates a TalosConfig owned by the Machine.
func createTalosConfig(ctx context.Context, t *testing.T, c client.Client, namespaceName string, spec bootstrapv1alpha3.TalosConfigSpec) *bootstrapv1alpha3.TalosConfig {
	t.Helper()

	talosConfigName := generateName(t, "talosconfig")
	talosConfig := &bootstrapv1alpha3.TalosConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespaceName,
			Name:      talosConfigName,
		},
		Spec: spec,
	}

	require.NoError(t, c.Create(ctx, talosConfig))

	// TODO that should not be needed
	if !skipCleanup {
		t.Cleanup(func() {
			t.Logf("Deleting TalosConfig %q ...", talosConfigName)
			assert.NoError(t, c.Delete(context.Background(), talosConfig)) // not ctx because it can be already canceled
		})
	}

	return talosConfig
}

// waitForReady waits for TalosConfig to be reconciled (ready).
func waitForReady(ctx context.Context, t *testing.T, c client.Client, talosConfig *bootstrapv1alpha3.TalosConfig) {
	t.Helper()

	key := types.NamespacedName{
		Namespace: talosConfig.Namespace,
		Name:      talosConfig.Name,
	}

	for ctx.Err() == nil {
		err := c.Get(ctx, key, talosConfig)
		require.NoError(t, err)

		if talosConfig.Status.Ready {
			break
		}

		t.Log("Waiting ...")
		sleepCtx(ctx, 3*time.Second)
	}

	if ctx.Err() != nil {
		return
	}

	assert.True(t, conditions.IsTrue(talosConfig, bootstrapv1alpha3.DataSecretAvailableCondition))

	if talosConfig.Spec.GenerateType == talosmachine.TypeInit.String() || talosConfig.Spec.GenerateType == talosmachine.TypeControlPlane.String() {
		// wait for additional condition

		for ctx.Err() == nil {
			err := c.Get(ctx, key, talosConfig)
			require.NoError(t, err)

			if conditions.IsTrue(talosConfig, bootstrapv1alpha3.ClientConfigAvailableCondition) {
				break
			}

			t.Log("Waiting for talosconfig...")
			sleepCtx(ctx, 3*time.Second)
		}
	}
}

// waitForEndpointsClusterClientConfig waits for cluster-wide talosconfig to have a number of endpoints.
func waitForEndpointsClusterClientConfig(ctx context.Context, t *testing.T, c client.Client, cluster *capiv1.Cluster, numEndpoints int) {
	t.Helper()

	var secret corev1.Secret

	key := client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-talosconfig",
	}

	for ctx.Err() == nil {
		err := c.Get(ctx, key, &secret)
		require.NoError(t, err)

		clientConfig, err := talosclientconfig.FromString(string(secret.Data["talosconfig"]))
		require.NoError(t, err)

		require.NotNil(t, clientConfig.Contexts[clientConfig.Context])

		if len(clientConfig.Contexts[clientConfig.Context].Endpoints) == numEndpoints {
			break
		}

		t.Logf("Waiting for %d endpoints...", numEndpoints)
		sleepCtx(ctx, time.Second)
	}
}

// validateClientConfig validates talosctl configuration.
func validateClientConfig(t *testing.T, config *talosclientconfig.Config, endpoints ...string) {
	t.Helper()

	assert.Len(t, config.Contexts, 1)
	assert.NotEmpty(t, config.Context)
	context := config.Contexts[config.Context]
	require.NotNil(t, context)

	if endpoints == nil {
		endpoints = []string{}
	}

	assert.Equal(t, endpoints, context.Endpoints)
	assert.Empty(t, context.Nodes)

	assert.NotEmpty(t, context.CA)
	assert.NotEmpty(t, context.Crt)
	assert.NotEmpty(t, context.Key)
}

type runtimeMode struct {
	requiresInstall bool
}

func (m runtimeMode) String() string {
	return fmt.Sprintf("runtimeMode(%v)", m.requiresInstall)
}

func (m runtimeMode) RequiresInstall() bool {
	return m.requiresInstall
}

func (m runtimeMode) InContainer() bool {
	return false
}

// check interface
var _ validation.RuntimeMode = runtimeMode{}
