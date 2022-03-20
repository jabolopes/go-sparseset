package sparseset_test

import (
	"fmt"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/jabolopes/go-sparseset"
	"golang.org/x/exp/slices"
)

type MyValue struct {
	value int
}

func ExampleSet() {
	const pageSize = 1 << 10
	const nullValue = 1 << 20
	set := sparseset.New[string](pageSize, nullValue)
	*set.Add(10) = "hello"
	*set.Add(20) = "world"

	// Properties.
	fmt.Printf("Length %d\n", set.Length())

	// Removal.
	set.Remove(20)

	// Lookup.
	if value, ok := set.Get(10); ok {
		fmt.Printf("Value %d = %s\n", 10, *value)
	} else {
		fmt.Println("Set does not contain id")
	}

	// Traversal.
	for iterator := sparseset.Iterate(set); ; {
		key, value, ok := iterator.Next()
		if !ok {
			break
		}

		fmt.Printf("Traverse %d = %s\n", key, *value)
	}

	// Joining 2 sets.
	set2 := sparseset.New[int](pageSize, nullValue)
	*set2.Add(10) = 1
	*set2.Add(20) = 2

	for iterator := sparseset.Join(set, set2); ; {
		key, value1, value2, ok := iterator.Next()
		if !ok {
			break
		}

		fmt.Printf("Join %d = %s, %d\n", key, *value1, *value2)
	}

	// Output: Length 2
	// Value 10 = hello
	// Traverse 10 = hello
	// Join 10 = hello, 1
}

func TestAdd(t *testing.T) {
	const pageSize = 1 << 10
	const nullValue = 1 << 20
	set := sparseset.New[MyValue](pageSize, nullValue)

	value := set.Add(10)
	value.value = 10
	if value == nil {
		t.Errorf("Add(10) = %v; want non-%v", value, nil)
	}

	if got := set.Length(); got != 1 {
		t.Errorf("Length() = %d; want %d", got, 1)
	}

	wantValues := []MyValue{{10}}
	if got := set.Values(); !slices.Equal(got, wantValues) {
		t.Errorf("Values() = %v; want %v", got, wantValues)
	}

	if got, ok := set.Get(10); got != value || !ok {
		t.Errorf("Get(10) = %v, %v; want %v, %v", got, ok, value, true)
	}

	value.value = 100
	if got := set.Add(10); !reflect.DeepEqual(got, value) {
		t.Errorf("Add(10) = %v; want %v", got, value)
	}

	if got := set.Length(); got != 1 {
		t.Errorf("Length() = %d; want %d", got, 1)
	}

	wantValues = []MyValue{{100}}
	if got := set.Values(); !slices.Equal(got, wantValues) {
		t.Errorf("Values() = %v; want %v", got, wantValues)
	}

	if got, ok := set.Get(10); got != value || !ok {
		t.Errorf("Get(10) = %v, %v; want %v, %v", got, ok, value, true)
	}
}

func TestAdd_Quick(t *testing.T) {
	const pageSize = 1 << 10
	const nullValue = 1 << 20
	set := sparseset.New[MyValue](pageSize, nullValue)

	add := func(key int) bool {
		length := set.Length()

		value := set.Add(key)
		got, ok := set.Get(key)

		if key >= 0 && key < nullValue {
			return value != nil && got == value && ok && set.Length() == length+1
		}

		return value == nil && got == nil && !ok && set.Length() == length
	}

	if err := quick.Check(add, nil); err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	const pageSize = 1 << 10
	const nullValue = 1 << 20

	called := false
	options := sparseset.Options[MyValue]{func(value *MyValue) {
		called = true
	}}
	set := sparseset.NewWithOptions[MyValue](pageSize, nullValue, options)

	_ = set.Add(10)
	if got := set.Length(); got != 1 {
		t.Errorf("Length() = %d; want %d", got, 1)
	}

	set.Remove(10)

	if got, ok := set.Get(10); got != nil || ok {
		t.Errorf("Get(10) = %v, %v; want %v, %v", got, ok, nil, false)
	}

	if got := set.Length(); got != 0 {
		t.Errorf("Length() = %d; want %d", got, 0)
	}

	if !called {
		t.Errorf("called = %v; want %v", called, true)
	}
}

func TestRemove_Quick(t *testing.T) {
	const pageSize = 1 << 10
	const nullValue = 1 << 20

	called := false
	options := sparseset.Options[MyValue]{func(value *MyValue) {
		called = true
	}}
	set := sparseset.NewWithOptions[MyValue](pageSize, nullValue, options)

	remove := func(key int) bool {
		called = false
		value := set.Add(key)

		gotLength := set.Length()

		if value == nil {
			if gotLength != 0 {
				t.Logf("Length() = %d; want %d", gotLength, 0)
				return false
			}
		} else {
			if gotLength != 1 {
				t.Logf("Length() = %d; want %d", gotLength, 1)
				return false
			}
		}

		set.Remove(key)

		if got, ok := set.Get(key); got != nil || ok {
			t.Logf("Get(%d) = %v, %v; want %v, %v", key, got, ok, nil, false)
			return false
		}

		if got := set.Length(); got != 0 {
			t.Logf("Length() = %d; want %d", gotLength, 0)
			return false
		}

		if value != nil {
			if !called {
				t.Errorf("called = %v; want %v", called, true)
			}
		}

		return true
	}

	if err := quick.Check(remove, nil); err != nil {
		t.Error(err)
	}
}

func TestAddRemove_Sequence(t *testing.T) {
	const n = 1000

	const pageSize = 1 << 10
	const nullValue = 1 << 20
	set := sparseset.New[MyValue](pageSize, nullValue)

	for i := 0; i < n; i++ {
		value := set.Add(i)
		value.value = i
	}

	if got := set.Length(); got != n {
		t.Errorf("Length() = %d; want %d", got, n)
	}

	wantValues := []MyValue{}
	for i := 0; i < n; i++ {
		wantValue := MyValue{i}
		if got, ok := set.Get(i); *got != wantValue || !ok {
			t.Errorf("Get(%d) = %v; %v; want %v, %v", i, got, ok, wantValue, true)
		}
		wantValues = append(wantValues, wantValue)
	}

	if got := set.Values(); !slices.Equal(got, wantValues) {
		t.Errorf("Values() = %v; want %v", got, wantValues)
	}

	for i := 0; i < 2; i++ {
		set.Remove(i)

		wantValues := []MyValue{}
		for j := i + 1; j < n; j++ {
			wantValue := MyValue{j}
			if got, ok := set.Get(j); got == nil || *got != wantValue || !ok {
				t.Errorf("Get(%d) = %v; %v; want %v, %v", j, got, ok, wantValue, true)
			}
			wantValues = append(wantValues, wantValue)
		}

		got := append([]MyValue{}, set.Values()...)
		slices.SortFunc(got, func(i, j MyValue) bool { return i.value < j.value })
		slices.SortFunc(wantValues, func(i, j MyValue) bool { return i.value < j.value })
		if !slices.Equal(got, wantValues) {
			t.Errorf("Values() = %v; want %v", got, wantValues)
		}
	}
}
