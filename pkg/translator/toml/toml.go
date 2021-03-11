package toml

import (
	"encoding/json"

	"cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cuemod/translator"
	"github.com/pelletier/go-toml"
)

func init() {
	translator.Register(t{})
}

type t struct{}

func (t) Name() string {
	return "toml"
}

// FromJSON convert JSON raw to TOML
func (t) MarshalCueValue(cueValue cue.Value) ([]byte, error) {
	data, err := translator.ValueFromCueValue(cueValue)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch val := data.(type) {
	case []byte:
		if err := json.Unmarshal(val, &v); err != nil {
			return nil, err
		}
	case string:
		if err := json.Unmarshal([]byte(val), &v); err != nil {
			return nil, err
		}
	default:
		v = data
	}

	return toml.Marshal(v)
}
