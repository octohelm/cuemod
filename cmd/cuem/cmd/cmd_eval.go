package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/octohelm/cuemod/pkg/cli"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/octohelm/cuemod/pkg/cuex"
)

func init() {
	app.Add(&Eval{})
}

type EvalFlags struct {
	Output string `flag:"output,o" desc:"output filename and fmt"`
	Write  bool   `flag:"write,w" desc:"write"`
}

type Eval struct {
	cli.Name `args:"CUE_FILES..." desc:"evalWithPatches files"`
	EvalFlags
}

func (opts *Eval) Run(ctx context.Context, args []string) error {
	results, err := evalWithPatches(ctx, args, cuex.WithEncodingFromFileExt(filepath.Ext(opts.Output)))
	if err != nil {
		return err
	}

	if opts.Output != "" && opts.Write {
		if err := writeFile(opts.Output, results); err != nil {
			return err
		}
	} else {
		_, _ = io.Copy(os.Stdout, bytes.NewBuffer(results))
	}

	return nil
}

func evalWithPatches(ctx context.Context, fileOrPatches []string, options ...cuex.EvalOptionFunc) ([]byte, error) {
	runtime := cuemod.FromContext(ctx)

	cwd, _ := os.Getwd()
	for i := range fileOrPatches {
		if fileOrPatches[i][0] == '.' {
			fileOrPatches[i] = filepath.Join(cwd, fileOrPatches[i])
		}
	}

	return runtime.EvalWithPatches(ctx, fileOrPatches, options...)
}

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
