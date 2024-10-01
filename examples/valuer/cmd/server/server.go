package server

import (
	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/examples/valuer/cmd/server/run"
	"github.com/jimmykodes/gommand/examples/valuer/internal/conf"
)

func Cmd(config *conf.Config) *gommand.Command {
	cmd := &gommand.Command{
		Name: "server",
	}
	cmd.SubCommand(
		run.Cmd(config),
	)
	return cmd
}
