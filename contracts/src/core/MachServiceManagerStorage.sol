// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IMachServiceManager} from "../interfaces/IMachServiceManager.sol";

abstract contract MachServiceManagerStorage is IMachServiceManager {
    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    EnumerableSet.Bytes32Set internal _messageHashes;

    /// @notice Ethereum addresses of currently register operators
    EnumerableSet.AddressSet internal _operators;

    /// @notice address that is permissioned to confirm alerts
    address public alertConfirmer;

    /// @notice Set of operators that are allowed to register
    mapping(address => bool) internal _allowlist;

    /// @notice Whether or not the allowlist is enabled
    bool public allowlistEnabled = true;

    /// @notice when applied to a function, ensures that the function is only callable by the `alertConfirmer`.
    modifier onlyAlertConfirmer() {
        require(msg.sender == alertConfirmer, "onlyAlertConfirmer: not from alert confirmer");
        _;
    }
}
