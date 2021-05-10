package builtin

import (
	"bufio"
	"bytes"
	_ "embed"
)

//go:embed builtins.txt
var list []byte

var builtins = map[string]bool{}

func init() {
	scanner := bufio.NewScanner(bytes.NewBuffer(list))
	for scanner.Scan() {
		importPath := scanner.Text()
		if importPath != "" {
			builtins[importPath] = true
		}
	}
}

func IsBuiltIn(importPath string) bool {
	if _, ok := builtins[importPath]; ok {
		return true
	}
	return false
}
