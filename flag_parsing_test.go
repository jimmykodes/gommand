package gommand_test

import (
	"errors"
	"testing"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

func TestFlagParsing(t *testing.T) {
	type testCase struct {
		name    string
		args    []string
		flags   []flags.Flag
		check   func(*testing.T, *flags.FlagGetter)
		wantErr bool
	}

	tests := []testCase{
		{
			name:  "long flag space-separated",
			args:  []string{"--name", "alice"},
			flags: []flags.Flag{flags.StringFlag("name", "", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.String("name"); got != "alice" {
					t.Errorf("got %q, want %q", got, "alice")
				}
			},
		},
		{
			name:  "long flag = assignment",
			args:  []string{"--name=alice"},
			flags: []flags.Flag{flags.StringFlag("name", "", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.String("name"); got != "alice" {
					t.Errorf("got %q, want %q", got, "alice")
				}
			},
		},
		{
			name:  "short flag space-separated",
			args:  []string{"-n", "alice"},
			flags: []flags.Flag{flags.StringFlagS("name", 'n', "", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.String("name"); got != "alice" {
					t.Errorf("got %q, want %q", got, "alice")
				}
			},
		},
		{
			name:  "short flag = assignment",
			args:  []string{"-n=alice"},
			flags: []flags.Flag{flags.StringFlagS("name", 'n', "", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.String("name"); got != "alice" {
					t.Errorf("got %q, want %q", got, "alice")
				}
			},
		},
		{
			name:  "bool flag presence sets true",
			args:  []string{"--verbose"},
			flags: []flags.Flag{flags.BoolFlag("verbose", false, "verbose output")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.Bool("verbose"); !got {
					t.Errorf("got %v, want true", got)
				}
			},
		},
		{
			name:  "bool flag explicit false",
			args:  []string{"--verbose=false"},
			flags: []flags.Flag{flags.BoolFlag("verbose", true, "verbose output")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.Bool("verbose"); got {
					t.Errorf("got %v, want false", got)
				}
			},
		},
		{
			name: "multi-flag sets all bool flags",
			args: []string{"-abc"},
			flags: []flags.Flag{
				flags.BoolFlagS("alpha", 'a', false, "alpha"),
				flags.BoolFlagS("bravo", 'b', false, "bravo"),
				flags.BoolFlagS("charlie", 'c', false, "charlie"),
			},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				for _, name := range []string{"alpha", "bravo", "charlie"} {
					if got := fg.Bool(name); !got {
						t.Errorf("flag %q: got %v, want true", name, got)
					}
				}
			},
		},
		{
			name: "multi-flag with non-bool flag is error",
			args: []string{"-as"},
			flags: []flags.Flag{
				flags.BoolFlagS("alpha", 'a', false, "alpha"),
				flags.StringFlagS("str", 's', "", "str"),
			},
			wantErr: true,
		},
		{
			name: "multi-flag with = value is error",
			args: []string{"-ab=val"},
			flags: []flags.Flag{
				flags.BoolFlagS("alpha", 'a', false, "alpha"),
				flags.BoolFlagS("bravo", 'b', false, "bravo"),
			},
			wantErr: true,
		},
		{
			name:    "unknown long flag is error",
			args:    []string{"--undefined"},
			wantErr: true,
		},
		{
			name:    "unknown short flag is error",
			args:    []string{"-z"},
			wantErr: true,
		},
		{
			name:  "default value used when flag not provided",
			args:  []string{},
			flags: []flags.Flag{flags.StringFlag("name", "default", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				if got := fg.String("name"); got != "default" {
					t.Errorf("got %q, want %q", got, "default")
				}
			},
		},
		{
			name:  "IsSet false when using default",
			args:  []string{},
			flags: []flags.Flag{flags.StringFlag("name", "default", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				f := fg.Flag("name")
				if f == nil {
					t.Fatal("flag not found")
				}
				if f.IsSet() {
					t.Error("IsSet() should be false when using default")
				}
			},
		},
		{
			name:  "IsSet true after CLI provides value",
			args:  []string{"--name", "provided"},
			flags: []flags.Flag{flags.StringFlag("name", "default", "name")},
			check: func(t *testing.T, fg *flags.FlagGetter) {
				f := fg.Flag("name")
				if f == nil {
					t.Fatal("flag not found")
				}
				if !f.IsSet() {
					t.Error("IsSet() should be true after CLI provides value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer overwriteArgs(tt.args)()

			fs := flags.NewFlagSet().
				AddFlags(tt.flags...)

			var capturedGetter *flags.FlagGetter
			cmd := &gommand.Command{
				Name:         "test",
				FlagSet:      fs,
				SilenceHelp:  true,
				SilenceError: true,
				ArgValidator: gommand.ArgsAny(),
				Run: func(ctx *gommand.Context) error {
					capturedGetter = ctx.Flags()
					return nil
				},
			}

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, capturedGetter)
			}
		})
	}
}

func TestRequiredFlags(t *testing.T) {
	t.Run("required flag missing returns ErrMissingRequiredFlag", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		fs := flags.NewFlagSet().
			AddFlags(flags.StringFlag("output", "", "output path").Required())
		cmd := &gommand.Command{
			Name:         "test",
			FlagSet:      fs,
			SilenceHelp:  true,
			SilenceError: true,
			Run:          func(ctx *gommand.Context) error { return nil },
		}

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected error for missing required flag, got nil")
		}
		var mrf flags.ErrMissingRequiredFlag
		if !errors.As(err, &mrf) {
			t.Errorf("expected ErrMissingRequiredFlag, got %T: %v", err, err)
		}
	})

	t.Run("required flag satisfied via CLI", func(t *testing.T) {
		defer overwriteArgs([]string{"--output", "out.txt"})()

		fs := flags.NewFlagSet().
			AddFlags(flags.StringFlag("output", "", "output path").Required())
		cmd := &gommand.Command{
			Name:         "test",
			FlagSet:      fs,
			SilenceHelp:  true,
			SilenceError: true,
			Run:          func(ctx *gommand.Context) error { return nil },
		}
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("MarkRequired on unregistered flag returns ErrUnregisteredFlag", func(t *testing.T) {
		fs := flags.NewFlagSet()
		err := fs.MarkRequired("nonexistent")
		var target flags.ErrUnregisteredFlag
		if !errors.As(err, &target) {
			t.Errorf("expected ErrUnregisteredFlag, got %T: %v", err, err)
		}
	})
}
