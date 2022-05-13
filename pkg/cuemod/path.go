package cuemod

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

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
			}
		}
		i.SubPath, _ = subDir(i.Module, importPath)
		if i.SubPath != "" {
			i.Dir = path.Join(i.Dir, i.SubPath)
		}
	}

	return i
}

type Path struct {
	*Mod
	Module  string
	SubPath string
	Replace *ReplaceRule
}

func (i Path) String() string {
	s := strings.Builder{}

	if i.Replace != nil {
		s.WriteString(i.Replace.From)
		s.WriteString(" => ")
	}

	s.WriteString(i.Mod.String())

	if i.SubPath != "" {
		s.WriteString("/")
		s.WriteString(i.SubPath)
	}

	return s.String()
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

	if !i.Mod.Root {
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
	pkgRootPath := filepath.Join(root, pkgRoot, i.ImportPathRoot())

	if err := i.symlink(ctx, i.ImportPathRootDir(), pkgRootPath); err != nil {
		return err
	}

	if gen != "" {
		importPathDir := filepath.Join(root, pkgRoot, importPath)

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
	if _, err := os.Stat(from); err != nil {
		return errors.Wrapf(err, "invalid path %s", from)
	}
	return forceSymlink(from, to)
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

func (i *Path) ResolvedImportPath() string {
	return filepath.Join(i.Dir, i.SubPath)
}

func (i Path) ImportPathRoot() string {
	if i.Replace != nil {
		return i.Replace.From
	}
	return i.Repo
}

func (i Path) ImportPathRootDir() string {
	ipr := i.ImportPathRoot()
	ip := i.ImportPath()
	if ip == ipr {
		return i.Dir
	}
	rel, _ := subDir(ipr, ip)
	return i.Dir[0 : len(i.Dir)-len("/"+rel)]
}
