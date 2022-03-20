package sparseset

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type page2[Value any] struct {
	values    []Value
	numValues int
}

type PagedArray[Value constraints.Ordered] struct {
	pool      *sync.Pool
	pages     []page2[Value]
	pageSize  int
	nullValue Value
	length    int
}

func (a *PagedArray[Value]) NullValue() Value { return a.nullValue }
func (a *PagedArray[Value]) Length() int      { return a.length }

func (a *PagedArray[Value]) Get(index int) Value {
	if index < 0 {
		return a.nullValue
	}

	pageNum := index / a.pageSize
	pageOffset := index % a.pageSize

	if pageNum >= len(a.pages) {
		return a.nullValue
	}

	page := &a.pages[pageNum]
	if page.numValues == 0 {
		return a.nullValue
	}

	return page.values[pageOffset]
}

func (a *PagedArray[Value]) Set(index int, value Value) {
	if index < 0 || value >= a.nullValue {
		return
	}

	pageNum := index / a.pageSize
	pageOffset := index % a.pageSize

	if pageNum >= len(a.pages) {
		// Allocate pages between len(a.pages) and a.pages[pageNum].
		if diff := pageNum - len(a.pages); diff >= 0 {
			a.pages = append(a.pages, make([]page2[Value], diff+1)...)
		}
	}

	page := &a.pages[pageNum]
	if page.values == nil {
		// initialize the page.
		page.values = a.pool.Get().([]Value)
		page.numValues = 0

		for i := range page.values {
			page.values[i] = a.nullValue
		}
	}

	if page.values[pageOffset] == a.nullValue {
		page.numValues++
		a.length++
	}
	page.values[pageOffset] = value
}

// TODO: Reclaim pages at the end of the array that are empty.
func (a *PagedArray[Value]) Unset(index int) {
	if index < 0 {
		return
	}

	pageNum := index / a.pageSize
	pageOffset := index % a.pageSize

	if pageNum >= len(a.pages) {
		return
	}

	page := &a.pages[pageNum]
	if page.values == nil || page.values[pageOffset] == a.nullValue {
		return
	}

	page.numValues--
	a.length--
	page.values[pageOffset] = a.nullValue

	if page.numValues <= 0 {
		var values []Value
		values, page.values = page.values, nil
		a.pool.Put(values)
	}
}

func (a *PagedArray[Value]) Clear() {
	for _, page := range a.pages {
		if page.values != nil {
			var values []Value
			values, page.values = page.values, nil
			a.pool.Put(values)
		}
	}

	a.pages = nil
	a.length = 0
}

func NewPagedArray[Value constraints.Ordered](pageSize int, nullValue Value) *PagedArray[Value] {
	return &PagedArray[Value]{
		&sync.Pool{
			New: func() any {
				return make([]Value, pageSize)
			},
		},
		nil, /* pages */
		pageSize,
		nullValue,
		0, /* length */
	}
}
