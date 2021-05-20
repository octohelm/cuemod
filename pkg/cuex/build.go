package cuex

import (
	"cuelang.org/go/cue/build"
	cueparser "cuelang.org/go/cue/parser"
)

func InstanceFromRaw(src []byte) (*build.Instance, error) {
	inst := build.NewContext().NewInstance("", nil)
	f, err := cueparser.ParseFile("main.cue", src)
	if err != nil {
		return nil, err
	}
	if err := inst.AddSyntax(f); err != nil {
		return nil, err
	}
	return inst, nil
}

func InstanceFromTemplateAndOverwrites(template []byte, overwrites []byte) (*build.Instance, error) {
	t, err := InstanceFromRaw(template)
	if err != nil {
		return nil, err
	}

	if overwrites == nil {
		overwrites = []byte(`
import t "t"
t & {}
`)
	}

	m, err := InstanceFromRaw(overwrites)
	if err != nil {
		return nil, err
	}

	t.ImportPath = "t"

	m.Imports = append(m.Imports, t)

	return m, nil
}
