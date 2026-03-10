package gommand_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jimmykodes/gommand"
)

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
