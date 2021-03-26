package helm

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cuetoken "cuelang.org/go/cue/token"
	"k8s.io/apimachinery/pkg/util/yaml"

	cueast "cuelang.org/go/cue/ast"
	"github.com/octohelm/cuemod/pkg/extractor/core"
	"github.com/pkg/errors"

	apiextensions_v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func init() {
	core.Register(&Extractor{})
}

// Extractor from helm charts
//
// Targets:
// * gen values to value check
// * gen templates
type Extractor struct {
}

func (Extractor) Name() string {
	return "crd"
}

// never detect
func (Extractor) Detect(ctx context.Context, src string) (bool, map[string]string) {
	return false, nil
}

func (e *Extractor) Extract(ctx context.Context, src string) (files []*cueast.File, err error) {
	crdFiles, err := filepath.Glob(filepath.Join(src, "*.yaml"))
	if err != nil {
		return nil, errors.Wrapf(err, "find crd.yaml  failed from %s", src)
	}

	for i := range crdFiles {
		data, err := os.ReadFile(crdFiles[i])
		if err != nil {
			return nil, err
		}

		if trimmedContent := strings.TrimSpace(string(data)); trimmedContent != "" {
			decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(trimmedContent), 4096)

			for {
				crd := apiextensions_v1.CustomResourceDefinition{}

				if err := decoder.Decode(&crd); err != nil {
					if err == io.EOF {
						break
					}
					return nil, errors.Wrapf(err, "invalid crd failed: %s\n%s", crdFiles[i], trimmedContent)
				}

				if crd.Spec.Group == "" {
					continue
				}

				cueFile, err := e.fileFromCRD(&crd)
				if err != nil {
					return nil, err
				}

				files = append(files, cueFile)
			}
		}

	}

	return
}

func (e *Extractor) fileFromCRD(crd *apiextensions_v1.CustomResourceDefinition) (*cueast.File, error) {
	f := &cueast.File{}
	f.Filename = crd.Name + "_gen.cue"
	f.Decls = []cueast.Decl{
		&cueast.Package{Name: cueast.NewIdent("crd")},
	}

	decl := func(d cueast.Decl) {
		f.Decls = append(f.Decls, d)
	}

	decl(&cueast.Field{Label: cueast.NewIdent("#Group"), Value: cueast.NewString(crd.Spec.Group)})

	for _, v := range crd.Spec.Versions {
		decl(&cueast.Field{
			Label: cueast.NewIdent(v.Name),
			Value: &cueast.StructLit{
				Elts: []cueast.Decl{
					&cueast.Field{
						Label: cueast.NewIdent(crd.Spec.Names.Kind),
						Value: e.fromJSONSchema(v.Schema.OpenAPIV3Schema),
					},
				},
			}})
	}

	return f, nil
}

func (e Extractor) fromJSONSchema(s *apiextensions_v1.JSONSchemaProps) cueast.Expr {
	if len(s.AnyOf) > 0 {
		items := make([]cueast.Expr, len(s.AnyOf))

		for i := range items {
			items[i] = e.fromJSONSchema(&s.AnyOf[i])
		}

		return cueast.NewBinExpr(cuetoken.OR, items...)
	}

	switch s.Type {
	case "object":
		structLit := &cueast.StructLit{}

		//s.Required

		fields := make([]string, 0)
		required := map[string]bool{}

		for f := range s.Properties {
			fields = append(fields, f)
		}

		for _, f := range s.Required {
			required[f] = true
		}

		sort.Strings(fields)

		for _, f := range fields {
			p := s.Properties[f]

			field := &cueast.Field{Label: cueast.NewString(f), Value: e.fromJSONSchema(&p)}

			if p.Description != "" {
				addComments(field, &cueast.CommentGroup{Doc: true, List: []*cueast.Comment{
					toCueComment(p.Description),
				}})
			}

			if _, ok := required[f]; !ok {
				field.Token = cuetoken.COLON
				field.Optional = cuetoken.Blank.Pos()
			}

			structLit.Elts = append(structLit.Elts, field)
		}

		return structLit

	case "string":
		return cueast.NewIdent("string")
	case "integer":
		switch s.Format {
		case "int", "int8", "int16", "int32", "int64":
			return cueast.NewIdent(s.Format)
		}
		return cueast.NewIdent("int")
	case "number":
		switch s.Format {
		case "float":
			return cueast.NewIdent("float32")
		}
		return cueast.NewIdent("float64")
	case "boolean":
		return cueast.NewIdent("boolean")
	case "array":
		if s.Items == nil {
			return cueast.NewList(&cueast.Ellipsis{
				Type: any(),
			})
		}

		if s.Items.Schema != nil {
			elem := e.fromJSONSchema(s.Items.Schema)
			if elem == nil {
				return nil
			}
			return cueast.NewList(&cueast.Ellipsis{
				Type: elem,
			})
		}

		items := make([]cueast.Expr, len(s.Items.JSONSchemas))

		for i := range items {
			items[i] = e.fromJSONSchema(&s.Items.JSONSchemas[i])
		}

		return cueast.NewList()
	}

	return any()
}

func any() cueast.Expr {
	return cueast.NewIdent("_")
}

func addComments(node cueast.Node, comments ...*cueast.CommentGroup) {
	for i := range comments {
		cg := comments[i]
		if cg == nil {
			continue
		}
		cueast.AddComment(node, comments[i])
	}
}

func toCueComment(d string) *cueast.Comment {
	lines := strings.Split(d, "\n")

	c := &cueast.Comment{}

	buf := bytes.NewBuffer(nil)

	for i := range lines {
		if i > 0 {
			buf.WriteString("\n")
		}

		buf.WriteString("// ")
		buf.WriteString(lines[i])
	}

	c.Text = buf.String()

	return c
}
