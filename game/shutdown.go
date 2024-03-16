package game

import (
	"github.com/pkg/errors"
)

// Shutdown tears down the game and sets an error to return in Update
func (g *Game) Shutdown(msg string) {
	// close the websocket connection
	if g.wsClient != nil {
		g.wsClient.Close(msg)
	}

	// set internal error
	g.error = errors.New(msg)

	// stop the game
	g.Running = false

	g.logger.Info("goodbye!")
}
