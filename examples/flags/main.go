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
		if !ctx.Flags().Flag("insensitive").IsSet() {
			fmt.Println("insensitive used default value")
		}
		s, err := ctx.Flags().LookupStringSlice("strings")
		if err != nil {
			return err
		}
		fmt.Println(n, d, i, s)
		return nil
	},
}

func init() {
	rootCmd.Flags().Int("num", 10, "a number")
	rootCmd.Flags().BoolS("dry-run", 'd', false, "dry run")
	rootCmd.Flags().BoolS("insensitive", 'i', false, "case-insensitive")
	rootCmd.Flags().StringSliceS("strings", 's', []string{"test", "taco"}, "some strings")
	rootCmd.PersistentFlags().Int("mult", 100, "something")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
