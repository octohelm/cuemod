package kubernetes

import (
	"bytes"
	"context"
	"fmt"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"github.com/pmezard/go-difflib/difflib"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func (c *KubeClient) Diff(ctx context.Context, objs []manifest.Object) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for i := range objs {
		obj := objs[i]

		gvk := obj.GetObjectKind().GroupVersionKind()

		live, err := manifest.NewForGroupVersionKind(gvk)
		if err != nil {
			return nil, err
		}

		toCreated := false

		if err := c.Get(ctx, client.ObjectKeyFromObject(obj), live); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, err
			} else {
				toCreated = true
			}
		}

		if err := ApplyOne(ctx, c.Client, obj, true); err != nil {
			return nil, err
		}

		// ignore managedFields
		live.SetManagedFields(nil)
		obj.SetManagedFields(nil)

		liveManifest, err := yaml.Marshal(live)
		if err != nil {
			return nil, err
		}

		if toCreated {
			// created
			liveManifest = nil
		}

		mergedManifest, err := yaml.Marshal(obj)
		if err != nil {
			return nil, err
		}

		_ = difflib.WriteUnifiedDiff(buf, difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(liveManifest)),
			B:        difflib.SplitLines(string(mergedManifest)),
			FromFile: fmt.Sprintf("%s LIVE", manifest.Identity(live)),
			ToFile:   fmt.Sprintf("%s MERGED", manifest.Identity(obj)),
		})
	}

	return buf.Bytes(), nil
}
