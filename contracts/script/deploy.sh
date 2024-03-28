#!/bin/bash

forge script ./script/EigenLayerDeployer.s.sol --broadcast -vvvv \
    --private-key $OWNER_PRIVATE \
    --rpc-url $RPC_URL

forge script ./script/MachServiceManagerDeployer.s.sol \
    --private-key $OWNER_PRIVATE \
    --broadcast -vvvv --rpc-url $RPC_URL

