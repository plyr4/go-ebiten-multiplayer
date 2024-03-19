# go-ebiten-multiplayer

A small game multiplayer made with [Go](https://go.dev/), [ebitengine](https://ebitengine.org/) and websockets.

## Play

Web and linux demos are **_coming soon_**! For now, you need to compile and run the code locally to play.

## Development

How to run and debug the game from source.

### Prerequisites

At a minimum: [make](https://www.gnu.org/software/make/manual/make.html), [Docker](https://docs.docker.com/get-docker/), [docker-buildx](https://docs.docker.com/reference/cli/docker/buildx/).

For local dev: [Go](https://go.dev/), [ebitengine](https://ebitengine.org/) (make sure the environment test passes).

### Run via Docker

- Build and run via Docker and web assembly.
    ```bash
    make up
    ```
- Visit http://localhost:8080.

### Run via Go

- Run directly via Go.
    ```bash
    # run server in the background
    make srv &
    # run the client
    make clt
    ```

- Disable multiplayer by using `make clt-local` or by setting `CLIENT_MULTIPLAYER` to `false`.
    ```bash
    make clt-local
    # or
    export CLIENT_MULTIPLAYER=false
    make clt
    ```

## Troubleshooting

- Ports stuck after running and killing the server? Run `make kill-srv` then try again.

## TODOs

- [ ] strong client uuids
- [ ] websocket security
- [ ] sprite animations
- [ ] dynamic animations
- [ ] server self-cleanup
- [ ] multiplayer lobbies
- [ ] ui
- [ ] player customization
