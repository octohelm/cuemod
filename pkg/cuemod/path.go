package cuemod

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/module"

	"github.com/octohelm/cuemod/pkg/cueify"
	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/pkg/errors"
)

func PathFor(mod *Mod, importPath string) *Path {
	i := &Path{}
	i.Mod = *mod

	if importPath != "" {
		// when use replace with custom forked repo
		// mod.Module may not same as mod.Repo
		// then should to overwrite this path.Module with mode.Repo
		if !strings.HasPrefix(importPath, i.Module) {
			if strings.HasPrefix(importPath, i.Repo) {
				i.Module = i.Repo
			}
		}

		i.Module = i.ModuleRoot()
		i.SubPath, _ = subDir(i.Module, importPath)
	}

	return i
}

type Path struct {
	Mod
	SubPath string
	Replace *ReplaceRule
}

func (i Path) String() string {
	s := strings.Builder{}

	if i.Replace != nil {
		s.WriteString(i.Replace.From)

		if i.SubPath != "" {
			s.WriteString("/")
			s.WriteString(i.SubPath)
		}

		s.WriteString(" => ")
	}

	s.WriteString(i.Module)

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

	resolveDest := func(base string, withPathMajor bool, genSource bool) string {
		// for generated code
		//
		// pkg/
		//   .cuemod/<import_path> -> <GOMODCACHE>/<import_path>@<version>/
		if genSource {
			return filepath.Join(root, "cue.mod/pkg/.cuemod", base)
		}
		// support multi versions of pkg
		//
		// gen/
		//   <import_path>/v2  -> <GOMODCACHE>/<import_path>@<version>/
		// pkg/
		//   <import_path> -> <GOMODCACHE>/<import_path>@<version>/
		if withPathMajor {
			return filepath.Join(root, "cue.mod/gen", base)
		}
		return filepath.Join(root, "cue.mod/pkg", base)
	}

	_, pathMajor, ok := module.SplitPathVersion(i.ImportPathRoot())
	withPathMajor := ok && pathMajor != ""

	if err := i.symlink(ctx, i.ImportPathRootDir(), resolveDest(i.ImportPathRoot(), withPathMajor, gen != "")); err != nil {
		return err
	}

	if gen != "" {
		if err := cueify.ExtractToDir(
			ctx,
			gen,
			resolveDest(i.ImportPath(), withPathMajor, gen != ""),
			resolveDest(i.ImportPath(), withPathMajor, false),
		); err != nil {
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
	dir := i.Dir
	if i.SubPath != "" {
		dir = filepath.Join(dir, i.SubPath)
	}
	return dir[0 : len(dir)-len("/"+rel)]
}
