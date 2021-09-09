// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"context"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
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

// setupSuite setups the whole test suite.
func setupSuite(t *testing.T) (context.Context, client.Client) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping in -short mode.")
	}

	ctx := context.Background()

	if !skipCleanup {
		// cancel context on first Ctrl+C, kill on second
		var stop context.CancelFunc
		ctx, stop = signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
		t.Cleanup(stop)
		go func() {
			<-ctx.Done()
			t.Log("Stopping...")
			stop()
		}()

		// reserve 1 minute for cleanup if possible
		deadline, ok := t.Deadline()
		if ok && time.Until(deadline) > 70*time.Second {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, deadline.Add(-60*time.Second))
			t.Cleanup(cancel)
		}
	}

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	installCAPI(ctx, t)
	restCfg := startTestEnv(ctx, t)

	c, err := client.New(restCfg, client.Options{})
	require.NoError(t, err)

	stopCAPI(ctx, t, c)

	waitForCAPIAvailability(ctx, t, c)

	// TODO(aleksi): make one manager per test / namespace (move to setupTest)?

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
	}()

	t.Log("Setup done.")

	return ctx, c
}

// setupTest setups one per-test (subtest) namespace.
func setupTest(ctx context.Context, t *testing.T, c client.Client) string {
	t.Helper()

	namespace := generateName(t, "ns")

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err := c.Create(ctx, ns)
	require.NoError(t, err)

	if !skipCleanup {
		t.Cleanup(func() {
			opts := &client.DeleteOptions{
				GracePeriodSeconds: pointer.ToInt64(0),
			}

			t.Logf("Deleting namespace %q ...", namespace)
			assert.NoError(t, c.Delete(context.Background(), ns, opts)) // not ctx because it can be already canceled
		})
	}

	return namespace
}

// installCAPI installs core CAPI components.
//
// Context cancelation is honored.
func installCAPI(ctx context.Context, t *testing.T) {
	t.Helper()

	// Run InitImages / Init in the goroutine to handle context cancelation.
	// t.FailNow() should be called in the main goroutine.
	initErr := make(chan error, 1)
	go func() {
		clusterctlClient, err := clusterctlclient.New("")
		if err != nil {
			initErr <- err
			return
		}

		t.Log("Getting CAPI core components versions ...")

		initOpts := clusterctlclient.InitOptions{
			BootstrapProviders:      []string{clusterctlclient.NoopProvider},
			InfrastructureProviders: []string{clusterctlclient.NoopProvider},
			ControlPlaneProviders:   []string{clusterctlclient.NoopProvider},
		}
		images, err := clusterctlClient.InitImages(initOpts)
		if err != nil {
			initErr <- err
			return
		}

		t.Logf("Installing CAPI core components: %s ...", strings.Join(images, ", "))

		_, err = clusterctlClient.Init(initOpts)
		initErr <- err
	}()

	var err error
	select {
	case err = <-initErr:
	case <-ctx.Done():
		err = ctx.Err()
	}

	require.NoError(t, err, "failed to install CAPI core components")

	t.Log("Done installing CAPI core components.")
}

// startTestEnv starts envtest environment: installs CRDs, etc.
//
// Context cancelation is honored.
func startTestEnv(ctx context.Context, t *testing.T) *rest.Config {
	t.Helper()

	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
		CRDInstallOptions: envtest.CRDInstallOptions{
			ErrorIfPathMissing: true,
			MaxTime:            20 * time.Second,
			PollInterval:       time.Second,
			CleanUpAfterUse:    !skipCleanup,
		},
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			// TODO paths?

			MaxTime:      10 * time.Second,
			PollInterval: time.Second,
		},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    pointer.ToBool(true),
	}

	// Run Start in the goroutine to handle context cancelation.
	// t.FailNow() should be called in the main goroutine.
	type result struct {
		cfg *rest.Config
		err error
	}
	startErr := make(chan result, 1)
	go func() {
		if !skipCleanup {
			t.Cleanup(func() {
				t.Log("Stopping test-env ...")

				if err := testEnv.Stop(); err != nil {
					t.Logf("Failed to stop test-env: %s.", err)
					return
				}

				t.Logf("Test-env stopped.")
			})
		}

		t.Log("Starting test-env ...")

		cfg, err := testEnv.Start()
		startErr <- result{cfg, err}
	}()

	var res result
	select {
	case res = <-startErr:
	case <-ctx.Done():
		res.err = ctx.Err()
	}

	require.NoError(t, res.err, "failed to start test-env")

	t.Logf("Test-env started: %s.", res.cfg.Host)
	return res.cfg
}

// stopCAPI stops CAPI components that we don't need so they don't interfere with our tests.
//
// Context cancelation is honored.
func stopCAPI(ctx context.Context, t *testing.T, c client.Client) {
	t.Helper()

	t.Log("Stopping CAPI components ...")

	var deployment appsv1.Deployment
	key := client.ObjectKey{Namespace: "capi-system", Name: "capi-controller-manager"}

	require.NoError(t, c.Get(ctx, key, &deployment))

	patchHelper, err := patch.NewHelper(&deployment, c)
	require.NoError(t, err)

	deployment.Spec.Replicas = pointer.ToInt32(0)

	require.NoError(t, patchHelper.Patch(ctx, &deployment))

	for {
		var deployment appsv1.Deployment

		require.NoError(t, c.Get(ctx, key, &deployment))

		if deployment.Status.Replicas == 0 {
			break
		}

		t.Logf("Waiting: %+v ...", deployment.Status)
		sleepCtx(ctx, 5*time.Second)
		if ctx.Err() != nil {
			t.Fatalf("Failed to stop CAPI components: %s.", ctx.Err())
		}
	}

	t.Log("Done stopping CAPI components.")
}

// waitForCAPIAvailability waits for needed CAPI components availability.
//
// Context cancelation is honored.
func waitForCAPIAvailability(ctx context.Context, t *testing.T, c client.Client) {
	t.Helper()

	t.Log("Waiting for CAPI availability ...")

	// TODO is is not entirely clear why we need webhooks

	key := client.ObjectKey{Namespace: "capi-webhook-system", Name: "capi-controller-manager"}

	for {
		var deployment appsv1.Deployment

		require.NoError(t, c.Get(ctx, key, &deployment))

		var available bool
		for _, cond := range deployment.Status.Conditions {
			if cond.Type != appsv1.DeploymentAvailable {
				continue
			}

			available = cond.Status == corev1.ConditionTrue
		}

		if available {
			break
		}

		t.Logf("Waiting: %+v ...", deployment.Status)
		sleepCtx(ctx, 5*time.Second)
		if ctx.Err() != nil {
			t.Fatalf("Failed to wait for CAPI availability: %s.", ctx.Err())
		}
	}

	t.Log("CAPI is available.")
}
