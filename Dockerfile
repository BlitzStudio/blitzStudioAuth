# Building the binary of the App
FROM golang:trixie AS build

WORKDIR /go/src/blitzStudioAuth

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} go build -o blitzStudioAuth .

FROM alpine:3.21 AS release

WORKDIR /app
COPY --from=build /go/src/blitzStudioAuth/blitzStudioAuth /app/app

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "/app/app"]