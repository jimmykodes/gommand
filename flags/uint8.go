package flags

import (
    "strconv"
)

type uint8Flag struct {
    *baseFlag

	defValue uint8
	value    uint8
}

func (f *uint8Flag) Type() FlagType { return Uint8FlagType }

func (f *uint8Flag) Value() any {
    if f.IsSet() {
        return f.value
    }
    return f.defValue
}

func (f *uint8Flag) Set(s string) error {
    v, err := strconv.ParseInt(s, 0, 8)
    if err != nil {
        return err
    }
    f.value = uint8(v)
    f.set = true
    return nil
}

func Uint8Flag(name string, value uint8, usage string) Flag {
    return &uint8Flag{
        baseFlag: &baseFlag{name: name, usage: usage},
        defValue: value,
    }
}

func Uint8FlagS(name string, shorthand rune, value uint8, usage string) Flag {
    return &uint8Flag{
        baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
        defValue: value,
    }
}

func (fs *FlagSet) Uint8(name string, value uint8, usage string) {
    fs.addFlag(Uint8Flag(name, value, usage))
}

func (fs *FlagSet) Uint8S(name string, shorthand rune, value uint8, usage string) {
    fs.addFlag(Uint8FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) uint8Val(name string) (uint8, error) {
    f, err := fs.flag(name, Uint8FlagType)
    if err != nil {
        return 0, err
    }
    return f.Value().(uint8), nil
}

func (g FlagGetter) LookupUint8(name string) (uint8, error) {
    return g.fs.uint8Val(name)
}

func (g FlagGetter) Uint8(name string) uint8 {
    v, _ := g.LookupUint8(name)
    return v
}
