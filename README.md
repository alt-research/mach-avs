# Mach AVS (M2)

AltLayer Mach AVS is a fast finality layer for Ethereum rollups. In Mach AVS , operators will run verifier nodes for rollups, monitor the rollup for fraudulent behavior, and submit a fraudulent alert to Mach AVS. Fraudulent alert can then be confirmed via sufficient quorum or ZK proofs.

# Use cases

1. Fast finality services for Web3 application
2. Circuit breaker for RPC nodes

# Components 

Mach AVS consists of the following componenet:
- Mach AVS service manager contract
- Mach AVS aggregator (for signature aggregation version)
- Mach AVS operator
- Mach verifier

## Architecture

### BLS Signature Aggregation Mode 
![BLS Mode](docs/images/EigenlayerMachAVSArch(BLS).jpg)

### ZK Proof Mode
![ZK Proof Mode](docs/images/EigenlayerMachAVSArch(ZK-OP).jpg)

## Mach AVS service manager contract

Mach AVS service manager contracts can be found in [contracts folder](contracts/src/core/)
- [Mach AVS for OP stack based rollup (ZK proof)](contracts/src/core/MachOptimismServiceManager.sol)
- [Mach AVS for all rollup stack (Signature aggregation)](contracts/src/core/MachServiceManager.sol)

### Dependencies 

Mach AVS uses [EigenLayer Middleware v0.1.2](https://github.com/Layr-Labs/eigenlayer-middleware/releases/tag/v0.1.2-holesky-init-deployment)

### Alert submission

- Submit alert: `confirmAlert()`

### Training wheels

Mach AVS includes operator allowlist which can be managed by contract owner. 

#### Enable/Disable operator allowlist
- Enable operator allowlist: `enableAllowlist()`
- Diosable operator allowlist: `disableAllowlist()`

#### Operator allowlist management 
- Add operator to whitelist: `addToAllowlist(address operator)`
- Remove operator from whitelist: `removeFromAllowlist(address operator)` 

## Mach AVS Aggregator (for Signature Aggregation Mode)

Mach AVS aggregator service can be found in [aggregator](aggregator/)

## Mach AVS Operator

Operator sample configuration file can be found at [config-files/operator.yaml](config-files/operator.yaml).

Operator can be run using the following:
```bash
./bin/mach-operator-signer --config <PATH_TO_CONFIG> 
```
Node operator client can be found in [operator](operator/). For more information on how to run operator, check out our guide at [here](contracts/script/README.md)

## Mach AVS Verifier

The verifier codebase is found at [https://github.com/alt-research/alt-mach-verifier](https://github.com/alt-research/alt-mach-verifier)

# Deployment script

Check out [scripts](contracts/script)
