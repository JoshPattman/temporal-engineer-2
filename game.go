package main

import (
	"ent"
	"math"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewGame() Screen {
	entities := ent.NewEntities()
	entities.Add(NewCamera())
	entities.Add(NewStation())
	entities.Add(NewCompass())
	entities.Add(NewBackground())
	asteroidBatch := NewBatchDraw("asteroid.png")
	entities.Add(asteroidBatch)
	for range 20 {
		ast := NewAsteroid(asteroidBatch)
		ast.SetPosition(pixel.V(1, 0).Scaled(rand.Float64() * 10).Rotated(rand.Float64() * math.Pi * 2))
		entities.Add(ast)
	}
	entities.Add(NewPlayer())
	return &Game{
		entities: entities,
	}
}

type Game struct {
	entities *ent.Entities
}

func (g *Game) Update(win *pixelgl.Window, dt float64) Screen {
	g.entities.Update(win, dt)
	return nil
}

func (g *Game) Draw(win *pixelgl.Window) {
	g.entities.PreDraw(win)

	// Get matrix to transform workd to screen pos
	camMat := pixel.IM.Scaled(pixel.ZV, 20).Moved(win.Bounds().Center())
	camera, ok := ent.First(
		ent.FilterEntitiesByType[CameraTarget](
			g.entities.ForTag("camera"),
		),
	)
	if ok {
		camMat = pixel.IM.Moved(camera.Position().Scaled(-1)).Chained(camMat)
	}

	// Draw all objects
	win.Clear(pixel.RGB(0.01, 0.01, 0.05))
	for e := range g.entities.ByDraw() {
		e.Draw(win, camMat)
	}
}
