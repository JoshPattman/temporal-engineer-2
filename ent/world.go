package ent

import (
	"iter"
	"slices"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// A collection of entities that can be indexed and updated in various ways.
type World struct {
	allEntities     *OrderedIndex[Entity]
	orderedByDraw   *OrderedIndex[Drawer]
	orderedByUpdate *OrderedIndex[Updater]
	physicsBodies   *OrderedIndex[PhysicsBody]
	byTags          map[string][]Entity
}

// Create a new, empty, world.
func NewWorld() *World {
	return &World{
		allEntities:     NewUnorderedIndex[Entity](),
		orderedByDraw:   NewOrderedIndex(func(d Drawer) int { return d.DrawLayer() }),
		orderedByUpdate: NewOrderedIndex(func(u Updater) int { return u.UpdateLayer() }),
		physicsBodies:   NewUnorderedIndex[PhysicsBody](),
		byTags:          make(map[string][]Entity, 0),
	}
}

// Add the entities to the world, adding it to all relevant indexes.
// The entity tags at this point in time will now be used of the entity.
// Each entity can only be added to the world once.
func (es *World) Add(toAdd ...Entity) {
	for _, e := range toAdd {
		if HasUUID(e) {
			continue
		}
		SetRandomUUID(e)
		es.allEntities.Add(e)
		es.orderedByDraw.AddUntyped(e)
		es.orderedByUpdate.AddUntyped(e)
		es.physicsBodies.AddUntyped(e)
		e.AfterAdd(es)
	}
}

// Remove the entity from the world.
// If the entity is not there, this will be a no-op.
func (es *World) Remove(toRemove ...Entity) {
	for _, e := range toRemove {
		if !HasUUID(e) {
			continue
		}
		ok := es.allEntities.Remove(e)
		if !ok {
			continue
		}
		es.orderedByDraw.RemoveUntyped(e)
		es.orderedByUpdate.RemoveUntyped(e)
		es.physicsBodies.RemoveUntyped(e)
		for tag, entities := range es.byTags {
			idx := slices.Index(entities, e)
			if idx != -1 {
				es.byTags[tag] = slices.Delete(es.byTags[tag], idx, idx+1)
				if len(es.byTags[tag]) == 0 {
					delete(es.byTags, tag)
				}
			}
		}
		ClearUUID(e)
	}
}

// Does the world contain this entity already?
func (es *World) Has(e Entity) bool {
	for e2 := range es.allEntities.All() {
		if e == e2 {
			return true
		}
	}
	return false
}

// Get all the entities for the given tag.
func (es *World) ForTag(tag string) iter.Seq[Entity] {
	if forTag, ok := es.byTags[tag]; ok {
		return slices.Values(forTag)
	} else {
		return func(yield func(Entity) bool) {}
	}
}

// Add the tags to the specific object.
func (es *World) AddTags(e Entity, tags ...string) {
	for _, tag := range tags {
		if _, ok := es.byTags[tag]; !ok {
			es.byTags[tag] = make([]Entity, 0)
		}
		if !slices.Contains(es.byTags[tag], e) {
			es.byTags[tag] = append(es.byTags[tag], e)
		}
	}
}

// Remove the tags from the specific object.
func (es *World) RemoveTags(e Entity, tags ...string) {
	for _, tag := range tags {
		idx := slices.Index(es.byTags[tag], e)
		if idx != -1 {
			es.byTags[tag] = slices.Delete(es.byTags[tag], idx, idx+1)
			if len(es.byTags[tag]) == 0 {
				delete(es.byTags, tag)
			}
		}
	}
}

// Update the world at the provided time interval.
// First, run all update steps.
// Then, add and remove all new entities.
// Finally, resolve physics then run collision handlers.
func (es *World) Update(win *pixelgl.Window, dt float64) {
	allToCreate := []Entity{}
	allToRemove := []Entity{}
	for e := range es.orderedByUpdate.All() {
		toCreate, toRemove := e.Update(win, es, dt)
		allToCreate = append(allToCreate, toCreate...)
		allToRemove = append(allToRemove, toRemove...)
	}
	for _, e := range allToCreate {
		if !es.Has(e) {
			es.Add(e)
		}
	}
	for _, e := range allToRemove {
		if es.Has(e) {
			es.Remove(e)
		}
	}

	fizBodies := slices.Collect(es.physicsBodies.All())
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
