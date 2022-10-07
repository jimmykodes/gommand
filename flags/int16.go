package flags

import (
    "strconv"
)

type int16Flag struct {
    *baseFlag

	defValue int16
	value    int16
}

func (f *int16Flag) Type() FlagType { return Int16FlagType }

func (f *int16Flag) Value() any {
    if f.IsSet() {
        return f.value
    }
    return f.defValue
}

func (f *int16Flag) Set(s string) error {
    v, err := strconv.ParseInt(s, 0, 16)
    if err != nil {
        return err
    }
    f.value = int16(v)
    f.set = true
    return nil
}

func Int16Flag(name string, value int16, usage string) Flag {
    return &int16Flag{
        baseFlag: &baseFlag{name: name, usage: usage},
        defValue: value,
    }
}

func Int16FlagS(name string, shorthand rune, value int16, usage string) Flag {
    return &int16Flag{
        baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
        defValue: value,
    }
}

func (fs *FlagSet) Int16(name string, value int16, usage string) {
    fs.addFlag(Int16Flag(name, value, usage))
}

func (fs *FlagSet) Int16S(name string, shorthand rune, value int16, usage string) {
    fs.addFlag(Int16FlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) int16Val(name string) (int16, error) {
    f, err := fs.flag(name, Int16FlagType)
    if err != nil {
        return 0, err
    }
    return f.Value().(int16), nil
}

func (g FlagGetter) LookupInt16(name string) (int16, error) {
    return g.fs.int16Val(name)
}

func (g FlagGetter) Int16(name string) int16 {
    v, _ := g.LookupInt16(name)
    return v
}
