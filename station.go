package main

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.Entity = &Station{}

func NewStation() *Station {
	sprites := GlobalSpriteManager.TiledSprites("station.png", 32, []TilePos{
		{0, 0},
		{1, 0},
		{2, 0},
	})
	return &Station{
		sprites: sprites,
	}
}

type Station struct {
	ent.EntityBase
	pos         pixel.Vec
	sprites     []*pixel.Sprite
	spriteTimer float64
}

// Draw implements ent.Entity.
func (s *Station) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	spriteIdx := int(s.spriteTimer*0.5) % len(s.sprites)
	s.sprites[spriteIdx].Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Chained(worldToScreen))
}

// DrawLayer implements ent.Entity.
func (s *Station) DrawLayer() int {
	return 1
}

// Tags implements ent.Entity.
func (s *Station) Tags() []string {
	return []string{"station"}
}

// Update implements ent.Entity.
func (s *Station) Update(win *pixelgl.Window, all *ent.World, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	s.spriteTimer += dt
	return nil, nil
}
