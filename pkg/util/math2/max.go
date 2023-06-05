package math2

import (
	"golang.org/x/exp/constraints"
)

func Max[N constraints.Ordered](a, b N) N {
	if a > b {
		return a
	}
	return b
}

func Min[N constraints.Ordered](a, b N) N {
	if a < b {
		return a
	}
	return b
}
