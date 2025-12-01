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

func Emit(w *World, b *Bus, d any) {
	newListeners := make([]EntityUUID, 0, len(b.listeners))
	for _, l := range b.listeners {
		// Only keep if the entity is active in the world or it is queued
		if !w.HasOrQueued(l) {
			continue
		}
		newListeners = append(newListeners, l)
		// If the entity is in the world now, send message immediately,
		// otherwise, queue it to be sent once the entity is added (so we dont drop signals).
		e, ok := w.WithUUID(l)
		if ok {
			e.HandleMessage(w, d)
		} else {
			w.queuedAddWaitingSignals[l] = append(w.queuedAddWaitingSignals[l], d)
		}
	}
	b.listeners = newListeners
}

func Subscribe(b *Bus, es ...EntityUUIDer) {
	for _, e := range es {
		if slices.Contains(b.listeners, e.UUID()) {
			continue
		}
		b.listeners = append(b.listeners, e.UUID())
	}
}

func Unsubscribe(b *Bus, es ...EntityUUIDer) {
	toDelete := make([]EntityUUID, len(es))
	for i, e := range es {
		toDelete[i] = e.UUID()
	}
	b.listeners = slices.DeleteFunc(b.listeners, func(eid EntityUUID) bool {
		return slices.Contains(toDelete, eid)
	})
}

func UnsubscribeAll(b *Bus) {
	b.listeners = nil
}
