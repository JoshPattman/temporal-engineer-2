package fiz

import "github.com/gopxl/pixel"

type BodyState struct {
	Position pixel.Vec
	Velocity pixel.Vec
}

type PhysicsBody interface {
	State() BodyState
	Shape() Shape
	Elasticity() float64
}

type ActivePhysicsBody interface {
	PhysicsBody
	SetState(BodyState)
	Mass() float64
	IsPhysicsActive() bool
}
