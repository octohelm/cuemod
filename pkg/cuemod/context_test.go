package cuemod_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/octohelm/cuemod/pkg/cuemodx"
	"github.com/octohelm/cuemod/pkg/cuex"
	"github.com/octohelm/cuemod/pkg/cuex/format"
	. "github.com/onsi/gomega"
)

func TestContext(t *testing.T) {
	cwd, _ := os.Getwd()

	ctx := logr.WithLogger(context.Background(), logr.StdLogger())
	ctx = cuemod.WithOpts(ctx, cuemod.OptVerbose(true))

	t.Run("mod a", func(t *testing.T) {
		r := cuemod.ContextFor(filepath.Join(cwd, "./testdata/a"))
		_ = r.Cleanup()

		t.Run("EvalContext", func(t *testing.T) {
			data, err := cuemodx.EvalContext(ctx, r, ".", cuex.WithEncoding(cuex.JSON))
			NewWithT(t).Expect(err).To(BeNil())
			fmt.Println(string(data))
			NewWithT(t).Expect(r.Mod.Require["k8s.io/api"].Version).To(Equal("v0.24.1"))
		})

		t.Run("EvalContext from exported single file", func(t *testing.T) {
			data, err := cuemodx.EvalContext(ctx, r, ".", cuex.WithEncoding(cuex.CUE))
			NewWithT(t).Expect(err).To(BeNil())

			_ = os.WriteFile("../../_output/debug.cue", data, os.ModePerm)

			inst, _ := cuex.InstanceFromRaw(data)
			ret, err := cuex.Eval(inst, cuex.WithEncoding(cuex.JSON))
			NewWithT(t).Expect(err).To(BeNil())
			fmt.Println(string(ret))
		})
	})

	t.Run("mod b", func(t *testing.T) {
		r := cuemod.ContextFor(filepath.Join(cwd, "./testdata/b"))
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
				NewWithT(t).Expect(len(files) > 0).To(BeTrue())

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
	})

	t.Run("mod dagger-example", func(t *testing.T) {
		t.Run("direct", func(t *testing.T) {
			r := cuemod.ContextFor(filepath.Join(cwd, "./testdata/dagger"))
			t.Log(r.Cleanup())

			t.Run("Get", func(t *testing.T) {
				err := r.Get(ctx, "./...")
				NewWithT(t).Expect(err).To(BeNil())
			})
		})

		t.Run("with custom replace", func(t *testing.T) {
			r := cuemod.ContextFor(filepath.Join(cwd, "./testdata/dagger-replace"))
			t.Log(r.Cleanup())

			t.Run("Get", func(t *testing.T) {
				err := r.Get(ctx, "./...")
				NewWithT(t).Expect(err).To(BeNil())
			})
		})
	})

}
