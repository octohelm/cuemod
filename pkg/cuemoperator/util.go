package cuemoperator

import (
	"strings"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

func NewReleaseSecret(namespace string, name string, data []byte) *corev1.Secret {
	s := &corev1.Secret{}
	s.Kind = "Secret"
	s.APIVersion = "v1"
	s.Labels = map[string]string{
		SecretTypeRelease: name,
	}

	s.Namespace = namespace
	s.Name = name + ".release.octohelm.tech"

	s.Type = SecretTypeRelease
	s.Data = map[string][]byte{
		"template": data,
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
