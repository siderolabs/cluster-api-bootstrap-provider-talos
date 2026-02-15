// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"

// Conditions and condition Reasons for the TalosConfig object

const (
	// TalosConfigReadyV1Beta2Condition condition is true when all ClientConfigAvailable and DataSecretAvailable conditions are true.
	TalosConfigReadyV1Beta2Condition = capiv1.ReadyV1Beta2Condition

	// TalosConfigNotReadyV1Beta2Reason surfaces when TalosConfig is ready.
	TalosConfigReadyV1Beta2Reason = capiv1.ReadyV1Beta2Reason

	// TalosConfigNotReadyV1Beta2Reason surfaces when TalosConfig is not ready.
	TalosConfigNotReadyV1Beta2Reason = capiv1.NotReadyV1Beta2Reason

	// TalosConfigReadyUnknownV1Beta2Reason surfaces when TalosConfig readiness is unknown.
	TalosConfigReadyUnknownV1Beta2Reason = capiv1.ReadyUnknownV1Beta2Reason
)

const (
	// ClientConfigAvailableV1Beta2Condition documents the status of the client config generation process.
	ClientConfigAvailableV1Beta2Condition = "ClientConfigAvailable"

	// ClientConfigAvailableV1Beta2Reason surfaces when generated Talos client config is available.
	ClientConfigAvailableV1Beta2Reason = capiv1.AvailableV1Beta2Reason

	// ClientConfigAvailableInternalErrorV1Beta2Reason surfaces unexpected failures when reading or generating
	// Talos client config.
	ClientConfigAvailableInternalErrorV1Beta2Reason = capiv1.InternalErrorV1Beta2Reason
)

const (
	// DataSecretAvailableCondition documents the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	DataSecretAvailableV1Beta2Condition = "DataSecretAvailable"

	// DataSecretAvailableV1Beta2Condition surfaces when Talos bootstrap secret is available.
	DataSecretAvailableV1Beta2Reason = capiv1.AvailableV1Beta2Reason

	// DataSecretNotAvailableV1Beta2Condition surfaces when Talos bootstrap is not available.
	DataSecretNotAvailableV1Beta2Reason = capiv1.NotAvailableV1Beta2Reason
)
