package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
)

func (g *Game) Run() error {
	g.Running = true

	g.logger.Info("preparing game")

	// window
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("title")

	// multiplayer
	if g.multiplayer {
		g.logger.Debug("running in multiplayer mode")

		go func() {
			// this should run forever
			g.error = g.RunMultiplayer()
			if g.error == nil {
				g.logger.Error("multiplayer ended without error, this should not happen")
			}
		}()
	} else {
		g.logger.Debug("running in local mode")
	}

	// run
	g.logger.Debug("starting game")

	err := ebiten.RunGame(g)
	if err != nil {
		return errors.Wrap(err, "game error")
	}

	return nil
}
