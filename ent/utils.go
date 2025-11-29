package ent

import (
	"iter"
	"math"

	"github.com/gopxl/pixel"
)

// Filter the iterator of entities to only those of the given type.
func OfType[T any](xs iter.Seq[Entity]) iter.Seq[T] {
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

type Positioner interface {
	Position() pixel.Vec
}

// Get the closest transform
func Closest[T Positioner](pos pixel.Vec, xs iter.Seq[T]) (T, bool) {
	var c T
	d := math.Inf(1)
	for item := range xs {
		dist := item.Position().To(pos).Len()
		if dist < d {
			d = dist
			c = item
		}
	}
	if math.IsInf(d, 1) {
		return *new(T), false
	} else {
		return c, true
	}
}
