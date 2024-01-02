package modutil

import (
	"context"
	"path"
	"strings"
	_ "unsafe"

	"github.com/pkg/errors"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch/codehost"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/vcs"
	"github.com/octohelm/cuemod/pkg/version"
)

//go:linkname newCodeRepo github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch.newCodeRepo
func newCodeRepo(code codehost.Repo, codeRoot, path string) (modfetch.Repo, error)

//go:linkname lookupCodeRepo github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch.lookupCodeRepo
func lookupCodeRepo(ctx context.Context, rr *vcs.RepoRoot) (codehost.Repo, error)

func finalLookupCodeRepo(ctx context.Context, rr *vcs.RepoRoot, localOk bool) (codehost.Repo, error) {
	if strings.ToLower(rr.VCS.Name) == "git" && localOk {
		return codehost.LocalGitRepo(ctx, path.Join(rr.Root, ".git"))
	}
	return lookupCodeRepo(ctx, rr)
}

type RevInfo = modfetch.RevInfo

func RevInfoFromDir(ctx context.Context, dir string) (*RevInfo, error) {
	rootDir, c, err := vcs.FromDir(dir, "", true)
	if err != nil {
		return nil, err
	}

	repo, err := c.RemoteRepo(c, rootDir)
	if err != nil {
		return nil, errors.Wrap(err, "resolve remote repo failed")
	}

	head, err := c.Status(c, rootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "stat faield")
	}

	rr := &vcs.RepoRoot{}
	rr.VCS = c
	rr.Root = rootDir
	rr.Repo = repo

	code, err := finalLookupCodeRepo(ctx, rr, true)
	if err != nil {
		return nil, err
	}

	importPath := rr.Root

	data, err := code.ReadFile(ctx, head.Revision, "go.mod", -1)
	if err == nil {
		f, err := modfile.ParseLax("go.mod", data, nil)
		if err != nil {
			return nil, err
		}

		// <import_path>/v2
		_, pathMajor, ok := module.SplitPathVersion(f.Module.Mod.Path)
		if ok && pathMajor != "" {
			importPath += pathMajor
		}
	}

	r, err := newCodeRepo(code, rr.Root, importPath)
	if err != nil {
		return nil, err
	}

	info, err := r.Stat(ctx, head.Revision)
	if err != nil {
		return nil, errors.Wrapf(err, "stat faield")
	}

	info.Version = version.Convert(info.Version, info.Time, info.Short, head.Uncommitted)

	return info, nil
}
