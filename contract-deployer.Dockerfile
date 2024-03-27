# syntax=docker/dockerfile:1.4

# From https://github.com/foundry-rs/foundry/blob/master/Dockerfile

FROM ghcr.io/foundry-rs/foundry:latest as foundry-client

RUN apk add --no-cache jq

ENTRYPOINT ["/bin/sh", "-c"]
