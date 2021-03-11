package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/octohelm/cuemod/pkg/cuemod"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type runFn = func(ctx context.Context, args []string) error

func setupRun(cmd *cobra.Command, opts interface{}, fn runFn) *cobra.Command {
	bindFlags(cmd.Flags(), opts)
	cmd.Run = runE(fn)
	return cmd
}

func setupPersistentPreRun(cmd *cobra.Command, opts interface{}, fn runFn) *cobra.Command {
	bindFlags(cmd.PersistentFlags(), opts)
	cmd.PersistentPreRun = runE(fn)
	return cmd
}

func runE(fn runFn) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if projectOpts.Verbose {
			ctx = cuemod.WithOpts(ctx, cuemod.OptVerbose(true))
		}

		if err := fn(ctx, args); err != nil {
			fmt.Println("execute failed: ")
			fmt.Println(err)
		}
	}
}

func bindFlags(flags *pflag.FlagSet, v interface{}) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("only support a ptr struct value, but got %#v", v))
	}

	var walk func(rv reflect.Value)

	walk = func(rv reflect.Value) {
		t := rv.Type()

		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			fv := rv.Field(i)

			if ft.Anonymous && ft.Type.Kind() == reflect.Struct {
				walk(fv)
				continue
			}

			if name, ok := ft.Tag.Lookup("name"); ok {
				parts := strings.Split(name, ",")

				name = parts[0]
				aliass := ""
				usage := ft.Tag.Get("usage")

				if len(parts) > 1 {
					aliass = parts[1]
				}

				switch v := fv.Interface().(type) {
				case string:
					flags.StringVarP(fv.Addr().Interface().(*string), name, aliass, v, usage)
				case []string:
					flags.StringSliceVarP(fv.Addr().Interface().(*[]string), name, aliass, v, usage)
				case bool:
					flags.BoolVarP(fv.Addr().Interface().(*bool), name, aliass, v, usage)
				case []bool:
					flags.BoolSliceVarP(fv.Addr().Interface().(*[]bool), name, aliass, v, usage)
				default:
					panic(fmt.Errorf("unsupported field type %s", ft.Type))
				}
			}

		}
	}

	walk(rv.Elem())
}
