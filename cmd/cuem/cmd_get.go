package main

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cuemod"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		cmdGet(),
	)
}

func cmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "download dependencies",
	}

	o := cuemod.Opts{}

	return setupRun(cmd, &o, func(ctx context.Context, args []string) error {
		importPath := "."
		if len(args) > 0 {
			importPath = args[0]
		}
		return runtime.Get(cuemod.WithOpts(ctx, cuemod.OptUpgrade(o.Upgrade)), importPath)
	})
}
