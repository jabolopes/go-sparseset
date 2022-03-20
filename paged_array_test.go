package sparseset_test

import (
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/jabolopes/go-sparseset"
)

func TestConstructor(t *testing.T) {
	const pageSize = 1 << 10
	const nullValue = int(1e6)
	sparseset.NewPagedArray(pageSize, nullValue)
}

func TestLifecycle(t *testing.T) {
	const pageSize = 5
	const nullValue = 11
	array := sparseset.NewPagedArray(pageSize, nullValue)

	array.Set(6, 0)

	if got := array.Get(5); got != nullValue {
		t.Errorf("Get(%d) = %v; want %v", 10, got, nullValue)
	}

	want := 10
	for i := 0; i < 2; i++ {
		array.Set(5, want)

		if got := array.Get(5); got != want {
			t.Errorf("Get(%d) = %v; want %v", 10, got, want)
		}

		wantLength := 2
		if got := array.Length(); got != wantLength {
			t.Errorf("Length() = %v; want %v", got, wantLength)
		}
	}

	for i := 0; i < 2; i++ {
		array.Unset(5)

		if got := array.Get(5); got != nullValue {
			t.Errorf("Get(%d) = %v; want %v", 10, got, nullValue)
		}

		wantLength := 1
		if got := array.Length(); got != wantLength {
			t.Errorf("Length() = %v; want %v", got, wantLength)
		}
	}

	array.Set(5, want)
	array.Clear()

	if got := array.Get(5); got != nullValue {
		t.Errorf("Get(%d) = %v; want %v", 10, got, nullValue)
	}
}

func TestGet_Rand(t *testing.T) {
	const n = 1000

	const pageSize = 1 << 10
	const nullValue = int(1e6)
	array := sparseset.NewPagedArray(pageSize, nullValue)

	for i := 0; i < n; i++ {
		index := rand.Int()
		if got := array.Get(index); got != nullValue {
			t.Errorf("Get(%d) = %v; want %v", index, got, nullValue)
		}

		if got := array.Length(); got != 0 {
			t.Errorf("Length() = %v; want %v", got, 0)
		}
	}
}

func TestSet_Rand(t *testing.T) {
	const n = 1000

	const pageSize = 1 << 10
	const nullValue = int(1e6)
	array := sparseset.NewPagedArray(pageSize, nullValue)

	set := map[int]int{}

	for len(set) < n {
		index := int(rand.Int31())
		value := rand.Intn(nullValue + n)

		if value < nullValue {
			set[index] = value
		}

		array.Set(index, value)

		if got := array.Length(); got != len(set) {
			t.Errorf("Length() = %v; want %v", got, len(set))
		}
	}

	for index, want := range set {
		if got := array.Get(index); got != want {
			t.Errorf("Get(%d) = %v; want %v", index, got, want)
		}

		if got := array.Length(); got != len(set) {
			t.Errorf("Length() = %v; want %v", got, len(set))
		}
	}

	array.Clear()

	for index, _ := range set {
		if got := array.Get(index); got != nullValue {
			t.Errorf("Get(%d) = %v; want %v", index, got, nullValue)
		}
	}

	if got := array.Length(); got != 0 {
		t.Errorf("Length() = %v; want %v", got, 0)
	}
}

func TestUnset_Rand(t *testing.T) {
	const n = 1000

	const pageSize = 1 << 10
	const nullValue = int(1e6)
	array := sparseset.NewPagedArray(pageSize, nullValue)

	set := map[int]int{}
	unset := map[int]struct{}{}

	for len(set) < n {
		index := int(rand.Int31())
		value := rand.Intn(nullValue + n)

		if value < nullValue {
			set[index] = value
		}

		array.Set(index, value)

		if got := array.Length(); got != len(set) {
			t.Errorf("Length() = %v; want %v", got, len(set))
		}
	}

	{
		i := 0
		for index, _ := range set {
			if i%2 == 0 {
				unset[index] = struct{}{}
				array.Unset(index)
			}

			wantLength := len(set) - len(unset)
			if got := array.Length(); got != wantLength {
				t.Errorf("Length() = %v; want %v", got, wantLength)
			}

			i++
		}
	}

	for index, value := range set {
		want := value
		if _, ok := unset[index]; ok {
			want = nullValue
		}

		if got := array.Get(index); got != want {
			t.Errorf("Get(%d) = %v; want %v", index, got, want)
		}

		wantLength := len(set) - len(unset)
		if got := array.Length(); got != wantLength {
			t.Errorf("Length() = %v; want %v", got, wantLength)
		}
	}

	array.Clear()

	for index, _ := range set {
		if got := array.Get(index); got != nullValue {
			t.Errorf("Get(%d) = %v; want %v", index, got, nullValue)
		}
	}

	if got := array.Length(); got != 0 {
		t.Errorf("Length() = %v; want %v", got, 0)
	}
}

func TestGet_Quick(t *testing.T) {
	const pageSize = 1 << 10
	const nullValue = int(1e6)
	array := sparseset.NewPagedArray(pageSize, nullValue)

	get := func(index int) bool {
		return array.Get(index) == nullValue
	}

	if err := quick.Check(get, nil); err != nil {
		t.Error(err)
	}
}
