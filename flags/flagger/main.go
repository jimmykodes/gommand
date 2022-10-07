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
	//go:embed sliceFlag.gotmpl
	sliceTmplStr string

	funcMap = template.FuncMap{"pascal": strman.ToPascal, "camel": strman.ToCamel, "default": func(a, b string) string {
		if a != "" {
			return a
		}
		return b
	}}
	tmpl      = template.Must(template.New("flag").Funcs(funcMap).Parse(tmplStr))
	sliceTmpl = template.Must(template.New("sliceFlag").Funcs(funcMap).Parse(sliceTmplStr))
)

func intParseFunc(size string) string {
	s := "strconv.ParseInt(%s, 0, "
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
			convertParsed := t != "int" || s != "64"
			if err := createFlag(t+s, intParseFunc(s), "0", convertParsed, baseImport); err != nil {
				panic(err)
			}
		}
	}
	for _, s := range []string{"32", "64"} {
		if err := createFlag("float"+s, "strconv.ParseFloat(%s, "+s+")", "0", s != "64", baseImport); err != nil {
			panic(err)
		}
	}
	if err := createFlag("bool", "strconv.ParseBool(%s)", "false", false, baseImport); err != nil {
		panic(err)
	}
	if err := createFlag("string", "", `""`, false, nil); err != nil {
		panic(err)
	}
	if err := createFlag("time.Duration", "time.ParseDuration(%s)", `time.Duration(0)`, false, []string{"time"}); err != nil {
		panic(err)
	}
}

type data struct {
	Type          string
	Name          string
	ParseFunc     string
	Default       string
	Imports       []string
	ConvertParsed bool
}

func createFlag(t, parseFunc, defaultVal string, convert bool, imports []string) error {
	name := strings.ToLower(t[strings.LastIndex(t, ".")+1:])
	f, err := os.Create(fmt.Sprintf("../%s.go", name))
	if err != nil {
		return err
	}
	defer f.Close()
	sort.Strings(imports)
	if err := tmpl.Execute(f, data{
		Type:          t,
		Name:          name,
		ParseFunc:     singleParseFunc(parseFunc),
		ConvertParsed: convert,
		Default:       defaultVal,
		Imports:       imports,
	}); err != nil {
		return err
	}
	imports = append(imports, "strings")
	sliceF, err := os.Create(fmt.Sprintf("../%s_slice.go", name))
	if err != nil {
		return err
	}
	defer sliceF.Close()
	sort.Strings(imports)
	return sliceTmpl.Execute(sliceF, data{
		Type:          t,
		Name:          name,
		ParseFunc:     sliceParseFunc(parseFunc),
		ConvertParsed: convert,
		Default:       defaultVal,
		Imports:       imports,
	})
}

func singleParseFunc(f string) string {
	if f != "" {
		return fmt.Sprintf(f, "s")
	}
	return f
}
func sliceParseFunc(f string) string {
	if f != "" {
		return fmt.Sprintf(f, "piece")
	}
	return f
}
