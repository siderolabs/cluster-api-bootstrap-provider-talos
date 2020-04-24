package v1alpha2

import (
	bootstrapv1alpha3 "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this TalosConfig to the Hub version (v1alpha3).
func (src *TalosConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1alpha3.TalosConfig)
	if err := Convert_v1alpha2_TalosConfig_To_v1alpha3_TalosConfig(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &bootstrapv1alpha3.TalosConfig{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Status.DataSecretName = restored.Status.DataSecretName

	return nil
}

// ConvertFrom converts from the TalosConfig Hub version (v1alpha3) to this version.
func (dst *TalosConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1alpha3.TalosConfig)
	if err := Convert_v1alpha3_TalosConfig_To_v1alpha2_TalosConfig(src, dst, nil); err != nil {
		return nil
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this TalosConfigList to the Hub version (v1alpha3).
func (src *TalosConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1alpha3.TalosConfigList)
	return Convert_v1alpha2_TalosConfigList_To_v1alpha3_TalosConfigList(src, dst, nil)
}

// ConvertFrom converts from the TalosConfigList Hub version (v1alpha3) to this version.
func (dst *TalosConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1alpha3.TalosConfigList)
	return Convert_v1alpha3_TalosConfigList_To_v1alpha2_TalosConfigList(src, dst, nil)
}

// ConvertTo converts this TalosConfigTemplate to the Hub version (v1alpha3).
func (src *TalosConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1alpha3.TalosConfigTemplate)
	return Convert_v1alpha2_TalosConfigTemplate_To_v1alpha3_TalosConfigTemplate(src, dst, nil)
}

// ConvertFrom converts from the TalosConfigTemplate Hub version (v1alpha3) to this version.
func (dst *TalosConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1alpha3.TalosConfigTemplate)
	return Convert_v1alpha3_TalosConfigTemplate_To_v1alpha2_TalosConfigTemplate(src, dst, nil)
}

// ConvertTo converts this TalosConfigTemplateList to the Hub version (v1alpha3).
func (src *TalosConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1alpha3.TalosConfigTemplateList)
	return Convert_v1alpha2_TalosConfigTemplateList_To_v1alpha3_TalosConfigTemplateList(src, dst, nil)
}

// ConvertFrom converts from the TalosConfigTemplateList Hub version (v1alpha3) to this version.
func (dst *TalosConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1alpha3.TalosConfigTemplateList)
	return Convert_v1alpha3_TalosConfigTemplateList_To_v1alpha2_TalosConfigTemplateList(src, dst, nil)
}

// Convert_v1alpha2_TalosConfigStatus_To_v1alpha3_TalosConfigStatus converts this TalosConfigStatus to the Hub version (v1alpha3).
func Convert_v1alpha2_TalosConfigStatus_To_v1alpha3_TalosConfigStatus(in *TalosConfigStatus, out *bootstrapv1alpha3.TalosConfigStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_TalosConfigStatus_To_v1alpha3_TalosConfigStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Error fields to the Failure fields
	out.FailureMessage = in.ErrorMessage
	out.FailureReason = in.ErrorReason

	return nil
}

// Convert_v1alpha3_TalosConfigStatus_To_v1alpha2_TalosConfigStatus converts from the Hub version (v1alpha3) of the TalosConfigStatus to this version.
func Convert_v1alpha3_TalosConfigStatus_To_v1alpha2_TalosConfigStatus(in *bootstrapv1alpha3.TalosConfigStatus, out *TalosConfigStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_TalosConfigStatus_To_v1alpha2_TalosConfigStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Failure fields to the Error fields
	out.ErrorMessage = in.FailureMessage
	out.ErrorReason = in.FailureReason

	return nil
}

// Convert_v1alpha2_TalosConfigSpec_To_v1alpha3_TalosConfigSpec converts this TalosConfigSpec to the Hub version (v1alpha3).
func Convert_v1alpha2_TalosConfigSpec_To_v1alpha3_TalosConfigSpec(in *TalosConfigSpec, out *bootstrapv1alpha3.TalosConfigSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha2_TalosConfigSpec_To_v1alpha3_TalosConfigSpec(in, out, s)
}

// Convert_v1alpha3_TalosConfigSpec_To_v1alpha2_TalosConfigSpec converts from the Hub version (v1alpha3) of the TalosConfigSpec to this version.
func Convert_v1alpha3_TalosConfigSpec_To_v1alpha2_TalosConfigSpec(in *bootstrapv1alpha3.TalosConfigSpec, out *TalosConfigSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_TalosConfigSpec_To_v1alpha2_TalosConfigSpec(in, out, s)
}
