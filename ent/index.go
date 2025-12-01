package ent

import (
	"iter"
	"slices"
	"sort"
)

func NewOrderedIndex[T EntityUUIDer](orderFunc func(T) int) *Index[T] {
	return &Index[T]{
		orderedItems:  make([]orderedItem[T], 0),
		containsUUIDs: make(map[EntityUUID]struct{}),
		orderOf:       orderFunc,
	}
}

func NewUnorderedIndex[T EntityUUIDer]() *Index[T] {
	return NewOrderedIndex(func(t T) int { return 0 })
}

type orderedItem[T any] struct {
	item  T
	order int
}

type Index[T EntityUUIDer] struct {
	orderedItems  []orderedItem[T]
	containsUUIDs map[EntityUUID]struct{}
	orderOf       func(T) int
}

func (oi *Index[T]) AddUntyped(item any) bool {
	itemTyped, ok := item.(T)
	if !ok {
		return false
	}
	return oi.Add(itemTyped)
}

func (oi *Index[T]) RemoveUntyped(item any) bool {
	itemTyped, ok := item.(T)
	if !ok {
		return false
	}
	return oi.Remove(itemTyped)
}

func (oi *Index[T]) HasUntyped(item any) bool {
	itemTyped, ok := item.(T)
	if !ok {
		return false
	}
	return oi.Has(itemTyped)
}

func (oi *Index[T]) Add(item T) bool {
	id := item.UUID()
	if _, ok := oi.containsUUIDs[id]; ok {
		return false
	}
	oi.containsUUIDs[id] = struct{}{}
	order := oi.orderOf(item)
	insertIndex := sort.Search(
		len(oi.orderedItems),
		func(i int) bool {
			return order > oi.orderedItems[i].order
		},
	)
	oi.orderedItems = slices.Insert(
		oi.orderedItems,
		insertIndex,
		orderedItem[T]{item, order},
	)
	return true
}

func (oi *Index[T]) Remove(item T) bool {
	id := item.UUID()
	if _, ok := oi.containsUUIDs[id]; !ok {
		return false
	}
	delete(oi.containsUUIDs, id)
	order := oi.orderOf(item)
	// We should only start the search from the position where the orders start to equal (can skip all orders below that)
	startSearchIndex := sort.Search(
		len(oi.orderedItems),
		func(i int) bool {
			return order >= oi.orderedItems[i].order
		},
	)
	for i := startSearchIndex; i < len(oi.orderedItems); i++ {
		checkItem := oi.orderedItems[i]
		if checkItem.item.UUID() == id {
			oi.orderedItems = slices.Delete(
				oi.orderedItems,
				i, i+1,
			)
			return true
		}
	}
	panic("somthing has gone wrong with maintaining the index")
}

func (index *Index[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range index.orderedItems {
			if !yield(item.item) {
				return
			}
		}
	}
}

func (index *Index[T]) Has(item T) bool {
	_, ok := index.containsUUIDs[item.UUID()]
	return ok
}

func (index Index[T]) Len() int {
	return len(index.orderedItems)
}
