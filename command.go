package gommand

import (
	"context"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Usage       string
	Description string
	Args        ArgValidator
	Run         func(*Context) error
	parent      *Command
	commands    map[string]*Command
}

func (c *Command) ExecuteContext(ctx context.Context) error {
	args := os.Args[1:]
	// todo: parse flags

	err := c.execute(&Context{
		Context: ctx,
		args:    args,
	})
	if err != nil {
		// todo: print usage and error
	}
	return err
}

func (c *Command) Execute() error {
	return c.ExecuteContext(context.Background())
}

func (c *Command) SubCommand(cmd *Command, cmds ...*Command) {
	c.subCommand(cmd)
	for _, command := range cmds {
		c.subCommand(command)
	}
}

func (c *Command) subCommand(cmd *Command) {
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}
	c.commands[cmd.Name] = cmd
	cmd.parent = c
}

func (c *Command) execute(ctx *Context) error {
	if c.commands != nil {
		if len(ctx.Args()) > 0 {
			next := c.commands[ctx.args[0]]
			if next != nil {
				ctx.args = ctx.args[1:]
				return next.execute(ctx)
			}
		}
		// cannot have a command with subcommands also have its own run func. because reasons
		return fmt.Errorf("early termination: %s", c.Name)
	}
	validator := c.Args
	if validator == nil {
		// default to allowing no args unless specified otherwise.
		validator = ArgsNone()
	}
	if !validator(ctx.Args()) {
		return fmt.Errorf("invalid args")
	}
	return c.Run(ctx)
}
