package main

import (
	"context"

	"github.com/octohelm/cuemod/cmd/cuem/cmd"
)

func main() {
	if err := cmd.Run(context.Background()); err != nil {
		panic(err)
	}
}
