package gommand

type ArgValidator func(s []string) bool

func ArgsExact(n int) ArgValidator {
	return func(s []string) bool {
		return len(s) == n
	}
}

func ArgsNone() ArgValidator {
	return ArgsExact(0)
}

func ArgsMin(bound int) ArgValidator {
	return func(s []string) bool {
		return len(s) >= bound
	}
}

func ArgsMax(bound int) ArgValidator {
	return func(s []string) bool {
		return len(s) <= bound
	}
}

func ArgsBetween(lower, upper int) ArgValidator {
	return func(s []string) bool {
		return lower <= len(s) && len(s) <= upper
	}
}

func ArgsEvery(validators ...ArgValidator) ArgValidator {
	return func(s []string) bool {
		for _, v := range validators {
			if !v(s) {
				return false
			}
		}
		return true
	}
}

func ArgsSome(validators ...ArgValidator) ArgValidator {
	return func(s []string) bool {
		for _, v := range validators {
			if v(s) {
				return true
			}
		}
		return false
	}
}

func ArgsAny() ArgValidator {
	return func([]string) bool {
		return true
	}
}
