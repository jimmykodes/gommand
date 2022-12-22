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
		PersistentFlags: []flags.Flag{
			flags.IntFlag("mult", 100, "something"),
		},
		Run: func(ctx *gommand.Context) error {
			fmt.Println("root called")
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
			fmt.Println(n, d, i, s)
			return nil
		},
	}
	subCommands = []*gommand.Command{
		{
			Name: "sub1",
			Flags: []flags.Flag{
				flags.StringFlag("context", "", "some important context for you to understand"),
			},
			PersistentFlags: []flags.Flag{
				flags.IntFlag("port", 12, "a port of some kind"),
			},
			Run: func(context *gommand.Context) error {
				fmt.Println("sub1 called")
				return nil
			},
		},
		{
			Name: "sub2",
			Flags: []flags.Flag{
				flags.BoolFlag("treat", false, "true if you get a treat, false if you get a trick"),
			},
			Run: func(context *gommand.Context) error {
				fmt.Println("sub1 called")
				return nil
			},
		},
	}
)

func init() {
	rootCmd.SubCommand(subCommands...)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
