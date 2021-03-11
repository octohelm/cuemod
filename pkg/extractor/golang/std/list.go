//go:generate sh ./list.sh
package std

import (
	"bufio"
	"bytes"
	_ "embed"
)

//go:embed list.txt
var list []byte

var stds = map[string]bool{}

func init() {
	scanner := bufio.NewScanner(bytes.NewBuffer(list))
	for scanner.Scan() {
		pkg := scanner.Text()
		if pkg != "" {
			stds[pkg] = true
		}
	}
}

func IsStd(importPath string) bool {
	if _, ok := stds[importPath]; ok {
		return true
	}
	return false
}
