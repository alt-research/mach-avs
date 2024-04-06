// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";

contract MachServiceManagerImplDeployer is Script {
    function run() external {
        address avsDirectory = vm.envAddress("AVS_DIRECTORY");
        address registryCoordinator = vm.envAddress("REGISTRY_COORDINATOR");
        address stakeRegistry = vm.envAddress("STAKE_REGISTRY");
        uint256 chainId = vm.envUint("CHAIN_ID");
        vm.startBroadcast();
        // 1. deploy new implementation contract
        MachServiceManager machServiceManagerImplementation = new MachServiceManager(
            IAVSDirectory(avsDirectory),
            IRegistryCoordinator(registryCoordinator),
            IStakeRegistry(stakeRegistry),
            chainId
        );
        vm.stopBroadcast();
    }
}
