# layers
FROM caddy/caddy:alpine as caddy

# base
FROM golang:1.22.1

# caddy
COPY --from=caddy /etc/caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy
COPY Caddyfile /etc/caddy/Caddyfile

# go
COPY . /go/src/
WORKDIR /go/src/

# expose ports
EXPOSE 8080

# run
ADD entrypoint.sh /bin/entrypoint
CMD ["/bin/entrypoint"]