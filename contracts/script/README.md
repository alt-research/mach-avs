# Deploy script

## 1. Boot a test environment by anvil

```bash
anvil
```

```bash
export OPERATOR_ADDR=0x957D781ab2Bc6D27Fde0a0b427ebF46ee1395661
```

### 1.1 Deploy Eigenlayer contracts

```bash
cd contracts

forge script script/EigenLayerDeployer.s.sol --broadcast -vvvv \
    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
    --rpc-url http://localhost:8545
```

the contract address is in `eigenlayer_deploy_output.json`.

```bash
export UNDERLAYING_TOKEN=0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44
```

### 1.1 Send token to operator

```bash
cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --value 50ether $OPERATOR_ADDR

cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 $UNDERLAYING_TOKEN 'transfer(address, uint256) (bool)' $OPERATOR_ADDR 100000000000
```

```bash
cast call $UNDERLAYING_TOKEN 'balanceOf(address) (uint256)' 0x957D781ab2Bc6D27Fde0a0b427ebF46ee1395661
100000000000 [1e11]
```

### 1.2 Deploy Avs contracts

```bash
STRATEGY=0x4A679253410272dd5232B3Ff7cF5dbB88f295319 AVS_DIRECTORY=0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9 DELEGATION_MANAGER=0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 \
forge script script/MachServiceManagerDeployer.s.sol \
--private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
--broadcast -vvvv --rpc-url http://localhost:8545
```

### 1.3 Reg to eigenlayer

First should change the addresses in `operator.yaml`:

```yaml
# This is the address of the slasher which is deployed in the anvil saved state
# The saved eigenlayer state is located in tests/anvil/credible_squaring_avs_deployment_output.json
# TODO(samlaf): automate updating these addresses when we deploy new contracts
avs_registry_coordinator_address: 0xa82fF9aFd8f496c3d6ac40E2a0F282E47488CFc9
operator_state_retriever_address: 0x0E801D84Fa97b50751Dbf25036d067dCf18858bF

# EigenLayer Delegation Manager contract address
# This will be provided by EigenLayer team
el_delegation_manager_address: 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707
```

Then reg:

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml rel
```

### 1.4 Add `strategyBaseTVLLimits` to `strategyWhitelister`

```bash
export STRATEGYMANAGER_ADDR=0x0165878A594ca255338adfa4d48449f69242Eb8F
```

Add `strategyWhitelister`, use `strategyManager`

get:

```bash
cast call $STRATEGYMANAGER_ADDR 'strategyWhitelister() (address)'
0x0000000000000000000000000000000000000000
```

set by owner:

```bash
cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 $STRATEGYMANAGER_ADDR \
'setStrategyWhitelister(address)' 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

get will be owner:

```bash
cast call 0x0165878A594ca255338adfa4d48449f69242Eb8F 'strategyWhitelister() (address)'
0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

set the `strategyBaseTVLLimits` to `strategyWhitelister`

> Note the `strategyBaseTVLLimits` address is in json.

```bash
cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
$STRATEGYMANAGER_ADDR \
'addStrategiesToDepositWhitelist(address[], bool[])' \
'[0x4A679253410272dd5232B3Ff7cF5dbB88f295319]' '[false]' 
```

### 1.5 Deposit tokens into a strategy

> `--strategy-addr` need use `strategyBaseTVLLimits`

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml d --strategy-addr 0x4A679253410272dd5232B3Ff7cF5dbB88f295319 --amount 10000000
```

### 1.6 Reg to AVS

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml r   
```

### 1.7 Boot operator and aggregator

```bash
./bin/mach-aggregator --config ./config-files/aggregator.yaml --ecdsa-private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --avs-deployment ./contracts/script/output/machavs_deploy_output.json
```

```bash
./bin/mach-operator-signer --config ./config-files/operator.yaml 
```

### 1.8 Call a fake alert for test

the port for operator is:

```yaml
operator_server_ip_port_addr: localhost:8091
```

```bash
curl --noproxy '*' -H "Content-Type: application/json" \
  -d '{"id":2, "jsonrpc":"2.0", "method": "alert_blockMismatch", "params":{"invalid_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "expect_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "l2_block_number": 2000}}' \
  http://localhost:8091
```