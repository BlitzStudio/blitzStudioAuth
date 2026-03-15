# Building the binary of the App
FROM golang:trixie AS build

WORKDIR /go/src/ettiHelper

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
ARG TARGETARCH
RUN go build -o app .


FROM alpine:latest AS release

WORKDIR /app
COPY --from=build /go/src/ettiHelper/app /app/app

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "/app/app"]