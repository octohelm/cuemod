package cuemod

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue/build"
	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/pkg/errors"
)

func RuntimeFor(root string) *Runtime {
	vm := &Runtime{
		cache: newCache(),
	}

	if !filepath.IsAbs(root) {
		cwd, _ := os.Getwd()
		root = filepath.Join(cwd, root)
	}

	mod := &Mod{}
	mod.Dir = root

	ctx := context.Background()

	if _, err := mod.LoadInfo(ctx); err != nil {
		panic(err)
	}

	vm.mod = mod

	vm.cache.Collect(ctx, mod)

	return vm
}

type Runtime struct {
	mod   *Mod
	cache *cache
}

func (r *Runtime) Cleanup() error {
	if err := os.RemoveAll(filepath.Join(r.CueModRoot(), "gen")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(r.CueModRoot(), "pkg")); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) completePath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(r.mod.Dir, p)
}

func (r *Runtime) ListCue(fromPath string) ([]string, error) {
	files := make([]string, 0)

	walkSubDir := strings.HasSuffix(fromPath, "/...")

	if walkSubDir {
		fromPath = fromPath[0 : len(fromPath)-4]
	}

	start := filepath.Join(r.mod.Dir, fromPath)

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if path == start {
			return nil
		}

		// skip cue.mod
		if isSubDirFor(path, r.CueModRoot()) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			// skip sub dir which is cuemod root
			if _, err := os.Stat(filepath.Join(path, modfile.ModFilename)); err == nil {
				return filepath.SkipDir
			}

			if walkSubDir {
				return nil
			}
			return filepath.SkipDir
		}

		if filepath.Ext(path) == ".cue" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *Runtime) Eval(ctx context.Context, filename string, encoding Encoding) ([]byte, error) {
	filename = r.completePath(filename)
	inst := r.build(ctx, filename)
	results, err := Eval(inst, encoding)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Runtime) Get(ctx context.Context, i string) error {
	if i[0] == '.' {
		return r.autoImport(ctx, i)
	}
	return r.download(ctx, i)
}

func (r *Runtime) Resolve(ctx context.Context, importPath string, importedAt string) (string, error) {
	resolvedImportPath, err := r.mod.ResolveImportPath(ctx, r.cache, importPath, "")
	if err != nil {
		return "", errors.Wrapf(err, "resolve import `%s` failed", importPath)
	}

	indirect := isSubDirFor(importedAt, r.CueModRoot()) && !isSubDirFor(importedAt, filepath.Join(r.CueModRoot(), "usr", r.mod.Module))

	if err := r.setRequireFromImportPath(ctx, resolvedImportPath, indirect); err != nil {
		return "", err
	}

	dir := resolvedImportPath.ResolvedImportPath()
	return dir, nil
}

func (r *Runtime) CueModRoot() string {
	return filepath.Join(r.mod.Dir, "cue.mod")
}

func (r *Runtime) setRequireFromImportPath(ctx context.Context, p *Path, indirect bool) error {
	modVersion := p.ModVersion

	if mv := r.cache.RepoVersion(p.Repo); mv.Version != "" {
		modVersion = mv
	}

	if err := p.SymlinkOrImport(ctx, r.mod.Dir); err != nil {
		return err
	}

	r.mod.SetRequire(p.Repo, modVersion, indirect)

	return r.syncFiles()
}

func (r *Runtime) syncFiles() error {
	if err := writeFile(filepath.Join(r.mod.Dir, modfile.ModFilename), r.mod.ModFile.Bytes()); err != nil {
		return nil
	}
	if err := writeFile(filepath.Join(r.mod.Dir, ModSumFilename), r.cache.ModuleSum()); err != nil {
		return nil
	}
	return writeFile(filepath.Join(r.CueModRoot(), ".gitignore"), []byte(`gen/
pkg/
`))
}

func (r *Runtime) build(ctx context.Context, filename string) *build.Instance {
	return Build(filename, OptRoot(r.mod.Dir), OptImportFunc(func(importPath string, importedAt string) (resolvedDir string, err error) {
		return r.Resolve(ctx, importPath, importedAt)
	}))
}

func (r *Runtime) autoImport(ctx context.Context, fromPath string) error {
	cueFiles, err := r.ListCue(fromPath)
	if err != nil {
		return err
	}

	for i := range cueFiles {
		_ = r.build(ctx, cueFiles[i])
	}

	return err
}

func (r *Runtime) download(ctx context.Context, importPath string) error {
	importPathAndVersion := strings.Split(importPath, "@")

	importPath, version := importPathAndVersion[0], ""
	if len(importPathAndVersion) > 1 {
		version = importPathAndVersion[1]
	}

	p, err := r.mod.ResolveImportPath(ctx, r.cache, importPath, version)
	if err != nil {
		return err
	}

	return r.setRequireFromImportPath(ctx, p, true)
}
