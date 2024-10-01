package cmd

import (
	"os"
	"path/filepath"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/examples/valuer/cmd/server"
	"github.com/jimmykodes/gommand/examples/valuer/internal/conf"
	"github.com/jimmykodes/gommand/flags"
)

func Root() *gommand.Command {
	var config conf.Config
	return cmd(&config)
}

func cmd(config *conf.Config) *gommand.Command {
	cmd := &gommand.Command{
		Name: "conf",
		PersistentFlagSet: flags.NewFlagSet().AddFlags(
			flags.StringFlag("config", "", "path to config file. (default $HOME/.config/conf/config.json)"),
		),
		PersistentPreRun: func(ctx *gommand.Context) error {
			configFile := ctx.Flags().String("config")
			if configFile == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				configFile = filepath.Join(home, ".config", "conf", "config.json")
			}
			return config.Load(configFile)
		},
	}
	cmd.SubCommand(
		server.Cmd(config),
	)
	return cmd
}
