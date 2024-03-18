// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IMachServiceManager} from "../interfaces/IMachServiceManager.sol";

abstract contract MachServiceManagerStorage is IMachServiceManager {
    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    EnumerableSet.UintSet internal _l2Blocks;

    /// @notice address that is permissioned to confirm alerts
    address public alertConfirmer;

    /// @notice when applied to a function, ensures that the function is only callable by the `alertConfirmer`.
    modifier onlyAlertConfirmer() {
        require(msg.sender == alertConfirmer, "onlyAlertConfirmer: not from alert confirmer");
        _;
    }
}
