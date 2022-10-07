package flags

import (
	"strconv"
	"strings"
)

type int64SliceFlag struct {
	*baseFlag

	defValue []int64
	value    []int64
}

func (f *int64SliceFlag) Type() FlagType { return Int64SliceFlagType }

func (f *int64SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *int64SliceFlag) Set(s string) error {
	var (
		pieces = strings.Split(s, sliceSeparator)
		v      = make([]int64, len(pieces))
	)
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 64)
		if err != nil {
			return err
		}
		v[i] = val
	}
	f.value = v
	f.set = true
	return nil
}

func Int64SliceFlag(name string, value []int64, usage string) Flag {
	return &int64SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Int64SliceFlagS(name string, shorthand rune, value []int64, usage string) Flag {
	return &int64SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Int64Slice(name string, value []int64, usage string) {
	fs.addFlag(Int64SliceFlag(name, value, usage))
}

func (fs *FlagSet) Int64SliceS(name string, shorthand rune, value []int64, usage string) {
	fs.addFlag(Int64SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int64SliceVal(name string) ([]int64, error) {
	f, err := fs.flag(name, Int64SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]int64), nil
}

func (g FlagGetter) LookupInt64Slice(name string) ([]int64, error) {
	return g.fs.int64SliceVal(name)
}

func (g FlagGetter) Int64Slice(name string) []int64 {
	v, _ := g.LookupInt64Slice(name)
	return v
}
