package ent

import (
	"github.com/google/uuid"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// An object that can be drawn to the screen.
type Drawer interface {
	UUIDer
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
	UUIDer
	// Called once per frame to update behaviour.
	Update(win *pixelgl.Window, world *World, dt float64) (toCreate, toDestroy []Entity)
	// Called to get the update layer for this entity.
	// Higher values will be updated first.
	// Should NEVER change after entity has been created.
	UpdateLayer() int
}

type UUIDer interface {
	UUID() string
}

// An object that can be added to a world.
type Entity interface {
	UUIDer
	SetUUID(string)
	AfterAdd(*World)
}

// You may optionally inherit this to give an entity some default behaviour.
type EntityBase struct {
	uuid string
}

func (e *EntityBase) UUID() string {
	return e.uuid
}
func (e *EntityBase) SetUUID(id string) {
	e.uuid = id
}

func (*EntityBase) Update(win *pixelgl.Window, world *World, dt float64) (toCreate, toDestroy []Entity) {
	return nil, nil
}
func (*EntityBase) PreDraw(win *pixelgl.Window)                                        {}
func (*EntityBase) Draw(win *pixelgl.Window, world *World, worldToScreen pixel.Matrix) {}
func (*EntityBase) UpdateLayer() int                                                   { return 0 }
func (*EntityBase) DrawLayer() int                                                     { return 0 }
func (*EntityBase) AfterAdd(*World)                                                    {}

func SetRandomUUID(e Entity) {
	e.SetUUID(uuid.NewString())
}

func ClearUUID(e Entity) {
	e.SetUUID("")
}

func HasUUID(e Entity) bool {
	return e.UUID() != ""
}
