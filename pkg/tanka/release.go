package tanka

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	releasev1alpha "github.com/octohelm/cuemod/pkg/api/release/v1alpha"
)

func ProcessResources(r *releasev1alpha.Release, exprs process.Matchers) (manifest.List, error) {
	raw := r.Spec

	if raw == nil {
		return manifest.List{}, nil
	}

	// Scan for everything that looks like a Kubernetes object
	extracted, err := process.Extract(raw)
	if err != nil {
		return nil, err
	}

	// Unwrap *List types
	if err := process.Unwrap(extracted); err != nil {
		return nil, err
	}

	out := make(manifest.List, 0, len(extracted))
	for _, m := range extracted {
		out = append(out, m)
	}

	// set default namespace
	out = process.Namespace(out, r.Metadata.Namespace)

	// tanka.dev/** labels
	out = InjectLabels(out, r)

	// Perhaps filter for kind/name expressions
	if len(exprs) > 0 {
		out = process.Filter(out, exprs)
	}

	// Best-effort dependency sort
	process.Sort(out)

	return out, nil
}

func InjectLabels(list manifest.List, r *releasev1alpha.Release) manifest.List {
	for i, m := range list {
		if m.Kind() == "Namespace" {
			// don't inject label to namespace
			continue
		}

		m.Metadata().Labels()[process.LabelEnvironment] = r.Metadata.Name
		list[i] = m
	}
	return list
}
