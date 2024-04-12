package poly

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	IMG_TILE_WHITE  = ebiten.NewImage(3, 3)
	IMG_PIXEL_WHITE *ebiten.Image
)

func init() {
	IMG_TILE_WHITE.Fill(color.White)
	IMG_PIXEL_WHITE = IMG_TILE_WHITE.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}
