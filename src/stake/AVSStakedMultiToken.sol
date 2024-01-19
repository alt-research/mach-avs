// SPDX-License-Identifier: agpl-3.0
pragma solidity ^0.8.23;

import {StakedMultiToken} from "./StakedMultiToken.sol";

contract AVSStakedMultiToken is StakedMultiToken {
    error NotAVS();

    address public immutable avs;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(address avs_) {
        avs = avs_;
        _disableInitializers();
    }

    /// @notice Stakes tokens.
    /// @param to The address of the user for whom the tokens are being staked
    /// @param id The identifier of the staking pool
    /// @param amount The amount of tokens to be staked
    function stake(address to, uint256 id, uint256 amount) public override {
        if (avs != _msgSender()) {
            revert NotAVS();
        }
        super.stake(to, id, amount);
    }
}
