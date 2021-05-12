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
