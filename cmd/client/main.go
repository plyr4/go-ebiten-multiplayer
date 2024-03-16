package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/plyr4/go-ebiten-multiplayer/game"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	ctx := context.Background()

	// send os signals to sigs channel
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// create the game
	g, err := game.New(
		game.WithContext(ctx),
		// todo: create a uuid
		game.WithUUID("1234"),
		// todo: cli or flags
		game.WithMultiplayer(os.Getenv("CLIENT_MULTIPLAYER") == "true"),
	)
	if err != nil {
		logrus.Error(err)

		return
	}

	// late shutdown if necessary
	defer func() {
		if g.Running {
			g.Shutdown("shutting down game")

			return
		}
	}()

	// run the game
	err = g.Run()
	if err != nil {
		logrus.Error(err)

		return
	}
}
