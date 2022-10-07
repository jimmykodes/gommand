package flags

import (
    "strconv"
)

type boolFlag struct {
    *baseFlag

	defValue bool
	value    bool
}

func (f *boolFlag) Type() FlagType { return BoolFlagType }

func (f *boolFlag) Value() any {
    if f.IsSet() {
        return f.value
    }
    return f.defValue
}

func (f *boolFlag) Set(s string) error {
    v, err := strconv.ParseBool(s)
    if err != nil {
        return err
    }
    f.value = v
    f.set = true
    return nil
}

func BoolFlag(name string, value bool, usage string) Flag {
    return &boolFlag{
        baseFlag: &baseFlag{name: name, usage: usage},
        defValue: value,
    }
}

func BoolFlagS(name string, shorthand rune, value bool, usage string) Flag {
    return &boolFlag{
        baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
        defValue: value,
    }
}

func (fs *FlagSet) Bool(name string, value bool, usage string) {
    fs.addFlag(BoolFlag(name, value, usage))
}

func (fs *FlagSet) BoolS(name string, shorthand rune, value bool, usage string) {
    fs.addFlag(BoolFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) boolVal(name string) (bool, error) {
    f, err := fs.flag(name, BoolFlagType)
    if err != nil {
        return false, err
    }
    return f.Value().(bool), nil
}

func (g FlagGetter) LookupBool(name string) (bool, error) {
    return g.fs.boolVal(name)
}

func (g FlagGetter) Bool(name string) bool {
    v, _ := g.LookupBool(name)
    return v
}
