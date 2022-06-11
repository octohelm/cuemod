package embedstdlib

import (
	"embed"

	"github.com/octohelm/cuemod/pkg/cuemod/stdlib"
)

var (
	//go:embed std.x.io
	FS embed.FS
)

func init() {
	stdlib.Register(FS, "v0.3.0", "std.x.io")
}
