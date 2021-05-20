package cuemoperator

import (
	"strings"

	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
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

func allWatchableGroupVersionKinds(conf *rest.Config) (gvks []schema.GroupVersionKind, err error) {
	dc, e := discovery.NewDiscoveryClientForConfig(conf)
	if e != nil {
		return nil, e
	}

	preferredResources, err := dc.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	for i := range preferredResources {
		pr := preferredResources[i]

		if len(pr.APIResources) == 0 {
			continue
		}

		gv, err := schema.ParseGroupVersion(pr.GroupVersion)
		if err != nil {
			continue
		}

		for _, resource := range pr.APIResources {
			if len(resource.Verbs) == 0 {
				continue
			}

			if strings.Contains(strings.Join(resource.Verbs, ","), "watch") {
				gvks = append(gvks, gv.WithKind(resource.Kind))
			}
		}
	}
	return
}
