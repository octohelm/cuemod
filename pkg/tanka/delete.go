package tanka

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/term"
)

type DeleteOpts struct {
	AutoApprove bool `name:"auto-approve" usage:"skips the interactive approval"`
	Force       bool `name:"force,f" usage:"ignores any warnings kubectl might have"`
	Validate    bool `name:"validate" usage:"set to false ignores invalid Kubernetes schemas"`
}

func (l *LoadResult) Delete(opts DeleteOpts) error {
	kube, err := l.Connect()
	if err != nil {
		return err
	}
	defer kube.Close()

	// show diff
	// static differ will never fail and always return something if input is not nil
	diff, err := kubernetes.StaticDiffer(false)(l.Resources)

	if err != nil {
		fmt.Println("Error diffing:", err)
	}

	// in case of non-fatal error diff may be nil
	if diff != nil {
		b := term.Colordiff(*diff)
		fmt.Print(b.String())
	}

	// prompt for confirmation
	if opts.AutoApprove {
	} else if err := confirmPrompt("Deleting from", l.Env.Spec.Namespace, kube.Info()); err != nil {
		return err
	}

	return kube.Delete(l.Resources, kubernetes.DeleteOpts{
		Force:    opts.Force,
		Validate: opts.Validate,
	})
}
