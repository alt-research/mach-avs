// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {MockRegistryCoordinator} from "./MockRegistryCoordinator.sol";
import {InvalidOperator} from "../error/Errors.sol";

contract MockMachServiceManager is OwnableUpgradeable {
    MockRegistryCoordinator public registryCoordinator;

    constructor(MockRegistryCoordinator __registryCoordinator) {
        registryCoordinator = __registryCoordinator;
    }

    function registerOperatorToAVS(address operator) public {
        if (_msgSender() != operator) {
            revert InvalidOperator();
        }
    }

    /**
     * return the sender of this call.
     * if the call came through our trusted forwarder, return the original sender.
     * otherwise, return `msg.sender`.
     * should be used in the contract anywhere instead of msg.sender
     */
    function _msgSender() internal view override returns (address ret) {
        if (msg.data.length >= 20 && isTrustedForwarder(msg.sender)) {
            // At this point we know that the sender is a trusted forwarder,
            // so we trust that the last bytes of msg.data are the verified sender address.
            // extract sender address from the end of msg.data
            assembly {
                ret := shr(96, calldataload(sub(calldatasize(), 20)))
            }
        } else {
            ret = msg.sender;
        }
    }

    function isTrustedForwarder(address _forwarder) internal view returns (bool) {
        return address(registryCoordinator) == _forwarder;
    }
}
