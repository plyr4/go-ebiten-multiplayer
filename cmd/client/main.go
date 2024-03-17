package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/plyr4/go-ebiten-multiplayer/game"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	ctx := context.Background()

	// generate a unique id for this session
	id := newUUID()

	// send os signals to sigs channel
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// create the game
	g, err := game.New(
		game.WithContext(ctx),
		game.WithUUID(id),
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

// see: https://stackoverflow.com/questions/44859156/get-permanent-mac-address
func newUUID() string {
	ifas, err := net.Interfaces()
	if err != nil {
		return uuid.New().String()
	}

	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}

	if len(as) == 0 {
		return uuid.New().String()
	}

	return as[0]
}
