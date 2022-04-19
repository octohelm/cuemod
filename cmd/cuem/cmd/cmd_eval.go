package cmd

import (
	"bytes"
	"context"
	"fmt"
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
	cli.Name `args:"CUE_FILES..." desc:"eval files"`
	EvalFlags
}

func (opts *Eval) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing input")
	}

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, args[0])

	format := cuex.YAML

	switch v := filepath.Ext(opts.Output); v {
	case ".yaml":
		format = cuex.YAML
	case ".json":
		format = cuex.JSON
	case ".cue":
		format = cuex.CUE
	default:
		panic(fmt.Errorf("unsupport output format %s", v))
	}

	results, err := cuemod.FromContext(ctx).Eval(ctx, path, format)
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

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
