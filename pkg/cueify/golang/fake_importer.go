package golang

import (
	"go/importer"
	"go/types"
	"regexp"
	"strings"
)

func newFakeImporter() types.Importer {
	return &fakeImporter{Importer: importer.Default()}
}

type fakeImporter struct {
	Importer types.Importer
}

func (f *fakeImporter) Import(importPath string) (*types.Package, error) {
	if pkg, err := f.Importer.Import(importPath); err == nil {
		return pkg, nil
	}

	pkg := types.NewPackage(importPath, pkgNameFromImportPath(importPath))
	pkg.MarkComplete()
	return pkg, nil
}

var reV = regexp.MustCompile("v([0-9]+)$")

// pkg name from import path
func pkgNameFromImportPath(importPath string) string {
	parts := strings.Split(importPath, "/")
	name := parts[len(parts)-1]

	// xxx/v2
	if len(parts) > 1 && reV.MatchString(name) && !(name == "v1" || name == "v0") {
		name = parts[len(parts)-2]
	}

	// like yaml.v2
	if names := strings.Split(name, "."); len(names) == 2 {
		if reV.MatchString(names[1]) {
			name = names[0]
		}
	}

	return name
}
