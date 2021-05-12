package kubernetes

import (
	"context"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
)

func (c *KubeClient) DeleteResources(ctx context.Context, list manifest.List) error {
	finalList := manifest.ProcessManifests(list, manifest.IgnoreNamespace())

	manifest.SortByKind(finalList, manifest.UninstallOrder)

	for i := range finalList {
		m := finalList[i]
		if m == nil {
			continue
		}

		if err := c.Delete(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
