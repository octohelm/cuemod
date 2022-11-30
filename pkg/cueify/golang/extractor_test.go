package golang_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-courier/logr/slog"

	"github.com/octohelm/cuemod/pkg/cueify"

	"github.com/go-courier/logr"
	"github.com/onsi/gomega"
)

func TestExtractor(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), slog.Logger(slog.Default()))

	_ = os.RemoveAll("./testdata/gen")

	err := cueify.ExtractToDir(
		ctx,
		"go",
		"./testdata/src",
		"./testdata/gen",
	)

	gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
}
