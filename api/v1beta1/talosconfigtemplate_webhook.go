// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1beta1

import (
	"context"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *TalosConfigTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithValidator(r).
		Complete()
}

//+kubebuilder:webhook:verbs=update,path=/validate-bootstrap-cluster-x-k8s-io-v1beta1-talosconfigtemplate,mutating=false,failurePolicy=fail,groups=bootstrap.cluster.x-k8s.io,resources=talosconfigtemplates,versions=v1beta1,name=vtalosconfigtemplate.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

var _ admission.Validator[*TalosConfigTemplate] = &TalosConfigTemplate{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type
func (r *TalosConfigTemplate) ValidateCreate(ctx context.Context, obj *TalosConfigTemplate) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type
func (r *TalosConfigTemplate) ValidateUpdate(ctx context.Context, oldObj *TalosConfigTemplate, newObj *TalosConfigTemplate) (admission.Warnings, error) {
	old := oldObj
	r = newObj

	if !cmp.Equal(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("TalosConfigTemplate.Spec is immutable")
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type
func (r *TalosConfigTemplate) ValidateDelete(ctx context.Context, obj *TalosConfigTemplate) (admission.Warnings, error) {
	return nil, nil
}
