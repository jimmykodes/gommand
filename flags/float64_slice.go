package flags

import (
	"strconv"
	"strings"
)

type float64SliceFlag struct {
	*baseFlag

	defValue []float64
	value    []float64
}

func (f *float64SliceFlag) Type() FlagType { return Float64SliceFlagType }

func (f *float64SliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *float64SliceFlag) Set(s string) error {
	pieces := strings.Split(s, sliceSeparator)
	v := make([]float64, len(pieces))
	for i, piece := range pieces {
		val, err := strconv.ParseFloat(piece, 64)
		if err != nil {
			return err
		}
		v[i] = val
	}
	f.value = v
	f.set = true
	return nil
}

func Float64SliceFlag(name string, value []float64, usage string) Flag {
	return &float64SliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Float64SliceFlagS(name string, shorthand rune, value []float64, usage string) Flag {
	return &float64SliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Float64Slice(name string, value []float64, usage string) {
	fs.addFlag(Float64SliceFlag(name, value, usage))
}

func (fs *FlagSet) Float64SliceS(name string, shorthand rune, value []float64, usage string) {
	fs.addFlag(Float64SliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) float64SliceVal(name string) ([]float64, error) {
	f, err := fs.flag(name, Float64SliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]float64), nil
}

func (g FlagGetter) LookupFloat64Slice(name string) ([]float64, error) {
	return g.fs.float64SliceVal(name)
}

func (g FlagGetter) Float64Slice(name string) []float64 {
	v, _ := g.LookupFloat64Slice(name)
	return v
}
