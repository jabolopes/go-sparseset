package sparseset_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jabolopes/go-maths"
	"github.com/jabolopes/go-sparseset"
	"golang.org/x/exp/slices"
)

type join3Result[A, B, C any] struct {
	key int
	a   A
	b   B
	c   C
	ok  bool
}

func join3All[A, B, C any](iterator *sparseset.Join3Iterator[A, B, C]) []join3Result[A, B, C] {
	results := []join3Result[A, B, C]{}
	for {
		key, a, b, c, ok := iterator.Next()
		if !ok {
			break
		}

		results = append(results, join3Result[A, B, C]{key, *a, *b, *c, ok})
	}

	for i := 0; i < 10; i++ {
		key, a, b, c, ok := iterator.Next()
		if !ok {
			break
		}

		if key != 0 || a != nil || b != nil || c != nil || ok != false {
			panic(fmt.Sprintf("Next() = %v, %v, %v, %v, %v; want %v, %v, %v, %v, %v", key, a, b, c, ok, 0, nil, nil, nil, false))
		}
	}

	return results
}

func TestJoin3(t *testing.T) {
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

	t.Run("Join A-B-C", func(t *testing.T) {
		want := make([]join3Result[string, int, float32], maths.Min(maths.Min(len(data1), len(data2)), len(data3)))
		for i := range want {
			want[i] = join3Result[string, int, float32]{i, data1[i], data2[i], data3[i], true}
		}

		if got := join3All(sparseset.Join3(set1, set2, set3)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join B-A-C", func(t *testing.T) {
		want := make([]join3Result[int, string, float32], maths.Min(maths.Min(len(data1), len(data2)), len(data3)))
		for i := range want {
			want[i] = join3Result[int, string, float32]{i, data2[i], data1[i], data3[i], true}
		}

		if got := join3All(sparseset.Join3(set2, set1, set3)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join C-B-A", func(t *testing.T) {
		want := make([]join3Result[float32, int, string], maths.Min(maths.Min(len(data1), len(data2)), len(data3)))
		for i := range want {
			want[i] = join3Result[float32, int, string]{i, data3[i], data2[i], data1[i], true}
		}

		if got := join3All(sparseset.Join3(set3, set2, set1)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join A-A-A", func(t *testing.T) {
		want := make([]join3Result[string, string, string], len(data1))
		for i := range want {
			want[i] = join3Result[string, string, string]{i, data1[i], data1[i], data1[i], true}
		}

		if got := join3All(sparseset.Join3(set1, set1, set1)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})
}

func TestJoin3_RemoveRandomValues(t *testing.T) {
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

	set3 := sparseset.New[float32](4096, 1<<20)
	for i, value := range data3 {
		*set3.Add(i) = value
	}

	deleted3 := []int{}
	if len(data3) > 0 {
		for i := 0; i < 3; i++ {
			n := rand.Intn(len(data3))
			set3.Remove(n)
			deleted3 = append(deleted3, n)
		}
	}

	want := []join3Result[string, int, float32]{}
	for i := 0; i < maths.Min(maths.Min(len(data1), len(data2)), len(data3)); i++ {
		if slices.Contains(deleted1, i) || slices.Contains(deleted2, i) || slices.Contains(deleted3, i) {
			continue
		}

		want = append(want, join3Result[string, int, float32]{i, data1[i], data2[i], data3[i], true})
	}

	got := join3All(sparseset.Join3(set1, set2, set3))

	slices.SortFunc(want, func(r1, r2 join3Result[string, int, float32]) bool { return r1.key < r2.key })
	slices.SortFunc(got, func(r1, r2 join3Result[string, int, float32]) bool { return r1.key < r2.key })

	if !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestJoin3_AllSetsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)
	set2 := sparseset.New[int](4096, 1<<20)
	set3 := sparseset.New[float32](4096, 1<<20)

	want := []join3Result[string, int, float32]{}
	if got := join3All(sparseset.Join3(set1, set2, set3)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestJoin3_OneSetIsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)

	var data2 []int
	faker.FakeData(&data2)

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	var data3 []float32
	faker.FakeData(&data3)

	set3 := sparseset.New[float32](4096, 1<<30)
	for i, value := range data3 {
		*set3.Add(i) = value
	}

	want := []join3Result[string, int, float32]{}
	if got := join3All(sparseset.Join3(set1, set2, set3)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}
