package modfile

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
)

const ModFilename = "cue.mod/module.cue"

type ModVersion struct {
	Version string
	VcsRef  string
}

func (mv ModVersion) Exactly() bool {
	return mv.Version != "" && mv.VcsRef == ""
}

type Requirement struct {
	ModVersion
	Indirect bool
}

type ReplaceTarget struct {
	VersionedPathIdentity
	Import string
}

type ModFile struct {
	// Module name
	Module string

	// Replace
	// version limit
	Replace map[VersionedPathIdentity]ReplaceTarget
	// Require same as go root
	// require { module: version }
	// indirect require { module:: version }
	Require map[string]Requirement

	comments map[string][]*ast.CommentGroup
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

		m.writeRequires(buf, modules, func(r *Requirement) bool {
			return !r.Indirect
		})

		m.writeRequires(buf, modules, func(r *Requirement) bool {
			return r.Indirect
		})
	}

	if len(m.Replace) > 0 {
		replacements := make([]VersionedPathIdentity, 0)

		for r := range m.Replace {
			replacements = append(replacements, r)
		}

		sort.Slice(replacements, func(i, j int) bool {
			return replacements[i].Path < replacements[j].Path
		})

		m.writeReplaces(buf, replacements, func(r *ReplaceTarget) bool {
			return r.Import == ""
		})

		m.writeReplaces(buf, replacements, func(r *ReplaceTarget) bool {
			return r.Import != ""
		})
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	return data
}

func (m *ModFile) writeReplaces(w io.Writer, replacements []VersionedPathIdentity, filter func(replaceTarget *ReplaceTarget) bool) {
	fields := make([]interface{}, 0)

	for _, replaceFrom := range replacements {
		replaceTarget := m.Replace[replaceFrom]

		if !filter(&replaceTarget) {
			continue
		}

		i := replaceFrom.String()

		f := &ast.Field{Label: ast.NewString(i)}

		if replaceTarget.Path == replaceFrom.Path {
			f.Value = ast.NewString((&VersionedPathIdentity{ModVersion: replaceTarget.ModVersion}).String())
		} else {
			f.Value = ast.NewString(replaceTarget.String())
		}

		if cg, ok := m.comments["replace://"+i]; ok {
			for i := range cg {
				f.AddComment(cg[i])
			}
		}

		if replaceTarget.Import != "" {
			f.Attrs = append(f.Attrs, &ast.Attribute{Text: attr("import", replaceTarget.Import)})
		}

		fields = append(fields, f)
	}

	if len(fields) == 0 {
		return
	}

	data, _ := format.Node(ast.NewStruct(fields...))

	_, _ = fmt.Fprintf(w, `
replace: %s
`, string(data))
}

func (m *ModFile) writeRequires(w io.Writer, modules []string, filter func(r *Requirement) bool) {
	fields := make([]interface{}, 0)

	// direct require
	for _, module := range modules {
		r := m.Require[module]

		if !filter(&r) {
			continue
		}

		f := &ast.Field{Label: ast.NewString(module)}
		f.Value = ast.NewString(r.Version)

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

	if len(fields) > 0 {
		data, _ := format.Node(ast.NewStruct(fields...))
		_, _ = fmt.Fprintf(w, `
require: %s
`, string(data))
	}
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
