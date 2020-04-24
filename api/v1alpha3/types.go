package v1alpha3

// TalosConfigTemplateResource defines the Template structure
type TalosConfigTemplateResource struct {
	Spec TalosConfigSpec `json:"spec,omitempty"`
}
