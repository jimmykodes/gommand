package gommand

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/jimmykodes/gommand/flags"
	"github.com/jimmykodes/gommand/internal/lexer"
)

var (
	ErrNoRunner     = errors.New("gommand: command has no run function")
	ErrNoSubcommand = errors.New("gommand: must specify a subcommand")
	errShowHelp     = errors.New("show help")
	errShowVersion  = errors.New("show version")
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
	// Anything included after a space is expected to be usage descriptions
	// General syntax guidance
	//   ... indicates multiple of the preceding argument can be provided
	//   [ ] indicates optional arguments
	//   { } indicates a set of mutually exclusive required arguments
	//   |   indicates mutually exclusive arguments, where only one value in the set
	//       should be provided at a time. As described above, if the set of arguments
	//       are optional, the set should be enclosed in [ ] otherwise they should be
	//       enclosed in { }
	//
	// Example: create {--from-file file | --from-gcs bucket} [-d destination] file_name...
	Name string

	// Usage is the short explanation of the command
	Usage string

	// Description is the longer description of the command printed out by the help text
	Description string

	// Aliases are aliases for the current command.
	//
	// Ex:
	// c1 := &Command{Name: "items"}
	// c2 := &Command{Name: "list", Aliases: []string{"ls", "l"}}
	// c1.SubCommand(c2)
	//
	// items list
	// items ls
	// items l
	//
	// All are valid ways of executing the `list` command
	Aliases []string

	// Version is the value that will be printed when `--version` is passed to the command.
	// When retrieving the command version, the call tree is traversed backwards until a Command
	// is reached that has a non-zero value for the version. This means that it is possible
	// to version individual branches of the call tree, though this is not recommended. It is
	// intended to be set at the root of the tree, ideally through a package level var that can
	// be set using ldflags at build time
	// ie: go build -ldflags="cmd.Version=1.1.0"
	Version string

	// ArgValidator is an ArgValidator to be called on the args of the function being executed. This is called before any of
	// the functions for this command are called.
	// If this is not defined ArgsNone is used.
	ArgValidator ArgValidator

	// Flags are a slice of flags.Flag that will be used to initialize the command's FlagSet
	FlagSet *flags.FlagSet

	// PersistentFlags are a slice of flags.Flag that will be used to initialize the command's PersistentFlagSet
	PersistentFlagSet *flags.FlagSet

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
	// c1.SubCommand(c2)
	// c2.SubCommand(c3)
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

	// SilenceHelp will not print the help message if the command exits with an error.
	// This field will propagate to subcommands and cannot be overwritten by the child, so if any
	// point of a command's upstream lineage has the value set, the help message will be silenced
	SilenceHelp bool

	// SilenceError is like SilenceHelp but does not print the "Error: xxx" message when the command
	// exits with an error
	SilenceError bool

	parent   *Command
	commands commands

	err error
}

func (c *Command) ExecuteContext(ctx context.Context) error {
	cmdCtx := &Context{
		Context: ctx,
		args:    os.Args[1:],
	}
	err := c.execute(cmdCtx)
	if errors.Is(err, errShowHelp) {
		cmdCtx.cmd.help()
		return nil
	}
	if errors.Is(err, errShowVersion) {
		v := c._version()
		if v == "" {
			v = "N/A"
		}
		fmt.Println(v)
		return nil
	}

	if mErr := errors.Join(err, c.err); mErr != nil {
		if !cmdCtx.silenceHelp {
			cmdCtx.cmd.help()
		}
		if !cmdCtx.silenceError {
			fmt.Println("Error:", mErr)
		}
		return mErr
	}
	return nil
}

func (c *Command) Execute() error {
	return c.ExecuteContext(context.Background())
}

func (c *Command) SubCommand(cmds ...*Command) {
	for _, command := range cmds {
		c.subCommand(command)
	}
}

func (c *Command) name() (name string) {
	name, _, _ = strings.Cut(c.Name, " ")
	return
}

func (c *Command) _version() string {
	_c := c
	for _c.Version == "" {
		if _c.parent == nil {
			break
		}
		_c = _c.parent
	}
	return _c.Version
}

