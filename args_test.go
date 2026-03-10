package gommand_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/jimmykodes/gommand"
)

func Test_ArgValidator(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	ints := func(args []string) error {
		for _, arg := range args {
			_, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
		}
		return nil
	}

	tests := []struct {
		name      string
		args      []string
		validator gommand.ArgValidator
		wanterr   bool
	}{
		{
			name:      "invalid args custom validator",
			args:      []string{"1", "a2", "3"},
			validator: ints,
			wanterr:   true,
		},
		{
			name:      "valid args custom validator",
			args:      []string{"1", "2", "3"},
			validator: ints,
			wanterr:   false,
		},
		{
			name:      "unconfigured validator accepts no args",
			args:      []string{"1", "2", "3"},
			validator: nil,
			wanterr:   true,
		},
		{
			name:      "negative numbers valid as args",
			args:      []string{"1", "-2", "3"},
			validator: gommand.ArgsAny(),
			wanterr:   false,
		},

		{
			name:      "gommand.ArgsExact valid",
			args:      []string{"asdf"},
			validator: gommand.ArgsExact(1),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsExact invalid",
			args:      []string{},
			validator: gommand.ArgsExact(1),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsNone valid",
			args:      []string{},
			validator: gommand.ArgsNone(),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsNone invalid",
			args:      []string{"test"},
			validator: gommand.ArgsNone(),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsMin valid exact",
			args:      []string{"1"},
			validator: gommand.ArgsMin(1),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsMin valid zero",
			args:      []string{},
			validator: gommand.ArgsMin(0),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsMin valid greater than",
			args:      []string{"1", "2"},
			validator: gommand.ArgsMin(1),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsMin invalid",
			args:      []string{},
			validator: gommand.ArgsMin(1),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsMax valid equal",
			args:      []string{"a"},
			validator: gommand.ArgsMax(1),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsMax valid less than",
			args:      []string{"a"},
			validator: gommand.ArgsMax(2),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsMax valid zero",
			args:      []string{},
			validator: gommand.ArgsMax(0),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsMax invalid",
			args:      []string{"a", "b", "c"},
			validator: gommand.ArgsMax(2),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsBetween valid zero lower",
			args:      []string{},
			validator: gommand.ArgsBetween(0, 1),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsBetween valid zero upper",
			args:      []string{},
			validator: gommand.ArgsBetween(0, 0),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsBetween valid in range",
			args:      []string{"a"},
			validator: gommand.ArgsBetween(0, 2),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsBetween valid lower bound",
			args:      []string{"a"},
			validator: gommand.ArgsBetween(1, 3),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsBetween valid upper bound",
			args:      []string{"a", "b", "c"},
			validator: gommand.ArgsBetween(1, 3),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsBetween invalid lower",
			args:      []string{"a"},
			validator: gommand.ArgsBetween(2, 4),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsBetween invalid larger",
			args:      []string{"a", "b", "c", "d"},
			validator: gommand.ArgsBetween(1, 3),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsEvery valid",
			args:      []string{"1", "2"},
			validator: gommand.ArgsEvery(gommand.ArgsMin(2), ints),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsEvery invalid",
			args:      []string{"1", "a"},
			validator: gommand.ArgsEvery(gommand.ArgsMin(2), ints),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsSome valid custom",
			args:      []string{"1"},
			validator: gommand.ArgsSome(gommand.ArgsMin(2), ints),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsSome valid builtin",
			args:      []string{"a", "b"},
			validator: gommand.ArgsSome(gommand.ArgsMin(2), ints),
			wanterr:   false,
		},
		{
			name:      "gommand.ArgsSome invalid",
			args:      []string{"a"},
			validator: gommand.ArgsSome(gommand.ArgsMin(2), ints),
			wanterr:   true,
		},
		{
			name:      "gommand.ArgsAny valid",
			args:      []string{},
			validator: gommand.ArgsAny(),
			wanterr:   false,
		},
	}

	args_match := func(t *testing.T, args []string) func(ctx *gommand.Context) error {
		return func(ctx *gommand.Context) error {
			t.Helper()
			for i, want := range args {
				if got := ctx.Arg(i); got != want {
					t.Fatalf("invalid args at index %d: got %q want %q", i, got, want)
				}
			}
			return nil
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &gommand.Command{
				Name:         "test",
				ArgValidator: tt.validator,
				Run:          args_match(t, tt.args),
			}
			os.Args = append([]string{"testing"}, tt.args...)
			err := cmd.Execute()
			if (err != nil) != tt.wanterr {
				t.Fatal("unexpected error result:", err)
			}
		})
	}
}
