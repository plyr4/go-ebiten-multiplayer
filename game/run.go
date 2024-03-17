package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"github.com/plyr4/go-ebiten-multiplayer/entity/player"
)

func (g *Game) Run() error {
	g.Running = true

	g.logger.Info("preparing game")

	// window
	ebiten.SetWindowSize(constants.WINDOW_WIDTH, constants.WINDOW_HEIGHT)
	ebiten.SetWindowTitle(constants.WINDOW_TITLE)

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

	// player
	p, err := player.New()
	if err != nil {
		return errors.Wrap(err, "unable to create player")
	}

	g.player = p

	// entities
	g.entities = append(g.entities, p)

	g.logger.Debug("starting game")

	err = ebiten.RunGame(g)
	if err != nil {
		return errors.Wrap(err, "game error")
	}

	return nil
}
