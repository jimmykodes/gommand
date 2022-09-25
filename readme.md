# Gommand

A simple library for making CLI tools.

Acknowledging my elders: this tool is heavily inspired by Cobra/Viper, with just a touch of urfave/cli. 

---

## YAGCLL

Maybe a more fitting name would be: Yet Another Go Command Line Library. With libraries like Cobra/Viper and urfave/cli, why do we need Yet Another One? 
The decision to write my own library was made for 2 reasons:

- First, and most importantly: I just wanted to.
- Second: All the clis I have made using Cobra/Viper require the same boiler plate to configure flags from env vars. Rather than continuing to require the extra steps,
I just made these behaviors the default.
- Third (because who doesn't love an off-by-one error): Viper has a large dependency graph because it does a lot of things I never need, namely: flag values from config files.

To summarize: I wanted a smaller dependency footprint that required less boiler plate for my base use case. 

## Usage

For more usage examples, check the examples directory.

```golang
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimmykodes/gommand"
)

var rootCmd = &gommand.Command{
	Name:         "sum",
	Usage:        "sum [...n]",
	Description:  "sum all provided numbers",
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

```
