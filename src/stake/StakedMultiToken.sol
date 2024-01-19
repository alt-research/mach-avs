// SPDX-License-Identifier: agpl-3.0
pragma solidity ^0.8.23;

import {EnumerableSetUpgradeable} from "@openzeppelin-upgrades/contracts/utils/structs/EnumerableSetUpgradeable.sol";
import {SafeERC20, IERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ERC1155SupplyUpgradeable, ERC1155Upgradeable} from "@openzeppelin-upgrades/contracts/token/ERC1155/extensions/ERC1155SupplyUpgradeable.sol";
import {ISlasher as IAccelSlasher} from "accel/contracts/src/L1/interfaces/ISlasher.sol";
import {DistributionManager, AccessControlUpgradeable} from "./DistributionManager.sol";

/// @dev The staked token should be deployed on Ethereum.
/// This is adapted from https://github.com/bgd-labs/aave-stk-v1-5/blob/8867dd5b1137d4d46acd9716fe98759cb16b1606/src/contracts/StakedTokenV3.sol
// solhint-disable not-rely-on-time, var-name-mixedcase
// slither-disable-start timestamp
contract StakedMultiToken is
    IAccelSlasher,
    ERC1155SupplyUpgradeable,
    DistributionManager
{
    using EnumerableSetUpgradeable for EnumerableSetUpgradeable.UintSet;
    using SafeERC20 for IERC20;

    error ZeroAddress();
    error NotSlashingAdmin();
    error NotSettleSlashingAdmin();
    error NotCooldownAdmin();
    error NotUnstakeWindowAdmin();
    error InvalidID();
    error InvalidStartIndex();
    error AlreadyInitialized();
    error InvalidBPS();
    error ZeroExchangeRate();
    error ZeroAmount();
    error Not216Bits();
    error MultiplicationOverflow();
    error SlashingOngoing();
    error MaxSupply();
    error InvalidBalanceOnCooldown();
    error InsufficientCooldown();
    error UnstakeWindowFinished();
    error InvalidZeroMaxRedeemable();
    error AmountLessThanMinimum();
    error ShareLessThanMinimum();
    error RemainingLessThanMinimum();

    bytes32 public constant SLASHING_ADMIN_ROLE =
        keccak256("SLASHING_ADMIN_ROLE");
    bytes32 public constant SETTLE_SLASHING_ADMIN_ROLE =
        keccak256("SETTLE_SLASHING_ADMIN_ROLE");
    bytes32 public constant COOLDOWN_ADMIN_ROLE =
        keccak256("COOLDOWN_ADMIN_ROLE");
    bytes32 public constant UNSTAKE_WINDOW_ADMIN_ROLE =
        keccak256("UNSTAKE_WINDOW_ADMIN_ROLE");

    /// @dev MAX_BPS the maximum number of basis points.
    /// 10000 basis points are equivalent to 100%.
    uint256 public constant MAX_BPS = 1e4;

    /// @dev The exchange rate is initialized as 1e18 which reflects a 1:1 peg.
    /// It is only adjusted up & down based on slash and returnFunds actions.
    uint216 public constant INITIAL_EXCHANGE_RATE = 1e18;
    uint256 public constant EXCHANGE_RATE_UNIT = 1e18;

    uint256 public constant MAX_SUPPLY = (type(uint256).max) / 1e18;

    /// @dev 1 unit of the staked asset (1 ether for 18 decimal assets).
    /// This lower bound is to prevent spam & avoid exchangeRate issues
    /// as returnFunds can be called permissionless an attacker could spam returnFunds(1) to produce exchangeRate snapshots making voting expensive
    uint256 public constant LOWER_BOUND = 1e18;

    IERC20 public stakedToken;
    IERC20 public rewardToken;
    /// @dev Seconds available to redeem once the cooldown period is fullfilled
    uint256 public unstakeWindow;
    /// @dev Address to pull from the rewards, needs to have approved this contract
    address public rewardVault;

    /// @dev Seconds between starting cooldown and being able to withdraw
    uint256 public cooldownSeconds;
    /// @dev The maximum amount of funds that can be slashed at any given time. It defaults to 3000 (30%).
    uint256 public maxSlashableBPS;
    /// @dev Mirror of latest snapshot value for cheaper access
    mapping(uint256 => uint216) public exchangeRates;
    mapping(uint256 => mapping(address => uint256)) public rewardsBalances;
    EnumerableSetUpgradeable.UintSet private _ids;

    struct CooldownSnapshot {
        uint40 timestamp;
        uint216 amount;
    }
    mapping(uint256 => mapping(address => CooldownSnapshot)) public cooldowns;

    event InitializeAsset(uint256 id);
    event Stake(
        address indexed from,
        address indexed onBehalfOf,
        uint256 id,
        uint256 assets,
        uint256 shares
    );
    event Redeem(
        address indexed from,
        address indexed to,
        uint256 id,
        uint256 assets,
        uint256 shares
    );

    event RewardsAccrued(address user, uint256 id, uint256 amount);
    event RewardsClaimed(
        address indexed from,
        address indexed to,
        uint256 id,
        uint256 amount
    );

    event Cooldown(address indexed user, uint256 indexed id, uint256 amount);

    event SetMaxSlashableBPS(uint256 bps);
    event SetCooldownSeconds(uint256 cooldownSeconds);
    event SetUnstakeWindow(uint256 unstakeWindow);
    event SetExchangeRate(uint256 id, uint216 exchangeRate);
    event Slash(uint256 id, address indexed beneficiary, uint256 amount);
    event ReturnFunds(uint256 id, uint256 amount);
    event SettleSlashing(uint256 id);

    modifier onlySlashingAdmin() {
        if (!hasRole(SLASHING_ADMIN_ROLE, _msgSender())) {
            revert NotSlashingAdmin();
        }
        _;
    }
    modifier onlySettleSlashingAdmin() {
        if (!hasRole(SETTLE_SLASHING_ADMIN_ROLE, _msgSender())) {
            revert NotSettleSlashingAdmin();
        }
        _;
    }

    modifier onlyCooldownAdmin() {
        if (!hasRole(COOLDOWN_ADMIN_ROLE, _msgSender())) {
            revert NotCooldownAdmin();
        }
        _;
    }
    modifier onlyUnstakeWindowAdmin() {
        if (!hasRole(UNSTAKE_WINDOW_ADMIN_ROLE, _msgSender())) {
            revert NotUnstakeWindowAdmin();
        }
        _;
    }

    modifier onlyValidID(uint256 id) {
        if (!isValidID(id)) {
            revert InvalidID();
        }
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address emissionManager_,
        uint128 emissionPerSecond_,
        uint128 distributionDuration_,
        address stakedToken_,
        address rewardToken_,
        uint256 cooldownSeconds_,
        uint256 unstakeWindow_,
        address rewardsVault_,
        address slashingAdmin_,
        address settleSlashingAdmin_,
        address cooldownAdmin_,
        address unstakeWindowAdmin_
    ) external initializer {
        if (
            emissionManager_ == address(0) ||
            stakedToken_ == address(0) ||
            rewardToken_ == address(0) ||
            rewardsVault_ == address(0) ||
            slashingAdmin_ == address(0) ||
            settleSlashingAdmin_ == address(0) ||
            cooldownAdmin_ == address(0) ||
            unstakeWindowAdmin_ == address(0)
        ) {
            revert ZeroAddress();
        }

        __ERC1155Supply_init();
        __init_DistributionManager_(
            emissionManager_,
            emissionPerSecond_,
            distributionDuration_
        );

        stakedToken = IERC20(stakedToken_);
        rewardToken = IERC20(rewardToken_);
        cooldownSeconds = cooldownSeconds_;
        unstakeWindow = unstakeWindow_;
        rewardVault = rewardsVault_;

        _grantRole(SLASHING_ADMIN_ROLE, slashingAdmin_);
        _grantRole(SETTLE_SLASHING_ADMIN_ROLE, settleSlashingAdmin_);
        _grantRole(COOLDOWN_ADMIN_ROLE, cooldownAdmin_);
        _grantRole(UNSTAKE_WINDOW_ADMIN_ROLE, unstakeWindowAdmin_);

        // Set initial value in initializer
        maxSlashableBPS = 3e3;
    }

    /// @notice Initializes asset
    /// @dev ACCEL USE CASE:
    /// Each node has a unique address which can be converted to a uint256 i.e uint256(uint160(addr)).
    /// For example, a node address like 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
    /// corresponds to the ID 0x000000000000000000000000f39Fd6e51aad88F6F4ce6aB8827279cffFb92266.
    function initializeAsset(uint256 id) external onlyEmissionAdmin {
        if (exchangeRates[id] != 0) {
            revert AlreadyInitialized();
        }

        // slither-disable-next-line unused-return
        _ids.add(id);
        _updateExchangeRate(id, INITIAL_EXCHANGE_RATE);
        _updateDistribution(id, totalSupply(id));
        emit InitializeAsset(id);
    }

    /// @notice Returns the total number of IDs.
    /// @return The total number of IDs.
    function totalIDs() public view returns (uint256) {
        return _ids.length();
    }

    /// @notice Checks if the given ID is valid and exists in the contract.
    /// @param id The ID to check.
    /// @return True if the ID is valid and exists, false otherwise.
    function isValidID(uint256 id) public view returns (bool) {
        return _ids.contains(id);
    }

    /// @notice Returns an array of IDs starting from the specified index up to the query size.
    /// @param start The start index to retrieve IDs.
    /// @param querySize The number of IDs to retrieve.
    /// @return An array of IDs.
    function queryIDs(
        uint256 start,
        uint256 querySize
    ) external view returns (uint256[] memory) {
        uint256 length = totalIDs();

        if (start >= length) {
            revert InvalidStartIndex();
        }

        uint256 end = start + querySize;

        if (end > length) {
            end = length;
        }

        uint256[] memory output = new uint256[](end - start);

        for (uint256 i = start; i < end; ) {
            output[i - start] = _ids.at(i);

            unchecked {
                ++i;
            }
        }

        return output;
    }

    /// @notice Sets the max slashable BPS
    /// @param bps must be strictly greater than 0 and less than 10000.
    /// Otherwise the exchange rate calculation would result in 0 division.
    function setMaxSlashableBPS(uint256 bps) external onlySlashingAdmin {
        if (bps == 0 || bps >= MAX_BPS) {
            revert InvalidBPS();
        }

        maxSlashableBPS = bps;
        emit SetMaxSlashableBPS(bps);
    }

    /// @notice Sets the cooldown seconds
    /// @param cooldownSeconds_ the new amount of cooldown seconds
    function setCooldownSeconds(
        uint256 cooldownSeconds_
    ) external onlyCooldownAdmin {
        cooldownSeconds = cooldownSeconds_;
        emit SetCooldownSeconds(cooldownSeconds_);
    }

    /// @notice Sets the unstake window seconds
    /// @param unstakeWindow_ the new amount of unstake window seconds
    function setUnstakeWindow(
        uint256 unstakeWindow_
    ) external onlyUnstakeWindowAdmin {
        unstakeWindow = unstakeWindow_;
        emit SetUnstakeWindow(unstakeWindow_);
    }

    /// @notice Updates the exchangeRate and emits events accordingly
    /// @param id Token id
    /// @param newExchangeRate the new exchange rate
    function _updateExchangeRate(
        uint256 id,
        uint216 newExchangeRate
    ) internal virtual {
        if (newExchangeRate == 0) {
            revert ZeroExchangeRate();
        }

        exchangeRates[id] = newExchangeRate;
        emit SetExchangeRate(id, newExchangeRate);
    }

    /// @notice Calculates the exchange rate based on totalAssets and totalShares
    /// @dev always rounds up to ensure 100% backing of shares by rounding in favor of the contract
    /// @param totalAssets The total amount of assets staked
    /// @param totalShares The total amount of shares
    /// @return exchangeRate as 18 decimal precision uint216
    function _getExchangeRate(
        uint256 totalAssets,
        uint256 totalShares
    ) internal pure returns (uint216) {
        uint256 value = ((totalShares * EXCHANGE_RATE_UNIT) + totalAssets - 1) /
            totalAssets;

        if (value > type(uint216).max) {
            revert Not216Bits();
        }

        return uint216(value);
    }

    /// @notice Returns the exact amount of shares that would be received for the provided number of assets
    /// @param id token ID
    /// @param assets the number of assets to stake
    /// @return uint256 shares the number of shares that would be received
    function previewStake(
        uint256 id,
        uint256 assets
    ) public view returns (uint256) {
        if (assets > (type(uint256).max) / exchangeRates[id]) {
            revert MultiplicationOverflow();
        }

        return (assets * exchangeRates[id]) / EXCHANGE_RATE_UNIT;
    }

    /// @dev returns the exact amount of assets that would be redeemed for the provided number of shares
    /// @param id token ID
    /// @param shares the number of shares to redeem
    /// @return uint256 assets the number of assets that would be redeemed
    function previewRedeem(
        uint256 id,
        uint256 shares
    ) public view returns (uint256) {
        return (EXCHANGE_RATE_UNIT * shares) / exchangeRates[id];
    }

    /// @dev Returns the accrued rewards for a staker on a specific ID.
    /// @param staker The address of the staker.
    /// @param id The ID for which to check the accrued rewards.
    /// @return The accrued rewards for the staker and ID.
    function getAccruedRewards(
        address staker,
        uint256 id
    ) external view returns (uint256) {
        uint256 ditributionIndex = _getDistributionIndex(
            distributionIndices[id],
            distributionTimestamps[id],
            totalSupply(id),
            distributionPaused[id]
        );
        uint256 userIndex = userIndices[id][staker];

        uint256 accruedRewards = _getAccruedRewards(
            balanceOf(staker, id),
            ditributionIndex,
            userIndex
        );

        return accruedRewards;
    }

    /// @notice Stakes tokens.
    /// @param to The address of the user for whom the tokens are being staked
    /// @param id The identifier of the staking pool
    /// @param amount The amount of tokens to be staked
    function stake(
        address to,
        uint256 id,
        uint256 amount
    ) public virtual onlyValidID(id) {
        address from = _msgSender();

        if (distributionPaused[id]) {
            revert SlashingOngoing();
        }

        if (amount == 0) {
            revert ZeroAmount();
        }

        uint256 accruedRewards = _updateUser(
            to,
            id,
            balanceOf(to, id),
            totalSupply(id)
        );

        if (accruedRewards != 0) {
            rewardsBalances[id][to] += accruedRewards;
            emit RewardsAccrued(to, id, accruedRewards);
        }

        uint256 sharesToMint = previewStake(id, amount);
        // This is to prevent arithmetic overflow in the `previewRedeem` function.
        // Without this, DoS attack is possible.
        if (totalSupply(id) + sharesToMint > MAX_SUPPLY) {
            revert MaxSupply();
        }

        _mint(to, id, sharesToMint, "");
        emit Stake(from, to, id, amount, sharesToMint);
        stakedToken.safeTransferFrom(from, address(this), amount);
    }

    /// @notice Initiates the cooldown period for a user's staked tokens
    /// @param id The identifier of the staking pool
    function cooldown(uint256 id) external onlyValidID(id) {
        address from = _msgSender();
        uint256 amount = balanceOf(from, id);

        if (amount == 0) {
            revert InvalidBalanceOnCooldown();
        }

        cooldowns[id][from] = CooldownSnapshot({
            timestamp: uint40(block.timestamp),
            amount: uint216(amount)
        });

        emit Cooldown(from, id, amount);
    }

    /// @notice Redeems staked tokens, and stop earning rewards
    /// @param to Address to redeem to
    /// @param id Token ID
    /// @param amount Amount to redeem
    function redeem(
        address to,
        uint256 id,
        uint256 amount
    ) external onlyValidID(id) {
        address from = _msgSender();

        if (amount == 0) {
            revert ZeroAmount();
        }

        CooldownSnapshot memory cooldownSnapshot = cooldowns[id][from];
        if (!distributionPaused[id]) {
            if (
                block.timestamp <= cooldownSnapshot.timestamp + cooldownSeconds
            ) {
                revert InsufficientCooldown();
            }

            if (
                block.timestamp -
                    (cooldownSnapshot.timestamp + cooldownSeconds) >
                unstakeWindow
            ) {
                revert UnstakeWindowFinished();
            }
        }

        uint256 balanceOfFrom = balanceOf(from, id);
        uint256 maxRedeemable = distributionPaused[id]
            ? balanceOfFrom
            : cooldownSnapshot.amount;

        if (maxRedeemable == 0) {
            revert InvalidZeroMaxRedeemable();
        }

        uint256 amountToRedeem = (amount > maxRedeemable)
            ? maxRedeemable
            : amount;

        _updateCurrentUnclaimedRewards(from, id, balanceOfFrom, true);

        uint256 underlyingToRedeem = previewRedeem(id, amountToRedeem);

        _burn(from, id, amountToRedeem);

        if (cooldownSnapshot.timestamp != 0) {
            if (cooldownSnapshot.amount - amountToRedeem == 0) {
                delete cooldowns[id][from];
            } else {
                cooldowns[id][from].amount -= uint216(amountToRedeem);
            }
        }
        emit Redeem(from, to, id, underlyingToRedeem, amountToRedeem);

        IERC20(stakedToken).safeTransfer(to, underlyingToRedeem);
    }

    /// @notice Executes a slashing of the underlying of a certain amount, transferring the seized funds
    /// to beneficiary. Decreasing the amount of underlying will automatically adjust the exchange rate.
    /// A call to `slash` will start a slashing event which has to be settled via `settleSlashing`.
    /// As long as the slashing event is ongoing, stake and slash are deactivated.
    /// - MUST NOT be called when a previous slashing is still ongoing
    /// @param slashedUser the slashed user (the token id)
    /// @param beneficiary the address where seized funds will be transferred
    /// @param bps the basis points of the underlying to be slashed
    /// - if the amount bigger than maximum allowed, the maximum will be slashed instead.
    /// @return amount the amount slashed
    function slash(
        address slashedUser,
        address beneficiary,
        uint256 bps
    )
        external
        onlySlashingAdmin
        onlyValidID(uint256(uint160(slashedUser)))
        returns (uint256)
    {
        uint256 id = uint256(uint160(slashedUser));
        if (distributionPaused[id]) {
            revert SlashingOngoing();
        }

        if (bps == 0) {
            revert InvalidBPS();
        }
        if (bps > maxSlashableBPS) {
            bps = maxSlashableBPS;
        }

        uint256 currentShares = totalSupply(id);
        uint256 balance = previewRedeem(id, currentShares);

        if (balance > (type(uint256).max) / maxSlashableBPS) {
            revert MultiplicationOverflow();
        }

        /// bps is greater than 0 and less than 1e4
        uint256 amount = (balance * bps) / MAX_BPS;

        if (balance - amount < LOWER_BOUND) {
            revert RemainingLessThanMinimum();
        }

        _setDistributionPaused(id, true);

        _updateExchangeRate(
            id,
            _getExchangeRate(balance - amount, currentShares)
        );
        emit Slash(id, beneficiary, amount);

        stakedToken.safeTransfer(beneficiary, amount);

        return amount;
    }

    /// @notice Pulls STAKE_TOKEN and distributes them amongst current stakers by altering the exchange rate.
    /// This method is permissionless and intended to be used after a slashing event to return potential excess funds.
    /// @param amount the token id
    /// @param amount amount of STAKE_TOKEN to pull.
    function returnFunds(uint256 id, uint256 amount) external onlyValidID(id) {
        if (amount < LOWER_BOUND) {
            revert AmountLessThanMinimum();
        }

        uint256 currentShares = totalSupply(id);

        if (currentShares < LOWER_BOUND) {
            revert ShareLessThanMinimum();
        }

        uint256 assets = previewRedeem(id, currentShares);
        _updateExchangeRate(
            id,
            _getExchangeRate(assets + amount, currentShares)
        );
        emit ReturnFunds(id, amount);

        stakedToken.safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @notice Settles an ongoing slashing event
    /// @param id the token id
    function settleSlashing(
        uint256 id
    ) external onlySettleSlashingAdmin onlyValidID(id) {
        _setDistributionPaused(id, false);
        emit SettleSlashing(id);
    }

    /// @notice Claims an `amount` of `rewardToken` to the address `to`
    /// @param to Address to stake for
    /// @param id Token ID
    /// @param amount Amount to stake
    function claimRewards(
        address to,
        uint256 id,
        uint256 amount
    ) external onlyValidID(id) {
        address from = _msgSender();
        uint256 newTotalRewards = _updateCurrentUnclaimedRewards(
            from,
            id,
            balanceOf(from, id),
            false
        );

        uint256 amountToClaim = (amount > newTotalRewards)
            ? newTotalRewards
            : amount;

        // slither-disable-next-line incorrect-equality
        if (amountToClaim == 0) {
            revert ZeroAmount();
        }

        rewardsBalances[id][from] = newTotalRewards - amountToClaim;
        emit RewardsClaimed(from, to, id, amountToClaim);

        // slither-disable-next-line arbitrary-send-erc20
        rewardToken.safeTransferFrom(rewardVault, to, amountToClaim);
    }

    /// @dev Internal ERC1155 _safeTransferFrom of the tokenized staked tokens
    /// @param from Address to transfer from
    /// @param to Address to transfer to
    /// @param id token ID
    /// @param amount Amount to transfer
    function _safeTransferFrom(
        address from,
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) internal override {
        uint256 balanceOfFrom = balanceOf(from, id);
        // Sender
        _updateCurrentUnclaimedRewards(from, id, balanceOfFrom, true);

        // Recipient
        if (from != to) {
            uint256 balanceOfTo = balanceOf(to, id);
            _updateCurrentUnclaimedRewards(to, id, balanceOfTo, true);

            CooldownSnapshot memory prevCooldown = cooldowns[id][from];
            if (prevCooldown.timestamp != 0) {
                // if cooldown was set and whole balance of sender was transferred - clear cooldown
                if (balanceOfFrom == amount) {
                    delete cooldowns[id][from];
                } else if (balanceOfFrom - amount < prevCooldown.amount) {
                    cooldowns[id][from].amount = uint216(
                        balanceOfFrom - amount
                    );
                }
            }
        }

        super._safeTransferFrom(from, to, id, amount, data);
    }

    /// @dev Updates the user state related with his accrued rewards
    /// @param user Address of the user
    /// @param id Token ID
    /// @param userBalance The current balance of the user
    /// @param updateStorage Boolean flag used to update or not the rewardsBalances of the user
    /// @return The unclaimed rewards that were added to the total accrued
    function _updateCurrentUnclaimedRewards(
        address user,
        uint256 id,
        uint256 userBalance,
        bool updateStorage
    ) internal returns (uint256) {
        uint256 accruedRewards = _updateUser(
            user,
            id,
            userBalance,
            totalSupply(id)
        );
        uint256 unclaimedRewards = rewardsBalances[id][user] + accruedRewards;

        if (accruedRewards != 0) {
            if (updateStorage) {
                rewardsBalances[id][user] = unclaimedRewards;
            }
            emit RewardsAccrued(user, id, accruedRewards);
        }

        return unclaimedRewards;
    }

    function supportsInterface(
        bytes4 interfaceId
    )
        public
        view
        virtual
        override(AccessControlUpgradeable, ERC1155Upgradeable)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
