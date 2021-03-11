package jsonnet

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
	gojsonnet "github.com/google/go-jsonnet"
	gojsonnetast "github.com/google/go-jsonnet/ast"
	"github.com/google/go-jsonnet/toolutils"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/octohelm/cuemod/pkg/extractor/core"
)

type File struct {
	Name    string
	Data    []byte
	Imports map[string]string
}

func (f *File) String() string {
	return fmt.Sprintf("%v", f.Imports)
}

func Load(ctx context.Context, root string) (*loader, error) {
	l := &loader{root: root}

	f := filepath.Join(l.root, "jsonnetfile.json")

	if data, err := os.ReadFile(f); err == nil {
		jf := &spec.JsonnetFile{}

		if err := jf.UnmarshalJSON(data); err != nil {
			return nil, err
		}

		l.jsonnetfile = jf
	}

	jsonnetfiles, err := filepath.Glob(filepath.Join(l.root, "*.*sonnet"))
	if err != nil {
		return nil, err
	}

	for i := range jsonnetfiles {
		rel, _ := filepath.Rel(l.root, jsonnetfiles[i])
		if err := l.Import(ctx, rel, filepath.Join(l.root, "main.jsonnet"), true); err != nil {
			return nil, err
		}
	}

	return l, nil
}

type loader struct {
	root        string
	files       map[string]*File
	jsonnetfile *spec.JsonnetFile
}

// jsonnet bundler could rename the import path, should use the origin repo path
func (l *loader) refillImportPath(importPath string) string {
	if l.jsonnetfile != nil {
		for _, d := range l.jsonnetfile.Dependencies {
			if d.Source.GitSource != nil {
				repo := filepath.Join(d.Source.GitSource.Host, d.Source.GitSource.User, d.Source.GitSource.Repo, d.Source.GitSource.Subdir)

				if strings.HasPrefix(importPath+"/", d.LegacyName()+"/") {
					rel, _ := filepath.Rel(d.LegacyName()+"/", importPath+"/")
					return filepath.Join(repo+"/", rel)
				}
			}
		}
	}

	return importPath
}

const filesVar = "_files"

func (l *loader) Extract() (files []*cueast.File, err error) {

	for fieldName := range l.files {
		cuefile := &cueast.File{}
		cuefile.Filename = strings.TrimLeft(strings.Replace(fieldName, "/", "__", -1)+"_gen.cue", "_")

		cuefile.Decls = []cueast.Decl{
			&cueast.Package{Name: cueast.NewIdent(core.SafeIdentifierFromImportPath(l.root))},
		}

		f := l.files[fieldName]

		deps := map[string]string{}
		fields := make([]interface{}, 0)

		imports := make([]interface{}, 0)
		importPaths := make([]string, 0)

		for k := range f.Imports {
			importPaths = append(importPaths, k)
		}

		sort.Strings(importPaths)

		for _, importPath := range importPaths {
			resolved := f.Imports[importPath]

			f := &cueast.Field{Label: cueast.NewString(importPath)}

			// lock
			if _, ok := l.files[resolved]; ok {
				f.Value = &cueast.IndexExpr{X: cueast.NewIdent(filesVar), Index: cueast.NewString(resolved)}
			} else {
				dirname := l.refillImportPath(filepath.Dir(importPath))

				filename := filepath.Base(importPath)

				pkgName := core.SafeIdentifierFromImportPath(dirname)

				deps[fmt.Sprintf("%s:%s", dirname, pkgName)] = pkgName

				f.Value = &cueast.IndexExpr{X: cueast.NewIdent(pkgName), Index: cueast.NewString(filename)}
			}

			imports = append(imports, f)
		}

		field := &cueast.Field{}
		field.Label = cueast.NewString(f.Name)
		field.Value = cueast.NewStruct(
			&cueast.Field{Label: cueast.NewString("imports"), Value: cueast.NewStruct(imports...)},
			&cueast.Field{Label: cueast.NewString("data"), Value: newBytes(f.Data)},
		)

		cueast.SetRelPos(field, cuetoken.NewSection)

		fields = append(fields, field)

		if len(deps) > 0 {
			importDecl := &cueast.ImportDecl{}

			importPaths := make([]string, 0)

			for key := range deps {
				importPaths = append(importPaths, key)
			}

			for _, importPath := range importPaths {
				importDecl.Specs = append(importDecl.Specs, &cueast.ImportSpec{
					Path: cueast.NewString(importPath),
				})
			}

			cuefile.Decls = append(cuefile.Decls, importDecl)
		}

		cuefile.Decls = append(cuefile.Decls,
			&cueast.EmbedDecl{Expr: cueast.NewIdent(filesVar)},
			&cueast.Field{Label: cueast.NewIdent(filesVar), Value: cueast.NewStruct(fields...)},
		)

		files = append(files, cuefile)
	}

	return
}

func newBytes(data []byte) *cueast.BasicLit {
	return &cueast.BasicLit{
		Kind:     cuetoken.STRING,
		ValuePos: cuetoken.NoPos,
		Value:    "'''\n" + strings.Replace(string(data), "\\", "\\\\", -1) + "'''",
	}
}

func (l *loader) fileID(filename string) string {
	if len(filename) > 0 && (filename[0] == '.' || strings.HasPrefix(filepath.Join(".", filename), filepath.Join(".", l.root))) {
		rel, _ := filepath.Rel(l.root, filename)
		return rel
	}
	return filename
}

func (l *loader) Import(ctx context.Context, importPath string, importedAt string, canParse bool) error {
	if l.files == nil {
		l.files = map[string]*File{}
	}

	fileID := l.fileID(importedAt)

	filename := filepath.Join(filepath.Dir(importedAt), importPath)

	resolvedPath := l.fileID(filename)

	defer func() {
		if _, ok := l.files[fileID]; ok {
			if l.files[fileID].Imports == nil {
				l.files[fileID].Imports = map[string]string{}
			}

			l.files[fileID].Imports[importPath] = resolvedPath
		}
	}()

	var data []byte

	if stub, ok := stubs[resolvedPath]; ok {
		canParse = true
		data = []byte(fmt.Sprintf(`(import '%s')`, stub))
	} else {
		d, err := os.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		data = d
	}

	if canParse {
		data = dropDocUtil(data)
	}

	file := &File{}
	file.Name = l.fileID(filename)

	if _, ok := l.files[file.Name]; !ok {
		file.Data = data

		l.files[file.Name] = file

		if canParse {
			node, err := gojsonnet.SnippetToAST(filename, string(data))
			if err != nil {
				return err
			}
			if err := l.walk(ctx, node, filename); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *loader) walk(ctx context.Context, node gojsonnetast.Node, importedAt string) error {
	switch n := (node).(type) {
	case *gojsonnetast.Import:
		return l.Import(ctx, n.File.Value, importedAt, true)
	case *gojsonnetast.ImportStr:
		return l.Import(ctx, n.File.Value, importedAt, false)
	default:
		for _, n := range toolutils.Children(node) {
			if err := l.walk(ctx, n, importedAt); err != nil {
				return err
			}
		}
	}
	return nil
}

var stubs = map[string]string{
	"k.libsonnet": "github.com/jsonnet-libs/k8s-alpha/1.19/main.libsonnet",
}
