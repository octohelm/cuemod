package builtin

import (
	"cuelang.org/go/cue"
)

var r = cue.NewRuntime()

func IsBuiltIn(importPath string) bool {
	return r.IsBuiltinPackage(importPath)
}
