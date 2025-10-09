package thread

import "golang.org/x/exp/constraints"

func clamp[T constraints.Ordered](x, lo, hi T) T { return max(lo, min(x, hi)) }

func min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}
	return b
}
