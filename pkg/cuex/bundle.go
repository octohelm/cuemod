package cuex

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

// BundleToRaw bundle instance to single cue file
func BundleToRaw(inst *build.Instance) ([]byte, error) {
	sf := &bundler{
		stds:    map[string]*cueast.ImportSpec{},
		imports: map[string]*cueast.Field{},
	}

	f, err := sf.Export(inst)
	if err != nil {
		return nil, err
	}

	return format.Node(f, format.Simplify())
}

type bundler struct {
	stds map[string]*cueast.ImportSpec

	imports      map[string]*cueast.Field
	importOrders []string
}

func (sf *bundler) importPkg(importPath string, f *cueast.Field) {
	if _, ok := sf.imports[importPath]; !ok {
		sf.imports[importPath] = f

		cueast.AddComment(f, &cueast.CommentGroup{
			Doc: true,
			List: []*cueast.Comment{{
				Text: "// " + importPath,
			}},
		})

		sf.importOrders = append([]string{importPath}, sf.importOrders...)
	}
}

func (sf *bundler) importDecl() *cueast.ImportDecl {
	stds := make([]string, 0)

	for i := range sf.stds {
		stds = append(stds, i)
	}

	if len(stds) == 0 {
		return nil
	}

	sort.Strings(stds)

	d := &cueast.ImportDecl{}

	d.Specs = make([]*cueast.ImportSpec, len(stds))

	for i, importPath := range stds {
		d.Specs[i] = sf.stds[importPath]
	}

	return d
}

func (sf *bundler) importAliases() []cueast.Decl {
	decls := make([]cueast.Decl, 0)

	for _, importPath := range sf.importOrders {
		decls = append(decls, sf.imports[importPath])
	}

	return decls
}

func (sf *bundler) Export(inst *build.Instance) (*cueast.File, error) {
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

func (sf *bundler) Walk(inst *build.Instance) (*cueast.File, error) {
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
						sf.stds[importPath] = cueast.NewImport(id, importPath)
					} else {
						for _, dep := range inst.Imports {
							if dep.ImportPath == importPath {
								f, err := sf.Walk(dep)
								if err != nil {
									return nil, err
								}

								id := cueast.NewIdent(toSafeID(importPath))

								sf.importPkg(importPath, &cueast.Field{
									Label: id,
									Value: &cueast.StructLit{
										Elts: f.Decls,
									},
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
				if id, ok := node.(*cueast.Ident); ok && id.Node != nil {
					if _, ok := id.Node.(*cueast.ImportSpec); ok {
						for n, uniqPkgName := range importAliases {
							if n == id.Name {
								id.Name = uniqPkgName
							}
						}
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
