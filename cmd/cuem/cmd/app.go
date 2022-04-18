package cmd

import (
	"context"
	"os"

	"github.com/octohelm/cuemod/internal/version"
	"github.com/octohelm/cuemod/pkg/cli"
)

var app = cli.NewApp("cuemod", version.Version, &ProjectFlags{})

func Run(ctx context.Context) error {
	return app.Run(ctx, os.Args)
}
