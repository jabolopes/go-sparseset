package sparseset_test

import (
	"cmp"
	"fmt"
	"math/rand"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jabolopes/go-maths"
	"github.com/jabolopes/go-sparseset"
	"golang.org/x/exp/slices"
)

type join4Result[A, B, C, D any] struct {
	key int
	a   A
	b   B
	c   C
	d   D
	ok  bool
}

func join4All[A, B, C, D any](iterator *sparseset.Join4Iterator[A, B, C, D]) []join4Result[A, B, C, D] {
	results := []join4Result[A, B, C, D]{}
	for {
		key, a, b, c, d, ok := iterator.Next()
		if !ok {
			break
		}

		results = append(results, join4Result[A, B, C, D]{key, *a, *b, *c, *d, ok})
	}

	for i := 0; i < 10; i++ {
		key, a, b, c, d, ok := iterator.Next()
		if !ok {
			break
		}

		if key != 0 || a != nil || b != nil || c != nil || d != nil || ok != false {
			panic(fmt.Sprintf("Next() = %v, %v, %v, %v, %v, %v; want %v, %v, %v, %v, %v, %v", key, a, b, c, d, ok, 0, nil, nil, nil, nil, false))
		}
	}

	return results
}

func TestJoin4(t *testing.T) {
	var data1 []string
	faker.FakeData(&data1)

	var data2 []int
	faker.FakeData(&data2)

	var data3 []float32
	faker.FakeData(&data3)

	var data4 []float64
	faker.FakeData(&data4)

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

	set4 := sparseset.New[float64](4096, 1<<20)
	for i, value := range data4 {
		*set4.Add(i) = value
	}

	t.Run("Join A-B-C-D", func(t *testing.T) {
		want := make([]join4Result[string, int, float32, float64], maths.Min(maths.Min(maths.Min(len(data1), len(data2)), len(data3)), len(data4)))
		for i := range want {
			want[i] = join4Result[string, int, float32, float64]{i, data1[i], data2[i], data3[i], data4[i], true}
		}

		if got := join4All(sparseset.Join4(set1, set2, set3, set4)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join B-A-C-D", func(t *testing.T) {
		want := make([]join4Result[int, string, float32, float64], maths.Min(maths.Min(maths.Min(len(data1), len(data2)), len(data3)), len(data4)))
		for i := range want {
			want[i] = join4Result[int, string, float32, float64]{i, data2[i], data1[i], data3[i], data4[i], true}
		}

		if got := join4All(sparseset.Join4(set2, set1, set3, set4)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join C-B-A-D", func(t *testing.T) {
		want := make([]join4Result[float32, int, string, float64], maths.Min(maths.Min(maths.Min(len(data1), len(data2)), len(data3)), len(data4)))
		for i := range want {
			want[i] = join4Result[float32, int, string, float64]{i, data3[i], data2[i], data1[i], data4[i], true}
		}

		if got := join4All(sparseset.Join4(set3, set2, set1, set4)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join D-C-B-A", func(t *testing.T) {
		want := make([]join4Result[float64, float32, int, string], maths.Min(maths.Min(maths.Min(len(data1), len(data2)), len(data3)), len(data4)))
		for i := range want {
			want[i] = join4Result[float64, float32, int, string]{i, data4[i], data3[i], data2[i], data1[i], true}
		}

		if got := join4All(sparseset.Join4(set4, set3, set2, set1)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join A-A-A-A", func(t *testing.T) {
		want := make([]join4Result[string, string, string, string], len(data1))
		for i := range want {
			want[i] = join4Result[string, string, string, string]{i, data1[i], data1[i], data1[i], data1[i], true}
		}

		if got := join4All(sparseset.Join4(set1, set1, set1, set1)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})
}

func TestJoin4_RemoveRandomValues(t *testing.T) {
	var data1 []string
	faker.FakeData(&data1)

	var data2 []int
	faker.FakeData(&data2)

	var data3 []float32
	faker.FakeData(&data3)

	var data4 []float64
	faker.FakeData(&data4)

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

	set4 := sparseset.New[float64](4096, 1<<20)
	for i, value := range data4 {
		*set4.Add(i) = value
	}

	deleted4 := []int{}
	if len(data4) > 0 {
		for i := 0; i < 3; i++ {
			n := rand.Intn(len(data4))
			set4.Remove(n)
			deleted4 = append(deleted4, n)
		}
	}

	want := []join4Result[string, int, float32, float64]{}
	for i := 0; i < maths.Min(maths.Min(maths.Min(len(data1), len(data2)), len(data3)), len(data4)); i++ {
		if slices.Contains(deleted1, i) || slices.Contains(deleted2, i) || slices.Contains(deleted3, i) || slices.Contains(deleted4, i) {
			continue
		}

		want = append(want, join4Result[string, int, float32, float64]{i, data1[i], data2[i], data3[i], data4[i], true})
	}

	got := join4All(sparseset.Join4(set1, set2, set3, set4))

	slices.SortFunc(want, func(r1, r2 join4Result[string, int, float32, float64]) int { return cmp.Compare(r1.key, r2.key) })
	slices.SortFunc(got, func(r1, r2 join4Result[string, int, float32, float64]) int { return cmp.Compare(r1.key, r2.key) })

	if !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestJoin4_AllSetsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)
	set2 := sparseset.New[int](4096, 1<<20)
	set3 := sparseset.New[float32](4096, 1<<20)
	set4 := sparseset.New[float64](4096, 1<<20)

	want := []join4Result[string, int, float32, float64]{}
	if got := join4All(sparseset.Join4(set1, set2, set3, set4)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestJoin4_OneSetIsEmpty(t *testing.T) {
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

	var data4 []float64
	faker.FakeData(&data4)

	set4 := sparseset.New[float64](4096, 1<<30)
	for i, value := range data4 {
		*set4.Add(i) = value
	}

	want := []join4Result[string, int, float32, float64]{}
	if got := join4All(sparseset.Join4(set1, set2, set3, set4)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}
