package flags

import (
	"strconv"
)

type float64Flag struct {
	*baseFlag

	defValue float64
	value    float64
}

func (f *float64Flag) Type() FlagType { return Float64FlagType }

func (f *float64Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *float64Flag) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	f.value = v
	f.set = true
	return nil
}

func Float64Flag(name string, value float64, usage string) Flag {
	return &float64Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Float64FlagS(name string, shorthand rune, value float64, usage string) Flag {
	return &float64Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Float64(name string, value float64, usage string) {
	fs.addFlag(Float64Flag(name, value, usage))
}

func (fs *FlagSet) Float64S(name string, shorthand rune, value float64, usage string) {
	fs.addFlag(Float64FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) float64Val(name string) (float64, error) {
	f, err := fs.flag(name, Float64FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(float64), nil
}

func (g FlagGetter) LookupFloat64(name string) (float64, error) {
	return g.fs.float64Val(name)
}

func (g FlagGetter) Float64(name string) float64 {
	v, _ := g.LookupFloat64(name)
	return v
}