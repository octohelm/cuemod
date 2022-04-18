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

				if IsFileImport(importPath) {
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

func Build(path string, optionFns ...OptionFunc) *Instance {
	c := &load.Config{}
	for i := range optionFns {
		optionFns[i](c)
	}
	rel, _ := filepath.Rel(c.Dir, path)
	return load.Instances([]string{"./" + rel}, c)[0]
}

func IsFileImport(p string) bool {
	return strings.HasPrefix(p, "file/")
}

func FixFileImport(p string) string {
	if filepath.IsAbs(p) {
		return "file" + p + ":cue"
	}
	return p
}
