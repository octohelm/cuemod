package extractor

import (
	"context"

	"github.com/octohelm/cuemod/pkg/extractor/core"

	_ "github.com/octohelm/cuemod/pkg/extractor/golang"
	_ "github.com/octohelm/cuemod/pkg/extractor/helm"
	_ "github.com/octohelm/cuemod/pkg/extractor/jsonnet"
)

func Detect(ctx context.Context, src string) (string, map[string]string) {
	return core.SharedExtractors.Detect(ctx, src)
}

func ExtractToDir(ctx context.Context, name string, src string, gen string) error {
	return core.SharedExtractors.ExtractToDir(ctx, name, src, gen)
}
