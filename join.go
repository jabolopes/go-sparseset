package sparseset

type JoinIterator[A, B any] struct {
	get func() (int, *A, *B, bool)
}

func (i *JoinIterator[A, B]) Next() (int, *A, *B, bool) {
	return i.get()
}

func Join[A, B any](set1 *Set[A], set2 *Set[B]) *JoinIterator[A, B] {
	var get func() (int, *A, *B, bool)

	if len(set1.dense) <= len(set2.dense) {
		iterator := Iterate(set1)
		get = func() (int, *A, *B, bool) {
			for {
				key, a, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, false
				}

				b, ok := set2.Get(key)
				if !ok {
					continue
				}

				return key, a, b, true
			}
		}
	} else {
		iterator := Iterate(set2)
		get = func() (int, *A, *B, bool) {
			for {
				key, b, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, false
				}

				a, ok := set1.Get(key)
				if !ok {
					continue
				}

				return key, a, b, true
			}
		}
	}

	return &JoinIterator[A, B]{get}
}
