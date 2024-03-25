// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity ^0.8.12;

import {IRiscZeroVerifier} from "./IRiscZeroVerifier.sol";
import {CallbackAuthorization} from "./IBonsaiRelay.sol";

/// @title IMachOptimism
/// @notice The Interface for a Mach optimism contract.
interface IMachOptimism {
    event AlertBlockMismatch(bytes32 invalidOutputRoot, bytes32 expectOutputRoot, uint256 indexed l2BlockNumber);
    event AlertBlockOutputOracleMismatch(
        uint256 indexed invalidOutputIndex,
        bytes32 invalidOutputRoot,
        bytes32 expectOutputRoot,
        uint256 indexed l2BlockNumber
    );
    event SubmittedBlockProve(uint256 indexed invalidOutputIndex, bytes32 OutputRoot, uint256 indexed l2BlockNumber);

    event AlertDelete(
        uint256 indexed invalidOutputIndex,
        bytes32 expectOutputRoot,
        bytes32 OutputRoot,
        uint256 indexed l2BlockNumber,
        address indexed submitter
    );
    event AlertReset(
        uint256 indexed invalidOutputIndex,
        bytes32 invalidOutputRoot,
        bytes32 expectOutputRoot,
        uint256 indexed l2BlockNumber,
        address fromSubmitter,
        address indexed toSubmitter
    );

    /**
     * @notice Emitted when an operator is added to the MachServiceManagerAVS.
     * @param operator The address of the operator
     */
    event OperatorAdded(address indexed operator);

    /**
     * @notice Emitted when an operator is removed from the MachServiceManagerAVS.
     * @param operator The address of the operator
     */
    event OperatorRemoved(address indexed operator);

    /**
     * @notice Emitted when an operator is added to the allowlist.
     * @param operator The operator
     */
    event OperatorAllowed(address operator);

    /**
     * @notice Emitted when an operator is removed from the allowlist.
     * @param operator The operator
     */
    event OperatorDisallowed(address operator);

    /**
     * @notice Emitted when the allowlist is enabled.
     */
    event AllowlistEnabled();

    /**
     * @notice Emitted when the allowlist is disabled.
     */
    event AllowlistDisabled();

    struct L2OutputAlert {
        uint256 l2BlockNumber;
        uint256 invalidOutputIndex;
        bytes32 invalidOutputRoot;
        bytes32 expectOutputRoot;
        address submitter;
    }

    /// Returns the datas for alert by its index.
    function getAlert(uint256 index) external view returns (L2OutputAlert memory);

    /// @notice Return the latest alert 's block number, if not exist, just return 0.
    ///         TODO: we can add more view functions to get details info about alert.
    ///         This function just used for verifier check if need commit more
    ///         alerts to contract.
    function latestAlertBlockNumber() external view returns (uint256);

    /// @notice Return the latest no proved alert 's block number, if not exist, just return 0.
    function latestUnprovedBlockNumber() external view returns (uint256);

    /// @notice Submit alert for verifier found a op block output mismatch.
    ///         It just a warning without any prove, the prover verifier should
    ///         submit a prove to ensure the alert is valid.
    ///         This alert can for the blocks which had not proposal its output
    ///         root to layer1, this block may not the checkpoint.
    /// @param invalidOutputRoot the invalid output root verifier got from op-devnet.
    /// @param expectOutputRoot the output root calc by verifier.
    /// @param l2BlockNumber the layer2 block 's number.
    function alertBlockMismatch(bytes32 invalidOutputRoot, bytes32 expectOutputRoot, uint256 l2BlockNumber) external;

    /// @notice Submit alert for verifier found a op block output root mismatch.
    ///         It just a warning without any prove, the prover verifier should
    ///         submit a prove to ensure the alert is valid.
    ///         This alert only for the proposed output root by proposer,
    ///         so we just submit the index for this output root.
    /// @param invalidOutputIndex the invalid output root index.
    /// @param expectOutputRoot the output root calc by verifier.
    function alertBlockOutputOracleMismatch(uint256 invalidOutputIndex, bytes32 expectOutputRoot) external;

    /// @notice Submit a bonsai prove receipt to mach contract.
    function submitProve(
        bytes32 imageId_,
        bytes calldata journal,
        bytes calldata seal,
        bytes32 postStateDigest,
        uint256 perL2OutputIndex
    ) external;
}