func (c *Command) help() {
	if c.Description != "" {
		fmt.Println(c.Description)
		fmt.Println()
	} else if c.Usage != "" {
		fmt.Println(c.Usage)
		fmt.Println()
	}

	if v := c._version(); v != "" {
		fmt.Println("Version:", v)
		fmt.Println()
	}

	fmt.Println("Usage:")
	usage := []string{c.Name}
	for parent := c.parent; parent != nil; parent = parent.parent {
		usage = append([]string{parent.name()}, usage...)
	}
	fmt.Print("  ", strings.Join(usage, " "))

	if len(c.commands) > 0 {
		fmt.Print(" [commands]")
	}
	fmt.Println()
	fmt.Println()

	if len(c.Aliases) > 0 {
		fmt.Println("Aliases:")
		fmt.Printf("  %s\n\n", strings.Join(c.Aliases, ", "))
	}

	fs := flags.NewFlagSet(flags.WithHelpFlag()).AddFlagSet(c.FlagSet)
	pfs := flags.NewFlagSet()

	for p := c; p != nil; p = p.parent {
		pfs.AddFlagSet(p.PersistentFlagSet)
	}

	if len(c.commands) > 0 {
		fmt.Println("Available Commands:")
		fmt.Println(c.commands)
	}

	fsStr := fs.Repr()
	pfsStr := pfs.Repr()

	fmt.Println("Flags:")
	fmt.Println(fsStr)

	if pfsStr != "" {
		fmt.Println()
		fmt.Println("Global Flags:")
		fmt.Println(pfsStr)
	}
}

func (c *Command) subCommand(cmd *Command) {
	// todo: how to handle commands with the same name
	// 	warn? silently overwrite? return err?
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}
	c.commands[cmd.name()] = cmd
	for _, alias := range cmd.Aliases {
		c.commands[alias] = cmd
	}
	cmd.parent = c
}

func (c *Command) hasSubCommands() bool {
	return len(c.commands) > 0
}

func (c *Command) execute(ctx *Context) error {
	// ################
	// Append any persistent configs
	// ################
	ctx.cmd = c
	ctx.addPersistentFlags(c.PersistentFlagSet)

	// append pre run functions to be executed in order
	ctx.preRuns = append(ctx.preRuns, func(ctx *Context) error { return nil })
	if c.PersistentPreRun != nil {
		ctx.preRuns[ctx.depth] = c.PersistentPreRun
	}

	// append post run functions, to be defered in order
	if c.PersistentPostRun != nil {
		ctx.postRuns = append(ctx.postRuns, c.PersistentPostRun)
	}

	if c.DeferPost {
		ctx.deferPost = true
	}

	if c.SilenceError {
		ctx.silenceError = true
	}

	if c.SilenceHelp {
		ctx.silenceHelp = true
	}

	// ################
	// Walk the command tree
	// ################

	if c.hasSubCommands() {
		if len(ctx.args) == 0 {
			return fmt.Errorf("%s: %w", c.Name, ErrNoSubcommand)
		}
		next := c.commands[ctx.args[0]]
		if next != nil {
			ctx.depth++
			ctx.args = ctx.args[1:]
			return next.execute(ctx)
		}
	}

	// ################
	// Process Flags
	// ################

	fs := flags.NewFlagSet()
	fs.AddFlagSet(ctx.persistentFlags())
	fs.AddFlagSet(c.FlagSet)

	ctx.flagGetter = flags.NewFlagGetter(fs)

	var (
		argLexer    = lexer.New(ctx.args)
		args        []string
		collectArgs bool
	)
	for {
		token := argLexer.Read()
		if token == nil {
			break
		}
		switch token.Type {
		case lexer.ValueType:
			collectArgs = true
			args = append(args, token.Value)
		case lexer.MultiFlagType:
			if collectArgs {
				return fmt.Errorf("gommand: invalid flag position: flags must come before args: -%s", token.Name)
			}
			if token.Value != "" {
				return fmt.Errorf("gommand: invalid multi-flag: cannot assign value to multi-flag: -%s", token.Name)
			}
			for _, chr := range token.Name {
				if chr == 'h' {
					return errShowHelp
				}
				f := fs.FromShort(chr)
				if f == nil {
					return fmt.Errorf("gommand: missing flag: -%s", string(chr))
				}
				if f.Type() != flags.BoolFlagType {
					return fmt.Errorf("gommand: invalid multi-flag: -%s is not a boolean flag", string(chr))
				}
				_ = f.Set("true")
			}
		default:
			if collectArgs {
				prefix := "-"
				if token.Type == lexer.LongFlagType {
					prefix += "-"
				}
				return fmt.Errorf("gommand: invalid flag position: flags must come before args: %s%s", prefix, token.Name)
			}
			var f flags.Flag
			switch token.Type {
			case lexer.ShortFlagType:
				if token.Name == "h" {
					return errShowHelp
				}
				f = fs.FromShort(rune(token.Name[0]))
			case lexer.LongFlagType:
				if token.Name == "help" {
					return errShowHelp
				}
				if token.Name == "version" {
					return errShowVersion
				}
				f = fs.FromName(token.Name)
			}

			if f == nil {
				prefix := "-"
				if token.Type == lexer.LongFlagType {
					prefix += "-"
				}
				return fmt.Errorf("gommand: missing flag: %s%s", prefix, token.Name)
			}

			var setErr error
			if token.Value == "" {
				if f.Type() == flags.BoolFlagType {
					setErr = f.Set("true")
				} else {
					if peekToken := argLexer.Peek(); peekToken != nil && peekToken.Type == lexer.ValueType {
						setErr = f.Set(argLexer.Read().Value)
					}
				}
			} else {
				setErr = f.Set(token.Value)
			}
			if setErr != nil {
				return setErr
			}
		}
	}
	// ################
	// Validate args
	// ################

	ctx.args = args

	validator := c.ArgValidator
	if validator == nil {
		// default to allowing no args unless specified otherwise.
		validator = ArgsNone()
	}
	if err := validator(ctx.args); err != nil {
		return fmt.Errorf("gommand: invalid args: %w", err)
	}

	// ################
	// Run the things!
	// ################
	return c.run(ctx)
}

