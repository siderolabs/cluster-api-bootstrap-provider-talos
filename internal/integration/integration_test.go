// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	// +kubebuilder:scaffold:imports
)

func TestIntegration(t *testing.T) {
	ctx, c := setupSuite(t)

	// namespaced objects
	var (
		clusterName     = "test-cluster"
		machineName     = "test-machine"
		dataSecretName  = "test-secret"
		talosConfigName = "test-config"
	)

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()
		namespaceName := setupTest(ctx, t, c)

		cluster := &capiv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespaceName,
				Name:      clusterName,
			},
			Spec: capiv1.ClusterSpec{
				ClusterNetwork: &capiv1.ClusterNetwork{
					Pods: &capiv1.NetworkRanges{
						CIDRBlocks: []string{"192.168.0.0/16"},
					},
					ServiceDomain: "cluster.local",
					Services: &capiv1.NetworkRanges{
						CIDRBlocks: []string{"10.128.0.0/12"},
					},
				},
			},
		}
		require.NoError(t, c.Create(ctx, cluster), "can't create a cluster")

		cluster.Status.InfrastructureReady = true
		require.NoError(t, c.Status().Update(ctx, cluster))

		machine := &capiv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespaceName,
				Name:      machineName,
			},
			Spec: capiv1.MachineSpec{
				ClusterName: cluster.Name,
				Bootstrap: capiv1.Bootstrap{
					DataSecretName: &dataSecretName,
				},
			},
		}

		require.NoError(t, controllerutil.SetOwnerReference(cluster, machine, scheme.Scheme))
		require.NoError(t, c.Create(ctx, machine))

		config := &bootstrapv1alpha3.TalosConfig{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespaceName,
				Name:      talosConfigName,
			},
			Spec: bootstrapv1alpha3.TalosConfigSpec{
				GenerateType: "init",
			},
		}
		require.NoError(t, controllerutil.SetOwnerReference(machine, config, scheme.Scheme))

		err := c.Create(ctx, config)
		require.NoError(t, err)

		for ctx.Err() == nil {
			key := types.NamespacedName{
				Namespace: namespaceName,
				Name:      talosConfigName,
			}

			err = c.Get(ctx, key, config)
			require.NoError(t, err)

			if config.Status.Ready {
				break
			}

			t.Logf("Config: %+v", config)
			time.Sleep(5 * time.Second)
		}
	})
}
