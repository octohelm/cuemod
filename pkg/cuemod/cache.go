package cuemod

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/octohelm/cuemod/pkg/modutil"
	"github.com/pkg/errors"
	"golang.org/x/mod/module"
	"golang.org/x/tools/go/vcs"
)

func newCache() *cache {
	return &cache{
		replace:      map[modfile.PathMayWithVersion]replaceTargetWithMod{},
		mods:         map[module.Version]*Mod{},
		repoVersions: map[string]modfile.ModVersion{},
	}
}

type cache struct {
	replace map[modfile.PathMayWithVersion]replaceTargetWithMod
	// { [<module>@<version>]: *Mod }
	mods map[module.Version]*Mod
	// { [<repo>]:latest-version }
	repoVersions map[string]modfile.ModVersion
}

const ModSumFilename = "cue.mod/module.sum"

func (c *cache) ModuleSum() []byte {
	buf := bytes.NewBuffer(nil)

	moduleVersions := make([]module.Version, 0)

	for moduleVersion := range c.mods {
		moduleVersions = append(moduleVersions, moduleVersion)
	}

	module.Sort(moduleVersions)

	for _, n := range moduleVersions {
		m := c.mods[n]

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
	moduleVersion := module.Version{Path: mod.Module, Version: mod.Version}

	if mod.Repo == "" {
		mod.Repo = mod.Module
	}

	c.mods[moduleVersion] = mod

	// cached moduel@tag too
	if mod.VcsVersion != "" {
		c.mods[module.Version{Path: mod.Module, Version: mod.VcsVersion}] = mod
	}

	c.SetRepoVersion(mod.Repo, mod.ModVersion)

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

	if mod, ok := c.mods[module.Version{Path: repo, Version: version}]; ok && mod.Resolved() {
		root = mod
	} else {
		m, err := c.downloadIfNeed(ctx, repo, version)
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

		for _, m := range importPaths {
			if m == root.Module {
				break
			}

			if mod, ok := c.mods[module.Version{Path: m, Version: version}]; ok && mod.Resolved() {
				return mod, nil
			} else {
				rel, _ := subPath(root.Module, m)

				sub := Mod{}
				sub.Repo = root.Repo
				sub.Sum = root.Sum

				sub.Module = m
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

func (cache) downloadIfNeed(ctx context.Context, pkg string, version string) (*Mod, error) {

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
