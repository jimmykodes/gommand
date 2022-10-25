package flags

import (
	"strconv"
)

type int32Flag struct {
	*baseFlag

	defValue int32
	value    int32
}

func (f *int32Flag) Type() FlagType { return Int32FlagType }

func (f *int32Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int32Flag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return err
	}
	f.value = int32(v)
	f.set = true
	return nil
}

func Int32Flag(name string, value int32, usage string) Flag {
	return &int32Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int32FlagS(name string, shorthand rune, value int32, usage string) Flag {
	return &int32Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int32(name string, value int32, usage string) {
	fs.addFlag(Int32Flag(name, value, usage))
}

func (fs *FlagSet) Int32S(name string, shorthand rune, value int32, usage string) {
	fs.addFlag(Int32FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int32Val(name string) (int32, error) {
	f, err := fs.flag(name, Int32FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(int32), nil
}

func (g FlagGetter) LookupInt32(name string) (int32, error) {
	return g.fs.int32Val(name)
}

func (g FlagGetter) Int32(name string) int32 {
	v, _ := g.LookupInt32(name)
	return v
}
