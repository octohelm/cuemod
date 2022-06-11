package stdlib

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/octohelm/cuemod/pkg/cuemod/mod"
	"golang.org/x/mod/sumdb/dirhash"
)

var stdlibs = map[string]*stdlib{}

type stdlib struct {
	mod.Mod
	fs fs.ReadDirFS
}

func Register(dirFs fs.ReadDirFS, version string, pkgs ...string) {
	for _, p := range pkgs {
		m := &stdlib{}
		m.Repo = p
		m.Module = p
		m.Version = version
		m.Root = true
		m.fs = dirFs

		stdlibs[m.Repo] = m
	}
}

func RepoRootForImportPath(repo string) (string, bool) {
	for lib := range stdlibs {
		if isSubDirFor(repo, lib) {
			return lib, true
		}
	}
	return "", false
}

func Mount(ctx context.Context, importPath string, modRoot string) (*mod.Mod, error) {
	repo, ok := RepoRootForImportPath(importPath)
	if !ok {
		return nil, nil
	}

	noWrite := modRoot == ""
	if noWrite {
		modRoot = "<cue_mod_root>"
	}

	for lib := range stdlibs {
		if isSubDirFor(repo, lib) {
			found := stdlibs[lib]

			m := found.Mod
			m.Dir = filepath.Join(modRoot, "cue.mod/pkg/.stdlib", fmt.Sprintf("%s@%s", m.Repo, m.Version))

			if !noWrite {
				if _, err := os.Stat(m.Dir); os.IsNotExist(err) {
					_ = os.MkdirAll(m.Dir, os.ModePerm)

					err := fs.WalkDir(found.fs, repo, func(path string, d fs.DirEntry, err error) error {
						tmpPath := filepath.Join(m.Dir, "."+strings.TrimPrefix(path, repo))

						if d.IsDir() {
							_ = os.MkdirAll(tmpPath, os.ModePerm)
						} else {
							f, err := found.fs.Open(path)
							if err != nil {
								return err
							}
							defer f.Close()
							data, _ := io.ReadAll(f)
							if err = os.WriteFile(tmpPath, data, os.ModePerm); err != nil {
								return err
							}
						}

						return err
					})

					if err != nil {
						return nil, err
					}

					dirSum, err := dirhash.HashDir(m.Dir, "stdlib", dirhash.DefaultHash)
					if err != nil {
						return nil, err
					}
					m.Sum = dirSum
				}
			}

			return &m, nil
		}
	}

	return nil, nil
}

func isSubDirFor(targetpath string, root string) bool {
	targetpath = targetpath + "/"
	root = root + "/"
	return strings.HasPrefix(targetpath, root)
}
