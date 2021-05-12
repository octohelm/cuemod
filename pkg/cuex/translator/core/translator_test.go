package core_test

import (
	"fmt"
	"testing"

	"github.com/octohelm/cuemod/pkg/cuex"
)

func TestTranslator(t *testing.T) {
	inst, _ := cuex.InstanceFromRaw([]byte(`
name: test: """
{ "a": 1 }
""" @translate("toml")
`))

	data, err := cuex.Eval(inst, cuex.YAML)
	fmt.Println(string(data), err)
}
