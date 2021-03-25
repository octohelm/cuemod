package cuemod

import (
	"context"
	"testing"

	"github.com/go-courier/logr"

	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	ctx = WithOpts(ctx, OptVerbose(true))

	m := newCache()

	t.Run("should get go mod", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/grafana/jsonnet-libs/grafana", "master", nil)
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("github.com/grafana/jsonnet-libs/grafana"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/grafana/jsonnet-libs"))
	})

	t.Run("should get sub go mod", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/prometheus/node_exporter/docs/node-mixin", "master", nil)
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("github.com/prometheus/node_exporter/docs/node-mixin"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/prometheus/node_exporter/docs/node-mixin"))
	})
}
