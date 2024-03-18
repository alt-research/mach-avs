// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IMachOptimism} from "../interfaces/IMachOptimism.sol";

interface IMachServiceManager is IServiceManager {
    /**
     * @notice Emitted when the alert confirmer is changed.
     * @param previousAddress The address of the previous alert confirmer
     * @param newAddress The address of the new alert confirmer
     */
    event AlertConfirmerChanged(address previousAddress, address newAddress);

    /**
     * @notice Emitted when a Alert is confirmed.
     * @param alertHeaderHash The hash of the alert header
     * @param blockNumber The l2 block number
     */
    event AlertConfirmed(bytes32 indexed alertHeaderHash, uint256 blockNumber);
}
