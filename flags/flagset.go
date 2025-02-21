package flags

import (
	"errors"
	"maps"
	"sort"
	"strings"
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
	flags      map[string]Flag
	shortFlags map[rune]Flag
	sources    []Valuer
}

func (fs *FlagSet) FromName(name string) Flag {
	return fs.flags[name]
}

func (fs *FlagSet) FromShort(short rune) Flag {
	return fs.shortFlags[short]
}

func (fs *FlagSet) AddFlag(f Flag) *FlagSet {
	fs.flags[f.Name()] = f
	if f.Short() != 0 {
		fs.shortFlags[f.Short()] = f
	}
	f.AddSources(fs.sources...)
	return fs
}

func (fs *FlagSet) AddFlags(flags ...Flag) *FlagSet {
	for _, f := range flags {
		fs.AddFlag(f)
	}
	return fs
}

func (fs *FlagSet) AddSource(source Valuer) *FlagSet {
	fs.sources = append(fs.sources, source)
	return fs
}

func (fs *FlagSet) Repr() string {
	var (
		names    = make([]string, 0, len(fs.flags))
		maxLen   = 0
		hasShort = false
	)
	for n, flag := range fs.flags {
		names = append(names, n)
		if l := len(flag.Name()); l > maxLen {
			maxLen = l
		}
		if flag.Short() != 0 {
			hasShort = true
		}
	}
	sort.Strings(names)
	strs := make([]string, len(names))
	for i, name := range names {
		strs[i] = Stringer(fs.flags[name], maxLen, hasShort)
	}
	return strings.Join(strs, "\n")
}

func (fs *FlagSet) addHelpFlag() {
	fs.BoolS("help", 'h', false, "show this help message")
}

func (fs *FlagSet) AddFlagSet(set *FlagSet) *FlagSet {
	if set == nil {
		return fs
	}

	maps.Copy(fs.flags, set.flags)
	maps.Copy(fs.shortFlags, set.shortFlags)

	return fs
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

	if !f.IsSet() {
		for _, source := range f.Sources() {
			if val, ok := source.Value(name); ok {
				if err := f.Set(val); err != nil {
					return nil, err
				}
				break
			}
		}
	}

	if !f.IsSet() && f.IsRequired() {
		// was not set by the command line _or_ the environment
		return nil, ErrMissingRequiredFlag
	}

	return f, nil
}
