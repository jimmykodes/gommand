// Code generated by flagger; DO NOT EDIT.

package flags

import (
	"strconv"
)

var _ Flag = (*uint16Flag)(nil)

type uint16Flag struct {
	*baseFlag

	defValue uint16
	value    uint16
}

func (f *uint16Flag) Type() FlagType { return Uint16FlagType }

func (f *uint16Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint16Flag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return err
	}
	f.value = uint16(v)
	f.set = true
	return nil
}

func (f *uint16Flag) Required() Flag {
	f.req = true
	return f
}

func (f *uint16Flag) AddSources(sources ...Valuer) Flag {
	f.addSources(sources...)
	return f
}

func Uint16Flag(name string, value uint16, usage string) Flag {
	return &uint16Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint16FlagS(name string, shorthand rune, value uint16, usage string) Flag {
	return &uint16Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint16(name string, value uint16, usage string) {
	fs.AddFlag(Uint16Flag(name, value, usage))
}

func (fs *FlagSet) Uint16S(name string, shorthand rune, value uint16, usage string) {
	fs.AddFlag(Uint16FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint16Val(name string) (uint16, error) {
	f, err := fs.flag(name, Uint16FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(uint16), nil
}

func (g FlagGetter) LookupUint16(name string) (uint16, error) {
	return g.fs.uint16Val(name)
}

func (g FlagGetter) Uint16(name string) uint16 {
	v, _ := g.LookupUint16(name)
	return v
}
