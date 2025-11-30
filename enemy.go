package main

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewEnemy() *Enemy {
	enemy := &Enemy{
		sprite: GlobalSpriteManager.TiledSprites("enemy.png", 16, []TilePos{{2, 0}})[0],
	}
	enemy.SetState(enemy.State().WithPosition(pixel.V(3, 0)))
	return enemy
}

type Enemy struct {
	ent.MinimalEntity
	ent.MinimalActivePhysicsBody
	ent.MinimalDraw
	sprite *pixel.Sprite
}

func (e *Enemy) Shape() ent.Shape {
	return ent.Circle{Center: e.State().Position, Radius: 1}
}

func (e *Enemy) Draw(win *pixelgl.Window, world *ent.World, worldToScreen pixel.Matrix) {
	e.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.0/16.0).Chained(ent.PhysicsBodyMat(e)).Chained(worldToScreen))
}
