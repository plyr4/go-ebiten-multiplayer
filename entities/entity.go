package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plyr4/go-ebiten-multiplayer/internal"
)

type IEntity interface {
	Draw(*ebiten.Image) error
	Update() error
}

type Entity struct {
	X, Y float64
	*internal.Game
	// todo: children entities
}
