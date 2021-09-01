/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"context"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	clusterctlclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/talos-systems/cluster-api-bootstrap-provider-talos/controllers"
	// +kubebuilder:scaffold:imports
)

func setup(t *testing.T, doCleanup bool, namespace string) (context.Context, client.Client) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping in -short mode.")
	}

	// cancel context on first Ctrl+C, kill on second
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
	t.Cleanup(cancel)
	go func() {
		<-ctx.Done()
		t.Log("Stopping...")
		cancel()
	}()

	// reserve 1 minute for cleanup if possible
	if doCleanup {
		deadline, ok := t.Deadline()
		if ok && time.Until(deadline) > 70*time.Second {
			var stop context.CancelFunc
			ctx, stop = context.WithDeadline(ctx, deadline.Add(-60*time.Second))
			t.Cleanup(stop)
		}
	}

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	installCAPI(ctx, t)
	restCfg := startTestEnv(ctx, t, doCleanup)

	c, err := client.New(restCfg, client.Options{})
	require.NoError(t, err)

	stopCAPI(ctx, t, c)

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err = c.Create(ctx, ns)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, c.Delete(context.Background(), ns)) // not ctx because it can be already canceled
	})

	mgr, err := ctrl.NewManager(restCfg, ctrl.Options{})
	require.NoError(t, err)

	err = (&controllers.TalosConfigReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName(t.Name()),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: 10})
	require.NoError(t, err)

	go func() {
		assert.NoError(t, mgr.Start(ctx.Done()))
		cancel()
	}()

	t.Log("Setup done.")

	return ctx, c
}

// installCAPI installs core CAPI components. Context cancelation is honored.
func installCAPI(ctx context.Context, t *testing.T) {
	t.Helper()

	// t.FailNow() should be called in the main goroutine, so use assert, not require
	done := make(chan struct{})
	go func() {
		defer close(done)

		clusterctlClient, err := clusterctlclient.New("")
		if !assert.NoError(t, err) {
			return
		}

		t.Log("Getting CAPI core components versions ...")

		initOpts := clusterctlclient.InitOptions{
			BootstrapProviders:      []string{clusterctlclient.NoopProvider},
			InfrastructureProviders: []string{clusterctlclient.NoopProvider},
			ControlPlaneProviders:   []string{clusterctlclient.NoopProvider},
		}
		images, err := clusterctlClient.InitImages(initOpts)
		if !assert.NoError(t, err) {
			return
		}

		t.Logf("Installing CAPI core components: %s ...", strings.Join(images, ", "))

		_, err = clusterctlClient.Init(initOpts)
		if !assert.NoError(t, err) {
			return
		}

		t.Log("Done installing CAPI core components.")
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	assert.NoError(t, ctx.Err())

	if t.Failed() {
		t.FailNow()
	}
}

// startTestEnv starts envtest environment: installs CRDs, etc. Context cancelation is honored.
func startTestEnv(ctx context.Context, t *testing.T, doCleanup bool) *rest.Config {
	t.Helper()

	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
		CRDInstallOptions: envtest.CRDInstallOptions{
			ErrorIfPathMissing: true,
			CleanUpAfterUse:    doCleanup,
		},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    pointer.BoolPtr(true),
	}

	// t.FailNow() should be called in the main goroutine, so send errors to channel
	done := make(chan struct{})
	var cfg *rest.Config
	go func() {
		defer close(done)

		if doCleanup {
			t.Cleanup(func() {
				t.Log("Stopping test-env ...")

				if !assert.NoError(t, testEnv.Stop()) {
					return
				}

				t.Log("Test-env stopped.")
			})
		}

		t.Log("Starting test-env ...")

		var err error
		cfg, err = testEnv.Start()
		assert.NoError(t, err)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	assert.NoError(t, ctx.Err())

	if t.Failed() {
		t.FailNow()
	}

	t.Logf("Test-env started: %s.", cfg.Host)
	return cfg
}

// stopCAPI stops CAPI components so they don't interfere with our tests.
func stopCAPI(ctx context.Context, t *testing.T, c client.Client) {
	t.Helper()

	t.Log("Stopping CAPI components ...")

	var deployment appsv1.Deployment

	require.NoError(t, c.Get(ctx, client.ObjectKey{Namespace: "capi-system", Name: "capi-controller-manager"}, &deployment))

	patchHelper, err := patch.NewHelper(&deployment, c)
	require.NoError(t, err)

	deployment.Spec.Replicas = pointer.Int32Ptr(0)

	require.NoError(t, patchHelper.Patch(ctx, &deployment))

	for ctx.Err() == nil {
		var deployment appsv1.Deployment

		require.NoError(t, c.Get(ctx, client.ObjectKey{Namespace: "capi-system", Name: "capi-controller-manager"}, &deployment))

		if deployment.Status.Replicas == 0 {
			break
		}

		t.Logf("Waiting: %+v ...", deployment.Status)

		time.Sleep(5 * time.Second)
	}

	assert.NoError(t, ctx.Err())

	if t.Failed() {
		t.FailNow()
	}

	t.Log("Done stopping CAPI components.")
}
