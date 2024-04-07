// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";

contract MachServiceManagerUpgrader is Script {
    function run() external {
        address machServiceManager = vm.envAddress("SERVICE_MANAGER");
        address newImplAddress = vm.envAddress("NEW_IMPL_ADDRESS");
        address proxyAdmin = vm.envAddress("PROXY_ADMIN");
        ProxyAdmin machAVSProxyAdmin = ProxyAdmin(proxyAdmin);
        vm.startBroadcast();
        machAVSProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(machServiceManager))), address(newImplAddress)
        );
        vm.stopBroadcast();
    }
}
