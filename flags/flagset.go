package flags

import (
	"errors"
	"os"

	"github.com/jimmykodes/strman"
)

var (
	ErrUnregisteredFlag    = errors.New("gommand: flag not defined")
	ErrInvalidFlagType     = errors.New("gommand: invalid flag type")
	ErrMissingRequiredFlag = errors.New("gommand: missing required flag")
)

func NewFlagSet(options ...FlagSetOption) *FlagSet {
	f := &FlagSet{flags: make(map[string]Flag), shortFlags: make(map[rune]Flag)}
	for _, option := range options {
		option.Apply(f)
	}
	return f
}

type FlagSet struct {
	envPrefix  string
	noEnv      bool
	flags      map[string]Flag
	shortFlags map[rune]Flag
}

func (fs *FlagSet) FromName(name string) Flag {
	return fs.flags[name]
}

func (fs *FlagSet) FromShort(short rune) Flag {
	return fs.shortFlags[short]
}

func (fs *FlagSet) addFlag(f Flag) {
	f.SetEnvPrefix(fs.envPrefix)
	fs.flags[f.Name()] = f
	if f.Short() != 0 {
		fs.shortFlags[f.Short()] = f
	}
}

func (fs *FlagSet) addHelpFlag() {
	fs.BoolS("help", 'h', false, "show this help message")
}

func (fs *FlagSet) AddFlagSet(set *FlagSet) {
	for name, flag := range set.flags {
		fs.flags[name] = flag
	}
	for short, flag := range set.shortFlags {
		fs.shortFlags[short] = flag
	}
}

func (fs *FlagSet) MarkRequired(name string) error {
	f, ok := fs.flags[name]
	if !ok {
		return ErrUnregisteredFlag
	}
	f.Required()
	return nil
}

func (fs *FlagSet) flag(name string, t FlagType) (Flag, error) {
	f, ok := fs.flags[name]
	if !ok {
		return nil, ErrUnregisteredFlag
	}
	if f.Type() != t {
		return nil, ErrInvalidFlagType
	}
	if !f.IsSet() && !fs.noEnv {
		// flag was not set from command line, and env lookup is not turned off
		if val, ok := fromEnv(f.EnvPrefix(), name); ok {
			if err := f.Set(val); err != nil {
				return nil, err
			}
		}
	}
	if !f.IsSet() && f.IsRequired() {
		// was not set by the command line _or_ the environment
		return nil, ErrMissingRequiredFlag
	}
	return f, nil
}

func fromEnv(prefix, name string) (string, bool) {
	p := prefix
	if p != "" {
		p = prefix + "_"
	}
	return os.LookupEnv(strman.ToScreamingSnake(p + name))
}
