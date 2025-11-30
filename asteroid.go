package main

import (
	"ent"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.ActivePhysicsBody = &Asteroid{}
var _ ent.Entity = &Asteroid{}

type AsteroidType uint8

const (
	NormalAsteroid AsteroidType = iota
	MineableAsteroid
)

func NewAsteroid(world *ent.World, typ AsteroidType) *Asteroid {
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
	ast.SetState(ast.State().WithPosition(pixel.V(rand.Float64()*100, rand.Float64()*100)))
	return ast
}

type Asteroid struct {
	ent.MinimalEntity
	ent.MinimalActivePhysicsBody
	ent.MinimalUpdater
	ent.MinimalDraw
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
		Center: a.State().Position,
		Radius: a.radius,
	}
}

func (a *Asteroid) Radius() float64 {
	return a.radius
}

func (a *Asteroid) AfterAdd(w *ent.World) {
	w.AddTags(a, a.tagName)
}

func (a *Asteroid) Update(win *pixelgl.Window, entities *ent.World, dt float64) ([]ent.Entity, []ent.Entity) {
	// Check if out of range of player, and delete if so
	player, ok := ent.First(
		ent.OfType[*Player](
			entities.ForTag("player"),
		),
	)
	if ok {
		dist := player.Position().To(a.State().Position).Len()
		if dist > 40 {
			return nil, []ent.Entity{a}
		}
	}
	return nil, nil
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
			ent.PhysicsBodyMat(a),
		).Chained(
			worldToScreen,
		),
	)
}

var _ ent.Entity = &AsteroidSpawner{}

func NewAsteroidSpawner() *AsteroidSpawner {
	return &AsteroidSpawner{}
}

type AsteroidSpawner struct {
	ent.MinimalEntity
	ent.MinimalUpdater
	timer float64
}

// Update implements ent.Entity.
func (a *AsteroidSpawner) Update(win *pixelgl.Window, world *ent.World, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	player, ok := ent.First(
		ent.OfType[*Player](
			world.ForTag("player"),
		),
	)
	if !ok {
		return nil, nil
	}
	a.timer += dt
	if a.timer > 0.2 {
		a.timer = 0
		var asteroid ent.ActivePhysicsBody

		if rand.Float64() > 0.2 {
			asteroid = NewAsteroid(world, NormalAsteroid)
		} else {
			asteroid = NewAsteroid(world, MineableAsteroid)
		}
		state := asteroid.State()
		state.Velocity = pixel.V(3+rand.Float64()*7, 0).Rotated(rand.Float64() * math.Pi * 2)
		state.Position = player.Position().Add(pixel.V(35, 0).Rotated(rand.Float64() * math.Pi * 2))
		asteroid.SetState(state)
		return []ent.Entity{
			asteroid.(ent.Entity),
		}, nil
	}
	return nil, nil
}
