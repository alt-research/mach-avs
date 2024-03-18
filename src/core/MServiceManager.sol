// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {MachServiceManagerStorage} from "./MServiceManagerStorage.sol";

contract MachServiceManager is MachServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker {
    struct AlertHeader {
        uint256 l2BlockNumber;
        bytes quorumNumbers; // each byte is a different quorum number
        bytes quorumThresholdPercentages; // every bytes is an amount less than 100 specifying the percentage of stake
            // the must have signed in the corresponding quorum in `quorumNumbers`
        uint32 referenceBlockNumber;
    }

    /// @notice when applied to a function, ensures that the function is only callable by the `alertConfirmer`.
    modifier onlyAlertConfirmer() {
        require(msg.sender == alertConfirmer, "onlyAlertConfirmer: not from alert confirmer");
        _;
    }

    constructor(
        IAVSDirectory __avsDirectory,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry
    )
        BLSSignatureChecker(__registryCoordinator)
        ServiceManagerBase(__avsDirectory, __registryCoordinator, __stakeRegistry)
    {
        _disableInitializers();
    }

    function confirmAlert(
        AlertHeader calldata alertHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyAlertConfirmer {}
}
