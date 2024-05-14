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
    event AlertRemoved(bytes32 messageHash, address sender);

    uint256 nonRandomNumber = 111;
    uint256 numNonSigners = 1;
    uint256 quorumBitmap = 1;
    bytes quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

    function setUp() public virtual {
        _deployMockEigenLayerAndAVS();

        msgHash = keccak256(
            abi.encode(IMachServiceManager.ReducedAlertHeader({messageHash: "foo", referenceBlockNumber: 201}))
        );

        _setAggregatePublicKeysAndSignature();
    }

    function test_SetConfirmer() public {
        address newConfirmer = address(42);
        vm.startPrank(proxyAdminOwner);
        assertTrue(serviceManager.alertConfirmer() == proxyAdminOwner, "mismatch");
        vm.expectEmit();
        emit AlertConfirmerChanged(proxyAdminOwner, newConfirmer);
        serviceManager.setConfirmer(newConfirmer);
        assertTrue(serviceManager.alertConfirmer() == newConfirmer, "mismatch");
        vm.stopPrank();
    }

    function test_SetConfirmer_RevertIfNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.setConfirmer(address(42));
    }

    function test_SetWhitelister() public {
        address newWhitelister = address(42);
        vm.startPrank(proxyAdminOwner);
        assertTrue(serviceManager.whitelister() == proxyAdminOwner, "mismatch");
        vm.expectEmit();
        emit WhitelisterChanged(proxyAdminOwner, newWhitelister);
        serviceManager.setWhitelister(newWhitelister);
        assertTrue(serviceManager.whitelister() == newWhitelister, "mismatch");
        vm.stopPrank();
    }

    function test_SetWhitelister_RevertIfNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.setWhitelister(address(42));
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

    function test_AddToAllowlist_RevertIfNotWhitelister() public {
        vm.expectRevert(NotWhitelister.selector);
        serviceManager.addToAllowlist(defaultOperator);
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

    function test_RemoveFromAllowlist_RevertIfNotWhitelister() public {
        vm.expectRevert(NotWhitelister.selector);
        serviceManager.removeFromAllowlist(defaultOperator);
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

    function test_DisableAllowlist_RevertIfNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.disableAllowlist();
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

    function test_EnableAllowlist_RevertIfNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.enableAllowlist();
    }

    function test_EnableAllowlist_RevertIfAlreadyEnabled() public {
        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(AlreadyEnabled.selector); // Expect the specific revert for trying to enable an already enabled allowlist
        serviceManager.enableAllowlist();
        vm.stopPrank();
    }

    function test_UpdateQuorumThresholdPercentage() public {
        vm.startPrank(proxyAdminOwner);
        assertTrue(serviceManager.quorumThresholdPercentage() == 66, "mismatch");
        vm.expectEmit();
        emit QuorumThresholdPercentageChanged(76);
        serviceManager.updateQuorumThresholdPercentage(76);
        assertTrue(serviceManager.quorumThresholdPercentage() == 76, "mismatch");
        vm.stopPrank();
    }

    function test_UpdateQuorumThresholdPercentage_RevertIfNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.updateQuorumThresholdPercentage(76);
    }

    function test_UpdateQuorumThresholdPercentage_RevertIfInvalidQuorumThresholdPercentage() public {
        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(InvalidQuorumThresholdPercentage.selector);
        serviceManager.updateQuorumThresholdPercentage(101);
        vm.stopPrank();
    }

    function test_UpdateRollupChainId() public {
        vm.startPrank(proxyAdminOwner);
        assertTrue(serviceManager.rollupChainId() == 1, "mismatch");
        vm.expectEmit();
        emit RollupChainIdUpdated(1, 42);
        serviceManager.updateRollupChainId(42);
        assertTrue(serviceManager.rollupChainId() == 42, "mismatch");
        vm.stopPrank();
    }

    function test_UpdateRollupChainId_RevertIfNotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.updateRollupChainId(42);
    }

    function test_confirmAlert() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

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
        assertEq(serviceManager.totalAlerts(), 0);
        assertFalse(serviceManager.contains("foo"));

        emit AlertConfirmed(msgHash, alertHeader.messageHash);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);

        assertEq(serviceManager.totalAlerts(), 1);
        assertTrue(serviceManager.contains("foo"));

        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfInvalidConfirmer() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });
        vm.expectRevert(InvalidConfirmer.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
    }

    function test_confirmAlert_RevertIfAlreadyAdded() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.expectRevert(AlreadyAdded.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfInvalidQuorumParam() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(5);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(InvalidQuorumParam.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfResolvedAlert() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        serviceManager.removeAlert("foo");

        vm.expectRevert(ResolvedAlert.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfInvalidReferenceBlockNum() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: uint32(block.number)
        });

        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(InvalidReferenceBlockNum.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfInvalidQuorumThresholdPercentage() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(101));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(InvalidQuorumThresholdPercentage.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfInsufficientThresholdPercentages() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, numNonSigners, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(65));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(InsufficientThresholdPercentages.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_confirmAlert_RevertIfInsufficientThreshold() public {
        vm.startPrank(proxyAdminOwner);
        serviceManager.disableAllowlist();
        vm.stopPrank();

        (
            uint32 referenceBlockNumber,
            BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) = _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(nonRandomNumber, 8, quorumBitmap);

        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(67));

        IMachServiceManager.AlertHeader memory alertHeader = IMachServiceManager.AlertHeader({
            messageHash: "foo",
            quorumNumbers: quorumNumbers,
            quorumThresholdPercentages: quorumThresholdPercentages,
            referenceBlockNumber: referenceBlockNumber
        });

        vm.startPrank(proxyAdminOwner);
        vm.expectRevert(InsufficientThreshold.selector);
        serviceManager.confirmAlert(alertHeader, nonSignerStakesAndSignature);
        vm.stopPrank();
    }

    function test_removeAlert() public {
        test_confirmAlert();
        vm.startPrank(proxyAdminOwner);
        assertEq(serviceManager.totalAlerts(), 1);
        assertTrue(serviceManager.contains("foo"));

        vm.expectEmit();
        emit AlertRemoved("foo", msg.sender);
        serviceManager.removeAlert("foo");

        assertFalse(serviceManager.contains("foo"));
        assertEq(serviceManager.totalAlerts(), 0);
        vm.stopPrank();
    }

    function test_removeAlert_RevertIfNotOwner() public {
        test_confirmAlert();
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManager.removeAlert("foo");
    }
}
