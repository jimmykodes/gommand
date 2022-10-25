package flags

import (
	"strconv"
)

type int64Flag struct {
	*baseFlag

	defValue int64
	value    int64
}

func (f *int64Flag) Type() FlagType { return Int64FlagType }

func (f *int64Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int64Flag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	f.value = v
	f.set = true
	return nil
}

func Int64Flag(name string, value int64, usage string) Flag {
	return &int64Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int64FlagS(name string, shorthand rune, value int64, usage string) Flag {
	return &int64Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int64(name string, value int64, usage string) {
	fs.addFlag(Int64Flag(name, value, usage))
}

func (fs *FlagSet) Int64S(name string, shorthand rune, value int64, usage string) {
	fs.addFlag(Int64FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int64Val(name string) (int64, error) {
	f, err := fs.flag(name, Int64FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(int64), nil
}

func (g FlagGetter) LookupInt64(name string) (int64, error) {
	return g.fs.int64Val(name)
}

func (g FlagGetter) Int64(name string) int64 {
	v, _ := g.LookupInt64(name)
	return v
}
