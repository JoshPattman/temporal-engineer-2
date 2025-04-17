package main

import (
	"ent"
	"fiz"
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

type CollisionListener interface {
	OnCollision(fiz.Collision)
}

func (g *Game) Update(win *pixelgl.Window, dt float64) Screen {
	allToCreate := []ent.Entity{}
	allToRemove := []ent.Entity{}
	for _, e := range g.entities.ByUpdate() {
		toCreate, toRemove := e.Update(win, g.entities, dt)
		allToCreate = append(allToCreate, toCreate...)
		allToRemove = append(allToRemove, toRemove...)
	}
	for _, e := range allToCreate {
		if !g.entities.Has(e) {
			g.entities.Add(e)
		}
	}
	for _, e := range allToRemove {
		if g.entities.Has(e) {
			g.entities.Remove(e)
		}
	}

	fizBodies := ent.FilterEntitiesByType[fiz.PhysicsBody](g.entities.ByUpdate())
	bs := make([]fiz.PhysicsBody, 0)
	for _, b := range fizBodies {
		bs = append(bs, b)
	}
	cols := fiz.StatelessCollisionPhysics(bs)

	for _, col := range cols {
		self, ok := col.Self.(CollisionListener)
		if ok {
			self.OnCollision(col)
		}
		col = col.ForOther()
		self, ok = col.Self.(CollisionListener)
		if ok {
			self.OnCollision(col)
		}
	}

	return nil
}

func (g *Game) Draw(win *pixelgl.Window) {
	win.Clear(pixel.RGB(0.01, 0.01, 0.05))
	camMat := pixel.IM.Scaled(pixel.ZV, 20).Moved(win.Bounds().Center())
	for _, cam := range g.entities.ForTag("camera") {
		cam, ok := cam.(CameraTarget)
		if !ok {
			continue
		}
		camMat = pixel.IM.Moved(cam.Position().Scaled(-1)).Chained(camMat)
	}

	for _, e := range g.entities.ByDraw() {
		e.Draw(win, camMat)
	}
}
