FROM golang:1.21 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod tidy && go mod verify

COPY . .

WORKDIR /usr/src/app/aggregator/cmd
RUN go build -v -o /usr/local/bin/aggregator ./...

FROM debian:latest
COPY --from=build /usr/local/bin/aggregator /usr/local/bin/aggregator
ENTRYPOINT [ "aggregator"]
CMD ["--config=/app/avs_config.yaml --avs-deployment /app/machavs_deploy_output.json"]
