// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity =0.8.12;
import {ISlasher} from "eigenlayer-contracts/src/contracts/interfaces/ISlasher.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import "eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {ServiceManagerBase, IRegistryCoordinator, IStakeRegistry} from "eigenlayer-middleware/src/ServiceManagerBase.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import "../Error.sol";

interface IVoteWeigher {
    /**
     * @notice This function computes the total weight of the @param operator in the quorum @param quorumNumber.
     * @dev reverts in the case that `quorumNumber` is greater than or equal to `quorumCount`
     */
    function weightOfOperatorForQuorum(
        uint8 quorumNumber,
        address operator
    ) external view returns (uint96);
}

contract VitalServiceManager is ServiceManagerBase {
    address public accel;

    event Freeze(address freezed);
    error NotAccel();

    constructor(
        IDelegationManager _delegationManager,
        IRegistryCoordinator _registryCoordinator,
        IStakeRegistry _stakeRegistry
    )
        ServiceManagerBase(
            _delegationManager,
            _registryCoordinator,
            _stakeRegistry
        )
    {}

    /// @notice Sets the Accel proxy address.
    /// @param accel_ The address of Accel proxy.
    function setAccel(address accel_) external onlyOwner {
        accel = accel_;
    }

    /// @notice Called in the event of challenge resolution, in order to forward a call to the Slasher, which 'freezes' the `operator`.
    /// @dev The Slasher contract is under active development and its interface expected to change.
    ///      We recommend writing slashing logic without integrating with the Slasher at this point in time.
    function freezeOperator(address operatorAddr) external {
        if (msg.sender != accel) {
            revert NotAccel();
        }
        emit Freeze(operatorAddr);
        // slasher.freezeOperator(operatorAddr);
    }

    function isEligibleForChallenge(address user) external view returns (bool) {
        // TODO: fix the naive logic
        return
            IVoteWeigher(address(_stakeRegistry)).weightOfOperatorForQuorum(
                0,
                user
            ) >= _stakeRegistry.minimumStakeForQuorum(0);
    }
}
