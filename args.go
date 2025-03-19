package gommand

import (
	"iter"
	"strconv"
)

type Args []string

func (a Args) String(idx int) string {
	if idx >= len(a) {
		return ""
	}
	return a[idx]
}

func (a Args) Bool(idx int) (bool, error) {
	s := a.String(idx)
	if s == "" {
		return false, nil
	}
	return strconv.ParseBool(s)
}

func (a Args) Int64(idx int) (int64, error) {
	s := a.String(idx)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func (a Args) Int(idx int) (int, error) {
	i, err := a.Int64(idx)
	return int(i), err
}

func (a Args) Float64(idx int) (float64, error) {
	s := a.String(idx)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func (a Args) Float32(idx int) (float32, error) {
	s := a.String(idx)
	if s == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

func (a Args) Bools() iter.Seq2[bool, error] {
	return func(yield func(bool, error) bool) {
		for idx := range a {
			if !yield(a.Bool(idx)) {
				return
			}
		}
	}
}

func (a Args) Ints() iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		for idx := range a {
			if !yield(a.Int(idx)) {
				return
			}
		}
	}
}

func (a Args) Int64s() iter.Seq2[int64, error] {
	return func(yield func(int64, error) bool) {
		for idx := range a {
			if !yield(a.Int64(idx)) {
				return
			}
		}
	}
}

func (a Args) Float64s() iter.Seq2[float64, error] {
	return func(yield func(float64, error) bool) {
		for idx := range a {
			if !yield(a.Float64(idx)) {
				return
			}
		}
	}
}

func (a Args) Float32s() iter.Seq2[float32, error] {
	return func(yield func(float32, error) bool) {
		for idx := range a {
			if !yield(a.Float32(idx)) {
				return
			}
		}
	}
}
