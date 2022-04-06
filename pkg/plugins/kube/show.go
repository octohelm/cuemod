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

func (l *LoadResult) Show(opts ShowOpts) (err error) {
	var w io.Writer

	writeToFiles := opts.Output != ""
	writeToSingleFile := strings.HasSuffix(opts.Output, ".yaml")

	if writeToSingleFile {
		if err := os.MkdirAll(filepath.Dir(opts.Output), os.ModePerm); err != nil {
			return err
		}
		file, _ := os.Create(opts.Output)
		defer file.Close()
		w = file
	} else if !writeToFiles {
		buf := bytes.NewBuffer(nil)
		defer func() {
			if buf.Len() > 0 {
				err = pageln(buf.String())
			}
		}()
		w = buf
	}

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

		if writeToFiles && !writeToSingleFile {
			if err := writeFile(filepath.Join(opts.Output, filename), data); err != nil {
				return err
			}
		} else {
			_, _ = io.WriteString(w, "---\n")
			_, _ = io.WriteString(w, "# "+filename+"\n")
			_, _ = w.Write(data)
		}
	}

	return nil
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
