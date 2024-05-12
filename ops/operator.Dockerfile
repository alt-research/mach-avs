FROM golang:1.21-bullseye as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

ENV GOPRIVATE=github.com/alt-research/avs-generic-aggregator

RUN go mod download && go mod tidy && go mod verify

COPY . .

WORKDIR /usr/src/app/operator/cmd
RUN go build -v -o /usr/local/bin/operator ./...

FROM debian:bullseye as app
COPY --from=build /usr/local/bin/operator /usr/local/bin/operator

RUN apt-get update && \
        apt-get install --no-install-recommends -y curl sudo daemontools jq ca-certificates && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "operator"]
