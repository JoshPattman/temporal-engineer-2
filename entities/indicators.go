package entities

import (
	"ent"
	"fmt"

	_ "embed"

	"github.com/golang/freetype/truetype"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
)

//go:embed upheavtt.ttf
var mainFont []byte

func NewSheildsIndicator() *statsIndicator {
	return NewStatsIndicator("bubble.png", 150, func(w *ent.World) int {
		player, ok := ent.First(
			ent.OfType[*Player](
				w.WithTag("player"),
			),
		)
		if !ok {
			return 0
		}
		return player.Shields()
	})
}

func NewMineralsIndicator() *statsIndicator {
	return NewStatsIndicator("minerals.png", 200, func(w *ent.World) int {
		player, ok := ent.First(
			ent.OfType[*Player](
				w.WithTag("player"),
			),
		)
		if !ok {
			return 0
		}
		return player.Minerals()
	})
}

func NewStatsIndicator(spriteName string, vPos float64, get func(*ent.World) int) *statsIndicator {
	sprite := GlobalSpriteManager.FullSprite(spriteName)
	sa := &statsIndicator{
		sprite: sprite,
		text:   text.New(pixel.ZV, sheildsAtlas()).AlignedTo(pixel.Right),
		vPos:   vPos,
		get:    get,
	}
	sa.text.Color = colornames.Skyblue
	return sa
}

type statsIndicator struct {
	ent.CoreEntity
	ent.WithDraw
	ent.WithUpdate
	sprite *pixel.Sprite
	value  int
	text   *text.Text
	get    func(*ent.World) int
	vPos   float64
}

func (c *statsIndicator) Update(win *pixelgl.Window, all *ent.World, dt float64) {
	c.value = c.get(all)
}

func (c *statsIndicator) Draw(win *pixelgl.Window, _ *ent.World, worldToScreen pixel.Matrix) {
	c.sprite.Draw(
		win,
		pixel.IM.Scaled(
			pixel.ZV, 30.0/c.sprite.Frame().W(),
		).Moved(
			pixel.V(30, c.vPos),
		),
	)
	c.text.Clear()
	fmt.Fprintf(c.text, "%d", c.value)
	c.text.Draw(
		win,
		pixel.IM.Moved(
			pixel.V(60, c.vPos+1),
		),
	)
}

func (c *statsIndicator) DrawLayer() int { return -10 }

func sheildsAtlas() *text.Atlas {
	ttf, err := truetype.Parse(mainFont)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 40,
	})
	return text.NewAtlas(face, text.ASCII)
}
