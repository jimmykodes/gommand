package flags

import (
	"os"
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
	StringSliceFlagType
	BoolSliceFlagType
	DurationSliceFlagType
	IntSliceFlagType
	Int8SliceFlagType
	Int16SliceFlagType
	Int32SliceFlagType
	Int64SliceFlagType
	UintSliceFlagType
	Uint8SliceFlagType
	Uint16SliceFlagType
	Uint32SliceFlagType
	Uint64SliceFlagType
	Float32SliceFlagType
	Float64SliceFlagType
)

var (
	sliceSeparator = getSliceSep()
)

func getSliceSep() string {
	sep := os.Getenv("GOMMAND_SLICE_SEPARATOR")
	if sep == "" {
		sep = ","
	}
	return sep
}

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

func Stringer(flag Flag, nameLen int, hasShort bool) string {
	var sb strings.Builder
	sb.WriteString("  ")
	if hasShort {
		if flag.Short() > 0 {
			sb.WriteString("-")
			sb.WriteRune(flag.Short())
			sb.WriteString(", ")
		} else {
			sb.WriteString("    ")
		}
	}

	sb.WriteString("--")
	sb.WriteString(flag.Name())
	sb.WriteString(strings.Repeat(" ", nameLen-len(flag.Name())))
	sb.WriteString("  ")
	sb.WriteString(flag.Usage())
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
