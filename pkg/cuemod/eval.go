package cuemod

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"github.com/octohelm/cuemod/pkg/cuemod/bundle"
	"sigs.k8s.io/yaml"

	"github.com/octohelm/cuemod/pkg/cuemod/translator"
	_ "github.com/octohelm/cuemod/pkg/translator"
)

type Encoding = build.Encoding

const (
	JSON = build.JSON
	YAML = build.YAML
	CUE  = build.CUE
)

func Eval(instance *build.Instance, encoding Encoding) ([]byte, error) {
	if encoding == "" {
		encoding = build.JSON
	}

	v := cuecontext.New().BuildInstance(instance)

	if err := v.Validate(cue.Final()); err != nil {
		return nil, err
	}

	data, err := encode(instance, v, encoding)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func encode(inst *build.Instance, v cue.Value, encoding Encoding) ([]byte, error) {
	switch encoding {
	case CUE:
		return bundle.ToRaw(inst)
	case JSON:
		return translator.MarshalCueValue(v)
	case YAML:
		data, err := translator.MarshalCueValue(v)
		if err != nil {
			return nil, err
		}
		return yaml.JSONToYAML(data)
	}
	return nil, fmt.Errorf("unsupoort encoding %s", encoding)
}
