############################# HELP MESSAGE #############################
# Make sure the help command stays first, so that it's printed by default when `make` is called without arguments
.PHONY: help tests
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

AGGREGATOR_ECDSA_PRIV_KEY=0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6
CHALLENGER_ECDSA_PRIV_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a

CHAINID=31337
# Make sure to update this if the strategy address changes
# check in contracts/script/output/${CHAINID}/credible_squaring_avs_deployment_output.json
STRATEGY_ADDRESS=0x7a2088a1bFc9d81c55368AE168C2C02570cB814F
DEPLOYMENT_FILES_DIR=contracts/script/output/${CHAINID}

PROTOS := ./api/proto
PROTO_GEN := ./api/grpc

-----------------------------: ## 

___CONTRACTS___: ## 

build-contracts: ## builds all contracts
	cd contracts && forge build

bindings: ## generates contract bindings
	cd contracts && bash generate-go-bindings.sh

__CLI__: ## 

clean:
	find $(PROTO_GEN) -name "*.pb.go" -type f | xargs rm -rf
	mkdir -p $(PROTO_GEN)

protoc: clean
	protoc -I $(PROTOS) \
	--go_out=$(PROTO_GEN) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_GEN) \
	--go-grpc_opt=paths=source_relative \
	$(PROTOS)/**/*.proto

lint:
	staticcheck ./...
	golangci-lint run

build: build-operator build-aggregator build-cli build-operator-proxy

build-operator:
	go build -o ./bin/mach-operator-signer ./legacy/operator/cmd 

build-operator-proxy:
	go build -o ./bin/mach-operator-proxy ./generic-operator-proxy/cmd

build-aggregator:
	go build -o ./bin/mach-aggregator ./legacy/aggregator/cmd 

build-cli:
	go build -o ./bin/mach-operator-cli ./legacy/cli 

_____HELPER_____: ## 
mocks: ## generates mocks for tests
	go install go.uber.org/mock/mockgen@v0.3.0
	go generate ./...

