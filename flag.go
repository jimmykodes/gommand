package gommand

import (
	"fmt"
	"strconv"
	"strings"
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

func flagStringer(flag Flag, nameLen int) string {
	var sb strings.Builder

	_, _ = fmt.Fprint(&sb, "  ")
	if flag.Short() > 0 {
		_, _ = fmt.Fprint(&sb, "-", string(byte(flag.Short())), ", ")
	} else {
		_, _ = fmt.Fprint(&sb, "    ")
	}

	_, _ = fmt.Fprint(
		&sb,
		"--",
		flag.Name(),
		strings.Repeat(" ", nameLen-len(flag.Name())),
		"  ",
		flag.Usage(),
	)
	return sb.String()
}

func StringFlag(name, value, usage string) Flag {
	return &stringFlag{
		name:     name,
		usage:    usage,
		defValue: value,
	}
}

func StringFlagS(name string, shorthand rune, value, usage string) Flag {
	return &stringFlag{
		name:     name,
		short:    shorthand,
		usage:    usage,
		defValue: value,
	}
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

type float64Flag struct {
	name      string
	short     rune
	usage     string
	set       bool
	req       bool
	envPrefix string

	defValue float64
	Value    float64
}

func (f *float64Flag) Type() FlagType    { return Float64FlagType }
func (f *float64Flag) Name() string      { return f.name }
func (f *float64Flag) Short() rune       { return f.short }
func (f *float64Flag) Usage() string     { return f.usage }
func (f *float64Flag) IsSet() bool       { return f.set }
func (f *float64Flag) IsRequired() bool  { return f.req }
func (f *float64Flag) EnvPrefix() string { return f.envPrefix }

func (f *float64Flag) Required() {
	f.req = true
}

func (f *float64Flag) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	f.Value = v
	f.set = true
	return nil
}

func (f *float64Flag) SetEnvPrefix(prefix string) {
	f.envPrefix = prefix
}

type stringFlag struct {
	name      string
	short     rune
	usage     string
	set       bool
	req       bool
	envPrefix string

	defValue string
	Value    string
}

func (f *stringFlag) Type() FlagType    { return StringFlagType }
func (f *stringFlag) Name() string      { return f.name }
func (f *stringFlag) Short() rune       { return f.short }
func (f *stringFlag) Usage() string     { return f.usage }
func (f *stringFlag) IsSet() bool       { return f.set }
func (f *stringFlag) IsRequired() bool  { return f.req }
func (f *stringFlag) EnvPrefix() string { return f.envPrefix }

func (f *stringFlag) Required() {
	f.req = true
}

func (f *stringFlag) Set(s string) error {
	f.Value = s
	f.set = true
	return nil
}

func (f *stringFlag) SetEnvPrefix(prefix string) {
	f.envPrefix = prefix
}
