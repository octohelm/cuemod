package cuemod

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cuelang.org/go/cue/load"
	"github.com/octohelm/cuemod/pkg/cuex"

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
	vm := &Context{
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

type Context struct {
	mod   *Mod
	cache *cache
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

func (r *Context) completePath(p string) string {
	if filepath.IsAbs(p) {
		if filepath.Ext(p) == ".cue" {
			return p
		}
		if !strings.HasPrefix(p, r.mod.Dir+"/") {
			return p
		}
		p, _ = filepath.Rel(r.mod.Dir, p)
	}
	if strings.HasPrefix(p, r.mod.Repo+"/") || p == r.mod.Repo {
		return p
	}
	return filepath.Join(r.mod.Repo, p)
}

func (r *Context) ListCue(fromPath string) ([]string, error) {
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

func (r *Context) Get(ctx context.Context, i string) error {
	if i[0] == '.' {
		return r.autoImport(ctx, i)
	}
	return r.download(ctx, i)
}

func (r *Context) CueModRoot() string {
	return filepath.Join(r.mod.Dir, "cue.mod")
}

func (r *Context) setRequireFromImportPath(ctx context.Context, p *Path, indirect bool) error {
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

func (r *Context) syncFiles() error {
	if err := writeFile(filepath.Join(r.mod.Dir, modfile.ModFilename), r.mod.ModFile.Bytes()); err != nil {
		return nil
	}
	if err := writeFile(filepath.Join(r.mod.Dir, ModSumFilename), r.cache.ModuleSum()); err != nil {
		return nil
	}
	return writeFile(filepath.Join(r.CueModRoot(), ".gitignore"), []byte(strings.TrimSpace(`
gen/
pkg/
`)))
}

func (r *Context) EvalWithPatches(ctx context.Context, inputs []string, options ...cuex.EvalOptionFunc) ([]byte, error) {
	imports := make([]string, 0)
	b := strings.Builder{}

	c := 0
	for i := range inputs {
		input := inputs[i]
		if input == "" {
			continue
		}

		if c > 0 {
			b.WriteString(" & ")
		}

		if input[0] == '{' {
			b.WriteString(input)
		} else {
			b.WriteString("i" + strconv.Itoa(len(imports)))
			imports = append(imports, r.completePath(input))
		}

		c++
	}

	t := bytes.NewBuffer(nil)

	for i := range imports {
		path := r.completePath(imports[i])

		t.WriteString(
			fmt.Sprintf(`import i%d "%s"
`, i, FixFileImport(path)),
		)
	}

	t.WriteString(b.String())

	// fake main file
	mainCue := filepath.Join(r.mod.Dir, "./main.cue")
	inst := r.Build(ctx, mainCue, OptOverlay(map[string]load.Source{
		mainCue: load.FromBytes(t.Bytes()),
	}))

	for _, p := range imports {
		i := r.Build(ctx, p)
		i.ImportPath = FixFileImport(p)
		inst.Imports = append(inst.Imports, i)
	}

	results, err := cuex.Eval(inst, options...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Context) Eval(ctx context.Context, filename string, options ...cuex.EvalOptionFunc) ([]byte, error) {
	filename = r.completePath(filename)
	inst := r.Build(ctx, filename)
	results, err := cuex.Eval(inst, options...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Context) Build(ctx context.Context, filename string, options ...OptionFunc) *build.Instance {
	return Build(
		filename,
		append([]OptionFunc{
			OptRoot(r.mod.Dir),
			OptImportFunc(func(importPath string, importedAt string) (resolvedDir string, err error) {
				return r.Resolve(ctx, importPath, importedAt)
			}),
		}, options...)...,
	)
}

func (r *Context) Resolve(ctx context.Context, importPath string, importedAt string) (string, error) {
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

func (r *Context) autoImport(ctx context.Context, fromPath string) error {
	cueFiles, err := r.ListCue(fromPath)
	if err != nil {
		return err
	}

	for i := range cueFiles {
		_ = r.Build(ctx, cueFiles[i])
	}

	return err
}

func (r *Context) download(ctx context.Context, importPath string) error {
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
