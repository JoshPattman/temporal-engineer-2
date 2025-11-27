package ent

import "iter"

// Filter the iterator of entities to only those of the given type.
func FilterEntitiesByType[T any](xs iter.Seq[Entity]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(xs)
		for {
			e, ok := next()
			if !ok {
				stop()
				return
			}
			t, ok := e.(T)
			if !ok {
				continue
			}
			if !yield(t) {
				stop()
				return
			}
		}
	}
}

// Get the first item, or if there are no items return false.
func First[T any](xs iter.Seq[T]) (T, bool) {
	for t := range xs {
		return t, true
	}
	return *new(T), false
}
