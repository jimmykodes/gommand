package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

var (
	rootCmd = &gommand.Command{
		Name:         "math",
		Usage:        "a collection of math commands",
		Version:      "1.0.0",
		SilenceError: true,
		PersistentFlagSet: flags.NewFlagSet().AddFlags(
			flags.StringFlag("host", "", "host address"),
			flags.IntFlag("port", 8080, "port number"),
			flags.BoolFlagS("serve", 's', false, "serve something to the host and port"),
		),
	}
	sumCmd = &gommand.Command{
		Name:         "sum n...",
		Usage:        "sum all provided integers",
		ArgValidator: gommand.ArgsMin(1),
		SilenceHelp:  true,
		Run: func(ctx *gommand.Context) error {
			var total int
			for _, s := range ctx.Args() {
				i, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				total += i
			}
			fmt.Println(total)
			return nil
		},
	}
	multCmd = &gommand.Command{
		Name:         "mult n1 n2",
		Usage:        "multiply the two provided integers",
		ArgValidator: gommand.ArgsExact(2),
		Run: func(ctx *gommand.Context) error {
			args := ctx.Args()
			n0, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			n1, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			fmt.Println(n0 * n1)
			return nil
		},
	}
)

func init() {
	rootCmd.SubCommand(sumCmd, multCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
