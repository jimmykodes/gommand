package gommand

import (
	"context"
)

type Context struct {
	context.Context
	args []string
}

func (c *Context) Args() []string {
	return c.args
}
