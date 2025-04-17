package main

import "github.com/gopxl/pixel"

type Transform struct {
	pos pixel.Vec
	rot float64
}

func (t *Transform) Position() pixel.Vec {
	return t.pos
}

func (t *Transform) SetPosition(p pixel.Vec) {
	t.pos = p
}

func (t *Transform) Rotation() float64 {
	return t.rot
}

func (t *Transform) SetRotation(r float64) {
	t.rot = r
}

func (t *Transform) Forward() pixel.Vec {
	return pixel.V(1, 0).Rotated(t.rot)
}

func (t *Transform) Mat() pixel.Matrix {
	return pixel.IM.Rotated(
		pixel.ZV, t.Rotation(),
	).Moved(
		t.Position(),
	)
}
