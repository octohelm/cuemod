package cue

import (
	"fmt"

	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue/build"
	"sigs.k8s.io/yaml"
)

type Encoding = build.Encoding

const (
	JSON = build.JSON
	YAML = build.YAML
)

func encode(v cue.Value, encoding Encoding) ([]byte, error) {
	switch encoding {
	case JSON:
		return v.MarshalJSON()
	case YAML:
		data, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}
		return yaml.JSONToYAML(data)
	}
	return nil, fmt.Errorf("unsupoort encoding %s", encoding)
}
