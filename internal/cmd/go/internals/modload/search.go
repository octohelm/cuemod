// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/octohelm/cuemod/internal/cmd/go/internals/cfg"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/fsys"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/imports"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/modindex"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/par"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/search"
	"github.com/octohelm/cuemod/internal/cmd/go/internals/trace"
	"github.com/octohelm/cuemod/internal/cmd/internals/pkgpattern"

	"golang.org/x/mod/module"
)

type stdFilter int8

const (
	omitStd = stdFilter(iota)
	includeStd
)

// matchPackages is like m.MatchPackages, but uses a local variable (rather than
// a global) for tags, can include or exclude packages in the standard library,
// and is restricted to the given list of modules.
func matchPackages(ctx context.Context, m *search.Match, tags map[string]bool, filter stdFilter, modules []module.Version) {
	ctx, span := trace.StartSpan(ctx, "modload.matchPackages")
	defer span.Done()

	m.Pkgs = []string{}

	isMatch := func(string) bool { return true }
	treeCanMatch := func(string) bool { return true }
	if !m.IsMeta() {
		isMatch = pkgpattern.MatchPattern(m.Pattern())
		treeCanMatch = pkgpattern.TreeCanMatchPattern(m.Pattern())
	}

	var mu sync.Mutex
	have := map[string]bool{
		"builtin": true, // ignore pseudo-package that exists only for documentation
	}
	addPkg := func(p string) {
		mu.Lock()
		m.Pkgs = append(m.Pkgs, p)
		mu.Unlock()
	}
	if !cfg.BuildContext.CgoEnabled {
		have["runtime/cgo"] = true // ignore during walk
	}

	type pruning int8
	const (
		pruneVendor = pruning(1 << iota)
		pruneGoMod
	)

	q := par.NewQueue(runtime.GOMAXPROCS(0))

	walkPkgs := func(root, importPathRoot string, prune pruning) {
		_, span := trace.StartSpan(ctx, "walkPkgs "+root)
		defer span.Done()

		root = filepath.Clean(root)
		err := fsys.Walk(root, func(path string, fi fs.FileInfo, err error) error {
			if err != nil {
				m.AddError(err)
				return nil
			}

			want := true
			elem := ""

			// Don't use GOROOT/src but do walk down into it.
			if path == root {
				if importPathRoot == "" {
					return nil
				}
			} else {
				// Avoid .foo, _foo, and testdata subdirectory trees.
				_, elem = filepath.Split(path)
				if strings.HasPrefix(elem, ".") || strings.HasPrefix(elem, "_") || elem == "testdata" {
					want = false
				}
			}

			name := importPathRoot + filepath.ToSlash(path[len(root):])
			if importPathRoot == "" {
				name = name[1:] // cut leading slash
			}
			if !treeCanMatch(name) {
				want = false
			}

			if !fi.IsDir() {
				if fi.Mode()&fs.ModeSymlink != 0 && want && strings.Contains(m.Pattern(), "...") {
					if target, err := fsys.Stat(path); err == nil && target.IsDir() {
						fmt.Fprintf(os.Stderr, "warning: ignoring symlink %s\n", path)
					}
				}
				return nil
			}

			if !want {
				return filepath.SkipDir
			}
			// Stop at module boundaries.
			if (prune&pruneGoMod != 0) && path != root {
				if fi, err := os.Stat(filepath.Join(path, "go.mod")); err == nil && !fi.IsDir() {
					return filepath.SkipDir
				}
			}

			if !have[name] {
				have[name] = true
				if isMatch(name) {
					q.Add(func() {
						if _, _, err := scanDir(root, path, tags); err != imports.ErrNoGo {
							addPkg(name)
						}
					})
				}
			}

			if elem == "vendor" && (prune&pruneVendor != 0) {
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			m.AddError(err)
		}
	}

	// Wait for all in-flight operations to complete before returning.
	defer func() {
		<-q.Idle()
		sort.Strings(m.Pkgs) // sort everything we added for determinism
	}()

	if filter == includeStd {
		walkPkgs(cfg.GOROOTsrc, "", pruneGoMod)
		if treeCanMatch("cmd") {
			walkPkgs(filepath.Join(cfg.GOROOTsrc, "cmd"), "cmd", pruneGoMod)
		}
	}

	if cfg.BuildMod == "vendor" {
		mod := MainModules.mustGetSingleMainModule()
		if modRoot := MainModules.ModRoot(mod); modRoot != "" {
			walkPkgs(modRoot, MainModules.PathPrefix(mod), pruneGoMod|pruneVendor)
			walkPkgs(filepath.Join(modRoot, "vendor"), "", pruneVendor)
		}
		return
	}

	for _, mod := range modules {
		if !treeCanMatch(mod.Path) {
			continue
		}

		var (
			root, modPrefix string
			isLocal         bool
		)
		if MainModules.Contains(mod.Path) {
			if MainModules.ModRoot(mod) == "" {
				continue // If there is no main module, we can't search in it.
			}
			root = MainModules.ModRoot(mod)
			modPrefix = MainModules.PathPrefix(mod)
			isLocal = true
		} else {
			var err error
			root, isLocal, err = fetch(ctx, mod)
			if err != nil {
				m.AddError(err)
				continue
			}
			modPrefix = mod.Path
		}
		if mi, err := modindex.GetModule(root); err == nil {
			walkFromIndex(mi, modPrefix, isMatch, treeCanMatch, tags, have, addPkg)
			continue
		} else if !errors.Is(err, modindex.ErrNotIndexed) {
			m.AddError(err)
		}

		prune := pruneVendor
		if isLocal {
			prune |= pruneGoMod
		}
		walkPkgs(root, modPrefix, prune)
	}

	return
}

// walkFromIndex matches packages in a module using the module index. modroot
// is the module's root directory on disk, index is the modindex.Module for the
// module, and importPathRoot is the module's path prefix.
func walkFromIndex(index *modindex.Module, importPathRoot string, isMatch, treeCanMatch func(string) bool, tags, have map[string]bool, addPkg func(string)) {
	index.Walk(func(reldir string) {
		// Avoid .foo, _foo, and testdata subdirectory trees.
		p := reldir
		for {
			elem, rest, found := strings.Cut(p, string(filepath.Separator))
			if strings.HasPrefix(elem, ".") || strings.HasPrefix(elem, "_") || elem == "testdata" {
				return
			}
			if found && elem == "vendor" {
				// Ignore this path if it contains the element "vendor" anywhere
				// except for the last element (packages named vendor are allowed
				// for historical reasons). Note that found is true when this
				// isn't the last path element.
				return
			}
			if !found {
				// Didn't find the separator, so we're considering the last element.
				break
			}
			p = rest
		}

		// Don't use GOROOT/src.
		if reldir == "" && importPathRoot == "" {
			return
		}

		name := path.Join(importPathRoot, filepath.ToSlash(reldir))
		if !treeCanMatch(name) {
			return
		}

		if !have[name] {
			have[name] = true
			if isMatch(name) {
				if _, _, err := index.Package(reldir).ScanDir(tags); err != imports.ErrNoGo {
					addPkg(name)
				}
			}
		}
	})
}

// MatchInModule identifies the packages matching the given pattern within the
// given module version, which does not need to be in the build list or module
// requirement graph.
//
// If m is the zero module.Version, MatchInModule matches the pattern
// against the standard library (std and cmd) in GOROOT/src.
func MatchInModule(ctx context.Context, pattern string, m module.Version, tags map[string]bool) *search.Match {
	match := search.NewMatch(pattern)
	if m == (module.Version{}) {
		matchPackages(ctx, match, tags, includeStd, nil)
	}

	LoadModFile(ctx) // Sets Target, needed by fetch and matchPackages.

	if !match.IsLiteral() {
		matchPackages(ctx, match, tags, omitStd, []module.Version{m})
		return match
	}

	root, isLocal, err := fetch(ctx, m)
	if err != nil {
		match.Errs = []error{err}
		return match
	}

	dir, haveGoFiles, err := dirInModule(pattern, m.Path, root, isLocal)
	if err != nil {
		match.Errs = []error{err}
		return match
	}
	if haveGoFiles {
		if _, _, err := scanDir(root, dir, tags); err != imports.ErrNoGo {
			// ErrNoGo indicates that the directory is not actually a Go package,
			// perhaps due to the tags in use. Any other non-nil error indicates a
			// problem with one or more of the Go source files, but such an error does
			// not stop the package from existing, so it has no impact on matching.
			match.Pkgs = []string{pattern}
		}
	}
	return match
}
