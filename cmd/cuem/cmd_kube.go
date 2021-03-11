package main

import (
	"context"
	"fmt"
	"github.com/octohelm/cuemod/pkg/cuemod"

	"github.com/grafana/tanka/pkg/process"
	"github.com/octohelm/cuemod/pkg/tanka"
	"github.com/spf13/cobra"
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

	filters := &tanka.FilterOpts{
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

func cmdKubeShow(filters *tanka.FilterOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "show <input>",
	}

	opts := tanka.ShowOpts{}

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

func cmdKubeApply(filters *tanka.FilterOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "apply <input>",
	}

	opts := tanka.ApplyOpts{
		Validate: true,
	}

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

func cmdKubeDelete(filters *tanka.FilterOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete <input>",
	}

	opts := tanka.DeleteOpts{
		Validate: true,
	}

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

func cmdKubePrune(filters *tanka.FilterOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use: "prune <input>",
	}

	opts := tanka.PruneOpts{}

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

func load(ctx context.Context, filename string, opts *tanka.FilterOpts) (*tanka.LoadResult, error) {
	filters, err := process.StrExps(opts.Targets...)
	if err != nil {
		return nil, err
	}

	jsonRaw, err := runtime.Eval(ctx, filename, cuemod.JSON)
	if err != nil {
		return nil, err
	}

	return tanka.Process(jsonRaw, filters)
}
