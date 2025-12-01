package entities

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewExplosion(at pixel.Vec, scale float64) *Explosion {
	sprites := GlobalSpriteManager.TiledSprites(
		"boom.png",
		36,
		[]TilePos{
			{0, 1},
			{1, 1},
			{2, 1},
			{3, 1},
			{4, 1},
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
			{4, 0},
		},
	)
	return &Explosion{
		pos:     at,
		timer:   0,
		sprites: sprites,
		scale:   scale,
	}
}

var _ ent.Entity = &Explosion{}

type Explosion struct {
	ent.CoreEntity
	ent.WithDraw
	ent.WithUpdate
	pos     pixel.Vec
	timer   float64
	sprites []*pixel.Sprite
	scale   float64
}

// Draw implements ent.Entity.
func (e *Explosion) Draw(win *pixelgl.Window, _ *ent.World, worldToScreen pixel.Matrix) {
	idx := int(e.timer / 0.5 * float64(len(e.sprites)))
	s := e.sprites[idx]
	s.Draw(
		win,
		pixel.IM.Scaled(pixel.ZV, 0.1*e.scale).Moved(e.pos).Chained(worldToScreen),
	)
}

func (e *Explosion) DrawLayer() int { return -1 }

// Update implements ent.Entity.
func (e *Explosion) Update(win *pixelgl.Window, all *ent.World, dt float64) {
	e.timer += dt
	if e.timer >= 0.5 {
		all.Remove(e)
	}
}

func (e *Explosion) Position() pixel.Vec {
	return e.pos
}
