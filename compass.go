package main

import (
	"ent"
	"math"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.Entity = &Compass{}

func NewCompass() *Compass {
	sprite := GlobalSpriteManager.FullSprite("compass.png")
	return &Compass{
		sprite: sprite,
	}
}

type Compass struct {
	sprite *pixel.Sprite
	angle  float64
}

// Draw implements ent.Entity.
func (c *Compass) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	c.sprite.Draw(win, pixel.IM.Rotated(
		pixel.ZV, c.angle+math.Pi/2,
	).Scaled(
		pixel.ZV, 5,
	).Moved(
		pixel.V(64, 64),
	),
	)
}

// DrawLayer implements ent.Entity.
func (c *Compass) DrawLayer() int {
	return -10
}

// Tags implements ent.Entity.
func (c *Compass) Tags() []string {
	return nil
}

// Update implements ent.Entity.
func (c *Compass) Update(win *pixelgl.Window, all *ent.Entities, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	players := all.ForTag("player")
	var player *Player
	for _, p := range players {
		if p, ok := p.(*Player); ok {
			player = p
			break
		}
	}
	if player == nil {
		return nil, nil
	}
	c.angle = player.pos.Angle()
	return nil, nil
}

// UpdateLayer implements ent.Entity.
func (c *Compass) UpdateLayer() int {
	return -10
}
