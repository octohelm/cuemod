package golang

import (
	"context"
	"go/ast"
	"go/build"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/runtime/schema"

	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cue/native"
	"github.com/octohelm/cuemod/pkg/extractor/golang/std"
	"github.com/pkg/errors"
)

type pkgExtractor struct {
	*build.Package

	fset   *token.FileSet
	Syntax []*ast.File

	Types     *types.Package
	TypesInfo types.Info

	CueTypes map[types.Type]cueast.Expr

	cueTypes map[ast.Expr]cueast.Expr

	GroupVersion *schema.GroupVersion
}

func (e *pkgExtractor) Extract(ctx context.Context) ([]*cueast.File, error) {
	if err := e.load(ctx, token.NewFileSet()); err != nil {
		return nil, err
	}

	files := make([]*cueast.File, 0)

	for i := range e.GoFiles {
		f := e.extractGoFile(ctx, e.GoFiles[i], e.Syntax[i])
		if f != nil {
			files = append(files, f)
		}
	}

	return files, nil
}

func (e *pkgExtractor) load(ctx context.Context, fset *token.FileSet) error {
	e.fset = fset
	e.Syntax = make([]*ast.File, len(e.GoFiles))

	for i := range e.GoFiles {
		gofile := filepath.Join(e.Dir, e.GoFiles[i])
		data, err := os.ReadFile(gofile)
		if err != nil {
			return err
		}
		f, err := parser.ParseFile(e.fset, gofile, data, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return err
		}
		e.Syntax[i] = f
	}

	conf := types.Config{
		Importer:                 newFakeImporter(),
		IgnoreFuncBodies:         true,
		DisableUnusedImportCheck: true,
		Error: func(err error) {
		},
	}

	e.TypesInfo.Defs = map[*ast.Ident]types.Object{}
	e.TypesInfo.Uses = map[*ast.Ident]types.Object{}
	e.TypesInfo.Types = map[ast.Expr]types.TypeAndValue{}

	pkgTypes, _ := conf.Check(e.ImportPath, e.fset, e.Syntax, &e.TypesInfo)
	//if err != nil {
	//	logr.FromContext(ctx).Debug("type checking error: %s", err)
	//}
	e.Types = pkgTypes

	// GroupVersion
	for i := range e.TypesInfo.Defs {
		if i.Name == "SchemeGroupVersion" && i.Obj != nil {
			if valueSpec, ok := i.Obj.Decl.(*ast.ValueSpec); ok {
				if len(valueSpec.Values) == 1 {
					if compositeLit, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
						if selectorExpr, ok := compositeLit.Type.(*ast.SelectorExpr); ok {
							if selectorExpr.Sel.Name == "GroupVersion" && len(compositeLit.Elts) == 2 {
								for _, elt := range compositeLit.Elts {
									if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
										if key, ok := keyValueExpr.Key.(*ast.Ident); ok {
											if tv, ok := e.TypesInfo.Types[keyValueExpr.Value]; ok {
												v, _ := strconv.Unquote(tv.Value.String())

												gv := &schema.GroupVersion{}

												switch key.Name {
												case "Group":
													gv.Group = v
												case "Version":
													gv.Version = v
												}

												e.GroupVersion = gv
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (e *pkgExtractor) extractGoFile(ctx context.Context, filename string, f *ast.File) *cueast.File {
	file := &cueast.File{}
	file.Filename = strings.Replace(filepath.Base(filename), filepath.Ext(filename), "_go_gen.cue", -1)
	pkgDecl := &cueast.Package{Name: cueast.NewIdent(e.Name)}

	genDecls := make([]cueast.Decl, 0)

	addDecl := func(decl cueast.Decl) {
		if decl == nil {
			return
		}
		genDecls = append(genDecls, decl)
	}

	imports := importPaths{}

	ctx = withImportPaths(ctx, imports)

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			for i := range decl.Specs {
				s := decl.Specs[i]

				newline := i == 0

				switch spec := s.(type) {
				case *ast.TypeSpec:
					if !ast.IsExported(spec.Name.Name) {
						continue
					}

					if spec.Doc == nil {
						spec.Doc = decl.Doc
					}

					addDecl(
						e.def(
							spec.Name.Name,
							e.cueTypeFromAstType(ctx, spec.Type, spec.Name),
							newline,
							e.cueCommentGroup(spec.Doc, true),
						),
					)
				case *ast.ValueSpec:
					if decl.Tok != token.CONST {
						continue
					}

					for i, name := range spec.Names {
						if !ast.IsExported(name.Name) {
							continue
						}

						if obj := e.objectOf(name); obj != nil {
							var astValue ast.Expr

							if len(spec.Values) != 0 {
								astValue = spec.Values[i]
							}

							if c, ok := obj.(*types.Const); ok {
								if v := e.cueValue(ctx, c, astValue); v != nil {

									if len(spec.Names) == 1 && spec.Doc == nil {
										spec.Doc = decl.Doc
									}

									addDecl(
										e.def(
											name.Name,
											v,
											newline, // new line when last
											e.cueCommentGroup(spec.Doc, true),
											e.cueCommentGroup(spec.Comment, false),
										),
									)
								}

							}
						}
					}
				}
			}
		}
	}

	if len(genDecls) == 0 {
		return nil
	}

	file.Decls = []cueast.Decl{pkgDecl}

	if importDecl := imports.toImportDecl(); len(importDecl.Specs) > 0 {
		file.Decls = append(file.Decls, importDecl)
	}

	file.Decls = append(file.Decls, genDecls...)

	return file
}

func (e *pkgExtractor) cueTypeFromAstType(ctx context.Context, astType ast.Expr, defined *ast.Ident) cueast.Expr {
	if e.cueTypes == nil {
		e.cueTypes = map[ast.Expr]cueast.Expr{}
	}

	if tpe, ok := e.cueTypes[astType]; ok {
		return tpe
	}

	tpe := e.makeRootType(ctx, astType, defined)

	e.cueTypes[astType] = tpe
	return tpe
}

func (e *pkgExtractor) makeRootType(ctx context.Context, astType ast.Expr, defined *ast.Ident) cueast.Expr {
	if defined != nil {
		if e.GroupVersion != nil {
			ctx = withGroupVersionKind(ctx, &schema.GroupVersionKind{
				Group:   e.GroupVersion.Group,
				Version: e.GroupVersion.Version,
				Kind:    defined.Name,
			})
		}

		if obj := e.objectOf(defined); obj != nil {
			tpe := e.makeTypeFromNamed(ctx, obj.Type().(*types.Named))

			if tpe != nil {
				return tpe
			}
		}
	}

	return e.makeTypeFromAstType(ctx, astType)
}

func (e *pkgExtractor) makeTypeFromNamed(ctx context.Context, named *types.Named) cueast.Expr {
	if altType := e.altType(ctx, named); altType != nil {
		return altType
	}

	if enums := e.enumsOf(named); len(enums) > 0 {
		return oneOf(enums...)
	}

	return nil
}

func (e *pkgExtractor) makeTypeFromAstType(ctx context.Context, astType ast.Expr) cueast.Expr {
	switch x := (astType).(type) {

	case *ast.SelectorExpr:
		return e.ref(ctx, x)
	case *ast.Ident:
		if tv, ok := e.TypesInfo.Types[astType]; ok {
			switch tv.Type.(type) {
			case *types.Basic:
				return e.ident(x.Name, false)
			case *types.Named:
				return e.ident(x.Name, true)
			}
		}
	case *ast.ArrayType:
		if elm, ok := x.Elt.(*ast.Ident); ok {
			if elm.Name == "byte" {
				return e.ident("bytes", false)
			}
		}

		elmType := e.cueTypeFromAstType(ctx, x.Elt, nil)

		if elmType == nil {
			return nil
		}

		// array
		if x.Len != nil {
			n, _ := strconv.Atoi(x.Len.(*ast.BasicLit).Value)

			return cueast.NewBinExpr(
				cuetoken.MUL,
				newInt(n),
				cueast.NewList(elmType),
			)
		}
		// slice
		return cueast.NewList(&cueast.Ellipsis{Type: elmType})
	case *ast.MapType:
		propType := e.cueTypeFromAstType(ctx, x.Key, nil)
		elemType := e.cueTypeFromAstType(ctx, x.Value, nil)

		if propType == nil || elemType == nil {
			return nil
		}

		f := &cueast.Field{
			Label: cueast.NewList(propType),
			Value: elemType,
		}

		cueast.SetRelPos(f, cuetoken.Blank)

		s := cueast.NewStruct(f)
		s.Lbrace = cuetoken.Blank.Pos()
		s.Rbrace = cuetoken.Blank.Pos()

		return s
	case *ast.StarExpr:
		// ptr
		underlying := e.cueTypeFromAstType(ctx, x.X, nil)
		if underlying == nil {
			return nil
		}
		return oneOf(cueast.NewNull(), underlying)
	case *ast.FuncType:
		return nil
	case *ast.InterfaceType:
		return any()
	case *ast.StructType:
		st := &cueast.StructLit{
			Lbrace: cuetoken.Blank.Pos(),
			Rbrace: cuetoken.Newline.Pos(),
		}
		e.addFieldsFromAstStructType(ctx, x, st)

		if len(st.Elts) == 0 {
			return nil
		}

		return st
	default:
		logr.FromContext(ctx).Warn(errors.Errorf("unsupported ast type %#v", x))
	}

	return nil
}

func (e *pkgExtractor) addFieldsFromAstStructType(ctx context.Context, x *ast.StructType, st *cueast.StructLit) {
	add := func(x cueast.Decl) {
		st.Elts = append(st.Elts, x)
	}

	indirect := func(tpe ast.Expr) ast.Expr {
		for {
			p, ok := tpe.(*ast.StarExpr)
			if !ok {
				break
			}
			tpe = p.X
		}
		return tpe
	}

	for i := range x.Fields.List {
		astField := x.Fields.List[i]
		fieldType := astField.Type

		tag := ""
		if astField.Tag != nil {
			tag, _ = strconv.Unquote(astField.Tag.Value)
		}

		names := make([]string, 0)

		for _, name := range astField.Names {
			if name.IsExported() {
				names = append(names, name.Name)
			}
		}

		anonymous := len(astField.Names) == 0

		if anonymous {
			switch x := fieldType.(type) {
			case *ast.Ident:
				if x.IsExported() {
					names = append(names, x.Name)
				}
			case *ast.SelectorExpr:
				if x.Sel.IsExported() {
					names = append(names, x.Sel.Name)
				}
			}
		}

		for _, goFieldName := range names {
			fieldName, omitempty, hasNamedTag := getName(goFieldName, tag)

			if fieldName == "-" {
				continue
			}

			if anonymous && (!hasNamedTag || isInline(tag)) {
				typ := indirect(fieldType)

				switch x := typ.(type) {
				case *ast.StructType:
					e.addFieldsFromAstStructType(ctx, x, st)
				case *ast.Ident, *ast.SelectorExpr:
					sel := e.ref(ctx, x)
					if sel != nil {
						embed := &cueast.EmbedDecl{Expr: sel}
						if i > 0 {
							cueast.SetRelPos(embed, cuetoken.NewSection)
						}
						add(embed)
					}
				default:
					logr.FromContext(ctx).Warn(errors.Errorf("unimplemented embedding %s for type %T", goFieldName, x))
				}

				continue
			}

			kind := cuetoken.COLON
			if omitempty {
				kind = cuetoken.OPTION
				fieldType = indirect(fieldType)
			}

			typ := e.cueTypeFromAstType(ctx, fieldType, nil)

			if typ == nil {
				logr.FromContext(ctx).Warn(errors.Errorf("drop field %s, unsupport type %T", goFieldName, fieldType))
				continue
			}

			var label cueast.Label

			if kind == cuetoken.ISA {
				label = e.ident(fieldName, true)
			} else {
				label = cueast.NewString(fieldName)
			}

			field := &cueast.Field{Label: label, Value: typ}

			addComments(field, e.cueCommentGroup(astField.Doc, true), e.cueCommentGroup(astField.Comment, false))

			if kind == cuetoken.OPTION {
				field.Token = cuetoken.COLON
				field.Optional = cuetoken.Blank.Pos()
			}

			add(field)
		}
	}
}

func (e *pkgExtractor) objectOf(ident *ast.Ident) types.Object {
	for i, obj := range e.TypesInfo.Defs {
		if i == ident {
			return obj
		}
	}
	return nil
}

func (e *pkgExtractor) cueCommentGroup(c *ast.CommentGroup, doc bool) *cueast.CommentGroup {
	if c == nil {
		return nil
	}

	cg := &cueast.CommentGroup{
		Doc:  doc,
		Line: !doc,
	}

	if cg.Line {
		cg.Position = 3 // after value
	}

	for _, comment := range c.List {
		cg.List = append(cg.List, &cueast.Comment{
			Text: comment.Text,
		})
	}

	return cg
}

func (e *pkgExtractor) cueValue(ctx context.Context, c *types.Const, value ast.Expr) cueast.Expr {
	v := c.Val()

	switch v.Kind() {
	case constant.Unknown:
		return e.ref(ctx, value)
	case constant.String:
		return cueast.NewLit(cuetoken.STRING, v.String())
	case constant.Int:
		return cueast.NewLit(cuetoken.INT, v.String())
	case constant.Float:
		return cueast.NewLit(cuetoken.FLOAT, v.String())
	case constant.Bool:
		b, _ := strconv.ParseBool(v.String())
		return cueast.NewBool(b)
	}

	logr.FromContext(ctx).Warn(errors.Errorf("invalid const value: %d %#v", v.Kind(), c))

	return nil
}

func (e *pkgExtractor) enumsOf(tpe *types.Named) []cueast.Expr {
	names := make([]string, 0)

	for ident, def := range e.TypesInfo.Defs {
		c, ok := def.(*types.Const)
		if !ok {
			continue
		}

		if c.Type() != tpe {
			continue
		}

		// skip private
		if !ident.IsExported() {
			continue
		}

		name := ident.Name

		// skip hidden
		if name[0] == '_' {
			continue
		}

		names = append(names, name)
	}

	sort.Strings(names)

	ids := make([]cueast.Expr, len(names))
	for i, name := range names {
		ids[i] = e.ident(name, true)
	}
	return ids
}

func (e *pkgExtractor) altType(ctx context.Context, typ *types.Named) cueast.Expr {
	methods := map[string]*types.Func{}

	for i := 0; i < typ.NumMethods(); i++ {
		fn := typ.Method(i)

		if !fn.Exported() {
			continue
		}

		methods[fn.Name()] = fn
	}

	for name, stringInterface := range stringInterfaces {
		if fn, ok := methods[name]; ok {
			if stringInterface.Equal(fn) {
				return cueast.NewIdent("string")
			}
		}
	}

	for name, topInterface := range topInterfaces {
		if fn, ok := methods[name]; ok {
			if topInterface.Equal(fn) {
				return any()
			}
		}
	}

	return nil
}

func (e *pkgExtractor) sel(ctx context.Context, name string, pkgName string, importPath string) cueast.Expr {
	if std.IsStd(importPath) {
		// only builtin pkg can be select
		if native.IsBuiltinPackage(importPath) {
			importPathsFromContext(ctx).add(importPath, pkgName)
			// std & builtin pkg don't use #prefix
			return &cueast.SelectorExpr{X: e.ident(pkgName, false), Sel: e.ident(name, false)}
		}
		return nil
	}

	sel := &cueast.SelectorExpr{X: e.ident(pkgName, false), Sel: e.ident(name, true)}

	importPathsFromContext(ctx).add(importPath, pkgName)

	if gvk := groupVersionKindFromContext(ctx); gvk != nil {
		if importPath == "k8s.io/apimachinery/pkg/apis/meta/v1" && name == "TypeMeta" {
			return allOf(sel, cueast.NewStruct(
				&cueast.Field{
					Label: cueast.NewString("apiVersion"),
					Value: cueast.NewString(gvk.GroupVersion().String()),
				},
				&cueast.Field{
					Label: cueast.NewString("kind"),
					Value: cueast.NewString(gvk.Kind),
				},
			))
		}
	}

	return sel
}

func (e *pkgExtractor) ref(ctx context.Context, expr ast.Expr) cueast.Expr {
	switch x := expr.(type) {
	case *ast.Ident:
		return e.ident(x.Name, true)
	case *ast.SelectorExpr:
		from := x.X.(*ast.Ident)
		if o, ok := e.TypesInfo.Uses[from]; ok {
			if pkgName, ok := o.(*types.PkgName); ok {
				return e.sel(ctx, x.Sel.Name, pkgName.Name(), pkgName.Imported().Path())
			}
		}
	}
	return nil
}

func isInline(tag string) bool {
	return hasFlag(tag, "json", "inline", 1) || hasFlag(tag, "yaml", "inline", 1)
}

func hasFlag(tag, key, flag string, offset int) bool {
	if t := reflect.StructTag(tag).Get(key); t != "" {
		split := strings.Split(t, ",")
		if offset >= len(split) {
			return false
		}
		for _, str := range split[offset:] {
			if str == flag {
				return true
			}
		}
	}
	return false
}

func getName(name string, tag string) (n string, omitempty bool, hasTag bool) {
	tags := reflect.StructTag(tag)
	for _, s := range []string{"json", "yaml"} {
		if tag, ok := tags.Lookup(s); ok {
			omitempty := false

			if p := strings.Index(tag, ","); p >= 0 {
				omitempty = strings.Contains(tag, "omitempty")
				tag = tag[:p]
			}
			if tag != "" {
				return tag, omitempty, true
			}
		}
	}
	return name, false, false
}

func (e *pkgExtractor) ident(name string, isDef bool) *cueast.Ident {
	if isDef {
		r := []rune(name)[0]
		name = "#" + name
		if !unicode.Is(unicode.Lu, r) {
			name = "_" + name
		}
	}
	return cueast.NewIdent(name)
}

func (e *pkgExtractor) def(name string, valueOrType cueast.Expr, newline bool, comments ...*cueast.CommentGroup) cueast.Decl {
	if valueOrType == nil {
		return nil
	}

	f := &cueast.Field{
		Label: e.ident(name, true),
		Value: valueOrType,
	}

	addComments(f, comments...)

	if newline {
		cueast.SetRelPos(f, cuetoken.NewSection)
	}

	return f
}
