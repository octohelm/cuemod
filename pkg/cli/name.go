package cli

import (
	"fmt"
	"strings"
)

type Name struct {
	Name      string
	Desc      string
	ValidArgs ValidArgs

	subcommands []Command
	parent      Command
}

func (n *Name) Naming() *Name {
	return n
}

func (n *Name) Add(c Command) {
	n.subcommands = append(n.subcommands, c)
}

func ParseValidArgs(s string) ValidArgs {
	if s == "" {
		return ValidArgs{}
	}

	v := make(ValidArgs, 0)

	args := strings.Split(s, " ")

	for i := range args {
		arg := strings.TrimSpace(args[i])
		v = append(v, arg)
	}

	return v
}

type ValidArgs []string

func (as ValidArgs) HasVariadic() bool {
	for _, a := range as {
		if strings.HasSuffix(a, "...") {
			return true
		}
	}
	return false
}

func (as ValidArgs) Validate(args []string) error {
	if as.HasVariadic() {
		if len(args) < len(as) {
			return fmt.Errorf("requires at least %d arg(s), only received %d", len(as), len(args))
		}
	}
	if len(as) != len(args) {
		return fmt.Errorf("accepts %d arg(s), received %d", len(as), len(args))
	}
	return nil

}
