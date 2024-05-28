#!/bin/bash

forge script ./script/EigenLayerDeployer.s.sol \
    --private-key $OWNER_PRIVATE \
    --broadcast -vvvv --slow --rpc-url $RPC_URL

forge script ./script/MachServiceManagerDeployer.s.sol \
    --private-key $OWNER_PRIVATE \
    --broadcast -vvvv --slow --rpc-url $RPC_URL
