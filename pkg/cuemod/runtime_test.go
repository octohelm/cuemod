package cuemod

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cuemod/format"
	. "github.com/onsi/gomega"
)

func TestRuntime(t *testing.T) {
	cwd, _ := os.Getwd()

	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	ctx = WithOpts(ctx, OptVerbose(true))

	t.Run("mod a", func(t *testing.T) {
		r := RuntimeFor(filepath.Join(cwd, "./testdata/a"))
		_ = r.Cleanup()

		t.Run("Eval", func(t *testing.T) {
			data, err := r.Eval(ctx, ".", JSON)
			NewWithT(t).Expect(err).To(BeNil())
			fmt.Println(string(data))
			NewWithT(t).Expect(r.mod.Require["k8s.io/api"].Version).To(Equal("v0.20.5"))
		})

		t.Run("Eval from exported single file", func(t *testing.T) {
			data, err := r.Eval(ctx, ".", CUE)
			NewWithT(t).Expect(err).To(BeNil())
			//fmt.Println(string(data))

			_ = os.WriteFile("../../_output/debug.cue", data, os.ModePerm)

			//f, err := cueparser.ParseFile("main.cue", data)
			//NewWithT(t).Expect(err).To(BeNil())

			//inst := build.NewContext().NewInstance(".", nil)
			//_ = inst.AddSyntax(f)
			//
			//ret, err := Eval(inst, JSON)
			//NewWithT(t).Expect(err).To(BeNil())
			//fmt.Println(string(ret))
		})
	})

	t.Run("mod b", func(t *testing.T) {
		r := RuntimeFor(filepath.Join(cwd, "./testdata/b"))
		_ = r.Cleanup()

		t.Run("ListCue", func(t *testing.T) {
			t.Run("one dir", func(t *testing.T) {
				files, err := r.ListCue(".")

				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(files).To(HaveLen(1))
			})

			t.Run("all", func(t *testing.T) {
				files, err := r.ListCue("./...")

				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(len(files) > 1).To(BeTrue())

				t.Run("Format", func(t *testing.T) {
					err := format.FormatFiles(ctx, files, format.FormatOpts{
						ReplaceFile: true,
						PrintNames:  true,
					})
					NewWithT(t).Expect(err).To(BeNil())
				})
			})
		})

		t.Run("Get", func(t *testing.T) {
			err := r.Get(ctx, "./...")
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("Eval", func(t *testing.T) {
			ret, err := r.Eval(ctx, "./main.cue", YAML)
			NewWithT(t).Expect(err).To(BeNil())
			t.Log(string(ret))
		})
	})
}
