package cmd

import (
	"context"

	xslog "log/slog"

	"github.com/go-courier/logr"
	"github.com/go-courier/logr/slog"
	"github.com/octohelm/cuemod/pkg/cuemod"
)

type ProjectFlags struct {
	Root string `flag:"!project,p" desc:"project root"`
	V    int    `flag:"!verbose,v" desc:"verbose level"`
}

func (v ProjectFlags) PreRun(ctx context.Context) context.Context {
	ctx = cuemod.InjectContext(ctx, cuemod.ContextFor(v.Root))
	if v.V > 0 {
		return logr.WithLogger(ctx, slog.Logger(slog.Default()))
	}
	return logr.WithLogger(ctx, slog.Logger(xslog.New(xslog.Default().Handler())))
}
