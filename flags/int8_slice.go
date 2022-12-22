package flags

import (
	"strconv"
	"strings"
)

type int8SliceFlag struct {
	*baseFlag

	defValue []int8
	value    []int8
}

func (f *int8SliceFlag) Type() FlagType { return Int8SliceFlagType }

func (f *int8SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int8SliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]int8, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 8)
		if err != nil {
			return err
		}
		v[i] = int8(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Int8SliceFlag(name string, value []int8, usage string) Flag {
	return &int8SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int8SliceFlagS(name string, shorthand rune, value []int8, usage string) Flag {
	return &int8SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int8Slice(name string, value []int8, usage string) {
	fs.AddFlag(Int8SliceFlag(name, value, usage))
}

func (fs *FlagSet) Int8SliceS(name string, shorthand rune, value []int8, usage string) {
	fs.AddFlag(Int8SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int8SliceVal(name string) ([]int8, error) {
	f, err := fs.flag(name, Int8SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]int8), nil
}

func (g FlagGetter) LookupInt8Slice(name string) ([]int8, error) {
	return g.fs.int8SliceVal(name)
}

func (g FlagGetter) Int8Slice(name string) []int8 {
	v, _ := g.LookupInt8Slice(name)
	return v
}
