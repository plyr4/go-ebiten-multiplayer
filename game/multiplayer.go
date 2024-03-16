package game

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
)

// RunMultiplayer maintains a connection to the server
// connect to the server
// maintain a connection
//
//   - -> ping (send client state)
//   - pong <- (receive server state)
//   - update the game state based on server response
func (g *Game) RunMultiplayer() error {
	g.logger.Infof("establishing multiplayer session")

	// networking
	g.wsClient = ws.New()
	err := g.wsClient.Connect()
	if err != nil {
		return err
	}

	sendErrs := 0
	recvErrs := 0

	// loop on the websocket connection forever
	// if multiplayer ends, the game will shutdown
	for {
		msg := new(ws.Msg)
		msg.ClientUpdate = &ws.ClientUpdate{
			Status: "client-ping",
			Foo:    1,
		}

		err := g.wsClient.Send(msg)
		if err != nil {
			sendErrs++
			g.logger.Errorf("error sending client update: %v", err)
		} else {
			sendErrs = 0
		}

		time.Sleep(2 * time.Second)

		msg = new(ws.Msg)
		_, err = g.wsClient.Receive(msg)
		if err != nil {
			recvErrs++
			g.logger.Errorf("error receiving server update: %v", err)
		} else {
			recvErrs = 0
		}

		if msg.Ping != nil {
			g.logger.Tracef("received ping: %v", msg)
		} else if msg.ClientUpdate != nil {
			g.logger.Tracef("received client update: %v", msg)
		} else if msg.ServerUpdate != nil {
			g.logger.Tracef("received server update: %v", msg)
			g.Debug.ConnectedPlayers = msg.ServerUpdate.ConnectedPlayers
		} else {
			g.logger.Tracef("received unknown message type: %v", msg)
		}

		time.Sleep(2 * time.Second)

		if sendErrs > 3 || recvErrs > 3 {
			g.logger.Error("too many websocket connection failures, attempting to reconnect")

			rErr := g.wsClient.Reconnect()
			if rErr != nil {
				rErr = errors.Wrap(rErr, "unable to reconnect")

				g.logger.Error(rErr)

				return fmt.Errorf("websocket connection failed: %v", rErr)
			}
			sendErrs = 0
			recvErrs = 0
		}

		if sendErrs == 0 && recvErrs == 0 {
			g.Debug.Roundtrips++
		}
	}
}
