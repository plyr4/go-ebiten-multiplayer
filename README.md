# go-ebiten-multiplayer

A small game made with [Go](https://go.dev/) and [ebitengine](https://ebitengine.org/).

## Development

Run the game locally.

### Prerequisites

At a minimum:
- [make](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://docs.docker.com/get-docker/)
- [docker-buildx](https://docs.docker.com/reference/cli/docker/buildx/)

For local dev:
- [Go](https://go.dev/)
- [ebitengine](https://ebitengine.org/) (make sure the environment test passes)

### Run

- Build and run via Docker.
    ```bash
    $ make up
    ```

_Or_

- Run via Go.
    ```bash
    # run server in the background
    make srv &
    # run the client
    make clt
    ```

- Visit http://localhost:8080 

- Disable multiplayer by setting `CLIENT_MULTIPLAYER` to `false`.
    ```bash
    export CLIENT_MULTIPLAYER=false
    ```

## Troubleshooting

Ports stuck after running and killing the server? Run `make kill-srv` then try again.
