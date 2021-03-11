package golang

import (
	"context"
	"go/build"
	"os"
	"path/filepath"

	"github.com/octohelm/cuemod/pkg/extractor/core"

	cueast "cuelang.org/go/cue/ast"
)

func init() {
	core.Register(&Extractor{})
}

// Extractor similar to cue go, but just only generate for one import path
//
// Targets:
//
// * gen const values
// * gen types
//   * k8s resources with meta_v1.TypeMeta should gen with { apiVersion, kind }
//
// Rules:
//
// * skip & drop imports from go std libs exclude cue builtins support.
type Extractor struct {
}

func (Extractor) Name() string {
	return "go"
}

func (Extractor) Detect(ctx context.Context, src string) (bool, map[string]string) {
	f, err := os.Lstat(filepath.Join(src, "go.mod"))
	if err == nil {
		return !f.IsDir(), nil
	}
	return false, nil
}

func (e *Extractor) Extract(ctx context.Context, src string) ([]*cueast.File, error) {
	pkg, err := build.ImportDir(src, build.IgnoreVendor)
	if err != nil {
		return nil, err
	}
	return (&pkgExtractor{Package: pkg}).Extract(ctx)
}
