package yaml

import (
	"cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cuex/translator"
	"sigs.k8s.io/yaml"
)

func init() {
	translator.Register(t{})
}

type t struct{}

func (t) Name() string {
	return "yaml"
}

// FromJSON convert JSON raw to YAML
func (t) MarshalCueValue(cueValue cue.Value) ([]byte, error) {
	data, err := translator.ValueFromCueValue(cueValue)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch val := data.(type) {
	case []byte:
		if err := yaml.Unmarshal(val, &v); err != nil {
			return nil, err
		}
	case string:
		if err := yaml.Unmarshal([]byte(val), &v); err != nil {
			return nil, err
		}
	default:
		v = data
	}

	return yaml.Marshal(v)
}
