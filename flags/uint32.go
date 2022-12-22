package flags

import (
	"strconv"
)

type uint32Flag struct {
	*baseFlag

	defValue uint32
	value    uint32
}

func (f *uint32Flag) Type() FlagType { return Uint32FlagType }

func (f *uint32Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint32Flag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return err
	}
	f.value = uint32(v)
	f.set = true
	return nil
}

func Uint32Flag(name string, value uint32, usage string) Flag {
	return &uint32Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint32FlagS(name string, shorthand rune, value uint32, usage string) Flag {
	return &uint32Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint32(name string, value uint32, usage string) {
	fs.AddFlag(Uint32Flag(name, value, usage))
}

func (fs *FlagSet) Uint32S(name string, shorthand rune, value uint32, usage string) {
	fs.AddFlag(Uint32FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint32Val(name string) (uint32, error) {
	f, err := fs.flag(name, Uint32FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(uint32), nil
}

func (g FlagGetter) LookupUint32(name string) (uint32, error) {
	return g.fs.uint32Val(name)
}

func (g FlagGetter) Uint32(name string) uint32 {
	v, _ := g.LookupUint32(name)
	return v
}
