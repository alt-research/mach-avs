# Script For Testing AVS Contracts

## Note
Private keys in this repository are for testing purposes and should not be reused for production environment.

## 1. Launch a Testing Environemt

Launch anvil:

```bash
anvil
```

Use this following env:

```bash
export OWNER_ADDR=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
export OWNER_PRIVATE=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
export RPC_URL=http://localhost:8545
```

Deploy contracts

```bash
cd contracts

forge script script/EigenLayerDeployer.s.sol --broadcast -vvvv \
    --private-key $OWNER_PRIVATE \
    --rpc-url $RPC_URL 

forge script script/MachServiceManagerDeployer.s.sol \
    --private-key $OWNER_PRIVATE \
    --broadcast -vvvv --rpc-url $RPC_URL
```

The address is stored in `./contracts/script/output`.

## 2. Register Operators

In this example, we will register 3 operators to Mach AVS

```bash
cd ..
bash ./scripts/register_operator.sh $OWNER_ADDR $OWNER_PRIVATE $METADATA_URI 0xE9A7669aC9eBE9b7E21E0A323FC3A6f34CE744eb test1
bash ./scripts/register_operator.sh $OWNER_ADDR $OWNER_PRIVATE $METADATA_URI 0x957D781ab2Bc6D27Fde0a0b427ebF46ee1395661 test2
bash ./scripts/register_operator.sh $OWNER_ADDR $OWNER_PRIVATE $METADATA_URI 0x91d45D72e36c5a6838f14f49D607e9b16eD33f58 test3
```

It will use the keys in `./config-files/key/`, we can create new key by:

```bash
eigenlayer operator keys create -i -k ecdsa testName
eigenlayer operator keys create -i -k bls testName
```

the key is in `~/.eigenlayer/operator_keys/`

Metadata example:
```json
{
  "name": "Some operator",
  "website": "https://www.example.com",
  "description": "I operate on some data",
  "logo": "https://www.example.com/logo.png",
  "twitter": "https://x.com/example"
}
```

## 3. Launch BLS Signature Aggregator

```bash
./bin/mach-aggregator --config ./config-files/aggregator.yaml --ecdsa-private-key $OWNER_PRIVATE --avs-deployment ./contracts/script/output/machavs_deploy_output.json
```

## 4. Boot operators

We just use the test keys:

```bash
./bin/mach-operator-signer --config ./config-files/tests/operator_test1.yaml 
./bin/mach-operator-signer --config ./config-files/tests/operator_test2.yaml 
./bin/mach-operator-signer --config ./config-files/tests/operator_test3.yaml 
```

## 5. Create test alert (for testing purposes)

```bash
curl --noproxy '*' -H "Content-Type: application/json" \
  -d '{"id":2, "jsonrpc":"2.0", "method": "alert_blockMismatch", "params":{"invalid_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "expect_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "l2_block_number": 2000}}' \
  http://localhost:8091

  curl --noproxy '*' -H "Content-Type: application/json" \
  -d '{"id":2, "jsonrpc":"2.0", "method": "alert_blockMismatch", "params":{"invalid_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "expect_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "l2_block_number": 2000}}' \
  http://localhost:8092

  curl --noproxy '*' -H "Content-Type: application/json" \
  -d '{"id":2, "jsonrpc":"2.0", "method": "alert_blockMismatch", "params":{"invalid_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "expect_output_root": "5FC8d32690cc91D4c39d9d3abcBD16989F875700000000000000000000000000", "l2_block_number": 2000}}' \
  http://localhost:8093
```

it will create alert for all the 3 operators. In the aggregator logs, we should be observing that signatures were received and it crossed the threshold needed.

```log
{"level":"info","ts":1711197400.2389536,"caller":"logging/zap_logger.go:49","msg":"Received response from blsAggregationService","blsAggServiceResp":{"Err":null,"TaskIndex":0,"TaskResponseDigest":[112,9,108,53,119,234,136,155,76,207,40,99,52,110,43,23,26,111,84,95,174,145,200,118,96,230,38,118,167,197,125,189],"NonSignersPubkeysG1":[],"QuorumApksG1":[{"X":"9245057968145484860804280814781348845784489795866678144111989947626410422422","Y":"4206160717496399935935556607824205312866665547119048439416182281445762835632"}],"SignersApkG2":{"X":{"A0":"5319687821092315421704966764687566991558776008265084532749523747246034678820","A1":"11223072676248240401971043350937406910510747595318510439500982143518357330882"},"Y":{"A0":"4763582069887216360727645755241906558814452159592383203334545463609598755585","A1":"21350356272487941664524832102750065802080684976181056332566094158938378720768"}},"SignersAggSigG1":{"g1_point":{"X":"16428293862790245810342020027165842618240383426525781313336816979755469179242","Y":"11576373886324819258408043941549120009185745867634683580045521586375033792509"}},"NonSignerQuorumBitmapIndices":[],"QuorumApkIndices":[3],"TotalStakeIndices":[3],"NonSignerStakeIndices":[[]]}}
{"level":"info","ts":1711197400.239083,"caller":"logging/zap_logger.go:49","msg":"Threshold reached. Sending aggregated response onchain.","taskIndex":0}
```
