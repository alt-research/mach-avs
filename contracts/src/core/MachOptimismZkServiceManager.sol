// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity ^0.8.12;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {ISlasher} from "eigenlayer-core/contracts/interfaces/ISlasher.sol";
import {ISignatureUtils} from "eigenlayer-core/contracts/interfaces/ISignatureUtils.sol";
import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IDelegationManager} from "eigenlayer-core/contracts/interfaces/IDelegationManager.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {ServiceManagerBase, IRegistryCoordinator, IStakeRegistry} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {MachOptimismZkServiceManagerStorage} from "./MachOptimismZkServiceManagerStorage.sol";
import {IMachOptimism, CallbackAuthorization, IRiscZeroVerifier} from "../interfaces/IMachOptimism.sol";
import {IMachOptimismL2OutputOracle} from "../interfaces/IMachOptimismL2OutputOracle.sol";
import "../error/Errors.sol";

contract MachOptimismZkServiceManager is
    IMachOptimism,
    MachOptimismZkServiceManagerStorage,
    ServiceManagerBase,
    Pausable
{
    using EnumerableSet for EnumerableSet.AddressSet;

    constructor(
        uint256 rollupChainID_,
        IAVSDirectory __avsDirectory,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry
    )
        ServiceManagerBase(__avsDirectory, __registryCoordinator, __stakeRegistry)
        MachOptimismZkServiceManagerStorage(block.chainid, rollupChainID_)
    {}

    modifier onlyValidOperator() {
        IBLSApkRegistry blsApkRegistry = _registryCoordinator.blsApkRegistry();
        bytes32 operatorId = blsApkRegistry.getOperatorId(msg.sender);
        if (operatorId == bytes32(0)) {
            revert NotOperator();
        }
        _;
    }

    /// @notice Initializes the contract with provided parameters.
    function initialize(
        address contractOwner,
        bytes32 imageId_,
        IMachOptimismL2OutputOracle l2OutputOracle_,
        IRiscZeroVerifier verifier_
    ) external initializer {
        __Ownable_init();
        _transferOwnership(contractOwner);
        if (address(l2OutputOracle_) == address(0)) {
            revert ZeroAddress();
        }
        if (address(verifier_) == address(0)) {
            revert ZeroAddress();
        }
        if (imageId_ == bytes32(0)) {
            revert ZeroValue();
        }
        l2OutputOracle = l2OutputOracle_;
        verifier = verifier_;
        imageId = imageId_;
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Admin Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /// @notice Set imageId.
    function setImageId(bytes32 imageId_) external onlyOwner {
        imageId = imageId_;
    }

    /// @notice Set verifier.
    function setRiscZeroVerifier(IRiscZeroVerifier verifier_) external onlyOwner {
        if (address(verifier_) == address(0)) {
            revert ZeroAddress();
        }

        verifier = verifier_;
    }

    /// @notice Clear alerts.
    function clearAlerts() external onlyOwner {
        delete l2OutputAlerts;
        provedIndex = 0;
    }

    /// @notice Clear block alerts up to a specific number.
    function clearBlockAlertsUpTo(uint256 l2BlockNumber) external onlyOwner {
        require(l2BlockNumber > 0, "Invalid l2BlockNumber");

        uint256 alertsLength = l2OutputAlerts.length;

        if (alertsLength == 0 || provedIndex == 0) {
            revert("No alerts to clear");
        }

        if (provedIndex > alertsLength) {
            revert("Invalid provedIndex");
        }

        // Iterate through the alerts and clear those up to l2BlockNumber
        for (uint256 i = provedIndex - 1; i < alertsLength; i++) {
            if (l2OutputAlerts[i].l2BlockNumber <= l2BlockNumber) {
                // Clear the alert by shifting the subsequent alerts
                for (uint256 j = i; j < alertsLength - 1; j++) {
                    l2OutputAlerts[j] = l2OutputAlerts[j + 1];
                }
                l2OutputAlerts.pop();
            } else {
                break; // Stop once we reach an alert with a higher l2BlockNumber
            }
        }

        provedIndex = l2OutputAlerts.length;
    }

    /**
     * @notice Add an operator to the allowlist.
     * @param operator The operator to add
     */
    function addToAllowlist(address operator) external onlyOwner {
        if (operator == address(0)) {
            revert ZeroAddress();
        }
        if (_allowlist[operator]) {
            revert AlreadyInAllowlist();
        }
        _allowlist[operator] = true;
        emit OperatorAllowed(operator);
    }

    /**
     * @notice Remove an operator from the allowlist.
     * @param operator The operator to remove
     */
    function removeFromAllowlist(address operator) external onlyOwner {
        if (!_allowlist[operator]) {
            revert NotInAllowlist();
        }
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
    ) public override(ServiceManagerBase) whenNotPaused onlyRegistryCoordinator {
        if (allowlistEnabled && !_allowlist[operator]) {
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
        override(ServiceManagerBase)
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

    /// @notice Submit alert for verifier found a op block output mismatch.
    ///         It just a warning without any prove, the prover verifier should
    ///         submit a prove to ensure the alert is valid.
    ///         This alert can for the blocks which had not proposal its output
    ///         root to layer1, this block may not the checkpoint.
    /// @param invalidOutputRoot the invalid output root verifier got from op-devnet.
    /// @param expectOutputRoot the output root calc by verifier.
    /// @param l2BlockNumber the layer2 block 's number.
    function alertBlockMismatch(bytes32 invalidOutputRoot, bytes32 expectOutputRoot, uint256 l2BlockNumber)
        external
        onlyValidOperator
    {
        // Make sure there are no other alert, OR the currently alert is not the earliest error.
        uint256 latestAlertBlockNum = latestAlertBlockNumber();
        if (latestAlertBlockNum != 0 && l2BlockNumber >= latestAlertBlockNum) {
            revert UselessAlert();
        }

        if (l2BlockNumber == 0 || invalidOutputRoot == expectOutputRoot) {
            revert InvalidAlert();
        }

        // Make sure the block have not proposal to layer1,
        // if had proposal ouput to layer1, should use `alertBlockOutputOracleMismatch`.
        if (l2BlockNumber <= l2OutputOracle.latestBlockNumber()) {
            revert InvalidAlertType();
        }

        emit AlertBlockMismatch(invalidOutputRoot, expectOutputRoot, l2BlockNumber);

        _pushAlert(invalidOutputRoot, expectOutputRoot, 0, l2BlockNumber, msg.sender);
    }

    /// @notice Submit alert for verifier found a op block output root mismatch.
    ///         It just a warning without any prove, the prover verifier should
    ///         submit a prove to ensure the alert is valid.
    ///         This alert only for the porposaled output root by proposer,
    ///         so we just sumit the index for this output root.
    /// @param invalidOutputIndex the invalid output root index.
    /// @param expectOutputRoot the output root calc by verifier.
    function alertBlockOutputOracleMismatch(uint256 invalidOutputIndex, bytes32 expectOutputRoot)
        external
        onlyValidOperator
    {
        if (invalidOutputIndex >= l2OutputOracle.latestOutputIndex()) {
            revert InvalidAlertType();
        }

        IMachOptimismL2OutputOracle.OutputProposal memory proposal = l2OutputOracle.getL2Output(invalidOutputIndex);

        uint256 l2BlockNumber = proposal.l2BlockNumber;

        // Make sure there are no other alert, OR the currently alert is not the earliest error.
        uint256 latestAlertBlockNum = latestAlertBlockNumber();

        if (latestAlertBlockNum != 0 && l2BlockNumber >= latestAlertBlockNum) {
            revert UselessAlert();
        }

        if (l2BlockNumber == 0 || proposal.outputRoot == expectOutputRoot) {
            revert InvalidAlert();
        }

        emit AlertBlockOutputOracleMismatch(invalidOutputIndex, proposal.outputRoot, expectOutputRoot, l2BlockNumber);

        _pushAlert(proposal.outputRoot, expectOutputRoot, invalidOutputIndex, l2BlockNumber, msg.sender);
    }

    /// @notice Submit a bonsai prove receipt to mach contract.
    function submitProve(
        bytes32 imageId_,
        bytes calldata journal,
        bytes calldata seal,
        bytes32 postStateDigest,
        uint256 l2OutputIndex
    ) external onlyValidOperator {
        uint256 alertsLength = l2OutputAlerts.length;

        if (alertsLength == 0 || provedIndex == 0) {
            revert NoAlert();
        }

        if (l2OutputIndex == 0) {
            revert InvalidIndex();
        }

        if (provedIndex > alertsLength) {
            revert InvalidProvedIndex();
        }

        if (imageId == bytes32(0)) {
            revert NotInitialized();
        }

        // receipt.meta.preStateDigest, which just is the imageId in risc0
        if (imageId != imageId_) {
            revert ProveImageIdMismatch();
        }

        if (journal.length == 0) {
            revert InvalidJournal();
        }

        if (!verifier.verify(seal, imageId, postStateDigest, sha256(journal))) {
            revert ProveVerifyFailed();
        }

        // Got the per l2 ouput root info by index
        IMachOptimismL2OutputOracle.OutputProposal memory checkpoint = l2OutputOracle.getL2Output(l2OutputIndex);
        if (checkpoint.l2BlockNumber == 0 || checkpoint.outputRoot == bytes32(0)) {
            revert InvalidCheckpoint();
        }

        // Now we can trust the receipt.
        // this data is defend in guest.
        // TODO: check block header and parent output root.
        uint256 l2BlockNumber = 0;
        bytes32 outputRoot = bytes32(0);
        bytes32 headerHash = bytes32(0);
        bytes32 checkpointOutputRoot = bytes32(0);
        uint256 parentCheckpointNumber = 0;

        (headerHash, l2BlockNumber, checkpointOutputRoot, parentCheckpointNumber, outputRoot) =
            abi.decode(journal, (bytes32, uint256, bytes32, uint256, bytes32));

        L2OutputAlert memory alert = l2OutputAlerts[provedIndex - 1];

        if (l2BlockNumber != alert.l2BlockNumber) {
            revert ProveBlockNumberMismatch();
        }

        if (parentCheckpointNumber != checkpoint.l2BlockNumber) {
            revert ParentCheckpointNumberMismatch();
        }

        if (checkpointOutputRoot != checkpoint.outputRoot) {
            revert ParentCheckpointOutputRootMismatch();
        }

        uint256 invalidOutputIndex = alert.invalidOutputIndex;

        // if the output root is not equal to the expectOutputRoot, the alert is invalid.
        // TODO: In the future, we need to slash the submitter. For now we just delete it.
        if (outputRoot != alert.expectOutputRoot) {
            if (outputRoot == alert.invalidOutputRoot) {
                if (provedIndex < alertsLength) {
                    for (uint256 i = provedIndex; i < alertsLength; i++) {
                        l2OutputAlerts[i] = l2OutputAlerts[i + 1];
                    }
                }

                l2OutputAlerts.pop();

                emit AlertDelete(invalidOutputIndex, alert.expectOutputRoot, outputRoot, l2BlockNumber, alert.submitter);
            } else {
                l2OutputAlerts[provedIndex - 1].expectOutputRoot = outputRoot;
                l2OutputAlerts[provedIndex - 1].submitter = msg.sender;

                emit AlertReset(
                    invalidOutputIndex, alert.invalidOutputRoot, outputRoot, l2BlockNumber, alert.submitter, msg.sender
                );
            }
        }

        provedIndex = provedIndex - 1;

        emit SubmittedBlockProve(invalidOutputIndex, outputRoot, l2BlockNumber);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                               View Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    ///  @notice Get the address for RegistryCoordinator,
    ///  it help the verifier to check if self is a valid operator.
    function getRegistryCoordinatorAddress() public view returns (address) {
        return address(_registryCoordinator);
    }

    /// @notice Get total alert length
    function getAlertsLength() public view returns (uint256) {
        return l2OutputAlerts.length;
    }

    /// @notice Return the latest alert 's block number, if not exist, just return 0.
    ///         TODO: we can add more view functions to get details info about alert.
    ///         This function just used for verifier check if need commit more
    ///         alerts to contract.
    function latestAlertBlockNumber() public view returns (uint256) {
        return l2OutputAlerts.length == 0 ? 0 : l2OutputAlerts[l2OutputAlerts.length - 1].l2BlockNumber;
    }

    /// @notice Get Alert by index
    function getAlert(uint256 index) external view returns (L2OutputAlert memory) {
        if (index >= l2OutputAlerts.length) {
            revert InvalidIndex();
        }
        return l2OutputAlerts[index];
    }

    /// @notice Return the latest no proved alert 's block number, if not exist, just return 0.
    function latestUnprovedBlockNumber() external view returns (uint256) {
        if (provedIndex == 0) {
            return 0;
        }

        if (provedIndex > l2OutputAlerts.length) {
            revert InvalidProvedIndex();
        }

        return l2OutputAlerts.length == 0 ? 0 : l2OutputAlerts[provedIndex - 1].l2BlockNumber;
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Internal Functions                          //
    //////////////////////////////////////////////////////////////////////////////

    /// @notice push new alert
    function _pushAlert(
        bytes32 invalidOutputRoot,
        bytes32 expectOutputRoot,
        uint256 invalidOutputIndex,
        uint256 l2BlockNumber,
        address sender
    ) private {
        l2OutputAlerts.push(
            L2OutputAlert({
                l2BlockNumber: l2BlockNumber,
                invalidOutputIndex: invalidOutputIndex,
                invalidOutputRoot: invalidOutputRoot,
                expectOutputRoot: expectOutputRoot,
                submitter: sender
            })
        );

        // For the proved output, if there are exist a early block alert
        // we will make it not proved! so we just set to `length`
        provedIndex = l2OutputAlerts.length;
    }
}
