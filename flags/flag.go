package flags

import (
	"fmt"
	"strings"
)

type FlagType int

const (
	UnknownFlagType FlagType = iota
	StringFlagType
	BoolFlagType
	DurationFlagType
	IntFlagType
	Int8FlagType
	Int16FlagType
	Int32FlagType
	Int64FlagType
	UintFlagType
	Uint8FlagType
	Uint16FlagType
	Uint32FlagType
	Uint64FlagType
	Float32FlagType
	Float64FlagType
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
	Value() any

	Set(string) error
	SetEnvPrefix(string)
}

func Stringer(flag Flag, nameLen int) string {
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

type baseFlag struct {
	name      string
	short     rune
	usage     string
	set       bool
	req       bool
	envPrefix string
}

func (f *baseFlag) Type() FlagType             { return UnknownFlagType }
func (f *baseFlag) Name() string               { return f.name }
func (f *baseFlag) Short() rune                { return f.short }
func (f *baseFlag) Usage() string              { return f.usage }
func (f *baseFlag) IsSet() bool                { return f.set }
func (f *baseFlag) IsRequired() bool           { return f.req }
func (f *baseFlag) EnvPrefix() string          { return f.envPrefix }
func (f *baseFlag) Required()                  { f.req = true }
func (f *baseFlag) SetEnvPrefix(prefix string) { f.envPrefix = prefix }
