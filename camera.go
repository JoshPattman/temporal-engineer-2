package main

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.Entity = &Camera{}

func NewCamera() *Camera {
	return &Camera{}
}

type Camera struct {
	pos pixel.Vec
}

func (c *Camera) Position() pixel.Vec {
	return c.pos
}

type CameraTarget interface {
	Position() pixel.Vec
}

// Draw implements ent.Entity.
func (c *Camera) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {}

// DrawLayer implements ent.Entity.
func (c *Camera) DrawLayer() int {
	return 0
}

// Tags implements ent.Entity.
func (c *Camera) Tags() []string {
	return []string{"camera"}
}

// Update implements ent.Entity.
func (c *Camera) Update(win *pixelgl.Window, all *ent.Entities, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	targets := all.ForTag("player_camera_target")
	var target CameraTarget = c
	for _, t := range targets {
		t, ok := t.(CameraTarget)
		if ok {
			target = t
			break
		}
	}
	c.pos = pixel.Lerp(c.pos, target.Position(), 0.05)
	return nil, nil
}

// UpdateLayer implements ent.Entity.
func (c *Camera) UpdateLayer() int {
	return 5
}
