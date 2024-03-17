package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plyr4/go-ebiten-multiplayer/input"
)

type IEntity interface {
	Draw(*ebiten.Image) error
	Update(*input.Input) error

	Position() (float64, float64)
}
