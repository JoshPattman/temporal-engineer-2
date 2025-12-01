package ent

import (
	"iter"
	"slices"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// A collection of entities that can be indexed and updated in various ways.
type World struct {
	byIDLookup              map[EntityUUID]Entity
	allEntities             *Index[Entity]
	orderedByDraw           *Index[Drawer]
	orderedByUpdate         *Index[Updater]
	physicsBodies           *Index[PhysicsBody]
	byTags                  map[string]*Index[Entity]
	queuedAdd               []Entity
	queuedAddWaitingSignals map[EntityUUID][]any
	queuedRemove            []Entity
}

// Create a new, empty, world.
func NewWorld() *World {
	return &World{
		byIDLookup:              make(map[EntityUUID]Entity),
		allEntities:             NewUnorderedIndex[Entity](),
		orderedByDraw:           NewOrderedIndex(func(d Drawer) int { return d.DrawLayer() }),
		orderedByUpdate:         NewOrderedIndex(func(u Updater) int { return u.UpdateLayer() }),
		physicsBodies:           NewUnorderedIndex[PhysicsBody](),
		byTags:                  make(map[string]*Index[Entity], 0),
		queuedAddWaitingSignals: make(map[EntityUUID][]any),
	}
}

// Add the entities to the world, adding it to all relevant indexes.
// The entity tags at this point in time will now be used of the entity.
// Each entity can only be added to the world once.
// Will also call AfterAdd, and will then send any queued signals.
func (es *World) Add(toAdd ...Entity) {
	for _, e := range toAdd {
		uid := e.UUID()
		if _, ok := es.byIDLookup[uid]; ok {
			continue
		}
		es.byIDLookup[uid] = e
		es.allEntities.Add(e)
		es.orderedByDraw.AddUntyped(e)
		es.orderedByUpdate.AddUntyped(e)
		es.physicsBodies.AddUntyped(e)
		e.AfterAdd(es)
		for _, sig := range es.queuedAddWaitingSignals[e.UUID()] {
			e.HandleMessage(sig)
		}
		delete(es.queuedAddWaitingSignals, e.UUID())
	}
}

// Queue the entities to be added to the world when appropriate.
func (w *World) Instantiate(toInstantiate ...Entity) {
	w.queuedAdd = append(w.queuedAdd, toInstantiate...)
}

// Remove the entity from the world.
// If the entity is not there, this will be a no-op.
func (es *World) Remove(toRemove ...Entity) {
	for _, e := range toRemove {
		ok := es.allEntities.Remove(e)
		if !ok {
			continue
		}
		delete(es.byIDLookup, e.UUID())
		es.orderedByDraw.RemoveUntyped(e)
		es.orderedByUpdate.RemoveUntyped(e)
		es.physicsBodies.RemoveUntyped(e)
		for tag, index := range es.byTags {
			if index.Remove(e) && es.byTags[tag].Len() == 0 {
				delete(es.byTags, tag)
			}
		}
	}
}

// Queue the entities to be removed to the world when appropriate.
func (w *World) Destroy(toDestroy ...Entity) {
	w.queuedRemove = append(w.queuedRemove, toDestroy...)
}

// Does the world contain this entity / uuid already?
func (es *World) Has(e EntityUUIDer) bool {
	_, ok := es.byIDLookup[e.UUID()]
	return ok
}

// Does the world contain this entity / uuid already, or is it queued to be added?
func (es *World) HasOrQueued(id EntityUUIDer) bool {
	_, ok := es.byIDLookup[id.UUID()]
	if ok {
		return true
	}
	for _, e := range es.queuedAdd {
		if e.UUID() == id {
			return true
		}
	}
	return false
}

// Get all the entities for the given tag.
func (es *World) ForTag(tag string) iter.Seq[Entity] {
	if forTag, ok := es.byTags[tag]; ok {
		return forTag.All()
	} else {
		return func(yield func(Entity) bool) {}
	}
}

// Get the specific entity with the given UUID
func (es *World) WithUUID(id EntityUUID) (Entity, bool) {
	e, ok := es.byIDLookup[id]
	return e, ok
}

// Add the tags to the specific object.
func (es *World) AddTags(e Entity, tags ...string) {
	for _, tag := range tags {
		if _, ok := es.byTags[tag]; !ok {
			es.byTags[tag] = NewUnorderedIndex[Entity]()
		}
		es.byTags[tag].Add(e)
	}
}

// Remove the tags from the specific object.
func (es *World) RemoveTags(e Entity, tags ...string) {
	for _, tag := range tags {
		index, ok := es.byTags[tag]
		if !ok {
			continue
		}
		if index.Remove(e) && es.byTags[tag].Len() == 0 {
			delete(es.byTags, tag)
		}
	}
}

// Update the world at the provided time interval.
// First, run all update steps.
// Then, add and remove all new entities.
// Finally, resolve physics then run collision handlers.
func (es *World) Update(win *pixelgl.Window, dt float64) {
	for e := range es.orderedByUpdate.All() {
		e.Update(win, es, dt)
	}
	for _, e := range es.queuedAdd {
		if !es.Has(e) {
			es.Add(e)
		}
	}
	es.queuedAdd = nil
	for _, e := range es.queuedRemove {
		if es.Has(e) {
			es.Remove(e)
		}
	}
	es.queuedRemove = nil

	fizBodies := slices.Collect(es.physicsBodies.All())
	for _, body := range fizBodies {
		body, ok := body.(ActivePhysicsBody)
		if ok && body.IsPhysicsActive() {
			body.PysicsUpdate(dt)
		}
	}
	cols := StatelessCollisionPhysics(fizBodies)

	for _, col := range cols {
		self, ok := col.Self.(CollisionListener)
		if ok {
			self.OnCollision(col)
		}
		col = col.ForOther()
		self, ok = col.Self.(CollisionListener)
		if ok {
			self.OnCollision(col)
		}
	}
}

// Call predraw on all entities, then call draw.
// Pass the provided world to screen mapping to all draw calls.
func (es *World) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	for e := range es.orderedByDraw.All() {
		e.PreDraw(win)
	}
	for e := range es.orderedByDraw.All() {
		e.Draw(win, es, worldToScreen)
	}
}
