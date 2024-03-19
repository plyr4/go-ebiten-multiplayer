package game

import (
	"time"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/entities"
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

	// create a websocket client
	c, err := ws.New(g.logger, ws.WithContext(g.ctx))
	if err != nil {
		return err
	}

	g.wsClient = c

	// todo: clean this up
	err = g.wsClient.Connect()
	for i := 0; err != nil; i++ {
		if i > 2 {
			i = 2
		}

		err = g.wsClient.Connect()
		if err == nil {
			break
		}

		if i == 2 {
			err = errors.Wrap(err, "unable to initialize connection")

			g.logger.Error(err)
		}

		t := 1 + i*i

		time.Sleep(time.Duration(t) * time.Second)
	}

	sendErrs := 0
	recvErrs := 0

	// maintain a connection
	// ping (client state) ->
	// <- pong (server state)
	for {
		msg := new(ws.Msg)
		msg.ClientUpdate = &ws.ClientUpdate{
			Status: "client-ping",
			Player: g.player.ToMultiplayerData(),
		}

		err := g.wsClient.Send(msg)
		if err != nil {
			sendErrs++
			g.logger.Errorf("error sending client update: %v", err)
		} else {
			sendErrs = 0
		}

		msg = new(ws.Msg)

		_, err = g.wsClient.Receive(msg)
		if err != nil {
			recvErrs++
			g.logger.Errorf("error receiving server update: %v", err)
		} else {
			recvErrs = 0
		}

		// todo: clean this up
		if msg.Ping != nil {
			g.logger.Tracef("received ping: %v", msg)
		} else if msg.ClientUpdate != nil {
			g.logger.Tracef("received client update: %v", msg)
		} else if msg.ServerUpdate != nil {
			g.logger.Tracef("received server update: %v", msg)

			g.ConnectedPlayers = msg.ServerUpdate.Players

			for _, e := range g.entities {
				np, ok := e.(*entities.NetworkPlayer)
				if ok {
					pd, ok := g.ConnectedPlayers[np.PlayerData.UUID]
					if ok {
						np.PlayerData = pd
						pd.ClientUpdated = true
					} else {
						// todo: remove this entity after N seconds of disconnection
						np.PlayerData.Connected = false
					}
				}
			}

			for _, pd := range g.ConnectedPlayers {
				if !pd.ClientUpdated && pd.UUID != g.UUID {
					np, err := entities.NewNetworkPlayer(g.Game, pd)
					if err != nil {
						g.logger.Errorf("error creating network player: %v", err)
					} else {
						g.entities = append(g.entities, np)
					}
				}
			}
		} else {
			g.logger.Tracef("received unknown message type: %v", msg)
		}

		// todo: clean this up
		if sendErrs > 3 || recvErrs > 3 {
			g.logger.Error("too many websocket connection failures, attempting to reconnect")

			for i := 0; ; i++ {
				if i > 2 {
					i = 2
				}

				rErr := g.wsClient.Reconnect()
				if rErr == nil {
					break
				}

				if i == 2 {
					rErr = errors.Wrap(rErr, "unable to reconnect")

					g.logger.Error(rErr)
				}

				t := 1 + i*i

				time.Sleep(time.Duration(t) * time.Second)
			}

			sendErrs = 0
			recvErrs = 0
		}

		if sendErrs == 0 && recvErrs == 0 {
			g.Roundtrips++
		}
	}
}
