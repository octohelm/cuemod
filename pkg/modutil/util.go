package modutil

import (
	"context"
	"strings"

	"golang.org/x/mod/module"

	"github.com/octohelm/cuemod/internal/cmd/go/internals/cfg"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/modfetch"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/modload"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/vcs"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/web"
)

type Module struct {
	Path    string
	Version string
	Error   string
	Dir     string
	Sum     string
}

func envOr(key, def string) string {
	val := cfg.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

func init() {
	cfg.GOPROXY = envOr("GOPROXY", "https://proxy.golang.org,direct")
	cfg.GOSUMDB = envOr("GOSUMDB", "sum.golang.org")
}

func RepoRootForImportPath(importPath string) (string, error) {
	r, err := vcs.RepoRootForImportPath(importPath, vcs.IgnoreMod, web.DefaultSecurity)
	if err != nil {
		return "", err
	}
	return r.Root, nil
}

// Get Module
func Get(ctx context.Context, path string, version string, verbose bool) *Module {
	mv := module.Version{Path: path, Version: version}
	p, err := modfetch.DownloadDir(ctx, mv)
	if err == nil {
		// found in cache
		return &Module{
			Path:    mv.Path,
			Version: mv.Version,
			Dir:     p,
			Sum:     modfetch.Sum(ctx, mv),
		}
	}

	modload.ForceUseModules = true
	cfg.BuildX = verbose

	requestPath := path + "@" + version
	if strings.HasSuffix(path, ".v3") {
		requestPath = path
	}

	found, err := modload.ListModules(ctx, []string{requestPath}, modload.ListVersions, "")
	if err != nil {
		panic(err)
	}
	if len(found) > 0 {
		info := found[0]

		m := &Module{
			Path:    info.Path,
			Version: info.Version,
		}

		if info.Error != nil {
			m.Error = info.Error.Err
		} else {
			m.Dir = info.Dir
			m.Sum = modfetch.Sum(ctx, module.Version{Path: m.Path, Version: m.Version})
		}
		return m
	}
	return nil
}

// Download Module
func Download(ctx context.Context, m *Module) {
	mv := module.Version{Path: m.Path, Version: m.Version}
	dir, err := modfetch.DownloadDir(ctx, mv)
	if err == nil {
		// found in cache
		m.Dir = dir
		m.Sum = modfetch.Sum(ctx, module.Version{Path: m.Path, Version: m.Version})
		return
	}

	dir, err = modfetch.Download(ctx, mv)
	if err != nil {
		m.Error = err.Error()
		return
	}
	m.Dir = dir
	m.Sum = modfetch.Sum(ctx, module.Version{Path: m.Path, Version: m.Version})
}
