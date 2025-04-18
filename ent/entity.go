package ent

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

type Entity interface {
	// Called once per frame to update behaviour.
	Update(win *pixelgl.Window, all *World, dt float64) (toCreate, toDestroy []Entity)
	// Called before any Draw method is called on any entity.
	// Should be used to ready self for drawing.
	PreDraw(win *pixelgl.Window)
	// Called to draw self to screen.
	Draw(win *pixelgl.Window, worldToScreen pixel.Matrix)
	// Gets the tags of this entity.
	// Should NEVER change after entity has been created.
	Tags() []string
	// Called to get the update layer for this entity.
	// Higher values will be updated first.
	// Should NEVER change after entity has been created.
	UpdateLayer() int
	// Called to get the draw layer for this entity.
	// Higher values will be drawn first (appear below other objects).
	// Should NEVER change after entity has been created.
	DrawLayer() int
}

// You may inherit this to give your entity some default behaviour.
// This is optional though
type EntityBase struct{}

func (*EntityBase) Update(win *pixelgl.Window, all *World, dt float64) (toCreate, toDestroy []Entity) {
	return nil, nil
}

func (*EntityBase) PreDraw(win *pixelgl.Window)                          {}
func (*EntityBase) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {}
func (*EntityBase) Tags() []string                                       { return nil }
func (*EntityBase) UpdateLayer() int                                     { return 0 }
func (*EntityBase) DrawLayer() int                                       { return 0 }
