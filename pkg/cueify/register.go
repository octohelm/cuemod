package cueify

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cueify/core"
	_ "github.com/octohelm/cuemod/pkg/cueify/crd"
	_ "github.com/octohelm/cuemod/pkg/cueify/golang"
)

func Detect(ctx context.Context, src string) (string, map[string]string) {
	return core.SharedExtractors.Detect(ctx, src)
}

func ExtractToDir(ctx context.Context, name string, src string, dist string) error {
	return core.SharedExtractors.ExtractToDir(ctx, name, src, dist)
}
