package main

import (
	"context"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/internal/version"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/spf13/cobra"
)

var (
	runtime     *cuemod.Runtime
	rootCmd     = cmdRoot()
	log         = logr.StdLogger()
	projectOpts = &ProjectOpts{}
)

type ProjectOpts struct {
	Root    string `name:"project,p" usage:"project root dir"`
	Verbose bool   `name:"verbose,v" usage:"verbose"`
}

func cmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cuem",
		Short:   "cue mod",
		Version: version.Version,
	}

	return setupPersistentPreRun(cmd, projectOpts, func(ctx context.Context, args []string) error {
		runtime = cuemod.RuntimeFor(projectOpts.Root)

		if projectOpts.Verbose {
			log.(interface{ SetLevel(lvl logr.Level) }).SetLevel(logr.TraceLevel)
		}

		return nil
	})
}
