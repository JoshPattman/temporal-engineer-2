package ent

import "github.com/gopxl/pixel"

type Transform interface {
	Position() pixel.Vec
	Angle() float64
}
