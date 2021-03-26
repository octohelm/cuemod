package helm_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/extractor"
	"github.com/onsi/gomega"
)

func TestExtractor(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	_ = os.RemoveAll("./testdata/gen")

	err := extractor.ExtractToDir(
		ctx,
		"crd",
		"./testdata/src",
		"./testdata/gen",
	)

	gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
}
