package game

import (
	"context"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
	"github.com/sirupsen/logrus"
)

type Game struct {
	uuid        string
	multiplayer bool
	wsClient    *ws.Client
	logger      *logrus.Entry
	ctx         context.Context
	error

	Running bool
	Debug
}

type Debug struct {
	Foo              string
	Roundtrips       int
	Frame            int
	ConnectedPlayers int
}

// New creates a new Game instance
func New(opts ...Opt) (*Game, error) {
	g := new(Game)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(g)
		if err != nil {
			return nil, err
		}
	}

	err := g.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "game is invalid")
	}

	// logging
	g.logger = logrus.NewEntry(logrus.StandardLogger()).WithFields(
		logrus.Fields{
			"ID": g.uuid,
		},
	)

	return g, nil
}

// Validate checks that the game is ready to run, returns error on misconfiguration
func (g *Game) Validate() error {
	if g.ctx == nil {
		return errors.New("missing context")
	}

	if len(g.uuid) == 0 {
		return errors.New("missing uuid")
	}

	return g.error
}

type Opt func(*Game) error

// WithContext sets the internal context
func WithContext(ctx context.Context) Opt {
	return func(g *Game) error {
		g.ctx = ctx

		return nil
	}
}

// WithMultiplayer sets whether the game should run in multiplayer mode
func WithMultiplayer(mp bool) Opt {
	return func(g *Game) error {
		g.multiplayer = mp

		return nil
	}
}

// WithUUID sets the internal unique identifier
func WithUUID(uuid string) Opt {
	return func(g *Game) error {
		g.uuid = uuid

		return nil
	}
}
