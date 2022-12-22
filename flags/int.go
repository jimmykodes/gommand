package flags

import (
	"strconv"
)

type intFlag struct {
	*baseFlag

	defValue int
	value    int
}

func (f *intFlag) Type() FlagType { return IntFlagType }

func (f *intFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *intFlag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	f.value = int(v)
	f.set = true
	return nil
}

func IntFlag(name string, value int, usage string) Flag {
	return &intFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func IntFlagS(name string, shorthand rune, value int, usage string) Flag {
	return &intFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int(name string, value int, usage string) {
	fs.AddFlag(IntFlag(name, value, usage))
}

func (fs *FlagSet) IntS(name string, shorthand rune, value int, usage string) {
	fs.AddFlag(IntFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) intVal(name string) (int, error) {
	f, err := fs.flag(name, IntFlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(int), nil
}

func (g FlagGetter) LookupInt(name string) (int, error) {
	return g.fs.intVal(name)
}

func (g FlagGetter) Int(name string) int {
	v, _ := g.LookupInt(name)
	return v
}
