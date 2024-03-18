#!/bin/sh
set -x

# wasm
export WASM_HTTP_PORT=${WASM_HTTP_PORT:-8090}

# server (defaults set in script due to wasm wrapper limitations)
export SERVER_WS_HOST="${SERVER_WS_HOST:-0.0.0.0:8091}"

# client (defaults set in code)
export CLIENT_WS_PROTOCOL=${CLIENT_WS_PROTOCOL:-ws}
export CLIENT_WS_PROTOCOL=${CLIENT_WS_PROTOCOL:-ws}
export CLIENT_WS_HOST="${CLIENT_WS_HOST}"
export CLIENT_WS_PATH="${CLIENT_WS_PATH}"
export CLIENT_MULTIPLAYER="${CLIENT_MULTIPLAYER:-true}"

go run cmd/server/main.go &

go run github.com/hajimehoshi/wasmserve@latest -allow-origin='*' -http=:$WASM_HTTP_PORT ./cmd/client/main.go &

caddy run
