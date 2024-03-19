// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";
import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {MachServiceManagerStorage} from "./MachServiceManagerStorage.sol";

contract MachServiceManager is MachServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker, Pausable {
    using EnumerableSet for EnumerableSet.UintSet;

    uint8 internal constant PAUSED_CONFIRM_ALERT = 0;

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

    function initialize(
        IPauserRegistry _pauserRegistry,
        uint256 _initialPausedStatus,
        address _initialOwner,
        address _batchConfirmer
    ) public initializer {
        _initializePauser(_pauserRegistry, _initialPausedStatus);
        _transferOwnership(_initialOwner);
        _setAlertConfirmer(_batchConfirmer);
    }

    function removeAlert(uint256 blockNumber) external onlyOwner {
        _l2Blocks.remove(blockNumber);
        emit AlertRemoved(blockNumber, _msgSender());
    }

    function confirmAlert(
        AlertHeader calldata alertHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyWhenNotPaused(PAUSED_CONFIRM_ALERT) onlyAlertConfirmer {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        require(
            tx.origin == msg.sender, "MachServiceManager.confirmAlert: header and nonsigner data must be in calldata"
        );
        // make sure the stakes against which the Batch is being confirmed are not stale
        require(
            alertHeader.referenceBlockNumber <= block.number,
            "MachServiceManager.confirmAlert: specified referenceBlockNumber is in future"
        );
        bytes32 hashedHeader = hashAlertHeader(alertHeader);

        // check the signature
        (QuorumStakeTotals memory quorumStakeTotals, bytes32 signatoryRecordHash) = checkSignatures(
            hashedHeader,
            alertHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
            alertHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        // check that signatories own at least a threshold percentage of each quourm
        for (uint256 i = 0; i < alertHeader.quorumThresholdPercentages.length; i++) {
            // we don't check that the quorumThresholdPercentages are not >100 because a greater value would trivially fail the check, implying
            // signed stake > total stake
            // signedStakeForQuorum[i] / totalStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= quorumThresholdPercentages[i]
            // => signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= totalStakeForQuorum[i] * quorumThresholdPercentages[i]
            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR
                    >= quorumStakeTotals.totalStakeForQuorum[i] * uint8(alertHeader.quorumThresholdPercentages[i]),
                "MachServiceManager.confirmAlert: signatories do not own at least threshold percentage of a quorum"
            );
        }

        // store alert
        _l2Blocks.add(alertHeader.l2BlockNumber);

        emit AlertConfirmed(hashedHeader, alertHeader.l2BlockNumber);
    }

    function totalAlerts() public view returns (uint256) {
        return _l2Blocks.length();
    }

    function contains(uint256 blockNumber) public view returns (bool) {
        return _l2Blocks.contains(blockNumber);
    }

    function queryBlockNumber(uint256 start, uint256 querySize) public view returns (uint256[] memory) {
        uint256 length = totalAlerts();

        if (start >= length) {
            revert InvalidStartIndex();
        }

        uint256 end = start + querySize;

        if (end > length) {
            end = length;
        }

        uint256[] memory output = new uint256[](end - start);

        for (uint256 i = start; i < end; ++i) {
            output[i - start] = _l2Blocks.at(i);
        }

        return output;
    }

    /// @notice changes the alert confirmer
    function _setAlertConfirmer(address _alertConfirmer) internal {
        address previousBatchConfirmer = alertConfirmer;
        alertConfirmer = _alertConfirmer;
        emit AlertConfirmerChanged(previousBatchConfirmer, alertConfirmer);
    }

    /**
     * @notice converts a alert header to a reduced alert header
     * @param alertHeader the alert header to convert
     */
    function convertAlertHeaderToReducedAlertHeader(AlertHeader memory alertHeader)
        internal
        pure
        returns (ReducedAlertHeader memory)
    {
        return ReducedAlertHeader({
            l2BlockNumber: alertHeader.l2BlockNumber,
            referenceBlockNumber: alertHeader.referenceBlockNumber
        });
    }

    /// @notice hash the alert header
    function hashAlertHeader(AlertHeader memory alertHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(convertAlertHeaderToReducedAlertHeader(alertHeader)));
    }
}
