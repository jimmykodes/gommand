package gommand

import (
	"context"

	"github.com/jimmykodes/gommand/flags"
)

type Context struct {
	context.Context
	cmd          *Command
	args         []string
	preRuns      []func(*Context) error
	postRuns     []func(*Context) error
	deferPost    bool
	silenceHelp  bool
	silenceError bool

	persistentFlags *flags.FlagSet

	flagGetter *flags.FlagGetter
}

func (c *Context) addPersistentFlags(fs *flags.FlagSet) {
	if c.persistentFlags == nil {
		c.persistentFlags = fs
	} else {
		c.persistentFlags.AddFlagSet(fs)
	}
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) Flags() *flags.FlagGetter {
	return c.flagGetter
}
