package main

import (
	_ "embed"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
)

//go:embed entities/upheavtt.ttf
var mainFont []byte

func titleAtlas() *text.Atlas {
	ttf, err := truetype.Parse(mainFont)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 65,
	})
	return text.NewAtlas(face, text.ASCII)
}

func infoAtlas() *text.Atlas {
	ttf, err := truetype.Parse(mainFont)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 35,
	})
	return text.NewAtlas(face, text.ASCII)
}

func NewMenu() Screen {
	titleText := text.New(pixel.ZV, titleAtlas()).AlignedTo(pixel.Center)
	titleText.Color = colornames.White
	titleText.Clear()
	titleText.Write([]byte("Temporal Engineer"))

	infoText := text.New(pixel.ZV, infoAtlas()).AlignedTo(pixel.Center)
	infoText.Color = colornames.Grey
	infoText.Clear()
	infoText.Write([]byte("Space to Play\n\nEscape to Quit"))

	return &Menu{
		titleText: titleText,
		infoText:  infoText,
	}
}

type Menu struct {
	titleText     *text.Text
	infoText      *text.Text
	titleRotation float64
	timer         float64
}

// Draw implements Screen.
func (m *Menu) Draw(win *pixelgl.Window) {
	win.Clear(pixel.RGB(0.016, 0.071, 0.137))
	m.titleText.Draw(
		win,
		pixel.IM.Rotated(
			pixel.ZV,
			m.titleRotation,
		).Moved(
			win.Bounds().Center(),
		),
	)
	m.infoText.Draw(
		win,
		pixel.IM.Moved(
			win.Bounds().Center(),
		).Moved(
			pixel.V(0, -m.titleText.LineHeight*2),
		),
	)
}

// Update implements Screen.
func (m *Menu) Update(win *pixelgl.Window, dt float64) Screen {
	if win.JustPressed(pixelgl.KeySpace) {
		return NewGame()
	}
	m.timer += dt
	m.titleRotation = math.Sin(m.timer*1.5) * 0.05
	return nil
}
