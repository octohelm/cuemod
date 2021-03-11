package cue

import (
	"context"
	"strconv"
	"strings"

	"cuelang.org/go/cue/ast"

	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue/load"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue/parser"
	"github.com/octohelm/cuemod/pkg/cue/native"

	_ "github.com/octohelm/cuemod/pkg/extension"
)

type RuntimeOpts struct {
	ModuleRoot string
	Importer   Importer
}

type OptionFunc = func(c *RuntimeOpts)

func OptModuleRoot(moduleRoot string) OptionFunc {
	return func(c *RuntimeOpts) {
		c.ModuleRoot = moduleRoot
	}
}

func OptImporter(importer Importer) OptionFunc {
	return func(c *RuntimeOpts) {
		c.Importer = importer
	}
}

func NewRuntime(optionFns ...OptionFunc) *Runtime {
	r := &Runtime{}

	for i := range optionFns {
		optionFns[i](&r.Opts)
	}

	return r
}

type Runtime struct {
	Opts RuntimeOpts
}

func (r *Runtime) Eval(ctx context.Context, args []string, encoding Encoding) ([][]byte, error) {
	instances := r.build(ctx, args)

	results := make([][]byte, len(instances))

	for i := range results {
		v := instances[i].Value()

		if err := v.Validate(cue.Final()); err != nil {
			return nil, err
		}

		data, err := encode(v, encoding)
		if err != nil {
			return nil, err
		}
		results[i] = []byte(data)
	}

	return results, nil
}

func (r *Runtime) build(ctx context.Context, args []string) []*cue.Instance {
	c := &load.Config{}

	if v := r.Opts.ModuleRoot; v != "" {
		c.ModuleRoot = v
	}

	if importer := r.Opts.Importer; importer != nil {
		// todo ugly hack, mv to loaderFunc if possible
		c.ParseFile = func(filename string, src interface{}) (*ast.File, error) {
			f, err := parser.ParseFile(filename, src)
			if err != nil {
				return nil, err
			}

			for i := range f.Imports {
				importPath, _ := strconv.Unquote(f.Imports[i].Path.Value)

				importPath = strings.Split(importPath, ":")[0]

				// skip builtin
				if native.IsBuiltinPackage(importPath) {
					continue
				}

				_, err := importer.Import(ctx, importPath, filename)
				if err != nil {
					return nil, err
				}
			}

			return f, err
		}
	}

	return cue.Build(load.Instances(args, c))
}
