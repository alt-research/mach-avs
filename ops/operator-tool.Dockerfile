FROM golang:1.21-bullseye as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod tidy && go mod verify

COPY . .

RUN make build-cli

RUN cp ./bin/mach-operator-cli /usr/local/bin/mach-operator-cli

FROM debian:bullseye as app
COPY --from=build /usr/local/bin/mach-operator-cli /usr/local/bin/mach-operator-cli

RUN apt-get update && \
        apt-get install --no-install-recommends -y curl sudo daemontools jq ca-certificates && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "mach-operator-cli"]
