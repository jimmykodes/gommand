package run

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/examples/valuer/internal/conf"
	"github.com/jimmykodes/gommand/flags"
)

func Cmd(config *conf.Config) *gommand.Command {
	cmd := &gommand.Command{
		Name: "run",
		FlagSet: flags.NewFlagSet().
			AddSource(flags.ValuerFunc(config.SubPath("server"))).
			AddFlags(
				flags.StringFlag("addr", ":8080", "server address"),
			),
		Run: func(ctx *gommand.Context) error {
			mux := http.NewServeMux()
			svr := http.Server{
				Addr:    ctx.Flags().String("addr"),
				Handler: mux,
			}
			mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte{'o', 'k', '\n'}) })

			slog.InfoContext(ctx, "server running", "addr", svr.Addr)
			if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
	}
	return cmd
}
