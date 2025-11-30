package ent

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// Compose to provide basic behaviour to implement Entity.
type MinimalEntity struct {
	uuid string
}

func (e *MinimalEntity) UUID() string {
	return e.uuid
}
func (e *MinimalEntity) SetUUID(id string) {
	e.uuid = id
}
func (e *MinimalEntity) AfterAdd(*World) {}

// Compose additionally with MinimalEntity to provide basic behaviour to implement Drawer.
type MinimalDraw struct{}

func (e *MinimalDraw) PreDraw(win *pixelgl.Window)                                        {}
func (e *MinimalDraw) Draw(win *pixelgl.Window, world *World, worldToScreen pixel.Matrix) {}
func (e *MinimalDraw) DrawLayer() int                                                     { return 0 }

// Compose additionally with MinimalEntity to provide basic behaviour to implement Updater.
type MinimalUpdater struct{}

func (e *MinimalUpdater) Update(win *pixelgl.Window, world *World, dt float64) (toCreate, toDestroy []Entity) {
	return nil, nil
}
func (e *MinimalUpdater) UpdateLayer() int { return 0 }

// Compose additionally with MinimalEntity to provide basic behaviour to implement PhysicsBody.
type MinimalPhysicsBody struct {
	state BodyState
}

func (e *MinimalPhysicsBody) State() BodyState    { return e.state }
func (e *MinimalPhysicsBody) Shape() Shape        { return Circle{e.State().Position, 1} }
func (e *MinimalPhysicsBody) Elasticity() float64 { return 0.3 }
func (e *MinimalPhysicsBody) Position() pixel.Vec { return e.state.Position }
func (e *MinimalPhysicsBody) Angle() float64      { return e.state.Angle }

// Compose additionally with MinimalEntity to provide basic behaviour to implement ActivePhysicsBody.
type MinimalActivePhysicsBody struct {
	MinimalPhysicsBody
}

func (e *MinimalPhysicsBody) SetState(s BodyState)  { e.state = s }
func (e *MinimalPhysicsBody) Mass() float64         { return 1 }
func (e *MinimalPhysicsBody) IsPhysicsActive() bool { return true }
func (e *MinimalPhysicsBody) PysicsUpdate(dt float64) {
	EulerStateUpdate(e, BodyEffects{}, dt)
}
