package cmd

import (
	"context"
	"os"

	"github.com/octohelm/cuemod/pkg/cli"
	"github.com/octohelm/cuemod/pkg/version"
)

var app = cli.NewApp("cuemod", version.FullVersion(), &ProjectFlags{})

func Run(ctx context.Context) error {
	return app.Run(ctx, os.Args)
}
