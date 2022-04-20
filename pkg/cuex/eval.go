package cuex

import (
	"bytes"
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	cueerrors "cuelang.org/go/cue/errors"
	"sigs.k8s.io/yaml"

	_ "github.com/octohelm/cuemod/pkg/cuex/translator"
	"github.com/octohelm/cuemod/pkg/cuex/translator/core"
)

type Encoding = build.Encoding

const (
	JSON = build.JSON
	YAML = build.YAML
	CUE  = build.CUE
)

type EvalOptionFunc = func(o *EvalOption)

type EvalOption struct {
	Encoding Encoding
}

func (o EvalOption) SetDefaults() {
	if o.Encoding == "" {
		o.Encoding = build.JSON
	}
}

func WithEncoding(e build.Encoding) EvalOptionFunc {
	return func(o *EvalOption) {
		o.Encoding = e
	}
}

func WithEncodingFromFileExt(ext string) EvalOptionFunc {
	return func(o *EvalOption) {
		switch v := ext; v {
		case ".yaml":
			o.Encoding = YAML
		case ".json":
			o.Encoding = JSON
		case ".cue":
			o.Encoding = CUE
		default:
			panic(fmt.Errorf("unsupport output format %s", v))
		}
	}
}

func Eval(instance *build.Instance, options ...EvalOptionFunc) ([]byte, error) {
	o := &EvalOption{}

	for i := range options {
		options[i](o)
	}

	v := cuecontext.New().BuildInstance(instance)

	if err := v.Validate(cue.Final(), cue.Concrete(true)); err != nil {
		b := bytes.NewBuffer(nil)
		cueerrors.Print(b, err, nil)
		return nil, errors.New(b.String())
	}

	data, err := encode(instance, v, o.Encoding)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func encode(inst *build.Instance, v cue.Value, encoding Encoding) ([]byte, error) {
	switch encoding {
	case CUE:
		return BundleToRaw(inst)
	case JSON:
		return core.MarshalCueValue(v)
	case YAML:
		data, err := core.MarshalCueValue(v)
		if err != nil {
			return nil, err
		}
		return yaml.JSONToYAML(data)
	}
	return nil, fmt.Errorf("unsupoort encoding %s", encoding)
}
