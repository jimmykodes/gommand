package flags

import (
	"strconv"
	"strings"
)

type int32SliceFlag struct {
	*baseFlag

	defValue []int32
	value    []int32
}

func (f *int32SliceFlag) Type() FlagType { return Int32SliceFlagType }

func (f *int32SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int32SliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]int32, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 32)
		if err != nil {
			return err
		}
		v[i] = int32(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Int32SliceFlag(name string, value []int32, usage string) Flag {
	return &int32SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int32SliceFlagS(name string, shorthand rune, value []int32, usage string) Flag {
	return &int32SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int32Slice(name string, value []int32, usage string) {
	fs.addFlag(Int32SliceFlag(name, value, usage))
}

func (fs *FlagSet) Int32SliceS(name string, shorthand rune, value []int32, usage string) {
	fs.addFlag(Int32SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int32SliceVal(name string) ([]int32, error) {
	f, err := fs.flag(name, Int32SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]int32), nil
}

func (g FlagGetter) LookupInt32Slice(name string) ([]int32, error) {
	return g.fs.int32SliceVal(name)
}

func (g FlagGetter) Int32Slice(name string) []int32 {
	v, _ := g.LookupInt32Slice(name)
	return v
}
