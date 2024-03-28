#!/bin/bash

cd ./contracts

forge script ./script/EigenLayerDeployer.s.sol --broadcast -vvvv \
    --private-key $OWNER_PRIVATE \
    --rpc-url $RPC_URL

forge script ./script/MachServiceManagerDeployer.s.sol \
    --private-key $OWNER_PRIVATE \
    --broadcast -vvvv --rpc-url $RPC_URL

cd ..

OPERATOR_ADDR=0xE9A7669aC9eBE9b7E21E0A323FC3A6f34CE744eb
OPERATOR_KEY_NAME=test1

EIGENLAYER_ADDRESS_PATH='./contracts/script/output/eigenlayer_deploy_output.json'
AVS_ADDRESS_PATH='./contracts/script/output/machavs_deploy_output.json'

UNDERLAYING_TOKEN=$(cat $EIGENLAYER_ADDRESS_PATH | jq -r '.underlayingToken' )
REGISTRY_COORDINATOR_ADDR=$(cat $AVS_ADDRESS_PATH | jq -r '.registryCoordinator' )
OPERATOR_STATE_RETRIEVER_ADDR=$(cat $AVS_ADDRESS_PATH | jq -r '.operatorStateRetriever' )
STRATEGY_BASE_TVL_LIMITS_ADDR=$(cat $EIGENLAYER_ADDRESS_PATH | jq -r '.strategyBaseTVLLimits' )

cast send  -f $OWNER_ADDR --private-key $OWNER_PRIVATE --rpc-url $RPC_URL --value 2ether $OPERATOR_ADDR
cast send  -f $OWNER_ADDR --private-key $OWNER_PRIVATE $UNDERLAYING_TOKEN --rpc-url $RPC_URL 'transfer(address, uint256) (bool)' $OPERATOR_ADDR 100000000000 

# Reset the config
sed -i 's/avs_registry_coordinator_address: .\+/avs_registry_coordinator_address: '${REGISTRY_COORDINATOR_ADDR}'/g' ./ops/configs/operator-docker-compose.yaml
sed -i 's/operator_state_retriever_address: .\+/operator_state_retriever_address: '${OPERATOR_STATE_RETRIEVER_ADDR}'/g' ./ops/configs/operator-docker-compose.yaml
sed -i 's/ecdsa_private_key_store_path: .\+/ecdsa_private_key_store_path: .\/config-files\/key\/'${OPERATOR_KEY_NAME}'.ecdsa.key.json/g' ./ops/configs/operator-docker-compose.yaml
sed -i 's/bls_private_key_store_path: .\+/bls_private_key_store_path: .\/config-files\/key\/'${OPERATOR_KEY_NAME}'.bls.key.json/g' ./ops/configs/operator-docker-compose.yaml
sed -i 's/metadata_uri: .\+/metadata_uri: '${METADATA_URI}'/g' ./ops/configs/operator-docker-compose.yaml

./bin/mach-operator-cli --config ./ops/configs/operator-docker-compose.yaml rel
./bin/mach-operator-cli --config ./ops/configs/operator-docker-compose.yaml d --strategy-addr $STRATEGY_BASE_TVL_LIMITS_ADDR --amount 10000000
./bin/mach-operator-cli --config ./ops/configs/operator-docker-compose.yaml r

# Send eth to operator addr for make a new block
cast send  -f $OWNER_ADDR --private-key $OWNER_PRIVATE --rpc-url $RPC_URL --value 1ether $OPERATOR_ADDR
