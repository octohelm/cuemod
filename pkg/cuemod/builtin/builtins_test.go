package builtin

import (
	"go/build"
	"os"
	"strings"
	"testing"

	_ "cuelang.org/go/pkg"
)

func TestGen(t *testing.T) {
	pkg, _ := build.Import("cuelang.org/go/pkg", "", build.ImportComment)

	list := make([]string, 0)

	for _, importPath := range pkg.Imports {
		if strings.HasPrefix(importPath, "cuelang.org/go/pkg/") {
			list = append(list, importPath[len("cuelang.org/go/pkg/"):])
		}
	}

	_ = os.WriteFile("builtins.txt", []byte(strings.Join(list, "\n")), os.ModePerm)
}
