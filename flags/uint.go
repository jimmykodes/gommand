package flags

import (
	"strconv"
)

type uintFlag struct {
	*baseFlag

	defValue uint
	value    uint
}

func (f *uintFlag) Type() FlagType { return UintFlagType }

func (f *uintFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uintFlag) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	f.value = uint(v)
	f.set = true
	return nil
}

func UintFlag(name string, value uint, usage string) Flag {
	return &uintFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func UintFlagS(name string, shorthand rune, value uint, usage string) Flag {
	return &uintFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint(name string, value uint, usage string) {
	fs.addFlag(UintFlag(name, value, usage))
}

func (fs *FlagSet) UintS(name string, shorthand rune, value uint, usage string) {
	fs.addFlag(UintFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uintVal(name string) (uint, error) {
	f, err := fs.flag(name, UintFlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(uint), nil
}

func (g FlagGetter) LookupUint(name string) (uint, error) {
	return g.fs.uintVal(name)
}

func (g FlagGetter) Uint(name string) uint {
	v, _ := g.LookupUint(name)
	return v
}
