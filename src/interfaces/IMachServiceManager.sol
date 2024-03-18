// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IMachOptimism} from "../interfaces/IMachOptimism.sol";

interface IMachServiceManager is IServiceManager, IMachOptimism {
    error ZeroAddress();
    error InvalidIndex();
    error InvalidProvedIndex();
    error InvalidAlertType();
    error UselessAlert();
    error InvalidAlert();
    error AlreadyInitialized();
    error NotInitialized();
    error ZeroValue();
    error InvalidCheckpoint();
    error ProveImageIdMismatch();
    error ProveBlockNumberMismatch();
    error ProveOutputRootMismatch();
    error ParentCheckpointNumberMismatch();
    error ParentCheckpointOutputRootMismatch();
    error ProveVerifyFailed();
    error InvalidJournal();
    error NoAlert();
    error NotOperator();
}
