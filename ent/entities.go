package ent

import (
	"iter"
	"slices"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// A collection of entities that can be indexed and updated in various ways.
type World struct {
	allEntities     *index[Entity]
	orderedByDraw   *index[Drawer]
	orderedByUpdate *index[Updater]
	physicsBodies   *index[PhysicsBody]
	byTags          map[string][]Entity
}

// Create a new, empty, world.
func NewWorld() *World {
	return &World{
		allEntities:     newUnorderedIndex[Entity](),
		orderedByDraw:   newOrderedIndex(func(d Drawer) int { return d.DrawLayer() }),
		orderedByUpdate: newOrderedIndex(func(u Updater) int { return u.UpdateLayer() }),
		physicsBodies:   newUnorderedIndex[PhysicsBody](),
		byTags:          make(map[string][]Entity, 0),
	}
}

// Add the entities to the world, adding it to all relevant indexes.
// The entity tags at this point in time will now be used of the entity.
// Each entity can only be added to the world once.
func (es *World) Add(toAdd ...Entity) {
	for _, e := range toAdd {
		ok := es.allEntities.tryAdd(e)
		if !ok {
			continue
		}
		es.orderedByDraw.tryAdd(e)
		es.orderedByUpdate.tryAdd(e)
		es.physicsBodies.tryAdd(e)
		{
			for _, tag := range e.Tags() {
				if _, ok := es.byTags[tag]; !ok {
					es.byTags[tag] = make([]Entity, 0)
				}
				es.byTags[tag] = append(es.byTags[tag], e)
			}
		}
	}
}

// Remove the entity from the world.
// If the entity is not there, this will be a no-op.
func (es *World) Remove(e Entity) {
	ok := es.allEntities.tryRemove(e)
	if !ok {
		return
	}
	es.orderedByDraw.tryRemove(e)
	es.orderedByUpdate.tryRemove(e)
	es.physicsBodies.tryRemove(e)
	{
		for _, tag := range e.Tags() {
			idx := slices.Index(es.byTags[tag], e)
			if idx == -1 {
				panic("was not in entities")
			}
			es.byTags[tag] = slices.Delete(es.byTags[tag], idx, idx+1)
		}
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
		e.Draw(win, worldToScreen)
	}
}

func newOrderedIndex[T comparable](order func(T) int) *index[T] {
	return &index[T]{
		items: make([]T, 0),
		layer: order,
	}
}

func newUnorderedIndex[T comparable]() *index[T] {
	return &index[T]{
		items: make([]T, 0),
		layer: func(t T) int { return 0 },
	}
}

type index[T comparable] struct {
	items []T
	layer func(T) int
}

func (index *index[T]) tryAdd(item any) bool {
	itemTyped, ok := item.(T)
	if !ok {
		return false
	}

	if slices.Contains(index.items, itemTyped) {
		return false
	}

	i := 0
	for i < len(index.items) && index.layer(itemTyped) <= index.layer(index.items[i]) {
		i += 1
	}

	index.items = slices.Insert(index.items, i, itemTyped)
	return true
}

func (index *index[T]) tryRemove(item any) bool {
	itemTyped, ok := item.(T)
	if !ok {
		return false
	}
	i := slices.Index(index.items, itemTyped)
	if i == -1 {
		return false
	}
	index.items = slices.Delete(index.items, i, i+1)
	return true
}

func (index *index[T]) All() iter.Seq[T] {
	return slices.Values(index.items)
}
