package ent

import (
	"math"

	"github.com/gopxl/pixel"
)

// A sum type of different shapes supported by the physics engine.
type Shape interface {
	shape()
	EffectArea() (pixel.Vec, float64)
}

func (Circle) shape()     {}
func (Line) shape()       {}
func (MultiShape) shape() {}

// A circle shape.
type Circle struct {
	Center pixel.Vec
	Radius float64
}

func (c Circle) EffectArea() (pixel.Vec, float64) {
	return c.Center, c.Radius
}

type MultiShape struct {
	Shapes []Shape
}

func (ms MultiShape) EffectArea() (pixel.Vec, float64) {
	if len(ms.Shapes) == 0 {
		panic("must have at least one shape in a multi shape")
	}
	// Quick algorithm - find the encomapssing rect, then create a circle that encompasses that
	minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, s := range ms.Shapes {
		sCenter, sRad := s.EffectArea()
		minX = min(minX, sCenter.X-sRad)
		minY = min(minY, sCenter.Y-sRad)
		maxX = max(maxX, sCenter.X+sRad)
		maxY = max(maxY, sCenter.Y+sRad)
	}
	rect := pixel.R(minX, minY, maxX, maxY)
	return rect.Center(), rect.Size().Len() / 2
}

// A line shape.
type Line struct {
	A pixel.Vec
	B pixel.Vec
}

func (l Line) EffectArea() (pixel.Vec, float64) {
	return l.A.Add(l.B).Scaled(0.5), l.A.To(l.B).Len() / 2
}

// A collision of two shapes.
type shapeCollision struct {
	collided bool
	normal   pixel.Vec
	overlap  float64
	point    pixel.Vec
}

// Compute the collision of two shapes of any type.
func collideShapes(a, b Shape) shapeCollision {
	// Quick culling check before typecasting
	// Are the shapes too far apart?
	p1, r1 := a.EffectArea() // TODO: Its a bit inneficiant to keep computing these as they wont change
	p2, r2 := b.EffectArea()
	dist2 := p1.To(p2).SqLen()
	if dist2 > math.Pow(r1+r2, 2) {
		return shapeCollision{}
	}
	// Collide shapes based on type
	var col shapeCollision
	var ok bool
	col, ok = checkAndCollideSym(a, b, collideCircleCircle)
	if ok {
		return col
	}
	col, ok = checkAndCollideAsym(a, b, collideMultiOther)
	if ok {
		return col
	}
	col, ok = checkAndCollideAsym(a, b, collideCircleLine)
	if ok {
		return col
	}
	panic("collision not supported between those")
}

// Helper function to run the collision function if the shapes are of the correct type.
func checkAndCollideSym[T Shape](a Shape, b Shape, f func(T, T) shapeCollision) (shapeCollision, bool) {
	aT, ok := a.(T)
	if !ok {
		return shapeCollision{}, false
	}
	bU, ok := b.(T)
	if !ok {
		return shapeCollision{}, false
	}
	return f(aT, bU), true
}

// Helper function to run the collision function if the shapes are of the correct type.
func checkAndCollideAsym[T, U Shape](a Shape, b Shape, f func(T, U) shapeCollision) (shapeCollision, bool) {
	isRightTypes := true
	aT, ok := a.(T)
	if !ok {
		isRightTypes = false
	}
	bU, ok := b.(U)
	if !ok {
		isRightTypes = false
	}
	if isRightTypes {
		return f(aT, bU), true
	}
	isRightTypes = true
	aU, ok := a.(U)
	if !ok {
		isRightTypes = false
	}
	bT, ok := b.(T)
	if !ok {
		isRightTypes = false
	}
	if isRightTypes {
		col := f(bT, aU)
		col.normal = col.normal.Scaled(-1)
		return col, true
	}
	return shapeCollision{}, false
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
	point := a.Center.Add(normal.Scaled(a.Radius - overlapDist/2))
	return shapeCollision{
		collided: true,
		normal:   normal,
		overlap:  overlapDist,
		point:    point,
	}
}

func collideMultiOther(a MultiShape, b Shape) shapeCollision {
	for _, shape := range a.Shapes {
		col := collideShapes(shape, b)
		if col.collided {
			return col
		}
	}
	return shapeCollision{}
}

func collideCircleLine(a Circle, b Line) shapeCollision {
	closestLinePoint := pixel.Line(b).Closest(a.Center)
	dist := a.Center.To(closestLinePoint).Len()
	if dist > a.Radius {
		return shapeCollision{}
	}
	overlapDist := a.Radius - dist
	normal := a.Center.To(closestLinePoint).Scaled(1.0 / dist)
	return shapeCollision{
		collided: true,
		normal:   normal,
		overlap:  overlapDist,
		point:    closestLinePoint,
	}
}
