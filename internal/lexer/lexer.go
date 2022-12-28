package lexer

import "strings"

func New(strs []string) *Lexer {
	return &Lexer{strs: strs}
}

type TokenType int

const (
	UnknownType TokenType = iota
	LongFlagType
	ShortFlagType
	MultiFlagType
	ValueType
)

func (t TokenType) String() string {
	return [...]string{
		"UnknownType",
		"LongFlagType",
		"ShortFlagType",
		"MultiFlagType",
		"ValueType",
	}[t]
}

type Token struct {
	Type  TokenType
	Name  string
	Value string
}

type Lexer struct {
	strs []string
	pos  int
}

func (l *Lexer) Read() *Token {
	t := l.Peek()
	if t != nil {
		l.pos++
	}
	return t
}

func (l *Lexer) Peek() *Token {
	if l.pos >= len(l.strs) {
		return nil
	}
	flag, value, _ := strings.Cut(l.strs[l.pos], "=")
	if strings.HasPrefix(flag, "--") {
		return &Token{
			Type:  LongFlagType,
			Name:  strings.TrimPrefix(flag, "--"),
			Value: value,
		}
	}
	if strings.HasPrefix(flag, "-") {
		if len(flag) == 2 {
			return &Token{
				Type:  ShortFlagType,
				Name:  strings.TrimPrefix(flag, "-"),
				Value: value,
			}
		}
		return &Token{
			Type:  MultiFlagType,
			Name:  strings.TrimPrefix(flag, "-"),
			Value: value,
		}
	}

	return &Token{
		Type:  ValueType,
		Value: l.strs[l.pos],
	}
}
