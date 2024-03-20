# Operator

## Send token to operator

```bash
cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --value 50ether 0x957D781ab2Bc6D27Fde0a0b427ebF46ee1395661

cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44 'transfer(address, uint256) (bool)' 0x957D781ab2Bc6D27Fde0a0b427ebF46ee1395661 100000000000
```

## Reg to eigenlayer

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml rel 
```

## Add `strategyBaseTVLLimits` to `strategyWhitelister`

Add `strategyWhitelister`, use `strategyManager`

get:

```bash
cast call 0x0165878A594ca255338adfa4d48449f69242Eb8F 'strategyWhitelister() (address)'
0x0000000000000000000000000000000000000000
```

set by owner:

```bash
cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 0x0165878A594ca255338adfa4d48449f69242Eb8F 'setStrategyWhitelister(address)' 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

get will be owner:

```bash
cast call 0x0165878A594ca255338adfa4d48449f69242Eb8F 'strategyWhitelister() (address)'
0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

set the `strategyBaseTVLLimits` to `strategyWhitelister`

```bash
cast send  -f 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 0x0165878A594ca255338adfa4d48449f69242Eb8F 'addStrategiesToDepositWhitelist(address[], bool[])' '[0x4A679253410272dd5232B3Ff7cF5dbB88f295319]' '[false]' 
```

## Deposit tokens into a strategy

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml d --strategy-addr 0x4A679253410272dd5232B3Ff7cF5dbB88f295319 --amount 10000000
```

## Reg

```bash
./bin/mach-operator-cli --config ./config-files/operator.yaml r   
```

## Boot

```bash
./bin/mach-aggregator --config ./config-files/aggregator.yaml --ecdsa-private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --avs-deployment ./contracts/script/output/machavs_deploy_output.json
```

```bash
./bin/mach-operator-signer --config ./config-files/operator.yaml 
```
