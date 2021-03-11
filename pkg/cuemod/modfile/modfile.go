package modfile

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"cuelang.org/go/cue/format"

	"cuelang.org/go/cue/ast"
)

const ModFilename = "cue.mod/module.cue"

type ModVersion struct {
	Version    string
	VcsVersion string
}

type Require struct {
	ModVersion
	Indirect bool
}

type ModFile struct {
	// Module name
	Module string

	// Replace
	// version limit
	Replace map[PathMayWithVersion]ReplaceTarget
	// Require same as go root
	// require { module: version }
	// indirect require { module:: version }
	Require map[string]Require

	comments map[string][]*ast.CommentGroup
}

type ReplaceTarget struct {
	PathMayWithVersion
	Gen string
}

func (m *ModFile) String() string {
	return string(m.Bytes())
}

func (m *ModFile) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(buf, "module: %s\n", strconv.Quote(m.Module))

	if len(m.Require) > 0 {
		modules := make([]string, 0)

		for module := range m.Require {
			modules = append(modules, module)
		}

		sort.Strings(modules)

		fields := make([]interface{}, 0)

		for _, module := range modules {
			r := m.Require[module]

			f := &ast.Field{Label: ast.NewString(module)}
			f.Value = ast.NewString(r.Version)

			if r.VcsVersion != "" && r.VcsVersion != r.Version {
				f.Attrs = append(f.Attrs, &ast.Attribute{Text: attr("vcs", r.VcsVersion)})
			}

			if r.Indirect {
				f.Attrs = append(f.Attrs, &ast.Attribute{Text: attr("indirect")})
			}

			if cg, ok := m.comments["require://"+module]; ok {
				for i := range cg {
					f.AddComment(cg[i])
				}
			}

			fields = append(fields, f)
		}

		data, _ := format.Node(ast.NewStruct(fields...))
		_, _ = fmt.Fprintf(buf, `
require: %s
`, string(data))
	}

	if len(m.Replace) > 0 {
		replacements := make([]PathMayWithVersion, 0)

		for r := range m.Replace {
			replacements = append(replacements, r)
		}

		sort.Slice(replacements, func(i, j int) bool {
			return replacements[i].Path < replacements[j].Path
		})

		fields := make([]interface{}, 0)

		for _, replaceFrom := range replacements {
			i := replaceFrom.String()

			replaceTarget := m.Replace[replaceFrom]

			f := &ast.Field{Label: ast.NewString(i)}

			if replaceTarget.Path == replaceFrom.Path {
				f.Value = ast.NewString((&PathMayWithVersion{Version: replaceTarget.Version}).String())
			} else {
				f.Value = ast.NewString(replaceTarget.String())
			}

			if cg, ok := m.comments["replace://"+i]; ok {
				for i := range cg {
					f.AddComment(cg[i])
				}
			}

			if replaceTarget.Gen != "" {
				f.Attrs = append(f.Attrs, &ast.Attribute{Text: attr("gen", replaceTarget.Gen)})
			}

			fields = append(fields, f)
		}

		data, _ := format.Node(ast.NewStruct(fields...))

		_, _ = fmt.Fprintf(buf, `
replace: %s
`, string(data))
	}

	data, err := format.Source(buf.Bytes(), format.Simplify())
	if err != nil {
		panic(err)
	}

	return data
}

func attr(key string, values ...string) string {
	buf := bytes.NewBufferString("@")
	buf.WriteString(key)
	buf.WriteString("(")

	for i := range values {
		if i > 0 {
			buf.WriteString(",")
		}

		buf.WriteString(strconv.Quote(values[i]))
	}

	buf.WriteString(")")

	return buf.String()
}
