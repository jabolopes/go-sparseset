package sparseset_test

import (
	"cmp"
	"math/rand"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jabolopes/go-sparseset"
	"golang.org/x/exp/slices"
)

func TestSort(t *testing.T) {
	var data []string
	faker.FakeData(&data)

	set := sparseset.New[string](4096, 1<<20)
	for i, value := range data {
		*set.Add(i) = value
	}

	if len(data) > 0 {
		for i := 0; i < 3; i++ {
			n := rand.Intn(len(data))
			set.Remove(n)
		}
	}

	want := iterateAll(sparseset.Iterate(set))
	slices.SortStableFunc(want, func(x, y iterateResult[string]) int { return cmp.Compare(x.a, y.a) })

	sparseset.SortStableFunc(set, func(keyA int, a *string, keyB int, b *string) int {
		if want := data[keyA]; *a != want {
			t.Errorf("a = %v; want %v", *a, want)
		}

		if want := data[keyB]; *b != want {
			t.Errorf("a = %v; want %v", *b, want)
		}

		return cmp.Compare(*a, *b)
	})

	if got := iterateAll(sparseset.Iterate(set)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}

func TestSort_EmptySet(t *testing.T) {
	set := sparseset.New[string](4096, 1<<20)

	sparseset.SortStableFunc(set, func(_ int, a *string, _ int, b *string) int {
		t.Fatalf("compare() called")
		return 0
	})

	want := []iterateResult[string]{}
	if got := iterateAll(sparseset.Iterate(set)); !slices.Equal(got, want) {
		t.Errorf("results = %v; want %v", got, want)
	}
}
