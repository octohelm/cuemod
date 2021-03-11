package cuemod

import "context"

type Opts struct {
	Upgrade bool `name:"upgrade,u" usage:"need upgrade dependencies"`
	Verbose bool `name:"verbose,v" usage:"verbose"`
}

type OptFn = func(o *Opts)

func OptUpgrade(b bool) OptFn {
	return func(o *Opts) {
		o.Upgrade = b
	}
}

func OptVerbose(v bool) OptFn {
	return func(o *Opts) {
		o.Verbose = v
	}
}

type contextKeyForOpts int

func OptsFromContext(ctx context.Context) Opts {
	if o, ok := ctx.Value(contextKeyForOpts(0)).(Opts); ok {
		return o
	}
	return Opts{}
}

func WithOpts(ctx context.Context, fns ...OptFn) context.Context {
	o := OptsFromContext(ctx)

	for _, fn := range fns {
		if fn != nil {
			fn(&o)
		}
	}

	return context.WithValue(ctx, contextKeyForOpts(0), o)
}
