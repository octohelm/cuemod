// Code generated by go generate. DO NOT EDIT.

//go:generate rm pkg.go
//go:generate go run ../../gen/gen.go

package sha1

import (
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals/core/adt"
	internal "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/internals"
)

func init() {
	internal.Register("crypto/sha1", pkg)
}

var _ = adt.TopKind // in case the adt package isn't used

var pkg = &internal.Package{
	Native: []*internal.Builtin{{
		Name:  "Size",
		Const: "20",
	}, {
		Name:  "BlockSize",
		Const: "64",
	}, {
		Name: "Sum",
		Params: []internal.Param{
			{Kind: adt.BytesKind | adt.StringKind},
		},
		Result: adt.BytesKind | adt.StringKind,
		Func: func(c *internal.CallCtxt) {
			data := c.Bytes(0)
			if c.Do() {
				c.Ret = Sum(data)
			}
		},
	}},
}
