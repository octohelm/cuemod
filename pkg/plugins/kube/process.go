package kube

import (
	"encoding/json"
	"fmt"

	"github.com/octohelm/cuemod/pkg/kubernetes"
	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	"github.com/octohelm/cuemod/pkg/term"

	releasev1alpha "github.com/octohelm/cuemod/pkg/api/release/v1alpha"
	"github.com/pkg/errors"

	"github.com/fatih/color"
)

type Opts struct {
	AsTemplate bool     `name:"as-template"  usage:"As release secret"`
	Targets    []string `name:"target,t"  usage:"Regex filter on '<kind>/<name>'. See https://tanka.dev/output-filtering"`
}

func ReleaseFromJSONRaw(data []byte) (*releasev1alpha.Release, error) {
	release := &releasev1alpha.Release{}

	if err := json.Unmarshal(data, release); err != nil {
		return nil, err
	}

	return release, nil
}

func Process(release *releasev1alpha.Release, filters manifest.Matchers) (*LoadResult, error) {
	processed, err := ProcessResources(release, filters)
	if err != nil {
		return nil, err
	}
	return &LoadResult{Release: release, Resources: processed}, nil
}

type LoadResult struct {
	Release   *releasev1alpha.Release
	Resources manifest.List
}

func (l *LoadResult) Connect() (*kubernetes.KubeClient, error) {
	release := *l.Release

	// check release is complete
	s := ""

	if context, ok := release.Metadata.Labels["context"]; !ok || context == "" {
		s += "  * metadata.labels.context: No Kubernetes context"
	}
	if release.Metadata.Namespace == "" {
		s += "  * metadata.namespace: Default namespace missing"
	}

	if s != "" {
		return nil, fmt.Errorf("Your Release's spec seems incomplete:\n%s\n\nPlease see https://github.com/octohelm/cuemod#plugin-kube for reference", s)
	}

	// connect client
	kube, err := kubernetes.NewClientForContext(release.Metadata.Labels["context"])
	if err != nil {
		return nil, errors.Wrap(err, "connecting to Kubernetes")
	}

	return kube, nil
}

func confirmPrompt(action string, r *releasev1alpha.Release, info string) error {
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	return term.Confirm(
		fmt.Sprintf(
			`%s namespace '%s' of %s (%s).`,
			action,
			alert(r.Metadata.Namespace),
			alert(r.Metadata.Labels["context"]),
			info,
		),
		"yes",
	)
}
