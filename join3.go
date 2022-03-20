package sparseset

type Join3Iterator[A, B, C any] struct {
	get func() (int, *A, *B, *C, bool)
}

func (i *Join3Iterator[A, B, C]) Next() (int, *A, *B, *C, bool) {
	return i.get()
}

func Join3[A, B, C any](set1 *Set[A], set2 *Set[B], set3 *Set[C]) *Join3Iterator[A, B, C] {
	var get func() (int, *A, *B, *C, bool)

	if len(set1.dense) <= len(set2.dense) && len(set1.dense) <= len(set3.dense) {
		iterator := Iterate(set1)
		get = func() (int, *A, *B, *C, bool) {
			for {
				key, a, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, false
				}

				b, ok := set2.Get(key)
				if !ok {
					continue
				}

				c, ok := set3.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, true
			}
		}
	} else if len(set2.dense) <= len(set1.dense) && len(set2.dense) <= len(set3.dense) {
		iterator := Iterate(set2)
		get = func() (int, *A, *B, *C, bool) {
			for {
				key, b, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, false
				}

				a, ok := set1.Get(key)
				if !ok {
					continue
				}

				c, ok := set3.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, true
			}
		}
	} else {
		iterator := Iterate(set3)
		get = func() (int, *A, *B, *C, bool) {
			for {
				key, c, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, false
				}

				a, ok := set1.Get(key)
				if !ok {
					continue
				}

				b, ok := set2.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, true
			}
		}
	}

	return &Join3Iterator[A, B, C]{get}
}
