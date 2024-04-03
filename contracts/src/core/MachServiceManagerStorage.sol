// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity ^0.8.12;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

abstract contract MachServiceManagerStorage {
    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    /// @notice Rollup chain id, it is different from block.chainid
    uint256 public immutable rollupChainId;

    EnumerableSet.Bytes32Set internal _messageHashes;

    /// @notice Ethereum addresses of currently register operators
    EnumerableSet.AddressSet internal _operators;

    /// @notice Set of operators that are allowed to register
    mapping(address => bool) public allowlist;

    /// @notice address that is permissioned to confirm alerts
    address public alertConfirmer;

    /// @notice Whether or not the allowlist is enabled
    bool public allowlistEnabled;

    /// @notice Minimal quorum threshold percentage
    uint8 public quorumThresholdPercentage;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[46] private __GAP;

    constructor(uint256 _rollupChainId) {
        rollupChainId = _rollupChainId;
    }
}
