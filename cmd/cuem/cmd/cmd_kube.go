package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/octohelm/cuemod/pkg/cli"
	"github.com/octohelm/cuemod/pkg/cuemod"
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
	cli.Name `args:"INPUT" desc:"show rendered kube manifests"`
	kube.Opts
	kube.ShowOpts
}

func (o *Show) Run(ctx context.Context, args []string) error {
	show := func(input string) error {
		lr, err := load(ctx, input, &o.Opts)
		if err != nil {
			return err
		}
		return lr.Show(o.ShowOpts)
	}

	if o.Output != "" {
		for _, arg := range args {
			if err := show(arg); err != nil {
				return err
			}
		}
		return nil
	}

	return show(args[0])
}

type Apply struct {
	cli.Name `args:"INPUT" desc:"apply kube manifests"`
	kube.Opts
	kube.ApplyOpts
}

func (o *Apply) Run(ctx context.Context, args []string) error {
	lr, err := load(ctx, args[0], &o.Opts)
	if err != nil {
		return err
	}
	return lr.Apply(o.ApplyOpts)
}

type Delete struct {
	cli.Name `args:"INPUT" desc:"delete kube manifests"`
	kube.Opts
	kube.DeleteOpts
}

func (o *Delete) Run(ctx context.Context, args []string) error {
	lr, err := load(ctx, args[0], &o.Opts)
	if err != nil {
		return err
	}
	return lr.Delete(o.DeleteOpts)
}

type Prune struct {
	cli.Name `args:"INPUT" desc:"prune kube manifests"`
	kube.Opts
	kube.PruneOpts
}

func (o *Prune) Run(ctx context.Context, args []string) error {
	lr, err := load(ctx, args[0], &o.Opts)
	if err != nil {
		return err
	}
	return lr.Prune(o.PruneOpts)
}

func load(ctx context.Context, filename string, opts *kube.Opts) (*kube.LoadResult, error) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, filename)
	runtime := cuemod.FromContext(ctx)

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

	return kube.Process(release, filters)
}
