package gommand_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

// overwriteArgs replaces os.Args for the duration of a test.
// Usage: defer overwriteArgs([]string{"arg1", "arg2"})()
func overwriteArgs(args []string) func() {
	originalArgs := os.Args
	os.Args = append([]string{"test"}, args...)
	return func() { os.Args = originalArgs }
}

func TestSimpleCommand(t *testing.T) {
	defer overwriteArgs([]string{"1", "2", "3"})()

	cmd := &gommand.Command{
		Name:         "sum  [...n]",
		ArgValidator: gommand.ArgsAny(),
		Run: func(ctx *gommand.Context) error {
			var total int
			for _, s := range ctx.Args() {
				i, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				total += i
			}
			if total != 6 {
				t.Fatalf("unexpected sum: got %v - want 6", total)
			}
			return nil
		},
	}

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestContextCancel(t *testing.T) {
	defer overwriteArgs([]string{})()
	cmd := &gommand.Command{
		Name: "test",
		Run: func(ctx *gommand.Context) error {
			tick := time.NewTicker(time.Second)
			select {
			case <-tick.C:
				return fmt.Errorf("ticker tripped instead of context")
			case <-ctx.Done():
				return nil
			}
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)
	go func() { errChan <- cmd.ExecuteContext(ctx) }()
	cancel()
	err := <-errChan
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

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

func TestUnknownValueRouting(t *testing.T) {
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

func TestSilenceHelp(t *testing.T) {
	errRun := errors.New("run error")

	t.Run("SilenceHelp=true suppresses help on error", func(t *testing.T) {
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

	t.Run("SilenceHelp propagates to subcommands", func(t *testing.T) {
		defer overwriteArgs([]string{"sub"})()

		var stderr bytes.Buffer
		root := &gommand.Command{
			Name:         "root",
			SilenceHelp:  true,
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

func TestVersion(t *testing.T) {
	noopRun := func(ctx *gommand.Context) error { return nil }

	tests := []struct {
		name    string
		build   func() *gommand.Command
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "root version printed",
			build: func() *gommand.Command {
				return &gommand.Command{Name: "app", Version: "1.2.3", Run: noopRun}
			},
			args: []string{"--version"},
			want: "1.2.3",
		},
		{
			name: "no version prints N/A",
			build: func() *gommand.Command {
				return &gommand.Command{Name: "app", Run: noopRun}
			},
			args: []string{"--version"},
			want: "N/A",
		},
		{
			name: "subcommand inherits root version",
			build: func() *gommand.Command {
				root := &gommand.Command{Name: "app", Version: "2.0.0"}
				sub := &gommand.Command{Name: "sub", Run: noopRun}
				root.SubCommand(sub)
				return root
			},
			args: []string{"sub", "--version"},
			want: "2.0.0",
		},
		{
			name: "--version does not return an error",
			build: func() *gommand.Command {
				return &gommand.Command{Name: "app", Version: "1.0.0", Run: noopRun}
			},
			args:    []string{"--version"},
			want:    "1.0.0",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer overwriteArgs(tt.args)()
			var buf bytes.Buffer
			cmd := tt.build()
			err := cmd.Execute(gommand.WithStdout(&buf))
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			got := strings.TrimRight(buf.String(), "\n")
			if got != tt.want {
				t.Fatalf("version output: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorAccumulationDeferPost(t *testing.T) {
	errRun := errors.New("run error")
	errPost := errors.New("post error")
	errPersistentPost := errors.New("persistent post error")

	t.Run("all post-run errors are joined when DeferPost=true", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		cmd := &gommand.Command{
			Name:              "cmd",
			DeferPost:         true,
			Run:               func(*gommand.Context) error { return errRun },
			PostRun:           func(*gommand.Context) error { return errPost },
			PersistentPostRun: func(*gommand.Context) error { return errPersistentPost },
		}

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errRun) {
			t.Errorf("expected joined error to contain errRun, got: %v", err)
		}
		if !errors.Is(err, errPost) {
			t.Errorf("expected joined error to contain errPost, got: %v", err)
		}
		if !errors.Is(err, errPersistentPost) {
			t.Errorf("expected joined error to contain errPersistentPost, got: %v", err)
		}
	})

	t.Run("post-run errors returned even when Run succeeds", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		cmd := &gommand.Command{
			Name:              "cmd",
			DeferPost:         true,
			Run:               func(*gommand.Context) error { return nil },
			PostRun:           func(*gommand.Context) error { return errPost },
			PersistentPostRun: func(*gommand.Context) error { return errPersistentPost },
		}

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected error from post-run failures, got nil")
		}
		if !errors.Is(err, errPost) {
			t.Errorf("expected joined error to contain errPost, got: %v", err)
		}
		if !errors.Is(err, errPersistentPost) {
			t.Errorf("expected joined error to contain errPersistentPost, got: %v", err)
		}
	})

	t.Run("PersistentPostRun errors from multiple levels are all joined", func(t *testing.T) {
		defer overwriteArgs([]string{"child"})()

		errRoot := errors.New("root post error")
		errChild := errors.New("child post error")

		root := &gommand.Command{
			Name:              "root",
			DeferPost:         true,
			PersistentPostRun: func(*gommand.Context) error { return errRoot },
		}
		root.SubCommand(&gommand.Command{
			Name:              "child",
			Run:               func(*gommand.Context) error { return nil },
			PersistentPostRun: func(*gommand.Context) error { return errChild },
		})

		err := root.Execute()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errRoot) {
			t.Errorf("expected joined error to contain errRoot, got: %v", err)
		}
		if !errors.Is(err, errChild) {
			t.Errorf("expected joined error to contain errChild, got: %v", err)
		}
	})
}
