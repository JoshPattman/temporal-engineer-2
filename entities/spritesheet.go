package entities

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path"

	"github.com/gopxl/pixel"
)

var GlobalSpriteManager *SpriteManager

func NewSpriteManager(root string) *SpriteManager {
	return &SpriteManager{
		loaded: make(map[string]pixel.Picture),
		root:   root,
	}
}

type SpriteManager struct {
	loaded map[string]pixel.Picture
	root   string
}

func (s *SpriteManager) Picture(spritePath string) pixel.Picture {
	pic, ok := s.loaded[spritePath]
	if ok {
		return pic
	}
	f, err := os.Open(path.Join(s.root, spritePath))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	newPic := pixel.PictureDataFromImage(img)
	s.loaded[spritePath] = newPic
	return newPic
}

func (s *SpriteManager) FullSprite(spritePath string) *pixel.Sprite {
	pic := s.Picture(spritePath)
	return pixel.NewSprite(pic, pic.Bounds())
}

type TilePos struct {
	X int
	Y int
}

func (s *SpriteManager) TiledSprites(spritePath string, tileSize int, positions []TilePos) []*pixel.Sprite {
	pic := s.Picture(spritePath)
	sprites := make([]*pixel.Sprite, 0, len(positions))
	for _, p := range positions {
		intR := pixel.R(float64(p.X), float64(p.Y), float64(p.X+1), float64(p.Y+1))
		sprite := pixel.NewSprite(
			pic,
			pixel.Rect{
				Min: intR.Min.Scaled(float64(tileSize)),
				Max: intR.Max.Scaled(float64(tileSize)),
			},
		)
		sprites = append(sprites, sprite)
	}
	return sprites
}

func MustLoadPic(pngData []byte) pixel.Picture {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img)
}

func GetSprite(pic pixel.Picture, tileWidth float64, x int, y int) *pixel.Sprite {
	intR := pixel.R(float64(x), float64(y), float64(x+1), float64(y+1))
	return pixel.NewSprite(pic, pixel.Rect{
		Min: intR.Min.Scaled(tileWidth),
		Max: intR.Max.Scaled(tileWidth),
	})
}
