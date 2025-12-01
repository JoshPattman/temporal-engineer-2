package entities

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewCamera() *Camera {
	return &Camera{}
}

type Camera struct {
	ent.CoreEntity
	ent.WithUpdate
	pos pixel.Vec
}

func (c *Camera) Position() pixel.Vec {
	return c.pos
}

type CameraTarget interface {
	Position() pixel.Vec
}

func (c *Camera) AfterAdd(world *ent.World) {
	world.AddTags(c, "camera")
}

// Update implements ent.Entity.
func (c *Camera) Update(win *pixelgl.Window, all *ent.World, dt float64) {
	target, ok := ent.First(
		ent.OfType[CameraTarget](
			all.WithTag("player_camera_target"),
		),
	)
	if !ok {
		target = c
	}
	c.pos = pixel.Lerp(c.pos, target.Position(), 0.05)
}

// UpdateLayer implements ent.Entity.
func (c *Camera) UpdateLayer() int {
	return 5
}
