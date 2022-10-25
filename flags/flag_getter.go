package flags

func NewFlagGetter(fs *FlagSet) *FlagGetter {
	return &FlagGetter{fs: fs}
}

type FlagGetter struct {
	fs *FlagSet
}

func (g FlagGetter) Flag(name string) Flag {
	return g.fs.FromName(name)
}
