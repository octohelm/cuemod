package cuemod

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/octohelm/cuemod/pkg/modfetch"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/vcs"
)

func newCache() *cache {
	return &cache{
		replace:      map[modfile.PathMayWithVersion]replaceTargetWithMod{},
		mods:         map[string]*Mod{},
		repoVersions: map[string]modfile.ModVersion{},
	}
}

type cache struct {
	replace map[modfile.PathMayWithVersion]replaceTargetWithMod
	// { [<module>@<version>]: *Mod }
	mods map[string]*Mod
	// { [<repo>]:latest-version }
	repoVersions map[string]modfile.ModVersion
}

type replaceTargetWithMod struct {
	modfile.ReplaceTarget
	mod *Mod
}

func (c *cache) LookupReplace(importPath string) (matched modfile.PathMayWithVersion, replace modfile.ReplaceTarget, exists bool) {
	for _, path := range paths(importPath) {
		p := modfile.PathMayWithVersion{Path: path}
		if rp, ok := c.replace[p]; ok {
			return p, rp.ReplaceTarget, true
		}
	}
	return modfile.PathMayWithVersion{}, modfile.ReplaceTarget{}, false
}

func (c *cache) Collect(ctx context.Context, mod *Mod) {
	id := mod.String()

	if mod.RepoRoot == "" {
		mod.RepoRoot = mod.Module
	}

	c.mods[id] = mod

	// cached moduel@tag too
	if mod.VcsVersion != "" {
		c.mods[mod.Module+"@"+mod.VcsVersion] = mod
	}

	c.SetRepoVersion(mod.RepoRoot, mod.ModVersion)

	for repo, r := range mod.Require {
		c.SetRepoVersion(repo, r.ModVersion)
	}

	for k, replaceTarget := range mod.Replace {
		if currentReplaceTarget, ok := c.replace[k]; !ok {
			c.replace[k] = replaceTargetWithMod{mod: mod, ReplaceTarget: replaceTarget}
		} else {
			if replaceTarget.String() != currentReplaceTarget.PathMayWithVersion.String() {
				fmt.Printf(`
[WARNING] '%s' already replaced to 
	'%s' (using by module '%s'), but another module want to replace as 
	'%s' (requested by module %s)

`,
					k,
					currentReplaceTarget.PathMayWithVersion, currentReplaceTarget.mod,
					replaceTarget, mod,
				)
			}
		}
	}
}

func (c *cache) SetRepoVersion(module string, version modfile.ModVersion) {
	if mv, ok := c.repoVersions[module]; ok {
		if mv.Version == "" {
			c.repoVersions[module] = version
		} else if versionGreaterThan(version.Version, mv.Version) {
			c.repoVersions[module] = version
		} else if version.Version == mv.Version && version.VcsVersion != "" {
			// sync tag version
			mv.VcsVersion = version.VcsVersion
			c.repoVersions[module] = mv
		}

	} else {
		c.repoVersions[module] = version
	}
}

func (c *cache) RepoVersion(repo string) modfile.ModVersion {
	if v, ok := c.repoVersions[repo]; ok {
		return v
	}
	return modfile.ModVersion{}
}

type VersionFixer = func(repo string, version string) string

func (c *cache) Get(ctx context.Context, pkgImportPath string, version string, fixVersion VersionFixer) (*Mod, error) {
	repo, err := c.repoRoot(ctx, pkgImportPath)
	if err != nil {
		return nil, err
	}

	if fixVersion != nil {
		version = fixVersion(repo, version)
	}

	return c.get(ctx, repo, version, pkgImportPath)
}

const versionUpgrade = "upgrade"

func (c *cache) get(ctx context.Context, repo string, requestedVersion string, importPath string) (*Mod, error) {
	version := requestedVersion

	if version == "" {
		version = versionUpgrade
	}

	if OptsFromContext(ctx).Upgrade {
		version = versionUpgrade

		// when tag version exists, should upgrade with tag version
		if mv, ok := c.repoVersions[repo]; ok {
			if mv.VcsVersion != "" {
				version = mv.VcsVersion
			}
		}
	} else {
		// use the resolved version, when already resolved.
		if mv, ok := c.repoVersions[repo]; ok {
			if mv.VcsVersion != "" && mv.Version != "" && mv.VcsVersion == requestedVersion {
				version = mv.Version
			}
		}
	}

	// mod@version replace
	if r, ok := c.replace[modfile.PathMayWithVersion{Path: repo, Version: version}]; ok {
		repo, version = r.Path, r.Version
	}

	if version == "" {
		version = versionUpgrade
	}

	var root *Mod

	if mod, ok := c.mods[repo+"@"+version]; ok && mod.Resolved() {
		root = mod
	} else {
		logr.FromContext(ctx).Debug(fmt.Sprintf("get %s@%s", repo, version))

		m, err := c.download(ctx, repo, version)
		if err != nil {
			return nil, err
		}

		if version != versionUpgrade {
			m.VcsVersion = requestedVersion
		}

		root = m

		if _, err := root.LoadInfo(ctx); err != nil {
			return nil, err
		}

		c.Collect(ctx, root)
	}

	if root != nil {
		// sub dir may as mod.
		importPaths := paths(importPath)

		for _, module := range importPaths {
			if module == root.Module {
				break
			}

			if mod, ok := c.mods[module+"@"+version]; ok && mod.Resolved() {
				return mod, nil
			} else {
				rel, _ := subPath(root.Module, module)

				sub := Mod{}
				sub.RepoRoot = root.RepoRoot
				sub.RepoSum = root.RepoSum

				sub.Module = module
				sub.ModVersion = root.ModVersion

				sub.Dir = filepath.Join(root.Dir, rel)

				ok, err := sub.LoadInfo(ctx)
				if err != nil {
					// if dir contains go.mod, will be empty
					if os.IsNotExist(errors.Unwrap(err)); err != nil {
						return c.get(ctx, sub.Module, version, importPath)
					}
					return nil, err
				}

				if ok {
					c.Collect(ctx, &sub)
					return &sub, nil
				}
			}
		}
	}

	return root, nil
}

func (c *cache) repoRoot(ctx context.Context, importPath string) (string, error) {
	importPaths := paths(importPath)

	for _, p := range importPaths {
		if _, ok := c.repoVersions[p]; ok {
			return p, nil
		}
	}

	logr.FromContext(ctx).Debug(fmt.Sprintf("resolve %s", importPath))

	r, err := vcs.RepoRootForImportPath(importPath, true)
	if err != nil {
		return "", errors.Wrapf(err, "resolve `%s` failed", importPath)
	}

	c.SetRepoVersion(r.Root, modfile.ModVersion{})

	return r.Root, nil
}

func (cache) download(ctx context.Context, pkg string, version string) (*Mod, error) {
	info := modfetch.Fetch(ctx, pkg, version, OptsFromContext(ctx).Verbose)
	if info == nil {
		return nil, fmt.Errorf("can't found %s@%s", pkg, version)
	}

	if info.Error != "" {
		return nil, errors.New(info.Error)
	}

	mod := &Mod{}

	mod.Module = info.Path
	mod.RepoRoot = info.Path
	mod.Version = info.Version
	mod.RepoSum = info.Sum
	mod.Dir = info.Dir

	return mod, nil
}
