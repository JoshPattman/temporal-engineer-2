package main

import (
	"ent"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.ActivePhysicsBody = &Asteroid{}
var _ ent.Entity = &Asteroid{}

func NewAsteroid(world *ent.World) *Asteroid {
	batch, ok := ent.First(
		ent.FilterEntitiesByType[*BatchDraw](
			world.ForTag("asteroid_batch"),
		),
	)
	if !ok {
		panic("cannot add asteroid without an asteroid batch present")
	}
	sprite := GlobalSpriteManager.FullSprite("asteroid.png")
	return &Asteroid{
		Transform: Transform{
			pos: pixel.V(rand.Float64()*100, rand.Float64()*100),
		},
		sprite:   sprite,
		velocity: pixel.V(0.5, 0).Rotated(rand.Float64() * math.Pi * 2),
		radius:   rand.Float64()*1.5 + 0.5,
		batch:    batch.Batch,
	}
}

type Asteroid struct {
	ent.EntityBase
	Transform
	batch    *pixel.Batch
	sprite   *pixel.Sprite
	velocity pixel.Vec
	radius   float64
}

// Elasticity implements ent.ActivePhysicsBody.
func (a *Asteroid) Elasticity() float64 {
	return 0.3
}

// IsPhysicsActive implements ent.ActivePhysicsBody.
func (a *Asteroid) IsPhysicsActive() bool {
	return true
}

// Mass implements ent.ActivePhysicsBody.
func (a *Asteroid) Mass() float64 {
	return 1
}

// SetState implements ent.ActivePhysicsBody.
func (a *Asteroid) SetState(state ent.BodyState) {
	a.pos = state.Position
	a.velocity = state.Velocity
	a.rot = state.Angle
}

// Shape implements ent.ActivePhysicsBody.
func (a *Asteroid) Shape() ent.Shape {
	return ent.Circle{
		Center: a.pos,
		Radius: a.radius,
	}
}

// State implements ent.ActivePhysicsBody.
func (a *Asteroid) State() ent.BodyState {
	return ent.BodyState{
		Position: a.pos,
		Velocity: a.velocity,
		Angle:    a.rot,
	}
}

func (a *Asteroid) Radius() float64 {
	return a.radius
}

func (a *Asteroid) Update(win *pixelgl.Window, entities *ent.World, dt float64) ([]ent.Entity, []ent.Entity) {
	fx := ent.BodyEffects{}
	fx.Force = ent.CalculateDragForce(a.velocity, 0.5, 0)
	ent.EulerStateUpdate(a, fx, dt)
	return nil, nil
}

func (a *Asteroid) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	a.sprite.Draw(
		a.batch,
		pixel.IM.Scaled(
			pixel.ZV,
			a.radius*2.0/a.sprite.Frame().W(),
		).Chained(
			a.Mat(),
		).Chained(
			worldToScreen,
		),
	)
}

var _ ent.Entity = &AsteroidSpawner{}

func NewAsteroidSpawner() *AsteroidSpawner {
	return &AsteroidSpawner{}
}

type AsteroidSpawner struct {
	ent.EntityBase
	timer float64
}

// Update implements ent.Entity.
func (a *AsteroidSpawner) Update(win *pixelgl.Window, world *ent.World, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	a.timer += dt
	if a.timer > 1 {
		a.timer = 0
		return []ent.Entity{
			NewAsteroid(world),
		}, nil
	}
	return nil, nil
}
