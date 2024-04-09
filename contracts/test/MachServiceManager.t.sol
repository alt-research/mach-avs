// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "./AVSDeployer.sol";

contract MachServiceManagerTest is AVSDeployer {
    function setUp() public virtual {
        _deployMockEigenLayerAndAVS();
    }

    function test_AddToAllowlist() public {
        vm.startPrank(proxyAdminOwner);
        assertFalse(serviceManager.allowlist(defaultOperator), "mismatch");
        serviceManager.addToAllowlist(defaultOperator);
        assertTrue(serviceManager.allowlist(defaultOperator), "mismatch");
        vm.stopPrank();
    }
}
