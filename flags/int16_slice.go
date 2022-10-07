package flags

import (
	"strconv"
	"strings"
)

type int16SliceFlag struct {
	*baseFlag

	defValue []int16
	value    []int16
}

func (f *int16SliceFlag) Type() FlagType { return Int16SliceFlagType }

func (f *int16SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int16SliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]int16, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 16)
		if err != nil {
			return err
		}
		v[i] = int16(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Int16SliceFlag(name string, value []int16, usage string) Flag {
	return &int16SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int16SliceFlagS(name string, shorthand rune, value []int16, usage string) Flag {
	return &int16SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int16Slice(name string, value []int16, usage string) {
	fs.addFlag(Int16SliceFlag(name, value, usage))
}

func (fs *FlagSet) Int16SliceS(name string, shorthand rune, value []int16, usage string) {
	fs.addFlag(Int16SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int16SliceVal(name string) ([]int16, error) {
	f, err := fs.flag(name, Int16SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]int16), nil
}

func (g FlagGetter) LookupInt16Slice(name string) ([]int16, error) {
	return g.fs.int16SliceVal(name)
}

func (g FlagGetter) Int16Slice(name string) []int16 {
	v, _ := g.LookupInt16Slice(name)
	return v
}
