# Building the binary of the App
FROM golang:trixie AS build

# Install SQLite development libraries for Cgo in the build stage
# 'sqlite-dev' is the package name for Alpine
RUN apk add build-base
RUN apk add --no-cache sqlite-dev
# setting the workdir
WORKDIR /go/src/ettiHelper

# Copy all the Code and stuff to compile everything
COPY go.mod go.mod
COPY go.sum go.sum

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

COPY . .
ARG TARGETARCH
# Builds the application as a staticly linked one, to allow it to run on alpine
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o app .
RUN go build -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest AS release

WORKDIR /app
COPY --from=build /go/src/ettiHelper/app .

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

# Exposes port 3000 because our program listens on that port
EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "/app/app"]