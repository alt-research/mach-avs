// SPDX-License-Identifier: UNLICENSED
pragma solidity =0.8.12;

import "forge-std/Test.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import "../src/interfaces/IMachServiceManager.sol";
import "../src/core/MachServiceManagerRegistry.sol";
import "../src/error/Errors.sol";

contract MockServiceManager is ITotalAlerts {
    uint256 public alerts;

    function setAlerts(uint256 _alerts) public {
        alerts = _alerts;
    }

    function totalAlerts(uint256) external view override returns (uint256) {
        return alerts;
    }
}

contract MockServiceManagerLegacy is ITotalAlertsLegacy {
    uint256 public alerts;

    function setAlerts(uint256 _alerts) public {
        alerts = _alerts;
    }

    function totalAlerts() external view override returns (uint256) {
        return alerts;
    }
}

contract MachServiceManagerRegistryTest is Test {
    MachServiceManagerRegistry registry;
    MockServiceManager mockServiceManager;
    MockServiceManagerLegacy mockServiceManagerLegacy;

    event ServiceManagerRegistered(uint256 indexed rollupChainId, address serviceManager, address sender);
    event ServiceManagerDeregistered(uint256 indexed rollupChainId, address serviceManager, address sender);

    function setUp() public {
        // Deploy the registry
        registry = MachServiceManagerRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(new MachServiceManagerRegistry()), address(new ProxyAdmin()), ""
                )
            )
        );
        registry.initialize();

        // Deploy the mock service managers
        mockServiceManager = new MockServiceManager();
        mockServiceManagerLegacy = new MockServiceManagerLegacy();

        // Ensure registry is owned by this test contract for testing purposes
        registry.transferOwnership(address(this));
    }

    function test_RegisterServiceManager() public {
        uint256 rollupChainId = 1;

        vm.expectEmit();
        emit ServiceManagerRegistered(rollupChainId, address(mockServiceManager), address(this));

        registry.registerServiceManager(rollupChainId, address(mockServiceManager));

        assertEq(
            address(registry.serviceManagers(rollupChainId)), address(mockServiceManager), "Service manager mismatch"
        );
    }

    function test_RegisterServiceManager_RevertIfZeroAddress() public {
        uint256 rollupChainId = 1;

        vm.expectRevert(ZeroAddress.selector);
        registry.registerServiceManager(rollupChainId, address(0));
    }

    function test_RegisterServiceManager_RevertIfAlreadyAdded() public {
        uint256 rollupChainId = 1;

        registry.registerServiceManager(rollupChainId, address(mockServiceManager));

        vm.expectRevert(AlreadyAdded.selector);
        registry.registerServiceManager(rollupChainId, address(mockServiceManager));
    }

    function test_RegisterServiceManager_RevertIfNotOwner() public {
        uint256 rollupChainId = 1;

        // Transfer ownership to another address for this test
        registry.transferOwnership(address(0xdead));

        vm.expectRevert("Ownable: caller is not the owner");
        registry.registerServiceManager(rollupChainId, address(mockServiceManager));
    }

    function test_HasActiveAlerts() public {
        uint256 rollupChainId = 1;

        // Set alerts to a non-zero value
        mockServiceManager.setAlerts(1);

        registry.registerServiceManager(rollupChainId, address(mockServiceManager));

        bool hasAlerts = registry.hasActiveAlerts(rollupChainId);

        assertTrue(hasAlerts, "Expected to have active alerts");
    }

    function test_HasActiveAlerts_NoServiceManager() public {
        bool hasAlerts = registry.hasActiveAlerts(99); // no service manager registerd for the rollup ID 99
        assertFalse(hasAlerts, "Expected to have no active alerts");
    }

    function test_HasActiveAlerts_NoAlerts() public {
        uint256 rollupChainId = 1;

        // Set alerts to zero
        mockServiceManager.setAlerts(0);

        registry.registerServiceManager(rollupChainId, address(mockServiceManager));

        bool hasAlerts = registry.hasActiveAlerts(rollupChainId);

        assertFalse(hasAlerts, "Expected to have no active alerts");
    }

    function test_HasActiveAlerts_LegacyAlerts() public {
        uint256 rollupChainId = 1;

        // Set alerts to a non-zero value (simulating legacy alerts)
        mockServiceManagerLegacy.setAlerts(1);

        registry.registerServiceManager(rollupChainId, address(mockServiceManagerLegacy));

        bool hasAlerts = registry.hasActiveAlerts(rollupChainId);

        assertTrue(hasAlerts, "Expected to have active alerts in legacy mode");
    }

    function test_HasActiveAlerts_LegacyNoAlerts() public {
        uint256 rollupChainId = 1;

        // Set alerts to zero (simulating no legacy alerts)
        mockServiceManagerLegacy.setAlerts(0);

        registry.registerServiceManager(rollupChainId, address(mockServiceManagerLegacy));

        bool hasAlerts = registry.hasActiveAlerts(rollupChainId);

        assertFalse(hasAlerts, "Expected to have no active alerts in legacy mode");
    }

    function test_DeregisterServiceManager() public {
        uint256 rollupChainId = 1;

        registry.registerServiceManager(rollupChainId, address(mockServiceManager));

        vm.expectEmit();
        emit ServiceManagerDeregistered(rollupChainId, address(mockServiceManager), address(this));

        registry.deregisterServiceManager(rollupChainId, address(mockServiceManager));

        assertEq(address(registry.serviceManagers(rollupChainId)), address(0), "Service manager should be deregistered");
    }

    function test_DeregisterServiceManager_RevertIfNotAdded() public {
        uint256 rollupChainId = 1;

        vm.expectRevert(NotAdded.selector);
        registry.deregisterServiceManager(rollupChainId, address(mockServiceManager));
    }

    function test_DeregisterServiceManager_RevertIfNotOwner() public {
        uint256 rollupChainId = 1;

        // Register a service manager
        registry.registerServiceManager(rollupChainId, address(mockServiceManager));

        // Transfer ownership to another address for this test
        registry.transferOwnership(address(0xdead));

        vm.expectRevert("Ownable: caller is not the owner");
        registry.deregisterServiceManager(rollupChainId, address(mockServiceManager));
    }
}
