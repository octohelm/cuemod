package cuemod

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/octohelm/cuemod/pkg/cue"

	. "github.com/onsi/gomega"

	"github.com/go-courier/logr"
)

func TestRuntime(t *testing.T) {
	cwd, _ := os.Getwd()

	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	ctx = WithOpts(ctx, OptVerbose(true))

	t.Run("mod a", func(t *testing.T) {
		r := RuntimeFor(filepath.Join(cwd, "./testdata/a"))
		_ = os.RemoveAll(filepath.Join(r.CueModRoot()))

		t.Run("Eval", func(t *testing.T) {
			_, err := r.Eval(ctx, []string{"./main.cue"}, cue.YAML)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(r.mod.Require["k8s.io/api"].Version).To(Equal("v0.20.4"))
		})
	})

	t.Run("mod b", func(t *testing.T) {
		r := RuntimeFor(filepath.Join(cwd, "./testdata/b"))
		_ = os.RemoveAll(filepath.Join(r.CueModRoot()))

		t.Run("ListCue", func(t *testing.T) {
			t.Run("one dir", func(t *testing.T) {
				files, err := r.ListCue(".")

				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(files).To(HaveLen(2))
			})

			t.Run("all", func(t *testing.T) {
				files, err := r.ListCue("./...")

				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(len(files) > 2).To(BeTrue())

				t.Run("Format", func(t *testing.T) {
					err := cue.FormatFiles(ctx, files, cue.FormatOpts{
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
			ret, err := r.Eval(ctx, []string{"./main.cue"}, cue.YAML)
			NewWithT(t).Expect(err).To(BeNil())
			t.Log(string(ret[0]))
		})
	})
}
