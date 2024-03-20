# Mach M2

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