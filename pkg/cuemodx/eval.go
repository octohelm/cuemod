package cuemodx

import (
	"context"
	"fmt"
	"os"

	"cuelang.org/go/cue/load"
	cueparser "cuelang.org/go/cue/parser"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/octohelm/cuemod/pkg/cuex"
)

func EvalContextWithPatches(ctx context.Context, r *cuemod.Context, inputs []string, options ...cuex.EvalOptionFunc) ([]byte, error) {
	files := make([]string, 0)
	inputSources := map[string]load.Source{}

	c := 0
	for i := range inputs {
		input := inputs[i]
		if input == "" {
			continue
		}

		var s load.Source
		if input[0] == '{' {
			s = load.FromString(input)
			input = fmt.Sprintf("./input_____%d.cue", i)
		} else {
			info, err := os.Stat(input)
			if err != nil {
				return nil, err
			}

			if info.IsDir() {
				s = load.FromString(fmt.Sprintf(`
import t "%s"
t
`, r.CompletePath(input)))
				input = fmt.Sprintf("./input_%d.cue", i)
			} else {
				src, err := os.ReadFile(input)
				if err != nil {
					return nil, err
				}
				f, err := cueparser.ParseFile(input, src)
				if err != nil {
					return nil, err
				}
				s = load.FromFile(f)
			}
		}

		input = r.CompletePath(input)

		files = append(files, input)
		inputSources[input] = s

		c++
	}

	inst := r.Build(ctx, files, cuemod.OptOverlay(inputSources))

	results, err := cuex.Eval(inst, options...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func EvalContext(ctx context.Context, r *cuemod.Context, filename string, options ...cuex.EvalOptionFunc) ([]byte, error) {
	filename = r.CompletePath(filename)
	inst := r.Build(ctx, []string{filename})
	results, err := cuex.Eval(inst, options...)
	if err != nil {
		return nil, err
	}
	return results, nil
}
