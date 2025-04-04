// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "forge-std/Script.sol";
import "forge-std/console2.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";
import {IStakeRegistry, IDelegationManager} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IRewardsCoordinator} from "eigenlayer-core/contracts/interfaces/IRewardsCoordinator.sol";
import {IBLSSignatureChecker} from "eigenlayer-middleware/interfaces/IBLSSignatureChecker.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";
import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";
import {SocketRegistry} from "eigenlayer-middleware/SocketRegistry.sol";
import {IMachServiceManager} from "../src/interfaces/IMachServiceManager.sol";
import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";

contract MachServiceManagerImplDeployer is Script {
    function run() external {
        // Read contract addresses from JSON file for existing proxy addresses
        string memory proxyAddressesPath = vm.envOr("CONTRACT_ADDRESSES_PATH", string("./script/input/proxies.json"));
        string memory proxyAddressesJson = vm.readFile(proxyAddressesPath);

        // Read addresses from JSON
        address delegationManager = abi.decode(vm.parseJson(proxyAddressesJson, ".DELEGATION_MANAGER"), (address));
        address avsDirectory = abi.decode(vm.parseJson(proxyAddressesJson, ".AVS_DIRECTORY"), (address));
        address rewardsCoordinator = abi.decode(vm.parseJson(proxyAddressesJson, ".REWARDS_COORDINATOR"), (address));
        address registryCoordinatorProxy =
            abi.decode(vm.parseJson(proxyAddressesJson, ".REGISTRY_COORDINATOR"), (address));
        address apkRegistryProxy = abi.decode(vm.parseJson(proxyAddressesJson, ".APK_REGISTRY"), (address));
        address indexRegistryProxy = abi.decode(vm.parseJson(proxyAddressesJson, ".INDEX_REGISTRY"), (address));
        address stakeRegistryProxy = abi.decode(vm.parseJson(proxyAddressesJson, ".STAKE_REGISTRY"), (address));
        address machServiceManagerProxy =
            abi.decode(vm.parseJson(proxyAddressesJson, ".MACH_SERVICE_MANAGER"), (address));

        vm.startBroadcast();

        // Deploy apk registry implementation
        BLSApkRegistry apkRegistryImpl = new BLSApkRegistry(IRegistryCoordinator(registryCoordinatorProxy));

        // Deploy implementations for the component contracts
        IndexRegistry indexRegistryImpl = new IndexRegistry(IRegistryCoordinator(registryCoordinatorProxy));

        StakeRegistry stakeRegistryImpl =
            new StakeRegistry(IRegistryCoordinator(registryCoordinatorProxy), IDelegationManager(delegationManager));

        // Deploy socket registry which is needed for registry coordinator
        SocketRegistry socketRegistryContract = new SocketRegistry(RegistryCoordinator(registryCoordinatorProxy));

        // Deploy registry coordinator implementation
        RegistryCoordinator registryCoordinatorImpl = new RegistryCoordinator(
            IMachServiceManager(machServiceManagerProxy),
            IStakeRegistry(stakeRegistryProxy),
            IBLSApkRegistry(apkRegistryProxy),
            IIndexRegistry(indexRegistryProxy),
            socketRegistryContract
        );

        // Deploy BLS signature checker contract which is needed for mach service manager
        BLSSignatureChecker blsSignatureChecker =
            new BLSSignatureChecker(IRegistryCoordinator(registryCoordinatorProxy));

        // Deploy mach service manager implementation
        MachServiceManager machServiceManagerImpl = new MachServiceManager(
            IAVSDirectory(avsDirectory),
            IRewardsCoordinator(rewardsCoordinator),
            IRegistryCoordinator(registryCoordinatorProxy),
            IStakeRegistry(stakeRegistryProxy),
            blsSignatureChecker // Use the newly deployed BLSSignatureChecker
        );

        vm.stopBroadcast();

        // Log all deployed implementation addresses
        console2.log("====== Implementation Contracts Deployed ======");
        console2.log("Mach Service Manager Implementation:", address(machServiceManagerImpl));
        console2.log("Registry Coordinator Implementation:", address(registryCoordinatorImpl));
        console2.log("Index Registry Implementation:", address(indexRegistryImpl));
        console2.log("Stake Registry Implementation:", address(stakeRegistryImpl));
        console2.log("APK Registry Implementation:", address(apkRegistryImpl));
        console2.log("BLS Signature Checker:", address(blsSignatureChecker));
        console2.log("Socket Registry:", address(socketRegistryContract));
    }
}
