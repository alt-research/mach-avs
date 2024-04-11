// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "./AVSDeployer.sol";
import "../src/error/Errors.sol";

contract MachServiceManagerTest is AVSDeployer {
    event OperatorAllowed(address operator);
    event OperatorDisallowed(address operator);
    event AllowlistEnabled();
    event AllowlistDisabled();

    function setUp() public virtual {
        _deployMockEigenLayerAndAVS();
    }

    function test_AddToAllowlist() public {
        vm.startPrank(proxyAdminOwner);
        assertFalse(serviceManager.allowlist(defaultOperator), "mismatch");
        vm.expectEmit();
        emit OperatorAllowed(defaultOperator);
        serviceManager.addToAllowlist(defaultOperator);
        assertTrue(serviceManager.allowlist(defaultOperator), "mismatch");
        vm.stopPrank();
    }

    function test_AddToAllowlist_RevertIfZeroAddress() public {
        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(ZeroAddress.selector);
        serviceManager.addToAllowlist(address(0));
        vm.stopPrank();
    }

    function test_AddToAllowlist_RevertIfAlreadyInAllowlist() public {
        test_AddToAllowlist();
        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(AlreadyInAllowlist.selector);
        serviceManager.addToAllowlist(defaultOperator);
        vm.stopPrank();
    }

    function test_RemoveFromAllowlist() public {
        test_AddToAllowlist();
        vm.startPrank(proxyAdminOwner);
        assertTrue(serviceManager.allowlist(defaultOperator), "Operator should be in allowlist before removal");

        vm.expectEmit();
        emit OperatorDisallowed(defaultOperator);

        serviceManager.removeFromAllowlist(defaultOperator);

        assertFalse(serviceManager.allowlist(defaultOperator), "Operator should not be in allowlist after removal");
        vm.stopPrank();
    }

    function test_RemoveFromAllowlist_RevertIfNotInAllowlist() public {
        address nonListedOperator = address(0xdead);

        vm.startPrank(proxyAdminOwner);
        assertFalse(serviceManager.allowlist(nonListedOperator), "Operator should not be in allowlist");

        vm.expectRevert(NotInAllowlist.selector);
        serviceManager.removeFromAllowlist(nonListedOperator);
        vm.stopPrank();
    }

    function test_DisableAllowlist() public {
        vm.startPrank(proxyAdminOwner);
        assertTrue(serviceManager.allowlistEnabled(), "Allowlist should be enabled initially");

        vm.expectEmit(true, true, true, true); // Check all parameters of the event
        emit AllowlistDisabled();

        serviceManager.disableAllowlist();

        assertFalse(serviceManager.allowlistEnabled(), "Allowlist should be disabled after calling disableAllowlist");
        vm.stopPrank();
    }

    function test_DisableAllowlist_RevertIfAlreadyDisabled() public {
        vm.startPrank(proxyAdminOwner);

        // First, ensure the allowlist is disabled
        serviceManager.disableAllowlist();

        assertFalse(serviceManager.allowlistEnabled(), "Allowlist should already be disabled");
        vm.expectRevert(AlreadyDisabled.selector); // Expect the specific revert for trying to disable an already disabled allowlist
        serviceManager.disableAllowlist();
        vm.stopPrank();
    }

    function test_EnableAllowlist() public {
        vm.startPrank(proxyAdminOwner);

        // First, ensure the allowlist is disabled
        serviceManager.disableAllowlist();

        assertFalse(serviceManager.allowlistEnabled(), "Allowlist should be disabled initially");

        vm.expectEmit(true, true, true, true); // Check all parameters of the event
        emit AllowlistEnabled();

        serviceManager.enableAllowlist();

        assertTrue(serviceManager.allowlistEnabled(), "Allowlist should be enabled after calling enableAllowlist");
        vm.stopPrank();
    }

    function test_EnableAllowlist_RevertIfAlreadyEnabled() public {
        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(AlreadyEnabled.selector); // Expect the specific revert for trying to enable an already enabled allowlist
        serviceManager.enableAllowlist();
        vm.stopPrank();
    }
}
