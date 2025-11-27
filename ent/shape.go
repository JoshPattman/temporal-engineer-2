package ent

import (
	"math"

	"github.com/gopxl/pixel"
)

// A sum type of different shapes supported by the physics engine.
type Shape interface {
	shape()
}

// A circle shape.
type Circle struct {
	Center pixel.Vec
	Radius float64
}

func (Circle) shape() {}

// A collision of two shapes.
type shapeCollision struct {
	collided bool
	normal   pixel.Vec
	overlap  float64
}

// Compute the collision of two shapes of any type.
func collideShapes(a, b Shape) shapeCollision {
	var col shapeCollision
	var ok bool
	col, ok = checkAndCollide(a, b, collideCircleCircle)
	if ok {
		return col
	}
	panic("collision not supported between those")
}

// Helper function to run the collision function if the shapes are of the correct type.
func checkAndCollide[T, U Shape](a Shape, b Shape, f func(T, U) shapeCollision) (shapeCollision, bool) {
	aT, ok := a.(T)
	if !ok {
		return shapeCollision{}, false
	}
	bU, ok := b.(U)
	if !ok {
		return shapeCollision{}, false
	}
	return f(aT, bU), true
}

// helper function to collide two circles.
func collideCircleCircle(a, b Circle) shapeCollision {
	centerDelta := a.Center.To(b.Center)
	centerDist2 := centerDelta.SqLen()
	touchDist := a.Radius + b.Radius
	if centerDist2 >= touchDist*touchDist {
		return shapeCollision{
			collided: false,
		}
	}
	centerDist := math.Sqrt(centerDist2)
	overlapDist := touchDist - centerDist
	normal := centerDelta.Scaled(1.0 / centerDist)
	return shapeCollision{
		collided: true,
		normal:   normal,
		overlap:  overlapDist,
	}
}
