package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/sumdb/dirhash"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"github.com/pkg/errors"
)

type Extractor interface {
	// Name
	// Extractor name
	// used in mod.cue @import("")
	Name() string
	// check dir should use extractor.
	// if matched, return deps <repo>:<version> too.
	Detect(ctx context.Context, src string) (bool, map[string]string)
	// Extract convert dir to cue codes
	Extract(ctx context.Context, src string) ([]*cueast.File, error)
}

var SharedExtractors = Extractors{}

func Register(extractor Extractor) {
	SharedExtractors.Register(extractor)
}

type Extractors map[string]Extractor

func (extractors Extractors) Register(extractor Extractor) {
	extractors[extractor.Name()] = extractor
}

func (extractors Extractors) Detect(ctx context.Context, src string) (string, map[string]string) {
	for _, e := range extractors {
		if ok, deps := e.Detect(ctx, src); ok {
			return e.Name(), deps
		}
	}
	return "", nil
}

func (extractors Extractors) ExtractToDir(ctx context.Context, name string, src string, gen string) error {
	if extractor, ok := extractors[name]; ok {

		return extractors.do(ctx, extractor, src, gen)
	}
	return errors.Errorf("unsupport extractor `%s`", name)
}

func (Extractors) do(ctx context.Context, extractor Extractor, src string, gen string) error {
	sumFile := filepath.Join(gen, ".sum")
	sum, _ := os.ReadFile(sumFile)

	origin, err := os.Readlink(src)
	if err == nil {
		src = origin
	}

	dirSum, err := dirhash.HashDir(src, "", dirhash.DefaultHash)
	if err != nil {
		return err
	}

	if string(sum) == dirSum {
		// skip when dirSum same
		return nil
	}

	currentFiles, err := filepath.Glob(filepath.Join(gen, "*_gen.cue"))
	if err != nil {
		return err
	}

	shouldDelete := map[string]bool{}

	for _, f := range currentFiles {
		shouldDelete[filepath.Base(f)] = true
	}

	files, err := extractor.Extract(ctx, src)
	if err != nil {
		return err
	}

	for i := range files {
		shouldDelete[files[i].Filename] = false

		if err := writeCueFile(ctx, extractor.Name(), gen, files[i]); err != nil {
			return err
		}
	}

	for filename, ok := range shouldDelete {
		if ok {
			if err := os.RemoveAll(filepath.Join(gen, filename)); err != nil {
				return err
			}
		}
	}

	return writeFile(sumFile, []byte(dirSum))
}

func writeCueFile(ctx context.Context, name string, dir string, f *cueast.File) error {
	filename := filepath.Join(dir, f.Filename)

	cueast.AddComment(f.Decls[0], &cueast.CommentGroup{
		Doc: true,
		List: []*cueast.Comment{
			{Text: "// DO NOT EDIT THIS FILE DIRECTLY."},
			{Text: fmt.Sprintf("// generated by %s extractor.", name)},
		},
	})

	data, err := format.Node(f, format.Simplify())
	if err != nil {
		return err
	}

	return writeFile(filename, data)
}

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
