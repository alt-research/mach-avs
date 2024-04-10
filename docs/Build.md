# Mach AVS Build

## Install dependencies

Install go v1.12.1, first install [gvm](https://github.com/moovweb/gvm):

```bash
sudo apt-get install bison

bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```

Install go:

```bash
gvm install go1.21.1 
gvm use go1.21.1 --default
```

To build contract, need install [foundry](https://github.com/foundry-rs/foundry):

```bash
curl -L https://foundry.paradigm.xyz | bash

foundryup
```

## Build Operator and Aggregator

```bash
git clone https://github.com/alt-research/mach-avs.git
cd mach-avs
make build
```

will got the following output:

- ./bin/mach-operator-signer operator node for commit bls sig to aggregator
- ./bin/mach-aggregator aggregator for collect bls sig and commit to layer1
- ./bin/mach-operator-cli a tool for reg and dereg from avs

## Build Contract

```bash
git clone https://github.com/alt-research/mach-avs.git
cd mach-avs/contracts
git submodule update --init --recursive 
forge build 
```

Deploy the avs contract to testnet can see [Script For Testing AVS Contracts](../scripts/README.md).
