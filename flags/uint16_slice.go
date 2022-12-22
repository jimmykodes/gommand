package flags

import (
	"strconv"
	"strings"
)

type uint16SliceFlag struct {
	*baseFlag

	defValue []uint16
	value    []uint16
}

func (f *uint16SliceFlag) Type() FlagType { return Uint16SliceFlagType }

func (f *uint16SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint16SliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]uint16, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 16)
		if err != nil {
			return err
		}
		v[i] = uint16(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Uint16SliceFlag(name string, value []uint16, usage string) Flag {
	return &uint16SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint16SliceFlagS(name string, shorthand rune, value []uint16, usage string) Flag {
	return &uint16SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint16Slice(name string, value []uint16, usage string) {
	fs.AddFlag(Uint16SliceFlag(name, value, usage))
}

func (fs *FlagSet) Uint16SliceS(name string, shorthand rune, value []uint16, usage string) {
	fs.AddFlag(Uint16SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint16SliceVal(name string) ([]uint16, error) {
	f, err := fs.flag(name, Uint16SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]uint16), nil
}

func (g FlagGetter) LookupUint16Slice(name string) ([]uint16, error) {
	return g.fs.uint16SliceVal(name)
}

func (g FlagGetter) Uint16Slice(name string) []uint16 {
	v, _ := g.LookupUint16Slice(name)
	return v
}
