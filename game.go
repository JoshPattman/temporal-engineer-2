package main

import (
	"ent"
	"te2/entities"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewGame() Screen {
	world := ent.NewWorld()
	world.Add(
		entities.NewCamera(),
		entities.NewStation(),
		entities.NewCompass(),
		entities.NewBackground(),
		entities.NewBatchDraw("asteroid.png", "asteroid_batch"),
		entities.NewBatchDraw("asteroid-mineable.png", "mineable_asteroid_batch"),
		entities.NewPlayer(),
		entities.NewAsteroidSpawner(),
		entities.NewSheildsIndicator(),
		entities.NewMineralsIndicator(),
		entities.NewEnemy(),
	)
	return &Game{
		world: world,
	}
}

type Game struct {
	world *ent.World
}

func (g *Game) Update(win *pixelgl.Window, dt float64) Screen {
	g.world.Update(win, dt)
	return nil
}

func (g *Game) Draw(win *pixelgl.Window) {
	// Get matrix to transform workd to screen pos
	camMat := pixel.IM.Scaled(pixel.ZV, 20).Moved(win.Bounds().Center())
	camera, ok := ent.First(
		ent.OfType[entities.CameraTarget](
			g.world.ForTag("camera"),
		),
	)
	if ok {
		camMat = pixel.IM.Moved(camera.Position().Scaled(-1)).Chained(camMat)
	}

	// Draw all objects
	win.Clear(pixel.RGB(0.01, 0.01, 0.05))
	g.world.Draw(win, camMat)
}
