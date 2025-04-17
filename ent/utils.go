package ent

import "iter"

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

func All[T any](xs []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range xs {
			if !yield(item) {
				return
			}
		}
	}
}

func First[T any](xs iter.Seq[T]) (T, bool) {
	for t := range xs {
		return t, true
	}
	return *new(T), false
}
