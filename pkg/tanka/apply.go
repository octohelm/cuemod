package tanka

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/term"
)

type ApplyOpts struct {
	DiffStrategy string `name:"diff-strategy" usage:"force the diff-strategy to use. Automatically chosen if not set. values: native, subset"`
	AutoApprove  bool   `name:"auto-approve" usage:"skips the interactive approval"`
	Force        bool   `name:"force,f" usage:"ignores any warnings kubectl might have"`
	Validate     bool   `name:"validate" usage:"set to false ignores invalid Kubernetes schemas"`
}

func (l *LoadResult) Apply(opts ApplyOpts) error {
	kube, err := l.Connect()
	if err != nil {
		return err
	}
	defer kube.Close()

	// show diff
	diff, err := kube.Diff(l.Resources, kubernetes.DiffOpts{Strategy: opts.DiffStrategy})
	switch {
	case err != nil:
		// This is not fatal, the diff is not strictly required
		fmt.Println("Error diffing:", err)
	case diff == nil:
		tmp := "Warning: There are no differences. Your apply may not do anything at all."
		diff = &tmp
	}

	// in case of non-fatal error diff may be nil
	if diff != nil {
		b := term.Colordiff(*diff)
		fmt.Print(b.String())
	}

	// prompt for confirmation
	if opts.AutoApprove {
	} else if err := confirmPrompt("Applying to", l.Env.Spec.Namespace, kube.Info()); err != nil {
		return err
	}

	return kube.Apply(l.Resources, kubernetes.ApplyOpts{
		Force:    opts.Force,
		Validate: opts.Validate,
	})
}
