package main

import (
	"ent"
	"fiz"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ fiz.ActivePhysicsBody = &Asteroid{}

func NewAsteroid(batch *BatchDraw) *Asteroid {
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
	Transform
	batch    *pixel.Batch
	sprite   *pixel.Sprite
	velocity pixel.Vec
	radius   float64
}

// Elasticity implements fiz.ActivePhysicsBody.
func (a *Asteroid) Elasticity() float64 {
	return 0.3
}

// IsPhysicsActive implements fiz.ActivePhysicsBody.
func (a *Asteroid) IsPhysicsActive() bool {
	return true
}

// Mass implements fiz.ActivePhysicsBody.
func (a *Asteroid) Mass() float64 {
	return 1
}

// SetState implements fiz.ActivePhysicsBody.
func (a *Asteroid) SetState(state fiz.BodyState) {
	a.pos = state.Position
	a.velocity = state.Velocity
}

// Shape implements fiz.ActivePhysicsBody.
func (a *Asteroid) Shape() fiz.Shape {
	return fiz.Circle{
		Center: a.pos,
		Radius: a.radius,
	}
}

// State implements fiz.ActivePhysicsBody.
func (a *Asteroid) State() fiz.BodyState {
	return fiz.BodyState{
		Position: a.pos,
		Velocity: a.velocity,
	}
}

func (a *Asteroid) Radius() float64 {
	return a.radius
}

func (a *Asteroid) Update(win *pixelgl.Window, entities *ent.Entities, dt float64) ([]ent.Entity, []ent.Entity) {
	a.SetPosition(
		a.Position().Add(a.velocity.Scaled(dt)),
	)
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

func (a *Asteroid) UpdateLayer() int {
	return 1
}

func (a *Asteroid) DrawLayer() int {
	return 1
}

func (a *Asteroid) Tags() []string {
	return []string{}
}
