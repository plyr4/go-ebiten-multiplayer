package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	keys  []ebiten.Key
	Left  bool
	Right bool
	Up    bool
	Down  bool
}

func (i *Input) Reset() {
	i.Left = false
	i.Right = false
	i.Up = false
	i.Down = false
}

func (i *Input) Update() {
	i.Reset()

	i.keys = inpututil.AppendPressedKeys(i.keys[:0])

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		i.Left = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		i.Right = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		i.Down = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		i.Up = true
	}
}
