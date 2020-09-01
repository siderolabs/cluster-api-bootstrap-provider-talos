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

package controllers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	"github.com/talos-systems/cluster-api-bootstrap-provider-talos/pkg/constants"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	configmachine "github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	controllerName = "cabpt-controller"
)

// TalosConfigReconciler reconciles a TalosConfig object
type TalosConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type TalosConfigScope struct {
	Config  *bootstrapv1alpha3.TalosConfig
	Machine *capiv1.Machine
	Cluster *capiv1.Cluster
}

type TalosConfigBundle struct {
	BoostrapData string
	TalosConfig  string
}

type talosConfig struct {
	Context  string
	Contexts map[string]*talosConfigContext
}

type talosConfigContext struct {
	Target string
	CA     string
	Crt    string
	Key    string
}

func (r *TalosConfigReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&bootstrapv1alpha3.TalosConfig{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=talosconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=talosconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status;machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *TalosConfigReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, rerr error) {
	ctx := context.Background()
	log := r.Log.WithName(controllerName).
		WithName(fmt.Sprintf("namespace=%s", req.Namespace)).
		WithName(fmt.Sprintf("talosconfig=%s", req.Name))

	// Lookup the talosconfig config
	config := &bootstrapv1alpha3.TalosConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to get config")
		return ctrl.Result{}, err
	}

	// Look up the Machine that owns this talosconfig if there is one
	machine, err := util.GetOwnerMachine(ctx, r.Client, config.ObjectMeta)
	if err != nil {
		log.Error(err, "could not get owner machine")
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Waiting for Machine Controller to set OwnerRef on the talosconfig")
		return ctrl.Result{}, errors.New("no owner ref")
	}
	log = log.WithName(fmt.Sprintf("machine-name=%s", machine.Name))

	// Lookup the cluster the machine is associated with
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Error(err, "could not get cluster by machine metadata")
		return ctrl.Result{}, err
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(config, r)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Always attempt to Patch the KubeadmConfig object and status after each reconciliation.
	defer func() {
		if err := patchHelper.Patch(ctx, config); err != nil {
			log.Error(err, "failed to patch config")
			if rerr == nil {
				rerr = err
			}
		}
	}()

	// If the talosConfig doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(config, bootstrapv1alpha3.ConfigFinalizer)

	// Handle deleted machines
	if !config.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, config)
	}

	// bail super early if it's already ready
	if config.Status.Ready {
		log.Info("ignoring an already ready config")
		return ctrl.Result{}, nil
	}

	// Wait patiently for the infrastructure to be ready
	if !cluster.Status.InfrastructureReady {
		log.Info("Infrastructure is not ready, waiting until ready.")
		return ctrl.Result{}, errors.New("infra not ready")
	}

	tcScope := &TalosConfigScope{
		Config:  config,
		Machine: machine,
		Cluster: cluster,
	}

	var retData *TalosConfigBundle

	switch config.Spec.GenerateType {
	// Slurp and use user-supplied configs
	case "none":
		if config.Spec.Data == "" {
			return ctrl.Result{}, errors.New("failed to specify config data with none generate type")
		}
		retData, err = r.userConfigs(ctx, tcScope)
		if err != nil {
			return ctrl.Result{}, err
		}

	// Generate configs on the fly
	case "init", "controlplane", "join":
		retData, err = r.genConfigs(ctx, tcScope)
		if err != nil {
			return ctrl.Result{}, err
		}

	default:
		return ctrl.Result{}, errors.New("unknown generate type specified")
	}

	// Packet acts a fool if you don't prepend #!talos to the userdata
	// so we try to suss out if that's the type of machine getting created.
	if machine.Spec.InfrastructureRef.Kind == "PacketMachine" {
		retData.BoostrapData = "#!talos\n" + retData.BoostrapData
	}

	err = r.writeBootstrapData(ctx, tcScope, []byte(retData.BoostrapData))
	if err != nil {
		return ctrl.Result{}, err
	}

	config.Status.DataSecretName = pointer.StringPtr(tcScope.Machine.Name + "-bootstrap-data")
	config.Status.TalosConfig = retData.TalosConfig
	config.Status.Ready = true

	return ctrl.Result{}, nil
}

func (r *TalosConfigReconciler) reconcileDelete(ctx context.Context, config *bootstrapv1alpha3.TalosConfig) (ctrl.Result, error) {
	controllerutil.RemoveFinalizer(config, bootstrapv1alpha3.ConfigFinalizer)

	return ctrl.Result{}, nil
}

func genTalosConfigFile(clusterName string, certs *generate.Certs) (string, error) {
	talosConfig := &talosConfig{
		Context: clusterName,
		Contexts: map[string]*talosConfigContext{
			clusterName: {
				Target: "",
				CA:     base64.StdEncoding.EncodeToString(certs.OS.Crt),
				Crt:    base64.StdEncoding.EncodeToString(certs.Admin.Crt),
				Key:    base64.StdEncoding.EncodeToString(certs.Admin.Key),
			},
		},
	}

	talosConfigBytes, err := yaml.Marshal(talosConfig)
	if err != nil {
		return "", err
	}

	return string(talosConfigBytes), nil
}

