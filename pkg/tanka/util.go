package tanka

import (
	"os"
	"path/filepath"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}

func ignoreNamespace(list manifest.List) (l manifest.List) {
	for i := range list {
		m := list[i]
		if m.Kind() == "Namespace" {
			continue
		}
		l = append(l, m)
	}
	return
}
