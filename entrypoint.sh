#!/bin/sh
set -x

# wasm
export WASM_HTTP_PORT=${WASM_HTTP_PORT:-8090}

# server (defaults set in script due to wasm wrapper limitations)
export WS_SERVER_HOST="${WS_SERVER_HOST:-0.0.0.0:8091}"

# client (defaults set in code)
export WS_CLIENT_PROTOCOL=${WS_CLIENT_PROTOCOL:-ws}
export WS_CLIENT_HOST="${WS_CLIENT_HOST}"
export WS_CLIENT_PATH="${WS_CLIENT_PATH}"

go run server/main.go &

go run github.com/hajimehoshi/wasmserve@latest -allow-origin='*' -http=:$WASM_HTTP_PORT ./client/main.go &

caddy run
