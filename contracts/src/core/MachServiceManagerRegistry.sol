// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity =0.8.12;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {IMachServiceManager} from "../interfaces/IMachServiceManager.sol";
import {ZeroAddress, AlreadyAdded} from "../error/Errors.sol";

/// @title MachServiceManagerRegistry
/// @notice This contract allows the owner to register service managers for specific rollup chain IDs and check if they have active alerts.
contract MachServiceManagerRegistry is OwnableUpgradeable {
    // Mapping of rollup chain ID to service manager
    mapping(uint256 => IMachServiceManager) public serviceManagers;

    /// @notice Emitted when a service manager is registered
    /// @param rollupChainId_ The rollup chain ID
    /// @param serviceManager_ The registered service manager
    /// @param sender The address that registered the service manager
    event ServiceManagerRegistered(uint256 rollupChainId_, IMachServiceManager serviceManager_, address sender);

    /// @notice Initializes the contract and sets the deployer as the owner
    function initialize() external initializer {
        __Ownable_init();
    }

    /// @notice Registers a service manager for a specific rollup chain ID
    /// @param rollupChainId_ The rollup chain ID
    /// @param serviceManager_ The service manager to be registered
    /// @dev Reverts if the service manager address is zero or already registered
    function registerServiceManager(uint256 rollupChainId_, IMachServiceManager serviceManager_) external onlyOwner {
        if (address(serviceManager_) == address(0)) {
            revert ZeroAddress();
        }
        if (serviceManagers[rollupChainId_] == serviceManager_) {
            revert AlreadyAdded();
        }
        serviceManagers[rollupChainId_] = serviceManager_;
        emit ServiceManagerRegistered(rollupChainId_, serviceManager_, _msgSender());
    }

    /// @notice Checks if a service manager has active alerts
    /// @param rollupChainId_ The rollup chain ID
    /// @return True if there are active alerts, false otherwise
    function hasActiveAlerts(uint256 rollupChainId_) external view returns (bool) {
        return serviceManagers[rollupChainId_].totalAlerts(rollupChainId_) > 0;
    }
}
