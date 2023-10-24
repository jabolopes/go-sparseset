package sparseset_test

import (
	"cmp"
	"fmt"
	"math/rand"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jabolopes/go-sparseset"
	"golang.org/x/exp/slices"
)

type iterateResult[A any] struct {
	key int
	a   A
	ok  bool
}

func iterateAll[A any](iterator *sparseset.Iterator[A]) []iterateResult[A] {
	results := []iterateResult[A]{}
	for {
		key, a, ok := iterator.Next()
		if !ok {
			break
		}

		results = append(results, iterateResult[A]{key, *a, ok})
	}

	for i := 0; i < 10; i++ {
		key, a, ok := iterator.Next()
		if !ok {
			break
		}

		if key != 0 || a != nil || ok != false {
			panic(fmt.Sprintf("Next() = %v, %v, %v; want %v, %v, %v", key, a, ok, 0, nil, false))
		}
	}

	return results
}

func TestIterate(t *testing.T) {
	var data []string
	faker.FakeData(&data)

	set := sparseset.New[string](4096, 1<<20)
	for i, value := range data {
		*set.Add(i) = value
	}

	{
		want := make([]iterateResult[string], len(data))
		for i := range want {
			want[i] = iterateResult[string]{i, data[i], true}
		}

		if got := iterateAll(sparseset.Iterate(set)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	}

	{
		want := make([]sparseset.IteratorResult[string], len(data))
		for i := range want {
			want[i] = sparseset.IteratorResult[string]{i, &set.Values()[i]}
		}

		if got := sparseset.Iterate(set).Collect(); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	}
}

func TestIterate_RemoveRandomValues(t *testing.T) {
	var data []string
	faker.FakeData(&data)

	set := sparseset.New[string](4096, 1<<20)
	for i, value := range data {
		*set.Add(i) = value
	}

	deleted := []int{}
	if len(data) > 0 {
		for i := 0; i < 3; i++ {
			n := rand.Intn(len(data))
			set.Remove(n)
			deleted = append(deleted, n)
		}
	}

	{
		want := []iterateResult[string]{}
		for i := range data {
			if !slices.Contains(deleted, i) {
				want = append(want, iterateResult[string]{i, data[i], true})
			}
		}

		got := iterateAll(sparseset.Iterate(set))

		slices.SortFunc(want, func(r1, r2 iterateResult[string]) int { return cmp.Compare(r1.key, r2.key) })
		slices.SortFunc(got, func(r1, r2 iterateResult[string]) int { return cmp.Compare(r1.key, r2.key) })

		if !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	}

	{
		want := []sparseset.IteratorResult[string]{}
		for i := range data {
			if !slices.Contains(deleted, i) {
				value, _ := set.Get(i)
				want = append(want, sparseset.IteratorResult[string]{i, value})
			}
		}

		got := sparseset.Iterate(set).Collect()

		slices.SortFunc(want, func(r1, r2 sparseset.IteratorResult[string]) int { return cmp.Compare(r1.Key, r2.Key) })
		slices.SortFunc(got, func(r1, r2 sparseset.IteratorResult[string]) int { return cmp.Compare(r1.Key, r2.Key) })

		if !slices.Equal(got, want) {
			t.Errorf("results = %v\nwant %v", got, want)
		}
	}
}

func TestIterate_EmptySet(t *testing.T) {
	set := sparseset.New[string](4096, 1<<20)

	{
		want := []iterateResult[string]{}
		if got := iterateAll(sparseset.Iterate(set)); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	}

	{
		want := []sparseset.IteratorResult[string]{}
		if got := sparseset.Iterate(set).Collect(); !slices.Equal(got, want) {
			t.Errorf("results = %v; want %v", got, want)
		}
	}
}
