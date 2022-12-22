package flags

import (
	"strconv"
	"strings"
)

type boolSliceFlag struct {
	*baseFlag

	defValue []bool
	value    []bool
}

func (f *boolSliceFlag) Type() FlagType { return BoolSliceFlagType }

func (f *boolSliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *boolSliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]bool, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseBool(piece)
		if err != nil {
			return err
		}
		v[i] = val
	}
	f.value = v
	f.set = true
	return nil
}

func BoolSliceFlag(name string, value []bool, usage string) Flag {
	return &boolSliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func BoolSliceFlagS(name string, shorthand rune, value []bool, usage string) Flag {
	return &boolSliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) BoolSlice(name string, value []bool, usage string) {
	fs.AddFlag(BoolSliceFlag(name, value, usage))
}

func (fs *FlagSet) BoolSliceS(name string, shorthand rune, value []bool, usage string) {
	fs.AddFlag(BoolSliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) boolSliceVal(name string) ([]bool, error) {
	f, err := fs.flag(name, BoolSliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]bool), nil
}

func (g FlagGetter) LookupBoolSlice(name string) ([]bool, error) {
	return g.fs.boolSliceVal(name)
}

func (g FlagGetter) BoolSlice(name string) []bool {
	v, _ := g.LookupBoolSlice(name)
	return v
}
