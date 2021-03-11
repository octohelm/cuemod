package toml

import (
	"encoding/json"

	"github.com/octohelm/cuemod/pkg/cue/native"
	"github.com/pelletier/go-toml"
)

func init() {
	native.Register(pkg{})
}

type pkg struct{}

func (pkg) ImportPath() string {
	return "extension/toml"
}

// FromJSON convert JSON raw to TOML
func (pkg) FromJSON(data string) (string, error) {
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return "", err
	}
	t, err := toml.TreeFromMap(v)
	if err != nil {
		return "", err
	}
	return t.ToTomlString()
}
