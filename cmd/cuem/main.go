package main

import (
	"context"
	"os"

	"github.com/go-courier/logr"
)

func main() {
	ctx := logr.WithLogger(context.Background(), log)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
