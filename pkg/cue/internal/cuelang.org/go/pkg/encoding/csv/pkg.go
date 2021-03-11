// Code generated by go generate. DO NOT EDIT.

//go:generate rm pkg.go
//go:generate go run ../../gen/gen.go

package csv

import (
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals/core/adt"
	internal "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/internals"
)

func init() {
	internal.Register("encoding/csv", pkg)
}

var _ = adt.TopKind // in case the adt package isn't used

var pkg = &internal.Package{
	Native: []*internal.Builtin{{
		Name: "Encode",
		Params: []internal.Param{
			{Kind: adt.TopKind},
		},
		Result: adt.StringKind,
		Func: func(c *internal.CallCtxt) {
			x := c.Value(0)
			if c.Do() {
				c.Ret, c.Err = Encode(x)
			}
		},
	}, {
		Name: "Decode",
		Params: []internal.Param{
			{Kind: adt.BytesKind | adt.StringKind},
		},
		Result: adt.ListKind,
		Func: func(c *internal.CallCtxt) {
			r := c.Reader(0)
			if c.Do() {
				c.Ret, c.Err = Decode(r)
			}
		},
	}},
}
