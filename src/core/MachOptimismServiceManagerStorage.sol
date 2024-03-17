// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IMachOptimisimServiceManager} from "../interfaces/IMachOptimisimServiceManager.sol";
import {IMachOptimismL2OutputOracle} from "../interfaces/IMachOptimismL2OutputOracle.sol";
import {IMachOptimism, CallbackAuthorization, IRiscZeroVerifier} from "../interfaces/IMachOptimism.sol";

abstract contract MachOptimismServiceManagerStorage is IMachOptimisimServiceManager {
    // The imageId for risc0 guest code.
    bytes32 public imageId;

    IRiscZeroVerifier public verifier;

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

    IMachOptimismL2OutputOracle public l2OutputOracle;

    /// @notice push new alert
    function _pushAlert(
        bytes32 invalidOutputRoot,
        bytes32 expectOutputRoot,
        uint256 invalidOutputIndex,
        uint256 l2BlockNumber,
        address sender
    ) internal {
        l2OutputAlerts.push(
            IMachOptimism.L2OutputAlert({
                l2BlockNumber: l2BlockNumber,
                invalidOutputIndex: invalidOutputIndex,
                invalidOutputRoot: invalidOutputRoot,
                expectOutputRoot: expectOutputRoot,
                submitter: sender
            })
        );

        // For the proved output, if there are exist a early block alert
        // we will make it not proved! so we just set to `length`
        provedIndex = l2OutputAlerts.length;
    }
}
