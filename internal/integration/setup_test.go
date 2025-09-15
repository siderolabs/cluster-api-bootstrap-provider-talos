// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	"context"
	"net"
	"net/url"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/siderolabs/go-pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	clusterctlclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientcfg "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	bootstrapv1alpha3 "github.com/siderolabs/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	"github.com/siderolabs/cluster-api-bootstrap-provider-talos/controllers"
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
	restCfg, testEnv := startTestEnv(ctx, t)

	c, err := client.New(restCfg, client.Options{})
	require.NoError(t, err)

	waitForCAPIAvailability(ctx, t, c)

	mgr, err := ctrl.NewManager(restCfg, ctrl.Options{
		WebhookServer: webhook.NewServer(
			webhook.Options{
				CertDir: testEnv.WebhookInstallOptions.LocalServingCertDir,
				Host:    testEnv.WebhookInstallOptions.LocalServingHost,
				Port:    testEnv.WebhookInstallOptions.LocalServingPort,
			},
		),
	})
	require.NoError(t, err)

	err = (&controllers.TalosConfigReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName(t.Name()),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: 10})
	require.NoError(t, err)

	err = (&bootstrapv1alpha3.TalosConfigTemplate{}).SetupWebhookWithManager(mgr)
	require.NoError(t, err)

	err = (&bootstrapv1alpha3.TalosConfig{}).SetupWebhookWithManager(mgr)
	require.NoError(t, err)

	go func() {
		assert.NoError(t, mgr.Start(ctx))
	}()

	<-mgr.Elected()

	waitForWebhooks(ctx, t, testEnv)

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
				GracePeriodSeconds: pointer.To(int64(0)),
			}

			t.Logf("Deleting namespace %q ...", namespace)
			assert.NoError(t, c.Delete(context.Background(), ns, opts)) // not ctx because it can be already canceled

			iteration := -1

			for {
				iteration++

				var obj corev1.Namespace

				err = c.Get(context.Background(), types.NamespacedName{Name: namespace}, &obj)
				if err == nil {
					if iteration%10 == 0 {
						t.Log("Waiting for ns deletion", namespace)
					}

					// a bit of black magic here:
					//   as we don't set infrastructureRef on machines, capi controller will hang forever
					//   trying to delete infrastructure data for the machine
					//   this allows us to override that removing the finalizer(s)
					var machineList capiv1.MachineList

					err = c.List(context.Background(), &machineList, client.InNamespace(namespace))
					if err != nil {
						t.Log("error listing machines", err)

						continue
					}

					for _, machine := range machineList.Items {
						machine.Finalizers = nil

						if err = c.Update(context.Background(), &machine); err != nil {
							// conflicts might be ignored here, as eventually this will succeed
							if !apierrors.IsConflict(err) {
								t.Log("error updating machine's finalizers", err)
							}
						}
					}

					time.Sleep(time.Second)

					continue
				}

				if !apierrors.IsNotFound(err) {
					t.Log("error waiting for namespace deletion", err)
				}

				break
			}
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
		clusterctlClient, err := clusterctlclient.New(ctx, "")
		if err != nil {
			initErr <- err
			return
		}

		t.Log("Getting CAPI core components versions ...")

		initOpts := clusterctlclient.InitOptions{
			BootstrapProviders:      []string{clusterctlclient.NoopProvider},
			InfrastructureProviders: []string{"docker"},
			ControlPlaneProviders:   []string{clusterctlclient.NoopProvider},
		}

		if false {
			// TODO: InitImages is broken in upstream, see https://github.com/kubernetes-sigs/cluster-api/issues/6986
			images, err := clusterctlClient.InitImages(ctx, initOpts)
			if err != nil {
				initErr <- err
				return
			}

			t.Logf("Installing CAPI core components: %s ...", strings.Join(images, ", "))
		}

		_, err = clusterctlClient.Init(ctx, initOpts)
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
func startTestEnv(ctx context.Context, t *testing.T) (*rest.Config, *envtest.Environment) {
	t.Helper()

	cfg, err := clientcfg.GetConfig()
	require.NoError(t, err)

	u, err := url.Parse(cfg.Host)
	require.NoError(t, err)

	// this is pure hack to support docker-based clusters
	// Docker says control plane endpoint is https://10.5.0.2:6443 (which is first node address)
	// we need 10.5.0.1 which is the bridge IP, which should be available both for the test and for the API server
	if u.Hostname() == "10.5.0.2" {
		u.Host = "10.5.0.1"
	}

	testEnv := &envtest.Environment{
		Config:            cfg,
		CRDDirectoryPaths: []string{filepath.Join("..", "..", Artifacts, "bootstrap-talos", Tag)},
		CRDInstallOptions: envtest.CRDInstallOptions{
			ErrorIfPathMissing: true,
			MaxTime:            20 * time.Second,
			PollInterval:       time.Second,
			CleanUpAfterUse:    !skipCleanup,
		},
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths:            []string{"../../config/webhook/manifests.yaml"},
			MaxTime:          10 * time.Second,
			LocalServingHost: u.Hostname(),
			PollInterval:     time.Second,
		},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    pointer.To(true),
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
	return res.cfg, testEnv
}

func waitForWebhooks(ctx context.Context, t *testing.T, testEnv *envtest.Environment) {
	host := testEnv.WebhookInstallOptions.LocalServingHost
	port := testEnv.WebhookInstallOptions.LocalServingPort

	t.Logf("Waiting for webhook port %d to be open prior to running tests", port)

	timeout := 1 * time.Second

	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			t.Fatalf("Failed to wait for webhook availability: %s.", ctx.Err())
		}

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)
		if err != nil {
			t.Logf("Webhook port is not ready, will retry in %v: %s", timeout, err)

			continue
		}

		conn.Close() //nolint:errcheck

		t.Logf("Webhook port is now open. Continuing with tests...")

		break
	}
}

// waitForCAPIAvailability waits for needed CAPI components availability.
//
// Context cancelation is honored.
func waitForCAPIAvailability(ctx context.Context, t *testing.T, c client.Client) {
	t.Helper()

	t.Log("Waiting for CAPI availability ...")

	key := client.ObjectKey{Namespace: "capi-system", Name: "capi-controller-manager"}

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
