package translator_test

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"github.com/octohelm/cuemod/pkg/cuemod"
)

func TestTranslator(t *testing.T) {
	i := instance([]byte(`
name: test: """
{ "a": 1 }
""" @translate("toml")
`))

	data, err := cuemod.Eval(i, cuemod.YAML)
	fmt.Println(string(data), err)
}

func instance(src []byte) *cue.Instance {
	p := build.NewContext().NewInstance("dir", nil)
	_ = p.AddFile("main.cue", src)
	_ = p.Complete()
	return cue.Build([]*build.Instance{p})[0]
}
