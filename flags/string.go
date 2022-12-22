package flags

type stringFlag struct {
	*baseFlag

	defValue string
	value    string
}

func (f *stringFlag) Type() FlagType { return StringFlagType }

func (f *stringFlag) Value() any {
	if f.IsSet() {
		return f.value
	}
	return f.defValue
}

func (f *stringFlag) Set(s string) error {
	f.value = s
	f.set = true
	return nil
}

func StringFlag(name string, value string, usage string) Flag {
	return &stringFlag{
		baseFlag: &baseFlag{name: name, usage: usage},
		defValue: value,
	}
}

func StringFlagS(name string, shorthand rune, value string, usage string) Flag {
	return &stringFlag{
		baseFlag: &baseFlag{name: name, short: shorthand, usage: usage},
		defValue: value,
	}
}

func (fs *FlagSet) String(name string, value string, usage string) {
	fs.AddFlag(StringFlag(name, value, usage))
}

func (fs *FlagSet) StringS(name string, shorthand rune, value string, usage string) {
	fs.AddFlag(StringFlagS(name, shorthand, value, usage))
}

func (fs *FlagSet) stringVal(name string) (string, error) {
	f, err := fs.flag(name, StringFlagType)
	if err != nil {
		return "", err
	}
	return f.Value().(string), nil
}

func (g FlagGetter) LookupString(name string) (string, error) {
	return g.fs.stringVal(name)
}

func (g FlagGetter) String(name string) string {
	v, _ := g.LookupString(name)
	return v
}
