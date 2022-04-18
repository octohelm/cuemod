package cmd

import (
	"context"
	"os"

	"github.com/octohelm/cuemod/pkg/cuemod"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ProjectFlags struct {
	Root string `flag:"!project,p" desc:"project root"`
	V    int    `flag:"!verbose,v" desc:"verbose level"`
}

func (v ProjectFlags) PreRun(ctx context.Context) context.Context {
	ctx = cuemod.InjectContext(ctx, cuemod.ContextFor(v.Root))

	if v.V > 0 {
		return logr.NewContext(ctx, NewLogger(-v.V))
	}
	return logr.NewContext(ctx, NewLogger(0))
}

func NewLogger(lvl int) (l logr.Logger) {
	defer func() {
		ctrl.SetLogger(l)
	}()

	return zapr.NewLoggerWithOptions(
		func(opts ...zap.Option) *zap.Logger {
			c := zap.NewProductionConfig()
			if os.Getenv("GOENV") == "DEV" {
				c = zap.NewDevelopmentConfig()
			}
			c.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
			l, _ := c.Build(opts...)
			return l
		}(zap.IncreaseLevel(zap.NewAtomicLevelAt(zapcore.Level(lvl)))),
	)
}
