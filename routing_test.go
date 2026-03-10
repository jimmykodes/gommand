package gommand_test

import (
	"errors"
	"testing"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

func TestSubcommandRouting(t *testing.T) {
	t.Run("basic dispatch", func(t *testing.T) {
		defer overwriteArgs([]string{"sub"})()

		var ran bool
		root := &gommand.Command{Name: "root"}
		sub := &gommand.Command{
			Name: "sub",
			Run:  func(ctx *gommand.Context) error { ran = true; return nil },
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ran {
			t.Error("subcommand Run was not called")
		}
	})

	t.Run("deep chain root→mid→leaf", func(t *testing.T) {
		defer overwriteArgs([]string{"mid", "leaf"})()

		var order []string
		root := &gommand.Command{Name: "root"}
		mid := &gommand.Command{
			Name:             "mid",
			PersistentPreRun: func(ctx *gommand.Context) error { order = append(order, "mid"); return nil },
		}
		leaf := &gommand.Command{
			Name: "leaf",
			Run:  func(ctx *gommand.Context) error { order = append(order, "leaf"); return nil },
		}
		mid.SubCommand(leaf)
		root.SubCommand(mid)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order) != 2 || order[0] != "mid" || order[1] != "leaf" {
			t.Errorf("unexpected order: %v", order)
		}
	})

	t.Run("alias dispatch", func(t *testing.T) {
		for _, alias := range []string{"ls", "l"} {
			alias := alias
			t.Run(alias, func(t *testing.T) {
				defer overwriteArgs([]string{alias})()

				var ran bool
				root := &gommand.Command{Name: "root"}
				sub := &gommand.Command{
					Name:    "list",
					Aliases: []string{"ls", "l"},
					Run:     func(ctx *gommand.Context) error { ran = true; return nil },
				}
				root.SubCommand(sub)

				if err := root.Execute(); err != nil {
					t.Fatalf("unexpected error via alias %q: %v", alias, err)
				}
				if !ran {
					t.Errorf("Run not called via alias %q", alias)
				}
			})
		}
	})

	t.Run("ErrNoRunner when no Run and no matching subcommand", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		root := &gommand.Command{Name: "root", SilenceHelp: true, SilenceError: true}

		err := root.Execute()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, gommand.ErrNoRunner) {
			t.Errorf("expected ErrNoRunner, got %v", err)
		}
	})

	t.Run("mixed token order: flag before subcommand, flag and arg after", func(t *testing.T) {
		defer overwriteArgs([]string{"--global", "gval", "sub", "--local", "lval", "myarg"})()

		pfs := flags.NewFlagSet().AddFlag(flags.StringFlag("global", "", "global flag"))
		lfs := flags.NewFlagSet().AddFlag(flags.StringFlag("local", "", "local flag"))

		var gotGlobal, gotLocal, gotArg string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name:         "sub",
			FlagSet:      lfs,
			ArgValidator: gommand.ArgsExact(1),
			Run: func(ctx *gommand.Context) error {
				gotGlobal = ctx.Flags().String("global")
				gotLocal = ctx.Flags().String("local")
				gotArg = ctx.Arg(0)
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotGlobal != "gval" {
			t.Errorf("global: got %q, want %q", gotGlobal, "gval")
		}
		if gotLocal != "lval" {
			t.Errorf("local: got %q, want %q", gotLocal, "lval")
		}
		if gotArg != "myarg" {
			t.Errorf("arg: got %q, want %q", gotArg, "myarg")
		}
	})
}
