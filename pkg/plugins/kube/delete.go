package kube

import (
	"context"
	"fmt"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
)

type DeleteOpts struct {
	AutoApprove bool `flag:"auto-approve" desc:"skips the interactive approval"`
}

func (l *LoadResult) Delete(opts DeleteOpts) error {
	kube, err := l.Connect()
	if err != nil {
		return err
	}

	resources := manifest.ProcessManifests(l.Resources,
		manifest.IgnoreNamespace(),
	)

	for i := range resources {
		fmt.Printf("to delete %s\n", manifest.Identity(resources[i]))
	}

	// prompt for confirmation
	if opts.AutoApprove {
	} else if err := confirmPrompt("Deleting from", l.Release, kube.Info()); err != nil {
		return err
	}

	if err := kube.DeleteResources(context.Background(), resources); err != nil {
		return err
	}

	for _, m := range resources {
		fmt.Printf("%s deleted\n", manifest.Identity(m))
	}

	return nil
}
