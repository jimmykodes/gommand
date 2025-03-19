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

func (c *Context) Args() Args {
	return c.args
}

// Arg will return the command line argument at the given index
// Returns an empty string if idx is out of range
func (c *Context) Arg(idx int) string {
	return Args(c.args).String(idx)
}

func (c *Context) Flags() *flags.FlagGetter {
	return c.flagGetter
}
