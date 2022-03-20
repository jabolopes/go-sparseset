package sparseset

import (
	"golang.org/x/exp/slices"
)

// SortStableFunc sorts the Set according to 'compare'. The 'compare' function
// receives the 'left-hand-side' ID and Value and the 'right-hand-side' ID and
// Value. The 'compare' function should call methods on the Set (e.g.,
// Set.Get()) since SortStableFunc modifies the Set.
func SortStableFunc[T any](set *Set[T], compare func(int, *T, int, *T) bool) {
	slices.SortStableFunc(set.dense, func(i, j int) bool {
		iPos := set.index.Get(i)
		jPos := set.index.Get(j)
		return compare(i, &set.store[iPos], j, &set.store[jPos])
	})

	for i := 0; i < len(set.dense); i++ {
		for pos, next := i, set.index.Get(set.dense[i]); pos != next; {
			pos1 := set.index.Get(set.dense[pos])
			pos2 := set.index.Get(set.dense[next])

			set.store[pos1], set.store[pos2] = set.store[pos2], set.store[pos1]
			set.index.Set(set.dense[pos], pos)

			pos = next
			next = set.index.Get(set.dense[pos])
		}
	}
}
