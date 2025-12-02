package ent

import (
	"slices"
)

func NewBus() *Bus {
	return &Bus{}
}

type Bus struct {
	listeners []EntityUUID
}

// Causes a message to be sent to all subscribed entities on the bus.
func Emit(w *World, b *Bus, d any) {
	b.listeners = emitHelper(w, d, b.listeners...)
}

// Causes a message to be sent to the specified entities directly, bypassing any bus but still calling their HandleMessage function.
func EmitDirectly(w *World, d any, es ...EntityUUIDer) {
	emitHelper(w, d, es...)
}

func emitHelper[T EntityUUIDer](w *World, d any, es ...T) []T {
	newListeners := make([]T, 0, len(es))
	for _, l := range es {
		// Only keep if the entity is active in the world or it is queued
		if !w.HasOrQueued(l) {
			continue
		}
		newListeners = append(newListeners, l)
		// If the entity is in the world now, send message immediately,
		// otherwise, queue it to be sent once the entity is added (so we dont drop signals).
		e, ok := w.WithUUID(l.UUID())
		if ok {
			e.HandleMessage(w, d)
		} else {
			w.queuedAddWaitingSignals[l.UUID()] = append(w.queuedAddWaitingSignals[l.UUID()], d)
		}
	}
	return newListeners
}

// Subscribes the specified entities to the bus.
func Subscribe(b *Bus, es ...EntityUUIDer) {
	for _, e := range es {
		if slices.Contains(b.listeners, e.UUID()) {
			continue
		}
		b.listeners = append(b.listeners, e.UUID())
	}
}

// Unsubscribes the specified entities from the bus.
func Unsubscribe(b *Bus, es ...EntityUUIDer) {
	toDelete := make([]EntityUUID, len(es))
	for i, e := range es {
		toDelete[i] = e.UUID()
	}
	b.listeners = slices.DeleteFunc(b.listeners, func(eid EntityUUID) bool {
		return slices.Contains(toDelete, eid)
	})
}

// Unsubscribes all entities from the bus.
func UnsubscribeAll(b *Bus) {
	b.listeners = nil
}
