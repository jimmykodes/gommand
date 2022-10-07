package flags

import (
	"strings"
)

type stringSliceFlag struct {
	*baseFlag

	defValue []string
	value    []string
}

func (f *stringSliceFlag) Type() FlagType { return StringSliceFlagType }

func (f *stringSliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *stringSliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	f.value = pieces
	f.set = true
	return nil
}

func StringSliceFlag(name string, value []string, usage string) Flag {
	return &stringSliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func StringSliceFlagS(name string, shorthand rune, value []string, usage string) Flag {
	return &stringSliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) StringSlice(name string, value []string, usage string) {
	fs.addFlag(StringSliceFlag(name, value, usage))
}

func (fs *FlagSet) StringSliceS(name string, shorthand rune, value []string, usage string) {
	fs.addFlag(StringSliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) stringSliceVal(name string) ([]string, error) {
	f, err := fs.flag(name, StringSliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]string), nil
}

func (g FlagGetter) LookupStringSlice(name string) ([]string, error) {
	return g.fs.stringSliceVal(name)
}

func (g FlagGetter) StringSlice(name string) []string {
	v, _ := g.LookupStringSlice(name)
	return v
}
