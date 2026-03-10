package gommand_test

import (
	"errors"
	"testing"

	"github.com/jimmykodes/gommand"
)

func TestErrorAccumulationDeferPost(t *testing.T) {
	errRun := errors.New("run error")
	errPost := errors.New("post error")
	errPersistentPost := errors.New("persistent post error")

	// When DeferPost=true and multiple post-run functions fail, all errors are
	// returned joined together — none are silently dropped.
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

	// When DeferPost=true and only post-run functions fail (Run succeeds),
	// those errors are still returned.
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

	// When DeferPost=true across a deep chain, all failing PersistentPostRun
	// hooks from every level are included in the joined error.
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
