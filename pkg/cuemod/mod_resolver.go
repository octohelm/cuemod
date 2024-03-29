package cuemod

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/octohelm/cuemod/pkg/cuemod/stdlib"
	"github.com/octohelm/cuemod/pkg/modutil"
	"github.com/pkg/errors"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

func newModResolver() *modResolver {
	return &modResolver{
		replace:      map[modfile.VersionedPathIdentity]replaceTargetWithMod{},
		mods:         map[module.Version]*Mod{},
		repoVersions: map[string]modfile.ModVersion{},
	}
}

type modResolver struct {
	root *Mod

	replace map[modfile.VersionedPathIdentity]replaceTargetWithMod
	// { [<module>@<version>]: *Mod }
	mods map[module.Version]*Mod
	// { [<repo>]:latest-version }
	repoVersions map[string]modfile.ModVersion
}

const ModSumFilename = "cue.mod/module.sum"

func (r *modResolver) ModuleSum() []byte {
	buf := bytes.NewBuffer(nil)

	moduleVersions := make([]module.Version, 0)

	for moduleVersion, m := range r.mods {
		if m.Root {
			moduleVersions = append(moduleVersions, moduleVersion)
		}
	}

	sort.Slice(moduleVersions, func(i, j int) bool {
		return moduleVersions[i].Path < moduleVersions[j].Path
	})

	for _, n := range moduleVersions {
		m := r.mods[n]

		if m.Version != "" {
			_, _ = fmt.Fprintf(buf, "%s %s %s\n", m.Module, m.Version, m.Sum)
		}
	}

	return buf.Bytes()
}

type replaceTargetWithMod struct {
	modfile.ReplaceTarget
	mod *Mod
}

func (r *modResolver) ResolveImportPath(ctx context.Context, root *Mod, importPath string, version string) (p *Path, e error) {
	// self import '<root.module>/dir/to/sub'
	if isSubDirFor(importPath, root.Module) {
		return PathFor(root, importPath), nil
	}

	if matched, replaceTarget, ok := r.LookupReplace(importPath); ok {
		// xxx => ../xxx
		// only works for root
		if replaceTarget.IsLocalReplace() {
			mod := &Mod{Dir: filepath.Join(root.Dir, replaceTarget.Path)}

			mod.Version = "v0.0.0"
			if _, err := mod.LoadInfo(ctx); err != nil {
				return nil, err
			}

			if mod.Module == "" {
				mod.Module = filepath.Join(root.Module, replaceTarget.Path)
			}

			r.Collect(ctx, mod)
			return PathFor(mod, importPath), nil
		}

		// a[@latest] => b@latest
		// must strict version
		replacedImportPath := replaceImportPath(replaceTarget.Path, matched.Path, importPath)

		// when version never upgrade for upgrade
		if replaceTarget.Exactly() {
			ctx = WithOpts(ctx, OptUpgrade(false))
		}

		mod, err := r.Get(ctx, replacedImportPath, replaceTarget.ModVersion)
		if err != nil {
			return nil, err
		}

		return PathFor(mod, replacedImportPath).WithReplace(matched.Path, replaceTarget), nil
	}

	mod, err := r.Get(ctx, importPath, modfile.ModVersion{Version: version})
	if err != nil {
		return nil, err
	}

	return PathFor(mod, importPath), nil
}

func (r *modResolver) LookupReplace(importPath string) (matched modfile.VersionedPathIdentity, replace modfile.ReplaceTarget, exists bool) {
	for _, path := range paths(importPath) {
		p := modfile.VersionedPathIdentity{Path: path}
		if rp, ok := r.replace[p]; ok {
			return p, rp.ReplaceTarget, true
		}
	}
	return modfile.VersionedPathIdentity{}, modfile.ReplaceTarget{}, false
}

func (r *modResolver) Collect(ctx context.Context, mod *Mod) {
	if r.root == nil {
		r.root = mod
	}

	moduleVersion := module.Version{Path: mod.Module, Version: mod.Version}

	if mod.Repo == "" {
		mod.Repo = mod.Module
	}

	r.mods[moduleVersion] = mod

	// cached moduel@tag too
	if mod.VcsRef != "" {
		r.mods[module.Version{Path: mod.Module, Version: mod.VcsRef}] = mod
	}

	r.SetRepoVersion(mod.Repo, mod.ModVersion)

	for repo, req := range mod.Require {
		r.SetRepoVersion(repo, req.ModVersion)
	}

	for k, replaceTarget := range mod.Replace {
		// only work for root mod
		if replaceTarget.IsLocalReplace() && mod != r.root {
			return
		}

		// never modify replaced
		if _, ok := r.replace[k]; !ok {
			r.replace[k] = replaceTargetWithMod{mod: mod, ReplaceTarget: replaceTarget}
		}
	}
}

func (r *modResolver) SetRepoVersion(module string, version modfile.ModVersion) {
	if mv, ok := r.repoVersions[module]; ok {
		if mv.Version == "" {
			mv.Version = version.Version
		} else if versionGreaterThan(version.Version, mv.Version) {
			mv.Version = version.Version
		}

		// sync tag version
		if version.VcsRef != "" {
			mv.VcsRef = version.VcsRef
		}

		r.repoVersions[module] = mv
	} else {
		r.repoVersions[module] = version
	}
}

