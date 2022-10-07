package flags

import (
	"strconv"
	"strings"
)

type intSliceFlag struct {
	*baseFlag

	defValue []int
	value    []int
}

func (f *intSliceFlag) Type() FlagType { return IntSliceFlagType }

func (f *intSliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *intSliceFlag) Set(s string) error {
	var (
		pieces = strings.Split(s, sliceSeparator)
		v      = make([]int, len(pieces))
	)
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 64)
		if err != nil {
			return err
		}
		v[i] = int(val)
	}
	f.value = v
	f.set = true
	return nil
}

func IntSliceFlag(name string, value []int, usage string) Flag {
	return &intSliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func IntSliceFlagS(name string, shorthand rune, value []int, usage string) Flag {
	return &intSliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) IntSlice(name string, value []int, usage string) {
	fs.addFlag(IntSliceFlag(name, value, usage))
}

func (fs *FlagSet) IntSliceS(name string, shorthand rune, value []int, usage string) {
	fs.addFlag(IntSliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) intSliceVal(name string) ([]int, error) {
	f, err := fs.flag(name, IntSliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]int), nil
}

func (g FlagGetter) LookupIntSlice(name string) ([]int, error) {
	return g.fs.intSliceVal(name)
}

func (g FlagGetter) IntSlice(name string) []int {
	v, _ := g.LookupIntSlice(name)
	return v
}
