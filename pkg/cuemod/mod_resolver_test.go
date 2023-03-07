package cuemod

import (
	"context"
	"testing"

	"github.com/go-courier/logr/slog"

	"github.com/octohelm/cuemod/pkg/cuemod/modfile"

	"github.com/go-courier/logr"

	. "github.com/onsi/gomega"

	_ "github.com/octohelm/cuemod/pkg/cuemod/testdata/embedstdlib"
)

func TestModResolver(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), slog.Logger(slog.Default()))

	ctx = WithOpts(ctx, OptVerbose(true))

	m := newModResolver()

	t.Run("should resolve stdlib", func(t *testing.T) {
		mod, err := m.Get(ctx, "std.x.io/a", modfile.ModVersion{VcsRef: "main"})
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(mod.Module).To(Equal("std.x.io"))
		NewWithT(t).Expect(mod.Repo).To(Equal("std.x.io"))
	})
}
