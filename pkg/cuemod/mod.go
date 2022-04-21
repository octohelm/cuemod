package cuemod

import (
	"context"
	"os"

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

func (m *Mod) FixVersion(repo string, version string) string {
	if m.Require != nil {
		if r, ok := m.Require[repo]; ok {
			if r.VcsVersion != "" && r.Version == "v0.0.0" {
				return r.VcsVersion
			}
			return r.Version
		}
	}
	return version
}
