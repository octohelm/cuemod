package helm

import (
	"context"
	encodingjson "encoding/json"
	"os"
	"path/filepath"
	"strings"

	cuetoken "cuelang.org/go/cue/token"

	"cuelang.org/go/encoding/json"

	cueast "cuelang.org/go/cue/ast"
	"github.com/octohelm/cuemod/pkg/extractor/core"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
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
	return "helm"
}

func (Extractor) Detect(ctx context.Context, src string) (bool, map[string]string) {
	f, err := os.Lstat(filepath.Join(src, "Chart.yaml"))
	if err == nil {
		return !f.IsDir(), nil
	}
	return false, nil
}

func (e *Extractor) Extract(ctx context.Context, src string) (files []*cueast.File, err error) {
	c, err := loader.LoadDir(src)
	if err != nil {
		return nil, errors.Wrapf(err, "load helm chart failed from %s", src)
	}

	return e.extractChart(c)
}

func (e *Extractor) extractChart(helmChart *chart.Chart) (files []*cueast.File, err error) {
	appendFile := func(scope string, field cueast.Expr) {
		f := &cueast.File{}
		f.Filename = scope + "_gen.cue"
		f.Decls = []cueast.Decl{
			&cueast.Package{Name: cueast.NewIdent("chart")},
			&cueast.Field{
				Label: cueast.NewIdent(scope),
				Value: field,
			},
		}

		files = append(files, f)
	}

	appendTemplateFile := func(filename string, data []byte) {
		f := &cueast.File{}
		f.Filename = strings.ReplaceAll(filename, "/", "__") + "_gen.cue"
		f.Decls = []cueast.Decl{
			&cueast.Package{Name: cueast.NewIdent("chart")},
			&cueast.Field{
				Label: cueast.NewIdent("_files"),
				Value: &cueast.StructLit{
					Lbrace: cuetoken.Blank.Pos(),
					Rbrace: cuetoken.Blank.Pos(),
					Elts: []cueast.Decl{
						&cueast.Field{
							Label: cueast.NewString(filename),
							Value: newBytes(data),
						},
					},
				},
			},
		}

		files = append(files, f)
	}

	valuesFile, err := fromValues(helmChart.Values)
	if err != nil {
		return nil, err
	}
	appendFile("values", valuesFile)

	metadataFile, err := fromValues(helmChart.Metadata)
	if err != nil {
		return nil, err
	}
	appendFile("metadata", metadataFile)

	helmFiles := make([]cueast.Expr, 0)
	for _, f := range append(helmChart.Templates, helmChart.Files...) {

		helmFiles = append(helmFiles, cueast.NewStruct(
			&cueast.Field{Label: cueast.NewString("name"), Value: cueast.NewString(f.Name)},
			&cueast.Field{Label: cueast.NewString("data"), Value: &cueast.IndexExpr{
				X:     cueast.NewIdent("_files"),
				Index: cueast.NewString(f.Name),
			}},
		))

		appendTemplateFile(f.Name, f.Data)
	}

	appendFile("files", cueast.NewList(helmFiles...))

	return
}

func newBytes(data []byte) *cueast.BasicLit {
	return &cueast.BasicLit{
		Kind:     cuetoken.STRING,
		ValuePos: cuetoken.NoPos,
		Value:    "'''\n" + strings.Replace(string(data), "\\", "\\\\", -1) + "'''",
	}
}

func fromValues(v interface{}) (cueast.Expr, error) {
	data, err := encodingjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.Extract("", data)
}
