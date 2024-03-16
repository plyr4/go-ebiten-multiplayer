package ws

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/shared/constants"

	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var DEFAULT_WS_CONFIG = Config{
	Protocol:   "ws",
	Host:       "localhost:8080",
	ClientPath: "client",
}

type Config struct {
	Protocol   string
	Host       string
	ClientPath string
}

// ValidateAndNullify fills a websocket Client with defaults when necessary
func (cfg *Config) ValidateAndNullify() error {
	if len(cfg.Protocol) == 0 {
		logrus.Warnf("using default protocol: %s", DEFAULT_WS_CONFIG.Protocol)
		cfg.Protocol = DEFAULT_WS_CONFIG.Protocol
	}

	if len(cfg.Host) == 0 {
		logrus.Warnf("using default host: %s", DEFAULT_WS_CONFIG.Host)
		cfg.Host = DEFAULT_WS_CONFIG.Host
	}

	if len(cfg.ClientPath) == 0 {
		logrus.Warnf("using default client path: %s", DEFAULT_WS_CONFIG.ClientPath)
		cfg.ClientPath = DEFAULT_WS_CONFIG.ClientPath
	}

	return nil
}

// Client houses a websocket connection
type Client struct {
	ctx        context.Context
	config     *Config
	connection *websocket.Conn
	logger     *logrus.Entry
}

// New creates a new Client from the environment
func New() *Client {
	c := new(Client)

	// context
	c.WithContext(context.Background())

	// config
	cfg := &Config{
		Protocol:   os.Getenv("CLIENT_WS_PROTOCOL"),
		Host:       os.Getenv("CLIENT_WS_HOST"),
		ClientPath: os.Getenv("CLIENT_WS_PATH"),
	}
	cfg.ValidateAndNullify()
	c.config = cfg

	// logging
	logger := logrus.NewEntry(logrus.StandardLogger()).WithFields(
		logrus.Fields{
			"protocol": cfg.Protocol,
			"host":     cfg.Host,
			"path":     cfg.ClientPath,
		},
	)
	c.WithLogger(logger)

	return c
}

// WithContext attaches a context to the Client
func (c *Client) WithContext(ctx context.Context) {
	c.ctx = ctx
}

// WithLogger attaches a logger to the Client
func (c *Client) WithLogger(l *logrus.Entry) {
	c.logger = l
}

// crafts a websocket URL
func (c *Client) URL() string {
	return strings.Join(
		[]string{
			c.config.Protocol, "://", c.config.Host, "/", c.config.ClientPath,
		},
		"")
}

// create a websocket connection
func (c *Client) Connect() error {
	c.logger.Infof("connecting to: %s", c.URL())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := websocket.DialOptions{
		Subprotocols: []string{
			constants.CLIENT_SUBPROTOCOL,
		},
	}

	var err error
	c.connection, _, err = websocket.Dial(ctx, c.URL(), &opts)
	if err != nil || c.connection == nil {
		if c.connection == nil {
			err = errors.New("connection is nil")
		}

		e := errors.Wrap(err, "unable to dial")

		c.logger.Error(e)

		return e
	}

	return nil
}

func (c *Client) Reconnect() error {
	c.logger.Trace("reconnecting ws client")

	if c.connection != nil {
		c.connection.CloseNow()
	}

	return c.Connect()
}

func (c *Client) Close(msg string) error {
	c.logger.Tracef("closing ws client: %s", msg)

	if c.connection != nil {
		return c.connection.Close(websocket.StatusNormalClosure, msg)
	}

	return errors.New("cannot close a nil connection")
}

func (c *Client) Send(msg interface{}) error {
	c.logger.Tracef("sending msg: %v", msg)

	err := wsjson.Write(c.ctx, c.connection, msg)
	if err != nil {
		e := errors.Wrap(err, "unable to write")

		c.logger.Error(e)

		return e
	}

	c.logger.Tracef("sent msg: %v", msg)

	return nil
}

func (c *Client) Receive(msg interface{}) (interface{}, error) {
	c.logger.Trace("receiving msg")

	err := wsjson.Read(c.ctx, c.connection, &msg)
	if err != nil {
		e := errors.Wrap(err, "unable to read")

		c.logger.Error(e)

		return nil, e
	}

	c.logger.Tracef("received msg: %v", msg)

	return msg, nil
}
