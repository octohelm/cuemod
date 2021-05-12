package kubernetes

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"
	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"

	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestClient(t *testing.T) {
	c, err := NewClient()
	NewWithT(t).Expect(err).To(BeNil())

	allGroupVersionKinds, err := c.AllListableGroupVersionKinds()
	NewWithT(t).Expect(err).To(BeNil())

	t.Run("Release Flow", func(t *testing.T) {
		t.Run("Apply & Delete", func(t *testing.T) {
			data, _ := os.ReadFile("./testdata/web.json")
			var r releasev1alpha1.Release
			_ = json.Unmarshal(data, &r)

			manifests, err := manifest.Process(r, nil)
			NewWithT(t).Expect(err).To(BeNil())

			t.Run("Apply", func(t *testing.T) {
				err := c.ApplyResources(context.Background(), manifests)
				NewWithT(t).Expect(err).To(BeNil())

				list, err := c.ListAll(
					context.Background(),
					allGroupVersionKinds,
					client.InNamespace(r.Namespace),
					client.MatchingLabels(map[string]string{
						manifest.LabelRelease: r.Name,
					}),
				)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(list).To(HaveLen(4)) // Endpoints
			})

			t.Run("Delete", func(t *testing.T) {
				err := c.DeleteResources(context.Background(), manifests)
				NewWithT(t).Expect(err).To(BeNil())

				list, err := c.ListAll(
					context.Background(),
					allGroupVersionKinds,
					client.InNamespace(r.Namespace),
					client.MatchingLabels(map[string]string{
						manifest.LabelRelease: r.Name,
					}),
				)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(list).To(HaveLen(0))
			})
		})

		t.Run("Diff & Prune", func(t *testing.T) {
			data, _ := os.ReadFile("./testdata/web.json")
			var r releasev1alpha1.Release
			_ = json.Unmarshal(data, &r)

			{
				manifests, err := manifest.Process(r, nil)
				NewWithT(t).Expect(err).To(BeNil())

				err = c.ApplyResources(context.Background(), manifests)
				NewWithT(t).Expect(err).To(BeNil())
			}

			t.Run("Diff", func(t *testing.T) {
				manifests, err := manifest.Process(r, nil)
				NewWithT(t).Expect(err).To(BeNil())

				manifestsWithoutService := manifest.ProcessManifests(manifests, func(m manifest.Object) manifest.Object {
					if m.GetObjectKind().GroupVersionKind().Kind == "Service" {
						return nil
					}
					return m
				})

				_, err = c.Diff(context.Background(), manifestsWithoutService)
				NewWithT(t).Expect(err).To(BeNil())
			})

			t.Run("Purge", func(t *testing.T) {
				liveResources, err := c.ListAll(
					context.Background(),
					allGroupVersionKinds,
					client.InNamespace(r.Namespace),
					client.MatchingLabels(map[string]string{
						manifest.LabelRelease: r.Name,
					}),
				)
				NewWithT(t).Expect(err).To(BeNil())

				manifests, err := manifest.Process(r, nil)
				NewWithT(t).Expect(err).To(BeNil())

				manifestsWithoutService := manifest.ProcessManifests(manifests, func(m manifest.Object) manifest.Object {
					if m.GetObjectKind().GroupVersionKind().Kind == "Service" {
						return nil
					}
					return m
				})

				err = c.DeleteResources(context.Background(), liveResources.Orphaned(manifestsWithoutService))
				NewWithT(t).Expect(err).To(BeNil())

				{
					liveResources, err := c.ListAll(
						context.Background(),
						allGroupVersionKinds,
						client.InNamespace(r.Namespace),
						client.MatchingLabels(map[string]string{
							manifest.LabelRelease: r.Name,
						}),
					)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(liveResources).To(HaveLen(2))
				}
			})
		})
	})

	t.Run("client", func(t *testing.T) {
		list := &appsv1.DeploymentList{}
		err := c.List(context.Background(), list, client.InNamespace("default"))
		NewWithT(t).Expect(err).To(BeNil())
	})
}
