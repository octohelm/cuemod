package cmd

import (
	"context"

	"github.com/octohelm/cuemod/pkg/cli"
	"github.com/octohelm/cuemod/pkg/cuex"
	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"github.com/octohelm/cuemod/pkg/plugins/kube"
)

func init() {
	k := &K{}

	k.Add(&Show{})
	k.Add(&Apply{})
	k.Add(&Prune{})
	k.Add(&Delete{})

	app.Add(k)

}

type K struct {
	cli.Name `desc:"kube commands"`
}

func (o *K) Run(ctx context.Context, args []string) error {
	return nil
}

type Show struct {
	cli.Name `args:"INPUT_AND_PATCHES..." desc:"show rendered kube manifests"`
	kube.Opts
	kube.ShowOpts
}

func (o *Show) Run(ctx context.Context, args []string) error {
	show := func(args []string) error {
		lr, err := load(ctx, args, &o.Opts)
		if err != nil {
			return err
		}
		return lr.Show(o.ShowOpts)
	}
	return show(args)
}

type Apply struct {
	cli.Name `args:"INPUT_AND_PATCHES..." desc:"apply kube manifests"`
	kube.Opts
	kube.ApplyOpts
}

func (o *Apply) Run(ctx context.Context, args []string) error {
	lr, err := load(ctx, args, &o.Opts)
	if err != nil {
		return err
	}
	return lr.Apply(o.ApplyOpts)
}

type Delete struct {
	cli.Name `args:"INPUT_AND_PATCHES..." desc:"delete kube manifests"`
	kube.Opts
	kube.DeleteOpts
}

func (o *Delete) Run(ctx context.Context, args []string) error {
	lr, err := load(ctx, args, &o.Opts)
	if err != nil {
		return err
	}
	return lr.Delete(o.DeleteOpts)
}

type Prune struct {
	cli.Name `args:"INPUT_AND_PATCHES..." desc:"prune kube manifests"`
	kube.Opts
	kube.PruneOpts
}

func (o *Prune) Run(ctx context.Context, args []string) error {
	lr, err := load(ctx, args, &o.Opts)
	if err != nil {
		return err
	}
	return lr.Prune(o.PruneOpts)
}

func load(ctx context.Context, fileOrPatches []string, opts *kube.Opts) (*kube.LoadResult, error) {
	jsonRaw, err := evalWithPatches(ctx, fileOrPatches, cuex.WithEncoding(cuex.JSON))
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

	return kube.Process(release, filters)
}
