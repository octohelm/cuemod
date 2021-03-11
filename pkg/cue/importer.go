package cue

import "context"

type Importer interface {
	Import(ctx context.Context, importPath string, importedAt string) (resolvedDir string, err error)
}
