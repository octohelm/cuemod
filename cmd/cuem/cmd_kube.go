package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/octohelm/cuemod/pkg/cuemoperator"

	"github.com/spf13/cobra"

	"github.com/octohelm/cuemod/pkg/cuex"
	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"github.com/octohelm/cuemod/pkg/plugins/kube"
)

func init() {
	rootCmd.AddCommand(
		cmdKube(),
	)
}

func cmdKube() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kube",
		Aliases: []string{"k"},
	}

	filters := &kube.Opts{
		Targets: []string{},
	}

	bindFlags(cmd.PersistentFlags(), filters)

	cmd.AddCommand(
		cmdKubeApply(filters),
		cmdKubeShow(filters),
		cmdKubeDelete(filters),
		cmdKubePrune(filters),
	)

	return cmd
}

func cmdKubeShow(filters *kube.Opts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "show <input>",
	}

	opts := kube.ShowOpts{}

	return setupRun(cmd, &opts, func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing input")
		}

		show := func(input string) error {
			lr, err := load(ctx, input, filters)
			if err != nil {
				return err
			}
			return lr.Show(opts)
		}

		if opts.Output != "" {
			for _, arg := range args {
				if err := show(arg); err != nil {
					return err
				}
			}
			return nil
		}

		return show(args[0])
	})
}

func cmdKubeApply(filters *kube.Opts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "apply <input>",
	}

	opts := kube.ApplyOpts{}

	return setupRun(cmd, &opts, func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing input")
		}

		lr, err := load(ctx, args[0], filters)
		if err != nil {
			return err
		}

		return lr.Apply(opts)
	})
}

func cmdKubeDelete(filters *kube.Opts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete <input>",
	}

	opts := kube.DeleteOpts{}

	return setupRun(cmd, &opts, func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing input")
		}

		lr, err := load(ctx, args[0], filters)
		if err != nil {
			return err
		}

		return lr.Delete(opts)
	})
}

func cmdKubePrune(filters *kube.Opts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "prune <input>",
	}

	opts := kube.PruneOpts{}

	return setupRun(cmd, &opts, func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing input")
		}

		lr, err := load(ctx, args[0], filters)
		if err != nil {
			return err
		}

		return lr.Prune(opts)
	})
}

func load(ctx context.Context, filename string, opts *kube.Opts) (*kube.LoadResult, error) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, filename)

	jsonRaw, err := runtime.Eval(ctx, path, cuex.JSON)
	if err != nil {
		return nil, err
	}

	release, err := kube.ReleaseFromJSONRaw(jsonRaw)
	if err != nil {
		return nil, err
	}

	filters, err := manifest.StrExps(opts.Targets...)
	if err != nil {
		return nil, err
	}

	if opts.AsTemplate {
		cueTemplate, err := runtime.Eval(ctx, path, cuex.CUE)
		if err != nil {
			return nil, err
		}

		s := cuemoperator.NewReleaseTemplate(release.Namespace, release.Name, cueTemplate)

		return &kube.LoadResult{
			Release:   release,
			Resources: []manifest.Object{s},
		}, nil
	}

	return kube.Process(release, filters)
}
