// Copyright 2020 CUE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: make this package public in cuelang.org/go/encoding
// once stabalized.

package encoding

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/token"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue/build"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue/format"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/cue/parser"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/encoding/json"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/encoding/jsonschema"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/encoding/openapi"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/encoding/protobuf"
	internal "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals/filetypes"
	"github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/internals/third_party/yaml"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Decoder struct {
	cfg            *Config
	closer         io.Closer
	next           func() (ast.Expr, error)
	interpretFunc  interpretFunc
	interpretation build.Interpretation
	expr           ast.Expr
	file           *ast.File
	filename       string // may change on iteration for some formats
	id             string
	index          int
	err            error
}

type interpretFunc func(*cue.Instance) (file *ast.File, id string, err error)

// ID returns a canonical identifier for the decoded object or "" if no such
// identifier could be found.
func (i *Decoder) ID() string {
	return i.id
}
func (i *Decoder) Filename() string { return i.filename }

// Interpretation returns the current interpretation detected by Detect.
func (i *Decoder) Interpretation() build.Interpretation {
	return i.interpretation
}
func (i *Decoder) Index() int { return i.index }
func (i *Decoder) Done() bool { return i.err != nil }

func (i *Decoder) Next() {
	if i.err != nil {
		return
	}
	// Decoder level
	i.file = nil
	i.expr, i.err = i.next()
	i.index++
	if i.err != nil {
		return
	}
	i.doInterpret()
}

func (i *Decoder) doInterpret() {
	// Interpretations
	if i.interpretFunc != nil {
		var r cue.Runtime
		i.file = i.File()
		inst, err := r.CompileFile(i.file)
		if err != nil {
			i.err = err
			return
		}
		i.file, i.id, i.err = i.interpretFunc(inst)
	}
}

func toFile(x ast.Expr) *ast.File {
	return internal.ToFile(x)
}

func valueToFile(v cue.Value) *ast.File {
	return internal.ToFile(v.Syntax())
}

func (i *Decoder) File() *ast.File {
	if i.file != nil {
		return i.file
	}
	return toFile(i.expr)
}

func (i *Decoder) Err() error {
	if i.err == io.EOF {
		return nil
	}
	return i.err
}

func (i *Decoder) Close() {
	i.closer.Close()
}

type Config struct {
	Mode filetypes.Mode

	// Out specifies an overwrite destination.
	Out    io.Writer
	Stdin  io.Reader
	Stdout io.Writer

	PkgName string // package name for files to generate

	Force     bool // overwrite existing files.
	Strict    bool
	Stream    bool // will potentially write more than one document per file
	AllErrors bool

	EscapeHTML bool
	ProtoPath  []string
	Format     []format.Option
	ParseFile  func(name string, src interface{}) (*ast.File, error)
}

