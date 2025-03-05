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
	depth        int

	persistentFlagSets []*flags.FlagSet

	flagGetter *flags.FlagGetter
}

func (c *Context) addPersistentFlags(fs *flags.FlagSet) {
	c.persistentFlagSets = append(c.persistentFlagSets, fs)
}

func (c *Context) persistentFlags() *flags.FlagSet {
	fs := flags.NewFlagSet()
	for _, pfs := range c.persistentFlagSets {
		fs.AddFlagSet(pfs)
	}
	return fs
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) Flags() *flags.FlagGetter {
	return c.flagGetter
}
