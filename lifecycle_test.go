package gommand_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/jimmykodes/gommand"
)

func TestLifecycleOrder(t *testing.T) {
	t.Run("single command: PersistentPreRun→PreRun→Run→PostRun→PersistentPostRun", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var order []string
		record := func(s string) func(*gommand.Context) error {
			return func(*gommand.Context) error { order = append(order, s); return nil }
		}

		cmd := &gommand.Command{
			Name:              "cmd",
			PersistentPreRun:  record("PersistentPreRun"),
			PreRun:            record("PreRun"),
			Run:               record("Run"),
			PostRun:           record("PostRun"),
			PersistentPostRun: record("PersistentPostRun"),
		}

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{"PersistentPreRun", "PreRun", "Run", "PostRun", "PersistentPostRun"}
		if !slices.Equal(order, want) {
			t.Errorf("order = %v, want %v", order, want)
		}
	})

	// PersistentPreRun hooks run in FIFO order (root → mid → leaf) before the
	// leaf's PreRun and Run. PersistentPostRun hooks run in LIFO order
	// (leaf → mid → root) after the leaf's PostRun.
	t.Run("deep chain root→mid→leaf: FIFO PersistentPreRun, LIFO PersistentPostRun", func(t *testing.T) {
		defer overwriteArgs([]string{"mid", "leaf"})()

		var order []string
		record := func(s string) func(*gommand.Context) error {
			return func(*gommand.Context) error { order = append(order, s); return nil }
		}

		root := &gommand.Command{
			Name:              "root",
			PersistentPreRun:  record("root:PersistentPreRun"),
			PersistentPostRun: record("root:PersistentPostRun"),
		}
		mid := &gommand.Command{
			Name:              "mid",
			PersistentPreRun:  record("mid:PersistentPreRun"),
			PersistentPostRun: record("mid:PersistentPostRun"),
		}
		leaf := &gommand.Command{
			Name:              "leaf",
			PersistentPreRun:  record("leaf:PersistentPreRun"),
			PreRun:            record("leaf:PreRun"),
			Run:               record("leaf:Run"),
			PostRun:           record("leaf:PostRun"),
			PersistentPostRun: record("leaf:PersistentPostRun"),
		}
		mid.SubCommand(leaf)
		root.SubCommand(mid)

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{
			"root:PersistentPreRun", "mid:PersistentPreRun", "leaf:PersistentPreRun",
			"leaf:PreRun",
			"leaf:Run",
			"leaf:PostRun",
			"leaf:PersistentPostRun", "mid:PersistentPostRun", "root:PersistentPostRun",
		}
		if !slices.Equal(order, want) {
			t.Errorf("order = %v, want %v", order, want)
		}
	})
}

