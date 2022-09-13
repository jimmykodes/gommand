package gommand

type FlagSetOption interface {
	Apply(*FlagSet)
}

var _ FlagSetOption = FlagSetOptionFunc(nil)

type FlagSetOptionFunc func(*FlagSet)

func (f FlagSetOptionFunc) Apply(s *FlagSet) { f(s) }

func WithEnvPrefix(prefix string) FlagSetOptionFunc {
	return func(s *FlagSet) {
		s.envPrefix = prefix
	}
}

func WithNoEnv() FlagSetOptionFunc {
	return func(s *FlagSet) {
		s.noEnv = true
	}
}
