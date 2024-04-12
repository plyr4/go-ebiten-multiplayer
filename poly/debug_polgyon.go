package poly

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PolygonDrawOpts struct {
	color.RGBA
	DstX, DstY float32
	Sides      int
	Radius     float64
	Rotation   float32
	*ebiten.DrawTrianglesOptions
}

func Polygon(renderTarget *ebiten.Image, opts *PolygonDrawOpts) {
	vertices := []ebiten.Vertex{}
	for i := 0; i < opts.Sides; i++ {
		rate := float64(i) / float64(opts.Sides)

		x := float32(opts.Radius*math.Cos(2*math.Pi*rate)) + opts.DstX
		y := float32(opts.Radius*math.Sin(2*math.Pi*rate)) + opts.DstY

		// rotate the point
		x, y = RotateAbout(x, y, opts.DstX, opts.DstY, opts.Rotation)

		// add a vertex for each point of the polygon
		vertices = append(vertices, ebiten.Vertex{
			DstX:   x,
			DstY:   y,
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(opts.R),
			ColorG: float32(opts.G),
			ColorB: float32(opts.B),
			ColorA: float32(opts.A),
		})
	}

	// add a vertex for the center of the polygon
	vertices = append(vertices, ebiten.Vertex{
		DstX:   opts.DstX,
		DstY:   opts.DstY,
		SrcX:   0,
		SrcY:   0,
		ColorR: float32(opts.R),
		ColorG: float32(opts.G),
		ColorB: float32(opts.B),
		ColorA: float32(opts.A),
	})

	indices := []uint16{}
	for i := 0; i < opts.Sides; i++ {
		// connect each vertex to the next and also the center
		// because they have to be triangles
		indices = append(indices, uint16(i), uint16(i+1)%uint16(opts.Sides), uint16(opts.Sides))
	}

	// draw the polygon
	renderTarget.DrawTriangles(vertices, indices,
		// source texture
		IMG_TILE_WHITE.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image),
		opts.DrawTrianglesOptions)
}

func RotateAbout(x, y, centerX, centerY, theta float32) (float32, float32) {
	translatedX := x - centerX
	translatedY := y - centerY

	cosTheta := float32(math.Cos(float64(theta)))
	sinTheta := float32(math.Sin(float64(theta)))

	rotatedX := translatedX*cosTheta - translatedY*sinTheta
	rotatedY := translatedX*sinTheta + translatedY*cosTheta

	newX := rotatedX + centerX
	newY := rotatedY + centerY

	return newX, newY
}
