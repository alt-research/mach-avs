// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity =0.8.12;
import "eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {IBLSRegistryCoordinatorWithIndices, ServiceManagerBase, IBLSRegistryCoordinatorWithIndices, ISlasher} from "eigenlayer-middleware/src/ServiceManagerBase.sol";
import "./Error.sol";
import {IMachOptimism, CallbackAuthorization, IRiscZeroVerifier} from "./interfaces/IMachOptimism.sol";
import {IMachOptimismL2OutputOracle} from "./interfaces/IMachOptimismL2OutputOracle.sol";

contract ServiceManager is IMachOptimism, ServiceManagerBase {
    IMachOptimismL2OutputOracle public l2OutputOracle;
    IRiscZeroVerifier public verifier;
    // The imageId for risc0 guest code.
    bytes32 public imageId;

    // Alerts for blocks, the tail is for earliest block.
    // For the proved output, if there are exist a early block alert
    // we will make it not proved!
    IMachOptimism.L2OutputAlert[] internal l2OutputAlerts;
    // The next index for no proved alert,
    // `l2OutputAlerts[provedIndex - 1]` is the first no proved alerts,
    // if is 0, means all alert is proved,
    // if provedIndex == l2OutputAlerts.length, means all alert is not proved,
    // the prover just need prove the earliest no proved alert,
    uint256 public provedIndex;
    event Freeze(address freezed);

    constructor(
        IBLSRegistryCoordinatorWithIndices _registryCoordinator,
        ISlasher _slasher
    ) ServiceManagerBase(_registryCoordinator, _slasher) {}

    modifier onlyValidOperator() {
        bytes32 operatorId = registryCoordinator.getOperatorId(msg.sender);
        if (operatorId == bytes32(0)) {
            revert NotOperator();
        }
        _;
    }

    /// @notice Called in the event of challenge resolution, in order to forward a call to the Slasher, which 'freezes' the `operator`.
    /// @dev The Slasher contract is under active development and its interface expected to change.
    ///      We recommend writing slashing logic without integrating with the Slasher at this point in time.
    function freezeOperator(address operatorAddr) external override onlyOwner {
        emit Freeze(operatorAddr);
        // slasher.freezeOperator(operatorAddr);
    }

    /// @notice Initializes the contract with provided parameters.
    function initialize(
        IPauserRegistry pauserRegistry_,
        address initialOwner_,
        bytes32 imageId_,
        IMachOptimismL2OutputOracle l2OutputOracle_,
        IRiscZeroVerifier verifier_
    ) external {
        super.initialize(pauserRegistry_, initialOwner_);
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

    function setImageId(bytes32 imageId_) external onlyOwner {
        imageId = imageId_;
    }

    function clearAlerts() external onlyOwner {
        delete l2OutputAlerts;
        provedIndex = 0;
    }

    function getAlert(
        uint256 index
    ) external view returns (L2OutputAlert memory) {
        if (index >= l2OutputAlerts.length) {
            revert InvalidIndex();
        }
        return l2OutputAlerts[index];
    }

    function getAlertsLength() public view returns (uint256) {
        return l2OutputAlerts.length;
    }

    /// @notice Return the latest alert 's block number, if not exist, just return 0.
    ///         TODO: we can add more view functions to get details info about alert.
    ///         This function just used for verifier check if need commit more
    ///         alerts to contract.
    function latestAlertBlockNumber() public view returns (uint256) {
        return
            l2OutputAlerts.length == 0
                ? 0
                : l2OutputAlerts[l2OutputAlerts.length - 1].l2BlockNumber;
    }

    /// @notice Return the latest no proved alert 's block number, if not exist, just return 0.
    function latestUnprovedBlockNumber() external view returns (uint256) {
        if (provedIndex == 0) {
            return 0;
        }

        if (provedIndex > l2OutputAlerts.length) {
            revert InvalidProvedIndex();
        }

        return
            l2OutputAlerts.length == 0
                ? 0
                : l2OutputAlerts[provedIndex - 1].l2BlockNumber;
    }

    /// @notice Submit alert for verifier found a op block output mismatch.
    ///         It just a warning without any prove, the prover verifier should
    ///         submit a prove to ensure the alert is valid.
    ///         This alert can for the blocks which had not proposal its output
    ///         root to layer1, this block may not the checkpoint.
    /// @param invalidOutputRoot the invalid output root verifier got from op-devnet.
    /// @param expectOutputRoot the output root calc by verifier.
    /// @param l2BlockNumber the layer2 block 's number.
    function alertBlockMismatch(
        bytes32 invalidOutputRoot,
        bytes32 expectOutputRoot,
        uint256 l2BlockNumber
    ) external onlyValidOperator {
        // Make sure there are no other alert, OR the currently alert is not the earliest error.
        uint256 latestAlertBlockNumber = latestAlertBlockNumber();
        if (
            latestAlertBlockNumber != 0 &&
            l2BlockNumber >= latestAlertBlockNumber
        ) {
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

        emit AlertBlockMismatch(
            invalidOutputRoot,
            expectOutputRoot,
            l2BlockNumber
        );

        _pushAlert(
            invalidOutputRoot,
            expectOutputRoot,
            0,
            l2BlockNumber,
            msg.sender
        );
    }

    /// @notice Submit alert for verifier found a op block output root mismatch.
    ///         It just a warning without any prove, the prover verifier should
    ///         submit a prove to ensure the alert is valid.
    ///         This alert only for the porposaled output root by proposer,
    ///         so we just sumit the index for this output root.
    /// @param invalidOutputIndex the invalid output root index.
    /// @param expectOutputRoot the output root calc by verifier.
    function alertBlockOutputOracleMismatch(
        uint256 invalidOutputIndex,
        bytes32 expectOutputRoot
    ) external onlyValidOperator {
        if (invalidOutputIndex >= l2OutputOracle.latestOutputIndex()) {
            revert InvalidAlertType();
        }

        IMachOptimismL2OutputOracle.OutputProposal
            memory proposal = l2OutputOracle.getL2Output(invalidOutputIndex);

        uint256 l2BlockNumber = proposal.l2BlockNumber;

        // Make sure there are no other alert, OR the currently alert is not the earliest error.
        uint256 latestBlockNumber = latestAlertBlockNumber();

        if (latestBlockNumber != 0 && l2BlockNumber >= latestBlockNumber) {
            revert UselessAlert();
        }

        if (l2BlockNumber == 0 || proposal.outputRoot == expectOutputRoot) {
            revert InvalidAlert();
        }

        emit AlertBlockOutputOracleMismatch(
            invalidOutputIndex,
            proposal.outputRoot,
            expectOutputRoot,
            l2BlockNumber
        );

        _pushAlert(
            proposal.outputRoot,
            expectOutputRoot,
            invalidOutputIndex,
            l2BlockNumber,
            msg.sender
        );
    }

    /// @notice Submit a bonsai prove receipt to mach contract.
    function submitProve(
        bytes32 imageId_,
        bytes calldata journal,
        bytes calldata seal,
        bytes32 postStateDigest
    ) external onlyValidOperator {
        uint256 alertsLength = l2OutputAlerts.length;

        if (alertsLength == 0 || provedIndex == 0) {
            revert NoAlert();
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

        // Now we can trust the receipt.
        // this data is defend in guest.
        // TODO: check block header and parent output root.
        uint256 l2BlockNumber = 0;
        bytes32 outputRoot = bytes32(0);
        bytes32 headerHash = bytes32(0);
        bytes32 parentHeaderHash = bytes32(0);

        (headerHash, l2BlockNumber, parentHeaderHash, outputRoot) = abi.decode(
            journal,
            (bytes32, uint256, bytes32, bytes32)
        );

        L2OutputAlert memory alert = l2OutputAlerts[provedIndex - 1];

        if (l2BlockNumber != alert.l2BlockNumber) {
            revert ProveBlockNumberMismatch();
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

                _clearBlockAlertsUpTo(l2BlockNumber);

                emit AlertDelete(
                    invalidOutputIndex,
                    alert.expectOutputRoot,
                    outputRoot,
                    l2BlockNumber,
                    alert.submitter
                );
            } else {
                l2OutputAlerts[provedIndex - 1].expectOutputRoot = outputRoot;
                l2OutputAlerts[provedIndex - 1].submitter = msg.sender;

                emit AlertReset(
                    invalidOutputIndex,
                    alert.invalidOutputRoot,
                    outputRoot,
                    l2BlockNumber,
                    alert.submitter,
                    msg.sender
                );
            }
        }

        provedIndex = provedIndex - 1;

        emit SubmittedBlockProve(invalidOutputIndex, outputRoot, l2BlockNumber);
    }

    function _clearBlockAlertsUpTo(uint256 l2BlockNumber) internal {
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
