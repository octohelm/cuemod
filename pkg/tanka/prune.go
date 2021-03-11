package tanka

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/term"
)

type PruneOpts struct {
	AutoApprove bool `name:"auto-approve" usage:"skips the interactive approval"`
	Force       bool `name:"force,f" usage:"ignores any warnings kubectl might have"`
}

func (l *LoadResult) Prune(opts PruneOpts) error {
	kube, err := l.Connect()
	if err != nil {
		return err
	}
	defer kube.Close()

	// find orphaned resources
	orphaned, err := kube.Orphaned(l.Resources)
	if err != nil {
		return err
	}

	if len(orphaned) == 0 {
		fmt.Println("Nothing found to prune.")
		return nil
	}

	// print diff
	diff, err := kubernetes.StaticDiffer(false)(orphaned)
	if err != nil {
		// static diff can't fail normally, so unlike in apply, this is fatal
		// here
		return err
	}
	fmt.Print(term.Colordiff(*diff).String())

	// prompt for confirm
	if opts.AutoApprove {
	} else if err := confirmPrompt("Pruning from", l.Env.Spec.Namespace, kube.Info()); err != nil {
		return err
	}

	// delete resources
	return kube.Delete(orphaned, kubernetes.DeleteOpts{
		Force: opts.Force,
	})
}
