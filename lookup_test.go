package sparseset_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jabolopes/go-maths"
	"github.com/jabolopes/go-sparseset"
	"golang.org/x/exp/slices"
)

func TestLookup(t *testing.T) {
	var data1 []string
	faker.FakeData(&data1)

	var data2 []int
	faker.FakeData(&data2)

	set1 := sparseset.New[string](4096, 1<<20)
	for i, value := range data1 {
		*set1.Add(i) = value
	}

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	t.Run("Lookup A-B", func(t *testing.T) {
		for i := 0; i < maths.Max(len(data1), len(data2))+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			var wantB *int
			if i < len(data2) {
				wantB = &data2[i]
			}

			wantOk := wantA != nil && wantB != nil

			if a, b, ok := sparseset.Lookup(i, set1, set2); !reflect.DeepEqual(a, wantA) || !reflect.DeepEqual(b, wantB) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", a, b, ok, wantA, wantB, wantOk)
			}
		}
	})

	t.Run("Lookup B-A", func(t *testing.T) {
		for i := 0; i < maths.Max(len(data1), len(data2))+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			var wantB *int
			if i < len(data2) {
				wantB = &data2[i]
			}

			wantOk := wantA != nil && wantB != nil

			if b, a, ok := sparseset.Lookup(i, set2, set1); !reflect.DeepEqual(a, wantA) || !reflect.DeepEqual(b, wantB) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", b, a, ok, wantB, wantA, wantOk)
			}
		}
	})

	t.Run("Lookup A-A", func(t *testing.T) {
		for i := 0; i < len(data1)+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			wantOk := wantA != nil

			if a1, a2, ok := sparseset.Lookup(i, set1, set1); !reflect.DeepEqual(a1, wantA) || !reflect.DeepEqual(a2, wantA) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", a1, a2, ok, wantA, wantA, wantOk)
			}
		}
	})

	t.Run("Lookup B-B", func(t *testing.T) {
		for i := 0; i < len(data2)+1; i++ {
			var wantB *int
			if i < len(data2) {
				wantB = &data2[i]
			}

			wantOk := wantB != nil

			if b1, b2, ok := sparseset.Lookup(i, set2, set2); !reflect.DeepEqual(b1, wantB) || !reflect.DeepEqual(b2, wantB) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", b1, b2, ok, wantB, wantB, wantOk)
			}
		}
	})
}

func TestLookup_RemoveRandomValues(t *testing.T) {
	var data1 []string
	faker.FakeData(&data1)

	var data2 []int
	faker.FakeData(&data2)

	set1 := sparseset.New[string](4096, 1<<20)
	for i, value := range data1 {
		*set1.Add(i) = value
	}

	deleted1 := []int{}
	if len(data1) > 0 {
		for i := 0; i < 3; i++ {
			n := rand.Intn(len(data1))
			set1.Remove(n)
			deleted1 = append(deleted1, n)
		}
	}

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	deleted2 := []int{}
	if len(data2) > 0 {
		for i := 0; i < 3; i++ {
			n := rand.Intn(len(data2))
			set2.Remove(n)
			deleted2 = append(deleted2, n)
		}
	}

	for i := 0; i < maths.Max(len(data1), len(data2))+1; i++ {
		var wantA *string
		if i < len(data1) && !slices.Contains(deleted1, i) {
			wantA = &data1[i]
		}

		var wantB *int
		if i < len(data2) && !slices.Contains(deleted2, i) {
			wantB = &data2[i]
		}

		wantOk := wantA != nil && wantB != nil

		if a, b, ok := sparseset.Lookup(i, set1, set2); !reflect.DeepEqual(a, wantA) || !reflect.DeepEqual(b, wantB) || ok != wantOk {
			t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", a, b, ok, wantA, wantB, wantOk)
		}
	}
}

func TestLookup_BothSetsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)
	set2 := sparseset.New[int](4096, 1<<20)

	if a, b, ok := sparseset.Lookup(0, set1, set2); a != nil || b != nil || ok {
		t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", a, b, ok, nil, nil, nil)
	}
}

func TestLookup_OneSetIsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)

	var data2 []int
	faker.FakeData(&data2)

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	wantB, _ := set2.Get(0)

	if a, b, ok := sparseset.Lookup(0, set1, set2); a != nil || b != wantB || ok {
		t.Errorf("Lookup() = %v, %v, %v; %v, %v, %v", a, b, ok, nil, wantB, nil)
	}
}
