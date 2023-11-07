package merge

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

func Mean[N Number, E any](items []E, f func(E) (N, bool)) (N, bool) {
	var t float64
	for _, item := range items {
		if _, ok := f(item); ok {
			t++
		}
	}
	if t == 0 {
		return 0, false
	}

	var res float64
	for _, item := range items {
		if v, ok := f(item); ok {
			res += float64(v) / t
		}
	}
	return N(res), true
}

func Mode[N Number, E any](items []E, f func(E) (N, bool)) (N, bool) {
	var vals = make(map[N]int)
	// count values
	for _, item := range items {
		if v, ok := f(item); ok {
			_, ok := vals[v]
			if !ok {
				vals[v] = 0
			}
			vals[v]++
		}
	}
	// find max
	max := 0
	fail := false
	var res N
	for v, c := range vals {
		if c > max {
			max = c
			res = v
			fail = false
		} else if c == max {
			// if more than 1 value has the same count, fail
			fail = true
		}
	}

	return N(res), fail
}

func Max[N Number, E any](items []E, f func(E) (N, bool)) (N, bool) {
	var res N
	var c int
	for _, item := range items {
		if v, ok := f(item); ok {
			if v > res {
				c++
				res = v
			}
		}
	}
	return res, c > 0
}

func Sum[N Number, E any](items []E, f func(E) (N, bool)) (N, bool) {
	var res N
	var c int
	for _, item := range items {
		if v, ok := f(item); ok {
			c++
			res += v
		}
	}
	return res, c > 0
}

func Ptr[T any](v T, ok bool) *T {
	if ok {
		return &v
	}
	return nil
}
