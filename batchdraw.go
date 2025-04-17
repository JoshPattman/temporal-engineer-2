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
	ent.EntityBase
	Batch *pixel.Batch
}

// PreDraw implements ent.Entity.
func (b *BatchDraw) PreDraw(win *pixelgl.Window) {
	b.Batch.Clear()
}

// Draw implements ent.Entity.
func (b *BatchDraw) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	b.Batch.Draw(win)
}

// DrawLayer implements ent.Entity.
func (b *BatchDraw) DrawLayer() int {
	return -3
}
