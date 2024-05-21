// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity =0.8.12;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

abstract contract MachServiceManagerStorage {
    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    // slot 0
    /// @notice Allowed rollup chain IDs
    mapping(uint256 => bool) public rollupChainIDs;

    // Slot 1
    mapping(uint256 => EnumerableSet.Bytes32Set) internal _messageHashes;

    // Slot 2, 3
    /// @notice Ethereum addresses of currently register operators
    EnumerableSet.AddressSet internal _operators;

    // Slot 4
    /// @notice Set of operators that are allowed to register
    mapping(address => bool) public allowlist;

    // Slot 5
    /// @notice address that is permissioned to confirm alerts
    address public alertConfirmer;

    /// @notice Whether or not the allowlist is enabled
    bool public allowlistEnabled;

    /// @notice Minimal quorum threshold percentage
    uint8 public quorumThresholdPercentage;

    // slot 6
    /// @notice Resolved message hashes, prevent aggregator from replay any resolved alert
    mapping(uint256 => EnumerableSet.Bytes32Set) internal _resolvedMessageHashes;

    // slot 7
    /// @notice Role for whitelisting operators
    address public whitelister;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[42] private __GAP;
}
