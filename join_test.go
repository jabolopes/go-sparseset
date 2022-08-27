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

type joinResult[A, B any] struct {
	key int
	a   A
	b   B
	ok  bool
}

func joinAll[A, B any](iterator *sparseset.JoinIterator[A, B]) []joinResult[A, B] {
	results := []joinResult[A, B]{}
	for {
		key, a, b, ok := iterator.Next()
		if !ok {
			break
		}

		results = append(results, joinResult[A, B]{key, *a, *b, ok})
	}

	for i := 0; i < 10; i++ {
		key, a, b, ok := iterator.Next()
		if !ok {
			break
		}

		if key != 0 || a != nil || b != nil || ok != false {
			panic(fmt.Sprintf("Next() = %v, %v, %v, %v; want %v, %v, %v, %v", key, a, b, ok, 0, nil, nil, false))
		}
	}

	return results
}

func TestJoin(t *testing.T) {
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

	t.Run("Join A-B", func(t *testing.T) {
		want := make([]joinResult[string, int], maths.Min(len(data1), len(data2)))
		for i := range want {
			want[i] = joinResult[string, int]{i, data1[i], data2[i], true}
		}

		if got := joinAll(sparseset.Join(set1, set2)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join B-A", func(t *testing.T) {
		want := make([]joinResult[int, string], maths.Min(len(data1), len(data2)))
		for i := range want {
			want[i] = joinResult[int, string]{i, data2[i], data1[i], true}
		}

		if got := joinAll(sparseset.Join(set2, set1)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join A-A", func(t *testing.T) {
		want := make([]joinResult[string, string], len(data1))
		for i, value := range data1 {
			want[i] = joinResult[string, string]{i, value, value, true}
		}

		if got := joinAll(sparseset.Join(set1, set1)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})

	t.Run("Join B-B", func(t *testing.T) {
		want := make([]joinResult[int, int], len(data2))
		for i, value := range data2 {
			want[i] = joinResult[int, int]{i, value, value, true}
		}

		if got := joinAll(sparseset.Join(set2, set2)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	})
}

func TestJoin_RemoveRandomValues(t *testing.T) {
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

	want := []joinResult[string, int]{}
	for i := 0; i < maths.Min(len(data1), len(data2)); i++ {
		if slices.Contains(deleted1, i) || slices.Contains(deleted2, i) {
			continue
		}

		want = append(want, joinResult[string, int]{i, data1[i], data2[i], true})
	}

	got := joinAll(sparseset.Join(set1, set2))

	slices.SortFunc(want, func(r1, r2 joinResult[string, int]) bool { return r1.key < r2.key })
	slices.SortFunc(got, func(r1, r2 joinResult[string, int]) bool { return r1.key < r2.key })

	if !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestJoin_AllSetsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)
	set2 := sparseset.New[int](4096, 1<<20)

	want := []joinResult[string, int]{}
	if got := joinAll(sparseset.Join(set1, set2)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestJoin_OneSetIsEmpty(t *testing.T) {
	set1 := sparseset.New[string](4096, 1<<20)

	var data2 []int
	faker.FakeData(&data2)

	set2 := sparseset.New[int](4096, 1<<20)
	for i, value := range data2 {
		*set2.Add(i) = value
	}

	want := []joinResult[string, int]{}
	if got := joinAll(sparseset.Join(set1, set2)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}
