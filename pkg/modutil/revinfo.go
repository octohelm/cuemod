package modutil

import (
	"context"
	"fmt"
	"time"
	_ "unsafe"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch/codehost"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/vcs"
)

//go:linkname newCodeRepo github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch.newCodeRepo
func newCodeRepo(code codehost.Repo, codeRoot, path string) (modfetch.Repo, error)

//go:linkname lookupCodeRepo github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch.lookupCodeRepo
func lookupCodeRepo(ctx context.Context, rr *vcs.RepoRoot) (codehost.Repo, error)

func finalLookupCodeRepo(ctx context.Context, rr *vcs.RepoRoot, localOk bool) (codehost.Repo, error) {
	if rr.VCS.Name == "Git" && localOk {
		return codehost.LocalGitRepo(ctx, rr.Root)
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
		return nil, err
	}

	head, err := c.Status(c, rootDir)
	if err != nil {
		return nil, err
	}

	rr := &vcs.RepoRoot{}
	rr.VCS = c
	rr.Root = rootDir
	rr.Repo = repo

	code, err := finalLookupCodeRepo(ctx, rr, true)
	if err != nil {
		return nil, err
	}

	r, err := newCodeRepo(code, rr.Root, rr.Root)
	if err != nil {
		return nil, err
	}

	info, err := r.Stat(ctx, head.Revision)
	if err != nil {
		return nil, err
	}

	version := info.Version
	base, err := module.PseudoVersionBase(version)
	if err == nil {
		version = base
	}
	if head.Uncommitted {
		version += "-dirty"
	}

	info.Version = PseudoVersion(version, info.Time, info.Short, !module.IsPseudoVersion(info.Version))

	return info, nil
}

func PseudoVersion(version string, t time.Time, rev string, exact bool) string {
	major := semver.Major(version)
	if major == "" {
		major = "v0"
	}

	if exact {
		build := semver.Build(version)
		segment := fmt.Sprintf("%s-%s", t.UTC().Format(module.PseudoVersionTimestampFormat), rev)
		version = semver.Canonical(version)
		if version == "" {
			version = major + ".0.0"
		}
		return version + "-" + segment + build
	}

	return module.PseudoVersion(
		major,
		version,
		t,
		rev,
	)
}
