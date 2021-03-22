package tanka

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

var interactive = term.IsTerminal(int(os.Stdout.Fd()))

type ShowOpts struct {
	Output string `name:"output,o" usage:"output base dir"`
}

func (l *LoadResult) Show(opts ShowOpts) error {
	if opts.Output != "" {
		for i := range l.Resources {
			r := l.Resources[i]

			filename := fmt.Sprintf(
				"%s/%s/%s.%s.%s.yaml",
				l.Release.Metadata.Name,
				l.Release.Metadata.Namespace,
				r.Metadata().Name(),
				r.Kind(),
				strings.Replace(r.APIVersion(), "/", ".", -1),
			)

			if err := writeFile(filepath.Join(opts.Output, filename), []byte(r.String())); err != nil {
				return err
			}
		}
		return nil
	}

	return pageln(l.Resources.String())
}

func pageln(i ...interface{}) error {
	return fPageln(strings.NewReader(fmt.Sprint(i...)))
}

func fPageln(r io.Reader) error {
	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		// --RAW-CONTROL-CHARS  Honors colors from diff. Must be in all caps, otherwise display issues occur.
		// --quit-if-one-screen Closer to the git experience.
		// --no-init            Don't clear the screen when exiting.
		pager = "less --RAW-CONTROL-CHARS --quit-if-one-screen --no-init"
	}

	if interactive && pager != "" {
		cmd := exec.Command("sh", "-c", pager)
		cmd.Stdin = r
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			// Fallthrough on failure so that the contents of the reader are copied to stdout.
		} else {
			return nil
		}
	}

	_, err := io.Copy(os.Stdout, r)
	return err
}