func (r *modResolver) RepoVersion(repo string) modfile.ModVersion {
	if v, ok := r.repoVersions[repo]; ok {
		return v
	}
	return modfile.ModVersion{}
}

func (r *modResolver) Get(ctx context.Context, pkgImportPath string, modVersion modfile.ModVersion) (*Mod, error) {
	repo, err := r.repoRoot(ctx, pkgImportPath)
	if err != nil {
		return nil, err
	}
	return r.get(ctx, repo, modVersion, pkgImportPath)
}

func (r *modResolver) get(ctx context.Context, repo string, requestedVersion modfile.ModVersion, importPath string) (*Mod, error) {
	// fix /v2
	if p, m, ok := module.SplitPathVersion(repo); ok {
		if requestedVersion.VcsRef == "" {
			requestedVersion.VcsRef = m
		}
	} else {
		repo = p
	}

	if requestedVersion.VcsRef == "" && requestedVersion.Version != "" && !semver.IsValid(requestedVersion.Version) {
		requestedVersion.VcsRef = requestedVersion.Version
		requestedVersion.Version = ""
	}

	if requestedVersion.VcsRef == "" {
		requestedVersion.VcsRef = "latest"
	}

	version := requestedVersion.Version

	forUpgrade := OptsFromContext(ctx).Upgrade

	if forUpgrade || version == "" {
		if requestedVersion.VcsRef != "" {
			version = requestedVersion.VcsRef
		}
	}

	if !forUpgrade {
		if mv, ok := r.repoVersions[repo]; ok {
			if mv.Version != "" {
				version = mv.Version
			}
		}
	}

	var root *Mod

	if mod, ok := r.mods[module.Version{Path: repo, Version: version}]; ok && mod.Resolved() {
		// resolved
		root = mod
	} else {
		m, err := r.resolveMod(ctx, repo, version)
		if err != nil {
			return nil, err
		}

		if requestedVersion.VcsRef != "latest" {
			m.VcsRef = requestedVersion.VcsRef
		}

		root = m
		// fetched repo always root
		root.Root = true

		if _, err := root.LoadInfo(ctx); err != nil {
			return nil, err
		}

		r.Collect(ctx, root)
	}

	if root != nil {
		// sub dir may a mod.
		importPaths := paths(importPath)

		for _, m := range importPaths {
			if m == root.Module {
				break
			}

			if mod, ok := r.mods[module.Version{Path: m, Version: version}]; ok && mod.Resolved() {
				return mod, nil
			} else {
				subPath, _ := subDir(root.Module, m)

				sub := Mod{}
				sub.Sum = root.Sum
				sub.Repo = root.Repo
				sub.SubPath = subPath

				sub.Module = m
				sub.ModVersion = root.ModVersion

				sub.Dir = filepath.Join(root.Dir, subPath)

				ok, err := sub.LoadInfo(ctx)
				if err != nil {
					// if subPath contains go.mod, it will be empty
					if os.IsNotExist(errors.Unwrap(err)) {
						return r.get(ctx, sub.Module, requestedVersion, importPath)
					}
					return nil, err
				}

				if ok {
					r.Collect(ctx, &sub)
					return &sub, nil
				}
			}
		}
	}

	return root, nil
}

func (r *modResolver) repoRoot(ctx context.Context, importPath string) (string, error) {
	importPaths := paths(importPath)

	for _, p := range importPaths {
		if _, ok := r.repoVersions[p]; ok {
			return p, nil
		}
	}

	logr.FromContext(ctx).Debug(fmt.Sprintf("resolve %s", importPath))

	var root string

	if stdrepo, ok := stdlib.RepoRootForImportPath(importPath); ok {
		root = stdrepo
	} else {
		re, err := modutil.RepoRootForImportPath(importPath)
		if err != nil {
			return "", errors.Wrapf(err, "resolve `%s` failed", importPath)
		}
		root = re
	}

	r.SetRepoVersion(root, modfile.ModVersion{})

	return root, nil
}

func (r *modResolver) resolveMod(ctx context.Context, pkg string, version string) (*Mod, error) {
	cuemModRoot := ""
	if r.root != nil {
		cuemModRoot = r.root.Dir
	}

	stdm, err := stdlib.Mount(ctx, pkg, cuemModRoot)
	if err != nil {
		return nil, err
	}

	if stdm != nil {
		return stdm, nil
	}

	info := modutil.Get(ctx, pkg, version, OptsFromContext(ctx).Verbose)
	if info == nil {
		return nil, fmt.Errorf("can't found %s@%s", pkg, version)
	}

	if info.Error != "" {
		return nil, errors.New(info.Error)
	}

	if info.Dir == "" {
		logr.FromContext(ctx).Debug(fmt.Sprintf("get %s@%s", pkg, version))
		modutil.Download(ctx, info)
		if info.Error != "" {
			return nil, errors.New(info.Error)
		}
	}

	mod := &Mod{}

	mod.Module = info.Path
	mod.Version = info.Version
	mod.Repo = info.Path
	mod.Sum = info.Sum
	mod.Dir = info.Dir

	return mod, nil
}
