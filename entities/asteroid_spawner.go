package entities

import (
	"ent"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewAsteroidSpawner() *AsteroidSpawner {
	return &AsteroidSpawner{}
}

type AsteroidSpawner struct {
	ent.CoreEntity
	ent.WithUpdate
	timer float64
}

// Update implements ent.Entity.
func (a *AsteroidSpawner) Update(win *pixelgl.Window, world *ent.World, dt float64) {
	player, ok := ent.First(
		ent.OfType[*Player](
			world.ForTag("player"),
		),
	)
	if !ok {
		return
	}
	a.timer += dt
	if a.timer > 0.2 {
		a.timer = 0
		var asteroid *Asteroid

		if rand.Float64() > 0.2 {
			asteroid = NewAsteroid(NormalAsteroid)
		} else {
			asteroid = NewAsteroid(MineableAsteroid)
		}
		asteroid.SetVelocity(pixel.V(3+rand.Float64()*7, 0).Rotated(rand.Float64() * math.Pi * 2))
		asteroid.SetPosition(player.Position().Add(pixel.V(35, 0).Rotated(rand.Float64() * math.Pi * 2)))
		world.Instantiate(asteroid)
	}
}
