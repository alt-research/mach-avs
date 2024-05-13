// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "forge-std/Test.sol";

import "./BLSAVSDeployer.sol";
import "../src/error/Errors.sol";
import "../src/interfaces/IMachServiceManager.sol";

contract MachServiceManagerTest is BLSAVSDeployer {
    event OperatorAllowed(address operator);
    event OperatorDisallowed(address operator);
    event AllowlistEnabled();
    event AllowlistDisabled();
    event AlertConfirmerChanged(address previousAddress, address newAddress);
    event WhitelisterChanged(address previousAddress, address newAddress);
    event QuorumThresholdPercentageChanged(uint8 thresholdPercentages);
    event RollupChainIdUpdated(uint256 previousRollupChainId, uint256 newRollupChainId);
    event AlertConfirmed(bytes32 indexed alertHeaderHash, bytes32 messageHash);

    function setUp() public virtual {
        _deployMockEigenLayerAndAVS();

        msgHash = keccak256(
            abi.encode(IMachServiceManager.ReducedAlertHeader({messageHash: "foo", referenceBlockNumber: 201}))
        );

        _setAggregatePublicKeysAndSignature();
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

    function test_confirmAlert() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        uint256 nonRandomNumber = 111;
        uint256 numNonSigners = 1;
        uint256 quorumBitmap = 1;
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        (
            BLSSignatureChecker.QuorumStakeTotals memory quorumStakeTotals,
            /* bytes32 signatoryRecordHash */
        ) = serviceManager.checkSignatures(msgHash, quorumNumbers, referenceBlockNumber, nonSignerStakesAndSignature);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        vm.expectEmit();
        emit AlertConfirmed(msgHash, alertHeader.messageHash);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);

        vm.expectRevert(AlreadyAdded.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);

        vm.stopPrank();
    }
}
