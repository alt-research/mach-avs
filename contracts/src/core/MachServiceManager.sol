// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {MachServiceManagerStorage} from "./MachServiceManagerStorage.sol";

contract MachServiceManager is MachServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker, Pausable {
    using EnumerableSet for EnumerableSet.Bytes32Set;
    using EnumerableSet for EnumerableSet.AddressSet;

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

    //////////////////////////////////////////////////////////////////////////////
    //                              Admin Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Add an operator to the allowlist.
     * @param operator The operator to add
     */
    function addToAllowlist(address operator) external onlyOwner {
        require(operator != address(0), "MachServiceManager.addToAllowlist: zero address");
        require(!_allowlist[operator], "MachServiceManager.addToAllowlist: already in allowlist");
        _allowlist[operator] = true;
        emit OperatorAllowed(operator);
    }

    /**
     * @notice Remove an operator from the allowlist.
     * @param operator The operator to remove
     */
    function removeFromAllowlist(address operator) external onlyOwner {
        require(_allowlist[operator], "MachServiceManager.removeFromAllowlist: not in allowlist");
        _allowlist[operator] = false;
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

    function removeAlert(bytes32 messageHash) external onlyOwner {
        _messageHashes.remove(messageHash);
        emit AlertRemoved(messageHash, _msgSender());
    }

    //////////////////////////////////////////////////////////////////////////////
    //                          Operator Registration                           //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Register an operator with the AVS. Forwards call to EigenLayer' AVSDirectory.
     * @param operatorSignature The signature, salt, and expiry of the operator's signature.
     */
    function registerOperator(ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature)
        external
        whenNotPaused
        onlyRegistryCoordinator
    {
        address operator = msg.sender;
        require(!allowlistEnabled || _allowlist[operator], "MachServiceManager.registerOperator: not allowed");
        // todo check strategy and stake
        _operators.add(operator);
        _avsDirectory.registerOperatorToAVS(operator, operatorSignature);
        emit OperatorAdded(operator);
    }

    /**
     * @notice Deregister an operator from the AVS. Forwards a call to EigenLayer's AVSDirectory.
     */
    function deregisterOperator() external whenNotPaused onlyRegistryCoordinator {
        address operator = msg.sender;
        _operators.remove(operator);
        _avsDirectory.deregisterOperatorFromAVS(operator);
        emit OperatorRemoved(operator);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Alert Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    function confirmAlert(
        AlertHeader calldata alertHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external whenNotPaused onlyAlertConfirmer {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        require(
            tx.origin == msg.sender, "MachServiceManager.confirmAlert: header and nonsigner data must be in calldata"
        );
        // make sure the stakes against which the Batch is being confirmed are not stale
        require(
            alertHeader.referenceBlockNumber <= block.number,
            "MachServiceManager.confirmAlert: specified referenceBlockNumber is in future"
        );
        bytes32 hashedHeader = alertHeader.messageHash;

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
        _messageHashes.add(alertHeader.messageHash);

        emit AlertConfirmed(hashedHeader, alertHeader.messageHash);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                               View Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    function totalAlerts() public view returns (uint256) {
        return _messageHashes.length();
    }

    function contains(bytes32 messageHash) public view returns (bool) {
        return _messageHashes.contains(messageHash);
    }

    function queryMessageHashes(uint256 start, uint256 querySize) public view returns (bytes32[] memory) {
        uint256 length = totalAlerts();

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
        return ReducedAlertHeader({messageHash: alertHeader.messageHash});
    }
}
