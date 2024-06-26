package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"github.com/plyr4/go-ebiten-multiplayer/entities"
)

func (g *Game) Run() error {
	g.Running = true

	g.logger.Info("preparing game")

	// window
	ebiten.SetWindowSize(constants.WINDOW_WIDTH, constants.WINDOW_HEIGHT)
	ebiten.SetWindowTitle(constants.WINDOW_TITLE)

	// boat
	b, err := entities.NewBoat(g.Game)
	if err != nil {
		return errors.Wrap(err, "unable to create boat")
	}
	g.entities = append(g.entities, b)

	// player
	p, err := entities.NewPlayer(g.Game)
	if err != nil {
		return errors.Wrap(err, "unable to create player")
	}

	g.Player = p
	b.Player = p

	g.entities = append(g.entities, p)

	// multiplayer
	if g.multiplayer {
		g.logger.Debug("running in multiplayer mode")

		go func() {
			// this should run forever
			err := g.RunMultiplayer()
			if err != nil {
				g.logger.Error("multiplayer ended without error, this should not happen")
			} else {
				g.logger.Error("multiplayer ended without error, this should not happen")
			}
		}()
	} else {
		g.logger.Debug("running in local mode")
	}

	g.logger.Debug("starting game")

	err = ebiten.RunGame(g)
	if err != nil {
		return errors.Wrap(err, "game error")
	}

	return nil
}
