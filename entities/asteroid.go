package entities

import (
	"ent"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewAsteroid(typ AsteroidType) *Asteroid {
	var batchName, tagName string
	var sprite *pixel.Sprite
	var resources int
	radius := rand.Float64()*1.5 + 0.5
	switch typ {
	case NormalAsteroid:
		batchName = "asteroid_batch"
		tagName = "asteroid"
		sprite = GlobalSpriteManager.FullSprite("asteroid.png")
		resources = 0
	case MineableAsteroid:
		batchName = "mineable_asteroid_batch"
		tagName = "mineable_asteroid"
		sprite = GlobalSpriteManager.FullSprite("asteroid-mineable.png")
		resources = int(radius * 3)
	}
	ast := &Asteroid{
		sprite:    sprite,
		velocity:  pixel.V(0.5, 0).Rotated(rand.Float64() * math.Pi * 2),
		radius:    radius,
		batchName: batchName,
		tagName:   tagName,
		resources: resources,
	}
	ast.SetPosition(pixel.V(rand.Float64()*100, rand.Float64()*100))
	return ast
}

type AsteroidType uint8

const (
	NormalAsteroid AsteroidType = iota
	MineableAsteroid
)

type Asteroid struct {
	ent.CoreEntity
	ent.WithActivePhysics
	ent.WithUpdate
	ent.WithDraw
	batchName string
	tagName   string
	sprite    *pixel.Sprite
	velocity  pixel.Vec
	radius    float64
	resources int
}

// Shape implements ent.ActivePhysicsBody.
func (a *Asteroid) Shape() ent.Shape {
	return ent.Circle{
		Center: a.Position(),
		Radius: a.radius,
	}
}

func (a *Asteroid) Radius() float64 {
	return a.radius
}

func (a *Asteroid) AfterAdd(w *ent.World) {
	w.AddTags(a, a.tagName)
}

func (a *Asteroid) Update(win *pixelgl.Window, entities *ent.World, dt float64) {
	// Check if out of range of player, and delete if so
	player, ok := ent.First(
		ent.OfType[*Player](
			entities.ForTag("player"),
		),
	)
	if ok {
		dist := player.Position().To(a.Position()).Len()
		if dist > 40 {
			entities.Destroy(a)
			return
		}
	}
}

func (a *Asteroid) Draw(win *pixelgl.Window, world *ent.World, worldToScreen pixel.Matrix) {
	batch, ok := ent.First(
		ent.OfType[*BatchDraw](
			world.ForTag(a.batchName),
		),
	)
	if !ok {
		return
	}
	a.sprite.Draw(
		batch.Batch,
		pixel.IM.Scaled(
			pixel.ZV,
			a.radius*2.0/a.sprite.Frame().W(),
		).Chained(
			ent.TransMat(a),
		).Chained(
			worldToScreen,
		),
	)
}
