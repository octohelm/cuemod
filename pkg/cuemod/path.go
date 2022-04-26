package cuemod

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/octohelm/cuemod/pkg/extractor"
)

func PathFor(mod *Mod, importPath string) *Path {
	i := &Path{}
	i.Mod = mod
	i.Module = mod.Module

	if importPath != "" {
		// when use replace with custom forked repo
		// mod.Module may not same as mod.Repo
		// then should to overwrite this path.Module with mode.Repo
		if !strings.HasPrefix(importPath, i.Module) {
			if strings.HasPrefix(importPath, i.Repo) {
				i.Module = i.Repo
			} else {
				panic(fmt.Errorf("import path `%s` not match %s", importPath, mod))
			}
		}

		i.SubPath, _ = subPath(i.Module, importPath)
	}

	return i
}

type Path struct {
	*Mod
	Module  string
	SubPath string
	Replace *ReplaceRule
}

type ReplaceRule struct {
	From string
	modfile.ReplaceTarget
}

func (i Path) WithReplace(from string, replaceTarget modfile.ReplaceTarget) *Path {
	i.Replace = &ReplaceRule{
		From:          from,
		ReplaceTarget: replaceTarget,
	}
	return &i
}

func (i *Path) SymlinkOrImport(ctx context.Context, root string) error {
	pkgRoot := "cue.mod/pkg"

	if root == i.Dir {
		// skip root dir
		return nil
	}

	gen := ""

	if i.Replace != nil && i.Replace.Import != "" {
		gen = i.Replace.Import
	}

	if gen != "" {
		pkgRoot = "cue.mod/gen/vendor"
	}

	importPath := i.ImportPath()

	repoRootDir := filepath.Join(root, pkgRoot, i.Repo)
	importPathDir := filepath.Join(root, pkgRoot, importPath)

	if err := i.symlink(ctx, i.RepoRootDir(), repoRootDir); err != nil {
		return err
	}

	if i.shouldReplace() {
		if err := i.symlink(ctx, i.ImportPathDir(), importPathDir); err != nil {
			return err
		}
	}

	if gen != "" {
		err := extractor.ExtractToDir(
			ctx,
			gen,
			importPathDir,
			filepath.Join(root, "cue.mod/gen", importPath),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Path) symlink(ctx context.Context, from string, to string) error {
	return filepath.Walk(from, func(subFrom string, info fs.FileInfo, err error) error {
		rel, _ := filepath.Rel(from, subFrom)
		subTo := filepath.Join(to, rel)

		if info.IsDir() {
			if strings.Contains(subTo, "/cue.mod/gen/vendor/") {
				if err := forceSymlink(subFrom, subTo); err != nil {
					return err
				}
				return filepath.SkipDir
			}

			ok, err := hasSubDir(subFrom)
			if err != nil {
				return err
			}

			// If no sub dir, could be safe to add link
			if !ok {
				if err := forceSymlink(subFrom, subTo); err != nil {
					return err
				}
				return filepath.SkipDir
			}
			return nil
		}

		if strings.Contains(subFrom, "/cue.mod/") {
			return filepath.SkipDir
		}

		if filepath.Ext(subFrom) != ".cue" {
			return nil
		}

		return forceSymlink(subFrom, subTo)
	})
}

func hasSubDir(path string) (ok bool, err error) {
	err = filepath.Walk(path, func(sub string, info fs.FileInfo, err error) error {
		if sub == path {
			return nil
		}
		if info.IsDir() {
			ok = true
			return filepath.SkipDir
		}
		return nil
	})
	return
}

func forceSymlink(from, to string) error {
	if err := os.RemoveAll(to); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	return os.Symlink(from, to)
}

func (i *Path) shouldReplace() bool {
	return i.Replace != nil && i.Replace.From != i.Replace.Path
}

func (i *Path) ImportPath() string {
	if i.shouldReplace() {
		return i.Replace.From + filepath.Join(i.Module, i.SubPath)[len(i.Replace.Path):]
	}
	return filepath.Join(i.Module, i.SubPath)
}

func (i *Path) ImportPathDir() string {
	if i.shouldReplace() {
		return filepath.Join(i.Dir, i.SubPath)
	}
	return i.RepoRootDir()
}

func (i *Path) ResolvedImportPath() string {
	return filepath.Join(i.Dir, i.SubPath)
}

func (i *Path) RepoRootDir() string {
	if i.Repo == i.Module {
		return i.Dir
	}
	rel, _ := subPath(i.Repo, i.Module)
	return i.Dir[0 : len(i.Dir)-len("/"+rel)]
}
