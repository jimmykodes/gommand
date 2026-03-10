package gommand_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/jimmykodes/gommand"
)

func TestSilenceHelp(t *testing.T) {
	errRun := errors.New("run error")

	// When SilenceHelp=true, an error from Run does not print help to stderr.
	t.Run("SilenceHelp=true suppresses help on error", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var stderr bytes.Buffer
		cmd := &gommand.Command{
			Name:        "app",
			SilenceHelp: true,
			SilenceError: true,
			Run:         func(*gommand.Context) error { return errRun },
		}

		_ = cmd.Execute(gommand.WithStderr(&stderr))
		if stderr.Len() != 0 {
			t.Errorf("expected empty stderr, got: %q", stderr.String())
		}
	})

	// Without SilenceHelp, an error from Run prints help to stderr.
	t.Run("SilenceHelp=false prints help on error", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var stderr bytes.Buffer
		cmd := &gommand.Command{
			Name:         "app",
			SilenceError: true,
			Run:          func(*gommand.Context) error { return errRun },
		}

		_ = cmd.Execute(gommand.WithStderr(&stderr))
		if stderr.Len() == 0 {
			t.Error("expected help text in stderr, got nothing")
		}
	})

	// SilenceHelp set on the root propagates to subcommands even when the
	// subcommand does not set it explicitly.
	t.Run("SilenceHelp propagates to subcommands", func(t *testing.T) {
		defer overwriteArgs([]string{"sub"})()

		var stderr bytes.Buffer
		root := &gommand.Command{
			Name:        "root",
			SilenceHelp: true,
			SilenceError: true,
		}
		root.SubCommand(&gommand.Command{
			Name: "sub",
			Run:  func(*gommand.Context) error { return errRun },
		})

		_ = root.Execute(gommand.WithStderr(&stderr))
		if stderr.Len() != 0 {
			t.Errorf("expected empty stderr due to propagated SilenceHelp, got: %q", stderr.String())
		}
	})
}

func TestSilenceError(t *testing.T) {
	errRun := errors.New("run error")

	// When SilenceError=true, the "Error: ..." line is not printed to stderr.
	t.Run("SilenceError=true suppresses error line", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var stderr bytes.Buffer
		cmd := &gommand.Command{
			Name:         "app",
			SilenceHelp:  true,
			SilenceError: true,
			Run:          func(*gommand.Context) error { return errRun },
		}

		_ = cmd.Execute(gommand.WithStderr(&stderr))
		if stderr.Len() != 0 {
			t.Errorf("expected empty stderr, got: %q", stderr.String())
		}
	})

	// Without SilenceError, the "Error: ..." line is printed to stderr.
	t.Run("SilenceError=false prints error line", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var stderr bytes.Buffer
		cmd := &gommand.Command{
			Name:        "app",
			SilenceHelp: true,
			Run:         func(*gommand.Context) error { return errRun },
		}

		_ = cmd.Execute(gommand.WithStderr(&stderr))
		if !strings.Contains(stderr.String(), "Error:") {
			t.Errorf("expected \"Error:\" in stderr, got: %q", stderr.String())
		}
	})

	// The two flags are independently controllable: SilenceError=true with
	// SilenceHelp=false still prints help but suppresses the error line.
	t.Run("SilenceError=true SilenceHelp=false: help printed, error line suppressed", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var stderr bytes.Buffer
		cmd := &gommand.Command{
			Name:         "app",
			SilenceError: true,
			Run:          func(*gommand.Context) error { return errRun },
		}

		_ = cmd.Execute(gommand.WithStderr(&stderr))
		out := stderr.String()
		if strings.Contains(out, "Error:") {
			t.Errorf("expected no \"Error:\" line, got: %q", out)
		}
		if len(out) == 0 {
			t.Error("expected help text in stderr, got nothing")
		}
	})
}
