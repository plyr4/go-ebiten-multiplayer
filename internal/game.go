package internal

import (
	"github.com/plyr4/go-ebiten-multiplayer/input"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
)

// Game represents game state shareable between packages
type Game struct {
	UUID    string
	Running bool
	Frame   int
	*input.Input

	// todo: debug cleanup
	Foo              string
	Roundtrips       int
	ConnectedPlayers map[string]*ws.PlayerData
}
