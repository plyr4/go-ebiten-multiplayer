
FROM caddy/caddy:alpine as caddy

# not a serious project, leaving this on latest
FROM golang

# copy what we need for caddy
ADD Caddyfile /etc/caddy/Caddyfile
COPY --from=caddy /etc/caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy

# copy go files to be served via wasm
COPY . /go/src/

# copy wasm wrapper that serves Go
COPY servewasm.sh /usr/bin/servewasm

# cd to the directory where the go files are
WORKDIR /go/src/

# expose ports
EXPOSE 8080
EXPOSE 8091

CMD [ "servewasm" ]
