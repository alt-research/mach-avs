# Mach AVS (M2)

Mach AVS is a fast finality layer for Ethereum rollups. In Mach AVS, operators 

1. verifier nodes for rollups
2. monitor the rollup for fraudulent behavior
3. submit a alert to the Mach AVS

# Use cases

1. Fast finality services for Web3 application
2. Act as a circuit breaker for RPC nodes

# Mach AVS contract

## Dependencies 

Mach AVS uses [EigenLayer Middleware v0.1.2](https://github.com/Layr-Labs/eigenlayer-middleware/releases/tag/v0.1.2-holesky-init-deployment)

## Training wheels

Mach AVS includes operator allowlist which can be managed by contract owner. 

## Enable/Disable operator allowlist
- Enable operator allowlist: `enableAllowlist()`
- Diosable operator allowlist: `disableAllowlist()`

## Operator allowlist management 
- Add operator to whitelist: `addToAllowlist(address operator)`
- Remove operator from whitelist: `removeFromAllowlist(address operator)` 

## Alert submission

- Submit alert from aggregator (Alert confirmer): `confirmAlert()`

# Deployment Process

## Deploy EigenLayer Contract

```shell
forge script script/EigenLayerDeployer.s.sol --rpc-url $rpc --private-key $key --broadcast -vvvv
```

## Create BLS key

Install `eigenlayer-cli`:

```bash
curl -sSfL https://raw.githubusercontent.com/layr-labs/eigenlayer-cli/master/scripts/install.sh | sh -s
```

Create bls key:

```bash
eigenlayer operator keys create --key-type bls test1
```