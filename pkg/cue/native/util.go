package native

import (
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals/core/runtime"
)

func IsBuiltinPackage(importPath string) bool {
	return runtime.SharedRuntime.IsBuiltinPackage(importPath)
}