func (c *Command) run(ctx *Context) (runErr error) {
	// if there is no Run command, no need to do pre/post run setup things
	if c.Run == nil {
		return ErrNoRunner
	}

	for depth, run := range ctx.preRuns {
		fs := ctx.persistentFlagSets[depth]
		if fs == nil {
			continue
		}
		fg := flags.NewFlagGetter(fs)
		for _, f := range fg.All() {
			if f.IsRequired() && !f.IsSet() {
				if err := flags.SetFromSources(f); err != nil {
					return err
				}
				if !f.IsSet() {
					return flags.ErrMissingRequiredFlag{Flag: f}
				}
			}
		}
		if err := run(ctx); err != nil {
			return err
		}
	}

	if c.FlagSet != nil {
		fg := flags.NewFlagGetter(c.FlagSet)
		for _, f := range fg.All() {
			if f.IsRequired() && !f.IsSet() {
				if err := flags.SetFromSources(f); err != nil {
					return err
				}
				if !f.IsSet() {
					return flags.ErrMissingRequiredFlag{Flag: f}
				}
			}
		}
	}

	if c.PreRun != nil {
		if err := c.PreRun(ctx); err != nil {
			return err
		}
	}

	defer func() {
		if runErr != nil && !ctx.deferPost {
			return
		}
		if c.PostRun != nil {
			if err := c.PostRun(ctx); err != nil {
				c.err = errors.Join(c.err, err)
				if !ctx.deferPost {
					return
				}
			}
		}
		for i := len(ctx.postRuns) - 1; i >= 0; i-- {
			if err := ctx.postRuns[i](ctx); err != nil {
				c.err = errors.Join(c.err, err)
				if !ctx.deferPost {
					return
				}
			}
		}
	}()

	defer func() {
		if p := recover(); p != nil {
			if !ctx.deferPost {
				panic(p)
			} else {
				runErr = errors.Join(runErr, fmt.Errorf("panic: %v", p))
			}
		}
	}()

	runErr = c.Run(ctx)
	return runErr
}

type commands map[string]*Command

func (c commands) String() string {
	var (
		sb     strings.Builder
		maxKey int
		keys   = make([]string, 0, len(c))
	)
	for name, command := range c {
		if slices.Contains(command.Aliases, name) {
			continue
		}
		keys = append(keys, name)
		if l := len(name); l > maxKey {
			maxKey = l
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		padding := maxKey - len(k)
		_, _ = fmt.Fprintf(&sb, "  %s%s  %s\n", k, strings.Repeat(" ", padding), c[k].Usage)
	}
	return sb.String()
}
