// SPDX-License-Identifier: agpl-3.0
pragma solidity ^0.8.23;

import {AccessControlUpgradeable} from "@openzeppelin-upgrades/contracts/access/AccessControlUpgradeable.sol";

/// @dev Accounting contract to manage staking distributions
/// This is adapted from https://github.com/bgd-labs/aave-stk-v1-5/blob/8867dd5b1137d4d46acd9716fe98759cb16b1606/src/contracts/AaveDistributionManager.sol
// solhint-disable not-rely-on-time
contract DistributionManager is AccessControlUpgradeable {
    error NotEmissionAdmin();

    bytes32 public constant EMISSION_ADMIN_ROLE =
        keccak256("EMISSION_ADMIN_ROLE");
    uint256 public constant PRECISION_FACTOR = 1e18;

    uint256 public distributionEnd;
    uint128 public emissionPerSecond;
    mapping(uint256 => uint128) public distributionTimestamps;
    mapping(uint256 => uint256) public distributionIndices;
    /// @dev Flag determining if there's an ongoing slashing event that needs to be settled
    mapping(uint256 => bool) public distributionPaused;
    mapping(uint256 => mapping(address => uint256)) public userIndices;

    /// @dev This empty reserved space is put in place to allow future versions to add new
    /// variables without shifting down storage in the inheritance chain.
    /// See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
    // slither-disable-next-line unused-state
    uint256[44] private __gap;

    event DistributionIndexUpdated(uint256 indexed id, uint256 index);
    event UserIndexUpdated(
        address indexed user,
        uint256 indexed id,
        uint256 index
    );

    // solhint-disable-next-line func-name-mixedcase
    function __init_DistributionManager_(
        address emissionManager_,
        uint128 emissionPerSecond_,
        uint256 distributionDuration_
    ) internal onlyInitializing {
        _grantRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _grantRole(EMISSION_ADMIN_ROLE, emissionManager_);

        emissionPerSecond = emissionPerSecond_;
        distributionEnd = block.timestamp + distributionDuration_;
    }

    modifier onlyEmissionAdmin() {
        if (!hasRole(EMISSION_ADMIN_ROLE, _msgSender())) {
            revert NotEmissionAdmin();
        }
        _;
    }

    function setEmissionPerSecond(
        uint128 emissionPerSecond_
    ) external onlyEmissionAdmin {
        emissionPerSecond = emissionPerSecond_;
    }

    function _setDistributionPaused(uint256 id, bool isPaused) internal {
        distributionPaused[id] = isPaused;
        distributionTimestamps[id] = uint128(block.timestamp);
    }

    /// @dev Updates the state of one distribution, mainly rewards index and timestamp
    /// @param id This used as key in the distribution
    /// @param totalStaked Current total of staked assets for this distribution
    /// @return The new distribution index
    function _updateDistribution(
        uint256 id,
        uint256 totalStaked
    ) internal returns (uint256) {
        uint256 oldIndex = distributionIndices[id];
        uint128 lastUpdateTimestamp = distributionTimestamps[id];

        // slither-disable-next-line incorrect-equality
        if (block.timestamp == lastUpdateTimestamp) {
            return oldIndex;
        }

        uint256 newIndex = _getDistributionIndex(
            oldIndex,
            lastUpdateTimestamp,
            totalStaked,
            distributionPaused[id]
        );

        if (newIndex != oldIndex) {
            distributionIndices[id] = newIndex;
            emit DistributionIndexUpdated(id, newIndex);
        }

        distributionTimestamps[id] = uint128(block.timestamp);

        return newIndex;
    }

    /// @dev Updates the state of an user in a distribution
    /// @param user The user's address
    /// @param id The id of the reference asset of the distribution
    /// @param stakedByUser Amount of tokens staked by the user in the distribution at the moment
    /// @param totalStaked Total tokens staked in the distribution
    /// @return The accrued rewards for the user until the moment
    function _updateUser(
        address user,
        uint256 id,
        uint256 stakedByUser,
        uint256 totalStaked
    ) internal returns (uint256) {
        uint256 newIndex = _updateDistribution(id, totalStaked);
        uint256 userIndex = userIndices[id][user];

        uint256 accruedRewards = 0;

        if (userIndex != newIndex) {
            if (stakedByUser != 0) {
                accruedRewards = _getAccruedRewards(
                    stakedByUser,
                    newIndex,
                    userIndex
                );
            }

            userIndices[id][user] = newIndex;
            emit UserIndexUpdated(user, id, newIndex);
        }

        return accruedRewards;
    }

    /// @dev Internal function for the calculation of user's rewards on a distribution
    /// @param stakedByUser Amount staked by the user on a distribution
    /// @param distributionIndex Current index of the distribution
    /// @param userIndex Index stored for the user, representation his staking moment
    /// @return The rewards
    function _getAccruedRewards(
        uint256 stakedByUser,
        uint256 distributionIndex,
        uint256 userIndex
    ) internal pure returns (uint256) {
        uint256 indexDelta = (distributionIndex - userIndex);
        return (stakedByUser * indexDelta) / PRECISION_FACTOR;
    }

    /// @dev Calculates the next value of an specific distribution index, with validations
    /// @param currentIndex Current index of the distribution
    /// @param lastUpdateTimestamp Last moment this distribution was updated
    /// @param totalSupply of tokens considered for the distribution
    /// @param paused boolean for paused state
    /// @return The new index.
    function _getDistributionIndex(
        uint256 currentIndex,
        uint128 lastUpdateTimestamp,
        uint256 totalSupply,
        bool paused
    ) internal view returns (uint256) {
        if (
            // slither-disable-next-line incorrect-equality
            emissionPerSecond == 0 ||
            totalSupply == 0 ||
            lastUpdateTimestamp == block.timestamp ||
            lastUpdateTimestamp >= distributionEnd ||
            paused
        ) {
            return currentIndex;
        }

        uint256 currentTimestamp = block.timestamp > distributionEnd
            ? distributionEnd
            : block.timestamp;

        uint256 timeDelta = currentTimestamp - lastUpdateTimestamp;

        uint256 newIndex = (emissionPerSecond * timeDelta * PRECISION_FACTOR) /
            totalSupply;

        return newIndex + currentIndex;
    }
}
