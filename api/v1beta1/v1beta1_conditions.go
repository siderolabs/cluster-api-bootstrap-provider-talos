// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1beta1

import capiv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

// Deprecated CAPI V1beta1 contract TalosConfig conditions. To be removed in a future version

const (
	// DataSecretAvailableCondition documents the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	DataSecretAvailableV1Beta1Condition capiv1.ConditionType = "DataSecretAvailable"

	// DataSecretAvailableReason documents a succeeded bootstrap secret generation process.
	//
	// NOTE: This is the only reason for DataSecretAvailableCondition becoming true.
	DataSecretAvailableV1Beta1Reason = "DataSecretAvailable"

	// WaitingForClusterInfrastructureReason (Severity=Info) document a bootstrap secret generation process
	// waiting for the cluster infrastructure to be ready.
	//
	// NOTE: Having the cluster infrastructure ready is a pre-condition for starting to create machines;
	// the TalosConfig controller ensure this pre-condition is satisfied.
	WaitingForClusterInfrastructureV1Beta1Reason = "WaitingForClusterInfrastructure"

	// DataSecretGenerationFailedReason (Severity=Warning) documents a TalosConfig controller detecting
	// an error while generating a data secret; those kind of errors are usually due to misconfigurations
	// and user intervention is required to get them fixed.
	DataSecretGenerationFailedV1Beta1Reason = "DataSecretGenerationFailed"
)

const (
	// ClientConfigAvailableCondition documents the status of the client config generation process.
	ClientConfigAvailableV1Beta1Condition capiv1.ConditionType = "ClientConfigAvailable"

	// ClientConfigGenerationFailedReason (Severity=Warning) documents a TalosConfig controller detecting
	// an error while generating a client config; those kind of errors are usually due to misconfigurations
	// and user intervention is required to get them fixed.
	ClientConfigGenerationFailedV1Beta1Reason = "ClientConfigGenerationFailed"
)
