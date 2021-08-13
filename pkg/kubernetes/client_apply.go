package kubernetes

import (
	"context"

	"github.com/stretchr/objx"
	"gopkg.in/square/go-jose.v2/json"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	if err := kubeC.Patch(ctx, obj, AutoMergePath(gvk, live)); err != nil {
		return err
	}
	if !dryRun {
		manifest.Annotated(obj, map[string]string{
			manifest.AnnotationReleaseStage: "patched",
		})
	}
	return nil
}

func AutoMergePath(gvk schema.GroupVersionKind, live client.Object) client.Patch {
	return &autoMergePath{
		gvk:  gvk,
		live: live,
	}
}

type autoMergePath struct {
	gvk  schema.GroupVersionKind
	live client.Object
}

func (a *autoMergePath) CanMergeStrategic() bool {
	if a.gvk.Group == corev1.GroupName || (a.gvk.Kind == "Service" || a.gvk.Kind == "PersistentVolumeClaim") {
		return false
	}

	if a.gvk.Group == corev1.GroupName || a.gvk.Group == appsv1.GroupName {
		return true
	}

	return false
}

func (a *autoMergePath) Type() types.PatchType {
	if a.CanMergeStrategic() {
		return types.StrategicMergePatchType
	}
	return types.MergePatchType
}

var deletableMaps = []string{
	"metadata.labels",
	"metadata.annotations",

	// QService
	"spec.envs",
}

func (a *autoMergePath) Data(obj client.Object) ([]byte, error) {
	if a.CanMergeStrategic() {
		if _, ok := obj.(*unstructured.Unstructured); ok {
			return client.MergeFromWithOptions(a.live, client.MergeFromWithOptimisticLock{}).Data(obj)
		}
		return client.StrategicMergeFrom(a.live, client.MergeFromWithOptimisticLock{}).Data(obj)
	}

	liveData, _ := json.Marshal(a.live)
	objData, _ := json.Marshal(obj)

	live := &unstructured.Unstructured{}
	merged := &unstructured.Unstructured{}

	_, _, _ = unstructured.UnstructuredJSONScheme.Decode(liveData, &a.gvk, live)
	_, _, _ = unstructured.UnstructuredJSONScheme.Decode(objData, &a.gvk, merged)

	merged.SetResourceVersion(live.GetResourceVersion())
	merged.SetUID(live.GetUID())

	l := objx.Map(live.Object)
	m := objx.Map(merged.Object)

	for _, path := range deletableMaps {
		if cur := l.Get(path); cur.IsObjxMap() {
			if target := m.Get(path); target.IsObjxMap() {
				c := cur.MustObjxMap()
				t := target.MustObjxMap()

				for k := range c {
					if _, ok := t[k]; !ok {
						// set nil for delete
						t[k] = nil
					}
				}
			}
		}
	}

	return json.Marshal(merged)
}
