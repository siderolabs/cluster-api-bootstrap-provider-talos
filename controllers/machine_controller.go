package controllers

import (
	"context"
	"fmt"
	machine2 "github.com/siderolabs/talos/pkg/machinery/api/machine"
	"time"

	"github.com/go-logr/logr"
	talosclient "github.com/siderolabs/talos/pkg/machinery/client"
	talosconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MachineReconciler reconciles a Machine object
type MachineReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	WatchFilterValue string
}

var MachineHookAnnotationTalosReset = capiv1.PreTerminateDeleteHookAnnotationPrefix + "/talos-reset"

func (r *MachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	r.Scheme = mgr.GetScheme()

	return ctrl.NewControllerManagedBy(mgr).
		For(&capiv1.Machine{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Complete(r)
}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get

func (r *MachineReconciler) Reconcile(ctx context.Context, req reconcile.Request) (_ reconcile.Result, rerr error) {
	log := r.Log.WithName(controllerName).
		WithName(fmt.Sprintf("namespace=%s", req.Namespace)).
		WithName(fmt.Sprintf("machine=%s", req.Name))

	// Lookup the talosconfig config
	machine := &capiv1.Machine{}
	if err := r.Client.Get(ctx, req.NamespacedName, machine); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get machine.")
		return ctrl.Result{}, err
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(machine, r.Client)
	if err != nil {
		log.Error(err, "Could not create a patchHelper for config.")
		return ctrl.Result{}, err
	}

	// Patch machine after each reconciliation
	defer func() {
		if err := patchHelper.Patch(ctx, machine); err != nil {
			log.Error(err, "Failed to patch machine.")
			if rerr == nil {
				rerr = err
			}
		}
	}()

	preTerminateDeleteHookCondition := conditions.Get(machine, capiv1.PreTerminateDeleteHookSucceededCondition)
	if !machine.ObjectMeta.DeletionTimestamp.IsZero() &&
		annotations.HasWithPrefix(MachineHookAnnotationTalosReset, machine.ObjectMeta.Annotations) &&
		preTerminateDeleteHookCondition != nil &&
		preTerminateDeleteHookCondition.Status == corev1.ConditionFalse {
		return r.resetMachine(ctx, machine)
	}

	return reconcile.Result{}, nil
}

func (r *MachineReconciler) resetMachine(ctx context.Context, machine *capiv1.Machine) (reconcile.Result, error) {
	log := r.Log.WithName(controllerName).
		WithName(fmt.Sprintf("namespace=%s", machine.GetNamespace())).
		WithName(fmt.Sprintf("machine=%s", machine.GetName()))

	talosClient, err := r.talosconfigForMachine(ctx, machine)
	if err != nil {
		log.Error(err, "Could not create talos client for machine")
		return reconcile.Result{Requeue: true}, err
	}

	defer talosClient.Close() //nolint:errcheck

	var address string

	// Prefer finding an InternalIP address for the machine first.
	// Fallback to finding an ExternalIP address for the machine
	// if no InternalIP is found.
	for _, addr := range machine.Status.Addresses {
		if addr.Type == capiv1.MachineInternalIP {
			address = addr.Address
			break
		}
		if addr.Type == capiv1.MachineExternalIP {
			address = addr.Address
		}
	}

	if address == "" {
		log.Error(nil, "No node addresses were found. Assuming node was never provisioned.")
		delete(machine.ObjectMeta.Annotations, MachineHookAnnotationTalosReset)
		return reconcile.Result{}, nil
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer ctxCancel()

	err = talosClient.ResetGeneric(talosclient.WithNode(ctxTimeout, address), &machine2.ResetRequest{
		Graceful: true,
		Reboot:   false,
	})
	if err != nil {
		log.Info("Failed to send Talos reset request to machine. Assuming node is already reset.")
		delete(machine.ObjectMeta.Annotations, MachineHookAnnotationTalosReset)
		return reconcile.Result{}, nil
	}
	log.Info("Talos node reset request successfully sent.")
	return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
}

func (r *MachineReconciler) talosconfigForMachine(ctx context.Context, machine *capiv1.Machine) (*talosclient.Client, error) {
	var (
		talosconfigSecret corev1.Secret
		clusterName       = machine.GetLabels()["cluster.x-k8s.io/cluster-name"]
	)

	if err := r.Client.Get(ctx,
		types.NamespacedName{
			Namespace: machine.GetNamespace(),
			Name:      clusterName + "-talosconfig",
		},
		&talosconfigSecret,
	); err != nil {
		return nil, err
	}

	t, err := talosconfig.FromBytes(talosconfigSecret.Data["talosconfig"])
	if err != nil {
		return nil, err
	}

	return talosclient.New(ctx, talosclient.WithConfig(t))
}
