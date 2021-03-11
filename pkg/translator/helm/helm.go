package helm

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cuemod/translator"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func init() {
	translator.Register(t{})
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
	chart.Chart
	Release ReleaseOption `json:"release"`
}

func (t) MarshalCueValue(value cue.Value) ([]byte, error) {
	c := ChartWithRelease{}

	if err := translator.UnmarshalCueValue(value, &c); err != nil {
		return nil, err
	}

	c.Chart.Templates = make([]*chart.File, 0)

	for i := range c.Chart.Files {
		f := c.Chart.Files[i]
		if strings.HasPrefix(f.Name, "templates/") {
			c.Chart.Templates = append(c.Chart.Templates, &chart.File{Name: f.Name, Data: f.Data})
		}
	}

	values, err := chartutil.CoalesceValues(&c.Chart, c.Chart.Values)
	if err != nil {
		return nil, err
	}

	valuesToRender, err := chartutil.ToRenderValues(&c.Chart, values, chartutil.ReleaseOptions{
		Name:      c.Release.Name,
		Namespace: c.Release.Namespace,
	}, nil)
	if err != nil {
		return nil, err
	}

	e := &engine.Engine{}

	renderedContentMap, err := e.Render(&c.Chart, valuesToRender)
	if err != nil {
		return nil, err
	}

	for _, c := range c.Chart.CRDObjects() {
		renderedContentMap[c.Filename] = string(c.File.Data)
	}

	manifests := map[string]map[string]interface{}{}

	for fileName, renderedContent := range renderedContentMap {
		if filepath.Ext(fileName) != ".yaml" || filepath.Ext(fileName) == ".yml" {
			continue
		}

		if strings.TrimSpace(renderedContent) != "" {
			decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(renderedContent), 4096)

			var manifest map[string]interface{}

			if err := decoder.Decode(&manifest); err != nil {
				if err == io.EOF {
					break
				}
				return nil, errors.Wrap(err, "helm template failed")
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

	return json.Marshal(manifests)
}
