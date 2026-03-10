package gommand_test

import (
	"errors"
	"testing"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

func TestPersistentFlags(t *testing.T) {
	t.Run("persistent flag accessible in subcommand", func(t *testing.T) {
		defer overwriteArgs([]string{"sub"})()

		pfs := flags.NewFlagSet().
			AddFlag(flags.StringFlag("global", "default", "a global flag"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name: "sub",
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("global")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "default" {
			t.Errorf("got %q, want %q", got, "default")
		}
	})

	t.Run("persistent flag set before subcommand token is available in subcommand", func(t *testing.T) {
		defer overwriteArgs([]string{"--global", "hello", "sub"})()

		pfs := flags.NewFlagSet().
			AddFlag(flags.StringFlag("global", "", "a global flag"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name: "sub",
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("global")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("persistent flag set after subcommand token is available in subcommand", func(t *testing.T) {
		defer overwriteArgs([]string{"sub", "--global", "world"})()

		pfs := flags.NewFlagSet().
			AddFlag(flags.StringFlag("global", "", "a global flag"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name: "sub",
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("global")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "world" {
			t.Errorf("got %q, want %q", got, "world")
		}
	})

	t.Run("required persistent flag missing returns error", func(t *testing.T) {
		defer overwriteArgs([]string{"sub"})()

		pfs := flags.NewFlagSet().
			AddFlag(flags.StringFlag("token", "", "required token").Required())

		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name: "sub",
			Run:  func(ctx *gommand.Context) error { return nil },
		}
		root.SubCommand(sub)

		err := root.Execute()
		if err == nil {
			t.Fatal("expected error for missing required persistent flag, got nil")
		}
		var mrf flags.ErrMissingRequiredFlag
		if !errors.As(err, &mrf) {
			t.Errorf("expected ErrMissingRequiredFlag, got %T: %v", err, err)
		}
	})

	t.Run("required persistent flag satisfied via CLI", func(t *testing.T) {
		defer overwriteArgs([]string{"--token", "abc123", "sub"})()

		pfs := flags.NewFlagSet().
			AddFlag(flags.StringFlag("token", "", "required token").Required())

		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name: "sub",
			Run:  func(ctx *gommand.Context) error { return nil },
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestPersistentFlagCollisions(t *testing.T) {
	// When a subcommand defines a local flag with the same name as a parent's
	// persistent flag, AddFlagSet adds persistent flags first then local flags,
	// so local flags win (maps.Copy overwrites). The tests below document this
	// precedence behaviour.

	t.Run("local flag shadows persistent flag in ctx.Flags()", func(t *testing.T) {
		// Neither flag is set via CLI; verify the local default wins.
		defer overwriteArgs([]string{"sub"})()

		pfs := flags.NewFlagSet().
			AddFlag(flags.StringFlag("output", "persistent-default", "output"))

		localFS := flags.NewFlagSet().
			AddFlag(flags.StringFlag("output", "local-default", "output"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name:    "sub",
			FlagSet: localFS,
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("output")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "local-default" {
			t.Errorf("got %q, want %q (local flag should shadow persistent)", got, "local-default")
		}
	})

	t.Run("local flag value wins when set after subcommand token", func(t *testing.T) {
		// --output is set after the subcommand token; the combined fs resolves
		// --output to the local flag, so the local flag gets the value.
		defer overwriteArgs([]string{"sub", "--output", "local-value"})()

		pfs := flags.NewFlagSet()
		pfs.AddFlags(flags.StringFlag("output", "persistent-default", "output"))

		localFS := flags.NewFlagSet()
		localFS.AddFlags(flags.StringFlag("output", "local-default", "output"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name:    "sub",
			FlagSet: localFS,
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("output")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "local-value" {
			t.Errorf("got %q, want %q", got, "local-value")
		}
	})

	t.Run("persistent value set before subcommand token is not visible via shadowed name", func(t *testing.T) {
		// --output is set before the subcommand token; root's execute sets the
		// persistent flag object. In sub's execute the combined fs maps "output"
		// to the local flag object (which is unset), so ctx.Flags() returns the
		// local default, not the persistent value.
		defer overwriteArgs([]string{"--output", "persistent-value", "sub"})()

		pfs := flags.NewFlagSet()
		pfs.AddFlags(flags.StringFlag("output", "persistent-default", "output"))

		localFS := flags.NewFlagSet()
		localFS.AddFlags(flags.StringFlag("output", "local-default", "output"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name:    "sub",
			FlagSet: localFS,
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("output")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "local-default" {
			t.Errorf("got %q, want %q (local flag shadows persistent value set before subcommand token)", got, "local-default")
		}
	})

	t.Run("short flag collision: local short flag shadows persistent short flag", func(t *testing.T) {
		defer overwriteArgs([]string{"sub", "-o", "local-value"})()

		pfs := flags.NewFlagSet()
		pfs.AddFlags(flags.StringFlagS("output", 'o', "persistent-default", "output"))

		localFS := flags.NewFlagSet()
		localFS.AddFlags(flags.StringFlagS("output", 'o', "local-default", "output"))

		var got string
		root := &gommand.Command{
			Name:              "root",
			PersistentFlagSet: pfs,
		}
		sub := &gommand.Command{
			Name:    "sub",
			FlagSet: localFS,
			Run: func(ctx *gommand.Context) error {
				got = ctx.Flags().String("output")
				return nil
			},
		}
		root.SubCommand(sub)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "local-value" {
			t.Errorf("got %q, want %q", got, "local-value")
		}
	})
}
