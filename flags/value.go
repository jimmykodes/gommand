package flags

import (
	"os"

	"github.com/jimmykodes/strman"
)

type Valuer interface {
	Value(name string) (string, bool)
}

type ValuerFunc func(string) (string, bool)

func (vf ValuerFunc) Value(s string) (string, bool) {
	return vf(s)
}

var (
	Environ Valuer = &environ{}
)

type environ struct{}

func (e *environ) Value(name string) (string, bool) {
	return os.LookupEnv(strman.ToScreamingSnake(name))
}
