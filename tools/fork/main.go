package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	importPaths := os.Args[1:]

	cwd, _ := os.Getwd()
	if !(strings.HasSuffix(cwd, "internal/third_party") || strings.HasSuffix(cwd, "internal")) {
		panic(cwd)
	}

	if err := cleanup(cwd); err != nil {
		panic(err)
	}

	prefixes := map[string]bool{}

	for _, importPath := range importPaths {
		prefixes[strings.Split(importPath, "/")[0]] = true
	}

	task, err := TaskFor("./", prefixes)
	if err != nil {
		panic(err)
	}

	for _, importPath := range importPaths {
		if err := task.Scan(importPath); err != nil {
			panic(err)
		}
	}

	if err := task.Sync(); err != nil {
		panic(err)
	}
}

func cleanup(cwd string) error {
	list, err := filepath.Glob(filepath.Join(cwd, "./*"))
	if err != nil {
		return err
	}

	for i := range list {
		p := list[i]

		f, err := os.Lstat(p)
		if err != nil {
			return err
		}

		if f.IsDir() && !strings.HasSuffix(f.Name(), "__gen__") {
			log.Println("REMOVE", p)
			if err := os.RemoveAll(p); err != nil {
				return err
			}
		}
	}

	return nil
}

func TaskFor(dir string, prefixes map[string]bool) (*Task, error) {
	if !filepath.IsAbs(dir) {
		cwd, _ := os.Getwd()
		dir = path.Join(cwd, dir)
	}

	d := dir

	for d != "/" {
		gmodfile := filepath.Join(d, "go.mod")

		if data, err := os.ReadFile(gmodfile); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		} else {
			f, _ := modfile.Parse(gmodfile, data, nil)

			rel, _ := filepath.Rel(d, dir)
			return &Task{
				PkgPath:  filepath.Join(f.Module.Mod.Path, rel),
				Dir:      filepath.Join(d, rel),
				Prefixes: prefixes,
			}, nil
		}

		d = filepath.Join(d, "../")
	}

	return nil, fmt.Errorf("missing go.mod")
}

type Task struct {
	Dir      string
	PkgPath  string
	Prefixes map[string]bool

	// map[importPath][filename]*parsedFile
	packages map[string]map[string]*parsedFile

	pkgInternals map[string]bool
	pkgUsed      map[string][]string
}

func (t *Task) Sync() error {
	needToForks := map[string]bool{}

	var findUsed func(importPath string) []string
	findUsed = func(importPath string) (used []string) {
		for _, importPath := range t.pkgUsed[importPath] {
			used = append(append(used, importPath), findUsed(importPath)...)
		}
		return
	}

	for pkgImportPath := range t.pkgInternals {
		needToForks[pkgImportPath] = true

		for _, p := range findUsed(pkgImportPath) {
			needToForks[p] = true
		}
	}

	for pkgImportPath := range needToForks {
		files := t.packages[pkgImportPath]

		for filename := range files {
			f := files[filename]

			astFile := f.file

			for _, i := range astFile.Imports {
				importPath, _ := strconv.Unquote(i.Path.Value)

				if needToForks[importPath] {
					_ = astutil.RewriteImport(
						f.fset, astFile,
						importPath,
						filepath.Join(t.PkgPath, t.replaceInternal(importPath)),
					)
				}
			}

			output := filepath.Join(t.Dir, t.replaceInternal(pkgImportPath), filepath.Base(filename))

			buf := bytes.NewBuffer(nil)
			if err := format.Node(buf, f.fset, astFile); err != nil {
				return err
			}
			if err := writeFile(output, buf.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Task) Scan(importPath string) error {
	if _, ok := t.packages[importPath]; ok {
		return nil
	}

	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return err
	}

	if err := t.scanPkg(pkg); err != nil {
		return err
	}

	log.Printf("scaned %s", importPath)

	return nil
}

func (t *Task) isInternalPkg(importPath string) bool {
	if strings.Contains(importPath, "internal/") {
		return true
	}
	return strings.HasSuffix(importPath, "internal") || strings.HasPrefix(importPath, "internal")
}

func (t *Task) scanPkg(pkg *build.Package) error {
	files, err := filepath.Glob(pkg.Dir + "/*.go")
	if err != nil {
		return err
	}

	if t.isInternalPkg(pkg.ImportPath) {
		if t.pkgInternals == nil {
			t.pkgInternals = map[string]bool{}
		}
		t.pkgInternals[pkg.ImportPath] = true
	}

	for _, f := range files {
		// skip test file
		if strings.HasSuffix(f, "_test.go") {
			continue
		}

		if err := t.scanGoFile(f, pkg); err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) scanGoFile(filename string, pkg *build.Package) error {
	f, err := newParsedFile(filename)
	if err != nil {
		return err
	}

	file := f.file

	pkgImportPath := pkg.ImportPath

	if pkg.Name == "" {
		if file.Name.Name != "main" {
			pkg.Name = file.Name.Name
		}
	}

	if file.Name.Name != pkg.Name {
		return nil
	}

	for _, i := range file.Imports {
		importPath, _ := strconv.Unquote(i.Path.Value)

		if t.pkgUsed == nil {
			t.pkgUsed = map[string][]string{}
		}

		t.pkgUsed[importPath] = append(t.pkgUsed[importPath], pkgImportPath)

		if t.Prefixes[strings.Split(importPath, "/")[0]] {
			if err := t.Scan(importPath); err != nil {
				return err
			}
		}
	}

	if t.packages == nil {
		t.packages = map[string]map[string]*parsedFile{}
	}

	if t.packages[pkgImportPath] == nil {
		t.packages[pkgImportPath] = map[string]*parsedFile{}
	}

	t.packages[pkgImportPath][filename] = f

	return nil
}

func newParsedFile(filename string) (*parsedFile, error) {
	f := &parsedFile{}

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	f.fset = token.NewFileSet()

	file, err := parser.ParseFile(f.fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	f.file = file

	return f, nil
}

type parsedFile struct {
	fset *token.FileSet
	file *ast.File
}

func (t *Task) replaceInternal(p string) string {
	if strings.HasSuffix(p, "internal") {
		return filepath.Join(filepath.Dir(p), "./internals")
	}

	return strings.Replace(
		p,
		"internal/",
		"internals/",
		-1)
}

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
