package ent

import "iter"

func FilterEntitiesByType[T any](xs iter.Seq2[int, Entity]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		next, stop := iter.Pull2(xs)
		for {
			i, e, ok := next()
			if !ok {
				stop()
				return
			}
			t, ok := e.(T)
			if !ok {
				continue
			}
			if !yield(i, t) {
				stop()
				return
			}
		}
	}
}
