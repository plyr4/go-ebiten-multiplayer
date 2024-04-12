package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"github.com/plyr4/go-ebiten-multiplayer/internal"
	"github.com/plyr4/go-ebiten-multiplayer/poly"
)

type Boat struct {
	Entity

	Movement
	*Player
}

func NewBoat(g *internal.Game) (*Boat, error) {
	e := new(Boat)

	e.Game = g

	// starting position
	e.X = constants.SCREEN_WIDTH / 2
	e.Y = constants.SCREEN_HEIGHT / 2

	// movement
	e.speed = 2

	return e, nil
}

func (e *Boat) Draw(renderTarget *ebiten.Image) error {
	var scale float32 = 0.2
	var speed float64 = 0.01
	r := float32(math.Cos(float64(e.Frame)*speed)) * scale

	opts := &ebiten.DrawTrianglesOptions{}

	x := float32(constants.SCREEN_WIDTH / 2)
	y := float32(constants.SCREEN_HEIGHT / 2)

	var path vector.Path
	width := float32(175)
	length := float32(250)
	sideLength := float32(0.75)
	sideWidth := float32(0.75)
	sternWidth := float32(0.85)

	length_ := length * 0.7
	width_ := width * 0.5
	// bow
	path.MoveTo(x+width_, y-length_)
	path.LineTo(x, y-length)
	path.LineTo(x-width_, y-length_)

	// left side
	path.LineTo(x-width*sideWidth, y)
	path.LineTo(x-width*sideWidth, y+length*sideLength)

	// stern
	path.LineTo(x-width_*sternWidth, y+length)
	path.LineTo(x+width_*sternWidth, y+length)

	// right side
	path.LineTo(x+width*sideWidth, y+length*sideLength)
	path.LineTo(x+width*sideWidth, y)

	path.Close()
	var vs []ebiten.Vertex
	var is []uint16

	op := &vector.StrokeOptions{}
	op.Width = 2
	// op.LineJoin = vector.LineJoinRound

	vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, op)
	for i := range vs {
		vs[i].DstX = (vs[i].DstX + float32(x)) / constants.WORLD_SCALE
		vs[i].DstY = (vs[i].DstY + float32(y)) / constants.WORLD_SCALE
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(0xff)
		vs[i].ColorG = float32(0xff)
		vs[i].ColorB = float32(0xff)
		vs[i].ColorA = 1

		// _ = r
		// r = 90
		x_, y_ := poly.RotateAbout(vs[i].DstX, vs[i].DstY, x, y, r)
		vs[i].DstX = x_
		vs[i].DstY = y_
	}

	opts.AntiAlias = true

	// if !line {
	// opts.FillRule = ebiten.EvenOdd
	// }

	renderTarget.DrawTriangles(vs, is, poly.IMG_PIXEL_WHITE, opts)

	// opts.Address = ebiten.AddressUnsafe
	// opts.Sides = 3
	// opts.Radius = 100
	// opts.RGBA = color.RGBA{0xff, 0xff, 0xff, 0xff}
	// poly.Polygon(renderTarget, opts)

	// opts.Sides = 3
	// opts.Radius = 98
	// opts.RGBA = color.RGBA{0, 0, 0, 0xff}
	// poly.Polygon(renderTarget, opts)

	return nil
}

func (e *Boat) Update() error {
	return nil
}
