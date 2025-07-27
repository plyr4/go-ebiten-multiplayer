# go-ebiten-multiplayer

A small multiplayer game prototype made with [Go](https://go.dev/), [ebitengine](https://ebitengine.org/) and websockets.

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

## Roadmap

The project roadmap has been converted to GitHub issues for better tracking and collaboration. See the [Issues page](https://github.com/plyr4/go-ebiten-multiplayer/issues) for current development tasks.

### Key Features in Development

- **Security** - Client UUID and WebSocket security improvements
- **Graphics & Animation** - Sprite animations and dynamic animation system  
- **Server** - Self-cleanup and maintenance features
- **Multiplayer** - Lobby system for organized game sessions
- **UI/UX** - User interface and player customization features

For detailed information about converting TODOs to issues, see [TODO_TO_ISSUES.md](./TODO_TO_ISSUES.md).
