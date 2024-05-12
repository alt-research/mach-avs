FROM golang:1.21-bullseye as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

ENV GOPRIVATE=github.com/alt-research/avs-generic-aggregator
ARG XDG_CONFIG_HOME=/root/.config/
RUN \
    --mount=type=secret,id=gh_hosts,target=/root/.config/gh/hosts.yml \
    --mount=type=secret,id=git_config,target=/root/.gitconfig \
    --mount=type=secret,id=git_credentials,target=/root/.git-credentials \
    <<EOF
    set -ex
    go version
EOF

RUN go mod download && go mod tidy && go mod verify

COPY . .

WORKDIR /usr/src/app/aggregator/cmd
RUN go build -v -o /usr/local/bin/aggregator ./...

FROM debian:bullseye as app
COPY --from=build /usr/local/bin/aggregator /usr/local/bin/aggregator

RUN apt-get update && \
        apt-get install --no-install-recommends -y curl sudo daemontools jq ca-certificates && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "aggregator"]
