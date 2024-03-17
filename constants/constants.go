package constants

import "time"

const (
	SERVER_WS_DEFAULT_HOST = "localhost:8091"
	// time wait between serving messages per-client
	SERVER_WS_LATENCY = 1 * time.Millisecond

	CLIENT_SUBPROTOCOL = "client"
	// time wait between pinging the server
	CLIENT_WS_LATENCY = 1 * time.Millisecond

	WINDOW_TITLE  = "Go Ebiten Multiplayer"
	WINDOW_WIDTH  = 1280
	WINDOW_HEIGHT = 720
)
