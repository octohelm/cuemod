package modfile

import (
	"os"
	"path/filepath"
	"strconv"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
)

func LoadModFile(dir string, m *ModFile) (bool, error) {
	f := filepath.Join(dir, ModFilename)

	if m.comments == nil {
		m.comments = map[string][]*ast.CommentGroup{}
	}

	if m.Replace == nil {
		m.Replace = map[PathMayWithVersion]ReplaceTarget{}
	}

	if m.Require == nil {
		m.Require = map[string]Require{}
	}

	data, err := os.ReadFile(f)

	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if len(data) > 0 {
		f, err := parser.ParseFile(ModFilename, data, parser.ParseComments)
		if err != nil {
			return false, err
		}

		for i := range f.Decls {
			decl := f.Decls[i]
			if field, ok := decl.(*ast.Field); ok {
				directive := stringValue(field.Label)

				if directive != "" {
					switch directive {
					case "module":
						if module := stringValue(field.Value); module != "" {
							m.Module = module
						}
					case "replace":
						if s, ok := field.Value.(*ast.StructLit); ok {
							for i := range s.Elts {
								if subField, ok := s.Elts[i].(*ast.Field); ok {
									from := stringValue(subField.Label)
									to := stringValue(subField.Value)

									if from != "" {
										cg := ast.Comments(subField)

										// from: xxx: xxx
										if s.Lbrace == token.NoPos {
											cg = ast.Comments(field)
										}

										m.comments[directive+"://"+from] = cg

										r, err := ParsePathMayWithVersion(from)
										if err != nil {
											return false, err
										}

										replaceTarget := ReplaceTarget{}

										if to != "" {
											if err := replaceTarget.UnmarshalText([]byte(to)); err != nil {
												return false, err
											}
										}

										for i := range subField.Attrs {
											k, v := subField.Attrs[i].Split()

											switch k {
											case "gen":
												value, _ := strconv.Unquote(v)
												replaceTarget.Gen = value
											}
										}

										if replaceTarget.Path == "" {
											replaceTarget.Path = r.Path
										}

										m.Replace[*r] = replaceTarget
									}
								}
							}
						}
					case "require":
						if s, ok := field.Value.(*ast.StructLit); ok {
							for i := range s.Elts {
								if subField, ok := s.Elts[i].(*ast.Field); ok {
									module := stringValue(subField.Label)
									version := stringValue(subField.Value)

									if module != "" && version != "" {
										cg := ast.Comments(subField)

										// require: xxx: xxx
										if s.Lbrace == token.NoPos {
											cg = ast.Comments(subField)
										}

										m.comments[directive+"://"+module] = cg

										r := Require{}
										r.Version = version

										for _, attr := range subField.Attrs {
											k, v := attr.Split()

											switch k {
											case "vcs":
												value, _ := strconv.Unquote(v)
												r.VcsVersion = value
											case "indirect":
												r.Indirect = true
											}
										}

										m.Require[module] = r
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return false, nil
}

func stringValue(node ast.Node) string {
	switch v := node.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.BasicLit:
		switch v.Kind {
		case token.STRING:
			s, _ := strconv.Unquote(v.Value)
			return s
		case token.INT, token.FLOAT, token.FALSE, token.TRUE:
			return v.Value
		}
	}

	return ""
}
