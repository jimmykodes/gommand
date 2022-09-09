package gommand

import (
	"context"
)

type Context struct {
	context.Context
	args      []string
	preRuns   []func(*Context) error
	postRuns  []func(*Context) error
	deferPost bool
}

func (c *Context) Args() []string {
	return c.args
}
