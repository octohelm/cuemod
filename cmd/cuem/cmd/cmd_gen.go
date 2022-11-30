package cmd

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cueify"

	"github.com/octohelm/cuemod/pkg/cli"
)

func init() {
	app.Add(&Gen{})
}

type Gen struct {
	cli.Name `args:"PATH" desc:"source path"`
	Output   string `flag:"output,o" desc:"generate output"`
	Import   string `flag:"import,i" desc:"declare language for generate. support values: crd | go"`
}

func (o *Gen) Run(ctx context.Context, args []string) error {
	return cueify.ExtractToDir(
		ctx,
		o.Import,
		args[0],
		o.Output,
	)
}
