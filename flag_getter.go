package gommand

type FlagGetter struct {
	fs *FlagSet
}

func (g FlagGetter) LookupInt(name string) (int, error) {
	return g.fs.intVal(name)
}

func (g FlagGetter) Int(name string) int {
	v, _ := g.LookupInt(name)
	return v
}

func (g FlagGetter) LookupBool(name string) (bool, error) {
	return g.fs.boolVal(name)
}

func (g FlagGetter) Bool(name string) bool {
	v, _ := g.LookupBool(name)
	return v
}
