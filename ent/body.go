package ent

import (
	"math"

	"github.com/gopxl/pixel"
)

// Describes the state of a physics-enabled body.
// This is applicable to both active and kinematic bodies.
type BodyState struct {
	Position        pixel.Vec
	Velocity        pixel.Vec
	Angle           float64
	AngularVelocity float64
}

// A body that has physics and the ability to affect other bodies in the world.
type PhysicsBody interface {
	UUIDer
	State() BodyState
	Shape() Shape
	Elasticity() float64
}

// A body that not only has physics, but is able to react to other bodies in the world.
type ActivePhysicsBody interface {
	PhysicsBody
	SetState(BodyState)
	Mass() float64
	IsPhysicsActive() bool
}

// Calculate the force of drag by summing the natural and linear drag forces, scaled by their multipliers.
// Natural drag is drag that follows the power of two rule.
// Linear drag is simplified drag that opposes motion linearly.
// Mixing the two can make your game feel nicest.
func CalculateDragForce(velocity pixel.Vec, naturalDrag, linearDrag float64) pixel.Vec {
	l := velocity.Len()
	natural := velocity.Scaled(-l * naturalDrag)
	linear := velocity.Scaled(-linearDrag)
	return natural.Add(linear)
}

// Calculate the torque of drag by summing the natural and linear drag forces, scaled by their multipliers.
// Natural drag is drag that follows the power of two rule.
// Linear drag is simplified drag that opposes motion linearly.
// Mixing the two can make your game feel nicest.
func CalculateDragTorque(angularVelocity float64, naturalDrag, linearDrag float64) float64 {
	natural := angularVelocity * -math.Abs(angularVelocity) * naturalDrag
	linear := angularVelocity * -linearDrag
	return natural + linear
}

// The forces and torques that are to be applied to an active physics body.
type BodyEffects struct {
	Force   pixel.Vec
	Impulse pixel.Vec
	Torque  float64
}

// Update an active physics body using euler rules with some effects and a time interval.
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
