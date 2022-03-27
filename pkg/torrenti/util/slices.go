package util

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func BinarySearchContain[E constraints.Ordered](x []E, target E) bool {
	idx := slices.BinarySearch(x, target)
	if len(x) <= idx {
		return false
	}
	return x[idx] == target
}
