package ent

import (
	"iter"
	"slices"
)

type Entities struct {
	orderedByDraw   []Entity
	orderedByUpdate []Entity
	byTags          map[string][]Entity
}

func NewEntities() *Entities {
	return &Entities{
		byTags: make(map[string][]Entity, 0),
	}
}

func (es *Entities) Add(e Entity) {
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
}

func (es *Entities) Remove(e Entity) {
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
}

func (es *Entities) Has(e Entity) bool {
	for _, e2 := range es.orderedByUpdate {
		if e == e2 {
			return true
		}
	}
	return false
}

func (es *Entities) ByDraw() iter.Seq2[int, Entity] {
	return slices.All(es.orderedByDraw)
}

func (es *Entities) ByUpdate() iter.Seq2[int, Entity] {
	return slices.All(es.orderedByUpdate)
}

func (es *Entities) ForTag(tag string) iter.Seq2[int, Entity] {
	if forTag, ok := es.byTags[tag]; ok {
		return slices.All(forTag)
	} else {
		return func(yield func(int, Entity) bool) {}
	}
}

func (es *Entities) ForTagSlice(tag string) []Entity {
	if forTag, ok := es.byTags[tag]; ok {
		return slices.Clone(forTag)
	} else {
		return nil
	}
}
