package cuemod

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	cueerrors "cuelang.org/go/cue/errors"
	cueload "cuelang.org/go/cue/load"

	"cuelang.org/go/cue/build"
	"github.com/octohelm/cuemod/pkg/cuemod/modfile"
	"github.com/pkg/errors"
)

type cuemodKey struct{}

func FromContext(c context.Context) *Context {
	return c.Value(cuemodKey{}).(*Context)
}

func InjectContext(c context.Context, cc *Context) context.Context {
	return context.WithValue(c, cuemodKey{}, cc)
}

func ContextFor(root string) *Context {
	vm := &Context{}

	if !filepath.IsAbs(root) {
		cwd, _ := os.Getwd()
		root = filepath.Join(cwd, root)
	}

	rootMod := &Mod{}
	rootMod.Dir = root

	vm.resolver = newModResolver()

	ctx := context.Background()

	if _, err := rootMod.LoadInfo(ctx); err != nil {
		panic(err)
	}

	vm.Mod = rootMod
	vm.resolver.Collect(ctx, rootMod)

	return vm
}

type Context struct {
	Mod      *Mod
	resolver *modResolver
}

func (r *Context) Cleanup() error {
	if err := os.RemoveAll(filepath.Join(r.CueModRoot(), "gen")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(r.CueModRoot(), "pkg")); err != nil {
		return err
	}
	return nil
}

func (r *Context) CompletePath(p string) string {
	if filepath.IsAbs(p) {
		if filepath.Ext(p) == ".cue" {
			return p
		}
		if !strings.HasPrefix(p, r.Mod.Dir+"/") {
			return p
		}
		p, _ = filepath.Rel(r.Mod.Dir, p)
		if strings.HasPrefix(p, r.Mod.Repo+"/") || p == r.Mod.Repo {
			return p
		}
		return filepath.Join(r.Mod.Repo, p)
	}
	if filepath.Ext(p) == ".cue" {
		return filepath.Join(r.Mod.Dir, p)
	}
	if strings.HasPrefix(p, r.Mod.Repo+"/") || p == r.Mod.Repo {
		return p
	}
	return filepath.Join(r.Mod.Repo, p)
}

func (r *Context) ListCue(fromPath string) ([]string, error) {
	files := make([]string, 0)

	walkSubDir := strings.HasSuffix(fromPath, "/...")

	if walkSubDir {
		fromPath = fromPath[0 : len(fromPath)-4]
	}

	start := filepath.Join(r.Mod.Dir, fromPath)

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

func (r *Context) Get(ctx context.Context, i string) error {
	if i[0] == '.' {
		return r.autoImport(ctx, i)
	}
	return r.download(ctx, i)
}

func (r *Context) CueModRoot() string {
	return filepath.Join(r.Mod.Dir, "cue.mod")
}

func (r *Context) setRequireFromImportPath(ctx context.Context, p *Path, indirect bool) error {
	modVersion := p.ModVersion

	if mv := r.resolver.RepoVersion(p.Repo); mv.Version != "" {
		modVersion = mv
	}

	if err := p.SymlinkOrImport(ctx, r.Mod.Dir); err != nil {
		return err
	}

	// only root mod could be in require
	if p.Root {
		r.Mod.SetRequire(p.ImportPathRoot(), modVersion, indirect)
	}

	return r.syncFiles()
}

func (r *Context) syncFiles() error {
	if err := writeFile(filepath.Join(r.Mod.Dir, modfile.ModFilename), r.Mod.ModFile.Bytes()); err != nil {
		return nil
	}
	if err := writeFile(filepath.Join(r.Mod.Dir, ModSumFilename), r.resolver.ModuleSum()); err != nil {
		return nil
	}
	return writeFile(filepath.Join(r.CueModRoot(), ".gitignore"), []byte(strings.TrimSpace(`
gen/
pkg/
`)))
}

func (r *Context) BuildConfig(ctx context.Context, options ...OptionFunc) *cueload.Config {
	return BuildConfig(append([]OptionFunc{
		OptRoot(r.Mod.Dir),
		OptImportFunc(func(importPath string, importedAt string) (resolvedDir string, err error) {
			return r.Resolve(ctx, importPath, importedAt)
		}),
	}, options...)...)
}

func (r *Context) Build(ctx context.Context, files []string, options ...OptionFunc) *build.Instance {
	return BuildInstances(r.BuildConfig(ctx, options...), files)[0]
}

func (r *Context) Resolve(ctx context.Context, importPath string, importedAt string) (string, error) {
	resolvedImportPath, err := r.resolver.ResolveImportPath(ctx, r.Mod, importPath, "")
	if err != nil {
		return "", errors.Wrapf(err, "resolve import `%s` failed", importPath)
	}

	indirect := isSubDirFor(importedAt, r.CueModRoot()) && !isSubDirFor(importedAt, filepath.Join(r.CueModRoot(), "usr", r.Mod.Module))

	if err := r.setRequireFromImportPath(ctx, resolvedImportPath, indirect); err != nil {
		return "", err
	}

	dir := resolvedImportPath.ResolvedImportPath()
	return dir, nil
}

func (r *Context) autoImport(ctx context.Context, fromPath string) error {
	cueFiles, err := r.ListCue(fromPath)
	if err != nil {
		return err
	}

	for i := range cueFiles {
		if inst := r.Build(ctx, []string{cueFiles[i]}); inst.Err != nil {
			cueerrors.Print(os.Stdout, err, nil)
			return inst.Err
		}
	}

	return err
}

func (r *Context) download(ctx context.Context, importPath string) error {
	importPathAndVersion := strings.Split(importPath, "@")

	importPath, version := importPathAndVersion[0], ""
	if len(importPathAndVersion) > 1 {
		version = importPathAndVersion[1]
	}

	if lang := OptsFromContext(ctx).Import; lang != "" {
		p := modfile.VersionedPathIdentity{Path: importPath}

		if v, ok := r.Mod.Replace[p]; ok {
			v.Import = lang
			r.Mod.Replace[p] = v
		} else {
			r.Mod.Replace[p] = modfile.ReplaceTarget{
				VersionedPathIdentity: p,
				Import:                lang,
			}
		}
	}

	p, err := r.resolver.ResolveImportPath(ctx, r.Mod, importPath, version)
	if err != nil {
		return err
	}

	return r.setRequireFromImportPath(ctx, p, true)
}