// NewDecoder returns a stream of non-rooted data expressions. The encoding
// type of f must be a data type, but does not have to be an encoding that
// can stream. stdin is used in case the file is "-".
func NewDecoder(f *build.File, cfg *Config) *Decoder {
	if cfg == nil {
		cfg = &Config{}
	}
	i := &Decoder{filename: f.Filename, cfg: cfg}
	i.next = func() (ast.Expr, error) {
		if i.err != nil {
			return nil, i.err
		}
		return nil, io.EOF
	}

	if file, ok := f.Source.(*ast.File); ok {
		i.file = file
		i.closer = ioutil.NopCloser(strings.NewReader(""))
		i.validate(file, f)
		return i
	}

	rc, err := reader(f, cfg.Stdin)
	i.closer = rc
	i.err = err
	if err != nil {
		return i
	}

	// For now we assume that all encodings require UTF-8. This will not be the
	// case for some binary protocols. We need to exempt those explicitly here
	// once we introduce them.
	// TODO: this code also allows UTF16, which is too permissive for some
	// encodings. Switch to unicode.UTF8Sig once available.
	t := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	r := transform.NewReader(rc, t)

	switch f.Interpretation {
	case "":
	case build.Auto:
		openAPI := openAPIFunc(cfg, f)
		jsonSchema := jsonSchemaFunc(cfg, f)
		i.interpretFunc = func(inst *cue.Instance) (file *ast.File, id string, err error) {
			switch i.interpretation = Detect(inst.Value()); i.interpretation {
			case build.JSONSchema:
				return jsonSchema(inst)
			case build.OpenAPI:
				return openAPI(inst)
			}
			return i.file, "", i.err
		}
	case build.OpenAPI:
		i.interpretation = build.OpenAPI
		i.interpretFunc = openAPIFunc(cfg, f)
	case build.JSONSchema:
		i.interpretation = build.JSONSchema
		i.interpretFunc = jsonSchemaFunc(cfg, f)
	default:
		i.err = fmt.Errorf("unsupported interpretation %q", f.Interpretation)
	}

	path := f.Filename
	switch f.Encoding {
	case build.CUE:
		if cfg.ParseFile == nil {
			i.file, i.err = parser.ParseFile(path, r, parser.ParseComments)
		} else {
			i.file, i.err = cfg.ParseFile(path, r)
		}
		i.validate(i.file, f)
		if i.err == nil {
			i.doInterpret()
		}
	case build.JSON, build.JSONL:
		i.next = json.NewDecoder(nil, path, r).Extract
		i.Next()
	case build.YAML:
		d, err := yaml.NewDecoder(path, r)
		i.err = err
		i.next = d.Decode
		i.Next()
	case build.Text:
		b, err := ioutil.ReadAll(r)
		i.err = err
		i.expr = ast.NewString(string(b))
	case build.Protobuf:
		paths := &protobuf.Config{
			Paths:   cfg.ProtoPath,
			PkgName: cfg.PkgName,
		}
		i.file, i.err = protobuf.Extract(path, r, paths)
	default:
		i.err = fmt.Errorf("unsupported encoding %q", f.Encoding)
	}

	return i
}

func jsonSchemaFunc(cfg *Config, f *build.File) interpretFunc {
	return func(i *cue.Instance) (file *ast.File, id string, err error) {
		id = f.Tags["id"]
		if id == "" {
			id, _ = i.Lookup("$id").String()
		}
		if id != "" {
			u, err := url.Parse(id)
			if err != nil {
				return nil, "", errors.Wrapf(err, token.NoPos, "invalid id")
			}
			u.Scheme = ""
			id = strings.TrimPrefix(u.String(), "//")
		}
		cfg := &jsonschema.Config{
			ID:      id,
			PkgName: cfg.PkgName,

			Strict: cfg.Strict,
		}
		file, err = jsonschema.Extract(i, cfg)
		// TODO: simplify currently erases file line info. Reintroduce after fix.
		// file, err = simplify(file, err)
		return file, id, err
	}
}

func openAPIFunc(c *Config, f *build.File) interpretFunc {
	cfg := &openapi.Config{PkgName: c.PkgName}
	return func(i *cue.Instance) (file *ast.File, id string, err error) {
		file, err = openapi.Extract(i, cfg)
		// TODO: simplify currently erases file line info. Reintroduce after fix.
		// file, err = simplify(file, err)
		return file, "", err
	}
}

func reader(f *build.File, stdin io.Reader) (io.ReadCloser, error) {
	switch s := f.Source.(type) {
	case nil:
		// Use the file name.
	case string:
		return ioutil.NopCloser(strings.NewReader(s)), nil
	case []byte:
		return ioutil.NopCloser(bytes.NewReader(s)), nil
	case *bytes.Buffer:
		// is io.Reader, but it needs to be readable repeatedly
		if s != nil {
			return ioutil.NopCloser(bytes.NewReader(s.Bytes())), nil
		}
	default:
		return nil, fmt.Errorf("invalid source type %T", f.Source)
	}
	// TODO: should we allow this?
	if f.Filename == "-" {
		return ioutil.NopCloser(stdin), nil
	}
	return os.Open(f.Filename)
}

