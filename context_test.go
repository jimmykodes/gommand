package gommand_test

import (
	"errors"
	"testing"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

func TestContextArgs(t *testing.T) {
	t.Run("ctx.Args() returns all positional arguments in order", func(t *testing.T) {
		defer overwriteArgs([]string{"a", "b", "c"})()

		var got []string
		cmd := &gommand.Command{
			Name:         "cmd",
			ArgValidator: gommand.ArgsAny(),
			Run:          func(ctx *gommand.Context) error { got = ctx.Args(); return nil },
		}

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"a", "b", "c"}
		if len(got) != len(want) {
			t.Fatalf("Args() = %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("Args()[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("ctx.Arg(i) returns the correct argument at a valid index", func(t *testing.T) {
		defer overwriteArgs([]string{"first", "second", "third"})()

		var got [3]string
		cmd := &gommand.Command{
			Name:         "cmd",
			ArgValidator: gommand.ArgsAny(),
			Run: func(ctx *gommand.Context) error {
				got[0] = ctx.Arg(0)
				got[1] = ctx.Arg(1)
				got[2] = ctx.Arg(2)
				return nil
			},
		}

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got[0] != "first" {
			t.Errorf("Arg(0) = %q, want %q", got[0], "first")
		}
		if got[1] != "second" {
			t.Errorf("Arg(1) = %q, want %q", got[1], "second")
		}
		if got[2] != "third" {
			t.Errorf("Arg(2) = %q, want %q", got[2], "third")
		}
	})

	t.Run("ctx.Arg(i) returns empty string for out-of-bounds index", func(t *testing.T) {
		defer overwriteArgs([]string{"only"})()

		var got string
		cmd := &gommand.Command{
			Name:         "cmd",
			ArgValidator: gommand.ArgsAny(),
			Run:          func(ctx *gommand.Context) error { got = ctx.Arg(99); return nil },
		}

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "" {
			t.Errorf("Arg(99) = %q, want empty string", got)
		}
	})
}

func TestContextFlags(t *testing.T) {
	// ctx.Flags() exposes local flags registered on the command.
	t.Run("ctx.Flags() exposes local flags", func(t *testing.T) {
		defer overwriteArgs([]string{"--name", "alice"})()

		fs := flags.NewFlagSet()
		fs.String("name", "", "")

		var got string
		cmd := &gommand.Command{
			Name:    "cmd",
			FlagSet: fs,
			Run: func(ctx *gommand.Context) error {
				v, err := ctx.Flags().LookupString("name")
				if err != nil {
					return err
				}
				got = v
				return nil
			},
		}

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "alice" {
			t.Errorf("Flags().LookupString(\"name\") = %q, want %q", got, "alice")
		}
	})

	// ctx.Flags() on a subcommand exposes persistent flags from the parent.
	t.Run("ctx.Flags() exposes persistent flags from parent", func(t *testing.T) {
		defer overwriteArgs([]string{"--global", "rootval", "sub"})()

		pfs := flags.NewFlagSet()
		pfs.String("global", "", "")

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		root.SubCommand(&gommand.Command{
			Name: "sub",
			Run: func(ctx *gommand.Context) error {
				v, err := ctx.Flags().LookupString("global")
				if err != nil {
					return err
				}
				got = v
				return nil
			},
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "rootval" {
			t.Errorf("Flags().LookupString(\"global\") = %q, want %q", got, "rootval")
		}
	})

	// ctx.Flags() exposes both local and persistent flags through a single FlagGetter.
	t.Run("ctx.Flags() exposes both local and persistent flags", func(t *testing.T) {
		defer overwriteArgs([]string{"--global", "gval", "sub", "--local", "lval"})()

		pfs := flags.NewFlagSet()
		pfs.String("global", "", "")

		lfs := flags.NewFlagSet()
		lfs.String("local", "", "")

		var gotGlobal, gotLocal string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		root.SubCommand(&gommand.Command{
			Name:    "sub",
			FlagSet: lfs,
			Run: func(ctx *gommand.Context) error {
				var err error
				gotGlobal, err = ctx.Flags().LookupString("global")
				if err != nil {
					return err
				}
				gotLocal, err = ctx.Flags().LookupString("local")
				return err
			},
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotGlobal != "gval" {
			t.Errorf("LookupString(\"global\") = %q, want %q", gotGlobal, "gval")
		}
		if gotLocal != "lval" {
			t.Errorf("LookupString(\"local\") = %q, want %q", gotLocal, "lval")
		}
	})

	// ctx.Flags() returns ErrUnregisteredFlag for an unknown flag name.
	t.Run("ctx.Flags() returns ErrUnregisteredFlag for unknown flag", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		cmd := &gommand.Command{
			Name: "cmd",
			Run: func(ctx *gommand.Context) error {
				_, err := ctx.Flags().LookupString("nope")
				return err
			},
		}

		err := cmd.Execute()
		var target flags.ErrUnregisteredFlag
		if !errors.As(err, &target) {
			t.Errorf("expected ErrUnregisteredFlag, got %T: %v", err, err)
		}
	})
}
