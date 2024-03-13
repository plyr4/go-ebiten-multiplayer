FROM caddy/caddy:alpine as caddy

ADD Caddyfile /etc/caddy/Caddyfile

# not a serious project, leaving this on latest
FROM golang

# copy what we need from caddy layer
COPY --from=caddy /etc/caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy

# copy go files to be served via wasm
COPY . /go/src/

# copy wasm wrapper that serves Go
COPY servewasm.sh /usr/bin/servewasm

# cd to the directory where the go files are
WORKDIR /go/src/

EXPOSE 8080

CMD [ "servewasm" ]