
# ------------------------------------------------------------------------------
# GOLANG Build Stage
# ------------------------------------------------------------------------------

FROM golang:1.12.4-alpine3.9 AS golang-builder

RUN apk add bash git gcc g++ libc-dev

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /usr/src/microservice_gokit_base

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app.bin

# ------------------------------------------------------------------------------
# Final Stage
# ------------------------------------------------------------------------------

FROM alpine:3.9

RUN addgroup -g 1000 docker

RUN adduser -D -s /bin/sh -u 1000 -G docker appuser

WORKDIR /microservice_gokit_base/bin/

COPY --from=golang-builder /usr/src/microservice_gokit_base/app.bin .

RUN chown appuser:docker app.bin

USER appuser

EXPOSE 8080

CMD ["./app.bin"]