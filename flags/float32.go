package flags

import (
	"strconv"
)

type float32Flag struct {
	*baseFlag

	defValue float32
	value    float32
}

func (f *float32Flag) Type() FlagType { return Float32FlagType }

func (f *float32Flag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *float32Flag) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	f.value = float32(v)
	f.set = true
	return nil
}

func Float32Flag(name string, value float32, usage string) Flag {
	return &float32Flag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func Float32FlagS(name string, shorthand rune, value float32, usage string) Flag {
	return &float32Flag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Float32(name string, value float32, usage string) {
	fs.AddFlag(Float32Flag(name, value, usage))
}

func (fs *FlagSet) Float32S(name string, shorthand rune, value float32, usage string) {
	fs.AddFlag(Float32FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) float32Val(name string) (float32, error) {
	f, err := fs.flag(name, Float32FlagType)
	if err != nil {
		return 0, err
	}
	return f.Value().(float32), nil
}

func (g FlagGetter) LookupFloat32(name string) (float32, error) {
	return g.fs.float32Val(name)
}

func (g FlagGetter) Float32(name string) float32 {
	v, _ := g.LookupFloat32(name)
	return v
}
