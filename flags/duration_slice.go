package flags

import (
	"strings"
	"time"
)

type durationSliceFlag struct {
	*baseFlag

	defValue []time.Duration
	value    []time.Duration
}

func (f *durationSliceFlag) Type() FlagType { return DurationSliceFlagType }

func (f *durationSliceFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *durationSliceFlag) Set(s string) error {
	var (
		pieces = strings.Split(s, sliceSeparator)
		v      = make([]time.Duration, len(pieces))
	)
	for i, piece := range pieces {
		val, err := time.ParseDuration(piece)
		if err != nil {
			return err
		}
		v[i] = val
	}
	f.value = v
	f.set = true
	return nil
}

func DurationSliceFlag(name string, value []time.Duration, usage string) Flag {
	return &durationSliceFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func DurationSliceFlagS(name string, shorthand rune, value []time.Duration, usage string) Flag {
	return &durationSliceFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) DurationSlice(name string, value []time.Duration, usage string) {
	fs.addFlag(DurationSliceFlag(name, value, usage))
}

func (fs *FlagSet) DurationSliceS(name string, shorthand rune, value []time.Duration, usage string) {
	fs.addFlag(DurationSliceFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) durationSliceVal(name string) ([]time.Duration, error) {
	f, err := fs.flag(name, DurationSliceFlagType)
	if err != nil {
		return nil, err
	}
	return f.Value().([]time.Duration), nil
}

func (g FlagGetter) LookupDurationSlice(name string) ([]time.Duration, error) {
	return g.fs.durationSliceVal(name)
}

func (g FlagGetter) DurationSlice(name string) []time.Duration {
	v, _ := g.LookupDurationSlice(name)
	return v
}
