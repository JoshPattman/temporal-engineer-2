package main

import (
	"ent"
	"fmt"

	"github.com/golang/freetype/truetype"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
)

var _ ent.Entity = &Compass{}

func NewSheildsIndicator() *ShieldsIndicator {
	sprite := GlobalSpriteManager.FullSprite("bubble.png")
	return &ShieldsIndicator{
		sprite: sprite,
		text:   text.New(pixel.ZV, sheildsAtlas()).AlignedTo(pixel.Center),
	}
}

type ShieldsIndicator struct {
	ent.EntityBase
	sprite  *pixel.Sprite
	sheilds int
	text    *text.Text
}

// Draw implements ent.Entity.
func (c *ShieldsIndicator) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	c.sprite.Draw(
		win,
		pixel.IM.Scaled(
			pixel.ZV, 1.2,
		).Moved(
			pixel.V(160, 64),
		),
	)
	c.text.Clear()
	c.text.WriteString(fmt.Sprintf("%d", c.sheilds))
	if c.sheilds == 1 {
		c.text.Color = colornames.Yellow
	} else if c.sheilds == 0 {
		c.text.Color = colornames.Red
	} else {
		c.text.Color = colornames.Skyblue
	}
	c.text.Draw(
		win,
		pixel.IM.Scaled(
			pixel.ZV, 1.2,
		).Moved(
			pixel.V(160, 66),
		),
	)
}

// DrawLayer implements ent.Entity.
func (c *ShieldsIndicator) DrawLayer() int {
	return -10
}

// Update implements ent.Entity.
func (c *ShieldsIndicator) Update(win *pixelgl.Window, all *ent.World, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	player, ok := ent.First(
		ent.OfType[*Player](
			all.ForTag("player"),
		),
	)
	if !ok {
		c.sheilds = 0
		return nil, nil
	}
	c.sheilds = player.Shields()
	return nil, nil
}

// UpdateLayer implements ent.Entity.
func (c *ShieldsIndicator) UpdateLayer() int {
	return -10
}

func sheildsAtlas() *text.Atlas {
	ttf, err := truetype.Parse(mainFont)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 45,
	})
	return text.NewAtlas(face, text.ASCII)
}
