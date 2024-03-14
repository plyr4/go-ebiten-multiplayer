package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"

	"github.com/plyr4/go-ebiten-multiplayer/shared/constants"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	err := run()
	if err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	host := os.Getenv("WS_SERVER_HOST")
	if len(host) == 0 {
		return errors.New("no host provided in environment variable WS_SERVER_HOST")
	}

	logrus.Tracef("running tcp server using host: %v", host)

	l, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	logrus.Tracef("listening on http://%v", l.Addr())

	// create a single handle server
	s := &http.Server{
		Handler: echoServer{
			logf: logrus.Debugf,
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

	return s.Shutdown(ctx)
}

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type echoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{
			constants.CLIENT_SUBPROTOCOL,
		},
	})
	if err != nil {
		s.logf("%v", err)
		return
	}

	defer c.CloseNow()

	if c.Subprotocol() != constants.CLIENT_SUBPROTOCOL {
		c.Close(websocket.StatusPolicyViolation,
			fmt.Sprintf("expected subprotocol %q but got %q", constants.CLIENT_SUBPROTOCOL, c.Subprotocol()),
		)
		return
	}

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		logrus.Trace("handling client message...")

		err = echo(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			logrus.Trace("recv close message from client")
			return
		}

		if err != nil {
			s.logf("failed to handle client message from %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return errors.Wrap(err, "failed to io.Copy")
	}

	err = w.Close()
	return err
}
