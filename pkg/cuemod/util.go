package cuemod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/semver"
)

func versionGreaterThan(v string, w string) bool {
	if w == "" {
		return false
	}
	return semver.Compare(v, w) > 0
}

func isSubDirFor(targetpath string, root string) bool {
	targetpath = targetpath + "/"
	root = root + "/"
	return strings.HasPrefix(targetpath, root)
}

func subDir(pkg string, importPath string) (string, error) {
	if isSubDirFor(importPath, pkg) {
		if len(importPath) > len(pkg)+1 {
			return importPath[len(pkg)+1:], nil
		}
		return "", nil
	}
	return "", fmt.Errorf("%s is not sub CompletePath of %s", importPath, pkg)
}

func replaceImportPath(to string, from string, importPath string) string {
	if from == importPath {
		return to
	}
	s, _ := subDir(from, importPath)
	return filepath.Join(to, s)
}

func paths(path string) []string {
	paths := make([]string, 0)
	d := path
	for {
		paths = append(paths, d)
		if !strings.Contains(d, "/") {
			break
		}
		d = filepath.Join(d, "../")
	}
	return paths
}

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
