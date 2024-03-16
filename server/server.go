package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/shared/constants"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// server status
var connectedPlayers = 0

// todo: move this implementation into a server package
type ClientServer struct{}

func (s ClientServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept the client connection
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{
			constants.CLIENT_SUBPROTOCOL,
		},
	})
	if err != nil {
		logrus.Errorf("%v", err)

		return
	}

	defer c.CloseNow()

	// check for client protocol
	if c.Subprotocol() != constants.CLIENT_SUBPROTOCOL {
		c.Close(websocket.StatusPolicyViolation,
			fmt.Sprintf("expected subprotocol %q but got %q", constants.CLIENT_SUBPROTOCOL, c.Subprotocol()),
		)

		return
	}

	connectedPlayers++
	defer func() { connectedPlayers-- }()

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		// receive messages from the client
		err = handleClientMessage(r.Context(), c, l)

		// handle closures
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			logrus.Tracef("received close message from client: %v", err)

			return
		}

		if websocket.CloseStatus(err) == websocket.StatusGoingAway {
			logrus.Tracef("received going away message from client: %v", err)

			return
		}

		if err != nil {
			logrus.Errorf("failed to handle client message from %v: %v", r.RemoteAddr, err)

			return
		}
	}
}

// handleClientMessage reads from the WebSocket connection then handles the incoming message
func handleClientMessage(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	// 10s to complete
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	logrus.Tracef("waiting %v before reading", l.Burst())

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	logrus.Trace("reading msg")

	msg := new(ws.Msg)

	err = wsjson.Read(ctx, c, msg)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	logrus.Tracef("got msg: %v", msg)

	// handle the message based on type
	// todo: implement this
	if msg.Ping != nil {
		logrus.Tracef("received ping: %v", msg)
	} else if msg.ClientUpdate != nil {
		logrus.Tracef("received client update: %v", msg)
	} else if msg.ServerUpdate != nil {
		logrus.Tracef("received server update: %v", msg)
	} else {
		logrus.Tracef("received unknown message type: %v", msg)
	}

	// respond to the client with a server update
	// todo: move this into the handler above and respond when asked
	msg = new(ws.Msg)

	su := ws.ServerUpdate{
		Status:           "ok",
		ConnectedPlayers: connectedPlayers,
	}

	msg.ServerUpdate = &su

	logrus.Tracef("sending server update: %v", su)

	err = wsjson.Write(ctx, c, msg)
	if err != nil {
		return errors.Wrap(err, "failed to write server update")
	}

	return err
}
