package sparseset_test

import (
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jabolopes/go-maths"
	"github.com/jabolopes/go-sparseset"
)

func TestLookup3(t *testing.T) {
	var data1 []string
	faker.FakeData(&data1)

	var data2 []int
	faker.FakeData(&data2)

	var data3 []float32
	faker.FakeData(&data3)

	set1 := sparseset.New[string](4096, 1<<20)
	for i, value := range data1 {
		*set1.Add(i) = value
	}

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	set3 := sparseset.New[float32](4096, 1<<20)
	for i, value := range data3 {
		*set3.Add(i) = value
	}

	t.Run("Lookup A-B-C", func(t *testing.T) {
		for i := 0; i < maths.Max(maths.Max(len(data1), len(data2)), len(data3))+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			var wantB *int
			if i < len(data2) {
				wantB = &data2[i]
			}

			var wantC *float32
			if i < len(data3) {
				wantC = &data3[i]
			}

			wantOk := wantA != nil && wantB != nil && wantC != nil

			if a, b, c, ok := sparseset.Lookup3(i, set1, set2, set3); !reflect.DeepEqual(a, wantA) || !reflect.DeepEqual(b, wantB) || !reflect.DeepEqual(c, wantC) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v, %v; %v, %v, %v, %v", a, b, c, ok, wantA, wantB, wantC, wantOk)
			}
		}
	})

	t.Run("Lookup B-A-C", func(t *testing.T) {
		for i := 0; i < maths.Max(maths.Max(len(data1), len(data2)), len(data3))+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			var wantB *int
			if i < len(data2) {
				wantB = &data2[i]
			}

			var wantC *float32
			if i < len(data3) {
				wantC = &data3[i]
			}

			wantOk := wantA != nil && wantB != nil && wantC != nil

			if b, a, c, ok := sparseset.Lookup3(i, set2, set1, set3); !reflect.DeepEqual(a, wantA) || !reflect.DeepEqual(b, wantB) || !reflect.DeepEqual(c, wantC) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v, %v; %v, %v, %v, %v", b, a, c, ok, wantB, wantA, wantC, wantOk)
			}
		}
	})

	t.Run("Lookup C-B-A", func(t *testing.T) {
		for i := 0; i < maths.Max(maths.Max(len(data1), len(data2)), len(data3))+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			var wantB *int
			if i < len(data2) {
				wantB = &data2[i]
			}

			var wantC *float32
			if i < len(data3) {
				wantC = &data3[i]
			}

			wantOk := wantA != nil && wantB != nil && wantC != nil

			if c, b, a, ok := sparseset.Lookup3(i, set3, set2, set1); !reflect.DeepEqual(a, wantA) || !reflect.DeepEqual(b, wantB) || !reflect.DeepEqual(c, wantC) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v, %v; %v, %v, %v, %v", c, b, a, ok, wantC, wantB, wantA, wantOk)
			}
		}
	})

	t.Run("Lookup A-A-A", func(t *testing.T) {
		for i := 0; i < len(data1)+1; i++ {
			var wantA *string
			if i < len(data1) {
				wantA = &data1[i]
			}

			wantOk := wantA != nil

			if a1, a2, a3, ok := sparseset.Lookup3(i, set1, set1, set1); !reflect.DeepEqual(a1, wantA) || !reflect.DeepEqual(a2, wantA) || !reflect.DeepEqual(a3, wantA) || ok != wantOk {
				t.Errorf("Lookup() = %v, %v, %v, %v; %v, %v, %v, %v", a1, a2, a3, ok, wantA, wantA, wantA, wantOk)
			}
		}
	})
}

func TestLookup3_AllSetsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)
	set2 := sparseset.New[int](4096, 1<<20)
	set3 := sparseset.New[float32](4096, 1<<20)

	if a, b, c, ok := sparseset.Lookup3(0, set1, set2, set3); a != nil || b != nil || c != nil || ok {
		t.Errorf("Lookup() = %v, %v, %v, %v; %v, %v, %v, %v", a, b, c, ok, nil, nil, nil, nil)
	}
}

func TestLookup3_OneSetIsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)

	var data2 []int
	faker.FakeData(&data2)

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	set3 := sparseset.New[float32](4096, 1<<20)

	wantB, _ := set2.Get(0)

	if a, b, c, ok := sparseset.Lookup3(0, set1, set2, set3); a != nil || b != wantB || c != nil || ok {
		t.Errorf("Lookup() = %v, %v, %v, %v; %v, %v, %v, %v", a, b, c, ok, nil, wantB, nil, nil)
	}
}
