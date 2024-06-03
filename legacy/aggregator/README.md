# Aggregator

config:

```yaml
# 'production' only prints info and above. 'development' also prints debug
environment: production

eth_rpc_url: https://ethereum-holesky-rpc.publicnode.com
eth_ws_url: wss://ethereum-holesky-rpc.publicnode.com

# address which the aggregator listens on for operator signed messages
aggregator_server_ip_port_address: 0.0.0.0:8090

# the layer1 chain id the avs contracts in
layer1_chain_id: 17000

# the layer2 chain id
layer2_chain_id: 20240219

```

envs:

```bash
AVS_REGISTRY_COORDINATOR_ADDRESS=0x1eA7D160d325B289bF981e0D7aB6Bf3261a0FFf2
OPERATOR_STATE_RETRIEVER_ADDRESS=0xBE1c904525910fdB49dB33b4960DF9aC9f603dC7
```
