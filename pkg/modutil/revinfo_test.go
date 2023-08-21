package modutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"golang.org/x/mod/module"
)

func TestRevInfoFromDir(t *testing.T) {
	cwd, _ := os.Getwd()
	info, err := RevInfoFromDir(context.Background(), cwd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info.Version)
	fmt.Println(module.PseudoVersionBase(info.Version))
}
