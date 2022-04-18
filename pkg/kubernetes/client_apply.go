package kubernetes

import (
	"context"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *KubeClient) ApplyResources(ctx context.Context, list manifest.List) error {
	for i := range list {
		if err := ApplyOne(ctx, c.Client, list[i], false); err != nil {
			return err
		}
	}
	return nil
}

const FieldOwner = client.FieldOwner("cuemod")

func ApplyOne(ctx context.Context, c client.Client, obj manifest.Object, dryRun bool) error {
	kubeC := c
	if dryRun {
		kubeC = client.NewDryRunClient(c)
	}

	if err := kubeC.Patch(ctx, obj, client.Apply, FieldOwner); err != nil {
		return err
	}

	if !dryRun {
		manifest.Annotated(obj, map[string]string{
			manifest.AnnotationReleaseStage: "patched",
		})
	}

	return nil
}
