package term

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

// Colordiff colorizes unified diff output (diff -u -N)
func Colordiff(d []byte) *bytes.Buffer {
	exps := map[string]func(b []byte) bool{
		"add":  regexp.MustCompile(`^\+.*`).Match,
		"del":  regexp.MustCompile(`^\-.*`).Match,
		"head": regexp.MustCompile(`^diff -u -N.*`).Match,
		"hid":  regexp.MustCompile(`^@.*`).Match,
	}

	buf := bytes.Buffer{}
	lines := bytes.Split(d, []byte("\n"))

	for _, l := range lines {
		switch {
		case exps["add"](l):
			_, _ = color.New(color.FgGreen).Fprintln(&buf, string(l))
		case exps["del"](l):
			_, _ = color.New(color.FgRed).Fprintln(&buf, string(l))
		case exps["head"](l):
			_, _ = color.New(color.FgBlue, color.Bold).Fprintln(&buf, string(l))
		case exps["hid"](l):
			_, _ = color.New(color.FgMagenta, color.Bold).Fprintln(&buf, string(l))
		default:
			_, _ = fmt.Fprintln(&buf, string(l))
		}
	}

	return &buf
}
