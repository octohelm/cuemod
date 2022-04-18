package kube

import (
	"context"
	"fmt"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"github.com/octohelm/cuemod/pkg/term"
)

type ApplyOpts struct {
	AutoApprove bool `flag:"auto-approve" desc:"skips the interactive approval"`
}

func (l *LoadResult) Apply(opts ApplyOpts) error {
	kc, err := l.Connect()
	if err != nil {
		return err
	}
	// show diff
	diff, err := kc.Diff(context.Background(), l.Resources.DeepCopy())
	switch {
	case err != nil:
		// This is not fatal, the diff is not strictly required
		fmt.Println("Error diffing:", err)
	case diff == nil:
		tmp := "Warning: There are no differences. Your apply may not do anything at all."
		diff = []byte(tmp)
	}

	if diff != nil {
		b := term.Colordiff(diff)
		fmt.Print(b.String())
	}

	if opts.AutoApprove {
	} else if err := confirmPrompt("Applying to", l.Release, kc.Info()); err != nil {
		return err
	}

	if err := kc.ApplyResources(context.Background(), l.Resources); err != nil {
		return err
	}

	for i := range l.Resources {
		fmt.Printf("%s %s\n", manifest.Identity(l.Resources[i]), l.Resources[i].GetAnnotations()[manifest.AnnotationReleaseStage])
	}

	return nil
}
