package tanka

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/grafana/tanka/pkg/term"
	"github.com/pkg/errors"
)

type FilterOpts struct {
	Targets []string `name:"target,t"  usage:"Regex filter on '<kind>/<name>'. See https://tanka.dev/output-filtering"`
}

func Process(data []byte, filters process.Matchers) (*LoadResult, error) {
	env := &v1alpha1.Environment{}

	if err := json.Unmarshal(data, env); err != nil {
		return nil, err
	}

	processed, err := process.Process(*env, filters)
	if err != nil {
		return nil, err
	}

	return &LoadResult{Env: env, Resources: processed}, nil
}

type LoadResult struct {
	Env       *v1alpha1.Environment
	Resources manifest.List
}

func (l *LoadResult) Connect() (*kubernetes.Kubernetes, error) {
	env := *l.Env

	// check env is complete
	s := ""
	if env.Spec.APIServer == "" {
		s += "  * spec.apiServer: No Kubernetes cluster endpoint specified"
	}
	if env.Spec.Namespace == "" {
		s += "  * spec.namespace: Default namespace missing"
	}
	if s != "" {
		return nil, fmt.Errorf("Your Environment's spec.json seems incomplete:\n%s\n\nPlease see https://tanka.dev/config for reference", s)
	}

	// connect client
	kube, err := kubernetes.New(env)
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
