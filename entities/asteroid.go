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
		toMiners:  ent.NewBus(),
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
	toMiners  *ent.Bus
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

func (a *Asteroid) ToMiners() *ent.Bus { return a.toMiners }

type MineAsteroid struct {
	From pixel.Vec
}

type AsteroidDestroyed struct{}

type AsteroidOutOfRange struct{}

type CheckOutOfMiningRange struct {
	From    pixel.Vec
	MaxDist float64
}

func (a *Asteroid) HandleMessage(world *ent.World, msg any) {
	switch msg := msg.(type) {
	case MineAsteroid:
		a.resources--
		destroy := a.resources <= 0
		if destroy {
			ent.Emit(world, a.toMiners, AsteroidDestroyed{})
			world.Destroy(a)
			world.Instantiate(NewExplosion(a.Position(), a.Radius()))
		} else {
			edgePos := a.Position().To(msg.From).Unit().Scaled(a.radius).Add(a.Position())
			world.Instantiate(NewExplosion(edgePos, 0.3))
		}
	case CheckOutOfMiningRange:
		if a.Position().To(msg.From).Len() > msg.MaxDist {
			ent.Emit(world, a.toMiners, AsteroidOutOfRange{})
			ent.UnsubscribeAll(a.toMiners)
		}
	}
}
