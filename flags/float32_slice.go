package flags

import (
	"strconv"
	"strings"
)

type float32SliceFlag struct {
	*baseFlag

	defValue []float32
	value    []float32
}

func (f *float32SliceFlag) Type() FlagType { return Float32SliceFlagType }

func (f *float32SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *float32SliceFlag) Set(s string) error {
	var (
		pieces = strings.Split(s, sliceSeparator)
		v      = make([]float32, len(pieces))
	)
	for i, piece := range pieces {
		val, err := strconv.ParseFloat(piece, 32)
		if err != nil {
			return err
		}
		v[i] = float32(val)
	}
	f.value = v
	f.set = true
	return nil
}

func Float32SliceFlag(name string, value []float32, usage string) Flag {
	return &float32SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Float32SliceFlagS(name string, shorthand rune, value []float32, usage string) Flag {
	return &float32SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Float32Slice(name string, value []float32, usage string) {
	fs.addFlag(Float32SliceFlag(name, value, usage))
}

func (fs *FlagSet) Float32SliceS(name string, shorthand rune, value []float32, usage string) {
	fs.addFlag(Float32SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) float32SliceVal(name string) ([]float32, error) {
	f, err := fs.flag(name, Float32SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]float32), nil
}

func (g FlagGetter) LookupFloat32Slice(name string) ([]float32, error) {
	return g.fs.float32SliceVal(name)
}

func (g FlagGetter) Float32Slice(name string) []float32 {
	v, _ := g.LookupFloat32Slice(name)
	return v
}
