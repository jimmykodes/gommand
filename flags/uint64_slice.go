package flags

import (
	"strconv"
	"strings"
)

type uint64SliceFlag struct {
	*baseFlag

	defValue []uint64
	value    []uint64
}

func (f *uint64SliceFlag) Type() FlagType { return Uint64SliceFlagType }

func (f *uint64SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint64SliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]uint64, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 64)
		if err != nil {
			return err
		}
		v[i] = uint64(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Uint64SliceFlag(name string, value []uint64, usage string) Flag {
	return &uint64SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint64SliceFlagS(name string, shorthand rune, value []uint64, usage string) Flag {
	return &uint64SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint64Slice(name string, value []uint64, usage string) {
	fs.addFlag(Uint64SliceFlag(name, value, usage))
}

func (fs *FlagSet) Uint64SliceS(name string, shorthand rune, value []uint64, usage string) {
	fs.addFlag(Uint64SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint64SliceVal(name string) ([]uint64, error) {
	f, err := fs.flag(name, Uint64SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]uint64), nil
}

func (g FlagGetter) LookupUint64Slice(name string) ([]uint64, error) {
	return g.fs.uint64SliceVal(name)
}

func (g FlagGetter) Uint64Slice(name string) []uint64 {
	v, _ := g.LookupUint64Slice(name)
	return v
}
