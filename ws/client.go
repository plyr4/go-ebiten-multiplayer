package ws

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"golang.org/x/time/rate"

	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var DEFAULT_WS_CONFIG = Config{
	Protocol:   "ws",
	Host:       constants.SERVER_WS_DEFAULT_HOST,
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
	ctx         context.Context
	config      *Config
	connection  *websocket.Conn
	logger      *logrus.Entry
	rateLimiter *rate.Limiter
}

// New creates a new Client from the environment
func New(logger *logrus.Entry, opts ...Opt) (*Client, error) {
	c := new(Client)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// config
	cfg := &Config{
		Protocol:   os.Getenv("CLIENT_WS_PROTOCOL"),
		Host:       os.Getenv("CLIENT_WS_HOST"),
		ClientPath: os.Getenv("CLIENT_WS_PATH"),
	}

	cfg.ValidateAndNullify()

	c.config = cfg

	// logging
	logger = logger.
		WithFields(
			logrus.Fields{
				"module":   "ws-client",
				"protocol": cfg.Protocol,
				"host":     cfg.Host,
				"path":     cfg.ClientPath,
			},
		)
	c.logger = logger

	// rate limiting
	latency := constants.CLIENT_WS_LATENCY

	c.rateLimiter = rate.NewLimiter(rate.Every(latency), 10)

	return c, nil
}

type Opt func(*Client) error

// WithContext sets the internal context
func WithContext(ctx context.Context) Opt {
	return func(g *Client) error {
		g.ctx = ctx

		return nil
	}
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
	c.logger.Tracef("connecting to: %s", c.URL())

	err := c.rateLimiter.Wait(c.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for rate limiter")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := websocket.DialOptions{
		Subprotocols: []string{
			constants.CLIENT_SUBPROTOCOL,
		},
	}

	c.connection, _, err = websocket.Dial(ctx, c.URL(), &opts)
	if err != nil || c.connection == nil {
		if c.connection == nil {
			if err == nil {
				err = errors.New("connection is nil")
			} else {
				err = errors.Wrap(err, "connection is nil")
			}
		}

		return errors.Wrap(err, "unable to dial")
	}

	return nil
}

func (c *Client) Reconnect() error {
	c.logger.Tracef("reconnecting to: %s", c.URL())

	err := c.rateLimiter.Wait(c.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for rate limiter")
	}

	if c.connection != nil {
		c.connection.CloseNow()
	}

	return c.Connect()
}

func (c *Client) Close(msg string) error {
	c.logger.Tracef("closing ws client: %s", msg)

	if c.connection == nil {
		return errors.New("connection is nil")
	}

	if c.connection != nil {
		return c.connection.Close(websocket.StatusNormalClosure, msg)
	}

	return nil
}

func (c *Client) IsConnected() bool {
	if c.connection == nil {
		return false
	}

	return c.connection.Ping(c.ctx) == nil
}

func (c *Client) Send(msg interface{}) error {
	c.logger.Tracef("sending msg: %v", msg)

	err := c.rateLimiter.Wait(c.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for rate limiter")
	}

	if c.connection == nil {
		return errors.New("connection is nil")
	}

	err = wsjson.Write(c.ctx, c.connection, msg)
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

	err := c.rateLimiter.Wait(c.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for rate limiter")
	}

	if c.connection == nil {
		return nil, errors.New("connection is nil")
	}

	err = wsjson.Read(c.ctx, c.connection, &msg)
	if err != nil {
		e := errors.Wrap(err, "unable to read")

		c.logger.Error(e)

		return nil, e
	}

	c.logger.Tracef("received msg: %v", msg)

	return msg, nil
}
