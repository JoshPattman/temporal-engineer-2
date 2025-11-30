package ent

import (
	"math"

	"github.com/gopxl/pixel"
)

type Transform interface {
	Position() pixel.Vec
	Angle() float64
}

type ActiveTransform interface {
	Transform
	SetPosition(pixel.Vec)
	SetAngle(float64)
}

type DynamicTransform interface {
	Transform
	Velocity() pixel.Vec
	AngularVelocity() float64
}

type ActiveDynamicTransform interface {
	ActiveTransform
	DynamicTransform
	SetVelocity(pixel.Vec)
	SetAngularVelocity(float64)
}

func Forward(a Transform) pixel.Vec {
	return pixel.V(1, 0).Rotated(a.Angle())
}

func TransMat(t Transform) pixel.Matrix {
	return pixel.IM.Rotated(pixel.ZV, t.Angle()).Moved(t.Position())
}

func VelocityAt(transform DynamicTransform, pt pixel.Vec) pixel.Vec {
	if pt == transform.Position() {
		return transform.Velocity()
	}
	dist := transform.Position().To(pt).Len()
	dir := transform.Position().To(pt).Scaled(1 / dist)
	distPerRad := dist // love maths
	return transform.Velocity().Add(dir.Rotated(math.Pi / 2).Scaled(transform.AngularVelocity() * distPerRad))
}
