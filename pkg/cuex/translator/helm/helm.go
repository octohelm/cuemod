package helm

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/octohelm/cuemod/pkg/cuex/translator/core"
)

func init() {
	core.Register(t{})
}

type t struct{}

func (t) Name() string {
	return "helm"
}

type ReleaseOption struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type ChartWithRelease struct {
	Metadata *chart.Metadata        `json:"metadata"`
	Files    []*chart.File          `json:"files"`
	Values   map[string]interface{} `json:"values"`
	Defaults map[string]interface{} `json:"defaults"`
	Release  ReleaseOption          `json:"release"`
}

func (t) MarshalCueValue(value cue.Value) ([]byte, error) {
	c := ChartWithRelease{}

	if err := core.UnmarshalCueValue(value, &c); err != nil {
		return nil, err
	}

	helmChart := &chart.Chart{}
	helmChart.Metadata = c.Metadata
	helmChart.Files = c.Files
	helmChart.Values = c.Defaults

	helmChart.Templates = make([]*chart.File, 0)

	for i := range helmChart.Files {
		f := helmChart.Files[i]
		if strings.HasPrefix(f.Name, "templates/") {
			helmChart.Templates = append(helmChart.Templates, &chart.File{Name: f.Name, Data: f.Data})
		}
	}

	values, err := chartutil.CoalesceValues(helmChart, c.Values)
	if err != nil {
		return nil, err
	}
	valuesToRender, err := chartutil.ToRenderValues(helmChart, values, chartutil.ReleaseOptions{
		Name:      c.Release.Name,
		Namespace: c.Release.Namespace,
	}, nil)
	if err != nil {
		return nil, err
	}

	e := &engine.Engine{}

	renderedContentMap, err := e.Render(helmChart, valuesToRender)
	if err != nil {
		return nil, err
	}

	for _, c := range helmChart.CRDObjects() {
		renderedContentMap[c.Filename] = string(c.File.Data)
	}

	manifests := map[string]map[string]interface{}{}

	for fileName, renderedContent := range renderedContentMap {
		if filepath.Ext(fileName) != ".yaml" || filepath.Ext(fileName) == ".yml" {
			continue
		}

		if trimmedContent := strings.TrimSpace(renderedContent); trimmedContent != "" {
			decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(trimmedContent), 4096)

			for {
				var manifest map[string]interface{}

				if err := decoder.Decode(&manifest); err != nil {
					if err == io.EOF {
						break
					}
					return nil, errors.Wrapf(err, "helm template failed: %s\n%s", fileName, trimmedContent)
				}

				if len(manifest) == 0 {
					continue
				}

				if kind, ok := manifest["kind"]; ok {
					if k, ok := kind.(string); ok {
						if metadata, ok := manifest["metadata"]; ok {
							if m, ok := metadata.(map[string]interface{}); ok {
								if name, ok := m["name"]; ok {
									if n, ok := name.(string); ok {
										if manifests[k] == nil {
											manifests[k] = map[string]interface{}{}
										}
										manifests[k][n] = manifest
									}
								}
							}

						}
					}
				}
			}
		}
	}

	return json.Marshal(manifests)
}
