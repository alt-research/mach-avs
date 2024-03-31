FROM golang:1.21 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod tidy && go mod verify

COPY . .

RUN make build-cli

RUN cp ./bin/mach-operator-cli /usr/local/bin/mach-operator-cli

FROM debian:latest as app
COPY --from=build /usr/local/bin/mach-operator-cli /usr/local/bin/mach-operator-cli
ENTRYPOINT [ "mach-operator-cli"]
