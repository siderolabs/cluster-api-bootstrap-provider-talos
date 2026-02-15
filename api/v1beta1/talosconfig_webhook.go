// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1beta1

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *TalosConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(r).
		Complete()
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-bootstrap-cluster-x-k8s-io-v1beta1-talosconfig,mutating=false,failurePolicy=fail,groups=bootstrap.cluster.x-k8s.io,resources=talosconfigs,versions=v1beta1,name=vtalosconfig.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.CustomValidator = &TalosConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TalosConfig) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	r = obj.(*TalosConfig)

	return nil, r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TalosConfig) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (admission.Warnings, error) {
	old := oldObj.(*TalosConfig)
	r = newObj.(*TalosConfig)

	if !cmp.Equal(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("TalosConfig.Spec is immutable")
	}

	return nil, r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TalosConfig) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (r *TalosConfig) validate() error {
	var allErrs field.ErrorList

	switch r.Spec.Hostname.Source {
	case "":
	case HostnameSourceMachineName:
	case HostnameSourceInfrastructureName:
	default:
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("hostname").Child("source"), r.Spec.Hostname.Source,
				fmt.Sprintf("valid values are: %q", []HostnameSource{HostnameSourceMachineName, HostnameSourceInfrastructureName}),
			),
		)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: GroupVersion.Group, Kind: "TalosConfig"},
		r.Name, allErrs)
}
