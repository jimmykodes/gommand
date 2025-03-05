package flags

import "iter"

func NewFlagGetter(fs *FlagSet) *FlagGetter {
	return &FlagGetter{fs: fs}
}

type FlagGetter struct {
	fs *FlagSet
}

func (g FlagGetter) Flag(name string) Flag {
	return g.fs.FromName(name)
}

func (g FlagGetter) All() iter.Seq2[string, Flag] {
	return func(yield func(string, Flag) bool) {
		for k, v := range g.fs.flags {
			if !yield(k, v) {
				return
			}
		}
	}
}
