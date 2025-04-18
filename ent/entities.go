package ent

import (
	"iter"
	"slices"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

type World struct {
	allEntities     *index[Entity]
	orderedByDraw   *index[Drawer]
	orderedByUpdate *index[Updater]
	physicsBodies   *index[PhysicsBody]
	byTags          map[string][]Entity
}

func NewWorld() *World {
	return &World{
		allEntities:     newUnorderedIndex[Entity](),
		orderedByDraw:   newOrderedIndex(func(d Drawer) int { return d.DrawLayer() }),
		orderedByUpdate: newOrderedIndex(func(u Updater) int { return u.UpdateLayer() }),
		physicsBodies:   newUnorderedIndex[PhysicsBody](),
		byTags:          make(map[string][]Entity, 0),
	}
}

func (es *World) Add(toAdd ...Entity) {
	for _, e := range toAdd {
		es.allEntities.tryAdd(e)
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

func (es *World) Remove(e Entity) {
	es.allEntities.tryRemove(e)
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

func (es *World) Has(e Entity) bool {
	for e2 := range es.allEntities.All() {
		if e == e2 {
			return true
		}
	}
	return false
}

func (es *World) ByDraw() iter.Seq[Drawer] {
	return es.orderedByDraw.All()
}

func (es *World) ByUpdate() iter.Seq[Updater] {
	return es.orderedByUpdate.All()
}

func (es *World) WithPhysics() iter.Seq[PhysicsBody] {
	return es.physicsBodies.All()
}

func (es *World) ForTag(tag string) iter.Seq[Entity] {
	if forTag, ok := es.byTags[tag]; ok {
		return All(forTag)
	} else {
		return func(yield func(Entity) bool) {}
	}
}

func (es *World) Update(win *pixelgl.Window, dt float64) {
	allToCreate := []Entity{}
	allToRemove := []Entity{}
	for e := range es.ByUpdate() {
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

	fizBodies := slices.Collect(es.WithPhysics())
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

func (es *World) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	for e := range es.ByDraw() {
		e.PreDraw(win)
	}
	for e := range es.ByDraw() {
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
	return All(index.items)
}
