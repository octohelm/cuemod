package kube

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"golang.org/x/term"
)

var interactive = term.IsTerminal(int(os.Stdout.Fd()))

type ShowOpts struct {
	Output string `name:"output,o" usage:"output base dir"`
}

func (l *LoadResult) Show(opts ShowOpts) error {
	buf := bytes.NewBuffer(nil)

	for i := range l.Resources {
		r := l.Resources[i]

		filename := fmt.Sprintf(
			"%s/%s/%s.%s.%s.yaml",
			l.Release.Namespace,
			l.Release.Name,
			r.GetName(),
			r.GetObjectKind().GroupVersionKind().Kind,
			strings.Replace(r.GetObjectKind().GroupVersionKind().GroupVersion().String(), "/", ".", -1),
		)

		data, _ := yaml.Marshal(r)

		if opts.Output != "" {
			if err := writeFile(filepath.Join(opts.Output, filename), data); err != nil {
				return err
			}
		} else {
			buf.Write(data)
		}
	}

	if buf.Len() == 0 {
		return nil
	}

	return pageln(buf.String())
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
