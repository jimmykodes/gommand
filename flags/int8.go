package flags

import (
	"strconv"
)

type int8Flag struct {
	*baseFlag

	defValue int8
	value    int8
}

func (f *int8Flag) Type() FlagType { return Int8FlagType }

func (f *int8Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int8Flag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return err
	}
	f.value = int8(v)
	f.set = true
	return nil
}

func Int8Flag(name string, value int8, usage string) Flag {
	return &int8Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int8FlagS(name string, shorthand rune, value int8, usage string) Flag {
	return &int8Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int8(name string, value int8, usage string) {
	fs.AddFlag(Int8Flag(name, value, usage))
}

func (fs *FlagSet) Int8S(name string, shorthand rune, value int8, usage string) {
	fs.AddFlag(Int8FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int8Val(name string) (int8, error) {
	f, err := fs.flag(name, Int8FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(int8), nil
}

func (g FlagGetter) LookupInt8(name string) (int8, error) {
	return g.fs.int8Val(name)
}

func (g FlagGetter) Int8(name string) int8 {
	v, _ := g.LookupInt8(name)
	return v
}
