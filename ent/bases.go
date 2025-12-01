package ent

import (
	"github.com/google/uuid"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// Compose to provide basic behaviour to implement Entity.
type CoreEntity struct {
	uuid EntityUUID
}

func (e *CoreEntity) UUID() EntityUUID {
	if e.uuid == "" {
		e.uuid = EntityUUID(uuid.NewString())
	}
	return e.uuid
}

func (e *CoreEntity) AfterAdd(*World) {}

func (e *CoreEntity) HandleMessage(*World, any) {}

// Compose additionally with MinimalEntity to provide basic behaviour to implement Drawer.
type WithDraw struct{}

func (e *WithDraw) PreDraw(win *pixelgl.Window)                                        {}
func (e *WithDraw) Draw(win *pixelgl.Window, world *World, worldToScreen pixel.Matrix) {}
func (e *WithDraw) DrawLayer() int                                                     { return 0 }

// Compose additionally with MinimalEntity to provide basic behaviour to implement Updater.
type WithUpdate struct{}

func (e *WithUpdate) Update(win *pixelgl.Window, world *World, dt float64) {}
func (e *WithUpdate) UpdateLayer() int                                     { return 0 }

// Compose additionally with MinimalEntity to provide basic behaviour to implement transform.
type WithTransform struct {
	position pixel.Vec
	angle    float64
}

func (t *WithTransform) Position() pixel.Vec {
	return t.position
}

func (t *WithTransform) Angle() float64 {
	return t.angle
}

func (t *WithTransform) SetPosition(p pixel.Vec) {
	t.position = p
}

func (t *WithTransform) SetAngle(a float64) {
	t.angle = a
}

// Compose additionally with MinimalEntity to provide basic behaviour to implement PhysicsBody.
type WithStaticPhysics struct {
	WithTransform
	velocity        pixel.Vec
	angularVelocity float64
}

func (e *WithStaticPhysics) Velocity() pixel.Vec {
	return e.velocity
}
func (e *WithStaticPhysics) SetVelocity(v pixel.Vec) {
	e.velocity = v
}
func (e *WithStaticPhysics) AngularVelocity() float64 {
	return e.angularVelocity
}
func (e *WithStaticPhysics) SetAngularVelocity(v float64) {
	e.angularVelocity = v
}
func (e *WithStaticPhysics) Shape() Shape        { return Circle{e.Position(), 1} }
func (e *WithStaticPhysics) Elasticity() float64 { return 0.3 }

// Compose additionally with MinimalEntity to provide basic behaviour to implement ActivePhysicsBody.
type WithActivePhysics struct {
	WithStaticPhysics
}

func (e *WithStaticPhysics) Mass() float64         { return 1 }
func (e *WithStaticPhysics) IsPhysicsActive() bool { return true }
func (e *WithStaticPhysics) PysicsUpdate(dt float64) {
	EulerStateUpdate(e, BodyEffects{}, dt)
}
