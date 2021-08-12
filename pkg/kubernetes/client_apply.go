package kubernetes

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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
	gvk := obj.GetObjectKind().GroupVersionKind()

	live, err := manifest.NewForGroupVersionKind(gvk)
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

	if err := kubeC.Patch(ctx, obj, PatchFor(gvk, live)); err != nil {
		return err
	}
	if !dryRun {
		manifest.Annotated(obj, map[string]string{
			manifest.AnnotationReleaseStage: "patched",
		})
	}
	return nil
}

func PatchFor(gvk schema.GroupVersionKind, live client.Object) client.Patch {
	if _, ok := live.(*unstructured.Unstructured); ok {
		return client.MergeFromWithOptions(live, client.MergeFromWithOptimisticLock{})
	}

	if gvk.Group == corev1.GroupName && gvk.Kind == "Service" {
		return client.Merge
	}

	// TODO handle more
	if gvk.Group == corev1.GroupName || gvk.Group == appsv1.GroupName {
		return client.StrategicMergeFrom(live)
	}

	return client.MergeFromWithOptions(live, client.MergeFromWithOptimisticLock{})
}
