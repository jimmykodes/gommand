package main

import (
	"fmt"
	"os"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

var (
	rootCmd = &gommand.Command{
		Name: "root",
		Flags: []flags.Flag{
			flags.IntFlag("num", 10, "a number"),
			flags.BoolFlagS("dry-run", 'd', false, "dry run"),
			flags.BoolFlagS("insensitive", 'i', false, "case-insensitive"),
			flags.StringSliceFlagS("strings", 's', []string{"test", "taco"}, "some strings"),
		},
		ArgValidator: gommand.ArgsAny(),
		PersistentFlags: []flags.Flag{
			flags.IntFlag("mult", 100, "something"),
		},
		Run: func(ctx *gommand.Context) error {
			fmt.Println("args", ctx.Args())
			n := ctx.Flags().Int("num")
			d, err := ctx.Flags().LookupBool("dry-run")
			if err != nil {
				return err
			}
			i, err := ctx.Flags().LookupBool("insensitive")
			if err != nil {
				return err
			}
			if !ctx.Flags().Flag("insensitive").IsSet() {
				fmt.Println("insensitive used default value")
			}
			s, err := ctx.Flags().LookupStringSlice("strings")
			if err != nil {
				return err
			}
			fmt.Println("num", n)
			fmt.Println("dry run", d)
			fmt.Println("insensitive", i)
			fmt.Println("strings", s)
			return nil
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
