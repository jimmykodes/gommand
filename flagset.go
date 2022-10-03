package gommand

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

type FlagSet struct {
	flags      map[string]Flag
	envPrefix  string
	noEnv      bool
	shortFlags map[rune]Flag
}

func NewFlagSet(options ...FlagSetOption) *FlagSet {
	f := &FlagSet{flags: make(map[string]Flag), shortFlags: make(map[rune]Flag)}
	for _, option := range options {
		option.Apply(f)
	}
	return f
}

func (fs *FlagSet) String(name, value, usage string) {
	fs.addFlag(StringFlag(name, value, usage))
}

func (fs *FlagSet) StringS(name string, shorthand rune, value, usage string) {
	fs.addFlag(StringFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) Int(name string, value int, usage string) {
	fs.addFlag(IntFlag(name, value, usage))
}

func (fs *FlagSet) IntS(name string, shorthand rune, value int, usage string) {
	fs.addFlag(IntFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) Bool(name string, value bool, usage string) {
	fs.addFlag(BoolFlag(name, value, usage))
}

func (fs *FlagSet) BoolS(name string, shorthand rune, value bool, usage string) {
	fs.addFlag(BoolFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) addFlag(f Flag) {
	f.SetEnvPrefix(fs.envPrefix)
	fs.flags[f.Name()] = f
	if f.Short() != 0 {
		fs.shortFlags[f.Short()] = f
	}
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

func (fs *FlagSet) intVal(name string) (int, error) {
	f, err := fs.flag(name, IntFlagType)
	if err != nil {
		return 0, err
	}
	tf := f.(*intFlag)
	if tf.IsSet() {
		return tf.Value, nil
	}
	return tf.defValue, nil
}

func (fs *FlagSet) boolVal(name string) (bool, error) {
	f, err := fs.flag(name, BoolFlagType)
	if err != nil {
		return false, err
	}
	tf := f.(*boolFlag)
	if tf.IsSet() {
		return tf.Value, nil
	}
	return tf.defValue, nil
}

func (fs *FlagSet) float64Val(name string) (float64, error) {
	f, err := fs.flag(name, Float64FlagType)
	if err != nil {
		return 0, err
	}
	tf := f.(*float64Flag)
	if tf.IsSet() {
		return tf.Value, nil
	}
	return tf.defValue, nil
}

func (fs *FlagSet) stringVal(name string) (string, error) {
	f, err := fs.flag(name, StringFlagType)
	if err != nil {
		return "", err
	}
	tf := f.(*stringFlag)
	if tf.IsSet() {
		return tf.Value, nil
	}
	return tf.defValue, nil
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
