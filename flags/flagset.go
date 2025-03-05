package flags

import (
	"fmt"
	"maps"
	"sort"
	"strings"
)

type ErrInvalidFlagType struct {
	Flag         Flag
	ExpectedType FlagType
}

func (e ErrInvalidFlagType) Error() string {
	return fmt.Sprintf("gommand: flag --%s is not of type %s", e.Flag.Name(), e.ExpectedType)
}

type ErrUnregisteredFlag struct {
	Name string
}

func (e ErrUnregisteredFlag) Error() string {
	return fmt.Sprintf("gommand: flag --%s not defined", e.Name)
}

type ErrMissingRequiredFlag struct {
	Flag Flag
}

func (e ErrMissingRequiredFlag) Error() string {
	return fmt.Sprintf("gommand: missing required flag: --%s", e.Flag.Name())
}

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
		return ErrUnregisteredFlag{Name: name}
	}
	f.Required()
	return nil
}

func (fs *FlagSet) flag(name string, t FlagType) (Flag, error) {
	f, ok := fs.flags[name]
	if !ok {
		return nil, ErrUnregisteredFlag{Name: name}
	}

	if f.Type() != t {
		return nil, ErrInvalidFlagType{Flag: f, ExpectedType: t}
	}

	if !f.IsSet() {
		if err := SetFromSources(f); err != nil {
			return nil, err
		}
	}

	if !f.IsSet() && f.IsRequired() {
		// was not set by the command line or any configured source
		return nil, ErrMissingRequiredFlag{Flag: f}
	}

	return f, nil
}
