package gommand_test

import (
	"errors"
	"testing"

	"github.com/jimmykodes/gommand"
)

func TestUnknownValueRouting(t *testing.T) {
	// A bare value that doesn't match any subcommand fails the default ArgsNone
	// validator and returns an error.
	t.Run("unknown token on command with subcommands returns error", func(t *testing.T) {
		defer overwriteArgs([]string{"typo"})()

		root := &gommand.Command{Name: "root"}
		root.SubCommand(&gommand.Command{
			Name: "real",
			Run:  func(ctx *gommand.Context) error { return nil },
		})

		if err := root.Execute(); err == nil {
			t.Fatal("expected an error for unknown token, got nil")
		}
	})

	// A bare value that doesn't match any subcommand is treated as a positional
	// arg. When the parent has ArgsAny and a Run func, execution proceeds
	// normally with the unknown token in ctx.Args().
	t.Run("unknown token captured as arg when parent has ArgsAny and a Run func", func(t *testing.T) {
		defer overwriteArgs([]string{"typo"})()

		var ranRoot bool
		root := &gommand.Command{
			Name:         "root",
			ArgValidator: gommand.ArgsAny(),
			Run:          func(ctx *gommand.Context) error { ranRoot = true; return nil },
		}
		root.SubCommand(&gommand.Command{
			Name: "real",
			Run:  func(ctx *gommand.Context) error { return nil },
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ranRoot {
			t.Error("expected root Run to be invoked")
		}
	})

	// A registered subcommand name always dispatches even when the parent also
	// has a Run func and ArgsAny.
	t.Run("registered subcommand name takes priority over positional arg", func(t *testing.T) {
		defer overwriteArgs([]string{"real"})()

		var ranRoot, ranSub bool
		root := &gommand.Command{
			Name:         "root",
			ArgValidator: gommand.ArgsAny(),
			Run:          func(ctx *gommand.Context) error { ranRoot = true; return nil },
		}
		root.SubCommand(&gommand.Command{
			Name: "real",
			Run:  func(ctx *gommand.Context) error { ranSub = true; return nil },
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ranSub {
			t.Error("subcommand Run was not invoked")
		}
		if ranRoot {
			t.Error("root Run was invoked instead of the subcommand")
		}
	})

	// A command with subcommands but no Run func receiving an unknown token
	// captures it as an arg and falls through to run(), which returns ErrNoRunner.
	t.Run("unknown token on command with subcommands and no Run returns ErrNoRunner", func(t *testing.T) {
		defer overwriteArgs([]string{"typo"})()

		root := &gommand.Command{
			Name:         "root",
			ArgValidator: gommand.ArgsAny(),
		}
		root.SubCommand(&gommand.Command{
			Name: "real",
			Run:  func(ctx *gommand.Context) error { return nil },
		})

		err := root.Execute()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, gommand.ErrNoRunner) {
			t.Errorf("expected ErrNoRunner, got %T: %v", err, err)
		}
	})

	// A bare arg that appears before a valid subcommand token returns an
	// ordering error. e.g. `math 1 sum 2 3` errors because `1` appears
	// before `sum`.
	t.Run("arg before subcommand token returns ordering error", func(t *testing.T) {
		defer overwriteArgs([]string{"1", "sum", "2", "3"})()

		math := &gommand.Command{Name: "math"}
		math.SubCommand(&gommand.Command{
			Name:         "sum",
			ArgValidator: gommand.ArgsAny(),
			Run:          func(ctx *gommand.Context) error { return nil },
		})

		if err := math.Execute(); err == nil {
			t.Fatal("expected ordering error, got nil")
		}
	})
}
