// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity ^0.8.12;

import {IMachOptimismL2OutputOracle} from "../interfaces/IMachOptimismL2OutputOracle.sol";
import {IRiscZeroVerifier} from "../interfaces/IMachOptimism.sol";
import {IMachOptimism} from "../interfaces/IMachOptimism.sol";

contract MachOptimismZkServiceManagerStorage {
    uint256 public immutable settlementChainID;
    uint256 public immutable rollupChainID;
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

    constructor(uint256 settlementChainID_, uint256 rollupChainID_) {
        settlementChainID = settlementChainID_;
        rollupChainID = rollupChainID_;
    }
}
