# Run Operator and Aggregator for Mach AVS

## Generate the keys for operator

Need use [eigenlayer-cli](https://github.com/Layr-Labs/eigenlayer-cli):

```bash
curl -sSfL https://raw.githubusercontent.com/layr-labs/eigenlayer-cli/master/scripts/install.sh | sh -s
export PATH=$PATH:~/bin
```

Can run eigenlayer:

```bash
 eigenlayer                                                                             
        
     _______ _                   _                              
    (_______|_)                 | |                             
     _____   _  ____  ____ ____ | |      ____ _   _  ____  ____ 
    |  ___) | |/ _  |/ _  )  _ \| |     / _  | | | |/ _  )/ ___)
    | |_____| ( ( | ( (/ /| | | | |____( ( | | |_| ( (/ /| |    
    |_______)_|\_|| |\____)_| |_|_______)_||_|\__  |\____)_|    
              (_____|                        (____/             
    NAME:
   eigenlayer - EigenLayer CLI

USAGE:
   eigenlayer [global options] command [command options] 

VERSION:
   0.6.2

COMMANDS:
   operator  Execute onchain operations for the operator
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   (c) 2024 EigenLabs
```

Then we can generate keys:

```bash
eigenlayer operator keys create --key-type ecdsa [keyname]
eigenlayer operator keys create --key-type bls [keyname]
```

- [keyname] - This will be the name of the created key file. It will be saved as <keyname>.ecdsa.key.json or <keyname>.bls.key.json.

The details can see [operator-guides](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#create-keys)


## Configuration For tools to register avs operator

```yaml
# this sets the logger level (true = info, false = debug)
production: true

# ETH RPC URL
eth_rpc_url: https://ethereum-holesky-rpc.publicnode.com
eth_ws_url: wss://ethereum-holesky-rpc.publicnode.com

# If you running this using eigenlayer CLI and the provided AVS packaging structure,
# this should be /operator_keys/ecdsa_key.json as the host path will be asked while running
#
# If you are running locally using go run main.go, this should be full path to your local ecdsa key file
ecdsa_private_key_store_path: ./config-files/key/opr.ecdsa.key.json

# If you running this using eigenlayer CLI and the provided AVS packaging structure,
# this should be /operator_keys/bls_key.json as the host path will be asked while running
#
# We are using bn254 curve for bls keys
#
# If you are running locally using go run main.go, this should be full path to your local bls key file
bls_private_key_store_path: ./config-files/key/opr.bls.key.json

# avs node spec compliance https://eigen.nethermind.io/docs/spec/intro
eigen_metrics_ip_port_address: localhost:9090
enable_metrics: true
node_api_ip_port_address: localhost:9010
enable_node_api: true

# EigenLayer Slasher contract address

# This is the address of the slasher which is deployed in the anvil saved state
avs_registry_coordinator_address: 0x32ff2a0a8C035dbfe28f38870e9f20c9391D7907
operator_state_retriever_address: 0xBfC0Ea531fa3db2aD36a3C4A205C02d4a2dd9fa0
```

> Note for operator tool, we need `ecdsa_private_key_store_path` for sign by ecdsa address.

For Mainnet:

- `eth_rpc_url`: https://ethereum-rpc.publicnode.com
- `eth_ws_url`: wss://ethereum-rpc.publicnode.com

The addresses of the contracts can see the [README](../README.md), the `avs_registry_coordinator_address` need use `RegistryCoordinator`,
the `operator_state_retriever_address` need use `OperatorStateRetriever`.


## Configuration

```yaml
# this sets the logger level (true = info, false = debug)
production: true

# ETH RPC URL
eth_rpc_url: https://ethereum-holesky-rpc.publicnode.com
eth_ws_url: wss://ethereum-holesky-rpc.publicnode.com

# If you running this using eigenlayer CLI and the provided AVS packaging structure,
# this should be /operator_keys/bls_key.json as the host path will be asked while running
#
# We are using bn254 curve for bls keys
#
# If you are running locally using go run main.go, this should be full path to your local bls key file
bls_private_key_store_path: ./config-files/key/opr.bls.key.json

# The operator 's ecdsa address, Note the operator node not need ecdsa key file, just set the address,
# can got this address by `eigenlayer operator keys list `
operator_ecdsa_address: 0xad6b95793dd4d2b8e184fb4666d1cfb14871a035

aggregator_jsonrpc_server_ip_port_address: http://localhost:8290

# avs node spec compliance https://eigen.nethermind.io/docs/spec/intro
eigen_metrics_ip_port_address: localhost:9090
enable_metrics: true
node_api_ip_port_address: localhost:9010
enable_node_api: true

# EigenLayer Slasher contract address

# This is the address of the slasher which is deployed in the anvil saved state
avs_registry_coordinator_address: 0x32ff2a0a8C035dbfe28f38870e9f20c9391D7907
operator_state_retriever_address: 0xBfC0Ea531fa3db2aD36a3C4A205C02d4a2dd9fa0

operator_server_ip_port_addr: localhost:8091

# the layer1 chain id the avs contracts in
layer1_chain_id: 17000

# the layer2 chain id
layer2_chain_id: 10

```

the `layer1_chain_id` and `layer2_chain_id` should match to the AVS.

## Register the operator to Mach AVS

First, the operator should registered into eigenlayer. this can see [Operator Registration](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#operator-registration)

> Note: if the key need password, should use `OPERATOR_BLS_KEY_PASSWORD` and `OPERATOR_ECDSA_KEY_PASSWORD` environment variables.

Reg:

```bash
 ./bin/mach-operator-cli --config ./config-files/operator.yaml register-operator-with-avs
```

Search the status for operator:

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml  print-operator-status

...

Printing operator status
{
 "EcdsaAddress": "0xAD6B95793DD4D2b8e184FB4666D1cfb14871A035",
 "PubkeysRegistered": true,
 "G1Pubkey": "E([498211989701534593628498974128726712526336918939770789545660245177948853517,19434346619705907282579203143605058653932187676054178921788041096426532277474])",
 "G2Pubkey": "E([18245660679142793594022062571710409567727998469085890965936881880506455511441+15102214288340415892913544550311375043597835055349681074674493756729384353913*u,20786923013429232749470594329455309882713682216761797910592284729738312371850+2073225601630511023793984244584721446296103509299903599739377312253320009110*u])",
 "RegisteredWithAvs": true,
 "OperatorId": "5a76fe9014f9cd296a69ac589c2bbd2c6a354c5e4c0c79ee35c5b8202b8523a2"
}
```

## Boot the operator

If is ok, can boot the operator.

```bash
./bin/mach-operator-signer --config ./config-files/operator.yaml
```

## Boot Mach verifier

We need boot a mach verifier for the operator to generate alert, can see the [Mach](https://github.com/alt-research/mach).

Note the mach 's alert config:

```yaml
##############################################
# Mach alerter configuration                 #
##############################################
[alerter]
# Only print the warning log if the alerter is disabled,
# Otherwise the alerts will be sent to the Mach AVS Operator
# Optional, default value: false
enable = false
# The JSONRPC endpoint of Mach AVS Operator
url = "http://127.0.0.1:8093"
```

the url should be the operator 's `operator_server_ip_port_addr` config.

## Boot the aggregator

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
./bin/mach-aggregator --config ./config-files/aggregator-holesky.yaml --ecdsa-private-key $PRIVATE_KEY
```

The `PRIVATE_KEY` should be the committer for Mach avs.
