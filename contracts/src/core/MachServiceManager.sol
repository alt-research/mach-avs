// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity ^0.8.12;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-core/contracts/interfaces/ISignatureUtils.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {MachServiceManagerStorage} from "./MachServiceManagerStorage.sol";
import {
    InvalidConfirmer,
    ZeroAddress,
    AlreadyInAllowlist,
    NotInAllowlist,
    InvalidReferenceBlockNum,
    InsufficientThreshold,
    InvalidStartIndex,
    InsufficientThresholdPercentages,
    InvalidSender,
    NotAllowed
} from "../error/Errors.sol";
import {IMachServiceManager} from "../interfaces/IMachServiceManager.sol";

/**
 * @title Primary entrypoint for procuring services from Altlayer Mach Service.
 * @author Altlayer, Inc.
 * @notice This contract is used for:
 * - whitelisting operators
 * - confirming the alert store by the aggregator with inferred aggregated signatures of the quorum
 */
contract MachServiceManager is
    IMachServiceManager,
    MachServiceManagerStorage,
    ServiceManagerBase,
    BLSSignatureChecker,
    Pausable
{
    using EnumerableSet for EnumerableSet.Bytes32Set;
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice when applied to a function, ensures that the function is only callable by the `alertConfirmer`.
    modifier onlyAlertConfirmer() {
        if (_msgSender() != alertConfirmer) {
            revert InvalidConfirmer();
        }
        _;
    }

    constructor(
        IAVSDirectory __avsDirectory,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry,
        uint256 __chainId
    )
        BLSSignatureChecker(__registryCoordinator)
        ServiceManagerBase(__avsDirectory, __registryCoordinator, __stakeRegistry)
        MachServiceManagerStorage(__chainId)
    {
        _disableInitializers();
    }

    function initialize(
        IPauserRegistry _pauserRegistry,
        uint256 _initialPausedStatus,
        address _initialOwner,
        address _alertConfirmer
    ) public initializer {
        _initializePauser(_pauserRegistry, _initialPausedStatus);
        _transferOwnership(_initialOwner);
        _setAlertConfirmer(_alertConfirmer);
        allowlistEnabled = true;
        quorumThresholdPercentage = 66;
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Admin Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Add an operator to the allowlist.
     * @param operator The operator to add
     */
    function addToAllowlist(address operator) external onlyOwner {
        if (operator == address(0)) {
            revert ZeroAddress();
        }
        if (allowlist[operator]) {
            revert AlreadyInAllowlist();
        }
        allowlist[operator] = true;
        emit OperatorAllowed(operator);
    }

    /**
     * @notice Remove an operator from the allowlist.
     * @param operator The operator to remove
     */
    function removeFromAllowlist(address operator) external onlyOwner {
        if (!allowlist[operator]) {
            revert NotInAllowlist();
        }
        allowlist[operator] = false;
        emit OperatorDisallowed(operator);
    }

    /**
     * @notice Enable the allowlist.
     */
    function enableAllowlist() external onlyOwner {
        allowlistEnabled = true;
        emit AllowlistEnabled();
    }

    /**
     * @notice Disable the allowlist.
     */
    function disableAllowlist() external onlyOwner {
        allowlistEnabled = false;
        emit AllowlistDisabled();
    }

    /**
     * @notice Remove an Alert.
     * @param messageHash The message hash of the alert
     */
    function removeAlert(bytes32 messageHash) external onlyOwner {
        _messageHashes.remove(messageHash);
        emit AlertRemoved(messageHash, _msgSender());
    }

    /**
     * @notice Update quorum threshold percentage
     * @param thresholdPercentage The new quorum threshold percentage
     */
    function updateQuorumThresholdPercentage(uint8 thresholdPercentage) external onlyOwner {
        quorumThresholdPercentage = thresholdPercentage;
        emit QuorumThresholdPercentageChanged(thresholdPercentage);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                          Operator Registration                           //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Register an operator with the AVS. Forwards call to EigenLayer' AVSDirectory.
     * @param operator The address of the operator to register.
     * @param operatorSignature The signature, salt, and expiry of the operator's signature.
     */
    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) public override(ServiceManagerBase, IServiceManager) whenNotPaused onlyRegistryCoordinator {
        if (allowlistEnabled && !allowlist[operator]) {
            revert NotInAllowlist();
        }
        _avsDirectory.registerOperatorToAVS(operator, operatorSignature);
        // we don't check if this operator has registered or not as AVSDirectory has such checking already
        _operators.add(operator);
        emit OperatorAdded(operator);
    }

    /**
     * @notice Deregister an operator from the AVS. Forwards a call to EigenLayer's AVSDirectory.
     * @param operator The address of the operator to register.
     */
    function deregisterOperatorFromAVS(address operator)
        public
        override(ServiceManagerBase, IServiceManager)
        whenNotPaused
        onlyRegistryCoordinator
    {
        _operators.remove(operator);
        _avsDirectory.deregisterOperatorFromAVS(operator);
        emit OperatorRemoved(operator);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Alert Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice This function is used for
     * - submitting alert,
     * - check that the aggregate signature is valid,
     * - and check whether quorum has been achieved or not.
     */
    function confirmAlert(
        AlertHeader calldata alertHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external whenNotPaused onlyAlertConfirmer {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        if (tx.origin != msg.sender) {
            revert InvalidSender();
        }
        // make sure the stakes against which the Batch is being confirmed are not stale
        if (alertHeader.referenceBlockNumber > block.number) {
            revert InvalidReferenceBlockNum();
        }
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
            uint8 currentQuorumThresholdPercentages = uint8(alertHeader.quorumThresholdPercentages[i]);
            if (currentQuorumThresholdPercentages < quorumThresholdPercentage) {
                revert InsufficientThresholdPercentages();
            }
            if (
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR
                    < quorumStakeTotals.totalStakeForQuorum[i] * currentQuorumThresholdPercentages
            ) {
                revert InsufficientThreshold();
            }
        }

        // store alert
        _messageHashes.add(alertHeader.messageHash);

        emit AlertConfirmed(hashedHeader, alertHeader.messageHash);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                               View Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /// @notice Returns the length of total alerts
    function totalAlerts() external view returns (uint256) {
        return _messageHashes.length();
    }

    /// @notice Checks if messageHash exists
    function contains(bytes32 messageHash) external view returns (bool) {
        return _messageHashes.contains(messageHash);
    }

    /// @notice Returns an array of messageHash
    function queryMessageHashes(uint256 start, uint256 querySize) external view returns (bytes32[] memory) {
        uint256 length = _messageHashes.length();

        if (start >= length) {
            revert InvalidStartIndex();
        }

        uint256 end = start + querySize;

        if (end > length) {
            end = length;
        }

        bytes32[] memory output = new bytes32[](end - start);
        for (uint256 i = start; i < end; ++i) {
            output[i - start] = _messageHashes.at(i);
        }

        return output;
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Internal Functions                          //
    //////////////////////////////////////////////////////////////////////////////

    /// @notice hash the alert header
    function hashAlertHeader(AlertHeader memory alertHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(convertAlertHeaderToReducedAlertHeader(alertHeader)));
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
            messageHash: alertHeader.messageHash,
            referenceBlockNumber: alertHeader.referenceBlockNumber
        });
    }
}