func TestLifecycleDeferPost(t *testing.T) {
	errStop := errors.New("stop")

	// PreRun returning an error halts execution before Run, PostRun, and
	// PersistentPostRun are called.
	t.Run("PreRun error stops Run and post-run functions", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var order []string
		record := func(s string) func(*gommand.Context) error {
			return func(*gommand.Context) error { order = append(order, s); return nil }
		}

		cmd := &gommand.Command{
			Name:              "cmd",
			PreRun:            func(*gommand.Context) error { return errStop },
			Run:               record("Run"),
			PostRun:           record("PostRun"),
			PersistentPostRun: record("PersistentPostRun"),
		}

		if err := cmd.Execute(); !errors.Is(err, errStop) {
			t.Fatalf("expected errStop, got %v", err)
		}
		if len(order) != 0 {
			t.Errorf("no functions should have run after PreRun error, but got: %v", order)
		}
	})

	// A PersistentPreRun in the chain returning an error halts execution before
	// any subsequent PersistentPreRun hooks, PreRun, Run, PostRun, or
	// PersistentPostRun are called.
	t.Run("PersistentPreRun error stops the entire chain", func(t *testing.T) {
		defer overwriteArgs([]string{"mid", "leaf"})()

		var order []string
		record := func(s string) func(*gommand.Context) error {
			return func(*gommand.Context) error { order = append(order, s); return nil }
		}

		root := &gommand.Command{
			Name:              "root",
			PersistentPreRun:  record("root:PersistentPreRun"),
			PersistentPostRun: record("root:PersistentPostRun"),
		}
		mid := &gommand.Command{
			Name:              "mid",
			PersistentPreRun:  func(*gommand.Context) error { return errStop },
			PersistentPostRun: record("mid:PersistentPostRun"),
		}
		leaf := &gommand.Command{
			Name:              "leaf",
			PersistentPreRun:  record("leaf:PersistentPreRun"),
			Run:               record("leaf:Run"),
			PostRun:           record("leaf:PostRun"),
			PersistentPostRun: record("leaf:PersistentPostRun"),
		}
		mid.SubCommand(leaf)
		root.SubCommand(mid)

		if err := root.Execute(); !errors.Is(err, errStop) {
			t.Fatalf("expected errStop, got %v", err)
		}
		// Only root's PersistentPreRun runs before mid's error fires.
		want := []string{"root:PersistentPreRun"}
		if !slices.Equal(order, want) {
			t.Errorf("order = %v, want %v", order, want)
		}
	})

	// When DeferPost is false (the default) and Run returns an error, PostRun
	// and PersistentPostRun are skipped.
	t.Run("Run error stops PostRun when DeferPost=false", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var order []string
		record := func(s string) func(*gommand.Context) error {
			return func(*gommand.Context) error { order = append(order, s); return nil }
		}

		cmd := &gommand.Command{
			Name:              "cmd",
			Run:               func(*gommand.Context) error { return errStop },
			PostRun:           record("PostRun"),
			PersistentPostRun: record("PersistentPostRun"),
		}

		if err := cmd.Execute(); !errors.Is(err, errStop) {
			t.Fatalf("expected errStop, got %v", err)
		}
		if len(order) != 0 {
			t.Errorf("no post-run functions should execute, but got: %v", order)
		}
	})

	// When DeferPost is true and Run returns an error, PostRun and
	// PersistentPostRun still execute.
	t.Run("DeferPost=true: PostRun and PersistentPostRun execute after Run error", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var order []string
		record := func(s string) func(*gommand.Context) error {
			return func(*gommand.Context) error { order = append(order, s); return nil }
		}

		cmd := &gommand.Command{
			Name:              "cmd",
			DeferPost:         true,
			Run:               func(*gommand.Context) error { return errStop },
			PostRun:           record("PostRun"),
			PersistentPostRun: record("PersistentPostRun"),
		}

		if err := cmd.Execute(); !errors.Is(err, errStop) {
			t.Fatalf("expected errStop, got %v", err)
		}
		want := []string{"PostRun", "PersistentPostRun"}
		if !slices.Equal(order, want) {
			t.Errorf("order = %v, want %v", order, want)
		}
	})

	// When DeferPost is true and Run panics, the panic is recovered and returned
	// as an error, and post-run functions still execute.
	t.Run("DeferPost=true: panic recovered as error and PostRun executes", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var ranPostRun bool
		cmd := &gommand.Command{
			Name:      "cmd",
			DeferPost: true,
			Run:       func(*gommand.Context) error { panic("oops") },
			PostRun:   func(*gommand.Context) error { ranPostRun = true; return nil },
		}

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected an error from recovered panic, got nil")
		}
		if !ranPostRun {
			t.Error("PostRun should have been called after panic with DeferPost=true")
		}
	})

	// When DeferPost is false and Run panics, the panic propagates normally
	// and PostRun is not called.
	t.Run("DeferPost=false: panic propagates and PostRun does not execute", func(t *testing.T) {
		defer overwriteArgs([]string{})()

		var ranPostRun bool
		cmd := &gommand.Command{
			Name:    "cmd",
			Run:     func(*gommand.Context) error { panic("oops") },
			PostRun: func(*gommand.Context) error { ranPostRun = true; return nil },
		}

		panicked := false
		func() {
			defer func() {
				if recover() != nil {
					panicked = true
				}
			}()
			_ = cmd.Execute()
		}()

		if !panicked {
			t.Error("expected Execute to panic")
		}
		if ranPostRun {
			t.Error("PostRun should not have been called after panic with DeferPost=false")
		}
	})

	// DeferPost set on a parent propagates to child commands. The child's
	// PostRun executes even when the child's Run returns an error.
	t.Run("DeferPost propagates from parent to child", func(t *testing.T) {
		defer overwriteArgs([]string{"child"})()

		var ranPostRun bool
		root := &gommand.Command{
			Name:      "root",
			DeferPost: true,
		}
		child := &gommand.Command{
			Name:    "child",
			Run:     func(*gommand.Context) error { return errStop },
			PostRun: func(*gommand.Context) error { ranPostRun = true; return nil },
		}
		root.SubCommand(child)

		if err := root.Execute(); !errors.Is(err, errStop) {
			t.Fatalf("expected errStop, got %v", err)
		}
		if !ranPostRun {
			t.Error("PostRun should have been called because DeferPost propagates from root")
		}
	})
}
