package ent

import (
	"iter"
	"slices"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

type World struct {
	orderedByDraw   []Entity
	orderedByUpdate []Entity
	physicsBodies   []PhysicsBody
	byTags          map[string][]Entity
}

func NewWorld() *World {
	return &World{
		byTags: make(map[string][]Entity, 0),
	}
}

func (es *World) Add(toAdd ...Entity) {
	for _, e := range toAdd {
		{
			inserted := false
			for i := range es.orderedByDraw {
				if e.DrawLayer() > es.orderedByDraw[i].DrawLayer() {
					es.orderedByDraw = slices.Insert(es.orderedByDraw, i, e)
					inserted = true
					break
				}
			}
			if !inserted {
				es.orderedByDraw = append(es.orderedByDraw, e)
			}
		}
		{
			inserted := false
			for i := range es.orderedByUpdate {
				if e.UpdateLayer() > es.orderedByUpdate[i].UpdateLayer() {
					es.orderedByUpdate = slices.Insert(es.orderedByUpdate, i, e)
					inserted = true
					break
				}
			}
			if !inserted {
				es.orderedByUpdate = append(es.orderedByUpdate, e)
			}
		}
		{
			for _, tag := range e.Tags() {
				if _, ok := es.byTags[tag]; !ok {
					es.byTags[tag] = make([]Entity, 0)
				}
				es.byTags[tag] = append(es.byTags[tag], e)
			}
		}
		if e, ok := e.(PhysicsBody); ok {
			es.physicsBodies = append(es.physicsBodies, e)
		}
	}
}

func (es *World) Remove(e Entity) {
	{
		idx := slices.Index(es.orderedByDraw, e)
		if idx == -1 {
			panic("was not in entities")
		}
		es.orderedByDraw = slices.Delete(es.orderedByDraw, idx, idx+1)
	}
	{
		idx := slices.Index(es.orderedByUpdate, e)
		if idx == -1 {
			panic("was not in entities")
		}
		es.orderedByUpdate = slices.Delete(es.orderedByUpdate, idx, idx+1)
	}
	{
		for _, tag := range e.Tags() {
			idx := slices.Index(es.byTags[tag], e)
			if idx == -1 {
				panic("was not in entities")
			}
			es.byTags[tag] = slices.Delete(es.byTags[tag], idx, idx+1)
		}
	}
	if e, ok := e.(PhysicsBody); ok {
		idx := slices.Index(es.physicsBodies, e)
		if idx == -1 {
			panic("was not in entities")
		}
		es.physicsBodies = slices.Delete(es.physicsBodies, idx, idx+1)
	}
}

func (es *World) Has(e Entity) bool {
	for _, e2 := range es.orderedByUpdate {
		if e == e2 {
			return true
		}
	}
	return false
}

func (es *World) ByDraw() iter.Seq[Entity] {
	return All(es.orderedByDraw)
}

func (es *World) ByUpdate() iter.Seq[Entity] {
	return All(es.orderedByUpdate)
}

func (es *World) WithPhysics() iter.Seq[PhysicsBody] {
	return All(es.physicsBodies)
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
