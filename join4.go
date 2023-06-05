package sparseset

type Join4Iterator[A, B, C, D any] struct {
	get func() (int, *A, *B, *C, *D, bool)
}

func (i *Join4Iterator[A, B, C, D]) Next() (int, *A, *B, *C, *D, bool) {
	return i.get()
}

func Join4[A, B, C, D any](set1 *Set[A], set2 *Set[B], set3 *Set[C], set4 *Set[D]) *Join4Iterator[A, B, C, D] {
	var get func() (int, *A, *B, *C, *D, bool)

	if len(set1.dense) <= len(set2.dense) && len(set1.dense) <= len(set3.dense) && len(set1.dense) <= len(set4.dense) {
		iterator := Iterate(set1)
		get = func() (int, *A, *B, *C, *D, bool) {
			for {
				key, a, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, nil, false
				}

				b, ok := set2.Get(key)
				if !ok {
					continue
				}

				c, ok := set3.Get(key)
				if !ok {
					continue
				}

				d, ok := set4.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, d, true
			}
		}
	} else if len(set2.dense) <= len(set1.dense) && len(set2.dense) <= len(set3.dense) && len(set2.dense) <= len(set4.dense) {
		iterator := Iterate(set2)
		get = func() (int, *A, *B, *C, *D, bool) {
			for {
				key, b, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, nil, false
				}

				a, ok := set1.Get(key)
				if !ok {
					continue
				}

				c, ok := set3.Get(key)
				if !ok {
					continue
				}

				d, ok := set4.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, d, true
			}
		}
	} else if len(set3.dense) <= len(set1.dense) && len(set3.dense) <= len(set2.dense) && len(set3.dense) <= len(set4.dense) {
		iterator := Iterate(set3)
		get = func() (int, *A, *B, *C, *D, bool) {
			for {
				key, c, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, nil, false
				}

				a, ok := set1.Get(key)
				if !ok {
					continue
				}

				b, ok := set2.Get(key)
				if !ok {
					continue
				}

				d, ok := set4.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, d, true
			}
		}
	} else {
		iterator := Iterate(set4)
		get = func() (int, *A, *B, *C, *D, bool) {
			for {
				key, d, ok := iterator.Next()
				if !ok {
					return 0, nil, nil, nil, nil, false
				}

				a, ok := set1.Get(key)
				if !ok {
					continue
				}

				b, ok := set2.Get(key)
				if !ok {
					continue
				}

				c, ok := set3.Get(key)
				if !ok {
					continue
				}

				return key, a, b, c, d, true
			}
		}
	}

	return &Join4Iterator[A, B, C, D]{get}
}

func EmptyJoin4Iterator[A, B, C, D any]() *Join4Iterator[A, B, C, D] {
	return &Join4Iterator[A, B, C, D]{func() (int, *A, *B, *C, *D, bool) {
		return 0, nil, nil, nil, nil, false
	}}
}
