// Code generated by go generate. DO NOT EDIT.

//go:generate rm pkg.go
//go:generate go run ../../gen/gen.go

package file

import (
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals/core/adt"
	internal "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/internals"
)

func init() {
	internal.Register("tool/file", pkg)
}

var _ = adt.TopKind // in case the adt package isn't used

var pkg = &internal.Package{
	Native: []*internal.Builtin{},
	CUE: `{
	Read: {
		$id:      "tool/file.Read"
		filename: !=""
		contents: *bytes | string
	}
	Append: {
		$id:         "tool/file.Append"
		filename:    !=""
		permissions: int | *438
		contents:    bytes | string
	}
	Create: {
		$id:         "tool/file.Create"
		filename:    !=""
		permissions: int | *438
		contents:    bytes | string
	}
	Glob: {
		$id:  "tool/file.Glob"
		glob: !=""
		files: [...string]
	}
}`,
}
