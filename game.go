package main

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewGame() Screen {
	world := ent.NewWorld()
	world.Add(
		NewCamera(),
		NewStation(),
		NewCompass(),
		NewBackground(),
		NewBatchDraw("asteroid.png", "asteroid_batch"),
		NewBatchDraw("asteroid-mineable.png", "mineable_asteroid_batch"),
		NewPlayer(),
		NewAsteroidSpawner(),
		NewSheildsIndicator(),
		NewMineralsIndicator(),
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
		ent.OfType[CameraTarget](
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
