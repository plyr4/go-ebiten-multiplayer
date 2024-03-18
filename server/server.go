package server

import (
	"context"
	"fmt"
	"net/http"
	"sort"
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
// todo: maintain an efficient sorted list of this data for responding faster to the client
var players = map[string]*ws.PlayerData{}
var mu sync.RWMutex

// todo: refactor most of this into ws package
// this package should be the actual server logic and how we manage players
type ClientServer struct {
	// exported because we set it somewhere else
	// todo: expose a WithLogger method someday
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
	// todo: refactor most of this into ws package under like ClientConnection or something
	// connected++
	// todo: clean up connected players
	// defer func() { connected-- }()

	// how do i identify the uuid of the client before sending or receiving messages???

	// todo: verify some kind of token that allows players to connect to the server
	// maybe a mac address can have 3 open windows or something

	// todo: make this configurable
	latency := constants.SERVER_WS_LATENCY
	rateLimiter := rate.NewLimiter(rate.Every(latency), 10)

	// identify client uuid earlier, using request headers
	clientUUID := ""

	for {
		// receive messages from the client
		clientUUID, err = s.handleClientMessage(r.Context(), conn, rateLimiter, clientUUID)
		if err != nil {
			s.Logger.Errorf("failed to handle client message: %v", err)

			// todo: should be able to handle disconnection better
			// the server shouldn't continuously grow larger when players are not connected
			// after X amount of time we should remove the player from the map
			// disconnected would show them as grayed out, then eventually ejected

			if len(clientUUID) > 0 {
				s.Logger.Infof("client disconnected: %v", clientUUID)
				mu.Lock()
				players[clientUUID].Connected = false
				mu.Unlock()
			}

			return
		}
	}
}

// handleClientMessage reads from the WebSocket connection then handles the incoming message
func (s ClientServer) handleClientMessage(ctx context.Context, conn *websocket.Conn, rateLimiter *rate.Limiter, clientUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// apply server latency per client
	err := rateLimiter.Wait(ctx)
	if err != nil {
		return clientUUID, errors.Wrap(err, "failed to wait for rate limiter")
	}

	msg := new(ws.Msg)

	err = wsjson.Read(ctx, conn, msg)
	if err != nil {
		return clientUUID, errors.Wrap(err, "failed to read")
	}

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
		players[msg.ClientUpdate.Player.UUID].Connected = true
		clientUUID = msg.ClientUpdate.Player.UUID

		mu.Unlock()

		// respond to the client with a server update
		msg = new(ws.Msg)

		su := ws.ServerUpdate{
			Status:  "ok",
			Players: []ws.PlayerData{},
		}

		mu.RLock()

		// sort the players by their UUID so that the order is consistent
		keys := make([]string, 0, len(players))
		for k := range players {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			p, ok := players[k]
			if !ok || p == nil {
				continue
			}

			su.Players = append(su.Players, *p)
		}

		mu.RUnlock()

		msg.ServerUpdate = &su

		s.Logger.Tracef("responding with server update: %v", su)

		err = wsjson.Write(ctx, conn, msg)
		if err != nil {
			return clientUUID, errors.Wrap(err, "failed to write server update")
		}
	} else if msg.ServerUpdate != nil {
		s.Logger.Tracef("received server update: %v", msg)
	} else {
		s.Logger.Tracef("received unknown message type: %v", msg)
	}

	// error and closure handling
	// not sure why this is necessary
	// todo: handle disconnects closures and errors better
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		return clientUUID, errors.Wrap(err, "received normal closure")
	}

	if websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return clientUUID, errors.Wrap(err, "received going away")
	}

	if websocket.CloseStatus(err) == websocket.StatusAbnormalClosure {
		return clientUUID, errors.Wrap(err, "received abnormal closure")
	}

	if websocket.CloseStatus(err) == websocket.StatusUnsupportedData {
		return clientUUID, errors.Wrap(err, "received unsupported data")
	}

	if websocket.CloseStatus(err) == websocket.StatusPolicyViolation {
		return clientUUID, errors.Wrap(err, "received policy violation")
	}

	if websocket.CloseStatus(err) == websocket.StatusMessageTooBig {
		return clientUUID, errors.Wrap(err, "received message too big")
	}

	return clientUUID, err
}
