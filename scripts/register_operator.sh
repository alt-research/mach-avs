#!/bin/bash

OWNER_ADDR=$1
OWNER_PRIVATE=$2
# Calc ADDR from OPERATOR_KEY_NAME
OPERATOR_ADDR=$3
OPERATOR_KEY_NAME=$4

EIGENLAYER_ADDRESS_PATH='./contracts/script/output/eigenlayer_deploy_output.json'
AVS_ADDRESS_PATH='./contracts/script/output/machavs_deploy_output.json'

echo $OPERATOR_ADDR
echo $RPC_URL

UNDERLAYING_TOKEN=$(cat $EIGENLAYER_ADDRESS_PATH | jq -r '.underlayingToken' )
REGISTRY_COORDINATOR_ADDR=$(cat $AVS_ADDRESS_PATH | jq -r '.registryCoordinator' )
OPERATOR_STATE_RETRIEVER_ADDR=$(cat $AVS_ADDRESS_PATH | jq -r '.operatorStateRetriever' )
STRATEGY_BASE_TVL_LIMITS_ADDR=$(cat $EIGENLAYER_ADDRESS_PATH | jq -r '.strategyBaseTVLLimits' )

echo $UNDERLAYING_TOKEN

cast send  -f $OWNER_ADDR --private-key $OWNER_PRIVATE --rpc-url $RPC_URL --value 10ether $OPERATOR_ADDR
cast send  -f $OWNER_ADDR --private-key $OWNER_PRIVATE $UNDERLAYING_TOKEN --rpc-url $RPC_URL 'transfer(address, uint256) (bool)' $OPERATOR_ADDR 100000000000 

# Reset the config
sed -i 's/avs_registry_coordinator_address: .\+/avs_registry_coordinator_address: '${REGISTRY_COORDINATOR_ADDR}'/g' ./config-files/operator.yaml
sed -i 's/operator_state_retriever_address: .\+/operator_state_retriever_address: '${OPERATOR_STATE_RETRIEVER_ADDR}'/g' ./config-files/operator.yaml
sed -i 's/ecdsa_private_key_store_path: .\+/ecdsa_private_key_store_path: .\/config-files\/key\/'${OPERATOR_KEY_NAME}'.ecdsa.key.json/g' ./config-files/operator.yaml
sed -i 's/bls_private_key_store_path: .\+/bls_private_key_store_path: .\/config-files\/key\/'${OPERATOR_KEY_NAME}'.bls.key.json/g' ./config-files/operator.yaml

./bin/mach-operator-cli --config ./config-files/operator.yaml rel
echo './bin/mach-operator-cli --config ./config-files/operator.yaml d'
./bin/mach-operator-cli --config ./config-files/operator.yaml d --strategy-addr $STRATEGY_BASE_TVL_LIMITS_ADDR --amount 10000000
echo 'reg avs'
./bin/mach-operator-cli --config ./config-files/operator.yaml r
