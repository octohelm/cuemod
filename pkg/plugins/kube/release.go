package kube

import (
	releasev1alpha "github.com/octohelm/cuemod/pkg/api/release/v1alpha"
	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
)

func ProcessResources(r *releasev1alpha.Release, exprs manifest.Matchers) (manifest.List, error) {
	raw := r.Spec

	if raw == nil {
		return manifest.List{}, nil
	}

	extracted, err := manifest.Extract(raw)
	if err != nil {
		return nil, err
	}

	out := make(manifest.List, 0, len(extracted))
	for _, m := range extracted {
		out = append(out, m)
	}

	out = manifest.ProcessManifests(out,
		manifest.WithReleaseName(r.Metadata.Name),
	)

	if len(exprs) > 0 {
		out = manifest.Filter(out, exprs)
	}

	manifest.SortByKind(out, manifest.InstallOrder)

	return out, nil
}
