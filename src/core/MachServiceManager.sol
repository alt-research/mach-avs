// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {IMachOptimism, CallbackAuthorization, IRiscZeroVerifier} from "../interfaces/IMachOptimism.sol";
import {IMachOptimismL2OutputOracle} from "../interfaces/IMachOptimismL2OutputOracle.sol";
import {MachServiceManagerStorage} from "./MachServiceManagerStorage.sol";

contract MachServiceManager is MachServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker {
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

    modifier onlyValidOperator() {
        IBLSApkRegistry blsApkRegistry = _registryCoordinator.blsApkRegistry();
        bytes32 operatorId = blsApkRegistry.getOperatorId(msg.sender);
        require(operatorId != bytes32(0), "onlyValidOperator: not valid operator");
        _;
    }

    /// @notice Initializes the contract with provided parameters.
    function initialize(
        address owner_,
        bytes32 imageId_,
        IMachOptimismL2OutputOracle l2OutputOracle_,
        IRiscZeroVerifier verifier_
    ) external initializer {
        __Ownable_init();
        _transferOwnership(owner_);
        if (address(l2OutputOracle_) == address(0) || address(verifier_) == address(0) || imageId_ == bytes32(0)) {
            revert ZeroAddress();
        }
        l2OutputOracle = l2OutputOracle_;
        verifier = verifier_;
        imageId = imageId_;
    }

    /// owner functions

    function setImageId(bytes32 imageId_) external onlyOwner {
        imageId = imageId_;
    }

    function setRiscZeroVerifier(IRiscZeroVerifier verifier_) external onlyOwner {
        if (address(verifier_) == address(0)) {
            revert ZeroAddress();
        }

        verifier = verifier_;
    }

    function clearAlerts() external onlyOwner {
        delete l2OutputAlerts;
        provedIndex = 0;
    }

    /// operator functions

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

    /// public functions

    /// @notice Return the latest alert 's block number, if not exist, just return 0.
    ///         TODO: we can add more view functions to get details info about alert.
    ///         This function just used for verifier check if need commit more
    ///         alerts to contract.
    function latestAlertBlockNumber() public view returns (uint256) {
        return l2OutputAlerts.length == 0 ? 0 : l2OutputAlerts[l2OutputAlerts.length - 1].l2BlockNumber;
    }

    /// external view functions

    function getAlert(uint256 index) external view returns (IMachOptimism.L2OutputAlert memory) {
        if (index >= l2OutputAlerts.length) {
            revert InvalidIndex();
        }
        return l2OutputAlerts[index];
    }

    function getAlertsLength() external view returns (uint256) {
        return l2OutputAlerts.length;
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
}
