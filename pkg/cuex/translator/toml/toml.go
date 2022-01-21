package toml

import (
	"cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cuex/translator/core"
	"github.com/pelletier/go-toml/v2"
	"sigs.k8s.io/yaml"
)

func init() {
	core.Register(t{})
}

type t struct{}

func (t) Name() string {
	return "toml"
}

// FromJSON convert JSON raw to TOML
func (t) MarshalCueValue(cueValue cue.Value) ([]byte, error) {
	data, err := core.ValueFromCueValue(cueValue)
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

	return toml.Marshal(convert(v))
}

func convert(v interface{}) interface{} {
	switch x := v.(type) {
	case float64:
		// if can int should be int
		if float64(int64(x)) == x {
			return int64(x)
		}
	case []interface{}:
		list := make([]interface{}, len(x))
		for i := range list {
			list[i] = convert(x[i])
		}
		return list
	case map[string]interface{}:
		m := map[string]interface{}{}
		for k, val := range x {
			m[k] = convert(val)
		}
		return m
	}
	return v
}
