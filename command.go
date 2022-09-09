package gommand

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/multierr"
)

// Command represents a command line command
//
// The order of functions is:
//   - PersistentPreRun -- see note on field
//   - PreRun
//   - Run
//   - PostRun
//   - PersistentPostRun -- see note of field
//
// if PersistentPreRun or PreRun return an error, execution is stopped and the error is returned
// if Run returns an error and DeferPost is false, execution is stopped and the error is returned
// if DeferPost is true, PostRun and PersistentPostRun will be executed even if Run returns an error
type Command struct {
	// Name is the name of the command.
	// This is ignored if the command is invoked using Execute or ExecuteContext, but if registered as a subcommand
	// of a function, the name defines how the subcommand is called.
	// ex:
	// c1 := &Command{Name: "entrypoint"}
	// c2 := &Command{Name: "my-sub-command", Run: func(*Context) error { fmt.Println("sub"); return nil }}
	//
	// c1.SubCommand(c2)
	//
	// func main() {
	// 		_ = c1.Execute()
	// }
	//
	// c2 would be called by running
	// entrypoint my-sub-command
	//
	// if spaces are included in the name, anything after the first space is discarded
	Name string

	// Usage is the short usage explanation string
	// todo: add more detail here
	Usage string

	// Description is the longer description of the command printed out by the help text
	Description string

	// Args is an ArgValidator to be called on the args of the function being executed. This is called before any of
	// the functions for this command are called.
	// If this is not defined ArgsNone is used.
	Args ArgValidator

	// Run is the core function the command should execute
	Run func(*Context) error

	// PreRun will run immediately before Run, if defined
	PreRun func(*Context) error
	// PostRun will run immediately after Run if defined and either Run exits with no error or DeferPost is true
	PostRun func(*Context) error

	// PersistentPreRun is a function that will run before PreRun and will be run for any subcommands of this command.
	//
	// PersistentPreRun commands are executed in FIFO order
	//
	// ex:
	// c1 := &Command{Name: "c1", PersistentPreRun: func(*Context) error { fmt.Println("c1"); return nil }}
	// c2 := &Command{Name: "c2", PersistentPreRun: func(*Context) error { fmt.Println("c2"); return nil }}
	// c3 := &Command{Name: "c3", Run: func(*Context) error { fmt.Println("c3"); return nil }}
	//
	// c1.SubCommand(c2)
	// c2.SubCommand(c3)
	//
	// When c3 is run, stdout will see
	// c1
	// c2
	// c3
	//
	// If any of the nested commands return an error, all execution is stopped and that error is returned.
	PersistentPreRun func(*Context) error

	// PersistentPostRun is a function that will run after PostRun and will be run for any subcommands of this command.
	//
	// PersistentPostRun commands are executed in LIFO order
	//
	// ex:
	// c1 := &Command{Name: "c1", PersistentPostRun: func(*Context) error { fmt.Println("c1"); return nil }}
	// c2 := &Command{Name: "c2", PersistentPostRun: func(*Context) error { fmt.Println("c2"); return nil }}
	// c3 := &Command{Name: "c3", Run: func(*Context) error { fmt.Println("c3"); return nil }}
	//
	// c1.subCommand(c2)
	// c2.subCommand(c3)
	//
	// when c3 is run, stdout will see
	// c3
	// c2
	// c1
	//
	// If any of the nested commands return an error, execution will stop and the error will be returned, unless DeferPost
	// is true, in which case the error will be recorded and returned at the end, but the remaining functions will still
	// execute
	PersistentPostRun func(*Context) error

	// DeferPost will defer PersistentPostRun and PostRun functions.
	// This will cause them to run even if the Run function exits with an error.
	// Setting this value is persistent, meaning any subcommands from where this is set will
	// also defer their post run functions
	DeferPost bool

	parent   *Command
	commands map[string]*Command
	errs     []error
}

func (c *Command) ExecuteContext(ctx context.Context) error {
	args := os.Args[1:]
	// todo: parse flags

	err := c.execute(&Context{
		Context: ctx,
		args:    args,
	})
	if mErr := multierr.Append(err, multierr.Combine(c.errs...)); mErr != nil {
		// todo: print usage and error(s)
		return mErr
	}
	return nil
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
	// todo: how to handle commands with the same name
	// 	warn? silently overwrite? return err?
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}
	c.commands[cmd.Name] = cmd
	cmd.parent = c
}

func (c *Command) execute(ctx *Context) error {
	// append pre run functions to be executed in order
	if c.PersistentPreRun != nil {
		ctx.preRuns = append(ctx.preRuns, c.PersistentPreRun)
	}
	// prepend post run functions, to be executed in reverse order
	if c.PersistentPostRun != nil {
		ctx.postRuns = append([]func(*Context) error{c.PersistentPostRun}, ctx.postRuns...)
	}

	if c.DeferPost {
		ctx.deferPost = true
	}

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

	for _, run := range ctx.preRuns {
		if err := run(ctx); err != nil {
			return err
		}
	}
	if err := c.PreRun(ctx); err != nil {
		return err
	}

	var runErr error
	defer func() {
		p := recover()
		if p != nil && !ctx.deferPost {
			// a panic happened, but DeferPost isn't set,
			panic(p)
		}
		if runErr != nil && !ctx.deferPost {
			return
		}
		if err := c.PostRun(ctx); err != nil {
			c.errs = append(c.errs, err)
			if !ctx.deferPost {
				return
			}
		}
		for _, f := range ctx.postRuns {
			if err := f(ctx); err != nil {
				c.errs = append(c.errs, err)
				if !ctx.deferPost {
					return
				}
			}
		}
	}()
	runErr = c.Run(ctx)
	return runErr
}
