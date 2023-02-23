package flags

type FlagSetOption interface {
	Apply(*FlagSet)
}

var _ FlagSetOption = FlagSetOptionFunc(nil)

type FlagSetOptionFunc func(*FlagSet)

func (f FlagSetOptionFunc) Apply(s *FlagSet) { f(s) }

func WithHelpFlag() FlagSetOptionFunc {
	return func(s *FlagSet) {
		s.addHelpFlag()
	}
}
