// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import {MachServiceManagerRegistry} from "../src/core/MachServiceManagerRegistry.sol";

// PROXY_ADMIN=$PROXY_ADMIN forge script ./script/MachServiceManagerRegistryDeployer.s.sol \
//     --private-key $PK \
//     --rpc-url $URL \
//     --etherscan-api-key $API_KEY \
//     --broadcast -vvvv --slow --verify
contract MachServiceManagerRegistryDeployer is Script {
    function run() external {
        vm.startBroadcast();
        address proxyAdmin = vm.envAddress("PROXY_ADMIN");
        MachServiceManagerRegistry registry = MachServiceManagerRegistry(
            address(new TransparentUpgradeableProxy(address(new MachServiceManagerRegistry()), address(proxyAdmin), ""))
        );
        registry.initialize();
        vm.stopBroadcast();
    }
}
