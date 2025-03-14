package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimmykodes/gommand"
)

var rootCmd = &gommand.Command{
	Name:         "sum  [...n]",
	Usage:        "sum all provided numbers",
	Version:      "1.0.0",
	ArgValidator: gommand.ArgsEvery(gommand.ArgsMin(1), ints),
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

func ints(s []string) error {
	for _, arg := range s {
		if _, err := strconv.Atoi(arg); err != nil {
			return fmt.Errorf("%s is not an integer", arg)
		}
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
