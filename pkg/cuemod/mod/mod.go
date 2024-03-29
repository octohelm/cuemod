package mod

import (
	"context"
	"os"
	"path/filepath"

	"golang.org/x/mod/semver"

	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/pkg/errors"
)

type Mod struct {
	modfile.ModFile
	modfile.ModVersion

	// Repo where module in vcs root
	Repo string
	// SubPath mod local sub path
	SubPath string
	// Dir module absolute dir
	Dir string
	// Root means this import path is mod root
	Root bool
	// Sum repo absolute dir sum
	Sum string
}

func (m *Mod) ModuleRoot() string {
	if m.SubPath != "" {
		return filepath.Join(m.Module, m.SubPath)
	}
	return m.Module
}

func (m *Mod) String() string {
	if m.Version == "" {
		return m.Module + "@v0.0.0"
	}
	return m.Module + "@" + m.Version
}

func (m *Mod) LoadInfo(ctx context.Context) (bool, error) {
	if m.Dir == "" || m.Dir[0] != '/' {
		return false, nil
	}

	if _, err := os.Stat(m.Dir); os.IsNotExist(err) {
		return false, errors.Wrapf(err, "%s not found", m.Dir)
	}

	exists, err := modfile.LoadModFile(m.Dir, &m.ModFile)
	if err != nil {
		return false, err
	}

	if exists {
		// module name should be from module.cue
		m.Module = m.ModFile.Module
		m.Root = true
	}

	return exists, nil
}

func (m *Mod) Resolved() bool {
	return m.Dir != ""
}

func (m *Mod) SetRequire(module string, modVersion modfile.ModVersion, indirect bool) {
	if module == m.Module {
		return
	}

	if m.Require == nil {
		m.Require = map[string]modfile.Requirement{}
	}

	r := modfile.Requirement{}

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

	if currentReplace, ok := m.Replace[modfile.VersionedPathIdentity{Path: module}]; ok {
		if currentReplace.IsLocalReplace() || currentReplace.Import != "" {
			return
		}

		currentReplace.Version = r.ModVersion.Version
		m.Replace[modfile.VersionedPathIdentity{Path: module}] = currentReplace
	}
}

func (m *Mod) FixVersion(repo string, version string) string {
	if m.Require != nil {
		if r, ok := m.Require[repo]; ok {
			if r.VcsRef != "" && r.Version == "v0.0.0" {
				return r.VcsRef
			}
			return r.Version
		}
	}
	return version
}

func versionGreaterThan(v string, w string) bool {
	if w == "" {
		return false
	}
	return semver.Compare(v, w) > 0
}
