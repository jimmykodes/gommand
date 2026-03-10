package gommand_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jimmykodes/gommand"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %s", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read pipe: %s", err)
	}
	return strings.TrimRight(string(out), "\n")
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
				return &gommand.Command{
					Name:    "app",
					Version: "1.2.3",
					Run:     noopRun,
				}
			},
			args: []string{"--version"},
			want: "1.2.3",
		},
		{
			name: "no version prints N/A",
			build: func() *gommand.Command {
				return &gommand.Command{
					Name: "app",
					Run:  noopRun,
				}
			},
			args: []string{"--version"},
			want: "N/A",
		},
		{
			name: "subcommand inherits root version",
			build: func() *gommand.Command {
				root := &gommand.Command{
					Name:    "app",
					Version: "2.0.0",
				}
				sub := &gommand.Command{
					Name: "sub",
					Run:  noopRun,
				}
				root.SubCommand(sub)
				return root
			},
			args: []string{"sub", "--version"},
			want: "2.0.0",
		},
		{
			name: "--version does not return an error",
			build: func() *gommand.Command {
				return &gommand.Command{
					Name:    "app",
					Version: "1.0.0",
					Run:     noopRun,
				}
			},
			args:    []string{"--version"},
			want:    "1.0.0",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer overwriteArgs(tt.args)()
			cmd := tt.build()
			var err error
			got := captureStdout(t, func() {
				err = cmd.Execute()
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("version output: got %q, want %q", got, tt.want)
			}
		})
	}
}
