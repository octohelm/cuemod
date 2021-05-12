package manifest

import (
	releasev1alpha "github.com/octohelm/cuemod/pkg/api/release/v1alpha"
)

func Process(cfg releasev1alpha.Release, exprs Matchers) (List, error) {
	raw := cfg.Spec

	if raw == nil {
		return List{}, nil
	}

	// Scan for everything that looks like a Kubernetes object
	extracted, err := Extract(raw)
	if err != nil {
		return nil, err
	}

	out := make(List, 0)

	for _, m := range extracted {
		out = append(out, m)
	}

	out = ProcessManifests(
		out,
		WithNamespace(cfg.Metadata.Namespace),
		WithReleaseName(cfg.Metadata.Name),
	)

	// Perhaps filter for kind/name expressions
	if len(exprs) > 0 {
		out = Filter(out, exprs)
	}

	SortByKind(out, InstallOrder)

	return out, nil
}

type ProcessFunc = func(m Object) Object

func ProcessManifests(list List, fns ...ProcessFunc) (processed List) {
	for i := range list {
		m := list[i]
		for j := range fns {
			if fn := fns[j]; fn != nil {
				m = fn(m)
			}
		}
		if m != nil {
			processed = append(processed, m)
		}
	}
	return
}
