package main

import "github.com/gopxl/pixel/pixelgl"

type Screen interface {
	Update(win *pixelgl.Window, dt float64) Screen
	Draw(win *pixelgl.Window)
}
