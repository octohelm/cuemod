package tanka

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	releasev1alpha "github.com/octohelm/cuemod/pkg/api/release/v1alpha"
	"github.com/pkg/errors"

	"github.com/fatih/color"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/term"
)

type FilterOpts struct {
	Targets []string `name:"target,t"  usage:"Regex filter on '<kind>/<name>'. See https://tanka.dev/output-filtering"`
}

func Process(data []byte, filters process.Matchers) (*LoadResult, error) {
	release := &releasev1alpha.Release{}

	if err := json.Unmarshal(data, release); err != nil {
		return nil, err
	}

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

func (l *LoadResult) Connect() (*kubernetes.Kubernetes, error) {
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

	ip, err := client.IPFromContext(release.Metadata.Labels["context"])
	if err != nil {
		return nil, fmt.Errorf("can't")
	}

	// connect client
	kube, err := kubernetes.New(v1alpha1.Environment{
		Metadata: v1alpha1.Metadata{
			Name:      release.Metadata.Name,
			Namespace: release.Metadata.Namespace,
		},
		Spec: v1alpha1.Spec{
			APIServer:    ip,
			Namespace:    release.Metadata.Namespace,
			InjectLabels: true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "connecting to Kubernetes")
	}

	return kube, nil
}

func confirmPrompt(action, namespace string, info client.Info) error {
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	return term.Confirm(
		fmt.Sprintf(`%s namespace '%s' of cluster '%s' at '%s' using context '%s'.`, action,
			alert(namespace),
			alert(info.Kubeconfig.Cluster.Name),
			alert(info.Kubeconfig.Cluster.Cluster.Server),
			alert(info.Kubeconfig.Context.Name),
		),
		"yes",
	)
}
