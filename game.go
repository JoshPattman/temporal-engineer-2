package main

import (
	"ent"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewGame() Screen {
	world := ent.NewWorld()
	asteroidBatch := NewBatchDraw("asteroid.png")
	world.Add(
		NewCamera(),
		NewStation(),
		NewCompass(),
		NewBackground(),
		asteroidBatch,
		NewPlayer(),
	)
	for range 20 {
		ast := NewAsteroid(asteroidBatch)
		ast.SetPosition(pixel.V(1, 0).Scaled(rand.Float64() * 10).Rotated(rand.Float64() * math.Pi * 2))
		world.Add(ast)
	}
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
		ent.FilterEntitiesByType[CameraTarget](
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
