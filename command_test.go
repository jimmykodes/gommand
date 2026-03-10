package gommand_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jimmykodes/gommand"
)

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

func overwriteArgs(args []string) func() {
	originalArgs := os.Args
	os.Args = append([]string{"test"}, args...)
	return func() { os.Args = originalArgs }
}
