package main

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var GlobalSpriteManager *SpriteManager

func main() {
	pixelgl.Run(run)
}

func run() {
	GlobalSpriteManager = NewSpriteManager("sprites")

	cfg := pixelgl.WindowConfig{
		Title:  "Temporal Engineer",
		VSync:  true,
		Bounds: pixel.R(0, 0, 800, 800),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var screen Screen
	screen = NewMenu()

	for !win.Closed() {
		newScreen := screen.Update(win, 1.0/60.0)
		screen.Draw(win)

		if newScreen != nil {
			screen = newScreen
		}
		win.Update()
	}
}
