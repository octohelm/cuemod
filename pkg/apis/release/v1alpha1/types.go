package v1alpha1

import (
	"bytes"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	SchemeBuilder.Register(&Release{}, &ReleaseList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReleaseSpec   `json:"spec,omitempty"`
	Status ReleaseStatus `json:"status,omitempty"`
}

type ReleaseTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// cue codes witch could return nested kube resources
	Data string `json:"data,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Release `json:"items"`
}

type ReleaseSpec map[string]interface{}

func (s ReleaseSpec) DeepCopy() ReleaseSpec {
	spec := ReleaseSpec{}

	buf := bytes.NewBuffer(nil)
	_ = json.NewEncoder(buf).Encode(s)
	_ = json.NewDecoder(buf).Decode(&spec)

	return spec
}

type ReleaseStatus struct {
	Resources []Resource `json:"resources,omitempty"`
}

type Resource struct {
	schema.GroupVersionKind
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
