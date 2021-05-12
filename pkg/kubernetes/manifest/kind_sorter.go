package manifest

import (
	"sort"

	"helm.sh/helm/v3/pkg/releaseutil"
)

type KindSortOrder = releaseutil.KindSortOrder

var (
	InstallOrder   = releaseutil.InstallOrder
	UninstallOrder = releaseutil.UninstallOrder
)

func SortByKind(manifests []Object, ordering KindSortOrder) []Object {
	ks := newKindSorter(manifests, ordering)
	sort.Stable(ks)
	return ks.manifests
}

type kindSorter struct {
	ordering  map[string]int
	manifests []Object
}

func newKindSorter(m []Object, s KindSortOrder) *kindSorter {
	o := make(map[string]int, len(s))
	for v, k := range s {
		o[k] = v
	}

	return &kindSorter{
		manifests: m,
		ordering:  o,
	}
}

func (k *kindSorter) Len() int { return len(k.manifests) }

func (k *kindSorter) Swap(i, j int) { k.manifests[i], k.manifests[j] = k.manifests[j], k.manifests[i] }

func (k *kindSorter) Less(i, j int) bool {
	a := k.manifests[i]
	b := k.manifests[j]
	first, aok := k.ordering[a.GetObjectKind().GroupVersionKind().Kind]
	second, bok := k.ordering[b.GetObjectKind().GroupVersionKind().Kind]

	if !aok && !bok {
		// if both are unknown then sort alphabetically by kind, keep original order if same kind
		if a.GetObjectKind().GroupVersionKind().Kind != b.GetObjectKind().GroupVersionKind().Kind {
			return a.GetObjectKind().GroupVersionKind().Kind < b.GetObjectKind().GroupVersionKind().Kind
		}
		return first < second
	}
	// unknown kind is last
	if !aok {
		return false
	}
	if !bok {
		return true
	}
	// sort different kinds, keep original order if same priority
	return first < second
}
