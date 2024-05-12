FROM golang:1.21-bullseye as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod tidy && go mod verify

COPY . .

WORKDIR /usr/src/app/generic-operator-proxy/cmd
RUN go build -v -o /usr/local/bin/generic-operator-proxy ./...

FROM debian:bullseye as app
COPY --from=build /usr/local/bin/generic-operator-proxy /usr/local/bin/mach-operator-proxy

RUN apt-get update && \
        apt-get install --no-install-recommends -y curl sudo daemontools jq ca-certificates && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "mach-operator-proxy"]
