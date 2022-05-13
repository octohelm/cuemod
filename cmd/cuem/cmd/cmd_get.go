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
	cli.Name `args:"IMPORT_PATH..." desc:"download dependencies"`
	Upgrade  bool   `flag:"upgrade,u" desc:"upgrade dependencies"`
	Import   string `flag:"import,i" desc:"declare language for generate. support values: crd | go"`
}

func (o *Get) Run(ctx context.Context, args []string) error {
	cc := cuemod.FromContext(ctx)

	for i := range args {
		p := args[i]
		err := cc.Get(
			cuemod.WithOpts(ctx,
				cuemod.OptUpgrade(o.Upgrade),
				cuemod.OptImport(o.Import),
				cuemod.OptVerbose(true),
			), p,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
