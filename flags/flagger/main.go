package main

import (
	_ "embed"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/jimmykodes/strman"
)

var (
	//go:embed flag.gotmpl
	tmplStr string

	funcMap = template.FuncMap{"pascal": strman.ToPascal, "camel": strman.ToCamel, "default": func(a, b string) string {
		if a != "" {
			return a
		}
		return b
	}}
	tmpl = template.Must(template.New("intFlag").Funcs(funcMap).Parse(tmplStr))
)

func intParseFunc(size string) string {
	s := "strconv.ParseInt(s, 0, "
	if size == "" {
		s += "64"
	} else {
		s += size
	}
	return s + ")"
}

var baseImport = []string{"strconv"}

func main() {
	for _, t := range []string{"int", "uint"} {
		for _, s := range []string{"", "8", "16", "32", "64"} {
			if err := createFlag(t+s, intParseFunc(s), "0", baseImport); err != nil {
				panic(err)
			}
		}
	}
	for _, s := range []string{"32", "64"} {
		if err := createFlag("float"+s, "strconv.ParseFloat(s, "+s+")", "0", baseImport); err != nil {
			panic(err)
		}
	}
	if err := createFlag("bool", "strconv.ParseBool(s)", "false", baseImport); err != nil {
		panic(err)
	}
	if err := createFlag("string", "", `""`, baseImport); err != nil {
		panic(err)
	}
	if err := createFlag("time.Duration", "time.ParseDuration(s)", `time.Duration(0)`, []string{"time"}); err != nil {
		panic(err)
	}
}

type data struct {
	Type      string
	Name      string
	ParseFunc string
	Default   string
	Imports   []string
}

func createFlag(t, parseFunc, defaultVal string, imports []string) error {
	sort.Strings(imports)
	name := strings.ToLower(t[strings.LastIndex(t, ".")+1:])
	f, err := os.Create(fmt.Sprintf("../%s.go", name))
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data{Type: t, Name: name, ParseFunc: parseFunc, Default: defaultVal, Imports: imports})
}
