package cuex

import (
	"fmt"
	"testing"
)

func TestInstanceFromRaw(t *testing.T) {
	inst, _ := InstanceFromTemplateAndOverwrites([]byte(`
name: test: *"1111" | string
`), []byte(`
import t "t"

t & {
	name: test: "2222"
}
`))
	data, err := Eval(inst, CUE)
	fmt.Println(string(data), err)

	data2, err2 := Eval(inst, YAML)
	fmt.Println(string(data2), err2)
}
