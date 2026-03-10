package gommand

import (
	"context"
	"io"
	"os"

	"github.com/jimmykodes/gommand/flags"
	"github.com/jimmykodes/gommand/internal/lexer"
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
	lexer        *lexer.Lexer

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

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

// Arg will return the command line argument at the given index
// Returns an empty string if idx is out of range
func (c *Context) Arg(idx int) string {
	if idx < len(c.args) {
		return c.args[idx]
	}
	return ""
}

func (c *Context) Flags() *flags.FlagGetter {
	return c.flagGetter
}

func (c *Context) Stdin() io.Reader {
	if c.stdin != nil {
		return c.stdin
	}
	return os.Stdin
}

func (c *Context) Stdout() io.Writer {
	if c.stdout != nil {
		return c.stdout
	}
	return os.Stdout
}

func (c *Context) Stderr() io.Writer {
	if c.stderr != nil {
		return c.stderr
	}
	return os.Stderr
}
