package helm_test

import (
	"context"
	"testing"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/extractor"
	"github.com/onsi/gomega"
)

func TestExtractor(t *testing.T) {
	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	err := extractor.ExtractToDir(
		ctx,
		"helm",
		"./testdata/src",
		"./testdata/gen",
	)

	gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
}
