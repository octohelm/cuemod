package cuemod

import (
	"path/filepath"
	"strconv"
	"strings"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/octohelm/cuemod/pkg/cuemod/builtin"
)

type OptionFunc = func(c *load.Config)

func OptRoot(dir string) OptionFunc {
	return func(c *load.Config) {
		c.Dir = dir
	}
}

func OptOverlay(overlay map[string]load.Source) OptionFunc {
	return func(c *load.Config) {
		c.Overlay = overlay
	}
}

type ImportFunc = func(importPath string, importedAt string) (resolvedDir string, err error)

func OptImportFunc(importFunc ImportFunc) OptionFunc {
	return func(c *load.Config) {
		c.ParseFile = func(filename string, src any) (*cueast.File, error) {
			f, err := parser.ParseFile(filename, src)
			if err != nil {
				return nil, err
			}

			for i := range f.Imports {
				importPath, _ := strconv.Unquote(f.Imports[i].Path.Value)

				// skip abs path and rel path
				if filepath.IsAbs(importPath) {
					continue
				}

				// "xxx/xxxx:xxx"
				importPath = strings.Split(importPath, ":")[0]

				// skip builtin
				if builtin.IsBuiltIn(importPath) {
					continue
				}

				_, err := importFunc(importPath, filename)
				if err != nil {
					return nil, err
				}
			}

			return f, nil
		}
	}
}

type Instance = build.Instance

func BuildConfig(optionFns ...OptionFunc) *load.Config {
	c := &load.Config{}
	for i := range optionFns {
		optionFns[i](c)
	}
	return c
}

func BuildInstances(c *load.Config, inputs []string) []*Instance {
	files := make([]string, len(inputs))

	for i, f := range inputs {
		if filepath.IsAbs(f) {
			rel, _ := filepath.Rel(c.Dir, f)
			files[i] = "./" + rel
		} else {
			files[i] = f
		}
	}

	return load.Instances(files, c)
}

func Build(inputs []string, optionFns ...OptionFunc) *Instance {
	c := BuildConfig(optionFns...)
	// load only support related path
	return BuildInstances(c, inputs)[0]
}