func shouldValidate(i *filetypes.FileInfo) bool {
	// TODO: We ignore attributes for now. They should be enabled by default.
	return false ||
		!i.Definitions ||
		!i.Data ||
		!i.Optional ||
		!i.Constraints ||
		!i.References ||
		!i.Cycles ||
		!i.KeepDefaults ||
		!i.Incomplete ||
		!i.Imports ||
		!i.Docs
}

type validator struct {
	allErrors bool
	count     int
	errs      errors.Error
	fileinfo  *filetypes.FileInfo
}

func (d *Decoder) validate(f *ast.File, b *build.File) {
	if d.err != nil {
		return
	}
	fi, err := filetypes.FromFile(b, filetypes.Input)
	if err != nil {
		d.err = err
		return
	}
	if !shouldValidate(fi) {
		return
	}

	v := validator{fileinfo: fi, allErrors: d.cfg.AllErrors}
	ast.Walk(f, v.validate, nil)
	d.err = v.errs
}

func (v *validator) validate(n ast.Node) bool {
	if v.count > 10 {
		return false
	}

	i := v.fileinfo

	// TODO: Cycles

	ok := true
	check := func(n ast.Node, option bool, s string, cond bool) {
		if !option && cond {
			v.errs = errors.Append(v.errs, errors.Newf(n.Pos(),
				"%s not allowed in %s mode", s, v.fileinfo.Form))
			v.count++
			ok = false
		}
	}

	// For now we don't make any distinction between these modes.

	constraints := i.Constraints && i.Incomplete && i.Optional && i.References

	check(n, i.Docs, "comments", len(ast.Comments(n)) > 0)

	switch x := n.(type) {
	case *ast.CommentGroup:
		check(n, i.Docs, "comments", len(ast.Comments(n)) > 0)
		return false

	case *ast.ImportDecl, *ast.ImportSpec:
		check(n, i.Imports, "imports", true)

	case *ast.Field:
		check(n, i.Definitions, "definitions",
			x.Token == token.ISA || internal.IsDefinition(x.Label))
		check(n, i.Data, "regular fields", internal.IsRegularField(x))
		check(n, constraints, "optional fields", x.Optional != token.NoPos)

		_, _, err := ast.LabelName(x.Label)
		check(n, constraints, "optional fields", err != nil)

		check(n, i.Attributes, "attributes", len(x.Attrs) > 0)
		ast.Walk(x.Value, v.validate, nil)
		return false

	case *ast.UnaryExpr:
		switch x.Op {
		case token.MUL:
			check(n, i.KeepDefaults, "default values", true)
		case token.SUB, token.ADD:
			// The parser represents negative numbers as an unary expression.
			// Allow one `-` or `+`.
			_, ok := x.X.(*ast.BasicLit)
			check(n, constraints, "expressions", !ok)
		case token.LSS, token.LEQ, token.EQL, token.GEQ, token.GTR,
			token.NEQ, token.NMAT, token.MAT:
			check(n, constraints, "constraints", true)
		default:
			check(n, constraints, "expressions", true)
		}

	case *ast.BinaryExpr, *ast.ParenExpr, *ast.IndexExpr, *ast.SliceExpr,
		*ast.CallExpr, *ast.Comprehension, *ast.ListComprehension,
		*ast.Interpolation:
		check(n, constraints, "expressions", true)

	case *ast.Ellipsis:
		check(n, constraints, "ellipsis", true)

	case *ast.Ident, *ast.SelectorExpr, *ast.Alias, *ast.LetClause:
		check(n, i.References, "references", true)

	default:
		// Other types are either always okay or handled elsewhere.
	}
	return ok
}

// simplify reformats a File. To be used as a wrapper for Extract functions.
//
// It currently does so by formatting the file using fmt.Format and then
// reparsing it. This is not ideal, but the package format does not provide a
// way to do so differently.
func simplify(f *ast.File, err error) (*ast.File, error) {
	if err != nil {
		return nil, err
	}
	// This needs to be a function that modifies f in order to maintain line
	// number information.
	b, err := format.Node(f, format.Simplify())
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(f.Filename, b, parser.ParseComments)
}
