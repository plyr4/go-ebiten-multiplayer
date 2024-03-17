package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"github.com/plyr4/go-ebiten-multiplayer/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	host := os.Getenv("SERVER_WS_HOST")

	if len(host) == 0 {
		logrus.Warnf("missing $SERVER_WS_HOST, using default: %s", constants.SERVER_WS_DEFAULT_HOST)
		host = constants.SERVER_WS_DEFAULT_HOST
	}

	logrus.Infof("running tcp server using: %v", host)

	l, err := net.Listen("tcp", host)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Tracef("listening on: http://%v", l.Addr())

	// create a single handle client server
	s := &http.Server{
		Handler: server.ClientServer{
			Logger: logrus.NewEntry(logrus.StandardLogger()).
				WithFields(logrus.Fields{
					"module": "ws-server",
					"host":   host,
				}),
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	// serve the http handler over tcp
	// send errors to errc channel
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	// send os signals to sigs channel
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	// goroutine handlers
	select {
	case err := <-errc:
		logrus.Errorf("failure serving: %v", err)
	case sig := <-sigs:
		logrus.Infof("terminating server, SIG: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
}
