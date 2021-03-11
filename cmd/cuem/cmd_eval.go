package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/errors"
	"github.com/octohelm/cuemod/pkg/cue"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		cmdEval(),
	)
}

type BuildOpts struct {
	Output string `name:"output,o" usage:"output filename"`
}

func cmdEval() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eval <input>",
		Short: "eval",
	}

	opts := BuildOpts{}

	return setupRun(cmd, &opts, func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing input")
		}

		results, err := runtime.Eval(ctx, args, cue.YAML)
		if err != nil {
			errors.Print(os.Stdout, err, nil)
			os.Exit(1)
			return nil
		}

		for i := range results {

			if o := opts.Output; o != "" {
				if err := writeFile(o, results[i]); err != nil {
					return err
				}
			} else {
				_, _ = io.Copy(os.Stdout, bytes.NewBuffer(results[i]))
			}
		}

		return nil
	})
}

func writeFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
