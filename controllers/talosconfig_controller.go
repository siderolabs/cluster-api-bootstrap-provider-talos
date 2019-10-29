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

	"github.com/go-logr/logr"
	bootstrapv1alpha2 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha2"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	controllerName = "cabpt-controller"
)

// TalosConfigReconciler reconciles a TalosConfig object
type TalosConfigReconciler struct {
	client.Client
	Log logr.Logger
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

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=talosconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=talosconfigs/status,verbs=get;update;patch\
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status;machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;events;configmaps,verbs=get;list;watch;create;update;patch;delete
func (r *TalosConfigReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, rerr error) {
	ctx := context.Background()
	log := r.Log.WithName(controllerName).
		WithName(fmt.Sprintf("namespace=%s", req.Namespace)).
		WithName(fmt.Sprintf("talosconfig=%s", req.Name))

	// Lookup the talosconfig config
	config := &bootstrapv1alpha2.TalosConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to get config")
		return ctrl.Result{}, err
	}

	// bail super early if it's already ready
	if config.Status.Ready {
		log.Info("ignoring an already ready config")
		return ctrl.Result{}, nil
	}

	// Look up the Machine that owns this talosconfig if there is one
	machine, err := util.GetOwnerMachine(ctx, r.Client, config.ObjectMeta)
	if err != nil {
		log.Error(err, "could not get owner machine")
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Waiting for Machine Controller to set OwnerRef on the talosconfig")
		return ctrl.Result{}, nil
	}
	log = log.WithName(fmt.Sprintf("machine-name=%s", machine.Name))

	// Ignore machines that already have bootstrap data
	if machine.Spec.Bootstrap.Data != nil {
		// TODO: mark the config as ready?
		return ctrl.Result{}, nil
	}

	// Lookup the cluster the machine is associated with
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Error(err, "could not get cluster by machine metadata")
		return ctrl.Result{}, err
	}

	// Wait patiently for the infrastructure to be ready
	if !cluster.Status.InfrastructureReady {
		log.Info("Infrastructure is not ready, waiting until ready.")
		return ctrl.Result{}, errors.New("infra not ready")
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

	// Determine what type of node this is
	machineType := generate.TypeJoin
	switch config.Spec.MachineType {
	case "init":
		machineType = generate.TypeInit
	case "controlplane":
		machineType = generate.TypeControlPlane
	}

	APIEndpointPort := strconv.Itoa(cluster.Status.APIEndpoints[0].Port)
	input, err := generate.NewInput(cluster.ObjectMeta.Name,
		"https://"+cluster.Status.APIEndpoints[0].Host+":"+APIEndpointPort,
		*machine.Spec.Version,
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	inputSecret, err := r.fetchInputSecret(ctx, config, cluster.ObjectMeta.Name)
	if machineType == generate.TypeInit && k8serrors.IsNotFound(err) {
		err = r.writeInputSecret(ctx, config, cluster.ObjectMeta.Name, input)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	certs := &generate.Certs{}
	kubeTokens := &generate.KubeadmTokens{}
	trustdInfo := &generate.TrustdInfo{}

	err = yaml.Unmarshal(inputSecret.Data["certs"], certs)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = yaml.Unmarshal(inputSecret.Data["kubeTokens"], kubeTokens)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = yaml.Unmarshal(inputSecret.Data["trustdInfo"], trustdInfo)
	if err != nil {
		return ctrl.Result{}, err
	}

	input.Certs = certs
	input.KubeadmTokens = kubeTokens
	input.TrustdInfo = trustdInfo

	talosConfig := &talosConfig{
		Context: input.ClusterName,
		Contexts: map[string]*talosConfigContext{
			input.ClusterName: {
				Target: "",
				CA:     base64.StdEncoding.EncodeToString(input.Certs.OS.Crt),
				Crt:    base64.StdEncoding.EncodeToString(input.Certs.Admin.Crt),
				Key:    base64.StdEncoding.EncodeToString(input.Certs.Admin.Key),
			},
		},
	}

	talosConfigBytes, err := yaml.Marshal(talosConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	data, err := generate.Config(machineType, input)
	if err != nil {
		return ctrl.Result{}, err
	}

	config.Status.BootstrapData = []byte(data)
	config.Status.TalosConfig = string(talosConfigBytes)
	config.Status.Ready = true

	return ctrl.Result{}, nil
}

func (r *TalosConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bootstrapv1alpha2.TalosConfig{}).
		Complete(r)
}
