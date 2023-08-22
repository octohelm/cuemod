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

			dirSum, err := HashDir(found.fs, ".", "", dirhash.DefaultHash)
			if err != nil {
				return nil, err
			}
			m.Sum = dirSum

			m.Dir = filepath.Join(modRoot, "cue.mod/pkg/.cuemod/std", fmt.Sprintf(
				"%s@%s-%s",
				m.Repo,
				m.Version,
				m.Sum,
			))

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

func HashDir(fs fs.ReadDirFS, dir string, prefix string, hash dirhash.Hash) (string, error) {
	files := make([]string, 0)

	if err := RangeFile(fs, dir, func(filename string) error {
		files = append(files, filepath.ToSlash(filepath.Join(prefix, filename)))
		return nil
	}); err != nil {
		return "", err
	}

	return hash(files, func(name string) (io.ReadCloser, error) {
		return fs.Open(filepath.Join(dir, strings.TrimPrefix(name, prefix)))
	})
}

func RangeFile(f fs.ReadDirFS, root string, each func(filename string) error) error {
	return fs.WalkDir(f, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel := path
		if root != "" && root != "." {
			rel, _ = filepath.Rel(root, path)
		}
		return each(rel)
	})

}
