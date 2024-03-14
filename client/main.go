package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/plyr4/go-ebiten-multiplayer/shared/constants"
	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var DEFAULT_CONNECTION = Connection{
	Protocol:   "ws",
	Host:       "localhost:8080",
	ClientPath: "client",
}

type Connection struct {
	Protocol   string
	Host       string
	ClientPath string
}

// fills a websocket connection with defaults when necessary
func (c *Connection) Defaultify() {
	if len(c.Protocol) == 0 {
		logrus.Warnf("using default protocol: %s", DEFAULT_CONNECTION.Protocol)
		c.Protocol = DEFAULT_CONNECTION.Protocol
	}
	if len(c.Host) == 0 {
		logrus.Warnf("using default host: %s", DEFAULT_CONNECTION.Host)
		c.Host = DEFAULT_CONNECTION.Host
	}
	if len(c.ClientPath) == 0 {
		logrus.Warnf("using default client path: %s", DEFAULT_CONNECTION.ClientPath)
		c.ClientPath = DEFAULT_CONNECTION.ClientPath
	}
}

// craft a websocket URL
func (c *Connection) URL() string {
	return strings.Join(
		[]string{
			c.Protocol, "://", c.Host, "/", c.ClientPath,
		},
		"")
}

func main() {
	ws := Connection{
		Protocol:   os.Getenv("WS_CLIENT_PROTOCOL"),
		Host:       os.Getenv("WS_CLIENT_HOST"),
		ClientPath: os.Getenv("WS_CLIENT_PATH"),
	}
	ws.Defaultify()

	logger := logrus.NewEntry(logrus.StandardLogger()).WithFields(
		logrus.Fields{
			"protocol": ws.Protocol,
			"host":     ws.Host,
			"path":     ws.ClientPath,
		},
	)
	logrus.SetLevel(logrus.TraceLevel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := websocket.DialOptions{
		Subprotocols: []string{
			constants.CLIENT_SUBPROTOCOL,
		},
	}

	c, _, err := websocket.Dial(ctx, ws.URL(), &opts)
	if err != nil {
		logger.Errorf("unable to dial websocket connection: %v", err)

		return
	}
	defer c.CloseNow()

	msg := "hey bb!"
	logger.Tracef("data send: %s", msg)

	err = wsjson.Write(ctx, c, msg)
	if err != nil {
		logger.Errorf("unable to write: %v", err)

		return
	}

	out := ""
	err = wsjson.Read(ctx, c, &out)
	if err != nil {
		logger.Errorf("unable to read: %v", err)

		return
	}

	logger.Tracef("data recv: %s", out)

	logger.Info("closing connection")

	c.Close(websocket.StatusNormalClosure, "")
}
