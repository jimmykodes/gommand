package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

type mapSource struct {
	m map[string]string
}

func (m *mapSource) Value(s string) (string, bool) {
	v, ok := m.m[s]
	return v, ok
}

var (
	myMap = mapSource{m: map[string]string{
		"sep": "*",
		"num": "3",
	}}
)

var (
	rootCmd = &gommand.Command{
		Name: "repeat",
		FlagSet: flags.NewFlagSet().AddSource(&myMap).AddFlags(
			flags.StringFlag("string", "", "the string to repeat"),
			flags.IntFlag("num", 10, "the number of times to repeat it"),
			flags.StringFlag("sep", "", "separation string in between repeated elements").AddSources(flags.Environ),
		),
		Run: func(ctx *gommand.Context) error {
			numElems := ctx.Flags().Int("num")
			strs := make([]string, numElems)
			val := ctx.Flags().String("string")

			for i := 0; i < len(strs); i++ {
				strs[i] = val
			}

			fmt.Println(strings.Join(strs, ctx.Flags().String("sep")))

			return nil
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
