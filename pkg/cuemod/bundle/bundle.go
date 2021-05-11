package bundle

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/format"
	"github.com/octohelm/cuemod/pkg/cuemod/builtin"
)

// ToRaw bundle instance to single cue file
func ToRaw(inst *build.Instance) ([]byte, error) {
	sf := &singleFile{Stds: map[string]*cueast.ImportSpec{}, Pkgs: map[string]*cueast.Field{}}
	f, err := sf.Export(inst)
	if err != nil {
		return nil, err
	}
	return format.Node(f, format.Simplify())
}

type singleFile struct {
	Stds map[string]*cueast.ImportSpec
	Pkgs map[string]*cueast.Field
}

func (sf *singleFile) importDecl() *cueast.ImportDecl {
	stds := make([]string, 0)

	for i := range sf.Stds {
		stds = append(stds, i)
	}

	if len(stds) == 0 {
		return nil
	}

	sort.Strings(stds)

	d := &cueast.ImportDecl{}

	d.Specs = make([]*cueast.ImportSpec, len(stds))

	for i, importPath := range stds {
		d.Specs[i] = sf.Stds[importPath]
	}

	return d
}

func (sf *singleFile) importAliases() []cueast.Decl {
	pkgAliases := make([]string, 0)

	for i := range sf.Pkgs {
		pkgAliases = append(pkgAliases, i)
	}

	if len(pkgAliases) == 0 {
		return nil
	}

	sort.Strings(pkgAliases)

	decls := make([]cueast.Decl, 0)

	for _, importPath := range pkgAliases {
		decls = append(decls, sf.Pkgs[importPath])
	}

	return decls
}

func (sf *singleFile) Export(inst *build.Instance) (*cueast.File, error) {
	f, err := sf.Walk(inst)
	if err != nil {
		return nil, err
	}

	decls := f.Decls

	f.Decls = make([]cueast.Decl, 0)

	if importDecl := sf.importDecl(); importDecl != nil {
		f.Decls = append(f.Decls, importDecl)
	}

	f.Decls = append(f.Decls, &cueast.StructLit{
		Elts: decls,
	})

	f.Decls = append(f.Decls, sf.importAliases()...)

	return f, nil
}

func (sf *singleFile) Walk(inst *build.Instance) (*cueast.File, error) {
	f := &cueast.File{
		Filename: fmt.Sprintf("%s/%s.cue", inst.ImportPath, inst.PkgName),
	}

	stmts := make([]cueast.Decl, 0)
	importAliases := map[string]string{}

	for _, file := range inst.Files {
		for _, d := range file.Decls {
			switch decl := d.(type) {
			case *cueast.Package:
				continue
			case *cueast.ImportDecl:
				for i := range decl.Specs {
					spec := decl.Specs[i]

					importPath, _ := strconv.Unquote(spec.Path.Value)

					if builtin.IsBuiltIn(importPath) {
						id := spec.Name
						if spec.Name != nil {
							id = cueast.NewIdent(spec.Name.Name)
						}
						sf.Stds[importPath] = cueast.NewImport(id, importPath)
					} else {
						for _, dep := range inst.Imports {
							if dep.ImportPath == importPath {
								f, err := sf.Walk(dep)
								if err != nil {
									return nil, err
								}

								id := cueast.NewIdent(toSafeID(importPath))

								sf.Pkgs[importPath] = &cueast.Field{
									Label: id,
									Value: &cueast.StructLit{
										Elts: f.Decls,
									},
								}

								cueast.AddComment(sf.Pkgs[importPath], &cueast.CommentGroup{
									Doc: true,
									List: []*cueast.Comment{{
										Text: "// " + importPath,
									}},
								})

								n := cueast.NewIdent(dep.PkgName)

								if p := strings.Split(importPath, ":"); len(p) == 2 {
									n = cueast.NewIdent(p[1])
								}

								if spec.Name != nil {
									n = cueast.NewIdent(spec.Name.Name)
								}

								importAliases[n.Name] = id.Name
							}
						}
					}
				}

			default:
				stmts = append(stmts, decl)
			}
		}
	}

	for _, stmt := range stmts {
		cueast.Walk(
			stmt,
			func(node cueast.Node) bool {
				if selectExpr, ok := node.(*cueast.SelectorExpr); ok {
					// xx.XX
					if id, ok := selectExpr.X.(*cueast.Ident); ok {
						for n, uniqPkgName := range importAliases {
							if n == id.Name {
								id.Name = uniqPkgName
							}
						}

						return false
					}
				}

				return true
			},
			nil,
		)

		f.Decls = append(f.Decls, stmt)
	}

	return f, nil
}

var re = regexp.MustCompile(`[^0-9A-Za-z_]`)

func toSafeID(importPath string) string {
	return "_" + re.ReplaceAllString(importPath, "_")
}
