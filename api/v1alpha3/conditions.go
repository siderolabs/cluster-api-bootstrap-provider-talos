// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

// Conditions and condition Reasons for the TalosConfig object

const (
	// DataSecretAvailableCondition documents the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	DataSecretAvailableCondition string = "DataSecretAvailable"

	// WaitingForClusterInfrastructureReason (Severity=Info) document a bootstrap secret generation process
	// waiting for the cluster infrastructure to be ready.
	//
	// NOTE: Having the cluster infrastructure ready is a pre-condition for starting to create machines;
	// the TalosConfig controller ensure this pre-condition is satisfied.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"

	// DataSecretGenerationFailedReason (Severity=Warning) documents a TalosConfig controller detecting
	// an error while generating a data secret; those kind of errors are usually due to misconfigurations
	// and user intervention is required to get them fixed.
	DataSecretGenerationFailedReason = "DataSecretGenerationFailed"
)

const (
	// ClientConfigAvailableCondition documents the status of the client config generation process.
	ClientConfigAvailableCondition string = "ClientConfigAvailable"

	// ClientConfigGenerationFailedReason (Severity=Warning) documents a TalosConfig controller detecting
	// an error while generating a client config; those kind of errors are usually due to misconfigurations
	// and user intervention is required to get them fixed.
	ClientConfigGenerationFailedReason = "ClientConfigGenerationFailed"
)
