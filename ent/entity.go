package ent

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

type Entity interface {
	Update(win *pixelgl.Window, all *Entities, dt float64) (toCreate, toDestroy []Entity)
	Draw(win *pixelgl.Window, worldToScreen pixel.Matrix)
	Tags() []string
	UpdateLayer() int
	DrawLayer() int
}
