package cuemod

import (
	"context"
	"testing"

	"github.com/go-courier/logr"

	. "github.com/onsi/gomega"
)

func TestModResolver(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	ctx = WithOpts(ctx, OptVerbose(true))

	m := newModResolver()

	t.Run("should get go mod", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/grafana/jsonnet-libs/grafana", "master", nil)
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("github.com/grafana/jsonnet-libs/grafana"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/grafana/jsonnet-libs"))
	})

	t.Run("should get sub go mod", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/open-telemetry/opentelemetry-go/exporters/prometheus", "main", nil)
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("github.com/open-telemetry/opentelemetry-go/exporters/prometheus"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/open-telemetry/opentelemetry-go/exporters/prometheus"))
	})

	t.Run("should delegate cue.mod/module.cue", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/dagger/dagger/pkg/dagger.io", "v0.2.10", nil)
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("dagger.io"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/dagger/dagger"))
	})
}
