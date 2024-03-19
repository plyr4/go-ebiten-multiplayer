package sprite

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type FrameOpts struct {
	CurrentGameFrame int
	FrameOX          int
	FrameOY          int
	FrameWidth       int
	FrameHeight      int
	FrameCount       int
}

func Sprite(sheet *ebiten.Image, opts *FrameOpts) *ebiten.Image {
	// center the sprite

	frame := (opts.CurrentGameFrame / 5) % opts.FrameCount

	// sprite sheet offset
	sx, sy := opts.FrameOX+frame*opts.FrameWidth, opts.FrameOY
	return sheet.SubImage(image.Rect(sx, sy, sx+opts.FrameWidth, sy+opts.FrameHeight)).(*ebiten.Image)
}
