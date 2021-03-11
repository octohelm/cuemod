package golang

import (
	"context"

	cueast "cuelang.org/go/cue/ast"
)

type importPaths map[string]string

func (i importPaths) toImportDecl() *cueast.ImportDecl {
	importDecl := &cueast.ImportDecl{}

	for importPath := range i {
		importDecl.Specs = append(importDecl.Specs, cueast.NewImport(cueast.NewIdent(i[importPath]), importPath))
	}

	return importDecl
}

func (i importPaths) add(importPath string, pkgName string) {
	i[importPath] = pkgName
}

type contextImportPaths struct {
}

func importPathsFromContext(ctx context.Context) importPaths {
	if i, ok := ctx.Value(contextImportPaths{}).(importPaths); ok {
		return i
	}
	return importPaths{}
}

func withImportPaths(ctx context.Context, i importPaths) context.Context {
	return context.WithValue(ctx, contextImportPaths{}, i)
}
