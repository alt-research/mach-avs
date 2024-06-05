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

contract MachServiceManagerRegistry is OwnableUpgradeable {
    // rollup chain ID => service manager
    mapping(uint256 => IMachServiceManager) public serviceManagers;

    event ServiceManagerRegistered(uint256 rollupChainId_, IMachServiceManager serviceManager_, address sender);

    function initialize() external initializer {
        __Ownable_init();
    }

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

    function hasActiveAlerts(uint256 rollupChainId_) external view returns (bool) {
        return serviceManagers[rollupChainId_].totalAlerts(rollupChainId_) > 0;
    }
}
