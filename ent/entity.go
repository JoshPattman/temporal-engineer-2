package ent

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

type EntityUUID string

func (e EntityUUID) UUID() EntityUUID { return e }

// An object that can be drawn to the screen.
type Drawer interface {
	EntityUUIDer
	// Called before any Draw method is called on any entity.
	// Should be used to ready self for drawing.
	PreDraw(win *pixelgl.Window)
	// Called to draw self to screen.
	Draw(win *pixelgl.Window, world *World, worldToScreen pixel.Matrix)
	// Called to get the draw layer for this entity.
	// Higher values will be drawn first (appear below other objects).
	// Should NEVER change after entity has been created.
	DrawLayer() int
}

// An object that can be updated on each update step.
type Updater interface {
	EntityUUIDer
	// Called once per frame to update behaviour.
	Update(win *pixelgl.Window, world *World, dt float64)
	// Called to get the update layer for this entity.
	// Higher values will be updated first.
	// Should NEVER change after entity has been created.
	UpdateLayer() int
}

type EntityUUIDer interface {
	UUID() EntityUUID
}

// An object that can be added to a world.
type Entity interface {
	EntityUUIDer
	AfterAdd(*World)
	HandleMessage(*World, any)
}
