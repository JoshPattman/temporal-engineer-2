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
	ent.EntityBase
	pos pixel.Vec
}

func (c *Camera) Position() pixel.Vec {
	return c.pos
}

type CameraTarget interface {
	Position() pixel.Vec
}

// Tags implements ent.Entity.
func (c *Camera) Tags() []string {
	return []string{"camera"}
}

// Update implements ent.Entity.
func (c *Camera) Update(win *pixelgl.Window, all *ent.Entities, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	target, ok := ent.First(
		ent.FilterEntitiesByType[CameraTarget](
			all.ForTag("player_camera_target"),
		),
	)
	if !ok {
		target = c
	}
	c.pos = pixel.Lerp(c.pos, target.Position(), 0.05)
	return nil, nil
}

// UpdateLayer implements ent.Entity.
func (c *Camera) UpdateLayer() int {
	return 5
}
