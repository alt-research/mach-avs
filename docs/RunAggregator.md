# Build and boot the aggregator

## Boot aggregator

Need configuration the aggregator:

```yaml
# 'production' only prints info and above. 'development' also prints debug
environment: development
eth_rpc_url: https://ethereum-holesky-rpc.publicnode.com
eth_ws_url: wss://ethereum-holesky-rpc.publicnode.com

# address which the aggregator json rpc listens on for operator signed messages
aggregator_jsonrpc_server_ip_port_address: 0.0.0.0:8290

rpc_vhosts: ["*", "localhost"]
rpc_cors: ["*", "localhost"]

# the layer1 chain id the avs contracts in
layer1_chain_id: 17000

# the layer2 chain id
layer2_chain_id: 10

# the QuorumNums we use, just no change
quorum_nums: [0]
```

Then can boot the aggregator:

```bash
AVS_REGISTRY_COORDINATOR_ADDRESS=0x32ff2a0a8C035dbfe28f38870e9f20c9391D7907 \
OPERATOR_STATE_RETRIEVER_ADDRESS=0xBfC0Ea531fa3db2aD36a3C4A205C02d4a2dd9fa0 \
./bin/mach-aggregator --config ./config-files/aggregator.yaml --ecdsa-private-key $PRIVATE_KEY
```

The `PRIVATE_KEY` should be the committer for Mach AVS.
