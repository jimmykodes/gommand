package gommand

import "fmt"

type ArgValidator func(s []string) error

func ArgsExact(n int) ArgValidator {
	return func(s []string) error {
		if len(s) != n {
			return fmt.Errorf("expected exactly %d arguments, got %d", n, len(s))
		}
		return nil
	}
}

func ArgsNone() ArgValidator {
	return ArgsExact(0)
}

func ArgsMin(bound int) ArgValidator {
	return func(s []string) error {
		if len(s) < bound {
			return fmt.Errorf("expected at least %d arguments, got %d", bound, len(s))
		}
		return nil
	}
}

func ArgsMax(bound int) ArgValidator {
	return func(s []string) error {
		if len(s) > bound {
			return fmt.Errorf("expected at most %d arguments, got %d", bound, len(s))
		}
		return nil
	}
}

func ArgsBetween(lower, upper int) ArgValidator {
	return func(s []string) error {
		if len(s) < lower || len(s) > upper {
			return fmt.Errorf("expected between %d and %d arguments, got %d", lower, upper, len(s))
		}
		return nil
	}
}

func ArgsEvery(validators ...ArgValidator) ArgValidator {
	return func(s []string) error {
		for _, v := range validators {
			if err := v(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func ArgsSome(validators ...ArgValidator) ArgValidator {
	return func(s []string) error {
		var errs []error
		for _, v := range validators {
			if err := v(s); err == nil {
				return nil
			} else {
				errs = append(errs, err)
			}
		}
		return fmt.Errorf("no validators passed: %v", errs)
	}
}

func ArgsAny() ArgValidator {
	return func([]string) error {
		return nil
	}
}
