package translator_test

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue/build"
	cueparser "cuelang.org/go/cue/parser"
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

func instance(src []byte) *build.Instance {
	p := build.NewContext().NewInstance("dir", nil)
	f, _ := cueparser.ParseFile("main.cue", src)
	_ = p.AddSyntax(f)
	_ = p.Complete()
	return p
}
