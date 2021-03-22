package v1alpha

type Release struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   Metadata               `json:"metadata"`
	Spec       map[string]interface{} `json:"spec"`
}

type Metadata struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}
