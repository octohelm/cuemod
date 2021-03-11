package jsonnet

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"

	cueast "cuelang.org/go/cue/ast"
	"github.com/octohelm/cuemod/pkg/extractor/core"
)

func init() {
	core.Register(&Extractor{})
}

// Extractor from jsonnet
type Extractor struct {
}

func (Extractor) Name() string {
	return "jsonnet"
}

func (Extractor) Detect(ctx context.Context, src string) (bool, map[string]string) {
	data, err := os.ReadFile(filepath.Join(src, "jsonnetfile.json"))
	if err == nil {
		deps := map[string]string{}

		jf := &spec.JsonnetFile{}

		if err := jf.UnmarshalJSON(data); err == nil {
			for _, d := range jf.Dependencies {
				if d.Source.GitSource != nil {
					repo := filepath.Join(d.Source.GitSource.Host, d.Source.GitSource.User, d.Source.GitSource.Repo)
					deps[repo] = d.Version
				}
			}
		}

		return true, deps
	}

	jsonnetfiles, err := filepath.Glob(filepath.Join(src, "*.*sonnet"))
	if err == nil {
		return len(jsonnetfiles) > 0, nil
	}

	return false, nil
}

func (e *Extractor) Extract(ctx context.Context, src string) (files []*cueast.File, err error) {
	l, err := Load(ctx, src)
	if err != nil {
		return nil, err
	}
	return l.Extract()
}