// userConfigs will fetch and make use of user-supplied bootstrap configs to return
func (r *TalosConfigReconciler) userConfigs(ctx context.Context, scope *TalosConfigScope) (*TalosConfigBundle, error) {
	retBundle := &TalosConfigBundle{}

	userConfig := &v1alpha1.Config{}
	err := yaml.Unmarshal([]byte(scope.Config.Spec.Data), userConfig)
	if err != nil {
		return retBundle, err
	}

	// Create the secret with kubernetes certs so a kubeconfig can be generated
	if userConfig.Machine().Type() == configmachine.TypeInit {
		err = r.writeK8sCASecret(ctx, scope, userConfig.Cluster().CA())
		if err != nil {
			return retBundle, err
		}
	}

	userConfigStr, err := userConfig.String()
	if err != nil {
		return retBundle, err
	}

	retBundle.BoostrapData = userConfigStr

	return retBundle, nil
}

// genConfigs will generate a bootstrap config and a talosconfig to return
func (r *TalosConfigReconciler) genConfigs(ctx context.Context, scope *TalosConfigScope) (*TalosConfigBundle, error) {
	retBundle := &TalosConfigBundle{}

	// Determine what type of node this is
	machineType := configmachine.TypeJoin
	switch scope.Config.Spec.GenerateType {
	case "init":
		machineType = configmachine.TypeInit
	case "controlplane":
		machineType = configmachine.TypeControlPlane
	}

	// Allow user to override default kube version.
	// This also handles version being formatted like "vX.Y.Z" instead of without leading 'v'
	// TrimPrefix returns the string unchanged if the prefix isn't present.
	k8sVersion := constants.DefaultKubeVersion
	if scope.Machine.Spec.Version != nil {
		k8sVersion = strings.TrimPrefix(*scope.Machine.Spec.Version, "v")
	}

	APIEndpointPort := strconv.Itoa(int(scope.Cluster.Spec.ControlPlaneEndpoint.Port))
	input, err := generate.NewInput(scope.Cluster.Name,
		"https://"+scope.Cluster.Spec.ControlPlaneEndpoint.Host+":"+APIEndpointPort,
		k8sVersion,
	)
	if err != nil {
		return retBundle, err
	}

	// Stash our generated input secrets so that we can reuse them for other nodes
	inputSecret, err := r.fetchSecret(ctx, scope.Config, scope.Cluster.Name+"-talos")
	if machineType == configmachine.TypeInit && k8serrors.IsNotFound(err) {
		inputSecret, err = r.writeInputSecret(ctx, scope, input)
		if err != nil {
			return retBundle, err
		}
	} else if err != nil {
		return retBundle, err
	}

	// Create the secret with kubernetes certs so a kubeconfig can be generated
	_, err = r.fetchSecret(ctx, scope.Config, scope.Cluster.Name+"-ca")
	if machineType == configmachine.TypeInit && k8serrors.IsNotFound(err) {
		err = r.writeK8sCASecret(ctx, scope, input.Certs.K8s)
		if err != nil {
			return retBundle, err
		}
	} else if err != nil {
		return retBundle, err
	}

	certs := &generate.Certs{}
	kubeSecrets := &generate.Secrets{}
	trustdInfo := &generate.TrustdInfo{}

	err = yaml.Unmarshal(inputSecret.Data["certs"], certs)
	if err != nil {
		return retBundle, err
	}

	err = yaml.Unmarshal(inputSecret.Data["kubeSecrets"], kubeSecrets)
	if err != nil {
		return retBundle, err
	}

	err = yaml.Unmarshal(inputSecret.Data["trustdInfo"], trustdInfo)
	if err != nil {
		return retBundle, err
	}

	input.Certs = certs
	input.Secrets = kubeSecrets
	input.TrustdInfo = trustdInfo

	tcString, err := genTalosConfigFile(input.ClusterName, input.Certs)
	if err != nil {
		return retBundle, err
	}

	retBundle.TalosConfig = tcString

	data, err := generate.Config(machineType, input)
	if err != nil {
		return retBundle, err
	}

	if scope.Cluster.Spec.ClusterNetwork.Pods != nil {
		data.ClusterConfig.ClusterNetwork.PodSubnet = scope.Cluster.Spec.ClusterNetwork.Pods.CIDRBlocks
	}
	if scope.Cluster.Spec.ClusterNetwork.Services != nil {
		data.ClusterConfig.ClusterNetwork.ServiceSubnet = scope.Cluster.Spec.ClusterNetwork.Services.CIDRBlocks
	}

	dataOut, err := data.String()
	if err != nil {
		return retBundle, err
	}

	retBundle.BoostrapData = dataOut

	return retBundle, nil
}
