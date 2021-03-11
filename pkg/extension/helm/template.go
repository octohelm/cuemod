package helm

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/octohelm/cuemod/pkg/cue/native"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func init() {
	native.Register(pkg{})
}

type pkg struct{}

func (pkg) ImportPath() string {
	return "extension/helm"
}

type ReleaseOption struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// Template
func (pkg) Template(c *chart.Chart, chartValues map[string]interface{}, releaseOptions ReleaseOption) (map[string]map[string]interface{}, error) {
	if c == nil {
		return nil, nil
	}

	c.Templates = make([]*chart.File, 0)

	for i := range c.Files {
		f := c.Files[i]
		if strings.HasPrefix(f.Name, "templates/") {
			c.Templates = append(c.Templates, &chart.File{Name: f.Name, Data: f.Data})
		}
	}

	values, err := chartutil.CoalesceValues(c, chartValues)
	if err != nil {
		return nil, err
	}
	valuesToRender, err := chartutil.ToRenderValues(c, values, chartutil.ReleaseOptions{
		Name:      releaseOptions.Name,
		Namespace: releaseOptions.Namespace,
	}, nil)
	if err != nil {
		return nil, err
	}

	e := &engine.Engine{}

	renderedContentMap, err := e.Render(c, valuesToRender)
	if err != nil {
		return nil, err
	}

	for _, c := range c.CRDObjects() {
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

	return manifests, nil
}
