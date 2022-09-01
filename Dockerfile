FROM golang:alpine as build

RUN apk add zeromq-dev build-base pkgconf
WORKDIR /sidearm

COPY . .
RUN go build

# Use a trimmed deployment image without the Go toolchain or source
FROM alpine:latest

RUN apk add zeromq
WORKDIR /

COPY --from=build /sidearm/sidearm /

ENTRYPOINT ["/sidearm"]
