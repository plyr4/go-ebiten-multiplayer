#!/bin/sh
set -x

WASM_HTTP_PORT=8090

go run github.com/hajimehoshi/wasmserve@latest -http=:$WASM_HTTP_PORT ./main.go &

caddy run --config /etc/caddy/Caddyfile
