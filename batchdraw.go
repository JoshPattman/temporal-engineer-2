package main

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.Entity = &BatchDraw{}

func NewBatchDraw(spritePath string) *BatchDraw {
	return &BatchDraw{
		Batch: pixel.NewBatch(
			&pixel.TrianglesData{},
			GlobalSpriteManager.Picture(spritePath),
		),
	}
}

type BatchDraw struct {
	Batch *pixel.Batch
}

// Draw implements ent.Entity.
func (b *BatchDraw) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	b.Batch.Draw(win)
}

// DrawLayer implements ent.Entity.
func (b *BatchDraw) DrawLayer() int {
	return -3
}

// Tags implements ent.Entity.
func (b *BatchDraw) Tags() []string {
	return nil
}

// Update implements ent.Entity.
func (b *BatchDraw) Update(win *pixelgl.Window, all *ent.Entities, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	b.Batch.Clear()
	return nil, nil
}

// UpdateLayer implements ent.Entity.
func (b *BatchDraw) UpdateLayer() int {
	return -1
}
