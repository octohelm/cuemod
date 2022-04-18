package cmd

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cli"
	"github.com/octohelm/cuemod/pkg/cuemod"
)

func init() {
	app.Add(&Get{})
}

type Get struct {
	cli.Name `args:"BASE_DIR" desc:"download dependencies"`
	Upgrade  bool `flag:"upgrade,u" desc:"upgrade dependencies"`
}

func (o *Get) Run(ctx context.Context, args []string) error {
	importPath := "."
	if len(args) > 0 {
		importPath = args[0]
	}
	return cuemod.FromContext(ctx).Get(cuemod.WithOpts(ctx, cuemod.OptUpgrade(o.Upgrade)), importPath)
}
