// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IMachOptimism} from "../interfaces/IMachOptimism.sol";

interface IMachServiceManager is IServiceManager {
    struct AlertHeader {
        bytes32 messageHash;
        bytes quorumNumbers; // each byte is a different quorum number
        bytes quorumThresholdPercentages; // every bytes is an amount less than 100 specifying the percentage of stake
            // the must have signed in the corresponding quorum in `quorumNumbers`
        uint32 referenceBlockNumber;
    }

    struct ReducedAlertHeader {
        bytes32 messageHash;
        uint32 referenceBlockNumber;
    }

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
     * @notice Emitted when the alert confirmer is changed.
     * @param previousAddress The address of the previous alert confirmer
     * @param newAddress The address of the new alert confirmer
     */
    event AlertConfirmerChanged(address previousAddress, address newAddress);

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

    /**
     * @notice Emitted when a Alert is confirmed.
     * @param alertHeaderHash The hash of the alert header
     * @param messageHash The message hash
     */
    event AlertConfirmed(bytes32 indexed alertHeaderHash, bytes32 messageHash);

    event AlertRemoved(bytes32 messageHash, address sender);
}
