package main

import (
	"os"

	"github.com/jimmykodes/gommand/examples/valuer/cmd"
)

func main() {
	root := cmd.Root()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
