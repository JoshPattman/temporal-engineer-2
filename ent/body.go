package ent

import (
	"math"

	"github.com/gopxl/pixel"
)

type BodyState struct {
	Position        pixel.Vec
	Velocity        pixel.Vec
	Angle           float64
	AngularVelocity float64
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

func CalculateDragForce(velocity pixel.Vec, naturalDrag, linearDrag float64) pixel.Vec {
	l := velocity.Len()
	natural := velocity.Scaled(-l * naturalDrag)
	linear := velocity.Scaled(-linearDrag)
	return natural.Add(linear)
}

func CalculateDragTorque(angularVelocity float64, naturalDrag, linearDrag float64) float64 {
	natural := angularVelocity * -math.Abs(angularVelocity) * naturalDrag
	linear := angularVelocity * -linearDrag
	return natural + linear
}

type BodyEffects struct {
	Force   pixel.Vec
	Impulse pixel.Vec
	Torque  float64
}

func EulerStateUpdate(body ActivePhysicsBody, effects BodyEffects, dt float64) {
	state := body.State()
	acceleration := effects.Force.Scaled(1.0 / body.Mass())
	state.Velocity = state.Velocity.Add(acceleration.Scaled(dt).Add(effects.Impulse))
	state.Position = state.Position.Add(state.Velocity.Scaled(dt))
	angularAcceleration := effects.Torque / body.Mass()
	state.AngularVelocity += angularAcceleration * dt
	state.Angle += state.AngularVelocity * dt
	body.SetState(state)
}
