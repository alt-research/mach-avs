// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";


contract MachOptimismServiceManager {

     constructor(
        IAVSDirectory __avsDirectory,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry
    )
        BLSSignatureChecker(__registryCoordinator)
        ServiceManagerBase(__avsDirectory, __registryCoordinator, __stakeRegistry)
    {
        _disableInitializers();
    }
}