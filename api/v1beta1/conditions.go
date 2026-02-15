// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1beta1

import capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

// Conditions and condition Reasons for the TalosConfig object

const (
	// TalosConfigReadyCondition condition is true when all ClientConfigAvailable and DataSecretAvailable conditions are true.
	TalosConfigReadyCondition = capiv1.ReadyCondition

	// TalosConfigReadyReason surfaces when TalosConfig is ready.
	TalosConfigReadyReason = capiv1.ReadyReason

	// TalosConfigNotReadyReason surfaces when TalosConfig is not ready.
	TalosConfigNotReadyReason = capiv1.NotReadyReason

	// TalosConfigReadyUnknownReason surfaces when TalosConfig readiness is unknown.
	TalosConfigReadyUnknownReason = capiv1.ReadyUnknownReason
)

const (
	// ClientConfigAvailableCondition documents the status of the client config generation process.
	ClientConfigAvailableCondition = "ClientConfigAvailable"

	// ClientConfigAvailableReason surfaces when generated Talos client config is available.
	ClientConfigAvailableReason = capiv1.AvailableReason

	// ClientConfigAvailableInternalErrorReason surfaces unexpected failures when reading or generating
	// Talos client config.
	ClientConfigAvailableInternalErrorReason = capiv1.InternalErrorReason
)

const (
	// DataSecretAvailableCondition documents the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	DataSecretAvailableCondition = "DataSecretAvailable"

	// DataSecretAvailableReason surfaces when Talos bootstrap secret is available.
	DataSecretAvailableReason = capiv1.AvailableReason

	// DataSecretNotAvailableCondition surfaces when Talos bootstrap secret is not available.
	DataSecretNotAvailableReason = capiv1.NotAvailableReason

	// DataSecretNotAvailableInternalErrorReason surfaces unexpected failures when generating Talos bootstrap secret.
	DataSecretNotAvailableInternalErrorReason = capiv1.InternalErrorReason
)
