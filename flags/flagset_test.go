package flags_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jimmykodes/gommand/flags"
)

func TestFlagGetter_LookupErrors(t *testing.T) {
	fs := flags.NewFlagSet()
	fs.AddFlags(
		flags.StringFlag("name", "default", "name"),
		flags.IntFlag("count", 0, "count"),
	)
	fg := flags.NewFlagGetter(fs)

	t.Run("unregistered flag returns ErrUnregisteredFlag", func(t *testing.T) {
		_, err := fg.LookupString("unknown")
		var target flags.ErrUnregisteredFlag
		if !errors.As(err, &target) {
			t.Errorf("expected ErrUnregisteredFlag, got %T: %v", err, err)
		}
	})

	t.Run("wrong type returns ErrInvalidFlagType", func(t *testing.T) {
		_, err := fg.LookupInt("name")
		var target flags.ErrInvalidFlagType
		if !errors.As(err, &target) {
			t.Errorf("expected ErrInvalidFlagType, got %T: %v", err, err)
		}
	})
}

func TestFlagSet_MarkRequired(t *testing.T) {
	t.Run("MarkRequired on unregistered flag returns ErrUnregisteredFlag", func(t *testing.T) {
		fs := flags.NewFlagSet()
		err := fs.MarkRequired("nonexistent")
		var target flags.ErrUnregisteredFlag
		if !errors.As(err, &target) {
			t.Errorf("expected ErrUnregisteredFlag, got %T: %v", err, err)
		}
	})

	t.Run("MarkRequired on registered flag succeeds", func(t *testing.T) {
		fs := flags.NewFlagSet()
		fs.AddFlags(flags.StringFlag("name", "", "name"))
		if err := fs.MarkRequired("name"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestFlagTypes(t *testing.T) {
	t.Run("int: valid value", func(t *testing.T) {
		f := flags.IntFlag("count", 0, "count")
		if err := f.Set("42"); err != nil {
			t.Fatalf("Set: %v", err)
		}
		if got := f.Value().(int); got != 42 {
			t.Errorf("got %d, want 42", got)
		}
	})

	t.Run("int: default value when not set", func(t *testing.T) {
		f := flags.IntFlag("count", 7, "count")
		if got := f.Value().(int); got != 7 {
			t.Errorf("got %d, want 7", got)
		}
	})

	t.Run("int: invalid value returns error", func(t *testing.T) {
		f := flags.IntFlag("count", 0, "count")
		if err := f.Set("notanint"); err == nil {
			t.Error("expected error for invalid int, got nil")
		}
	})

	t.Run("string: valid value", func(t *testing.T) {
		f := flags.StringFlag("name", "", "name")
		if err := f.Set("hello"); err != nil {
			t.Fatalf("Set: %v", err)
		}
		if got := f.Value().(string); got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("string: default value when not set", func(t *testing.T) {
		f := flags.StringFlag("name", "fallback", "name")
		if got := f.Value().(string); got != "fallback" {
			t.Errorf("got %q, want %q", got, "fallback")
		}
	})

	t.Run("bool: valid value", func(t *testing.T) {
		f := flags.BoolFlag("verbose", false, "verbose")
		if err := f.Set("true"); err != nil {
			t.Fatalf("Set: %v", err)
		}
		if got := f.Value().(bool); !got {
			t.Errorf("got %v, want true", got)
		}
	})

	t.Run("bool: invalid value returns error", func(t *testing.T) {
		f := flags.BoolFlag("verbose", false, "verbose")
		if err := f.Set("notabool"); err == nil {
			t.Error("expected error for invalid bool, got nil")
		}
	})

	t.Run("float64: valid value", func(t *testing.T) {
		f := flags.Float64Flag("rate", 0, "rate")
		if err := f.Set("3.14"); err != nil {
			t.Fatalf("Set: %v", err)
		}
		if got := f.Value().(float64); got != 3.14 {
			t.Errorf("got %v, want 3.14", got)
		}
	})

	t.Run("float64: invalid value returns error", func(t *testing.T) {
		f := flags.Float64Flag("rate", 0, "rate")
		if err := f.Set("notafloat"); err == nil {
			t.Error("expected error for invalid float64, got nil")
		}
	})

	t.Run("string_slice: comma-separated values", func(t *testing.T) {
		f := flags.StringSliceFlag("tags", nil, "tags")
		if err := f.Set("a,b,c"); err != nil {
			t.Fatalf("Set: %v", err)
		}
		got := f.Value().([]string)
		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("string_slice: default value when not set", func(t *testing.T) {
		def := []string{"x", "y"}
		f := flags.StringSliceFlag("tags", def, "tags")
		got := f.Value().([]string)
		if !reflect.DeepEqual(got, def) {
			t.Errorf("got %v, want %v", got, def)
		}
	})

	t.Run("int_slice: comma-separated values", func(t *testing.T) {
		f := flags.IntSliceFlag("ids", nil, "ids")
		if err := f.Set("1,2,3"); err != nil {
			t.Fatalf("Set: %v", err)
		}
		got := f.Value().([]int)
		want := []int{1, 2, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("int_slice: invalid value returns error", func(t *testing.T) {
		f := flags.IntSliceFlag("ids", nil, "ids")
		if err := f.Set("1,foo,3"); err == nil {
			t.Error("expected error for invalid int in slice, got nil")
		}
	})
}

func TestFlagSource(t *testing.T) {
	t.Run("flag-level source consulted when not set from CLI", func(t *testing.T) {
		f := flags.StringFlag("name", "", "name")
		f.AddSources(flags.ValuerFunc(func(s string) (string, bool) {
			if s == "name" {
				return "from-source", true
			}
			return "", false
		}))
		if err := flags.SetFromSources(f); err != nil {
			t.Fatalf("SetFromSources: %v", err)
		}
		if got := f.Value().(string); got != "from-source" {
			t.Errorf("got %q, want %q", got, "from-source")
		}
	})

	t.Run("CLI value takes precedence over source", func(t *testing.T) {
		fs := flags.NewFlagSet()
		f := flags.StringFlag("name", "", "name")
		f.AddSources(flags.ValuerFunc(func(s string) (string, bool) {
			return "from-source", true
		}))
		fs.AddFlags(f)

		// simulate CLI setting the flag
		if err := fs.FromName("name").Set("from-cli"); err != nil {
			t.Fatalf("Set: %v", err)
		}

		fg := flags.NewFlagGetter(fs)
		got, err := fg.LookupString("name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "from-cli" {
			t.Errorf("got %q, want %q", got, "from-cli")
		}
	})
}
