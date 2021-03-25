package cuemod

import (
	"context"
	"os"
	"path/filepath"

	"github.com/octohelm/cuemod/pkg/extractor"

	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/pkg/errors"
)

type Mod struct {
	modfile.ModFile

	modfile.ModVersion

	// Repo where module in vcs root
	Repo string
	// Sum repo absolute dir sum
	Sum string
	// Lang if set, will trigger extractor like `go`
	Lang string
	// Dir module absolute dir
	Dir string
}

func (m *Mod) String() string {
	if m.Version == "" {
		return m.Module + "@v0.0.0"
	}
	return m.Module + "@" + m.Version
}

func (m *Mod) LoadInfo(ctx context.Context) (bool, error) {
	if _, err := os.Stat(m.Dir); os.IsNotExist(err) {
		return false, errors.Wrapf(err, "%s not found", m.Dir)
	}

	modfileExists, err := modfile.LoadModFile(m.Dir, &m.ModFile)
	if err != nil {
		return false, err
	}

	// auto detect
	lang, deps := extractor.Detect(ctx, m.Dir)

	if lang != "" {
		m.Lang = lang
		modfileExists = true

		for repo, version := range deps {
			if m.Replace == nil {
				m.Replace = map[modfile.PathMayWithVersion]modfile.ReplaceTarget{}
			}

			m.Replace[modfile.PathMayWithVersion{Path: repo}] = modfile.ReplaceTarget{
				PathMayWithVersion: modfile.PathMayWithVersion{Path: repo, Version: version},
			}
		}
	}

	return modfileExists, nil
}

func (m *Mod) Resolved() bool {
	return m.Dir != ""
}

func (m *Mod) SetRequire(module string, modVersion modfile.ModVersion, indirect bool) {
	if module == m.Module {
		return
	}

	if m.Require == nil {
		m.Require = map[string]modfile.Require{}
	}

	r := modfile.Require{}
	r.ModVersion = modVersion
	r.Indirect = indirect

	if currentRequire, ok := m.Require[module]; ok {
		// always using greater one
		if versionGreaterThan(currentRequire.Version, r.Version) {
			r.ModVersion = currentRequire.ModVersion
		}

		if r.Indirect {
			r.Indirect = currentRequire.Indirect
		}
	}

	m.Require[module] = r
}

func (m *Mod) ResolveImportPath(ctx context.Context, cache *cache, importPath string, version string) (*Path, error) {
	// self import '<mod.module>/dir/to/sub'
	if isSubDirFor(importPath, m.Module) {
		return PathFor(m, importPath), nil
	}

	if matched, replaceTarget, ok := cache.LookupReplace(importPath); ok {
		// xxx => ../xxx
		if replaceTarget.IsLocalReplace() {
			mod := &Mod{Dir: filepath.Join(m.Dir, replaceTarget.Path)}

			mod.Version = "v0.0.0"
			if _, err := mod.LoadInfo(ctx); err != nil {
				return nil, err
			}

			if mod.Module == "" {
				mod.Module = filepath.Join(m.Module, replaceTarget.Path)
			}

			cache.Collect(ctx, mod)
			return PathFor(mod, importPath), nil
		}

		// a[@latest] => b@latest
		// must strict version
		replacedImportPath := replaceImportPath(replaceTarget.Path, matched.Path, importPath)

		ctxWithUpgradeDisabled := WithOpts(ctx, OptUpgrade(false))

		fixVersion := m.fixVersion

		if replaceTarget.Version != "" {
			fixVersion = nil
		}
		mod, err := cache.Get(ctxWithUpgradeDisabled, replacedImportPath, replaceTarget.Version, fixVersion)
		if err != nil {
			return nil, err
		}
		return PathFor(mod, replacedImportPath).WithReplace(matched.Path, replaceTarget), nil
	}

	mod, err := cache.Get(ctx, importPath, version, m.fixVersion)
	if err != nil {
		return nil, err
	}

	return PathFor(mod, importPath), nil
}

func (m *Mod) fixVersion(repo string, version string) string {
	if m.Require != nil {
		if r, ok := m.Require[repo]; ok {
			return r.Version
		}
	}
	return version
}
