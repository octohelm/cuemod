package main

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cuex/format"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		cmdFmt(),
	)
}

func cmdFmt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "format jsonnet codes",
	}

	formatOpts := format.FormatOpts{}

	return setupRun(cmd, &formatOpts, func(ctx context.Context, args []string) error {
		baseDir := "./"
		if len(args) > 0 {
			baseDir = args[0]
		}

		files, err := runtime.ListCue(baseDir)
		if err != nil {
			return err
		}

		return format.FormatFiles(cmd.Context(), files, formatOpts)
	})
}
