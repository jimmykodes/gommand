package flags

import (
	"strconv"
	"strings"
)

type uint8SliceFlag struct {
	*baseFlag

	defValue []uint8
	value    []uint8
}

func (f *uint8SliceFlag) Type() FlagType { return Uint8SliceFlagType }

func (f *uint8SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *uint8SliceFlag) Set(s string) error {
	var (
		pieces = strings.Split(s, sliceSeparator)
		v      = make([]uint8, len(pieces))
	)
	for i, piece := range pieces {
		val, err := strconv.ParseInt(piece, 0, 8)
		if err != nil {
			return err
		}
		v[i] = uint8(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Uint8SliceFlag(name string, value []uint8, usage string) Flag {
	return &uint8SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Uint8SliceFlagS(name string, shorthand rune, value []uint8, usage string) Flag {
	return &uint8SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Uint8Slice(name string, value []uint8, usage string) {
	fs.addFlag(Uint8SliceFlag(name, value, usage))
}

func (fs *FlagSet) Uint8SliceS(name string, shorthand rune, value []uint8, usage string) {
	fs.addFlag(Uint8SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint8SliceVal(name string) ([]uint8, error) {
	f, err := fs.flag(name, Uint8SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]uint8), nil
}

func (g FlagGetter) LookupUint8Slice(name string) ([]uint8, error) {
	return g.fs.uint8SliceVal(name)
}

func (g FlagGetter) Uint8Slice(name string) []uint8 {
	v, _ := g.LookupUint8Slice(name)
	return v
}
