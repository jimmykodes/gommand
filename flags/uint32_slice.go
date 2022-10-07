package flags

import (
	"strconv"
	"strings"
)

type uint32SliceFlag struct {
	*baseFlag

	defValue []uint32
	value    []uint32
}

func (f *uint32SliceFlag) Type() FlagType { return Uint32SliceFlagType }

func (f *uint32SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint32SliceFlag) Set(s string) error {
	var (
		pieces = strings.Split(s, sliceSeparator)
		v      = make([]uint32, len(pieces))
	)
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 32)
		if err != nil {
			return err
		}
		v[i] = uint32(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Uint32SliceFlag(name string, value []uint32, usage string) Flag {
	return &uint32SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint32SliceFlagS(name string, shorthand rune, value []uint32, usage string) Flag {
	return &uint32SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint32Slice(name string, value []uint32, usage string) {
	fs.addFlag(Uint32SliceFlag(name, value, usage))
}

func (fs *FlagSet) Uint32SliceS(name string, shorthand rune, value []uint32, usage string) {
	fs.addFlag(Uint32SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint32SliceVal(name string) ([]uint32, error) {
	f, err := fs.flag(name, Uint32SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]uint32), nil
}

func (g FlagGetter) LookupUint32Slice(name string) ([]uint32, error) {
	return g.fs.uint32SliceVal(name)
}

func (g FlagGetter) Uint32Slice(name string) []uint32 {
	v, _ := g.LookupUint32Slice(name)
	return v
}
