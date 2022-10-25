package flags

import (
	"time"
)

type durationFlag struct {
	*baseFlag

	defValue time.Duration
	value    time.Duration
}

func (f *durationFlag) Type() FlagType { return DurationFlagType }

func (f *durationFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *durationFlag) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	f.value = v
	f.set = true
	return nil
}

func DurationFlag(name string, value time.Duration, usage string) Flag {
	return &durationFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func DurationFlagS(name string, shorthand rune, value time.Duration, usage string) Flag {
	return &durationFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) {
	fs.addFlag(DurationFlag(name, value, usage))
}

func (fs *FlagSet) DurationS(name string, shorthand rune, value time.Duration, usage string) {
	fs.addFlag(DurationFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) durationVal(name string) (time.Duration, error) {
	f, err := fs.flag(name, DurationFlagType)
	if err != nil {
		return time.Duration(0), err
	}
	return f.Value().(time.Duration), nil
}

func (g FlagGetter) LookupDuration(name string) (time.Duration, error) {
	return g.fs.durationVal(name)
}

func (g FlagGetter) Duration(name string) time.Duration {
	v, _ := g.LookupDuration(name)
	return v
}
