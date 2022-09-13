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

func (g FlagGetter) LookupFloat64(name string) (float64, error) {
	return g.fs.float64Val(name)
}

func (g FlagGetter) Float64(name string) float64 {
	v, _ := g.LookupFloat64(name)
	return v
}

func (g FlagGetter) LookupString(name string) (string, error) {
	return g.fs.stringVal(name)
}

func (g FlagGetter) String(name string) string {
	v, _ := g.LookupString(name)
	return v
}
