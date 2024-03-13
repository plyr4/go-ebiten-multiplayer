#!/bin/sh
set -x

WASM_HTTP_PORT=8090

cd server

go run main.go &

cd ..

go run github.com/hajimehoshi/wasmserve@latest -http=:$WASM_HTTP_PORT ./client/main.go &


caddy run --config /etc/caddy/Caddyfile
