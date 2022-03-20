package sparseset

type Options[Value any] struct {
	DestroyValue func(*Value)
}

// Set is a sparse set with a value store.
//
// The Set can be accessed via read-only operations and iterators
// concurrently, but it cannot be concurrently accessed by readers and
// writers, and cannot be concurrently accessed by multiple writers.
//
// This is thread-compatible.
type Set[Value any] struct {
	// Sparse (paged) array. Stores positions (pos) by key.
	index *PagedArray[int]
	// Stores keys by position (pos) contiguously.
	dense []int
	// Stores values by position (pos) contiguously.
	store []Value
	// Destroy (or uninitializes) values when they are removed from the store.
	destroyValue func(*Value)
}

func (s *Set[Value]) Length() int     { return s.index.Length() }
func (s *Set[Value]) Values() []Value { return s.store }

func (s *Set[Value]) Add(key int) *Value {
	if key < 0 || key >= s.index.NullValue() {
		return nil
	}

	pos := s.index.Get(key)
	if pos != s.index.NullValue() {
		return &s.store[pos]
	}

	pos = len(s.store)

	s.dense = append(s.dense, key)

	var value Value
	s.store = append(s.store, value)

	s.index.Set(key, pos)
	return &s.store[pos]
}

func (s *Set[Value]) Remove(key int) {
	if key < 0 || key >= s.index.NullValue() {
		return
	}

	pos := s.index.Get(key)
	if pos == s.index.NullValue() {
		return
	}

	last := len(s.store) - 1
	var defaultValue Value

	if pos == last {
		// The value being removed is the store's last element, so to remove it from
		// the set, unset it from the index and clear its value.

		// Remove from index.
		s.index.Unset(key)

		// Destroy store value.
		s.destroyValue(&s.store[pos])

		// Facilitate GC.
		s.store[pos] = defaultValue

		// Remove element from the store.
		s.dense = s.dense[:last]
		s.store = s.store[:last]

		return
	}

	// The value being removed is in the middle of the store, so to remove it we
	// need to swap it so that it becomes the store's last element and then clear
	// it, otherwise this will break iteration since elements in the store won't
	// be contiguous anymore.

	s.index.Unset(key)
	s.index.Set(s.dense[last], pos)

	// Destroy store value.
	s.destroyValue(&s.store[pos])

	s.store[pos], s.store[last] = s.store[last], defaultValue
	s.dense[pos], s.dense[last] = s.dense[last], s.index.NullValue()

	s.dense = s.dense[:last]
	s.store = s.store[:last]
}

func (s *Set[Value]) Get(key int) (*Value, bool) {
	if key < 0 || key >= s.index.NullValue() {
		return nil, false
	}

	pos := s.index.Get(key)
	if pos == s.index.NullValue() {
		return nil, false
	}

	return &s.store[pos], true
}

func New[Value any](defaultPageSize, nullKey int) *Set[Value] {
	return NewWithOptions[Value](defaultPageSize, nullKey, Options[Value]{})
}

func NewWithOptions[Value any](defaultPageSize, nullKey int, options Options[Value]) *Set[Value] {
	if options.DestroyValue == nil {
		options.DestroyValue = func(*Value) {}
	}

	return &Set[Value]{
		NewPagedArray(defaultPageSize, nullKey),
		[]int{},
		[]Value{},
		options.DestroyValue,
	}
}
