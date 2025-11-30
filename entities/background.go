package entities

import (
	"ent"
	"math"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var _ ent.Entity = &Background{}

func NewBackground() *Background {
	s1 := GlobalSpriteManager.FullSprite("background001.png")
	s2 := GlobalSpriteManager.FullSprite("background002.png")
	s3 := GlobalSpriteManager.FullSprite("background003.png")
	return &Background{
		l1: s1,
		l2: s2,
		l3: s3,
		b1: pixel.NewBatch(&pixel.TrianglesData{}, s1.Picture()),
		b2: pixel.NewBatch(&pixel.TrianglesData{}, s2.Picture()),
		b3: pixel.NewBatch(&pixel.TrianglesData{}, s3.Picture()),
	}
}

type Background struct {
	ent.CoreEntity
	ent.WithDraw
	l1 *pixel.Sprite
	l2 *pixel.Sprite
	l3 *pixel.Sprite
	b1 *pixel.Batch
	b2 *pixel.Batch
	b3 *pixel.Batch
}

func drawLevel(s *pixel.Sprite, b *pixel.Batch, scale float64, win *pixelgl.Window, worldToScreen pixel.Matrix) {
	b.Clear()
	min := worldToScreen.Unproject(win.Bounds().Min)
	max := worldToScreen.Unproject(win.Bounds().Max)
	tileWorldWidth := 256.0 / 16.0 * scale
	minTile := min.Scaled(1.0 / tileWorldWidth)
	maxTile := max.Scaled(1.0 / tileWorldWidth)
	for x := math.Floor(minTile.X); x <= math.Ceil(maxTile.X); x += 1 {
		for y := math.Floor(minTile.Y); y <= math.Ceil(maxTile.Y); y += 1 {
			worldOffset := pixel.V(x, y).Scaled(tileWorldWidth)
			s.Draw(
				b,
				pixel.IM.Scaled(
					pixel.ZV, 1.0/16,
				).Moved(
					worldOffset.Scaled(1.0/scale),
				).Chained(
					worldToScreen.Scaled(pixel.ZV, scale),
				),
			)
		}
	}
	b.Draw(win)
}

// Draw implements ent.Entity.
func (b *Background) Draw(win *pixelgl.Window, _ *ent.World, worldToScreen pixel.Matrix) {
	drawLevel(b.l1, b.b1, 0.9, win, worldToScreen)
	drawLevel(b.l2, b.b2, 0.75, win, worldToScreen)
	drawLevel(b.l3, b.b3, 0.5, win, worldToScreen)
}

// DrawLayer implements ent.Entity.
func (b *Background) DrawLayer() int {
	return 100
}
