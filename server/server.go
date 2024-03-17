package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// server state
var players = map[string]*ws.PlayerData{}
var mu sync.RWMutex

// todo: move this implementation into a server package
type ClientServer struct {
	Logger *logrus.Entry
}

func (s ClientServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept the client connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{
			constants.CLIENT_SUBPROTOCOL,
		},
	})
	if err != nil {
		s.Logger.Errorf("%v", err)

		return
	}

	defer conn.CloseNow()

	// check for client protocol
	if conn.Subprotocol() != constants.CLIENT_SUBPROTOCOL {
		conn.Close(websocket.StatusPolicyViolation,
			fmt.Sprintf("expected subprotocol %q but got %q", constants.CLIENT_SUBPROTOCOL, conn.Subprotocol()),
		)

		return
	}

	// todo: capture the uuid of the client connection here, not in the client update message
	// todo: attach the player info to the logger fields
	// connected++
	// todo: clean up connected players
	// defer func() { connected-- }()

	// todo: make this configurable
	latency := time.Millisecond * 1

	rateLimiter := rate.NewLimiter(rate.Every(latency), 10)

	for {
		// receive messages from the client
		err = s.handleClientMessage(r.Context(), conn, rateLimiter)
		if err != nil {
			s.Logger.Errorf("failed to handle client message: %v", err)

			return
		}
	}
}

// handleClientMessage reads from the WebSocket connection then handles the incoming message
func (s ClientServer) handleClientMessage(ctx context.Context, conn *websocket.Conn, rateLimiter *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	s.Logger.Tracef("waiting %v before reading", rateLimiter.Burst())

	// apply server latency per client
	err := rateLimiter.Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for rate limiter")
	}

	s.Logger.Trace("reading msg from client")

	msg := new(ws.Msg)

	err = wsjson.Read(ctx, conn, msg)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	s.Logger.Tracef("received msg: %v", msg)

	// handle the message based on type
	// todo: implement this
	if msg.Ping != nil {
		s.Logger.Tracef("received ping: %v", msg)
	} else if msg.ClientUpdate != nil {
		s.Logger.Tracef("received client update: %v", msg)

		mu.Lock()

		// update this player
		_, ok := players[msg.ClientUpdate.Player.UUID]
		if !ok {
			players[msg.ClientUpdate.Player.UUID] = &msg.ClientUpdate.Player
		}

		players[msg.ClientUpdate.Player.UUID].X = msg.ClientUpdate.Player.X
		players[msg.ClientUpdate.Player.UUID].Y = msg.ClientUpdate.Player.Y

		mu.Unlock()

		// respond to the client with a server update
		msg = new(ws.Msg)

		su := ws.ServerUpdate{
			Status:  "ok",
			Players: []ws.PlayerData{},
		}

		mu.RLock()
		for _, p := range players {
			if p == nil {
				continue
			}

			su.Players = append(su.Players, *p)
		}
		mu.RUnlock()

		msg.ServerUpdate = &su

		s.Logger.Tracef("responding with server update: %v", su)

		err = wsjson.Write(ctx, conn, msg)
		if err != nil {
			return errors.Wrap(err, "failed to write server update")
		}
	} else if msg.ServerUpdate != nil {
		s.Logger.Tracef("received server update: %v", msg)
	} else {
		s.Logger.Tracef("received unknown message type: %v", msg)
	}

	// error and closure handling
	// not sure why this is necessary
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		return errors.Wrap(err, "received normal closure")
	}

	if websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return errors.Wrap(err, "received going away")
	}

	if websocket.CloseStatus(err) == websocket.StatusAbnormalClosure {
		return errors.Wrap(err, "received abnormal closure")
	}

	if websocket.CloseStatus(err) == websocket.StatusUnsupportedData {
		return errors.Wrap(err, "received unsupported data")
	}

	if websocket.CloseStatus(err) == websocket.StatusPolicyViolation {
		return errors.Wrap(err, "received policy violation")
	}

	if websocket.CloseStatus(err) == websocket.StatusMessageTooBig {
		return errors.Wrap(err, "received message too big")
	}

	return err
}
