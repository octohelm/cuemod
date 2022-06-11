package cuemod

import (
	"context"
	"testing"

	"github.com/octohelm/cuemod/pkg/cuemod/modfile"

	"github.com/go-courier/logr"

	. "github.com/onsi/gomega"

	_ "github.com/octohelm/cuemod/pkg/cuemod/testdata/embedstdlib"
)

func TestModResolver(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	ctx = WithOpts(ctx, OptVerbose(true))

	m := newModResolver()

	t.Run("should resolve stdlib", func(t *testing.T) {
		mod, err := m.Get(ctx, "std.x.io/a", modfile.ModVersion{VcsRef: "main"})
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(mod.Module).To(Equal("std.x.io"))
		NewWithT(t).Expect(mod.Repo).To(Equal("std.x.io"))
	})

	t.Run("should get sub when go.mod exists", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/open-telemetry/opentelemetry-go/exporters/prometheus", modfile.ModVersion{VcsRef: "main"})
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("github.com/open-telemetry/opentelemetry-go/exporters/prometheus"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/open-telemetry/opentelemetry-go/exporters/prometheus"))
	})

	t.Run("should delegate cue.mod/module.cue", func(t *testing.T) {
		mod, err := m.Get(ctx, "github.com/dagger/dagger/pkg/dagger.io", modfile.ModVersion{Version: "v0.2.10"})
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(mod.Module).To(Equal("dagger.io"))
		NewWithT(t).Expect(mod.Repo).To(Equal("github.com/dagger/dagger"))
	})
}
