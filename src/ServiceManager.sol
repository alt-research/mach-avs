// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {IBLSRegistryCoordinatorWithIndices, ServiceManagerBase, IBLSRegistryCoordinatorWithIndices, ISlasher} from "eigenlayer-middleware/src/ServiceManagerBase.sol";

/// @title AVS service manager.
contract ServiceManager is ServiceManagerBase {
    address public accel;

    event Freeze(address freezed);

    constructor(
        IBLSRegistryCoordinatorWithIndices _registryCoordinator,
        ISlasher _slasher
    ) ServiceManagerBase(_registryCoordinator, _slasher) {}

    function setAccel(address accel_) external onlyOwner {
        accel = accel_;
    }

    /// @notice Called in the event of challenge resolution, in order to forward a call to the Slasher, which 'freezes' the `operator`.
    /// @dev The Slasher contract is under active development and its interface expected to change.
    ///      We recommend writing slashing logic without integrating with the Slasher at this point in time.
    function freezeOperator(address operatorAddr) external override {
        require(msg.sender == accel, "NotAccel");
        emit Freeze(operatorAddr);
        // slasher.freezeOperator(operatorAddr);
    }
}
