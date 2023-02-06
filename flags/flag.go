package flags

import (
	"os"
	"strings"
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

//go:generate flagger
type Flag interface {
	Type() FlagType
	Name() string
	Short() rune
	Usage() string
	IsSet() bool
	IsRequired() bool
	EnvPrefix() string
	Value() any

	Required() Flag

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
func (f *baseFlag) SetEnvPrefix(prefix string) { f.envPrefix = prefix }
