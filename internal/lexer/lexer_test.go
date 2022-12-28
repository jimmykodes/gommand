package lexer_test

import (
	"reflect"
	"testing"

	"github.com/jimmykodes/gommand/internal/lexer"
)

func TestLexer(t *testing.T) {
	l := lexer.New([]string{
		"-f",
		"-lvg",
		"--port",
		"8080",
		"--host=0.0.0.0",
		"-t=taco",
		"-iec=false",
	})
	expected := []*lexer.Token{
		{
			Type: lexer.ShortFlagType,
			Name: "f",
		},
		{
			Type: lexer.MultiFlagType,
			Name: "lvg",
		},
		{
			Type: lexer.LongFlagType,
			Name: "port",
		},
		{
			Type:  lexer.ValueType,
			Value: "8080",
		},
		{
			Type:  lexer.LongFlagType,
			Name:  "host",
			Value: "0.0.0.0",
		},
		{
			Type:  lexer.ShortFlagType,
			Name:  "t",
			Value: "taco",
		},
		{
			Type:  lexer.MultiFlagType,
			Name:  "iec",
			Value: "false",
		},
	}

	for _, e := range expected {
		if token := l.Read(); !reflect.DeepEqual(token, e) {
			t.Errorf("invalid token. got %+v - want %+v", token, e)
		}
	}
}
