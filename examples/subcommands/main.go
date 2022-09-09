package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimmykodes/gommand"
)

var (
	rootCmd = &gommand.Command{
		Name:        "math",
		Description: "a collection of math commands",
	}
	sumCmd = &gommand.Command{
		Name:        "sum",
		Usage:       "sum [...n]",
		Description: "sum all supplied numbers",
		Args:        gommand.ArgsMin(1),
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
		Name:        "mult",
		Usage:       "mult n1 n2",
		Description: "multiply the two provided integers",
		Args:        gommand.ArgsExact(2),
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
		fmt.Println(err)
		os.Exit(1)
	}
}
