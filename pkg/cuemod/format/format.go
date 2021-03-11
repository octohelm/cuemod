package format

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/format"
	"github.com/go-courier/logr"
)

// FormatOpts
type FormatOpts struct {
	// ReplaceFile when enabled, will write formatted code to same file
	ReplaceFile bool `name:"write,w" usage:"write result to (source) file instead of stdout"`
	// PrintNames when enabled, will print formatted result of each file
	PrintNames bool `name:"list,l" `
}

// FormatFiles format jsonnet files
func FormatFiles(ctx context.Context, files []string, opt FormatOpts) error {
	log := logr.FromContext(ctx)

	writeFile := func(file string, data []byte) error {
		f, _ := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
		defer f.Close()
		_, err := io.Copy(f, bytes.NewBuffer(data))
		return err
	}

	cwd, _ := os.Getwd()

	for i := range files {
		file := files[i]

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		file, _ = filepath.Rel(cwd, file)

		formatted, changed, err := Format(data)
		if err != nil {
			return err
		}

		if changed {
			if opt.PrintNames {
				log.Info(fmt.Sprintf("`%s` formatted.", file))
			}

			if opt.ReplaceFile {
				if err := writeFile(file, formatted); err != nil {
					return err
				}
			} else {
				fmt.Printf(`
// %s 

%s
`, file, formatted)
			}
		} else {
			if opt.PrintNames {
				log.Info(fmt.Sprintf("`%s` no changes.", file))
			}
		}
	}

	return nil
}

func Format(src []byte) ([]byte, bool, error) {
	formatted, err := format.Source(src, format.Simplify())
	if err != nil {
		return nil, false, err
	}
	return formatted, !bytes.Equal(formatted, src), nil
}
