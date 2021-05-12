package cuemod

import (
	"path/filepath"
	"strconv"
	"strings"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/octohelm/cuemod/pkg/cuex/builtin"
)

type OptionFunc = func(c *load.Config)

func OptRoot(dir string) OptionFunc {
	return func(c *load.Config) {
		c.Dir = dir
	}
}

type ImportFunc = func(importPath string, importedAt string) (resolvedDir string, err error)

func OptImportFunc(importFunc ImportFunc) OptionFunc {
	return func(c *load.Config) {
		c.ParseFile = func(filename string, src interface{}) (*cueast.File, error) {
			f, err := parser.ParseFile(filename, src)
			if err != nil {
				return nil, err
			}

			for i := range f.Imports {
				importPath, _ := strconv.Unquote(f.Imports[i].Path.Value)

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

func Build(path string, optionFns ...OptionFunc) *Instance {
	c := &load.Config{}
	for i := range optionFns {
		optionFns[i](c)
	}
	rel, _ := filepath.Rel(c.Dir, path)
	return load.Instances([]string{"./" + rel}, c)[0]
}
