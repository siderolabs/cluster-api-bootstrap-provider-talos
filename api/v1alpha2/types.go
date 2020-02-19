package v1alpha2

// TalosConfigTemplateResource defines the Template structure
type TalosConfigTemplateResource struct {
	Spec TalosConfigSpec `json:"spec,omitempty"`
}
