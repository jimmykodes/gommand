package gommand

import (
	"context"

	"github.com/jimmykodes/gommand/flags"
)

type Context struct {
	context.Context
	args       []string
	preRuns    []func(*Context) error
	postRuns   []func(*Context) error
	deferPost  bool
	flagGetter *flags.FlagGetter
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) Flags() *flags.FlagGetter {
	return c.flagGetter
}
