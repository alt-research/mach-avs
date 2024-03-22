// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "forge-std/Test.sol";
import {MockMachServiceManager} from "../src/mock/MockMachServiceManager.sol";
import {MockRegistryCoordinator} from "../src/mock/MockRegistryCoordinator.sol";
import {MockNonRegistryCoordinator} from "../src/mock/MockNonRegistryCoordinator.sol";

contract MachServiceManagerTest is Test {
    MockMachServiceManager public machServiceManager;
    MockRegistryCoordinator public mockRegistryCoordinator;
    MockNonRegistryCoordinator public mockNonRegistryCoordinator;

    function setUp() public {
        mockRegistryCoordinator = new MockRegistryCoordinator();
        mockNonRegistryCoordinator = new MockNonRegistryCoordinator();
        machServiceManager = new MockMachServiceManager(mockRegistryCoordinator);
        mockRegistryCoordinator.setMachServiceManager(address(machServiceManager));
        mockNonRegistryCoordinator.setMachServiceManager(address(machServiceManager));
    }

    function testRegisterOperatorToAVS() public {
        mockRegistryCoordinator.registerOperatorToAVS(msg.sender);
    }

    function testNonRegisterOperatorToAVS() public {
        vm.expectRevert();
        mockNonRegistryCoordinator.registerOperatorToAVS(msg.sender);
    }
}
