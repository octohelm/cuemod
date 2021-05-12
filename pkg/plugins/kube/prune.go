package kube

import (
	"context"
	"fmt"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PruneOpts struct {
	AutoApprove bool `name:"auto-approve" usage:"skips the interactive approval"`
}

func (l *LoadResult) Prune(opts PruneOpts) error {
	kc, err := l.Connect()
	if err != nil {
		return err
	}

	allGroupVersionKinds, err := kc.AllListableGroupVersionKinds()
	if err != nil {
		return err
	}

	liveResources, err := kc.ListAll(
		context.Background(),
		allGroupVersionKinds,
		client.InNamespace(l.Release.Namespace),
		client.MatchingLabels(map[string]string{
			manifest.LabelRelease: l.Release.Name,
		}),
	)
	if err != nil {
		return err
	}

	orphaned := liveResources.Orphaned(l.Resources)

	if len(orphaned) == 0 {
		fmt.Println("Nothing found to prune.")
		return nil
	}

	for i := range orphaned {
		fmt.Printf("0 %s\n", manifest.Identity(orphaned[i]))
	}

	// prompt for confirm
	if opts.AutoApprove {
	} else if err := confirmPrompt("Pruning from", l.Release, kc.Info()); err != nil {
		return err
	}

	// delete resources
	return kc.DeleteResources(context.Background(), orphaned)
}
