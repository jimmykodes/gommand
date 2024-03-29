// Code generated by flagger; DO NOT EDIT.

package flags

import (
	"strconv"
)

var _ Flag = (*uint64Flag)(nil)

type uint64Flag struct {
	*baseFlag

	defValue uint64
	value    uint64
}

func (f *uint64Flag) Type() FlagType { return Uint64FlagType }

func (f *uint64Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint64Flag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	f.value = uint64(v)
	f.set = true
	return nil
}

func (f *uint64Flag) Required() Flag {
	f.req = true
	return f
}

func (f *uint64Flag) AddSources(sources ...Valuer) Flag {
	f.addSources(sources...)
	return f
}

func Uint64Flag(name string, value uint64, usage string) Flag {
	return &uint64Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint64FlagS(name string, shorthand rune, value uint64, usage string) Flag {
	return &uint64Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint64(name string, value uint64, usage string) {
	fs.AddFlag(Uint64Flag(name, value, usage))
}

func (fs *FlagSet) Uint64S(name string, shorthand rune, value uint64, usage string) {
	fs.AddFlag(Uint64FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint64Val(name string) (uint64, error) {
	f, err := fs.flag(name, Uint64FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(uint64), nil
}

func (g FlagGetter) LookupUint64(name string) (uint64, error) {
	return g.fs.uint64Val(name)
}

func (g FlagGetter) Uint64(name string) uint64 {
	v, _ := g.LookupUint64(name)
	return v
}
