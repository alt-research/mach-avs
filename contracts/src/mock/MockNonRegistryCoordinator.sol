// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {MockMachServiceManager} from "./MockMachServiceManager.sol";

contract MockNonRegistryCoordinator is OwnableUpgradeable {
    MockMachServiceManager public mockMachServiceManager;

    function setMachServiceManager(address machServiceManager) public {
        mockMachServiceManager = MockMachServiceManager(machServiceManager);
    }

    function registerOperatorToAVS(address operator) public {
        mockMachServiceManager.registerOperatorToAVS(operator);
    }
}
