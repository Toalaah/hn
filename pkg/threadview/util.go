package threadview

import "golang.org/x/exp/constraints"

func clamp[T constraints.Ordered](x, lo, hi T) T {
	return max(lo, min(x, hi))
}
