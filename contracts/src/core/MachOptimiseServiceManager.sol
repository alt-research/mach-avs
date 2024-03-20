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
import {IMachOptimismL2OutputOracle} from "../../interfaces/IMachOptimismL2OutputOracle.sol";
import {IRiscZeroVerifier} from "../../interfaces/IRiscZeroVerifier.sol";

contract MachOptimiseServiceManager is MachServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker, Pausable {
    using EnumerableSet for EnumerableSet.Bytes32Set;
    using EnumerableSet for EnumerableSet.AddressSet;

    uint8 internal constant PAUSED_CONFIRM_ALERT = 0;

    IRiscZeroVerifier public verifier;

    IMachOptimismL2OutputOracle public l2OutputOracle;

    // The imageId for risc0 guest code.
    bytes32 public imageId;

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
        address _batchConfirmer,
        bytes32 imageId_,
        IMachOptimismL2OutputOracle l2OutputOracle_,
        IRiscZeroVerifier verifier_
    ) public initializer {
        _initializePauser(_pauserRegistry, _initialPausedStatus);
        _transferOwnership(_initialOwner);
        _setAlertConfirmer(_batchConfirmer);
        l2OutputOracle = l2OutputOracle_;
        verifier = verifier_;
        imageId = imageId_;
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
    function deregisterOperator() external whenNotPaused {
        address operator = msg.sender;
        _operators.remove(operator);
        _avsDirectory.deregisterOperatorFromAVS(operator);
        emit OperatorRemoved(operator);
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Alert Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    function confirmAlert(
        bytes32 messageHash,
        bytes32 expectOutputRoot,
        bytes calldata journal,
        bytes calldata seal,
        bytes32 postStateDigest,
        uint256 l2OutputIndex
    ) external whenNotPaused onlyAlertConfirmer {
        require(
            verifier.verify(seal, imageId, postStateDigest, sha256(journal)),
            "MachServiceManager.confirmAlert: verify failed"
        );

        // Got the per l2 ouput root info by index
        IMachOptimismL2OutputOracle.OutputProposal memory checkpoint = l2OutputOracle.getL2Output(l2OutputIndex);
        require(
            checkpoint.l2BlockNumber != 0 && checkpoint.outputRoot != bytes32(0),
            "MachServiceManager.confirmAlert: invalid checkpoint"
        );
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
        require(outputRoot == expectOutputRoot, "MachServiceManager.confirmAlert: invalid outputRoot");

        // store alert
        _messageHashes.add(messageHash);

        emit AlertConfirmed(messageHash, messageHash);
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
