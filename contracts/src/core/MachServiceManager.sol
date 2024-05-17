// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity =0.8.12;

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
    NotWhitelister,
    ZeroAddress,
    AlreadyInAllowlist,
    NotInAllowlist,
    NoStatusChange,
    InvalidRollupChainID,
    InvalidReferenceBlockNum,
    InsufficientThreshold,
    InvalidStartIndex,
    InsufficientThresholdPercentages,
    InvalidSender,
    InvalidQuorumParam,
    InvalidQuorumThresholdPercentage,
    AlreadyAdded,
    ResolvedAlert,
    AlreadyEnabled,
    AlreadyDisabled
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

    /**
     * @dev Ensures that the function is only callable by the `alertConfirmer`.
     */
    modifier onlyAlertConfirmer() {
        if (_msgSender() != alertConfirmer) {
            revert InvalidConfirmer();
        }
        _;
    }

    /**
     * @dev Ensures that the function is only callable by the `whitelister`.
     */
    modifier onlyWhitelister() {
        if (_msgSender() != whitelister) {
            revert NotWhitelister();
        }
        _;
    }

    /**
     * @dev Ensures that the `rollupChainID` is valid.
     */
    modifier onlyValidRollupChainID(uint256 rollupChainID) {
        if (!rollupChainIDs[rollupChainID]) {
            revert InvalidRollupChainID();
        }
        _;
    }

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
        IPauserRegistry pauserRegistry_,
        uint256 initialPausedStatus_,
        address initialOwner_,
        address alertConfirmer_,
        address whitelister_,
        uint256[] calldata rollupChainIDs_
    ) public initializer {
        _initializePauser(pauserRegistry_, initialPausedStatus_);
        __ServiceManagerBase_init(initialOwner_);
        _setAlertConfirmer(alertConfirmer_);
        _setWhitelister(whitelister_);

        for (uint256 i; i < rollupChainIDs_.length; ++i) {
            _setRollupChainID(rollupChainIDs_[i], true);
        }

        allowlistEnabled = true;
        quorumThresholdPercentage = 66;
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Admin Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @inheritdoc IMachServiceManager
     */
    function addToAllowlist(address operator) external onlyWhitelister {
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
     * @inheritdoc IMachServiceManager
     */
    function removeFromAllowlist(address operator) external onlyWhitelister {
        if (!allowlist[operator]) {
            revert NotInAllowlist();
        }
        allowlist[operator] = false;
        emit OperatorDisallowed(operator);
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function enableAllowlist() external onlyOwner {
        if (allowlistEnabled) {
            revert AlreadyEnabled();
        } else {
            allowlistEnabled = true;
            emit AllowlistEnabled();
        }
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function disableAllowlist() external onlyOwner {
        if (!allowlistEnabled) {
            revert AlreadyDisabled();
        } else {
            allowlistEnabled = false;
            emit AllowlistDisabled();
        }
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function setConfirmer(address confirmer) external onlyOwner {
        _setAlertConfirmer(confirmer);
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function setWhitelister(address whitelister) external onlyOwner {
        _setWhitelister(whitelister);
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function setRollupChainID(uint256 rollupChainId, bool status) external onlyOwner {
        _setRollupChainID(rollupChainId, status);
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function removeAlert(uint256 rollupChainId, bytes32 messageHash)
        external
        onlyValidRollupChainID(rollupChainId)
        onlyOwner
    {
        bool ret = _messageHashes[rollupChainId].remove(messageHash);
        if (ret) {
            _resolvedMessageHashes[rollupChainId].add(messageHash);
            emit AlertRemoved(messageHash, _msgSender());
        }
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function updateQuorumThresholdPercentage(uint8 thresholdPercentage) external onlyOwner {
        if (thresholdPercentage > 100) {
            revert InvalidQuorumThresholdPercentage();
        }
        quorumThresholdPercentage = thresholdPercentage;
        emit QuorumThresholdPercentageChanged(thresholdPercentage);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                          Operator Registration                           //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @inheritdoc IServiceManager
     */
    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) public override(ServiceManagerBase, IServiceManager) whenNotPaused onlyRegistryCoordinator {
        if (allowlistEnabled && !allowlist[operator]) {
            revert NotInAllowlist();
        }
        // we don't check if this operator has registered or not as AVSDirectory has such checking already
        _operators.add(operator);
        // Stake requirement for quorum is checked in StakeRegistry.sol
        // https://github.com/Layr-Labs/eigenlayer-middleware/blob/dev/src/RegistryCoordinator.sol#L488
        // https://github.com/Layr-Labs/eigenlayer-middleware/blob/dev/src/StakeRegistry.sol#L84
        _avsDirectory.registerOperatorToAVS(operator, operatorSignature);
        emit OperatorAdded(operator);
    }

    /**
     * @inheritdoc IServiceManager
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
     * @inheritdoc IMachServiceManager
     */
    function confirmAlert(
        uint256 rollupChainId,
        AlertHeader calldata alertHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external whenNotPaused onlyAlertConfirmer onlyValidRollupChainID(rollupChainId) {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        if (tx.origin != msg.sender) {
            revert InvalidSender();
        }

        // check is it is the resolved alert before
        if (_resolvedMessageHashes[rollupChainId].contains(alertHeader.messageHash)) {
            revert ResolvedAlert();
        }

        // make sure the stakes against which the Batch is being confirmed are not stale
        if (alertHeader.referenceBlockNumber >= block.number) {
            revert InvalidReferenceBlockNum();
        }
        bytes32 hashedHeader = _hashAlertHeader(alertHeader);

        // check quorum parameters
        if (alertHeader.quorumNumbers.length != alertHeader.quorumThresholdPercentages.length) {
            revert InvalidQuorumParam();
        }

        // check the signature
        (QuorumStakeTotals memory quorumStakeTotals, bytes32 signatoryRecordHash) = checkSignatures(
            hashedHeader,
            alertHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
            alertHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        // check that signatories own at least a threshold percentage of each quourm
        for (uint256 i = 0; i < alertHeader.quorumThresholdPercentages.length; i++) {
            // signed stake > total stake
            // signedStakeForQuorum[i] / totalStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= quorumThresholdPercentages[i]
            // => signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= totalStakeForQuorum[i] * quorumThresholdPercentages[i]
            uint8 currentQuorumThresholdPercentages = uint8(alertHeader.quorumThresholdPercentages[i]);
            if (currentQuorumThresholdPercentages > 100) {
                revert InvalidQuorumThresholdPercentage();
            }
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
        bool success = _messageHashes[rollupChainId].add(alertHeader.messageHash);
        if (!success) {
            revert AlreadyAdded();
        }

        emit AlertConfirmed(hashedHeader, alertHeader.messageHash);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                               View Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @inheritdoc IMachServiceManager
     */
    function totalAlerts(uint256 rollupChainId) external view returns (uint256) {
        return _messageHashes[rollupChainId].length();
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function contains(uint256 rollupChainId, bytes32 messageHash) external view returns (bool) {
        return _messageHashes[rollupChainId].contains(messageHash);
    }

    /**
     * @inheritdoc IMachServiceManager
     */
    function queryMessageHashes(uint256 rollupChainId, uint256 start, uint256 querySize)
        external
        view
        returns (bytes32[] memory)
    {
        uint256 length = _messageHashes[rollupChainId].length();

        if (start >= length) {
            revert InvalidStartIndex();
        }

        uint256 end = start + querySize;

        if (end > length) {
            end = length;
        }

        bytes32[] memory output = new bytes32[](end - start);
        for (uint256 i = start; i < end; ++i) {
            output[i - start] = _messageHashes[rollupChainId].at(i);
        }

        return output;
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Internal Functions                          //
    //////////////////////////////////////////////////////////////////////////////

    /**
     *  @dev Hashes an alert header
     */
    function _hashAlertHeader(AlertHeader calldata alertHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(_convertAlertHeaderToReducedAlertHeader(alertHeader)));
    }

    /**
     * @dev Changes the alert confirmer
     */
    function _setAlertConfirmer(address _alertConfirmer) internal {
        address previousBatchConfirmer = alertConfirmer;
        alertConfirmer = _alertConfirmer;
        emit AlertConfirmerChanged(previousBatchConfirmer, alertConfirmer);
    }

    /**
     *  @dev Changes the whitelister
     */
    function _setWhitelister(address _whitelister) internal {
        address previousWhitelister = whitelister;
        whitelister = _whitelister;
        emit WhitelisterChanged(previousWhitelister, _whitelister);
    }

    /**
     * @dev Converts a alert header to a reduced alert header
     * @param alertHeader the alert header to convert
     */
    function _convertAlertHeaderToReducedAlertHeader(AlertHeader calldata alertHeader)
        internal
        pure
        returns (ReducedAlertHeader memory)
    {
        return ReducedAlertHeader({
            messageHash: alertHeader.messageHash,
            referenceBlockNumber: alertHeader.referenceBlockNumber
        });
    }

    function _setRollupChainID(uint256 rollupChainId, bool status) internal {
        if (rollupChainId < 1) {
            revert InvalidRollupChainID();
        }
        if (rollupChainIDs[rollupChainId] == status) {
            revert NoStatusChange();
        }
        rollupChainIDs[rollupChainId] = status;
        emit RollupChainIDUpdated(rollupChainId, status);
    }
}
