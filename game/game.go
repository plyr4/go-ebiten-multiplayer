package game

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
	"github.com/sirupsen/logrus"
)

type Game struct {
	wsClient *ws.Client

	Running bool

	error
	logger *logrus.Entry

	Debug
}

type Debug struct {
	Foo              string
	Roundtrips       int
	Frame            int
	ConnectedPlayers int
}

// WithLogger attaches a logger to the Game
func (g *Game) WithLogger(l *logrus.Entry) {
	g.logger = l
}

func New() *Game {
	g := new(Game)

	// todo: generate uuid from machine identifier (MAC address, etc.)
	// something that will allow the same player to reconnect
	// or... should multiple tabs of the same game be supported?
	uuid := "1234"

	// logging
	logger := logrus.NewEntry(logrus.StandardLogger()).WithFields(
		logrus.Fields{
			"ID": uuid,
		},
	)

	g.WithLogger(logger)
	return g
}

// Shutdown tears down the game and sets an error to return in Update
func (g *Game) Shutdown(msg string) {
	g.wsClient.Close(msg)
	g.error = errors.New(msg)
	g.Running = false
}

// Multiplayer maintains a connection to the server
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
			logrus.Errorf("error sending client update: %v", err)
		} else {
			sendErrs = 0
		}

		time.Sleep(2 * time.Second)

		msg = new(ws.Msg)
		_, err = g.wsClient.Receive(msg)
		if err != nil {
			recvErrs++
			logrus.Errorf("error receiving server update: %v", err)
		} else {
			recvErrs = 0
		}

		if msg.Ping != nil {
			logrus.Tracef("received ping: %v", msg)
		} else if msg.ClientUpdate != nil {
			logrus.Tracef("received client update: %v", msg)
		} else if msg.ServerUpdate != nil {
			logrus.Tracef("received server update: %v", msg)
			g.Debug.ConnectedPlayers = msg.ServerUpdate.ConnectedPlayers
		} else {
			logrus.Tracef("received unknown message type: %v", msg)
		}

		time.Sleep(2 * time.Second)

		if sendErrs > 3 || recvErrs > 3 {
			logrus.Error("too many websocket connection failures, attempting to reconnect")

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

func (g *Game) Update() error {
	if g.error != nil {
		return g.error
	}

	g.Debug.Frame++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen,
		"debug: "+g.Debug.Foo+
			"\nsuccessful server roundtrips: "+strconv.Itoa(g.Debug.Roundtrips)+
			"\ng.Debug.Frame: "+strconv.Itoa(g.Debug.Frame)+
			"\nconnected players: "+strconv.Itoa(g.Debug.ConnectedPlayers)+
			"\n\ntime: "+time.Now().Format(time.RFC3339))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) Run(ctx context.Context) error {
	logrus.Info("running game")

	g.Running = true

	// prep
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("<title>")

	// multiplayer
	go func() {
		// this should run forever
		g.error = g.RunMultiplayer()
		if g.error == nil {
			g.logger.Error("multiplayer ended without error, this should not happen")
		}
	}()

	// run
	err := ebiten.RunGame(g)
	if err != nil {
		return err
	}

	return nil
}
