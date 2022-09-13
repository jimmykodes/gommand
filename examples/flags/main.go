package main

import (
	"fmt"
	"os"

	"github.com/jimmykodes/gommand"
)

var rootCmd = &gommand.Command{
	Name: "root",
	Run: func(ctx *gommand.Context) error {
		fmt.Println("root called")
		n := ctx.Flags().Int("num")
		d, err := ctx.Flags().LookupBool("dry-run")
		if err != nil {
			return err
		}
		i, err := ctx.Flags().LookupBool("insensitive")
		if err != nil {
			return err
		}
		fmt.Println(n, d, i)
		return nil
	},
}

var subCmd = &gommand.Command{
	Name: "sub",
	Run: func(ctx *gommand.Context) error {
		m, err := ctx.Flags().LookupInt("mult")
		if err != nil {
			return err
		}
		fmt.Println(m)
		return nil
	},
}

func init() {
	rootCmd.SubCommand(subCmd)

	rootCmd.Flags().Int("num", 10, "a number")
	rootCmd.Flags().BoolS("dry-run", 'd', false, "dry run")
	rootCmd.Flags().BoolS("insensitive", 'i', false, "case-insensitive")
	rootCmd.PersistentFlags().Int("mult", 100, "something")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
