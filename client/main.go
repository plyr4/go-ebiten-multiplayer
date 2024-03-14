package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/plyr4/go-ebiten-multiplayer/game"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	// send os signals to sigs channel
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// send errors to errc channel
	errc := make(chan error, 1)
	go func() {
		g := game.New()
		errc <- g.Run()
		g.Shutdown("server shutdown")
	}()

	// goroutine handlers
	select {
	case err := <-errc:
		logrus.Errorf("failure serving: %v", err)

	case sig := <-sigs:
		logrus.Infof("terminating server, SIG: %v", sig)
	}
}
