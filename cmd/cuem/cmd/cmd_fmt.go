package cmd

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cli"
	"github.com/octohelm/cuemod/pkg/cuemod"

	"github.com/octohelm/cuemod/pkg/cuex/format"
)

func init() {
	app.Add(&Fmt{})
}

type Fmt struct {
	cli.Name `args:"BASE_DIR" desc:"base dir for files fmt"`
	format.FormatOpts
}

func (opts *Fmt) Run(ctx context.Context, args []string) error {
	baseDir := "./"
	if len(args) > 0 {
		baseDir = args[0]
	}

	files, err := cuemod.FromContext(ctx).ListCue(baseDir)
	if err != nil {
		return err
	}

	return format.FormatFiles(ctx, files, opts.FormatOpts)
}
