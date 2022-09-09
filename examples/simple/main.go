package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimmykodes/gommand"
)

var rootCmd = &gommand.Command{
	Name:        "sum",
	Usage:       "sum [...n]",
	Description: "sum all provided numbers",
	Args:        gommand.ArgsEvery(gommand.ArgsMin(1), ints),
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

func ints(s []string) bool {
	for _, arg := range s {
		if _, err := strconv.Atoi(arg); err != nil {
			return false
		}
	}
	return true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
