package cuemod

import (
	"context"
	"os"
	"path/filepath"

	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/octohelm/cuemod/pkg/extractor"
)

func PathFor(mod *Mod, importPath string) *Path {
	i := &Path{}
	i.Mod = mod

	if importPath != "" {
		i.SubPath, _ = subPath(i.Module, importPath)
	}

	return i
}

type Path struct {
	*Mod
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

func (i *Path) SymlinkOrGen(ctx context.Context, root string) error {
	pkgRoot := "cue.mod/pkg"

	if root == i.Dir {
		pkgRoot = "cue.mod/usr"
	}

	gen := ""

	if i.Lang != "" {
		gen = i.Lang
	}

	if i.Replace != nil && i.Replace.Gen != "" {
		gen = i.Replace.Gen
	}

	if gen != "" {
		pkgRoot = "cue.mod/gen/_src"
	}

	importPath := i.ImportPath()

	repoRootDir := filepath.Join(root, pkgRoot, i.RepoRoot)
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
	if err := os.RemoveAll(to); err != nil {
		return err
	}
	// make sure parent created
	if err := os.MkdirAll(filepath.Dir(to), 0777); err != nil {
		return err
	}
	if err := os.Symlink(from, to); err != nil {
		return err
	}
	return nil
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
	if i.RepoRoot == i.Module {
		return i.Dir
	}
	rel, _ := subPath(i.RepoRoot, i.Module)
	return i.Dir[0 : len(i.Dir)-len("/"+rel)]
}
