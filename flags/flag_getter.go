package flags

type FlagGetter struct {
	fs *FlagSet
}

func NewFlagGetter(fs *FlagSet) *FlagGetter {
	return &FlagGetter{fs: fs}
}
