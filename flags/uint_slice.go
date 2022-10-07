package flags

import (
	"strconv"
	"strings"
)

type uintSliceFlag struct {
	*baseFlag

	defValue []uint
	value    []uint
}

func (f *uintSliceFlag) Type() FlagType { return UintSliceFlagType }

func (f *uintSliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uintSliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]uint, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 64)
		if err != nil {
			return err
		}
		v[i] = uint(val)
	}
	f.value = v
	f.set = true
	return nil
}

func UintSliceFlag(name string, value []uint, usage string) Flag {
	return &uintSliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func UintSliceFlagS(name string, shorthand rune, value []uint, usage string) Flag {
	return &uintSliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) UintSlice(name string, value []uint, usage string) {
	fs.addFlag(UintSliceFlag(name, value, usage))
}

func (fs *FlagSet) UintSliceS(name string, shorthand rune, value []uint, usage string) {
	fs.addFlag(UintSliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uintSliceVal(name string) ([]uint, error) {
	f, err := fs.flag(name, UintSliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]uint), nil
}

func (g FlagGetter) LookupUintSlice(name string) ([]uint, error) {
	return g.fs.uintSliceVal(name)
}

func (g FlagGetter) UintSlice(name string) []uint {
	v, _ := g.LookupUintSlice(name)
	return v
}
