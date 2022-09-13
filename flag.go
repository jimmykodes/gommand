package gommand

import (
	"strconv"
)

type FlagType int

const (
	UnknownFlagType FlagType = iota
	StringFlagType
	IntFlagType
	Float64FlagType
	BoolFlagType
)

type Flag interface {
	Type() FlagType
	Name() string
	Short() rune
	Usage() string
	IsSet() bool
	Required()
	IsRequired() bool
	EnvPrefix() string

	Set(string) error
	SetEnvPrefix(string)
}

func IntFlag(name string, value int, usage string) Flag {
	return &intFlag{
		name:     name,
		usage:    usage,
		defValue: value,
	}
}

func IntFlagS(name string, shorthand rune, value int, usage string) Flag {
	return &intFlag{
		name:     name,
		short:    shorthand,
		usage:    usage,
		defValue: value,
	}
}

func BoolFlag(name string, value bool, usage string) Flag {
	return &boolFlag{
		name:     name,
		usage:    usage,
		defValue: value,
	}
}

func BoolFlagS(name string, shorthand rune, value bool, usage string) Flag {
	return &boolFlag{
		name:     name,
		short:    shorthand,
		usage:    usage,
		defValue: value,
	}
}

type rawFlag struct {
	name  string
	value string
}

type intFlag struct {
	name      string
	short     rune
	usage     string
	set       bool
	req       bool
	envPrefix string

	defValue int
	Value    int
}

func (f *intFlag) Type() FlagType    { return IntFlagType }
func (f *intFlag) Name() string      { return f.name }
func (f *intFlag) Short() rune       { return f.short }
func (f *intFlag) Usage() string     { return f.usage }
func (f *intFlag) IsSet() bool       { return f.set }
func (f *intFlag) IsRequired() bool  { return f.req }
func (f *intFlag) EnvPrefix() string { return f.envPrefix }

func (f *intFlag) Required() {
	f.req = true
}

func (f *intFlag) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	f.Value = v
	f.set = true
	return nil
}

func (f *intFlag) SetEnvPrefix(prefix string) {
	f.envPrefix = prefix
}

type boolFlag struct {
	name      string
	short     rune
	usage     string
	set       bool
	req       bool
	envPrefix string

	defValue bool
	Value    bool
}

func (f *boolFlag) Type() FlagType    { return BoolFlagType }
func (f *boolFlag) Name() string      { return f.name }
func (f *boolFlag) Short() rune       { return f.short }
func (f *boolFlag) Usage() string     { return f.usage }
func (f *boolFlag) IsSet() bool       { return f.set }
func (f *boolFlag) IsRequired() bool  { return f.req }
func (f *boolFlag) EnvPrefix() string { return f.envPrefix }

func (f *boolFlag) Required() {
	f.req = true
}

func (f *boolFlag) Set(s string) error {
	if s == "" {
		// bool flags are unique in that an empty string is treated as truthy
		f.Value = true
		f.set = true
		return nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	f.Value = v
	f.set = true
	return nil
}

func (f *boolFlag) SetEnvPrefix(prefix string) {
	f.envPrefix = prefix
}
