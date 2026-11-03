// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import (
	"strconv"
	"testing"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/utils/ptr"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/randfill"

	bsv1beta1 "github.com/siderolabs/cluster-api-bootstrap-provider-talos/api/v1beta1"
)

func TestFuzzyConversion(t *testing.T) {
	t.Run("for TalosConfig", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &bsv1beta1.TalosConfig{},
		Spoke:       &TalosConfig{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{talosConfigFuzzFuncs},
	}))

	t.Run("for TalosConfigTemplate", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &bsv1beta1.TalosConfigTemplate{},
		Spoke:       &TalosConfigTemplate{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{talosConfigTemplateFuzzFuncs},
	}))
}

func talosConfigFuzzFuncs(_ runtimeserializer.CodecFactory) []any {
	return []any{
		hubConfigPatches,
		spokeConfigPatches,
		hubTalosConfigStatus,
		spokeTalosConfigStatus,
	}
}

func talosConfigTemplateFuzzFuncs(_ runtimeserializer.CodecFactory) []any {
	return []any{
		hubConfigPatches,
		spokeConfigPatches,
	}
}

// Patch values are raw JSON, so they have to be valid JSON to survive a round trip.
func fuzzJSON(c randfill.Continue) apiextensions.JSON {
	return apiextensions.JSON{Raw: []byte(strconv.Quote(c.String(0)))}
}

func hubConfigPatches(in *bsv1beta1.ConfigPatches, c randfill.Continue) {
	c.FillNoCustom(in)

	in.Value = fuzzJSON(c)
}

func spokeConfigPatches(in *ConfigPatches, c randfill.Continue) {
	c.FillNoCustom(in)

	in.Value = fuzzJSON(c)
}

func hubTalosConfigStatus(in *bsv1beta1.TalosConfigStatus, c randfill.Continue) {
	c.FillNoCustom(in)

	// Up-conversion always creates the deprecated v1beta1 status struct.
	if in.Deprecated == nil {
		in.Deprecated = &bsv1beta1.TalosConfigDeprecatedStatus{}
	}

	if in.Deprecated.V1Beta1 == nil {
		in.Deprecated.V1Beta1 = &bsv1beta1.TalosConfigV1Beta1DeprecatedStatus{}
	}
}

func spokeTalosConfigStatus(in *TalosConfigStatus, c randfill.Continue) {
	c.FillNoCustom(in)

	// TalosConfig has no equivalent in the hub version and is dropped on up-conversion.
	in.TalosConfig = "" //nolint:staticcheck // deprecated, kept for conversion only

	// DataSecretName is a string in the hub version, so an empty value is down-converted to nil.
	if ptr.Deref(in.DataSecretName, "") == "" {
		in.DataSecretName = nil
	}

	// An empty v1beta2 status is dropped on up-conversion.
	if in.V1Beta2 != nil && len(in.V1Beta2.Conditions) == 0 {
		in.V1Beta2 = nil
	}
}
