package kubernetes

import (
	"context"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

func ApplyOne(ctx context.Context, c client.Client, obj manifest.Object, dryRun bool) error {
	live, err := manifest.NewForGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	if err != nil {
		return err
	}

	kubeC := c

	if dryRun {
		kubeC = client.NewDryRunClient(c)
	}

	if err := kubeC.Get(ctx, client.ObjectKeyFromObject(obj), live); err != nil {
		if apierrors.IsNotFound(err) {
			if err := kubeC.Create(ctx, obj); err != nil {
				return err
			}
			if !dryRun {
				manifest.Annotated(obj, map[string]string{
					manifest.AnnotationReleaseStage: "created",
				})
			}
			return nil
		}
		return err
	}

	if err := kubeC.Patch(ctx, obj, client.Merge); err != nil {
		return err
	}
	if !dryRun {
		manifest.Annotated(obj, map[string]string{
			manifest.AnnotationReleaseStage: "patched",
		})
	}
	return nil
}
