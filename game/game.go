package game

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
	"github.com/sirupsen/logrus"
)

type Game struct {
	logger *logrus.Entry

	wsClient *ws.Client

	Foo   string
	Frame int
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

func (g *Game) Shutdown(msg string) error {
	g.wsClient.Close(msg)
	return nil
}

func (g *Game) Update() error {
	// g.Frame++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "hello game: "+g.Foo+"\ngframe: "+strconv.Itoa(g.Frame))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) Run() error {
	logrus.Info("running game")

	// prep
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("<title>")

	// networking
	g.wsClient = ws.New()
	err := g.wsClient.Connect()
	if err != nil {
		return err
	}

	go func() {
		sendErrs := 0
		recvErrs := 0
		for {
			err := g.wsClient.Send(ws.Ping{})
			if err != nil {
				sendErrs++
				logrus.Errorf("error sending ping: %v", err)
			} else {
				sendErrs = 0
			}

			time.Sleep(2 * time.Second)

			_, err = g.wsClient.Receive(ws.Ping{})
			if err != nil {
				recvErrs++
				logrus.Errorf("error receiving ping: %v", err)
			} else {
				recvErrs = 0
			}

			time.Sleep(2 * time.Second)

			if sendErrs > 3 || recvErrs > 3 {
				logrus.Error("too many websocket connection failures, shutting down")

				g.Shutdown("websocket connection failed")

				break
			}
			if sendErrs == 0 && recvErrs == 0 {
				g.Frame++
			}
		}
	}()

	// run
	err = ebiten.RunGame(g)
	if err != nil {
		return err
	}

	return nil
}
