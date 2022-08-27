package sparseset

func Lookup[A, B any](key int, setA *Set[A], setB *Set[B]) (*A, *B, bool) {
	a, aOk := setA.Get(key)
	b, bOk := setB.Get(key)
	return a, b, aOk && bOk
}

func Lookup3[A, B, C any](key int, setA *Set[A], setB *Set[B], setC *Set[C]) (*A, *B, *C, bool) {
	a, b, ok := Lookup(key, setA, setB)
	c, cOk := setC.Get(key)
	return a, b, c, ok && cOk
}
