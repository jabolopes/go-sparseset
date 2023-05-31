package sparseset

type IteratorResult[A any] struct {
	Key   int
	Value *A
}

// Iterator can be used to traverse the keys and values of a Set. This iterator
// is read-only (see thread-safety notes on Set).
//
// This is thread-compatible.
type Iterator[A any] struct {
	get   func(int) (int, *A, bool)
	index int
}

// Next returns the next value for this iterator. Returns a key, a value and a
// boolean. If the boolean is true, then the key and value are valid, otherwise
// the key and value invalid (e.g., default initialized). If the boolean is
// false, then the end of the iteration has been reached and subsequent calls to
// Next() will not return any new elements.
func (i *Iterator[A]) Next() (int, *A, bool) {
	key, a, ok := i.get(i.index)
	if !ok {
		return 0, nil, false
	}

	i.index++
	return key, a, true
}

// Collect traverses the remaining elements and stores them in an array. This is
// more convenient than Next() but it performs memory allocations to create and
// resize the array. If memory allocations are considered expensive (e.g.,
// memory pressure, garbage collection, etc), then Next() should be preferred.
func (i *Iterator[A]) Collect() []IteratorResult[A] {
	results := []IteratorResult[A]{}
	for {
		key, value, ok := i.Next()
		if !ok {
			break
		}

		results = append(results, IteratorResult[A]{key, value})
	}
	return results
}

// Iterate returns an iterator that can be used to traverse all the keys and
// values of the set.
//
// TODO: Avoid allocating memory for the Iterator itself.
func Iterate[A any](set *Set[A]) *Iterator[A] {
	dense := set.dense
	store := set.store
	get := func(i int) (int, *A, bool) {
		if i < 0 || i >= len(dense) {
			return 0, nil, false
		}
		return dense[i], &store[i], true
	}

	return &Iterator[A]{get, 0}
}

func EmptyIterator[A any]() *Iterator[A] {
	return &Iterator[A]{func(int) (int, *A, bool) { return 0, nil, false }, 0}
}
