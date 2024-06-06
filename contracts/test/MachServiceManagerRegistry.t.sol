// SPDX-License-Identifier: UNLICENSED
pragma solidity =0.8.12;

import "forge-std/Test.sol";
import "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import "../src/interfaces/IMachServiceManager.sol";
import "../src/core/MachServiceManagerRegistry.sol";
import "../src/error/Errors.sol";

contract MachServiceManagerRegistryTest is Test {
    MachServiceManagerRegistry registry;
    IMachServiceManager mockServiceManager;

    event ServiceManagerRegistered(uint256 rollupChainId_, IMachServiceManager serviceManager_, address sender);

    function setUp() public {
        // Deploy the registry
        registry = new MachServiceManagerRegistry();
        registry.initialize();

        // Mocking the IMachServiceManager interface
        mockServiceManager = IMachServiceManager(address(1));

        // Ensure registry is owned by this test contract for testing purposes
        registry.transferOwnership(address(this));
    }

    function test_RegisterServiceManager() public {
        uint256 rollupChainId = 1;

        vm.expectEmit(true, true, true, true);
        emit ServiceManagerRegistered(rollupChainId, mockServiceManager, address(this));

        registry.registerServiceManager(rollupChainId, mockServiceManager);

        assertEq(
            address(registry.serviceManagers(rollupChainId)), address(mockServiceManager), "Service manager mismatch"
        );
    }

    function test_RegisterServiceManager_RevertIfZeroAddress() public {
        uint256 rollupChainId = 1;

        vm.expectRevert(ZeroAddress.selector);
        registry.registerServiceManager(rollupChainId, IMachServiceManager(address(0)));
    }

    function test_RegisterServiceManager_RevertIfAlreadyAdded() public {
        uint256 rollupChainId = 1;

        registry.registerServiceManager(rollupChainId, mockServiceManager);

        vm.expectRevert(AlreadyAdded.selector);
        registry.registerServiceManager(rollupChainId, mockServiceManager);
    }

    function test_RegisterServiceManager_RevertIfNotOwner() public {
        uint256 rollupChainId = 1;

        // Transfer ownership to another address for this test
        registry.transferOwnership(address(0xdead));

        vm.expectRevert("Ownable: caller is not the owner");
        registry.registerServiceManager(rollupChainId, mockServiceManager);
    }

    function test_HasActiveAlerts() public {
        uint256 rollupChainId = 1;

        // Mocking the behavior of totalAlerts to return a non-zero value
        vm.mockCall(
            address(mockServiceManager),
            abi.encodeWithSelector(IMachServiceManager.totalAlerts.selector, rollupChainId),
            abi.encode(1)
        );

        registry.registerServiceManager(rollupChainId, mockServiceManager);

        bool hasAlerts = registry.hasActiveAlerts(rollupChainId);

        assertTrue(hasAlerts, "Expected to have active alerts");
    }

    function test_HasActiveAlerts_NoAlerts() public {
        uint256 rollupChainId = 1;

        // Mocking the behavior of totalAlerts to return zero
        vm.mockCall(
            address(mockServiceManager),
            abi.encodeWithSelector(IMachServiceManager.totalAlerts.selector, rollupChainId),
            abi.encode(0)
        );

        registry.registerServiceManager(rollupChainId, mockServiceManager);

        bool hasAlerts = registry.hasActiveAlerts(rollupChainId);

        assertFalse(hasAlerts, "Expected to have no active alerts");
    }
}
