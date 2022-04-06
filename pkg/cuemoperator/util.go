package cuemoperator

import (
	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"
)

func NewReleaseTemplate(namespace string, name string, template []byte) *releasev1alpha1.Release {
	s := &releasev1alpha1.Release{}

	s.SetGroupVersionKind(releasev1alpha1.SchemeGroupVersion.WithKind("Release"))

	s.Namespace = namespace
	s.Name = name

	t := &releasev1alpha1.ReleaseTemplate{
		Template: string(template),
	}

	t.SetGroupVersionKind(releasev1alpha1.SchemeGroupVersion.WithKind("ReleaseTemplate"))

	s.Spec = map[string]interface{}{
		"template": t,
	}

	return s
}
