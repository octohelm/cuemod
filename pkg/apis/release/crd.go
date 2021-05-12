package release

import (
	"github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"
	"github.com/octohelm/cuemod/pkg/apiutil"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var CRDs = []*apiextensionsv1.CustomResourceDefinition{
	ReleaseCustomResourceDefinition(),
}

func ReleaseCustomResourceDefinition() *apiextensionsv1.CustomResourceDefinition {
	return apiutil.ToCRD(&apiutil.CustomResourceDefinition{
		GroupVersion: v1alpha1.SchemeGroupVersion,
		KindType:     &v1alpha1.Release{},
		ListKindType: &v1alpha1.ReleaseList{},
	})
}
