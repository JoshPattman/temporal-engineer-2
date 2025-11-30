package ent

import (
	"math"

	"github.com/gopxl/pixel"
)

// A body that has physics and the ability to affect other bodies in the world.
type PhysicsBody interface {
	UUIDer
	DynamicTransform
	Shape() Shape
	Elasticity() float64
}

// A body that not only has physics, but is able to react to other bodies in the world.
type ActivePhysicsBody interface {
	PhysicsBody
	EulerUpdateable
	IsPhysicsActive() bool
	PysicsUpdate(dt float64)
}

type EulerUpdateable interface {
	ActiveDynamicTransform
	Mass() float64
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
func EulerStateUpdate(body EulerUpdateable, effects BodyEffects, dt float64) {
	acceleration := effects.Force.Scaled(1.0 / body.Mass())
	body.SetVelocity(body.Velocity().Add(acceleration.Scaled(dt).Add(effects.Impulse)))
	body.SetPosition(body.Position().Add(body.Velocity().Scaled(dt)))
	angularAcceleration := effects.Torque / body.Mass()
	body.SetAngularVelocity(body.AngularVelocity() + angularAcceleration*dt)
	body.SetAngle(body.Angle() + body.AngularVelocity()*dt)
}
